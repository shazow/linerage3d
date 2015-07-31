package main

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type linerageWorld struct {
	scene    Scene
	bindings *Bindings

	light   *Light
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

	// Make skybox
	// TODO: Add closer, or use a texture loader
	scene.Add(NewSkybox(shaders.Get("skybox"), skyboxTex))

	// Make line
	shape := NewDynamicShape(6 * 4 * 10000)
	line := NewLine(shape)
	line.Add(0)
	line.Buffer(0)
	scene.Add(&Node{Shape: shape, shader: shaders.Get("line")})

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

	// Light
	light := &Light{color: mgl.Vec3{0.7, 0.5, 0.5}, position: line.position}
	scene.Add(light)

	return &linerageWorld{
		scene:    scene,
		bindings: bindings,

		light:   light,
		line:    line,
		emitter: emitter,
	}, err
}

func (world *linerageWorld) Focus() mgl.Vec3 {
	return world.line.position
}

func (world *linerageWorld) Tick(interval time.Duration) {
	var lineRotate float32

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
	world.light.MoveTo(world.line.position)
	world.emitter.MoveTo(world.line.position)
	world.emitter.Tick(interval)
}
