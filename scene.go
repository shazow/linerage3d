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

type Scene struct {
	Camera
	*Shader

	ambientColor mgl.Vec3
	lights       []Light
	nodes        []Node
	transform    *mgl.Mat4
}

func (scene *Scene) String() string {
	return fmt.Sprintf("%d nodes, %d lights, ambient %+v, camera:\n%s", len(scene.nodes), len(scene.lights), scene.ambientColor, scene.Camera)
}

func (scene *Scene) Draw() {
	shader := scene.Shader

	// Setup MVP
	projection, view := scene.Projection(), scene.View()
	gl.UniformMatrix4fv(shader.projection, projection[:])
	gl.UniformMatrix4fv(shader.view, view[:])

	// Light
	gl.Uniform3fv(shader.lightIntensities, []float32{1, 1, 1})
	gl.Uniform3fv(shader.lightPosition, []float32{1, 1, 1})

	for _, node := range scene.nodes {
		// TODO: Get model position from Shape, then translate
		//log.Printf("Drawing: %s", node.String())

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
