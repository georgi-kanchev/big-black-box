package assets

import (
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/file"
	"big-black-box/packages/utility/storage"
	"bytes"
	"encoding/xml"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Font internal.Font

func LoadFont(pngPath, xmlPath string) Font {
	var fontData = font{}
	var xmlContent = file.LoadText(xmlPath)
	if xmlContent == "" {
		return 0
	}
	storage.FromXML(xmlContent, &fontData)
	if fontData.Common.ScaleW == 0 || fontData.Common.ScaleH == 0 {
		return 0
	}

	var img, _, err = image.Decode(bytes.NewReader(file.LoadBytes(pngPath)))
	if err != nil {
		log.Println(err)
		return 0
	}

	var charsInXML = make(map[string]char)
	var runes = []rune(all)
	var asset = ebiten.NewImageFromImage(img)
	for _, c := range fontData.Chars.Chars {
		charsInXML[c.Char] = c
	}

	for i, r := range runes {
		if i >= 256 {
			break
		}

		var charStr = string(r)
		var glyph, _ = charsInXML[charStr]
		var packed uint32

		packed |= uint32(glyph.X) & 0x1FF             // 9 bits
		packed |= (uint32(glyph.Y) & 0x1FF) << 9      // 9 bits
		packed |= (uint32(glyph.Width) & 0x7F) << 18  // 7 bits
		packed |= (uint32(glyph.Height) & 0x7F) << 25 // 7 bits

		var c1 = color.RGBA{
			R: byte(packed),
			G: byte(packed >> 8),
			B: byte(packed >> 16),
			A: byte(packed >> 24),
		}
		var c2 = color.RGBA{
			R: byte(glyph.XOffset + 128),
			G: byte(glyph.YOffset + 128),
			B: byte(glyph.XAdvance),
			A: 255, // unused
		}

		asset.Set(i*2, 0, c1)
		asset.Set(i*2+1, 0, c2)
	}

	internal.Fonts = append(internal.Fonts, asset)
	return Font(len(internal.Fonts))
}

// private ========================================================

const core = " .,;:!?¡¿\"'()[]{}<>-/\\@#$%^&*_+=|~`" + "0123456789" + "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const latin = "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝßŒŠŽŁŃŚŹŻĆČĐŐŰàáâãäåæçèéêëìíîïñòóôõöøùúûüýÿœšžłńśźżćčđőűẞ"
const cyrillic = "АБВГДЕЁЖЗИЙКЛМНОПРСТУΦΧЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюяҐЄІЇґєії"
const all = core + latin + cyrillic

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
		Count int    `xml:"count,attr"`
		Chars []char `xml:"char"`
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
