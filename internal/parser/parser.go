package parser

import (
	"errors"
	"fmt"
	"go/parser"

	"github.com/mpanelo/go-tagger/internal/config"
	"golang.org/x/tools/go/buildutil"
)

// Parse is a wrapper over the ParseFile function from the package go/parser. If the modified is not nil,
// then we parse the archive and look for the provided filename. If it's found, then ParseFile will parse
// the archive, otherwise it will read from the file named filename.
func Parse(cfg *config.Config) (err error) {
	if cfg.Filename == "" {
		return errors.New("-file cannot be empty")
	}

	var contents any

	if cfg.Modified != nil {
		archive, err := buildutil.ParseOverlayArchive(cfg.Modified)
		if err != nil {
			return fmt.Errorf("failed to parse -modified archive: %w", err)
		}
		fc, ok := archive[cfg.Filename]
		if !ok {
			return fmt.Errorf("file %q not found in the -modified archive", cfg.Filename)
		}
		contents = fc
	}

	cfg.File, err = parser.ParseFile(cfg.Fset, cfg.Filename, contents, parser.ParseComments)
	return
}
