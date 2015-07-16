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
	camera   Camera
	shader   Shader
	shape    DynamicShape
	touchLoc geom.Point
	started  time.Time
	offset   int
	pos      [5]float32
	frames   int
}

func (e *Engine) Start() {
	var err error

	log.Println("Loading shaders...")

	e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
	if err != nil {
		panic(fmt.Sprintln("LoadProgram failed:", err))
	}
	e.pos = [5]float32{0, 0, 0, 1, 1}
	e.shape.vertices = append(e.shape.vertices, Quad(e.pos[0], e.pos[1], e.pos[2], e.pos[3], e.pos[4])...)
	//e.shape.normals = append(e.shape.normals, []float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}...)
	//e.shape.vertices = cubeData
	e.shape.Init(6 * 4 * 1000)
	e.shape.Buffer(0)

	gl.UseProgram(e.shader.program)
	e.shader.Bind()

	e.started = time.Now()

	log.Println("Starting.")
}

func (e *Engine) Stop() {
	gl.DeleteProgram(e.shader.program)
	gl.DeleteBuffer(e.shape.VBO)
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

	e.camera.Move(mgl.Vec3{3, 3, 3})
	e.camera.Pan(mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})

	// Setup MVP
	projection, view := e.camera.Projection(), e.camera.View()
	gl.UniformMatrix4fv(e.shader.projection, projection[:])
	gl.UniformMatrix4fv(e.shader.view, view[:])

	model := mgl.HomogRotate3D(float32(since.Seconds()), mgl.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(e.shader.model, model[:])

	normal := model.Mul4(view).Inv().Transpose()
	gl.UniformMatrix4fv(e.shader.normalMatrix, normal[:])

	// Light
	gl.Uniform3fv(e.shader.lightIntensities, []float32{1, 1, 1})
	gl.Uniform3fv(e.shader.lightPosition, []float32{1, 1, 1})

	if int(since.Seconds()) > e.offset {
		e.pos[1] += 1
		offset := len(e.shape.vertices) / vertexDim
		e.shape.vertices = append(e.shape.vertices, Quad(e.pos[0], e.pos[1], e.pos[2], e.pos[3], e.pos[4])...)
		//e.shape.normals = append(e.shape.normals, []float32{1.0, 0.0, 0.0}...)
		e.shape.Buffer(offset)
		e.offset++
	}

	//debug.DrawFPS(c)

	// Draw our shape
	e.shader.Draw(&e.shape)

	/*
		if glErr := gl.GetError(); glErr != 0 {
			fmt.Println("glErr", glErr)
		}
	*/
}

func main() {
	e := Engine{}
	app.Run(app.Callbacks{
		Start:  e.Start,
		Stop:   e.Stop,
		Draw:   e.Draw,
		Touch:  e.Touch,
		Config: e.Config,
	})

	log.SetOutput(os.Stdout)

}
