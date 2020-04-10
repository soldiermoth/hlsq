package hlsqlib

import (
	"fmt"
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
		`#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=1966314,BANDWIDTH=2164328,CODECS="hvc1.2.4.L123.B0,mp4a.40.2",RESOLUTION=960x540,FRAME-RATE=60.000000,CLOSED-CAPTIONS="cc",AUDIO="a1",SUBTITLES="sub1"`: Tag{
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
	for in, expected := range testcases {
		found, err := ParseLine(in)
		if err != nil {
			t.Fatalf("could not parse line=%q err=%q", in, err)
		}
		if !reflect.DeepEqual(found, expected) {
			var extraContext string
			fT, fIsT := found.(Tag)
			eT, eIsT := expected.(Tag)
			if fIsT && eIsT && (len(eT.Attrs) > 0 || len(fT.Attrs) > 0) {
				maxLen := len(fT.Attrs)
				if len(eT.Attrs) > maxLen {
					maxLen = len(eT.Attrs)
				}
				for i := 0; i < maxLen; i++ {
					var eA, fA Attr
					if i < len(eT.Attrs) {
						eA = eT.Attrs[i]
					}
					if i < len(fT.Attrs) {
						fA = fT.Attrs[i]
					}
					if !reflect.DeepEqual(eA, fA) {
						extraContext += fmt.Sprintf("\nAttr=#%d expectec=%#v got=%#v", i, eA, fA)
					}
				}
			}
			t.Fatalf("did not equal expected in=%q\nEXPECTED:\n%#v\nFOUND:\n%#v%s", in, expected, found, extraContext)
		}
	}
}
