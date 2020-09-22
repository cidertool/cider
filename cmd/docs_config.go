package cmd

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
	"os"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func newDocsConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Generate configuration file documentation for Cider.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runDocsConfigCmd,
	}
}

func runDocsConfigCmd(cmd *cobra.Command, args []string) error {
	var path string
	if len(args) == 0 {
		path = defaultDocsPath
	} else {
		path = args[0]
	}
	path = filepath.Join(path, "configuration2.md")

	log.WithField("path", path).Info("generating configuration documentation")
	err := genConfigMarkdown(path)
	if err != nil {
		log.Error("generation failed")
	} else {
		log.Info("generation completed successfully")
	}
	return err
}

func genConfigMarkdown(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := f.Close()
		if closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				log.Fatal(closeErr.Error())
			}
		}
	}()

	r, err := newRenderer()
	if err != nil {
		return err
	}

	return r.Render(f)
}

type stringWriterTo interface {
	io.StringWriter
	io.WriterTo
}

type docRenderer struct {
	Package *doc.Package
	Types   map[string]*doc.Type
	Values  map[string]*doc.Value
	Visited map[string]bool
	Queue   []*doc.Type
	buffer  stringWriterTo
}

func newRenderer() (*docRenderer, error) {
	r := new(docRenderer)

	pkg, err := openPackage()
	if err != nil {
		return nil, err
	}
	r.Package = pkg

	r.Types = make(map[string]*doc.Type)
	for _, typ := range pkg.Types {
		r.Types[typ.Name] = typ
	}

	r.Values = make(map[string]*doc.Value)
	r.Visited = make(map[string]bool)
	r.Queue = make([]*doc.Type, 0)
	r.buffer = new(bytes.Buffer)

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
	r.WriteString("# " + strings.Title(r.Package.Name) + "\n\n")
	r.WriteString(r.Package.Doc + "\n")

	r.Queue = append(r.Queue, projectType)

	for {
		typ := r.dequeueType()
		if typ == nil {
			break
		}
		r.renderDecl(typ.Decl, formatDoc(typ.Doc))
	}

	_, err := r.buffer.WriteTo(w)
	return err
}

func (r *docRenderer) renderDecl(decl *ast.GenDecl, doc string) {
	for _, spec := range decl.Specs {
		r.renderSpec(spec, doc)
	}
}

func (r *docRenderer) renderSpec(s ast.Spec, doc string) {
	switch spec := s.(type) {
	case *ast.TypeSpec:
		r.renderType(spec.Name.Name, doc, spec.Type)
	case *ast.ValueSpec:
		fmt.Println(spec)
	}
}

func (r *docRenderer) renderType(name string, doc string, expr ast.Expr) {
	switch t := expr.(type) {
	case *ast.ArrayType:
		r.enqueueTypeFromExprs(t.Elt)
	case *ast.StructType:
		r.WriteString("### " + name + "\n\n")
		if doc != "" {
			r.WriteString(doc + "\n\n")
		}
		for _, field := range t.Fields.List {
			typeName := formatTypeName(field.Type)
			tag, required := getTagValue(field.Tag)
			var requiredStr = " "
			if required && !strings.HasPrefix(typeName, "[") {
				requiredStr = "x"
			}
			if len(field.Names) != 0 {
				line := fmt.Sprintf(
					"- [%s] **%s: %s** - %s",
					requiredStr,
					tag,
					r.renderTypeLink(typeName, false),
					formatDoc(field.Doc.Text()),
				)
				r.WriteString(line)
				if !strings.HasSuffix(line, "\n") {
					r.WriteString("\n")
				}
			}
			r.enqueueTypeFromExprs(field.Type)
		}
	case *ast.MapType:
		r.enqueueTypeFromExprs(t.Key, t.Value)
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

func (r *docRenderer) hasType(name string) bool {
	_, ok := r.Types[name]
	return ok
}

func (r *docRenderer) hasValue(name string) bool {
	_, ok := r.Values[name]
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
	doc := strings.Builder{}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line == "" && i < len(lines)-1 {
			doc.WriteString("\n\n")
		} else {
			doc.WriteString(strings.TrimSpace(line) + " ")
		}
	}
	return doc.String()
}
