package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	StringRegex      = `^[a-zA-z0-9 !@#$€%^&*{}|()=/\;:+-_~'"<>,.? \t]*$`
	TitleStringRegex = `^[a-zA-z0-9 !()=[]:;+-_~'",.? \t]*$`
	TextStringRegex  = `^[a-zA-z0-9 !@#$€%^&*{}][|/\()=/\;:+-_~'"<>,.? \t]*$`

	SendingEmailNotification           = "We sent an email with a verification code to "
	UserConfirmationEmailTemplateName  = "verificationCode.html"
	UserConfirmationEmailTemplatePath  = "internal/user/data/repository/external/mail/template"
	ForgottenPasswordEmailTemplateName = "resetPassword.html"
	ForgottenPasswordEmailTemplatePath = "internal/user/data/repository/external/mail/template"

	InternalErrorNotification                  = "something went wrong, please repeat later"
	SendingEmailWithIntstructionsNotifications = "We sent you an email with needed instructions"
)

type Config struct {
	MongoURI          string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisURI          string `mapstructure:"REDIS_URL"`
	Port              string `mapstructure:"PORT"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	ClientOriginUrl string `mapstructure:"CLIENT_ORIGIN_URL"`

	EmailFrom    string `mapstructure:"EMAIL_FROM"`
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
