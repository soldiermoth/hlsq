package main

import (
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
		tagSpecific []hlsqlib.SpecificTagColors
		query       = flag.String("query", "", "Query")
		chomp       = flag.Bool("chomp", false, "Removes whitespace")
		demuxed     = flag.Bool("demuxed", false, "Set to demuxed colors")
		r           io.Reader
		err         error
	)
	flag.Var(flagColor{&colors.Tag}, "color-tag", "Tag")
	flag.Var(flagColor{&colors.Attr}, "color-attr", "Attr")
	flag.Parse()
	if *demuxed {
		colors.Tag = hlsqlib.White
		colors.Attr = hlsqlib.White
		colors.Raw = hlsqlib.Yellow
		tagSpecific = append(tagSpecific, func(s string) (hlsqlib.Colorizer, bool) {
			switch s {
			case "#EXT-X-DISCONTINUITY":
				return hlsqlib.NewColorizer(hlsqlib.Red), true
			}
			return "", false
		})
	}
	args := flag.Args()
	// Detect if stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		r = os.Stdin
		if len(args) == 1 {
			*query = args[0]
		} else if len(args) > 0 {
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
		scanner = hlsqlib.NewScanner(r)
		opts    = []hlsqlib.SerializeOption{
			hlsqlib.ColorLines(colors, tagSpecific...),
		}
	)
	if *chomp {
		opts = append(opts, hlsqlib.Chomp)
	}
	if query != nil && *query != "" {
		queryFunc, err := hlsqlib.ParseQuery(*query)
		if err != nil {
			log.Fatalf("could not parse expression %q %q", *query, err)
		}
		opts = append([]hlsqlib.SerializeOption{hlsqlib.AttrMatch(queryFunc)}, opts...)
	}
	for scanner.Scan() {
		line := scanner.Tag()
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
