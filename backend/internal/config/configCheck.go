package config

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}

func CheckConfig(config Config) error {
	if config.POSTGRES_DB == "" {
		return &ConfigError{Field: "POSTGRES_DB", Message: "POSTGRES_DB is required"}
	}
	if config.GITHUB_CLIENT_ID == "" {
		return &ConfigError{Field: "GITHUB_CLIENT_ID", Message: "GITHUB_CLIENT_ID is required"}
	}
	if config.GITHUB_CLIENT_SECRET == "" {
		return &ConfigError{Field: "GITHUB_CLIENT_SECRET", Message: "GITHUB_CLIENT_SECRET is required"}
	}
	if config.GITHUB_CALLBACK_URL == "" {
		return &ConfigError{Field: "GITHUB_CALLBACK_URL", Message: "GITHUB_CALLBACK_URL is required"}
	}
	if config.JWT_SECRET == "" {
		return &ConfigError{Field: "JWT_SECRET", Message: "JWT_SECRET is required"}
	}
	return nil
}
