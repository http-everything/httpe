package rules

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"http-everything/httpe/pkg/share/logger"
	"os"
	"strings"

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

func (r *Rules) Validate() (err error) {
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
			r.logger.Errorf("rule %d '%s' is missing a valid action in the 'do' section. Use one of '%s'.",
				i,
				rule.Name, strings.Join(ValidActions, ", "))
			hasErrors = true
		}
	}
	if hasErrors {
		return fmt.Errorf("invalid rules")
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
		r.logger.Errorf("schema validation against %s failed", SchemaURL)
		for _, desc := range result.Errors() {
			r.logger.Errorf("%s\n", desc)
		}
		return fmt.Errorf("invalid rules")
	}
	return nil
}

func (rule *Rule) Action() (action string) {
	if rule.Do == nil {
		return ""
	}
	if rule.Do.RunScript != "" {
		return RunScript
	}
	if rule.Do.SendEmail != "" {
		return SendEmail
	}
	if rule.Do.AnswerContent != "" {
		return AnswerContent
	}
	if rule.Do.AnswerFile != "" {
		return AnswerFile
	}
	if rule.Do.RedirectPermanent != "" {
		return RedirectPermanent
	}
	if rule.Do.RedirectTemporary != "" {
		return RedirectTemporary
	}
	if rule.Do.ServeDirectory != "" {
		return ServeDirectory
	}
	if len(rule.Do.RenderButtons) > 0 {
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

func (ruleResp *Respond) Headers(onError bool) map[string]string {
	if onError {
		return ruleResp.OnError.Headers
	}
	return ruleResp.OnSuccess.Headers
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
