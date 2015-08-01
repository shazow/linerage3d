package main

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

const turnSpeed = 0.05

type linerageWorld struct {
	scene    Scene
	bindings *Bindings

	line    *Line
	emitter Emitter
}

func LinerageWorld(scene Scene, bindings *Bindings, shaders Shaders) (World, error) {
	// Load shaders
	err := shaders.Load("line", "particle", "skybox")
	if err != nil {
		return nil, err
	}

	// Load textures
	skyboxTex, err := LoadTextureCube("square.png")
	if err != nil {
		return nil, err
	}

	// Add shader material
	shader := shaders.Get("line")
	shader.Use()
	gl.Uniform3fv(shader.Uniform("material.ambient"), []float32{0.1, 0.15, 0.4})
	gl.Uniform3fv(shader.Uniform("material.diffuse"), []float32{0.8, 0.6, 0.6})
	gl.Uniform3fv(shader.Uniform("material.specular"), []float32{1.0, 1.0, 1.0})
	//gl.Uniform1f(shader.Uniform("material.shininess"), 16.0)
	//gl.Uniform1f(shader.Uniform("material.refraction"), 1.0/1.52)

	gl.Uniform3fv(shader.Uniform("lights[0].color"), []float32{0.4, 0.2, 0.1})
	gl.Uniform1f(shader.Uniform("lights[0].intensity"), 0.0)

	//gl.Uniform3fv(shader.Uniform("lights[1].color"), []float32{0.4, 0.2, 0.1})
	//gl.Uniform1f(shader.Uniform("lights[1].intensity"), 0.0)

	// Make skybox
	// TODO: Add closer, or use a texture loader
	scene.Add(NewSkybox(shaders.Get("skybox"), skyboxTex))

	// Make line
	line := NewLine(shaders.Get("line"), 6*4*10000)
	line.Add(0)
	line.shape.Buffer(0)
	scene.Add(line)

	// Make particle emitter
	emitter := ParticleEmitter(mgl.Vec3{0, 1, 1}, 20, 1)
	scene.Add(&Node{Shape: emitter, shader: shaders.Get("particle")})

	/*
		// Cube for funsies:
		cube := NewStaticShape()
		cube.vertices = skyboxVertices
		cube.normals = skyboxNormals
		cube.indices = skyboxIndices
		cube.Buffer()
		scene.nodes = append(scene.nodes, Node{Shape: cube, shader: lineShader})
	*/

	/*
		// Reflective floor
		scene.nodes = append(scene.nodes, Node{
			Shape: NewFloor(Node{Shape: shape: lineShader}),
		})
	*/

	return &linerageWorld{
		scene:    scene,
		bindings: bindings,

		line:    line,
		emitter: emitter,
	}, err
}

func (world *linerageWorld) Focus() mgl.Vec3 {
	return world.line.position
}

func (world *linerageWorld) Tick(interval time.Duration) {
	var lineRotate float64

	if world.bindings.Pressed(KeyLineLeft) {
		lineRotate -= turnSpeed
	}
	if world.bindings.Pressed(KeyLineRight) {
		lineRotate += turnSpeed
	}

	// Spinny!
	//rotation := mgl.HomogRotate3D(float32(since.Seconds()), AxisFront)
	//world.scene.transform = &rotation

	world.line.Tick(interval, lineRotate)
	world.emitter.MoveTo(world.line.position)
	world.emitter.Tick(interval)
}
