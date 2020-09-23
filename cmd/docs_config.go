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
	"github.com/cidertool/asc-go/asc"
	"github.com/cidertool/cider/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	docsConfigFrontmatterTemplate = `---
layout: page
nav_order: %d
---

# Configuration
{: .no_toc }

`
	docsConfigTableOfContents = `
<details open markdown="block">
  <summary>
    Table of Contents
  </summary>
  {: .text-delta }
- TOC
{:toc}
</details>

`
	docsConfigTerminologyDisclaimer = `
- [x] An X here means the field is required.
- [ ] This field is optional and can be omitted.

`
)

// nolint: gochecknoglobals
var docsConfigExampleProject = config.Project{
	"My App": {
		BundleID:              "com.myproject.MyApp",
		PrimaryLocale:         "en-US",
		UsesThirdPartyContent: asc.Bool(false),
		Availability: &config.Availability{
			AvailableInNewTerritories: asc.Bool(false),
			Pricing: []config.PriceSchedule{
				{Tier: "0"},
			},
			Territories: []string{"USA"},
		},
		Categories: &config.Categories{
			Primary:   "SOCIAL_NETWORKING",
			Secondary: "GAMES",
			SecondarySubcategories: [2]string{
				"GAMES_SIMULATION",
				"GAMES_RACING",
			},
		},
		Localizations: config.AppLocalizations{
			"en-US": {
				Name:     "My App",
				Subtitle: "Not Your App",
			},
		},
		Versions: config.Version{
			Platform:             config.PlatformiOS,
			Copyright:            "2020 Me",
			EarliestReleaseDate:  nil,
			ReleaseType:          config.ReleaseTypeAfterApproval,
			PhasedReleaseEnabled: true,
			IDFADeclaration:      nil,
			Localizations: config.VersionLocalizations{
				"en-US": {
					Description:  "My App for cool people",
					Keywords:     "Apps, Cool, Mine",
					WhatsNewText: `Thank you for using My App! I bring you updates every week so this continues to be my app.`,
					PreviewSets: config.PreviewSets{
						config.PreviewTypeiPhone65: []config.Preview{
							{
								File: config.File{
									Path: "assets/store/iphone65/preview.mp4",
								},
							},
						},
					},
					ScreenshotSets: config.ScreenshotSets{
						config.ScreenshotTypeiPhone65: []config.File{
							{Path: "assets/store/iphone65/app.jpg"},
						},
					},
				},
			},
		},
		Testflight: config.Testflight{
			EnableAutoNotify: true,
			Localizations: config.TestflightLocalizations{
				"en-US": {
					Description: "My App for cool people using the beta",
				},
			},
		},
	},
}

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
	path = filepath.Join(path, "configuration.md")

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
	r.Header = func() string {
		return fmt.Sprintf(docsConfigFrontmatterTemplate, 4)
	}
	r.Footer = func() string {
		contents, err := ioutil.ReadFile(filepath.Join(filepath.Dir(path), "configuration-footer.md"))
		if err != nil {
			log.Error(err.Error())
			return ""
		}
		return string(contents)
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
