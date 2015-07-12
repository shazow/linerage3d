package main

import (
	"fmt"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

type Shape struct {
	buf gl.Buffer

	vertexCount int
}

type Shader struct {
	program      gl.Program
	vertCoord    gl.Attrib
	vertNormal   gl.Attrib
	vertTexCoord gl.Attrib

	projection   gl.Uniform
	view         gl.Uniform
	model        gl.Uniform
	normalMatrix gl.Uniform

	lightPosition    gl.Uniform
	lightIntensities gl.Uniform
}

type Engine struct {
	shader   Shader
	shape    Shape
	touchLoc geom.Point
	started  time.Time
}

func (e *Engine) Start() {
	var err error

	fmt.Println("Loading shaders...")

	e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
	if err != nil {
		panic(fmt.Sprintln("LoadProgram failed:", err))
	}

	e.shape.vertexCount = len(cubeData) / (3 + 3)

	e.shape.buf = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)
	gl.BufferData(gl.ARRAY_BUFFER, EncodeObject(cubeData), gl.STATIC_DRAW)

	e.shader.vertCoord = gl.GetAttribLocation(e.shader.program, "vertCoord")

	e.shader.projection = gl.GetUniformLocation(e.shader.program, "projection")
	e.shader.view = gl.GetUniformLocation(e.shader.program, "view")
	e.shader.model = gl.GetUniformLocation(e.shader.program, "model")
	e.shader.normalMatrix = gl.GetUniformLocation(e.shader.program, "normalMatrix")

	e.shader.lightIntensities = gl.GetUniformLocation(e.shader.program, "lightIntensities")
	e.shader.lightPosition = gl.GetUniformLocation(e.shader.program, "lightPosition")

	e.started = time.Now()

	fmt.Println("Starting.")
}

func (e *Engine) Stop() {
	gl.DeleteProgram(e.shader.program)
	gl.DeleteBuffer(e.shape.buf)
}

func (e *Engine) Config(new, old event.Config) {
	e.touchLoc = geom.Point{new.Width / 2, new.Height / 2}
}

func (e *Engine) Touch(t event.Touch, c event.Config) {
	e.touchLoc = t.Loc
}

func (e *Engine) Draw(c event.Config) {
	since := time.Now().Sub(e.started)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(e.shader.program)

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

	// Draw our shape
	gl.BindBuffer(gl.ARRAY_BUFFER, e.shape.buf)

	stride := 4 * (3 + 3) // 4 bytes in float, (3+2+3) values per vertex
	gl.EnableVertexAttribArray(e.shader.vertCoord)
	gl.VertexAttribPointer(e.shader.vertCoord, 3, gl.FLOAT, false, stride, 0)

	gl.EnableVertexAttribArray(e.shader.vertNormal)
	gl.VertexAttribPointer(e.shader.vertNormal, 3, gl.FLOAT, false, stride, 12)

	gl.DrawArrays(gl.TRIANGLES, 0, e.shape.vertexCount)

	gl.DisableVertexAttribArray(e.shader.vertCoord)

	debug.DrawFPS(c)
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
