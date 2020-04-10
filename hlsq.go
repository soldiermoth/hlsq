package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/soldiermoth/hlsq/hlsqlib"
)

func main() {
	var (
		colors = hlsqlib.ColorSettings{
			Tag:  hlsqlib.LightBlue,
			Attr: hlsqlib.Cyan,
		}
		chomp = flag.Bool("chomp", false, "Removes whitespace")
		r     io.Reader
		err   error
	)
	flag.Var(flagColor{&colors.Tag}, "color-tag", "Tag")
	flag.Var(flagColor{&colors.Attr}, "color-attr", "Attr")
	flag.Parse()
	args := flag.Args()
	// Detect if stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		r = os.Stdin
		if len(args) > 0 {
			log.Fatal("no arguments expected when reading from stdin")
		}
	} else if len(args) == 1 {
		if r, err = os.Open(args[0]); err != nil {
			log.Fatalf("could not open file %q err=%q", args[0], err)
		}
	} else {
		log.Fatal("expected 1 argument of the m3u8 file to process or to read from stdin")
	}

	var (
		scanner = bufio.NewScanner(r)
		opts    = []hlsqlib.SerializeOption{
			hlsqlib.ColorLines(colors),
		}
	)
	if *chomp {
		opts = append(opts, hlsqlib.Chomp)
	}
	for scanner.Scan() {
		line, err := hlsqlib.ParseLine(scanner.Text())
		if err != nil {
			log.Fatalf("error parsing line=%q %q", scanner.Text(), err)
		}

		hlsqlib.Serialize(os.Stdout, line, opts...)
	}
	fmt.Fprintln(os.Stdout)
}

type flagColor struct{ *hlsqlib.Color }

func (f flagColor) String() string {
	if f.Color == nil {
		return "Example"
	}
	return hlsqlib.NewColorizer(*f.Color).S("Example")
}
func (f flagColor) Set(s string) (err error) {
	if s != "" {
		*f.Color, err = hlsqlib.ParseColor(s)
	}
	return
}
