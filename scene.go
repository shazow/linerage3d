package main

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

type Light struct {
	color    mgl.Vec3
	position mgl.Vec3
}

// TODO: node tree with transforms
type Node struct {
	Shape
	transform *mgl.Mat4
}

func (node *Node) String() string {
	return fmt.Sprintf("<Shape of %d vertices; transform: %v>", node.Len(), node.transform)
}

type Skybox struct {
	*StaticShape
	shader Shader
}

func (shape *Skybox) Draw(camera Camera) {
	shader := shape.shader
	shader.Use()

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
}

type Scene struct {
	// TODO: Add a shader registry instead
	shader Shader
	skybox *Skybox

	ambientColor mgl.Vec3
	lights       []Light
	nodes        []Node
	transform    *mgl.Mat4
}

func (scene *Scene) String() string {
	return fmt.Sprintf("%d nodes, %d lights, ambient %+v", len(scene.nodes), len(scene.lights), scene.ambientColor)
}

func (scene *Scene) Draw(camera Camera) {
	if scene.skybox != nil {
		scene.skybox.Draw(camera)
	}

	shader := scene.shader
	shader.Use()

	// Setup MVP
	projection, view := camera.Projection(), camera.View()
	gl.UniformMatrix4fv(shader.Uniform("projection"), projection[:])
	gl.UniformMatrix4fv(shader.Uniform("view"), view[:])

	// Light
	gl.Uniform3fv(shader.Uniform("lightIntensities"), []float32{1, 1, 1})
	gl.Uniform3fv(shader.Uniform("lightPosition"), []float32{1, 1, 1})

	for _, node := range scene.nodes {
		model := transformModel(node.transform, scene.transform)
		gl.UniformMatrix4fv(shader.Uniform("model"), model[:])

		normal := model.Mul4(view).Inv().Transpose()
		gl.UniformMatrix4fv(shader.Uniform("normalMatrix"), normal[:])

		node.Draw(shader, camera)
	}
}

func transformModel(model *mgl.Mat4, scene *mgl.Mat4) mgl.Mat4 {
	if model == nil && scene == nil {
		return mgl.Ident4()
	}
	if model != nil && scene != nil {
		return model.Mul4(*scene)
	}
	if model != nil {
		return *model
	} else {
		return *scene
	}
}
