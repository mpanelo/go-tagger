package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"golang.org/x/tools/go/buildutil"
)

// Parse is a wrapper over the ParseFile function from the package go/parser. If the modified is not nil,
// then we parse the archive and look for the provided filename. If it's found, then ParseFile will parse
// the archive, otherwise it will read from the file named filename.
func Parse(filename string, modified io.Reader) (*token.FileSet, *ast.File, error) {
	var contents any

	if modified != nil {
		archive, err := buildutil.ParseOverlayArchive(modified)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse -modified archive: %w", err)
		}
		fc, ok := archive[filename]
		if !ok {
			return nil, nil, fmt.Errorf("file %q not found in the -modified archive", filename)
		}
		contents = fc
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filename, contents, parser.ParseComments)

	return fset, file, err
}
