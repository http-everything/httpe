package requesthandler

import (
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/answercontent"
	"http-everything/httpe/pkg/actions/answerfile"
	"http-everything/httpe/pkg/actions/runscript"
	"http-everything/httpe/pkg/auth"
	"http-everything/httpe/pkg/requestdata"
	"http-everything/httpe/pkg/response"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/logger"
	"net/http"
)

func Execute(rule rules.Rule, logger *logger.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Initialise a new http response writer.
		respWriter := response.New(w, rule.Respond, logger)

		// Authenticate, if requested by the rule
		if rule.With != nil {
			ok, err := auth.IsRequestAuthenticated(rule.With.AuthBasic, rule.With.AuthHashing, r)
			if err != nil {
				respWriter.InternalServerError(err)
				return
			}
			if !ok {
				respWriter.Unauthorised()
				return
			}
		}

		// Collect data from the request to be made available to the template engine and add to the response writer
		reqData, err := requestdata.Collect(r)
		if err != nil {
			respWriter.InternalServerError(err)
			return
		}
		respWriter.AddRequestData(reqData)

		//Create a container for the action that implements the action interface
		var actioner actions.Actioner

		// Hand over the request to the action specified by the rule defined by 'rule.Do' using switch case
		switch rule.Action() {
		case rules.RunScript:
			// Execute a script
			actioner = runscript.Script{}
		case rules.SendEmail:
			// Send an email
			actioner = runscript.Script{}
		case rules.AnswerContent:
			actioner = answercontent.AnswerContent{}
		case rules.AnswerFile:
			actioner = answerfile.AnswerFile{}
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
