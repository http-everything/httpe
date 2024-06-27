package sendemail_test

import (
	"fmt"
	"testing"

	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/actions/sendemail"
	"github.com/http-everything/httpe/pkg/config"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"
)

func TestSendEmailShouldFail(t *testing.T) {
	cases := []struct {
		name          string
		email         *rules.Email
		wantErrorBody string
	}{
		{
			name: "Bad to",
			email: &rules.Email{
				To:   "foo@bla",
				From: "user@example.com",
				Body: "1234abc",
			},
			wantErrorBody: "email to is not a valid email address",
		},
		{
			name: "Bad from",
			email: &rules.Email{
				To:   "user@example.com",
				From: "user@example com",
				Body: "1234abc",
			},
			wantErrorBody: "email from is not a valid email address",
		},
		{
			name: "Bad cc",
			email: &rules.Email{
				To:   "user@example.com",
				Cc:   "bad email.com",
				Body: "1234abc",
			},
			wantErrorBody: "email cc is not a valid email address",
		},
		{
			name: "Bad bcc",
			email: &rules.Email{
				To:   "user@example.com",
				Bcc:  "bad email.com",
				Body: "1234abc",
			},
			wantErrorBody: "email bcc is not a valid email address",
		},
	}

	reqData := requestdata.Data{}
	smtpConfig := &config.SMTPConfig{
		Server: "127.0.0.1",
		Port:   25,
		From:   "user@example.com",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rule := rules.Rule{
				SendEmail: tc.email,
			}
			//Create the actioner that implements the action interface
			var actioner actions.Actioner = sendemail.Email{SMTPConfig: smtpConfig}

			// Execute the action by calling the mandatory function Execute()
			actionResp, err := actioner.Execute(rule, reqData)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantErrorBody, actionResp.ErrorBody)
		})
	}
}

func TestSendEmailShouldSucceed(t *testing.T) {
	cases := []struct {
		name          string
		email         *rules.Email
		wantErrorBody string
	}{
		{
			name: "Good Email",
			email: &rules.Email{
				To:      "to1@example.com",
				Subject: "Test",
				Body:    "1",
				From:    "sender1@example.com",
			},
		},
		{
			name: "Good Email with Cc",
			email: &rules.Email{
				To:      "to2@example.com",
				Cc:      "cc@example.com",
				Subject: "Test",
				Body:    "2",
				From:    "sender2@example.com",
			},
		},
		{
			name: "Good Email with Bcc",
			email: &rules.Email{
				To:      "to3@example.com",
				Bcc:     "bcc@example.com",
				Subject: "Test",
				Body:    "3",
				From:    "sender3@example.com",
			},
		},
	}
	// Start smtpmock server
	smtpServer := smtpmock.New(smtpmock.ConfigurationAttr{
		LogToStdout:       true, // Enable logging for debugging
		LogServerActivity: false,
	})
	assert.Nil(t, smtpServer.Start())
	defer func() {
		err := smtpServer.Stop()
		assert.NoError(t, err)
	}()

	reqData := requestdata.Data{}
	smtpConfig := &config.SMTPConfig{
		Server: "127.0.0.1",
		Port:   smtpServer.PortNumber(),
		From:   "default@example.com",
	}

	for index, tc := range cases {
		rule := rules.Rule{
			SendEmail: tc.email,
		}
		//Create the actioner that implements the action interface
		var actioner actions.Actioner = sendemail.Email{SMTPConfig: smtpConfig}

		// Execute the action by calling the mandatory function Execute()
		actionResp, err := actioner.Execute(rule, reqData)

		assert.NoError(t, err)
		assert.Emptyf(t, actionResp.ErrorBody, "ErrorBody not empty")
		assert.Equal(t, "email sent", actionResp.SuccessBody)
		assert.Equal(t, 0, actionResp.Code)

		// Additionally, verify that an email was indeed sent to the SMTP mock
		msg := smtpServer.MessagesAndPurge()[0].MsgRequest()
		t.Logf("%d: Body Res: %s", index, msg)

		assert.Contains(t, msg, fmt.Sprintf("From: %s\r\n", rule.SendEmail.From))
		assert.Contains(t, msg, fmt.Sprintf("To: %s\r\n", rule.SendEmail.To))
		if rule.SendEmail.Cc != "" {
			assert.Contains(t, msg, fmt.Sprintf("Cc: %s\r\n", rule.SendEmail.Cc))
		}
		assert.Contains(t, msg, fmt.Sprintf("Subject: %s\r\n", rule.SendEmail.Subject))
		assert.Contains(t, msg, fmt.Sprintf("\r\n\r\n%s\r\n", rule.SendEmail.Body))
	}
}
