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
		r   io.Reader
		err error
	)
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
		tagC    = hlsqlib.NewColorizer(hlsqlib.LightBlue)
		attrC   = hlsqlib.NewColorizer(hlsqlib.Cyan)
	)
	for scanner.Scan() {
		line, err := hlsqlib.ParseLine(scanner.Text())
		if err != nil {
			log.Fatalf("error parsing line=%q %q", scanner.Text(), err)
		}
		if tag, ok := line.(hlsqlib.Tag); ok {
			tag.Name = tagC.S(tag.Name)
			line = tag
			for i, a := range tag.Attrs {
				tag.Attrs[i].Key = attrC.S(a.Key)
			}
		}
		hlsqlib.Serialize(os.Stdout, line)
	}
	fmt.Fprintln(os.Stdout)
}
