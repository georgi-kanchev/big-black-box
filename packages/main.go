package main

import (
	"big-black-box/packages/engine"
	"big-black-box/packages/geometry"
	"big-black-box/packages/graphics"
	"big-black-box/packages/input/keyboard"
	"big-black-box/packages/input/keyboard/key"
	"big-black-box/packages/utility/angle"
	"big-black-box/packages/utility/number"
	"big-black-box/packages/utility/time"
)

const core = " .,;:!?¡¿\"'()[]{}<>-/\\@#$%^&*_+=|~`" + "0123456789" + "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const latin = "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝßŒŠŽŁŃŚŹŻĆČĐŐŰàáâãäåæçèéêëìíîïñòóôõöøùúûüýÿœšžłńśźżćčđőűẞ"
const cyrillic = "АБВГДЕЁЖЗИЙКЛМНОПРСТУΦΧЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюяҐЄІЇґєії"
const all = core + latin + cyrillic

func main() {
	const offset = 250.0
	var cam = graphics.NewCamera()
	var obstacle = geometry.NewRoundedRectangle(400+offset, 400, 450, 250, 0, 0.5)
	var player = geometry.NewCircle(100+offset, 100, 40)
	var staticShapes = []geometry.Shape{
		geometry.NewCapsule(150+offset, 600, 350+offset, 750, 30),
		geometry.NewRoundedRectangle(offset-100, 400, 60, 400, 0, 0.2),
		geometry.NewRectangle(750+offset, 250, 150, 150, 45),
	}
	var rotSpeeds = []float32{0, 0.3, -0.2}

	engine.Run(func() {
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

		obstacle.Roundness = number.Sine(time.Running())

		player = obstacle.Collide(player)

		for _, s := range staticShapes {
			player = s.Collide(player)
		}

		var ang = angle.BetweenPoints(player.X, player.Y, obstacle.X, obstacle.Y)
		var hitX, hitY = obstacle.Raycast(player.X, player.Y, ang, 1000)

		//=================================================================

		cam.DrawShape(obstacle)
		for _, s := range staticShapes {
			cam.DrawShape(s)
		}
		cam.DrawShape(player)

		if !number.IsNaN(hitX) {
			cam.DrawShape(geometry.NewCircle(hitX, hitY, 6))
		}
		cam.DrawDebugInfo()
	})
}
