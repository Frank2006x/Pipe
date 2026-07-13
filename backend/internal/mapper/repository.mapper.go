package mapper

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/github"

	"github.com/jackc/pgx/v5/pgtype"
)

func MapGitHubRepository(userID int64, repository *github.Repository) sqlc.CreateRepositoryParams {
	description := pgtype.Text{
		String: repository.Description,
		Valid:  repository.Description != "",
	}

	return sqlc.CreateRepositoryParams{
		UserID:        userID,
		GithubRepoID:  repository.ID,
		Name:          repository.Name,
		FullName:      repository.FullName,
		HtmlUrl:       repository.HTMLURL,
		DefaultBranch: repository.DefaultBranch,
		Private:       repository.Private,
		WebhookID:     pgtype.Int8{},
		IsActive:      true,
		Owner:         repository.Owner.Login,
		Description:   description,
		CloneUrl:      repository.CloneURL,
		WebhookSecret: pgtype.Text{},
	}
}
