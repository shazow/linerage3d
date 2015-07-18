package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/config"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

const mouseSensitivity = 180

type Line struct {
	*DynamicShape

	offset int
	pos    [5]float32

	lastMove time.Duration
	rate     time.Duration
}

func (line *Line) Tick(since time.Duration) {
	if line.lastMove+line.rate > since {
		return
	}

	line.lastMove = since
	line.pos[1] += 1
	offset := len(line.vertices) / vertexDim
	line.vertices = append(line.vertices, Quad(line.pos[0], line.pos[1], line.pos[2], line.pos[3], line.pos[4])...)
	//e.shape.normals = append(e.shape.normals, []float32{1.0, 0.0, 0.0}...)
	line.Buffer(offset)
	line.offset++
}

type Engine struct {
	camera *EulerCamera
	scene  *Scene
	shader *Shader

	started time.Time

	line *Line

	touchLoc   geom.Point
	dragOrigin geom.Point
	dragging   bool
}

func (e *Engine) Start() {
	var err error

	log.Println("Loading shaders...")

	e.shader.program, err = LoadProgram("shader.v.glsl", "shader.f.glsl")
	if err != nil {
		panic(fmt.Sprintln("LoadProgram failed:", err))
	}

	e.camera.Rotate(mgl.Vec3{0, 0, 0})
	e.camera.Move(mgl.Vec3{30, 30, 30})

	shape := &DynamicShape{}
	e.scene.nodes = append(e.scene.nodes, Node{Shape: shape})

	e.line = &Line{
		DynamicShape: shape,
		rate:         time.Second / 10.0,
	}

	pos := [5]float32{0, 0, 0, 1, 1}
	shape.vertices = append(shape.vertices, Quad(pos[0], pos[1], pos[2], pos[3], pos[4])...)
	//shape.normals = append(shape.normals, []float32{0.0, 0.0, 0.0, 0.0, 0.0, 0.0}...)
	//shape.vertices = cubeData

	e.line.pos = pos
	e.line.Init(6 * 4 * 1000)
	e.line.Buffer(0)

	gl.UseProgram(e.shader.program)
	e.scene.Bind()

	e.started = time.Now()

	log.Println("Starting: ", e.scene.String())
}

func (e *Engine) Stop() {
	gl.DeleteProgram(e.shader.program)
	gl.DeleteBuffer(e.line.VBO)
}

func (e *Engine) Config(new, old config.Event) {
	e.touchLoc = geom.Point{new.Width / 2, new.Height / 2}
	e.camera.Perspective(0.785, float32(new.Width/new.Height), 0.1, 100.0)
}

func (e *Engine) Touch(t touch.Event, c config.Event) {
	if t.Type == touch.TypeStart {
		e.dragOrigin = t.Loc
		e.dragging = true
	} else if t.Type == touch.TypeEnd {
		e.dragging = false
	}
	e.touchLoc = t.Loc
	if e.dragging {
		deltaX, deltaY := float32(e.dragOrigin.X-e.touchLoc.X), float32(e.dragOrigin.Y-e.touchLoc.Y)
		e.camera.Rotate(mgl.Vec3{deltaX / mouseSensitivity, -deltaY / mouseSensitivity, 0})
		e.dragOrigin = t.Loc
	}
}

func (e *Engine) Draw(c config.Event) {
	since := time.Now().Sub(e.started)

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	//gl.Enable(gl.DEPTH_TEST)

	//gl.Disable(gl.CULL_FACE)
	//gl.DepthFunc(gl.LESS)
	//gl.SampleCoverage(4.0, false)

	// Spinny!
	rotation := mgl.HomogRotate3D(float32(since.Seconds()), mgl.Vec3{0, 1, 0})
	e.scene.transform = &rotation

	e.line.Tick(since)
	e.scene.Draw()

}

func main() {
	camera := EulerCamera{}
	shader := Shader{}
	engine := Engine{
		shader: &shader,
		camera: &camera,
		scene: &Scene{
			Camera: &camera,
			Shader: &shader,

			ambientColor: mgl.Vec3{0.5, 0.5, 0.5},
		},
	}

	log.SetOutput(os.Stdout)

	app.Main(func(a app.App) {
		var c config.Event
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					engine.Start()
				case lifecycle.CrossOff:
					engine.Stop()
				}
			case config.Event:
				engine.Config(e, c)
				c = e
			case paint.Event:
				engine.Draw(c)
				a.EndPaint()
			case touch.Event:
				engine.Touch(e, c)
			}
		}
	})
}
