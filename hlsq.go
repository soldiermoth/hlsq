package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/soldiermoth/hlsq/hlsqlib"
	"golang.org/x/text/unicode/norm"
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
		watch       = flag.Bool("watch", false, "Continuously watch the playlist (requires url)")
		url         = flag.String("url", "", "URL of the manifest to watch (required to watch)")
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
	} else if *watch && *url != "" {
		r = getManifest(*url)
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
	if *watch && *url != "" {
		opts = append(opts, hlsqlib.Chomp)
		var updatePeriod time.Duration
		seen := make(map[string]bool)
		ignoredTagsOnUpdate := []string{
			"#EXTM3U",
			"#EXT-X-VERSION",
			"#EXT-X-DISCONTINUITY-SEQUENCE",
			"#EXT-X-PUBLISHED-TIME",
			"#EXT-X-TARGETDURATION",
			"#EXT-X-MEDIA-SEQUENCE",
			"#EXT-X-PROGRAM-DATE-TIME",
			"#EXT-X-KEY",
		}
		for scanner.Scan() {
			line := scanner.Tag()
			if line.Name == "#EXT-X-TARGETDURATION" {
				d, _ := strconv.Atoi(line.Attrs[0].Key)
				updatePeriod = time.Duration(d) * time.Second
			}
			if line.Name == "#EXTINF" {
				segment := norm.NFC.Bytes([]byte(string(line.Trailing[0])))
				seen[string(segment)] = true
			}
			hlsqlib.Serialize(os.Stdout, line, opts...)
		}
		for range time.Tick(updatePeriod) {
			r = getManifest(*url)
			scanner = hlsqlib.NewScanner(r)
			updated := false
			for scanner.Scan() {
				line := scanner.Tag()
				if stringInSlice(line.Name, ignoredTagsOnUpdate) {
					continue
				}
				segment := norm.NFC.Bytes([]byte(string(line.Trailing[0])))
				if line.Name == "#EXTINF" && seen[string(segment)] {
					continue
				} else if line.Name == "#EXTINF" && !seen[string(segment)] {
					seen[string(segment)] = true
				}
				hlsqlib.Serialize(os.Stdout, line, opts...)
				updated = true
			}
			if !updated {
				fmt.Println("manifest did not change this period")
			}
		}
	} else {
		for scanner.Scan() {
			line := scanner.Tag()
			hlsqlib.Serialize(os.Stdout, line, opts...)
		}
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

func getManifest(url string) io.ReadCloser {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("could not fetch manifest: %v\n", err)
	}
	return resp.Body
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
