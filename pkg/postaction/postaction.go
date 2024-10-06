package postaction

import (
	"fmt"
	"github.com/http-everything/httpe/pkg/actions"
	"github.com/http-everything/httpe/pkg/actions/runscript"
	"github.com/http-everything/httpe/pkg/actions/sendemail"
	"github.com/http-everything/httpe/pkg/config"
	"github.com/http-everything/httpe/pkg/postactionresponsewriter"
	"github.com/http-everything/httpe/pkg/requestdata"
	"github.com/http-everything/httpe/pkg/rules"
	"github.com/http-everything/httpe/pkg/share/logger"
)

func Execute(postActionRule rules.Rule, reqData requestdata.Data, conf *config.Config, logger *logger.Logger) {
	if postActionRule.PostAction == nil {
		return
	}
	//Create a container for the action that implements the action interface
	var actioner actions.Actioner
	var rule rules.Rule
	prw := postactionresponsewriter.New(conf, logger)
	println("Starting postaction")
	if postActionRule.PostAction.RunScript != "" {
		fmt.Println("Starting postaction run.script")
		actioner = runscript.Script{}
		rule.RunScript = postActionRule.PostAction.RunScript
		rule.Args = postActionRule.PostAction.Args
		resp, err := actioner.Execute(rule, reqData)
		prw.AddActionResponse(rules.RunScript, resp, err)
	}

	if postActionRule.PostAction.SendEmail != nil {
		actioner = sendemail.Email{
			SMTPConfig: conf.SMTP,
		}
		rule.SendEmail = postActionRule.PostAction.SendEmail
		rule.Args = postActionRule.PostAction.Args
		resp, err := actioner.Execute(rule, reqData)
		prw.AddActionResponse(rules.SendEmail, resp, err)
	}
	prw.Write()
}
