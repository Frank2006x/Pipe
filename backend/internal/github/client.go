package github

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL = "https://github.com"
	apiURL  = "https://api.github.com"

	authorizeURL = baseURL + "/login/oauth/authorize"
	tokenURL     = baseURL + "/login/oauth/access_token"

	userEndpoint         = apiURL + "/user"
	repositoriesEndpoint = apiURL + "/user/repos"

	userAgent = "Pipe"
)

type GitHubClient interface {
	GetAuthURL(state string) (string, error)

	ExchangeCode(
		ctx context.Context,
		code string,
	) (*AccessTokenResponse, error)

	GetUser(
		ctx context.Context,
		accessToken string,
	) (*User, error)

	ListRepositories(
		ctx context.Context,
		accessToken string,
	) ([]sqlc.Repository, error)
}

type Client struct {
	httpClient *http.Client
	config     config.Config
}

func NewClient(config config.Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		config: config,
	}
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	endpoint string,
	body io.Reader,
	accessToken string,
) (*http.Request, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		endpoint,
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", userAgent)

	if accessToken != "" {
		req.Header.Set(
			"Authorization",
			fmt.Sprintf("Bearer %s", accessToken),
		)
	}

	return req, nil
}

func (c *Client) do(
	req *http.Request,
	out any,
) error {

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {

		var githubErr struct {
			Message string `json:"message"`
		}

		_ = json.NewDecoder(resp.Body).Decode(&githubErr)

		if githubErr.Message != "" {
			return fmt.Errorf("github: %s", githubErr.Message)
		}

		return fmt.Errorf(
			"github returned status %d",
			resp.StatusCode,
		)
	}

	if out == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) GetAuthURL(state string) (string, error) {
	u, err := url.Parse(authorizeURL)
	if err != nil {
		return "", err
	}

	q := u.Query()

	q.Set("client_id", c.config.GITHUB_CLIENT_ID)
	q.Set("redirect_uri", c.config.GITHUB_CALLBACK_URL)
	q.Set("scope", "repo user:email")
	q.Set("state", state)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *Client) ExchangeCode(
	ctx context.Context,
	code string,
) (*AccessTokenResponse, error) {

	data := url.Values{}
	data.Set("client_id", c.config.GITHUB_CLIENT_ID)
	data.Set("client_secret", c.config.GITHUB_CLIENT_SECRET)
	data.Set("code", code)
	data.Set("redirect_uri", c.config.GITHUB_CALLBACK_URL)

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		tokenURL,
		strings.NewReader(data.Encode()),
		"",
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	req.Header.Set("Accept", "application/json")

	var token AccessTokenResponse

	if err := c.do(req, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *Client) GetUser(
	ctx context.Context,
	token string,
) (*User, error) {

	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		userEndpoint,
		nil,
		token,
	)
	if err != nil {
		return nil, err
	}

	var user User

	if err := c.do(req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) ListRepositories(
	ctx context.Context,
	token string,
) ([]Repository, error) {

	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		repositoriesEndpoint,
		nil,
		token,
	)
	if err != nil {
		return nil, err
	}

	var repos []Repository

	if err := c.do(req, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}
