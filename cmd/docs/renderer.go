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

type stringWriterTo interface {
	io.StringWriter
	io.WriterTo
}

type docRenderer struct {
	Package *doc.Package
	Types   map[string]*doc.Type
	Values  map[string][]string
	Visited map[string]bool
	Queue   []*doc.Type

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

	r.Types = make(map[string]*doc.Type)
	r.Values = make(map[string][]string)
	for _, typ := range pkg.Types {
		r.Types[typ.Name] = typ

		if len(typ.Consts) > 0 {
			for _, cons := range typ.Consts {
				name, values := gatherConsts(cons)
				r.Values[name] = values
			}
		}
	}

	for _, cons := range pkg.Consts {
		name, values := gatherConsts(cons)
		r.Values[name] = values
	}

	r.Visited = make(map[string]bool)
	r.Queue = make([]*doc.Type, 0)
	r.buffer = new(bytes.Buffer)

	return r, nil
}

func gatherConsts(cons *doc.Value) (name string, values []string) {
	values = make([]string, len(cons.Decl.Specs))
	for i, s := range cons.Decl.Specs {
		if spec, ok := s.(*ast.ValueSpec); ok {
			name = spec.Type.(*ast.Ident).Name
			values[i] = spec.Values[0].(*ast.BasicLit).Value
		}
	}
	return name, values
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

func (r *docRenderer) WriteString(s string) {
	_, err := r.buffer.WriteString(s)
	if err != nil {
		log.Error(err.Error())
	}
}

func (r *docRenderer) Render(w io.Writer) error {
	projectType, ok := r.Types["Project"]
	if !ok ||
		projectType == nil ||
		len(projectType.Decl.Specs) == 0 {
		return errors.New("config.Project not found")
	}

	if r.Header != nil {
		r.WriteString(r.Header())
	}

	r.WriteString(r.Package.Doc)
	r.WriteString(docsConfigTerminologyDisclaimer)
	r.WriteString(docsConfigTableOfContents)
	r.WriteString("## Specification\n\n")

	r.Queue = append(r.Queue, projectType)

	for {
		typ := r.dequeueType()
		if typ == nil {
			break
		}
		r.renderDecl(typ.Decl, formatDoc(typ.Doc))
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

func (r *docRenderer) renderDecl(decl *ast.GenDecl, doc string) {
	for _, spec := range decl.Specs {
		r.renderSpec(spec, doc)
	}
}

func (r *docRenderer) renderSpec(s ast.Spec, doc string) {
	if spec, ok := s.(*ast.TypeSpec); ok {
		r.renderType(spec.Name.Name, doc, spec.Type)
	}
}

func (r *docRenderer) renderType(name string, doc string, expr ast.Expr) {
	switch t := expr.(type) {
	case *ast.ArrayType:
		r.enqueueTypeFromExprs(t.Elt)
	case *ast.StructType:
		r.renderTypePreamble(name, doc, 3)
		r.renderStruct(t)
	case *ast.MapType:
		r.renderTypePreamble(name, doc, 3)
		r.renderMap(t)
	case *ast.StarExpr:
		r.enqueueTypeFromExprs(t.X)
	case *ast.SelectorExpr:
		r.enqueueTypeFromExprs(t.Sel)
	case *ast.Ident:
		r.enqueueTypeFromExprs(t)
	default:
		log.Warnf("%s", t)
	}
	r.WriteString("\n")
}

func (r *docRenderer) renderTypePreamble(name string, doc string, level int) {
	r.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", level), name))
	if doc != "" {
		r.WriteString(doc + "\n\n")
	}
}

func (r *docRenderer) renderStruct(typ *ast.StructType) {
	for _, field := range typ.Fields.List {
		typeName := formatTypeName(field.Type)
		tag, required := getTagValue(field.Tag)
		var requiredStr = " "
		if required && !strings.HasPrefix(typeName, "[") {
			requiredStr = "x"
		}
		var options []string
		if values, ok := r.Values[typeName]; ok {
			typeName = "string"
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
			r.renderTypeLink(typeName, false),
			formatDoc(field.Doc.Text()),
			formatOptions(options, "", true),
		)
		r.WriteString(line)
		if !strings.HasSuffix(line, "\n") {
			r.WriteString("\n")
		}
		r.enqueueTypeFromExprs(field.Type)
	}
}

func (r *docRenderer) renderMap(typ *ast.MapType) {
	keyName := getTypeName(typ.Key)
	var options []string
	if values, ok := r.Values[keyName]; ok {
		options = values
	}
	r.WriteString(formatOptions(options, keyName, false) + "\n")
	r.enqueueTypeFromExprs(typ.Key, typ.Value)
}

func (r *docRenderer) hasType(name string) bool {
	_, ok := r.Types[name]
	return ok
}

func (r *docRenderer) enqueueTypeFromExprs(exprs ...ast.Expr) {
	for _, expr := range exprs {
		name := getTypeName(expr)
		if name == "" {
			continue
		}
		typ, ok := r.Types[name]
		if !ok {
			continue
		} else if r.Visited[name] {
			continue
		}
		r.Visited[name] = true
		r.Queue = append(r.Queue, typ)
	}
}

func (r *docRenderer) dequeueType() *doc.Type {
	if len(r.Queue) == 0 {
		return nil
	}
	typ := r.Queue[0]
	r.Queue[0] = nil
	r.Queue = r.Queue[1:]
	return typ
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

func (r *docRenderer) renderTypeLink(name string, plural bool) string {
	if r.hasType(name) {
		var pluralStr string
		if plural {
			pluralStr = "s"
		}
		return fmt.Sprintf("[%s%s](#%s)", name, pluralStr, strings.ToLower(name))
	}
	return name
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

func formatDoc(s string) string {
	if s == "" {
		log.Warnf("NO DOCS")
		return "NO_DOCS :shamebells:."
	}
	doc := strings.Builder{}
	lines := strings.Split(s, "\n")
	var inBlock bool
	for i, line := range lines {
		switch {
		case line == "" && i < len(lines)-1:
			doc.WriteString("\n\n")
		case line == "```yaml" && !inBlock:
			// TODO: make work for other block types
			inBlock = true
			fallthrough
		case inBlock:
			if line == "```" {
				inBlock = false
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
