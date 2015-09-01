package main

import (
	"image"

	"golang.org/x/mobile/gl"

	"github.com/shazow/linerage3d/collision"
)

func NewArenaNode(bounds image.Rectangle, shader Shader) *arena {
	shape := NewStaticShape()
	shape.vertices = []float32{
		float32(bounds.Min.X), 0, float32(bounds.Min.Y),
		float32(bounds.Max.X), 0, float32(bounds.Min.Y),
		float32(bounds.Max.X), 0, float32(bounds.Max.Y),
		float32(bounds.Max.X), 0, float32(bounds.Max.Y),
		float32(bounds.Min.X), 0, float32(bounds.Min.Y),
		float32(bounds.Min.X), 0, float32(bounds.Max.Y),
	}
	shape.Buffer()

	return newArena(bounds, &Node{
		Shape:  shape,
		shader: shader,
	})
}

func newArena(bounds image.Rectangle, node *Node) *arena {
	arena := &arena{
		Node:     node,
		Collider: collision.LinearCollider(bounds),
	}
	arena.Reset()
	return arena
}

type arena struct {
	*Node
	collision.Collider
}

func (shape *arena) Draw(camera Camera) {
	shader := shape.shader
	gl.Uniform3fv(shader.Uniform("material.ambient"), []float32{0.05, 0.0, 0.02})
	gl.Uniform3fv(shader.Uniform("lights[0].color"), []float32{0.2, 0.1, 0.1})
	gl.Uniform1f(shader.Uniform("lights[0].intensity"), 0.3)

	// TODO: Move this into a texture
	gl.Uniform1f(shader.Uniform("lights[1].intensity"), 25.0)
	gl.Uniform3fv(shader.Uniform("lights[1].position"), []float32{0, 20, 0})
	gl.Uniform3fv(shader.Uniform("lights[1].color"), []float32{0.05, 0.0, 0.1})

	shape.Node.Draw(camera)
}
