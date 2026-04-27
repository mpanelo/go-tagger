package modifytags

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"

	"github.com/mpanelo/go-tagger/internal/parser"
)

type Transform int

const (
	CamelCase = iota
	SnakeCase
	PascalCase
)

type Modification struct {
	Add       []string
	Remove    []string
	Transform Transform
	Override  bool
	Sort      bool
	Clear     bool
}

func (mod *Modification) Apply(pr *parser.ParseResult, start, end token.Pos) error {
	var errs []error

	rewriteRec := func(node ast.Node) bool {
		st, ok := node.(*ast.StructType)
		if !ok {
			return true
		}

		for _, f := range st.Fields.List {
			if !(start <= f.End() && f.Pos() <= end) {
				continue
			}

			fieldName := ""
			if len(f.Names) > 0 {
				fieldName = f.Names[0].Name
			}

			if f.Names == nil {
				ident, ok := f.Type.(*ast.Ident)
				if !ok {
					continue
				}
				fieldName = ident.Name
			}

			if fieldName == "" {
				continue
			}

			currTag := ""
			if f.Tag != nil {
				currTag = f.Tag.Value
			}

			res, err := mod.processField(fieldName, currTag)
			if err != nil {
				filename := pr.Fset.Position(f.Pos()).Filename
				line := pr.Fset.Position(f.Pos()).Line
				column := pr.Fset.Position(f.Pos()).Column
				errs = append(errs, fmt.Errorf("%s:%d:%d:%s", filename, line, column, err))
				continue
			}

			if res == "" {
				f.Tag = nil
			} else {
				if f.Tag == nil {
					f.Tag = &ast.BasicLit{}
				}
				f.Tag.Value = res
			}
		}

		return true
	}

	ast.Inspect(pr.File, rewriteRec)

	if len(errs) > 0 {
		return &RewriteErrors{Errs: errs}
	}
	return nil
}

func (mod *Modification) processField(fieldName, tagVal string) (string, error) {
	return "", nil
}

// RewriteErrors are errors that occurred while rewriting struct field tags.
type RewriteErrors struct {
	Errs []error
}

func (r *RewriteErrors) Error() string {
	var buf bytes.Buffer
	for _, e := range r.Errs {
		fmt.Fprintf(&buf, "%s\n", e.Error())
	}
	return buf.String()
}
