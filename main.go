package main

import (
	"fmt"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type Engine struct {
	shader   Shader
	shape    Shape
	touchLoc geom.Point
	started  time.Time
	offset   int
	pos      [5]float32
	frames   int
}

func (e *Engine) Start() {
	var err error

	fmt.Println("Loading shaders...")

	e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
	if err != nil {
		panic(fmt.Sprintln("LoadProgram failed:", err))
	}
	e.pos = [5]float32{0, 0, 0, 1, 1}
	e.shape.vertices = append(e.shape.vertices, Quad(e.pos[0], e.pos[1], e.pos[2], e.pos[3], e.pos[4])...)
	//e.shape.normals = append(e.shape.normals, []float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}...)
	//e.shape.vertices = cubeData
	e.shape.VBO = gl.CreateBuffer()
	e.shape.Buffer()

	e.shader.Bind()

	e.started = time.Now()

	fmt.Println("Starting.")
}

func (e *Engine) Stop() {
	gl.DeleteProgram(e.shader.program)
	gl.DeleteBuffer(e.shape.VBO)
}

func (e *Engine) Config(new, old event.Config) {
	e.touchLoc = geom.Point{new.Width / 2, new.Height / 2}
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

	// Setup MVP
	var m mgl.Mat4
	m = mgl.Perspective(0.785, float32(c.Width/c.Height), 0.1, 10.0)
	gl.UniformMatrix4fv(e.shader.projection, m[:])

	m = mgl.LookAtV(
		mgl.Vec3{3, 3, 3}, // eye
		mgl.Vec3{0, 0, 0}, // center
		mgl.Vec3{0, 1, 0}, // up
	)
	gl.UniformMatrix4fv(e.shader.view, m[:])

	modelView := m

	m = mgl.HomogRotate3D(float32(since.Seconds()), mgl.Vec3{0, 1, 0})
	gl.UniformMatrix4fv(e.shader.model, m[:])

	m = m.Mul4(modelView).Inv().Transpose()
	gl.UniformMatrix4fv(e.shader.normalMatrix, m[:])

	// Light
	gl.Uniform3fv(e.shader.lightIntensities, []float32{1, 1, 1})
	gl.Uniform3fv(e.shader.lightPosition, []float32{1, 1, 1})

	if int(since.Seconds()) > e.offset && false {
		e.pos[1] += 0.1
		offset := len(e.shape.vertices) / vertexDim
		e.shape.vertices = append(e.shape.vertices, Quad(e.pos[0], e.pos[1], e.pos[2], e.pos[3], e.pos[4])...)
		//e.shape.normals = append(e.shape.normals, []float32{1.0, 0.0, 0.0}...)
		e.shape.BufferSub(offset)
		e.offset++
		//e.shape.Buffer()
	}

	//debug.DrawFPS(c)

	// Draw our shape
	e.shader.Draw(&e.shape)
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
}
