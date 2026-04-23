package structfind

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type StructFinder struct {
	FileSet    *token.FileSet
	Line       string
	Offset     int
	StructName string
}

type structType struct {
	name string
	node *ast.StructType
}

var ErrInvalidLineSelection = errors.New("invalid line selection")

func newErr(errorMessage string) error {
	return fmt.Errorf("%w: %s", ErrInvalidLineSelection, errorMessage)
}

func (sf *StructFinder) LineSelection() (int, int, error) {
	if sf.Line == "" {
		return 0, 0, newErr("Line cannot be empty")
	}

	tokens := strings.Split(sf.Line, ",")

	if len(tokens) > 2 {
		return 0, 0, newErr(fmt.Sprintf("%d items provided, expected at most 2", len(tokens)))
	}

	var start, end int
	var err error

	start, err = strconv.Atoi(tokens[0])
	if err != nil {
		return 0, 0, newErr(err.Error())
	}

	end = start
	if len(tokens) == 2 {
		end, err = strconv.Atoi(tokens[1])
		if err != nil {
			return 0, 0, newErr(err.Error())
		}
	}

	if start > end {
		return 0, 0, newErr(fmt.Sprintf("start: %d, end: %d is not a valid range", start, end))
	}

	return start, end, nil
}

func (sf *StructFinder) OffsetSelection(node ast.Node) (int, int, error) {
	structs := sf.collectStructs(node)

	for _, st := range structs {
		start := sf.FileSet.Position(st.node.Pos())
		end := sf.FileSet.Position(st.node.End())

		if start.Offset < sf.Offset && sf.Offset < end.Offset {
			return start.Line, end.Line, nil
		}
	}

	return 0, 0, fmt.Errorf("offset %d is not within a struct", sf.Offset)
}

func (sf *StructFinder) StructSelection(node ast.Node) (int, int, error) {
	structs := sf.collectStructs(node)

	for _, st := range structs {
		if st.name == sf.StructName {
			startLine := sf.FileSet.Position(st.node.Pos()).Line
			endLine := sf.FileSet.Position(st.node.End()).Line
			return startLine, endLine, nil
		}
	}

	return 0, 0, fmt.Errorf("struct %s was not found", sf.StructName)
}

func (sf *StructFinder) collectStructs(node ast.Node) map[token.Pos]*structType {
	structs := make(map[token.Pos]*structType)

	collectStructs := func(n ast.Node) bool {
		t, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if t.Type == nil {
			return true
		}

		structName := t.Name.Name

		x, ok := t.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structs[x.Pos()] = &structType{
			name: structName,
			node: x,
		}

		return true
	}

	ast.Inspect(node, collectStructs)
	return structs
}
