package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/mpanelo/go-tagger/internal/config"
	"golang.org/x/tools/go/buildutil"
)

type ParseResult struct {
	Fset *token.FileSet
	File *ast.File
}

// Parse is a wrapper over the ParseFile function from the package go/parser. If the modified is not nil,
// then we parse the archive and look for the provided filename. If it's found, then ParseFile will parse
// the archive, otherwise it will read from the file named filename.
func Parse(cfg *config.Config) (pr *ParseResult, err error) {
	if cfg.Filename == "" {
		err = errors.New("-file cannot be empty")
		return
	}

	var contents any

	if cfg.Modified != nil {
		var archive map[string][]byte
		archive, err = buildutil.ParseOverlayArchive(cfg.Modified)
		if err != nil {
			err = fmt.Errorf("failed to parse -modified archive: %w", err)
			return
		}

		fc, ok := archive[cfg.Filename]
		if !ok {
			err = fmt.Errorf("file %q not found in the -modified archive", cfg.Filename)
			return
		}
		contents = fc
	}

	pr = &ParseResult{
		Fset: token.NewFileSet(),
	}
	pr.File, err = parser.ParseFile(pr.Fset, cfg.Filename, contents, parser.ParseComments)
	return
}
