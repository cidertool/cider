// Package template provides an interface for text templates to be used during pipes
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

// Template is used to apply text templates to strings to dynamically configure API values. See the documentation of
// text/template to see the valid template format.
type Template struct {
	fields Fields
}

// Fields is a heterogenous map type keyed by strings.
type Fields map[string]interface{}

// New returns a new template instance.
func New(ctx *context.Context) *Template {
	return &Template{
		Fields{
			projectNameKey: ctx.RawConfig.Name,
			versionKey:     ctx.Version,
			envKey:         ctx.Env,
			dateKey:        ctx.Date.UTC().Format(time.RFC3339),
			timestampKey:   ctx.Date.UTC().Unix(),
		},
	}
}

// WithFields merges the template's configured fields with the given Fields.
func (t *Template) WithFields(fields Fields) *Template {
	for key, value := range fields {
		t.fields[key] = value
	}
	return t
}

// WithEnv replaces the configured env of the template with the given key-value map.
func (t *Template) WithEnv(env map[string]string) *Template {
	t.fields[envKey] = env
	return t
}

// WithShellEnv replaces the configured env of the template with the given sequence of shell-style, e.g. "KEY=VALUE", strings.
func (t *Template) WithShellEnv(envs ...string) *Template {
	env := make(map[string]string)
	for _, e := range envs {
		parts := strings.SplitN(e, "=", 2)
		env[parts[0]] = parts[1]
	}
	return t.WithEnv(env)
}

// Apply takes the template string and processes it into its product string.
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
