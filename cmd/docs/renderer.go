package docs

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

// ErrRootTypeNotFound indicates an error when the root type of the configuration file
// could not be found in the project source.
var ErrRootTypeNotFound = errors.New("root type config.Project not found")

type stringWriterTo interface {
	io.StringWriter
	io.WriterTo
}

type renderTypeOptions struct {
	Name  string
	Level int
	Doc   string
	Type  ast.Expr
	// Parent *doc.Type
}

type docRenderer struct {
	Package       *doc.Package
	Types         map[string]*doc.Type
	Values        map[string][]string
	Visited       map[string]bool
	TypesToRender []renderTypeOptions

	Header func() string
	Footer func() string

	buffer stringWriterTo
}

func newRenderer() (*docRenderer, error) {
	r := new(docRenderer)

	pkg, err := openPackage()
	if err != nil {
		return nil, err
	}

	r.Package = pkg

	r.Types = make(map[string]*doc.Type, len(pkg.Types))
	r.Values = make(map[string][]string)
	r.Visited = make(map[string]bool, len(pkg.Types))
	r.TypesToRender = make([]renderTypeOptions, 0, len(r.Types))
	r.buffer = new(bytes.Buffer)

	err = r.gatherTypes()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func openPackage() (*doc.Package, error) {
	fset := token.NewFileSet()

	files, err := parseFilesInConfigPackage(fset)
	if err != nil {
		return nil, err
	}

	return doc.NewFromFiles(fset, files, "github.com/cidertool/cider/pkg/config")
}

func parseFilesInConfigPackage(fset *token.FileSet) (files []*ast.File, err error) {
	files = make([]*ast.File, 0)

	inDir, err := ioutil.ReadDir("pkg/config")
	if err != nil {
		return files, err
	}

	for _, finfo := range inDir {
		if finfo.IsDir() || filepath.Ext(finfo.Name()) != ".go" {
			continue
		}

		f, err := parser.ParseFile(fset, filepath.Join("pkg/config", finfo.Name()), nil, parser.ParseComments)
		if err != nil {
			return files, err
		}

		files = append(files, f)
	}

	return files, err
}

func (r *docRenderer) gatherTypes() error {
	for _, typ := range r.Package.Types {
		r.Types[typ.Name] = typ

		if len(typ.Consts) > 0 {
			for _, cons := range typ.Consts {
				name, values := r.gatherConsts(cons)
				r.Values[name] = values
			}
		}
	}

	for _, cons := range r.Package.Consts {
		name, values := r.gatherConsts(cons)
		r.Values[name] = values
	}

	root, ok := r.Types["Project"]
	if !ok ||
		root == nil ||
		len(root.Decl.Specs) == 0 {
		return ErrRootTypeNotFound
	}

	r.insertTypeTree(root, 0)

	return nil
}

func (r *docRenderer) insertTypeTree(t *doc.Type, level int) {
	getSubTypeNames := func(e ast.Expr) []string {
		switch t := e.(type) {
		case *ast.ArrayType:
			return []string{getTypeName(t)}
		case *ast.StructType:
			var fields = make([]string, len(t.Fields.List))
			for i, field := range t.Fields.List {
				fields[i] = getTypeName(field.Type)
			}

			return fields
		case *ast.MapType:
			return []string{getTypeName(t.Key), getTypeName(t.Value)}
		case *ast.StarExpr:
			return []string{getTypeName(t)}
		case *ast.SelectorExpr:
			return []string{getTypeName(t)}
		}

		return nil
	}
	decl := t.Decl

	for _, s := range decl.Specs {
		spec, ok := s.(*ast.TypeSpec)
		if !ok {
			continue
		} else if r.Visited[spec.Name.Name] {
			continue
		}

		r.Visited[spec.Name.Name] = true
		r.TypesToRender = append(r.TypesToRender, renderTypeOptions{
			Name:  spec.Name.Name,
			Level: level,
			Doc:   formatDoc(t.Doc),
			Type:  spec.Type,
		})

		for _, name := range getSubTypeNames(spec.Type) {
			typ, ok := r.Types[name]
			if !ok {
				continue
			}

			r.insertTypeTree(typ, level+1)
		}
	}
}

func (*docRenderer) gatherConsts(cons *doc.Value) (name string, values []string) {
	values = make([]string, len(cons.Decl.Specs))

	for i, s := range cons.Decl.Specs {
		if spec, ok := s.(*ast.ValueSpec); ok {
			name = spec.Type.(*ast.Ident).Name
			values[i] = spec.Values[0].(*ast.BasicLit).Value
		}
	}

	return name, values
}

func (r *docRenderer) WriteString(s string) {
	_, err := r.buffer.WriteString(s)
	if err != nil {
		log.Error(err.Error())
	}
}

func (r *docRenderer) Render(w io.Writer) error {
	if r.Header != nil {
		r.WriteString(r.Header())
	}

	r.WriteString(r.Package.Doc)
	r.WriteString(docsConfigTerminologyDisclaimer)
	r.WriteString(docsConfigTableOfContents)
	r.WriteString("## Specification\n\n")

	var topHeaderLevel = 3

	for _, opt := range r.TypesToRender {
		switch typ := opt.Type.(type) {
		case *ast.StructType:
			r.renderTypePreamble(opt.Name, opt.Doc, topHeaderLevel+opt.Level)
			r.renderStruct(typ)
		case *ast.MapType:
			r.renderTypePreamble(opt.Name, opt.Doc, topHeaderLevel+opt.Level)
			r.renderMap(typ)
		default:
			continue
		}
		r.WriteString("\n")
	}

	proj, err := yaml.Marshal(docsConfigExampleProject)
	if err != nil {
		return err
	}

	r.WriteString(fmt.Sprintf("## Full Example\n\n```yaml\n%s```\n\n", string(proj)))

	if r.Footer != nil {
		r.WriteString(r.Footer())
	}

	_, err = r.buffer.WriteTo(w)

	return err
}

func (r *docRenderer) renderTypePreamble(name string, doc string, level int) {
	var maxHeaderLevel = 6
	if level > maxHeaderLevel {
		level = maxHeaderLevel
	}

	r.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", level), name))

	if doc != "" {
		r.WriteString(doc + "\n\n")
	}
}

func (r *docRenderer) renderStruct(typ *ast.StructType) {
	for _, field := range typ.Fields.List {
		typeName := getTypeName(field.Type)
		typeNameFormatted := formatTypeName(field.Type)
		tag, required := getTagValue(field.Tag)

		var requiredStr = " "
		if required && !strings.HasPrefix(typeNameFormatted, "[") {
			requiredStr = "x"
		}

		var options []string

		if values, ok := r.Values[typeName]; ok {
			if identName := r.asIdentName(typeName); identName != "" {
				typeNameFormatted = identName
			} else {
				typeNameFormatted = "string"
			}

			options = values
		}

		if len(field.Names) == 0 {
			// embedded struct field
			typeName := getTypeName(field.Type)

			t, ok := r.Types[typeName]
			if !ok {
				continue
			}

			for _, s := range t.Decl.Specs {
				if spec, ok := s.(*ast.TypeSpec); ok {
					if f, ok := spec.Type.(*ast.StructType); ok {
						r.renderStruct(f)
					}
				}
			}

			continue
		}

		line := fmt.Sprintf(
			"- [%s] **%s: %s** – %s%s",
			requiredStr,
			tag,
			r.renderItemType(typeName, typeNameFormatted),
			formatDoc(field.Doc.Text()),
			formatOptions(options, "", true),
		)
		r.WriteString(line)

		if !strings.HasSuffix(line, "\n") {
			r.WriteString("\n")
		}
	}
}

func (r *docRenderer) renderMap(typ *ast.MapType) {
	keyName := getTypeName(typ.Key)

	var options []string
	if values, ok := r.Values[keyName]; ok {
		options = values
	}

	r.WriteString(formatOptions(options, keyName, false) + "\n")
}

func getTagValue(tag *ast.BasicLit) (val string, required bool) {
	val = strings.TrimSpace(tag.Value)
	quote := strings.Index(val, `"`)
	comma := strings.Index(val[quote:], ",")

	if comma == -1 {
		required = true
		comma = strings.LastIndex(val, `"`)
	} else {
		comma += quote
	}

	val = val[quote+1 : comma]

	return val, required
}

func (r *docRenderer) renderItemType(name, formatted string) string {
	if ident := r.asIdentName(name); ident != "" {
		return ident
	}

	if _, ok := r.Types[name]; ok {
		return fmt.Sprintf("[%s](#%s)", formatted, strings.ToLower(name))
	}

	return formatted
}

func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.ArrayType:
		return getTypeName(t.Elt)
	case *ast.StarExpr:
		return getTypeName(t.X)
	case *ast.SelectorExpr:
		return getTypeName(t.Sel)
	case *ast.Ident:
		return t.Name
	}

	return ""
}

func formatTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.ArrayType:
		return fmt.Sprintf("[%s]", formatTypeName(t.Elt))
	case *ast.MapType:
		return fmt.Sprintf("[%s: %s]", formatTypeName(t.Key), formatTypeName(t.Value))
	case *ast.StarExpr:
		return formatTypeName(t.X)
	case *ast.SelectorExpr:
		return formatTypeName(t.Sel)
	case *ast.Ident:
		return t.String()
	default:
		log.Warnf("%s", t)
		return "INVALID_EXPR"
	}
}

func (r *docRenderer) asIdentName(typeName string) string {
	typ, ok := r.Types[typeName]
	if !ok {
		return ""
	}

	if len(typ.Decl.Specs) == 0 {
		return ""
	}

	s := typ.Decl.Specs[0]

	spec, ok := s.(*ast.TypeSpec)
	if !ok {
		return ""
	}

	id, ok := spec.Type.(*ast.Ident)
	if !ok {
		return ""
	}

	return id.Name
}

func formatDoc(s string) string {
	if s == "" {
		log.Warnf("NO DOCS")
		return "NO_DOCS :shamebells:."
	}

	doc := strings.Builder{}
	lines := strings.Split(s, "\n")

	var inYamlBlock bool

	for i, line := range lines {
		switch {
		case line == "" && i < len(lines)-1:
			doc.WriteString("\n\n")
		case line == "```yaml" && !inYamlBlock:
			inYamlBlock = true
			fallthrough
		case inYamlBlock:
			if line == "```" {
				inYamlBlock = false
			}

			doc.WriteString(line + "\n")
		case line == ".":
			continue
		default:
			doc.WriteString(strings.TrimSpace(line) + " ")
		}
	}

	return doc.String()
}

func formatOptions(o []string, kind string, inline bool) string {
	if len(o) == 0 {
		return ""
	}

	if kind == "" {
		kind = "option"
	}

	var options string

	if inline {
		options = fmt.Sprintf(" `%s`.", strings.Join(o, "`, `"))
	} else {
		options = fmt.Sprintf("\n\n- `%s`", strings.Join(o, "`\n- `"))
	}

	return fmt.Sprintf(" Valid %ss:%s", kind, options)
}
