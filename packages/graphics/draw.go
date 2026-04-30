package graphics

import (
	"big-black-box/packages/assets"
	"big-black-box/packages/debug"
	"big-black-box/packages/geometry"
	"big-black-box/packages/internal"
	"big-black-box/packages/utility/color/palette"
	"big-black-box/packages/utility/text"
	"big-black-box/packages/utility/time"
	"big-black-box/packages/utility/time/unit"
)

func (v View) DrawShape(shape geometry.Shape) {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		Kind:         internal.KindShape,
		Shape:        internal.Shape(shape),
		Color:        palette.White,
		OutlineSize:  1,
		OutlineColor: palette.Green,
	})
}

func (v View) DrawFPS() {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		Kind:  internal.KindText,
		Shape: internal.Shape{X: 5, Y: 5},
		Text:  text.Start().String("FPS: ").Int(int(internal.FPS)).End(),
		Color: palette.White,
	})
}
func (v View) DrawDebugInfo() {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		Kind:  internal.KindText,
		Shape: internal.Shape{X: 5, Y: 5},
		Text: text.Start().String("FPS: ").Int(int(internal.FPS)).String(" TPS: ").Int(int(internal.TPS)).String("\n\n").
			String("Runtime: ").
			String(time.AsClock12(internal.Runtime, ":", unit.Hour|unit.Minute|unit.Second|unit.Millisecond, false)).
			String("\n\n").
			String(debug.MemoryUsage()).
			End(),
		Color: palette.White,
	})
}

func (v View) DrawImage(shape geometry.Shape, image assets.Image) {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		Kind:         internal.KindImage,
		Shape:        internal.Shape(shape),
		Color:        palette.White,
		Image:        internal.Image(image),
		OutlineSize:  1,
		OutlineColor: palette.Red,
	})
}

func (v View) DrawText(shape geometry.Shape, font assets.Font) {
	internal.Queue(internal.LayerDefault, internal.DrawItem{
		Kind:  internal.KindText,
		Shape: internal.Shape(shape),
		Color: palette.White,
		Font:  internal.Font(font),
	})
}
