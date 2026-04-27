package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var ErrNoArgs = errors.New("no args found")

// Config defines how tags should be modified.
type Config struct {
	Filename   string
	Output     string
	Write      bool
	Modified   io.Reader
	Offset     int
	StructName string
	Lines      string
	Start, End int
	Clear      bool

	// Tag-specific logic
	RemoveTags string
	AddTags    string
	Override   bool
	Sort       bool
}

// Parse handles flag initialization and returns a populated Config pointer.
func Parse() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.Filename, "filename", "", "File to be parsed")
	flag.BoolVar(&c.Write, "w", false, "Write result to (source) file instead of stdout")

	isModified := flag.Bool("modified", false, "read an archive of modified files from standard input")

	// Processing modes
	flag.IntVar(&c.Offset, "offset", 0, "Byte offset of the cursor position inside a struct")
	flag.StringVar(&c.Lines, "lines", "", "Line number or range (e.g. 4 or 4,8)")
	flag.StringVar(&c.StructName, "struct", "", "Struct name to be processed")

	// Tag flags
	flag.StringVar(&c.RemoveTags, "remove-tags", "", "Remove tags for comma separated list of keys")
	flag.StringVar(&c.AddTags, "add-tags", "", "Adds tags for comma separated list of keys")
	flag.BoolVar(&c.Override, "override", false, "Override current tags when adding tags")
	flag.BoolVar(&c.Sort, "sort", false, "Sort tags in increasing order")

	// Custom usage to avoid cluttering stderr on simple errors
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NFlag() == 0 {
		return nil, ErrNoArgs
	}

	if *isModified {
		c.Modified = os.Stdin
	}

	return c, nil
}
