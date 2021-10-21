package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func WriteDeprecated(dir string, comment string) error {
	return walkDir(dir, func(dir string, file fs.FileInfo) error {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") || strings.HasSuffix(file.Name(), "_test.go") {
			return nil
		}
		if strings.Contains(dir, "internal/") {
			return nil
		}
		path := dir + "/" + file.Name()
		return writeDeprecated(file.Name(), path, comment)
	})
}

func walkDir(dir string, f func(dir string, file fs.FileInfo) error) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = f(dir, file); err != nil {
			return err
		}
		if file.IsDir() {
			if err = walkDir(dir+"/"+file.Name(), f); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeDeprecated(name, path, deprecated string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, name, bs, parser.ParseComments)
	if err != nil {
		return err
	}

	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		f, ok := cursor.Node().(*ast.FuncDecl)
		if ok && f.Name.IsExported() {
			if f.Doc == nil {
				f.Doc = &ast.CommentGroup{List: []*ast.Comment{
					{Text: fmt.Sprintf("// %s ...", f.Name.Name)},
					{Text: "//"},
					{Text: "// Deprecated: " + deprecated},
				}}
			} else {
				lastText := f.Doc.List[len(f.Doc.List)-1].Text
				if lastText != "" {
					f.Doc.List = append(f.Doc.List, &ast.Comment{Text: "//"})
				}
				f.Doc.List = append(f.Doc.List, &ast.Comment{Text: "// Deprecated: " + deprecated})
			}
		}
		return true
	}, nil)

	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	return printer.Fprint(f, fset, file)
}
