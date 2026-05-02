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
	const offset = 250.0

	// Main view: camera follows player, MaskArea limits visible world to a 600x450 region.
	var view = graphics.NewView()
	view.MaskArea = graphics.Area{X: -300, Y: -225, Width: 600, Height: 450}

	// Minimap: overview of the whole playfield in the top-left corner.
	var minimap = graphics.View{
		X: 500, Y: 350,
		Zoom:       0.15,
		WindowArea: graphics.Area{X: 10, Y: 10, Width: 320, Height: 200},
	}

	var obstacle = geometry.NewRoundedRectangle(400+offset, 400, 450, 250, 0, 0.5)
	var player = geometry.NewCircle(100+offset, 100, 40)
	var staticShapes = []geometry.Shape{
		geometry.NewCapsule(150+offset, 600, 350+offset, 750, 30),
		geometry.NewRoundedRectangle(offset-100, 400, 60, 400, 0, 0.2),
		geometry.NewRectangle(750+offset, 250, 150, 150, 45),
	}
	var rotSpeeds = []float32{0, 0.3, -0.2}

	var img = assets.LoadImage("font.png")

	var a float32
	engine.Run(60, func() {
		const speed = 5.0
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
		a++

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
		view.DrawImage(geometry.NewRoundedRectangle(300, 300, 200, 200, a, 0), img)
		view.DrawFPS()

		// Minimap: same shapes, clipped to the WindowArea.
		minimap.DrawShape(obstacle)
		for _, s := range staticShapes {
			minimap.DrawShape(s)
		}
		minimap.DrawShape(player)
	})
}
