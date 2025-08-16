package default_email_templates

import "embed"

//go:embed *
var templatesFolder embed.FS

type EmailEventType string

const (
	EMAIL_CONFIRM          = "email_confirm"
	PASSWORD_RESET_CONFIRM = "password_reset_confirm"
)

type EmailPart string

const (
	BODY    EmailPart = "body.html"
	SUBJECT EmailPart = "subject.tmpl"
)

func GetTemplate(eventType EmailEventType, emailPart EmailPart) string {
	result, err := templatesFolder.ReadFile(string(eventType) + "_" + string(emailPart))
	if err != nil {
		return ""
	}
	return string(result)
}
