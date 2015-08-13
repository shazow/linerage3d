package main

import (
	"image"

	"golang.org/x/mobile/gl"

	"github.com/shazow/linerage3d/collision"
)

func NewArenaNode(bounds image.Rectangle, shader Shader) *arena {
	shape := NewStaticShape()
	shape.vertices = floorVertices
	shape.Buffer()

	return newArena(bounds, &Node{
		Shape:  shape,
		shader: shader,
	})
}

func newArena(bounds image.Rectangle, node *Node) *arena {
	arena := &arena{
		Node:     node,
		Collider: collision.GridCollider(bounds),
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
	gl.Uniform3fv(shader.Uniform("material.ambient"), []float32{0.2, 0.2, 0.2})
	gl.Uniform3fv(shader.Uniform("lights[0].color"), []float32{0.2, 0.1, 0.1})
	gl.Uniform1f(shader.Uniform("lights[0].intensity"), 0.3)

	shape.Node.Draw(camera)
}

var arenaVertices = []float32{
	// Floor
	-100, 0, -100,
	100, 0, -100,
	100, 0, 100,
	100, 0, 100,
	-100, 0, -100,
	-100, 0, 100,
}
