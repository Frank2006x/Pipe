package github

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
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`

	Owner         RepositoryOwner `json:"owner"`

	Description   string `json:"description"`
	DefaultBranch string `json:"default_branch"`

	Private       bool   `json:"private"`

	HTMLURL       string `json:"html_url"`
	CloneURL      string `json:"clone_url"`

	Language       string `json:"language"`
	Visibility     string `json:"visibility"`
	Archived       bool   `json:"archived"`
	Disabled       bool   `json:"disabled"`

	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	PushedAt       string `json:"pushed_at"`
}