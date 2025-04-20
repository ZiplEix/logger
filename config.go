package logger

import "github.com/golang-jwt/jwt/v5"

// Config define the configuration for the logger middleware
//
// It is recommended to use the default configuration and override only the values you need
type Config struct {
	// DatabasePath define the location where the sqlite database wil be stored
	//
	// Optional. Default: ./logs.sqlite
	DatabasePath string

	// WorkerBufferSize define the size of the worker buffer
	//
	// Optional. Default: 100
	WorkerBufferSize int

	// Path define the the base path to the logs view
	// (it is recommended to use a path that is not publicly accessible)
	// (e.g. /logs or /admin/logs)
	//
	// If set manually, it is recommended to also set the ExcludeRoutes and ExcludePatterns
	//
	// Optional. Default: /logs
	Path string

	// ExcludePaths define the paths to exclude from logging
	//
	// Optional. Default: ["logs", "/favicon.ico"]
	ExcludeRoutes []string

	// ExcludePaths define the paths to exclude from logging
	//
	// Optional. Default: ["/logs/*"]
	ExcludePatterns []string

	// ExcludePAram define a parameter that will exclude the request from logging
	//
	// Optional. Default: ""
	ExcludeParam string

	// LogDetailMember define the member name that will be used to store the log details in the fiber.Ctx local
	//
	// Optional. Default: "logDetail"
	LogDetailMember string

	// SecureByPassword define if the logs view should be secured by a password
	//
	// Optional. Default: true
	SecureByPassword *bool

	// Password define the password to secure the logs view
	// (it is recommended to store it in an environment variable)
	//
	// Optional. Default: "password"
	Password string

	// JwtSecret define the secret to sign the JWT token
	// (it is recommended to store it in an environment variable)
	//
	// Optional. Default: "secret"
	JwtSecret string

	// JwtExpireTime define the expiration time of the JWT token in seconds
	//
	// Optional. Default: 3600
	// (1 hour)
	JwtExpireTime int64

	// Theme define the theme of the logs view
	//
	// Optional. Default: ThemeConfig.Dark
	// (see https://daisyui.com/theme-generator/)
	Theme string

	// JwtSigningMethod define the signing method to use for the JWT token
	//
	// Optional. Default: jwt.SigningMethodHS256
	JwtSigningMethod jwt.SigningMethod

	// AuthTokenCookieName define the name of the cookie that will be used to store the JWT token
	//
	// Optional. Default: "auth_token"
	AuthTokenCookieName string

	// IncludeLogPageConnexion define if the connexion via password to the logs page should be logged
	//
	// Optional. Default: true
	IncludeLogPageConnexion *bool
}

var configDefault = Config{
	DatabasePath:            "./logs.sqlite",
	WorkerBufferSize:        100,
	Path:                    "/logs",
	ExcludeRoutes:           []string{"/logs", "/favicon.ico"},
	ExcludePatterns:         []string{"/logs/*"},
	ExcludeParam:            "",
	LogDetailMember:         "logDetail",
	SecureByPassword:        ptr(true),
	Password:                "password",
	JwtSecret:               "secret",
	JwtExpireTime:           3600,
	Theme:                   Theme.Dark,
	JwtSigningMethod:        jwt.SigningMethodHS256,
	AuthTokenCookieName:     "auth_token",
	IncludeLogPageConnexion: ptr(true),
}

func defaultConfig(config Config) Config {
	cfg := config

	if config.DatabasePath == "" {
		cfg.DatabasePath = configDefault.DatabasePath
	}

	if config.WorkerBufferSize == 0 {
		cfg.WorkerBufferSize = configDefault.WorkerBufferSize
	}

	if config.Path == "" {
		cfg.Path = configDefault.Path
	}

	if config.ExcludeRoutes == nil {
		cfg.ExcludeRoutes = configDefault.ExcludeRoutes
	}

	if config.ExcludePatterns == nil {
		cfg.ExcludePatterns = configDefault.ExcludePatterns
	}

	if config.ExcludeParam == "" {
		cfg.ExcludeParam = configDefault.ExcludeParam
	}

	if config.LogDetailMember == "" {
		cfg.LogDetailMember = configDefault.LogDetailMember
	}

	if config.SecureByPassword == nil {
		cfg.SecureByPassword = configDefault.SecureByPassword
	}

	if config.Password == "" {
		cfg.Password = configDefault.Password
	}

	if config.JwtSecret == "" {
		cfg.JwtSecret = configDefault.JwtSecret
	}

	if config.JwtExpireTime == 0 || config.JwtExpireTime < 0 {
		cfg.JwtExpireTime = configDefault.JwtExpireTime
	}

	if config.Theme == "" {
		cfg.Theme = configDefault.Theme
	}

	if config.JwtSigningMethod == nil {
		cfg.JwtSigningMethod = configDefault.JwtSigningMethod
	}

	if config.AuthTokenCookieName == "" {
		cfg.AuthTokenCookieName = configDefault.AuthTokenCookieName
	}

	if config.IncludeLogPageConnexion == nil {
		cfg.IncludeLogPageConnexion = configDefault.IncludeLogPageConnexion
	}

	return cfg
}
