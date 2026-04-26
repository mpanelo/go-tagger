package main

import (
	"fmt"
	"os"

	"github.com/mpanelo/go-tagger/internal/config"
	"github.com/mpanelo/go-tagger/internal/modifytags"
	"github.com/mpanelo/go-tagger/internal/parser"
	"github.com/mpanelo/go-tagger/internal/structfind"
)

func main() {
	var m Main
	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		os.Exit(1)
	}
}

type Main struct {
	cfg *config.Config
	mod *modifytags.Modification
}

func (m *Main) Run() error {
	var err error
	m.cfg = config.Parse()

	if err := parser.Parse(m.cfg); err != nil {
		return err
	}

	_, _, err = structfind.Find(m.cfg)
	if err != nil {
		return err
	}

	return nil
}
