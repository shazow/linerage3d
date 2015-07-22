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
	shape  *StaticShape
	shader *Shader
}

func (skybox *Skybox) Draw(camera Camera) {
	shader := skybox.shader
	shader.Use()

	gl.DepthMask(false)

	projection, view := camera.Projection(), camera.View().Mat3().Mat4()
	gl.UniformMatrix4fv(shader.projection, projection[:])
	gl.UniformMatrix4fv(shader.view, view[:])

	node := skybox.shape
	//gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, node.glTex)
	node.EnableAttrib(shader.vertCoord, shader.vertNormal, shader.vertTexCoord)
	gl.DrawArrays(gl.TRIANGLES, 0, node.Len())
	node.DisableAttrib(shader.vertCoord, shader.vertNormal, shader.vertTexCoord)

	gl.DepthMask(true)
}

type Scene struct {
	// TODO: Add a shader registry instead
	shader *Shader
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
	gl.UniformMatrix4fv(shader.projection, projection[:])
	gl.UniformMatrix4fv(shader.view, view[:])

	// Light
	gl.Uniform3fv(shader.lightIntensities, []float32{1, 1, 1})
	gl.Uniform3fv(shader.lightPosition, []float32{1, 1, 1})

	for _, node := range scene.nodes {
		model := transformModel(node.transform, scene.transform)
		gl.UniformMatrix4fv(shader.model, model[:])

		normal := model.Mul4(view).Inv().Transpose()
		gl.UniformMatrix4fv(shader.normalMatrix, normal[:])

		node.EnableAttrib(shader.vertCoord, shader.vertNormal, shader.vertTexCoord)
		gl.DrawArrays(gl.TRIANGLES, 0, node.Len())
		node.DisableAttrib(shader.vertCoord, shader.vertNormal, shader.vertTexCoord)
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
