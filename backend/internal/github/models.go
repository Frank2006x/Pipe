package github

import (
	"time"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type User struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

type RepositoryOwner struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`

	Owner RepositoryOwner `json:"owner"`

	Description   string `json:"description"`
	DefaultBranch string `json:"default_branch"`

	Private bool `json:"private"`

	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`

	Language   string `json:"language"`
	Visibility string `json:"visibility"`
	Archived   bool   `json:"archived"`
	Disabled   bool   `json:"disabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PushedAt  time.Time `json:"pushed_at"`
}

type WebhookConfig struct {
    URL         string `json:"url"`
    ContentType string `json:"content_type"`
    Secret      string `json:"secret"`
    InsecureSSL string `json:"insecure_ssl"`
}

type CreateWebhookRequest struct {
    Name   string         `json:"name"`
    Active bool           `json:"active"`
    Events []string       `json:"events"`
    Config WebhookConfig  `json:"config"`
}

type Webhook struct {
    ID int64 `json:"id"`
}
