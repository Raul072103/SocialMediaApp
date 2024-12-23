package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"io"
	"net/http"
	"time"
)

type mailTrapClient struct {
	fromEmail      string
	toEmailDefault string
	apiKey         string
}

func NewMailTrapClient(apiKey, fromEmail, toEmailDefault string) (*mailTrapClient, error) {
	if apiKey == "" {
		return &mailTrapClient{}, errors.New("api key is required for sending emails")
	}

	return &mailTrapClient{
		fromEmail:      fromEmail,
		toEmailDefault: toEmailDefault,
		apiKey:         apiKey,
	}, nil
}

func (m *mailTrapClient) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)

	var to *mail.Email
	if isSandbox {
		to = mail.NewEmail(FromName, m.toEmailDefault)
	} else {
		to = mail.NewEmail(username, email)
	}

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := []byte(fmt.Sprintf(`{
		"from":{"email": "%s"},
		"to":[{"email": "%s"}],
		"subject":"%s",
		"html": "%s"
	}`, from.Address, to.Address, subject.String(), prepareTemplateHtmlForJSON(body.String())))

	err = sendEmail(m.apiKey, message)
	if err != nil {
		return err
	}

	return nil
}

func sendEmail(token string, message []byte) error {
	httpHost := "https://send.api.mailtrap.io/api/send"

	// Set up request
	request, err := http.NewRequest(http.MethodPost, httpHost, bytes.NewBuffer(message))
	if err != nil {
		return err
	}

	// Set required headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = failedSendingEmail(body)
	if err != nil {
		return err
	}

	return nil
}
