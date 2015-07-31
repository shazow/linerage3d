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

type Drawable interface {
	Draw(Camera)
	Transform(*mgl.Mat4) mgl.Mat4
	UseShader(Shader) (Shader, bool)
}

// TODO: node tree with transforms
type Node struct {
	Shape
	transform *mgl.Mat4
	shader    Shader
}

func (node *Node) Draw(camera Camera) {
	node.Shape.Draw(node.shader, camera)
}

func (node *Node) UseShader(parent Shader) (Shader, bool) {
	if node.shader == nil || node.shader == parent {
		node.shader = parent
		return parent, false
	}
	node.shader.Use()
	return node.shader, true
}

func (node *Node) Transform(parent *mgl.Mat4) mgl.Mat4 {
	return MultiMul(node.transform, parent)
}

func (node *Node) String() string {
	return fmt.Sprintf("<Shape of %d vertices; transform: %v>", node.Len(), node.transform)
}

type Scene interface {
	Add(interface{})
	Draw(Camera)
	String() string
}

func NewScene() Scene {
	return &sliceScene{
		lights: []*Light{},
		nodes:  []Drawable{},
	}
}

type sliceScene struct {
	lights    []*Light
	nodes     []Drawable
	transform *mgl.Mat4
}

func (scene *sliceScene) String() string {
	return fmt.Sprintf("%d nodes, %d lights", len(scene.nodes), len(scene.lights))
}

func (scene *sliceScene) Add(item interface{}) {
	switch item := item.(type) {
	case *Light:
		scene.lights = append(scene.lights, item)
	default:
		scene.nodes = append(scene.nodes, item.(Drawable))
	}
}

func (scene *sliceScene) Draw(camera Camera) {
	// Setup MVP
	projection, view, position := camera.Projection(), camera.View(), camera.Position()

	light := scene.lights[0]
	var parentShader Shader
	for _, node := range scene.nodes {
		shader, changed := node.UseShader(parentShader)

		if changed {
			// TODO: Pre-load these into relevant shaders?
			gl.Uniform3fv(shader.Uniform("lightIntensities"), light.color[:])
			gl.Uniform3fv(shader.Uniform("lightPosition"), light.position[:])
			gl.UniformMatrix4fv(shader.Uniform("cameraPos"), position[:])
			gl.UniformMatrix4fv(shader.Uniform("view"), view[:])
			gl.UniformMatrix4fv(shader.Uniform("projection"), projection[:])
		}

		// TODO: Move these into node.Draw?
		model := node.Transform(scene.transform)
		normal := model.Mul4(view).Inv().Transpose()

		// Camera space
		gl.UniformMatrix4fv(shader.Uniform("model"), model[:])
		gl.UniformMatrix4fv(shader.Uniform("normalMatrix"), normal[:])

		node.Draw(camera)
	}
}
