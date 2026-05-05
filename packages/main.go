package main

import (
	"big-black-box/packages/assets"
	"big-black-box/packages/engine"
	"big-black-box/packages/geometry"
	"big-black-box/packages/graphics"
	"big-black-box/packages/input/keyboard"
	"big-black-box/packages/input/keyboard/key"
	"big-black-box/packages/utility/angle"
	"big-black-box/packages/utility/number"
	"big-black-box/packages/utility/time"
)

func main() {
	// Main view: camera follows player, MaskArea limits visible world to a 600x450 region.
	var view = graphics.NewView()

	// Minimap: overview of the whole playfield in the top-left corner.
	var minimap = graphics.View{
		X: 0, Y: 0,
		Zoom:       0.15,
		WindowArea: graphics.Area{X: 10, Y: 10, Width: 320 * 2, Height: 200 * 2},
	}

	var obstacle = geometry.NewRoundedRectangle(400*2, 400*2, 450*2, 250*2, 0, 0.5)
	var player = geometry.NewCircle(0, 0, 50*2)
	var staticShapes = []geometry.Shape{
		geometry.NewCapsule(150*2, 600*2, 350*2, 750*2, 30*2),
		geometry.NewRoundedRectangle(-100*2, 400*2, 60*2, 400*2, 0, 0.2),
		geometry.NewRectangle(750*2, 250*2, 150*2, 150*2, 45),
	}
	var rotSpeeds = []float32{0, 0.3, -0.2}

	var font = assets.LoadFont("font.png", "font.xml")

	var a float32
	engine.Run(60, func() {
		const speed = 10.0
		if keyboard.IsPressed(key.ArrowLeft) || keyboard.IsPressed(key.A) {
			player.X -= speed
		}
		if keyboard.IsPressed(key.ArrowRight) || keyboard.IsPressed(key.D) {
			player.X += speed
		}
		if keyboard.IsPressed(key.ArrowUp) || keyboard.IsPressed(key.W) {
			player.Y -= speed
		}
		if keyboard.IsPressed(key.ArrowDown) || keyboard.IsPressed(key.S) {
			player.Y += speed
		}

		if keyboard.IsPressed(key.ShiftLeft) {
			player.Width += 1
		}
		if keyboard.IsPressed(key.ControlLeft) {
			player.Width -= 1
		}

		obstacle.Angle += 1
		for i := range staticShapes {
			staticShapes[i].Angle += rotSpeeds[i]
		}
		a += 0.3

		obstacle.Roundness = (1 + number.Sine(time.Running())) / 2

		player = obstacle.Collide(player)
		for _, s := range staticShapes {
			player = s.Collide(player)
		}

		var ang = angle.BetweenPoints(player.X, player.Y, obstacle.X, obstacle.Y)
		var hitX, hitY = obstacle.Raycast(player.X, player.Y, ang, 1000)

		// Camera: player stays at view center.
		view.X = player.X
		view.Y = player.Y

		//=================================================================

		view.DrawShape(obstacle)
		for _, s := range staticShapes {
			view.DrawShape(s)
		}
		view.DrawShape(player)
		if !number.IsNaN(hitX) {
			view.DrawShape(geometry.NewCircle(hitX, hitY, 6))
		}
		view.DrawText(geometry.NewRoundedRectangle(0, 0, 1000, 1000, a, 0), font)
		view.DrawFPS()

		// Minimap: same shapes, clipped to the WindowArea.
		minimap.DrawShape(obstacle)
		for _, s := range staticShapes {
			minimap.DrawShape(s)
		}
		minimap.DrawShape(player)
	})
}
