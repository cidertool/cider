package template

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/aaronsky/applereleaser/pkg/context"
)

const (
	projectNameKey = "projectName"
	versionKey     = "version"
	envKey         = "env"
	dateKey        = "date"
	timestampKey   = "timestamp"
)

type Template struct {
	fields Fields
}

type Fields map[string]interface{}

func New(ctx *context.Context) *Template {
	return &Template{
		Fields{
			projectNameKey: ctx.Config.Name,
			versionKey:     ctx.Version,
			envKey:         ctx.Env,
			dateKey:        ctx.Date.UTC().Format(time.RFC3339),
			timestampKey:   ctx.Date.UTC().Unix(),
		},
	}
}

func (t *Template) WithFields(fields Fields) *Template {
	for key, value := range fields {
		t.fields[key] = value
	}
	return t
}

func (t *Template) WithEnv(env map[string]string) *Template {
	t.fields[envKey] = env
	return t
}

func (t *Template) WithShellEnv(envs ...string) *Template {
	env := make(map[string]string)
	for _, e := range envs {
		parts := strings.SplitN(e, "=", 2)
		env[parts[0]] = parts[1]
	}
	return t.WithEnv(env)
}

func (t *Template) Apply(s string) (string, error) {
	var out bytes.Buffer
	tmpl, err := template.New("tmpl").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"replace":    strings.ReplaceAll,
			"lowercased": strings.ToLower,
			"uppercased": strings.ToUpper,
			"titlecased": strings.ToTitle,
			"dir":        filepath.Dir,
			"abs":        filepath.Abs,
			"rel":        filepath.Rel,
		}).
		Parse(s)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&out, t.fields)
	return out.String(), err
}
