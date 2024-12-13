package mailer

import (
	"embed"
	"encoding/json"
	"errors"
	"strings"
)

const (
	FromName            = "RaulSocialMedia"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.gohtml"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}

func prepareTemplateHtmlForJSON(template string) string {
	trimmedStr := strings.TrimSpace(template)
	trimmedStr = strings.ReplaceAll(trimmedStr, "\"", "\\\"")
	trimmedStr = strings.ReplaceAll(strings.ReplaceAll(trimmedStr, "\r\n", ""), "\n", "")

	return trimmedStr
}

func failedSendingEmail(data []byte) error {
	type EmailResponse struct {
		Errors []string `json:"errors"`
	}

	var emailResponse EmailResponse

	err := json.Unmarshal(data, &emailResponse)
	if err != nil {
		return err
	}

	if len(emailResponse.Errors) != 0 {
		return errors.New(strings.Join(emailResponse.Errors, ","))
	}

	return nil
}
