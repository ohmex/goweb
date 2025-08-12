package config

import "os"

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
	Social        SocialConfig
}

type SocialConfig struct {
	Google GoogleConfig
	GitHub GitHubConfig
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func LoadAuthConfig() AuthConfig {
	return AuthConfig{
		AccessSecret:  os.Getenv("ACCESS_SECRET"),
		RefreshSecret: os.Getenv("REFRESH_SECRET"),
		Social: SocialConfig{
			Google: GoogleConfig{
				ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			},
			GitHub: GitHubConfig{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			},
		},
	}
}
