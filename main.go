package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/mpanelo/custom-go-tool/internal/parser"
	"github.com/mpanelo/custom-go-tool/internal/structfind"
)

// config defines how tags should be modified
type config struct {
	file     string
	output   string
	write    bool
	modified io.Reader // modified file content without saving to the file system

	offset     int
	structName string
	line       string
	start, end int

	remove        []string
	removeOptions []string

	add        []string
	addOptions []string
	override   bool

	transform   string
	sort        bool
	clear       bool
	clearOption bool
}

func main() {
	var m Main
	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		os.Exit(1)
	}
}

type Main struct {
	fset         *token.FileSet
	cfg          *config
	StructFinder *structfind.StructFinder
}

func (m *Main) Run() error {
	var err error
	m.cfg, err = m.getConfig()
	if err != nil {
		return err
	}
	// When only the usage is printed out, cfg will be nil. Instead of calling os.Exit(0)
	// just return early and exit main.
	if m.cfg == nil {
		return nil
	}

	var file *ast.File
	m.fset, file, err = parser.Parse(m.cfg.file, m.cfg.modified)
	if err != nil {
		return err
	}

	_, _, err = m.findSelection(file)
	if err != nil {
		return err
	}

	return nil
}

func (m *Main) findSelection(file *ast.File) (int, int, error) {
	if m.cfg.line != "" {
		return m.StructFinder.LineSelection()
	}

	if m.cfg.offset != 0 {
		return m.StructFinder.OffsetSelection(file)
	}

	if m.cfg.structName != "" {
		return m.StructFinder.StructSelection(file)
	}

	return 0, 0, fmt.Errorf("-line, -offset, or -struct must be provided")
}

func (m *Main) getConfig() (*config, error) {
	var (
		// file flags
		flagFile  = flag.String("file", "", "Filename to be parsed")
		flagWrite = flag.Bool("w", false,
			"Write result to (source) file instead of stdout")
		flagOutput = flag.String("format", "source", "Output format."+
			"By default it's the whole file. Options: [source, json]")
		flagModified = flag.Bool("modified", false, "read an archive of modified files from standard input")

		// processing modes
		flagOffset = flag.Int("offset", 0,
			"Byte offset of the cursor position inside a struct."+
				"Can be anwhere from the comment until closing bracket")
		flagLine = flag.String("line", "",
			"Line number of the field or a range of line. i.e: 4 or 4,8")
		flagStruct = flag.String("struct", "", "Struct name to be processed")

		// tag flags
		flagRemoveTags = flag.String("remove-tags", "",
			"Remove tags for the comma separated list of keys")
		flagClearTags = flag.Bool("clear-tags", false,
			"Clear all tags")
		flagAddTags = flag.String("add-tags", "",
			"Adds tags for the comma separated list of keys."+
				"Keys can contain a static value, i,e: json:foo")
		flagOverride  = flag.Bool("override", false, "Override current tags when adding tags")
		flagTransform = flag.String("transform", "snakecase",
			"Transform adds a transform rule when adding tags."+
				" Current options: [snakecase, camelcase, lispcase]")
		flagSort = flag.Bool("sort", false,
			"Sort sorts the tags in increasing order according to the key name")

		// option flags
		flagRemoveOptions = flag.String("remove-options", "",
			"Remove the comma separated list of options from the given keys, "+
				"i.e: json=omitempty,hcl=squash")
		flagClearOptions = flag.Bool("clear-options", false,
			"Clear all tag options")
		flagAddOptions = flag.String("add-options", "",
			"Add the options per given key. i.e: json=omitempty,hcl=squash")
	)

	// don't output full help information if something goes wrong
	flag.Usage = func() {}
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return nil, nil
	}

	cfg := &config{
		file:        *flagFile,
		line:        *flagLine,
		structName:  *flagStruct,
		offset:      *flagOffset,
		output:      *flagOutput,
		write:       *flagWrite,
		clear:       *flagClearTags,
		clearOption: *flagClearOptions,
		transform:   *flagTransform,
		sort:        *flagSort,
		override:    *flagOverride,
	}

	if *flagModified {
		cfg.modified = os.Stdin
	}

	if *flagAddTags != "" {
		cfg.add = strings.Split(*flagAddTags, ",")
	}

	if *flagAddOptions != "" {
		cfg.addOptions = strings.Split(*flagAddOptions, ",")
	}

	if *flagRemoveTags != "" {
		cfg.remove = strings.Split(*flagRemoveTags, ",")
	}

	if *flagRemoveOptions != "" {
		cfg.removeOptions = strings.Split(*flagRemoveOptions, ",")
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate validates whether the config is valid or not
func (c *config) validate() error {
	if c.file == "" {
		return errors.New("no file is passed")
	}

	if c.line == "" && c.offset == 0 && c.structName == "" {
		return errors.New("-line, -offset or -struct is not passed")
	}

	if c.line != "" && c.offset != 0 || c.line != "" && c.structName != "" || c.offset != 0 && c.structName != "" {
		return errors.New("-line, -offset or -struct cannot be used together. pick one")
	}

	if len(c.add) == 0 && len(c.addOptions) == 0 && !c.clear && !c.clearOption && len(c.removeOptions) == 0 &&
		len(c.remove) == 0 {
		return errors.New("one of [-add-tags, -add-options, -remove-tags, -remove-options, -clear-tags," +
			" -clear-options] should be defined")
	}

	return nil
}
