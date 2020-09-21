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

	pkg, err := openPackage()
	if err != nil {
		return err
	}

	return renderPackage(pkg, f)
}

func openPackage() (*doc.Package, error) {
	fset := token.NewFileSet()
	files, err := parsePkgConfigFiles(fset)
	if err != nil {
		return nil, err
	}
	return doc.NewFromFiles(fset, files, "github.com/cidertool/cider/pkg/config")
}

func parsePkgConfigFiles(fset *token.FileSet) (files []*ast.File, err error) {
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

func renderPackage(pkg *doc.Package, w io.Writer) error {
	var projectType *doc.Type
	for _, typ := range pkg.Types {
		if typ.Name == "Project" {
			projectType = typ
		}
	}
	if projectType == nil || len(projectType.Decl.Specs) == 0 {
		return errors.New("config.Project not found")
	}

	buf := new(bytes.Buffer)

	buf.WriteString("# " + strings.Title(pkg.Name) + "\n\n")

	buf.WriteString(pkg.Doc + "\n")

	for _, spec := range projectType.Decl.Specs {
		docForSpec(spec, buf)
	}

	// buf.WriteString("## Consts\n\n")

	// for _, group := range pkg.Consts {
	// 	if group.Doc != "" {
	// 		buf.WriteString(group.Doc + "\n")
	// 	}
	// 	for _, con := range group.Names {
	// 		buf.WriteString("- " + con + "\n")
	// 	}
	// 	buf.WriteString("\n")
	// }
	// buf.WriteString("\n")

	// buf.WriteString("## Types\n\n")

	// for _, typ := range pkg.Types {
	// 	buf.WriteString("### " + typ.Name + "\n\n")
	// 	if typ.Doc != "" {
	// 		buf.WriteString(typ.Doc)
	// 	}
	// 	for _, vari := range typ.Vars {
	// 		buf.WriteString(vari.Decl.Tok.String())
	// 	}
	// 	buf.WriteString("\n")
	// }
	buf.WriteString("\n")

	_, err := buf.WriteTo(w)
	return err
}

func docForSpec(spec ast.Spec, b io.StringWriter) {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		b.WriteString(s.Name.Name)
		docForType(s.Type, b)
	case *ast.ValueSpec:
		b.WriteString("value " + s.Names[0].Name + "\n")
	case *ast.ImportSpec:
		b.WriteString("import " + s.Name.Name + "\n")
	}
}

func docForType(expr ast.Expr, b io.StringWriter) {
	switch t := expr.(type) {
	case *ast.ArrayType:
		b.WriteString(" is an array of" + fmt.Sprint(t.Elt))
		b.WriteString("\n")
	case *ast.StructType:
		b.WriteString(" is a struct\n")
		for _, field := range t.Fields.List {
			docForType(field.Type, b)
		}
	case *ast.MapType:
		b.WriteString(" is a map of ")
		docForType(t.Key, b)
		b.WriteString(" to ")
		docForType(t.Value, b)
		b.WriteString("\n")
	default:
		fmt.Println(t)
	}
}
