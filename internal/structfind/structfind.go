package structfind

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/mpanelo/go-tagger/internal/config"
	"github.com/mpanelo/go-tagger/internal/parser"
)

type structType struct {
	name string
	node *ast.StructType
}

func Find(cfg *config.Config, pr *parser.ParseResult) (token.Pos, token.Pos, error) {
	if cfg.Lines != "" {
		return lineSelection(cfg, pr)
	}

	if cfg.Offset != 0 {
		return offsetSelection(cfg, pr)
	}

	if cfg.StructName != "" {
		return structSelection(cfg, pr)
	}

	return 0, 0, fmt.Errorf("-line, -offset, or -struct must be provided")
}

func lineSelection(cfg *config.Config, pr *parser.ParseResult) (token.Pos, token.Pos, error) {
	// TODO: Check there are structs within the lines or if line is within a struct definition

	tokens := strings.Split(cfg.Lines, ",")

	if len(tokens) > 2 {
		return 0, 0, fmt.Errorf("%d items provided, expected at most 2", len(tokens))
	}

	var start, end int
	var err error

	start, err = strconv.Atoi(tokens[0])
	if err != nil {
		return 0, 0, err
	}

	end = start
	if len(tokens) == 2 {
		end, err = strconv.Atoi(tokens[1])
		if err != nil {
			return 0, 0, err
		}
	}

	if start > end {
		return 0, 0, fmt.Errorf("start: %d, end: %d is not a valid range", start, end)
	}

	f := pr.Fset.File(pr.File.FileStart)
	startPos := f.LineStart(start)
	lineCount := f.LineCount()
	if start < 0 || start > lineCount || end > lineCount {
		return 0, 0, fmt.Errorf("outside file line range [%d, %d], got [%d, %d]", f.Base(), lineCount, start, end)
	}
	var endPos token.Pos
	if end == lineCount {
		endPos = f.Pos(f.Size())
	} else {
		endPos = f.LineStart(end+1) - 1
	}

	return startPos, endPos, nil
}

func offsetSelection(cfg *config.Config, parserResult *parser.ParseResult) (token.Pos, token.Pos, error) {
	structs := collectStructs(parserResult.File)

	for _, st := range structs {
		start := parserResult.Fset.Position(st.node.Pos())
		end := parserResult.Fset.Position(st.node.End())

		if start.Offset < cfg.Offset && cfg.Offset < end.Offset {
			return st.node.Pos(), st.node.End(), nil
		}
	}

	return 0, 0, fmt.Errorf("offset %d is not within a struct", cfg.Offset)
}

func structSelection(cfg *config.Config, parserResult *parser.ParseResult) (token.Pos, token.Pos, error) {
	structs := collectStructs(parserResult.File)

	for _, st := range structs {
		if st.name == cfg.StructName {
			return st.node.Pos(), st.node.End(), nil
		}
	}

	return 0, 0, fmt.Errorf("struct %s was not found", cfg.StructName)
}

func collectStructs(file *ast.File) map[token.Pos]*structType {
	structs := make(map[token.Pos]*structType)

	getAllStructs := func(n ast.Node) bool {
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

	ast.Inspect(file, getAllStructs)
	return structs
}
