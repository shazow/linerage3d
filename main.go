package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type Engine struct {
	camera *EulerCamera
	scene  *Scene
	shader *Shader

	touchLoc geom.Point
	started  time.Time
	offset   int
	pos      [5]float32
	frames   int

	line *DynamicShape
}

func (e *Engine) Start() {
	var err error

	log.Println("Loading shaders...")

	e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
	if err != nil {
		panic(fmt.Sprintln("LoadProgram failed:", err))
	}

	e.camera.Move(mgl.Vec3{3, 3, 3})
	e.camera.Pan(mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})

	shape := DynamicShape{}
	e.line = &shape
	e.scene.nodes = append(e.scene.nodes, Node{Shape: &shape})

	e.pos = [5]float32{0, 0, 0, 1, 1}
	shape.vertices = append(shape.vertices, Quad(e.pos[0], e.pos[1], e.pos[2], e.pos[3], e.pos[4])...)
	//shape.normals = append(shape.normals, []float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}...)
	//shape.vertices = cubeData

	shape.Init(6 * 4 * 1000)
	shape.Buffer(0)

	gl.UseProgram(e.shader.program)
	e.scene.Bind()

	e.started = time.Now()

	log.Println("Starting: ", e.scene.String())
}

func (e *Engine) Stop() {
	gl.DeleteProgram(e.shader.program)
	gl.DeleteBuffer(e.line.VBO)
}

func (e *Engine) Config(new, old event.Config) {
	e.touchLoc = geom.Point{new.Width / 2, new.Height / 2}
	e.camera.Perspective(0.785, float32(new.Width/new.Height), 0.1, 100.0)
}

func (e *Engine) Touch(t event.Touch, c event.Config) {
	e.touchLoc = t.Loc
}

func (e *Engine) Draw(c event.Config) {
	since := time.Now().Sub(e.started)

	gl.ClearColor(0, 0, 0, 1)
	//gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)
	//gl.Disable(gl.CULL_FACE)
	//gl.DepthFunc(gl.LESS)
	//gl.SampleCoverage(4.0, false)

	rotation := mgl.HomogRotate3D(float32(since.Seconds()), mgl.Vec3{0, 1, 0})
	e.scene.transform = &rotation

	if int(since.Seconds()) > e.offset {
		e.pos[1] += 1
		offset := len(e.line.vertices) / vertexDim
		e.line.vertices = append(e.line.vertices, Quad(e.pos[0], e.pos[1], e.pos[2], e.pos[3], e.pos[4])...)
		//e.shape.normals = append(e.shape.normals, []float32{1.0, 0.0, 0.0}...)
		e.line.Buffer(offset)
		e.offset++
	}

	//debug.DrawFPS(c)

	e.scene.Draw()

	/*
		if glErr := gl.GetError(); glErr != 0 {
			fmt.Println("glErr", glErr)
		}
	*/
}

func main() {
	camera := EulerCamera{}
	shader := Shader{}
	e := Engine{
		shader: &shader,
		camera: &camera,
		scene: &Scene{
			Camera: &camera,
			Shader: &shader,

			ambientColor: mgl.Vec3{0.5, 0.5, 0.5},
		},
	}

	app.Run(app.Callbacks{
		Start:  e.Start,
		Stop:   e.Stop,
		Draw:   e.Draw,
		Touch:  e.Touch,
		Config: e.Config,
	})

	log.SetOutput(os.Stdout)

}
