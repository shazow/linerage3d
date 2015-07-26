package main

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

// TODO: Load this from an .obj file in the asset repository?

var skyboxVertices = []float32{
	-1, 1, -1,
	-1, -1, -1,
	1, -1, -1,
	1, 1, -1,
	-1, -1, 1,
	-1, 1, 1,
	1, -1, 1,
	1, 1, 1,
}

var skyboxNormals = []float32{
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, -1.0,
}

var skyboxIndices = []uint8{
	0, 1, 2, 2, 3, 0,
	4, 1, 0, 0, 5, 4,
	2, 6, 7, 7, 3, 2,
	4, 5, 7, 7, 6, 4,
	0, 3, 7, 7, 5, 0,
	1, 4, 2, 2, 4, 6,
}

func NewSkybox(shader Shader, texture gl.Texture) Shape {
	skyboxShape := NewStaticShape()
	skyboxShape.vertices = skyboxVertices
	skyboxShape.indices = skyboxIndices
	skyboxShape.Buffer()
	skyboxShape.Texture = texture

	skybox := &Skybox{
		StaticShape: skyboxShape,
		shader:      shader,
	}

	return skybox
}

type Skybox struct {
	*StaticShape
	shader Shader
}

func (shape *Skybox) Draw(_ Shader, camera Camera) {
	shader := shape.shader
	shader.Use()

	gl.DepthFunc(gl.LEQUAL)
	gl.DepthMask(false)

	projection, view := camera.Projection(), camera.View().Mat3().Mat4()
	gl.UniformMatrix4fv(shader.Uniform("projection"), projection[:])
	gl.UniformMatrix4fv(shader.Uniform("view"), view[:])

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, shape.Texture)

	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)
	gl.EnableVertexAttribArray(shader.Attrib("vertCoord"))
	gl.VertexAttribPointer(shader.Attrib("vertCoord"), vertexDim, gl.FLOAT, false, shape.Stride(), 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, shape.IBO)
	gl.DrawElements(gl.TRIANGLES, len(shape.indices), gl.UNSIGNED_BYTE, 0)
	gl.DisableVertexAttribArray(shader.Attrib("vertCoord"))

	gl.DepthMask(true)
	gl.DepthFunc(gl.LESS)
}

var floorVertices = []float32{
	-100, 0, -100,
	100, 0, -100,
	100, 0, 100,
	100, 0, 100,
	-100, 0, -100,
	-100, 0, 100,
}

var floorNormals = []float32{
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
}

type Floor struct {
	Shape
	reflected []Node
	transform *mgl.Mat4
}

func (scene *Floor) Draw(shader Shader, camera Camera) {
	gl.Enable(gl.STENCIL_TEST)
	gl.StencilFunc(gl.ALWAYS, 1, 0xFF)
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)
	gl.StencilMask(0xFF)
	gl.DepthMask(false)
	gl.Clear(gl.STENCIL_BUFFER_BIT)

	// Draw floor
	gl.Uniform4fv(shader.Uniform("surfaceColor"), []float32{0, 0, 0, 0.3})
	scene.Shape.Draw(shader, camera)

	// Draw reflections
	gl.StencilFunc(gl.EQUAL, 1, 0xFF)
	gl.StencilMask(0x00)
	gl.DepthMask(true)

	view := camera.View()
	gl.Uniform4fv(shader.Uniform("surfaceColor"), []float32{0.2, 0.2, 0.2, 1})
	for _, node := range scene.reflected {
		model := transformModel(node.transform, scene.transform)
		gl.UniformMatrix4fv(shader.Uniform("model"), model[:])

		normal := model.Mul4(view).Inv().Transpose()
		gl.UniformMatrix4fv(shader.Uniform("normalMatrix"), normal[:])

		node.Draw(shader, camera)
	}

	gl.Disable(gl.STENCIL_TEST)
	gl.Uniform4fv(shader.Uniform("surfaceColor"), []float32{0, 0, 0, 0})
}

func NewFloor(reflected ...Node) Shape {
	floor := NewStaticShape()
	floor.vertices = floorVertices
	floor.normals = floorNormals
	floor.Buffer()
	flipped := mgl.Scale3D(1, -1, 1)
	return &Floor{
		Shape: floor, reflected: reflected,
		transform: &flipped,
	}
}
