package main

import (
	"image"
	"log"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// TODO: Load into here
type lineShader struct {
	Projection   [16]float32
	View         [16]float32
	Model        [16]float32
	NormalMatrix [16]float32

	Material struct {
		Ambient    [3]float32
		Diffuse    [3]float32
		Specular   [3]float32
		Shininess  float32
		Refraction float32
	}
	Lights []struct {
		Color     [3]float32
		Intensity float32
	}
}

const turnSpeed = 0.1

type linerageWorld struct {
	scene    Scene
	bindings *Bindings

	collisionToken *cellSegment
	arena          *arena
	line           *Line
	emitter        Emitter
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
	line := NewLine(shaders.Get("line"), 2*4*100000)
	line.Buffer(0)
	scene.Add(line)

	/*
		shader := shaders.Get("line")
		shader.Use()
		gl.Uniform3fv(shader.Uniform("lights[1].color"), []float32{1.0, 0.9, 0.9})
		gl.Uniform3fv(shader.Uniform("lights[1].position"), []float32{0.0, 20.0, 0.0})
		gl.Uniform1f(shader.Uniform("lights[1].intensity"), 1.0)
	*/

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
		scene.Add(NewFloor(shaders.Get("line"), line))
	*/

	arena := NewArenaNode(image.Rect(-10, -10, 10, 10), shaders.Get("line"))
	scene.Add(arena)

	bindings.On(KeyReload, func(_ KeyBinding) {
		log.Println("Reloading shaders.")
		err := shaders.Reload()
		if err != nil {
			log.Println("Shader reload error:", err)
		}
	})

	bindings.On(KeyDebug, func(_ KeyBinding) {
		log.Println("Segment: ", line.segments)
		log.Println(dumpArena(arena))
	})

	return &linerageWorld{
		scene:    scene,
		bindings: bindings,

		arena:   arena,
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

	var err error
	world.collisionToken, err = world.arena.Add(world.collisionToken, world.line.segments)
	if err != nil {
		n := len(world.line.segments) - 4
		if n < 0 {
			n = 0
		}
		log.Printf("Collision with %s\n\tLast segments: %v", err, world.line.segments[n:])
	}

}
