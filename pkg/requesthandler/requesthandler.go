package requesthandler

import (
	"net/http"

	"github.com/http-everything/httpe/pkg/actions/sendemail"
	"github.com/http-everything/httpe/pkg/config"

	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/actions/answercontent"
	"github.com/http-everything/httpe/pkg/actions/answerfile"
	"github.com/http-everything/httpe/pkg/actions/redirect"
	"github.com/http-everything/httpe/pkg/actions/renderbuttons"
	"github.com/http-everything/httpe/pkg/actions/runscript"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/response"
	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/share/logger"
)

const DefaultMaxRequestBody = "512KB"

func Execute(rule rules.Rule, logger *logger.Logger, smtpConfig *config.SMTPConfig) http.Handler {
	//return http.StripPrefix("/dir", http.FileServer(http.Dir("/Users/thorsten/tmp")))
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Initialise a new http response writer.
		respWriter := response.New(w, rule.Respond, logger)

		// Collect data from the request to be made available to the template engine and add to the response writer
		reqData, err := requestdata.Collect(r, rule.Args)
		if err != nil {
			respWriter.InternalServerError(err)
			return
		}
		respWriter.AddRequestData(reqData)

		//Create a container for the action that implements the action interface
		var actioner actions.Actioner

		// Hand over the request to the action specified by the rule defined by 'rule' using switch case
		switch rule.Action() {
		case rules.RunScript:
			// Execute a script
			actioner = runscript.Script{}
		case rules.SendEmail:
			// Send an email
			actioner = sendemail.Email{
				SMTPConfig: smtpConfig,
			}
		case rules.AnswerContent:
			actioner = answercontent.AnswerContent{}
		case rules.AnswerFile:
			actioner = answerfile.AnswerFile{}
		case rules.RedirectPermanent, rules.RedirectTemporary:
			actioner = redirect.Redirect{}
		case rules.RenderButtons:
			actioner = renderbuttons.RenderButtons{}

		default:
			// Do nothing, just create a response
			actioner = answercontent.AnswerContent{}
		}
		// Execute the action by calling the mandatory function Execute()
		actionResp, err := actioner.Execute(rule, reqData)
		if err != nil {
			respWriter.InternalServerErrorf("action %s: %s", rule.Action(), err)
			return
		}
		// Hand over the action response to our HTTP response writer
		respWriter.ActionResponse(actionResp)
	}
	return http.HandlerFunc(fn)
}
