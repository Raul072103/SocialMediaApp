package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

type mailTrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (*mailTrapClient, error) {
	if apiKey == "" {
		return &mailTrapClient{}, errors.New("api key is required for sending emails")
	}

	return &mailTrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m *mailTrapClient) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

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

	//message := gomail.NewMessage()
	//message.SetHeader("From", from.Address)
	//message.SetHeader("To", to.Address)
	//message.SetHeader("Subject", subject.String())
	//
	//message.AddAlternative("text/html", body.String())

	htmlMessage := strings.TrimSpace(body.String())

	fmt.Printf(htmlMessage)

	message := []byte(fmt.Sprintf(`{
		"from":{"email": "%s"},
		"to":[{"email": "%s"}],
		"subject":"%s",
		"text": "ceva_text",
		"html": "%s"
	}`, from.Address, to.Address, subject.String(), htmlMessage))

	err = sendEmailAsync(m.apiKey, message)
	if err != nil {
		return err
	}

	return nil
}

func sendEmailAsync(token string, message []byte) error {
	httpHost := "https://send.api.mailtrap.io/api/send"

	// Set up request
	request, err := http.NewRequest(http.MethodPost, httpHost, bytes.NewBuffer(message))
	if err != nil {
		return err
	}

	// Set required headers
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	// Send request asynchronously
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

	fmt.Println(string(body))
	return nil
}
