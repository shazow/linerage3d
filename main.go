package main

import (
	"fmt"
	"log"
	"math"
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

const mouseSensitivity = 0.01

func NewLine(shape *DynamicShape) *Line {
	return &Line{
		DynamicShape: shape,

		rate:   time.Second / 5.0,
		height: 1.0,
		length: 0.5,
	}
}

type Line struct {
	*DynamicShape

	lastMove time.Duration
	rate     time.Duration

	lastTurn mgl.Vec3
	position mgl.Vec3
	height   float32
	length   float32

	angle   float32
	turning bool
	offset  int
}

func (line *Line) Tick(since time.Duration, rotate float32) {
	if rotate != 0 {
		line.angle += rotate
		line.turning = true
	}

	if line.lastMove+line.rate > since {
		return
	}

	line.lastMove = since

	line.Add(line.angle)
	line.Buffer(line.offset)
}

func (line *Line) Add(angle float32) {
	turning := line.turning || line.angle != angle
	if turning {
		line.angle = angle
		line.lastTurn = line.position
	}

	// Rotate
	sin, cos := float32(math.Sin(float64(angle))), float32(math.Cos(float64(angle)))
	unit := mgl.Vec3{cos - sin, 0, sin + cos}

	// Normalize and reset height
	l := line.length / unit.Len()
	unit = mgl.Vec3{unit[0] * l, line.height, unit[2] * l}

	p1 := line.lastTurn
	p2 := line.position.Add(unit)
	quad := Quad(p1, p2)

	// Discard height and save
	p2[1] = 0
	line.position = p2

	if !turning && len(line.vertices) >= len(quad) {
		line.vertices = append(line.vertices[len(line.vertices)-len(quad):], quad...)
	} else {
		line.offset = len(line.vertices) / vertexDim
		line.vertices = append(line.vertices, quad...)
	}
	//e.shape.normals = append(e.shape.normals, []float32{1.0, 0.0, 0.0}...)

	line.turning = false
}

type Engine struct {
	camera *QuatCamera
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

	e.camera.MoveTo(mgl.Vec3{0, 10, -3})
	e.camera.RotateTo(mgl.Vec3{0, 0, 5})

	shape := &DynamicShape{}
	e.scene.nodes = append(e.scene.nodes, Node{Shape: shape})

	e.line = NewLine(shape)

	e.line.Add(0)
	e.line.Add(0)
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
	e.camera.SetPerspective(0.785, float32(new.Width/new.Height), 0.1, 100.0)
}

func (e *Engine) Touch(t touch.Event, c config.Event) {
	if t.Type == touch.TypeBegin {
		e.dragOrigin = t.Loc
		e.dragging = true
	} else if t.Type == touch.TypeEnd {
		e.dragging = false
		log.Println(e.camera.String())
	}
	e.touchLoc = t.Loc
	if e.dragging {
		deltaX, deltaY := float32(e.dragOrigin.X-e.touchLoc.X), float32(e.dragOrigin.Y-e.touchLoc.Y)
		e.camera.Rotate(mgl.Vec3{deltaY * mouseSensitivity, deltaX * mouseSensitivity, 0})
		e.dragOrigin = e.touchLoc
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
	rotation := mgl.HomogRotate3D(float32(since.Seconds()), AxisFront)
	e.scene.transform = &rotation

	e.line.Tick(since, 0.005)
	e.scene.Draw()

}

func main() {
	camera := NewQuatCamera()
	shader := &Shader{}
	engine := Engine{
		shader: shader,
		camera: camera,
		scene: &Scene{
			Camera: camera,
			Shader: shader,

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
