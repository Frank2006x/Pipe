package service

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/github"
	"context"
	"errors"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidCode  = errors.New("invalid code")
	ErrInvalidToken = errors.New("invalid token")
)

type AuthService struct {
	githubClient *github.Client
	jwtMaker     *auth.JwtMaker
	queries      *sqlc.Queries
}

func NewAuthService(githubClient *github.Client, jwtMaker *auth.JwtMaker, queries *sqlc.Queries) *AuthService {
	return &AuthService{
		githubClient: githubClient,
		jwtMaker:     jwtMaker,
		queries:      queries,
	}
}

func (s *AuthService) upsertUser(
	ctx context.Context,
	githubUser *github.User,
) (sqlc.User, error) {

	_, err := s.queries.GetUserByGithubID(ctx, githubUser.ID)
	if err == nil {
		log.Infof("User with GitHub ID %d already exists, updating user", githubUser.ID)
		updatedUser, err := s.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
			GithubID: githubUser.ID,
			Username: githubUser.Login,
			Email:    githubUser.Email,
			AvatarUrl: pgtype.Text{
				String: githubUser.AvatarURL,
				Valid:  githubUser.AvatarURL != "",
			},
		})

		if err != nil {
			return sqlc.User{}, err
		}

		return updatedUser, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return sqlc.User{}, err
	}
	log.Infof("User with GitHub ID %d does not exist, creating new user", githubUser.ID)
	newUser, err := s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		GithubID: githubUser.ID,
		Username: githubUser.Login,
		Email:    githubUser.Email,
		AvatarUrl: pgtype.Text{
			String: githubUser.AvatarURL,
			Valid:  githubUser.AvatarURL != "",
		},
	})

	if err != nil {
		return sqlc.User{}, err
	}

	return newUser, nil
}

func (s *AuthService) upsertToken(
	ctx context.Context,
	userID int64,
	accessToken string,
) error {

	_, err := s.queries.GetGithubToken(ctx, userID)
	if err == nil {

		_, err = s.queries.UpdateGithubToken(ctx, sqlc.UpdateGithubTokenParams{
			UserID:      userID,
			AccessToken: accessToken,
		})

		return err
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	_, err = s.queries.CreateGithubToken(ctx, sqlc.CreateGithubTokenParams{
		UserID:      userID,
		AccessToken: accessToken,
	})

	return err
}

func (s *AuthService) Callback(ctx context.Context, code string) (string, error) {

	token, err := s.githubClient.ExchangeCode(ctx, code)
	if err != nil {
		return "", err
	}

	githubUser, err := s.githubClient.GetUser(ctx, token.AccessToken)
	if err != nil {
		return "", err
	}

	dbUser, err := s.upsertUser(ctx, githubUser)
	if err != nil {
		return "", err
	}

	err = s.upsertToken(ctx, dbUser.ID, token.AccessToken)
	if err != nil {
		return "", err
	}

	jwtToken, err := s.jwtMaker.GenerateJWT(dbUser.ID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
