package sendemail

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"gopkg.in/gomail.v2"

	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/config"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/templating"
)

type Email struct {
	SMTPConfig *config.SMTPConfig
}

// Execute implements the actioner interface, being the final method executed by the action
func (e Email) Execute(rule rules.Rule, reqData requestdata.Data) (response actions.ActionResponse, err error) {
	body, err := templating.RenderString(rule.SendEmail.Body, reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}

	subject, err := templating.RenderString(rule.SendEmail.Subject, reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}

	to, err := templating.RenderString(rule.SendEmail.To, reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}
	cc, err := templating.RenderString(rule.SendEmail.Cc, reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}
	bcc, err := templating.RenderString(rule.SendEmail.Bcc, reqData)
	if err != nil {
		return actions.ActionResponse{}, err
	}

	var from string
	if rule.SendEmail.From != "" {
		from, err = templating.RenderString(rule.SendEmail.From, reqData)
		if err != nil {
			return actions.ActionResponse{}, err
		}
	} else if e.SMTPConfig.From != "" {
		from = e.SMTPConfig.From
	} else {
		return actions.ActionResponse{}, errors.New("no email from specified")
	}

	// Validate email addresses
	addrs := map[string]string{
		"from": from,
		"to":   to,
		"cc":   cc,
		"bcc":  bcc,
	}
	for label, addr := range addrs {
		if addr != "" && !govalidator.IsEmail(addr) {
			return actions.ActionResponse{
				Code:            1,
				ErrorHTTPStatus: http.StatusBadRequest,
				ErrorBody:       fmt.Sprintf("email %s is not a valid email address", label),
			}, nil
		}
	}

	err = e.sendEmail(from, to, subject, body, cc, bcc)
	if err != nil {
		return actions.ActionResponse{
			Code:            1,
			ErrorHTTPStatus: http.StatusBadRequest,
			ErrorBody:       fmt.Sprintf("SMTP connection error: %s", err),
		}, nil
	}
	return actions.ActionResponse{
		SuccessBody: "email sent",
	}, nil
}

// sendEmail sends an email using the local SMTP server.
func (e Email) sendEmail(from, to, subject, body, cc, bcc string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	m.SetHeader("To", to)

	if cc != "" {
		m.SetHeader("Cc", cc)
	}
	if bcc != "" {
		m.SetHeader("Bcc", bcc)
	}

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := gomail.NewDialer(e.SMTPConfig.Server, e.SMTPConfig.Port, e.SMTPConfig.Username, e.SMTPConfig.Password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail and return errors
	return d.DialAndSend(m)
}
