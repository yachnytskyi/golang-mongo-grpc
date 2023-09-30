package constant

import "time"

const (
	// Config paths.
	ConfigPath = "config/yaml/application.yaml"

	// Context timers.
	DefaultContextTimer = time.Duration(time.Second * 5)

	// Pagination.
	DefaultPage     = "1"
	DefaultLimit    = "10"
	MaxItemsPerPage = 100

	// Regex patterns.
	StringRegex      = `^[a-zA-z0-9 !@#$€%^&*{}|()=/\;:+-_~'"<>,.? \t]*$`
	TitleStringRegex = `^[a-zA-z0-9 !()=[]:;+-_~'",.? \t]*$`
	TextStringRegex  = `^[a-zA-z0-9 !@#$€%^&*{}][|/\()=/\;:+-_~'"<>,.? \t]*$`
	MinLength        = 4
	MaxLength        = 40

	// User Notifications.
	SendingEmailNotification = "We sent an email with a verification code to "

	// Error Messages.
	StringAllowedLength                        = "can be between %d and %d characters long"
	EmailAlreadyExists                         = "user with this email already exists"
	SendingEmailWithIntstructionsNotifications = "We sent you an email with needed instructions"
	InternalErrorNotification                  = "something went wrong, please repeat later"
	EntityNotFoundErrorNotification            = "please repeat it later"

	// Databases.
	MongoDB = "MongoDB"

	// Domains.
	UseCase = "UseCase"
)