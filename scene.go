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

func (light *Light) MoveTo(position mgl.Vec3) {
	light.position = position
}

// TODO: node tree with transforms
type Node struct {
	Shape
	transform *mgl.Mat4
	shader    Shader
}

func (node *Node) String() string {
	return fmt.Sprintf("<Shape of %d vertices; transform: %v>", node.Len(), node.transform)
}

type Scene struct {
	// TODO: Add a shader registry instead
	shader Shader
	skybox Shape

	lights    []*Light
	nodes     []Node
	transform *mgl.Mat4
}

func (scene *Scene) String() string {
	return fmt.Sprintf("%d nodes, %d lights, ambient %+v", len(scene.nodes), len(scene.lights))
}

func (scene *Scene) Draw(camera Camera) {
	if scene.skybox != nil {
		scene.skybox.Draw(scene.shader, camera)
	}

	shader := scene.shader

	// Setup MVP
	projection, view, position := camera.Projection(), camera.View(), camera.Position()

	for _, light := range scene.lights {
		// Light
		gl.Uniform3fv(shader.Uniform("lightIntensities"), light.color[:])
		gl.Uniform3fv(shader.Uniform("lightPosition"), light.position[:])
		fmt.Println(light.position)
		break
	}

	for _, node := range scene.nodes {
		shader = scene.shader
		if node.shader != nil {
			shader = node.shader
		}
		shader.Use()

		// TODO: Move these into node.Draw?
		model := transformModel(node.transform, scene.transform)
		normal := model.Mul4(view).Inv().Transpose()

		// Camera space
		gl.UniformMatrix4fv(shader.Uniform("model"), model[:])
		gl.UniformMatrix4fv(shader.Uniform("view"), view[:])
		gl.UniformMatrix4fv(shader.Uniform("projection"), projection[:])
		gl.UniformMatrix4fv(shader.Uniform("normalMatrix"), normal[:])
		gl.UniformMatrix4fv(shader.Uniform("cameraPos"), position[:])

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
