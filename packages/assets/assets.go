package assets

import "encoding/xml"

func LoadFont(pngPath, xmlPath string) {

}

// private ========================================================

type font struct {
	XMLName xml.Name `xml:"font"`
	Common  struct {
		ScaleW int `xml:"scaleW,attr"`
		ScaleH int `xml:"scaleH,attr"`
	} `xml:"common"`
	DistanceField struct {
		DistanceRange int `xml:"distanceRange,attr"`
	} `xml:"distanceField"`
	Chars struct {
		Count int       `xml:"count,attr"`
		Chars [256]char `xml:"char"`
	} `xml:"chars"`
}

type char struct {
	ID       int    `xml:"id,attr"`
	Char     string `xml:"char,attr"`
	X        int    `xml:"x,attr"`
	Y        int    `xml:"y,attr"`
	Width    int    `xml:"width,attr"`
	Height   int    `xml:"height,attr"`
	XOffset  int    `xml:"xoffset,attr"`
	YOffset  int    `xml:"yoffset,attr"`
	XAdvance int    `xml:"xadvance,attr"`
}
