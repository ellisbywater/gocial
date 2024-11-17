package mailer

import "embed"

const (
	FromName                    = "Gocial"
	maxRetries                  = 3
	UserEmailActivationTemplate = "user_invitation.go.tmpl"
)

// go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username string, email string, data any, isSandbox bool) error
}
