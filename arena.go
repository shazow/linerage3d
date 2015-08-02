package main

import "golang.org/x/mobile/gl"

func NewArena(shader Shader) Drawable {
	shape := NewStaticShape()
	shape.vertices = floorVertices
	shape.Buffer()

	return &arena{
		Node: Node{
			Shape:  shape,
			shader: shader,
		},
	}
}

type arena struct {
	Node
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
