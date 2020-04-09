package hlsqlib

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	testcases := map[string]Line{
		"#EXTM3U":            Tag{Name: "#EXTM3U"},
		"http://example.com": Raw("http://example.com"),
		"#EXT-X-VERSION:6":   Tag{Name: "#EXT-X-VERSION", Attrs: []Attr{{Key: "6"}}},
		`#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="a1",NAME="English",LANGUAGE="en-US",AUTOSELECT=YES,DEFAULT=YES,CHANNELS="2",URI="a1/prog_index.m3u8"`: Tag{
			Name: "#EXT-X-MEDIA",
			Attrs: []Attr{
				{Key: "TYPE", S: "AUDIO"},
				{Key: "GROUP-ID", Q: "a1"},
				{Key: "NAME", Q: "English"},
				{Key: "LANGUAGE", Q: "en-US"},
				{Key: "AUTOSELECT", B: true},
				{Key: "DEFAULT", B: true},
				{Key: "CHANNELS", Q: "2"},
				{Key: "URI", Q: "a1/prog_index.m3u8"},
			},
		},
	}
	for in, expected := range testcases {
		found, err := ParseLine(in)
		if err != nil {
			t.Fatalf("could not parse line=%q err=%q", in, err)
		}
		if !reflect.DeepEqual(found, expected) {
			t.Fatalf("did not equal expected in=%q\nEXPECTED:\n%#v\nFOUND:\n%#v", in, expected, found)
		}
	}
}
