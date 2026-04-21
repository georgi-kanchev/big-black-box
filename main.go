package main

import (
	"big-black-box/engine"
	"big-black-box/geometry"
	"big-black-box/graphics"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	var cam = graphics.NewCamera()
	var shapeA = geometry.Shape{X: 400, Y: 400, Width: 450, Height: 250, Angle: 20, Roundness: 0.5}
	var shapeB = geometry.Shape{Width: 120, Height: 120, Roundness: 1}

	engine.Run(func() {
		const speed = 5.0
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			shapeB.X -= speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			shapeB.X += speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			shapeB.Y -= speed
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			shapeB.Y += speed
		}

		currentAngle := shapeA.Angle
		currentRoundness := shapeA.Roundness
		newAngle := currentAngle + 1.0
		shapeA.Angle = newAngle
		shapeA.Roundness = currentRoundness

		shapeB = shapeA.Collide(shapeB)

		cam.RenderShape(shapeA)
		cam.RenderShape(shapeB)
	})
}
