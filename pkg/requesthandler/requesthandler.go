package requesthandler

import (
	"http-everything/httpe/pkg/actions"
	"http-everything/httpe/pkg/actions/script"
	"http-everything/httpe/pkg/auth"
	"http-everything/httpe/pkg/request"
	"http-everything/httpe/pkg/response"
	"http-everything/httpe/pkg/rules"
	"http-everything/httpe/pkg/share/logger"
	"log"
	"net/http"
)

func Execute(rule rules.Rule, logger *logger.Logger) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		rw := response.New(w, rule.Respond, request.BuildMetaData(r), logger)
		err := r.ParseForm()
		if err != nil {
			rw.InternalServerError(err)
			return
		}
		for k, v := range r.PostForm {
			for i, value := range v {
				log.Printf("k: %s, i: %d, value: %s \n", k, i, value)
			}
		}
		if rule.With != nil {
			ok, err := auth.IsRequestAuthenticated(rule.With.AuthBasic, rule.With.AuthHashing, r)
			if err != nil {
				rw.InternalServerError(err)
				return
			}
			if !ok {
				http.Error(w, "Unauthorised", http.StatusUnauthorized)
				return
			}
		}
		var a actions.Action //Create a container for the action that implements the action interface

		// Hand over the request to the action specified by the rule action defined by rule.Do using switch case
		do := rule.GetAction()
		switch do {
		case rules.RuleActionScript:
			// Execute a script
			a = script.Script{}
		case rules.RuleActionEmail:
			// Send an email
			a = script.Script{}
		default:
			// Do nothing, just create a response
			a = script.Script{}
		}
		// Execute the action by calling the mandatory function Execute()
		resp, err := a.Execute(rule)
		if err != nil {
			rw.InternalServerErrorf("action=%s: %s", do, err)
			return
		}
		rw.ActionResponse(resp)
	}
	return http.HandlerFunc(fn)
}
