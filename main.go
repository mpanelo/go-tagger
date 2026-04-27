package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mpanelo/go-tagger/internal/config"
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
	// mod *modifytags.Modification
}

func (m *Main) Run() error {
	cfg, err := config.Parse()
	if err != nil {
		if errors.Is(err, config.ErrNoArgs) {
			flag.Usage()
			os.Exit(0)
		}
	}

	pr, err := parser.Parse(cfg)
	if err != nil {
		return err
	}

	_, _, err = structfind.Find(cfg, pr)
	if err != nil {
		return err
	}

	return nil
}
