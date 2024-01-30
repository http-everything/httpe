package rules

const (
	SchemaURL         = "https://github.com/http-everything/httpe/main/pkg/rules/schema.json"
	RunScript         = "run.script"
	SendEmail         = "send.email"
	AnswerContent     = "answer.content"
	AnswerFile        = "answer.file"
	RedirectPermanent = "redirect.permanent"
	RedirectTemporary = "redirect.temporary"
	ServeDirectory    = "serve.directory"
)

var ValidActions = []string{
	RunScript,
	SendEmail,
	AnswerFile,
	AnswerContent,
	RedirectPermanent,
	RedirectTemporary,
	ServeDirectory,
}

type Rule struct {
	Name    string  `yaml:"name,omitempty" json:"name,omitempty"`
	On      *On     `yaml:"on" json:"on"`
	Do      *Do     `yaml:"do" json:"do"`
	With    *With   `yaml:"with" json:"with,omitempty"`
	Respond Respond `yaml:"respond" json:"respond"`
}

type On struct {
	Path    string   `yaml:"path,omitempty" json:"path,omitempty"`
	Methods []string `yaml:"methods,omitempty" json:"methods,omitempty"`
}

type Do struct {
	RunScript         string `yaml:"run.script,omitempty" json:"run.script,omitempty"`
	SendEmail         string `yaml:"send.email,omitempty" json:"send.email,omitempty"`
	AnswerContent     string `yaml:"answer.content,omitempty" json:"answer.content,omitempty"`
	AnswerFile        string `yaml:"answer.file,omitempty" json:"answer.file,omitempty"`
	RedirectPermanent string `yaml:"redirect.permanent,omitempty" json:"redirect.permanent,omitempty"`
	RedirectTemporary string `yaml:"redirect.temporary,omitempty" json:"redirect.temporary,omitempty"`
	ServeDirectory    string `yaml:"serve.directory,omitempty" json:"serve.directory,omitempty"`
	Args              Args   `yaml:"args" json:"args"`
}

type Args struct {
	Interpreter string `yaml:"interpreter" json:"interpreter"`
	Timeout     int    `yaml:"timeout" json:"timeout"`
	Cwd         string `yaml:"cwd" json:"cwd"`
}

type With struct {
	AuthBasic      []User `yaml:"auth_basic,omitempty" json:"auth_basic,omitempty"`
	AuthHashing    string `yaml:"auth_hashing,omitempty" json:"auth_hashing,omitempty"`
	MaxRequestBody string `yaml:"max_request_body,omitempty" json:"max_request_body,omitempty"`
}

type User struct {
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

type Respond struct {
	OnSuccess OnSuccess `yaml:"on_success" json:"on_success"`
	OnError   OnError   `yaml:"on_error" json:"on_error"`
}

type OnSuccess struct {
	HTTPStatus int     `yaml:"http_status" json:"http_status"`
	Body       string  `yaml:"body" json:"body"`
	Headers    Headers `yaml:"headers" json:"headers"`
}

type OnError struct {
	HTTPStatus int     `yaml:"http_status" json:"http_status"`
	Body       string  `yaml:"body" json:"body"`
	Headers    Headers `yaml:"headers" json:"headers"`
}

type Headers map[string]string
