package config

import (
	"flag"
	"os"
)

// holds common CLI flags for ingestion commands
type Flags struct {
	Path  string
	Clear bool
}

// parses CLI flags for the docs subcommand
func ParseDocsFlags() Flags {
	args := os.Args[2:]

	fs := flag.NewFlagSet("docs", flag.ExitOnError)
	path := fs.String("path", "./docs/strudel", "path to documentation directory")
	clear := fs.Bool("clear", false, "clear existing chunks before ingesting")
	fs.Parse(args)

	return Flags{Path: *path, Clear: *clear}
}

// parses CLI flags for the concepts subcommand
func ParseConceptsFlags() Flags {
	args := os.Args[2:]

	fs := flag.NewFlagSet("concepts", flag.ExitOnError)
	path := fs.String("path", "./docs/concepts", "path to concepts directory")
	clear := fs.Bool("clear", false, "clear existing concepts before ingesting")
	fs.Parse(args)

	return Flags{Path: *path, Clear: *clear}
}

// parses CLI flags for the examples subcommand
func ParseExamplesFlags() Flags {
	args := os.Args[2:]

	fs := flag.NewFlagSet("examples", flag.ExitOnError)
	path := fs.String("path", "./resources/strudel_examples.json", "path to examples JSON file")
	clear := fs.Bool("clear", false, "clear existing examples before ingesting")
	fs.Parse(args)

	return Flags{Path: *path, Clear: *clear}
}

// returns default flags for docs ingestion (used by "all" command)
func DefaultDocsFlags() Flags {
	return Flags{Path: "./docs/strudel", Clear: false}
}

// returns default flags for concepts ingestion (used by "all" command)
func DefaultConceptsFlags() Flags {
	return Flags{Path: "./docs/concepts", Clear: false}
}

// returns default flags for examples ingestion (used by "all" command)
func DefaultExamplesFlags() Flags {
	return Flags{Path: "./resources/strudel_examples.json", Clear: false}
}
