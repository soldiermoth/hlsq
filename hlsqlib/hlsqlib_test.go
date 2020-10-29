package hlsqlib

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestParseLine(t *testing.T) {
	testcases := map[string]Tag{
		"#EXTM3U":          {Name: "#EXTM3U"},
		"#EXT-X-VERSION:6": {Name: "#EXT-X-VERSION", Attrs: []Attr{{Key: "6"}}},
		`#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="a1",NAME="English",LANGUAGE="en-US",AUTOSELECT=YES,DEFAULT=YES,CHANNELS="2",URI="a1/prog_index.m3u8"`: {
			Name: "#EXT-X-MEDIA",
			Attrs: []Attr{
				{"TYPE", AttrEnum("AUDIO")},
				{"GROUP-ID", AttrString("a1")},
				{"NAME", AttrString("English")},
				{"LANGUAGE", AttrString("en-US")},
				{"AUTOSELECT", AttrBool(true)},
				{"DEFAULT", AttrBool(true)},
				{"CHANNELS", AttrString("2")},
				{"URI", AttrString("a1/prog_index.m3u8")},
			},
		},
		`#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=1966314,BANDWIDTH=2164328,CODECS="hvc1.2.4.L123.B0,mp4a.40.2",RESOLUTION=960x540,FRAME-RATE=60.000000,CLOSED-CAPTIONS="cc",AUDIO="a1",SUBTITLES="sub1"`: {
			Name: "#EXT-X-STREAM-INF",
			Attrs: []Attr{
				{"AVERAGE-BANDWIDTH", AttrInt(1966314)},
				{"BANDWIDTH", AttrInt(2164328)},
				{"CODECS", AttrString("hvc1.2.4.L123.B0,mp4a.40.2")},
				{"RESOLUTION", AttrEnum("960x540")},
				{"FRAME-RATE", AttrFloat(60)},
				{"CLOSED-CAPTIONS", AttrString("cc")},
				{"AUDIO", AttrString("a1")},
				{"SUBTITLES", AttrString("sub1")},
			},
		},
	}
	for in, e := range testcases {
		f, err := ParseTag(in)
		if err != nil {
			t.Fatalf("could not parse line=%q err=%q", in, err)
		}
		assertEqualTag(t, e, f)
	}
}

func TestScanner(t *testing.T) {
	const sample = `#EXTM3U
#EXT-X-VERSION:6
#EXT-X-INDEPENDENT-SEGMENTS

#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=2190673,BANDWIDTH=2523597,CODECS="avc1.640020,mp4a.40.2",RESOLUTION=960x540,FRAME-RATE=60.000,CLOSED-CAPTIONS="cc",AUDIO="a1",SUBTITLES="sub1"
v5/prog_index.m3u8

`

	scan := NewScanner(strings.NewReader(sample))
	for i, r := range []Tag{
		{Name: "#EXTM3U"},
		{Name: "#EXT-X-VERSION", Attrs: []Attr{{Key: "6"}}},
		{Name: "#EXT-X-INDEPENDENT-SEGMENTS", Trailing: []Raw{""}},
		{
			Name: "#EXT-X-STREAM-INF",
			Attrs: []Attr{
				{"AVERAGE-BANDWIDTH", AttrInt(2190673)},
				{"BANDWIDTH", AttrInt(2523597)},
				{"CODECS", AttrString("avc1.640020,mp4a.40.2")},
				{"RESOLUTION", AttrEnum("960x540")},
				{"FRAME-RATE", AttrFloat(60)},
				{"CLOSED-CAPTIONS", AttrString("cc")},
				{"AUDIO", AttrString("a1")},
				{"SUBTITLES", AttrString("sub1")},
			},
			Trailing: []Raw{"v5/prog_index.m3u8", ""},
		},
	} {
		if !scan.Scan() {
			log.Fatalf("Expected token @ row #%d", i+1)
		}
		assertEqualTag(t, r, scan.Tag())
	}
	if scan.Scan() {
		log.Fatalf("There should be no tokens left instead got %q", scan.Tag().Name)
	}
}

func assertEqualTag(t *testing.T, e, f Tag) {
	t.Helper()
	if reflect.DeepEqual(e, f) {
		return
	}
	if e.Name != f.Name {
		t.Fatalf("expected tag %q got %q", e.Name, f.Name)
	}
	if !reflect.DeepEqual(e.Attrs, f.Attrs) {
		msg := fmt.Sprintf("Got different attributes on tag %q", e.Name)
		attrCount := len(e.Attrs)
		if len(f.Attrs) > attrCount {
			attrCount = len(f.Attrs)
		}
		for i := 0; i < attrCount; i++ {
			var eA, fA Attr
			if i < len(e.Attrs) {
				eA = e.Attrs[i]
			}
			if i < len(f.Attrs) {
				fA = f.Attrs[i]
			}
			if !reflect.DeepEqual(eA, fA) {
				msg += fmt.Sprintf("\n\texpected %#v got %#v @ position %d", eA, fA, i)
			}
		}
		log.Fatal(msg)
	}
	if !reflect.DeepEqual(e.Trailing, f.Trailing) {
		log.Fatalf("Got different trailing lines on tag %q expected %#v got %#v", e.Name, e.Trailing, f.Trailing)
	}
	log.Fatalf(":shrug: \nExpected:\t%#v\nGot:\t\t%#v", e, f)
}
