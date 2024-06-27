package rules

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/http-everything/httpe/pkg/config"

	"github.com/http-everything/httpe/pkg/share/logger"

	"gopkg.in/yaml.v3"

	"github.com/xeipuuv/gojsonschema"
)

type Rules struct {
	Rules  *[]Rule `yaml:"rules" json:"rules"`
	logger *logger.Logger
}

type Cfg struct {
	Rules []Rule `yaml:"rules" json:"rules"`
}

//go:embed schema.json
var schemaJSON string

// Read reads the yaml file containing the rules
func Read(yamlFile string, logger *logger.Logger) (rules *Rules, err error) {
	yml, err := os.ReadFile(yamlFile)
	if err != nil {
		return &Rules{}, fmt.Errorf("error reading yaml file '%s': %w", yamlFile, err)
	}

	cfg := Cfg{}
	err = yaml.Unmarshal(yml, &cfg)
	if err != nil {
		return &Rules{}, fmt.Errorf("error parsing yaml file '%s': %w", yamlFile, err)
	}
	return &Rules{
		Rules:  &cfg.Rules,
		logger: logger,
	}, nil
}

func YamlToJSON(yamlFile string) string {
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return err.Error()
	}

	var data interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return err.Error()
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err.Error()
	}

	return string(jsonData)
}

func (r *Rules) Validate(smtpConfig *config.SMTPConfig) (err error) {
	// Convert the rules into JSON
	JSONConf, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if len(*r.Rules) == 0 {
		return errors.New("no rules specified")
	}

	// Validate the rule action. Doing it here and before the schema validation creates a more user-friendly error
	// message. While it is possible to do the below validation with the 'oneOf' statement of JSON schema validation,
	// the error message is confusing because only the first element of the oneOf list is named as required.
	var hasErrors = false
	for i, rule := range *r.Rules {
		if rule.Action() == "" {
			r.logger.PrintAndLogErrorf("rule %d '%s' is missing a valid action. Use one of '%s'.",
				i,
				rule.Name, strings.Join(ValidActions, ", "))
			hasErrors = true
		}
		if rule.Action() == SendEmail {
			if smtpConfig == nil {
				r.logger.PrintAndLogErrorf("rule %d: %s requires an smtp configuration in httpe configuration file",
					i,
					SendEmail,
				)
				hasErrors = true
			}
		}
	}
	if hasErrors {
		return fmt.Errorf("invalid rules: at least one rule is not valid")
	}

	// Validate against schema
	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	documentLoader := gojsonschema.NewStringLoader(string(JSONConf))
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("json schema validation failed: %w", err)
	}

	if result.Valid() {
		r.logger.Debugf("successfully validated against schema")
	} else {
		r.logger.PrintAndLogErrorf("schema validation against %s failed", SchemaURL)
		for _, desc := range result.Errors() {
			r.logger.PrintAndLogErrorf("%s\n", desc)
		}
		return fmt.Errorf("invalid rules: schema validation failed")
	}
	return nil
}

func (rule *Rule) Action() (action string) {
	if rule.RunScript != "" {
		return RunScript
	}
	if rule.SendEmail != nil {
		return SendEmail
	}
	if rule.AnswerContent != "" {
		return AnswerContent
	}
	if rule.AnswerFile != "" {
		return AnswerFile
	}
	if rule.RedirectPermanent != "" {
		return RedirectPermanent
	}
	if rule.RedirectTemporary != "" {
		return RedirectTemporary
	}
	if rule.ServeDirectory != "" {
		return ServeDirectory
	}
	if len(rule.RenderButtons) > 0 {
		return RenderButtons
	}
	return ""
}

func (rule *Rule) MaxRequestBody() string {
	if rule.With != nil {
		return rule.With.MaxRequestBody
	}
	return ""
}

func (ruleResp *Respond) Headers(onSuccess bool) map[string]string {
	if onSuccess {
		return ruleResp.OnSuccess.Headers
	}
	return ruleResp.OnError.Headers
}

func (rule *Rule) MatchesURI(URI string) bool {
	if rule.On.Path == URI {
		// Exact match
		return true
	}
	if !strings.HasPrefix(rule.On.Path, "{") && !strings.HasSuffix(rule.On.Path, "}") {
		// The rule is not using URL Placeholders, e.g. /foo/{key1}/{key2}/
		return false
	}
	return false
}
