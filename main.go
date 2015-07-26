package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/config"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
)

const mouseSensitivity = 0.01
const turnSpeed = 0.05
const moveSpeed = 0.1

type Engine struct {
	camera   *QuatCamera
	scene    *Scene
	bindings *Bindings

	started time.Time

	line *Line

	touchLoc   geom.Point
	dragOrigin geom.Point
	dragging   bool
	paused     bool
	following  bool

	followOffset mgl.Vec3
}

func (e *Engine) Start() {
	var err error

	// Setup scene shader
	e.scene.shader, err = NewShader("shader.v.glsl", "shader.f.glsl")
	if err != nil {
		fail(1, "Failed to load shaders:", err)
	}

	// Setup skybox
	skyboxTex, err := LoadTextureCube("square.png")
	if err != nil {
		fail(1, "Failed to load texture:", err)
	}
	skyboxShader, err := NewShader("skybox.v.glsl", "skybox.f.glsl")
	if err != nil {
		fail(1, "Failed to load shaders:", err)
	}
	e.scene.skybox = NewSkybox(skyboxShader, skyboxTex)

	e.followOffset = mgl.Vec3{0, 10, -3}
	e.camera.MoveTo(e.followOffset)
	e.camera.RotateTo(mgl.Vec3{0, 0, 5})

	shape := NewDynamicShape(6 * 4 * 1000)
	e.line = NewLine(shape)
	e.line.Add(0)
	e.line.Add(0)
	e.line.Buffer(0)
	e.scene.nodes = append(e.scene.nodes, Node{Shape: shape})

	/*
		// Cube for funsies:
		cube := NewStaticShape()
		cube.vertices = skyboxVertices
		cube.normals = skyboxNormals
		cube.indices = skyboxIndices
		cube.Buffer()
		e.scene.nodes = append(e.scene.nodes, Node{Shape: cube})
	*/

	/*
		// Reflective floor
		e.scene.nodes = append(e.scene.nodes, Node{
			Shape: NewFloor(Node{Shape: shape}),
		})
	*/

	// Toggle keys
	e.bindings.On(KeyPause, func(_ KeyBinding) {
		e.paused = !e.paused
		log.Println("Paused:", e.paused)
	})
	e.bindings.On(KeyCameraFollow, func(_ KeyBinding) {
		e.following = !e.following
		log.Println("Following:", e.following)
	})

	e.started = time.Now()

	log.Println("Starting: ", e.scene.String())
}

func (e *Engine) Stop() {
	e.scene.shader.Close()
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
		log.Println("vertices=", e.line.Len(), e.camera.String())
	}
	e.touchLoc = t.Loc
	if e.dragging {
		deltaX, deltaY := float32(e.dragOrigin.X-e.touchLoc.X), float32(e.dragOrigin.Y-e.touchLoc.Y)
		e.camera.Rotate(mgl.Vec3{deltaY * mouseSensitivity, deltaX * mouseSensitivity, 0})
		e.dragOrigin = e.touchLoc
	}
}

func (e *Engine) Press(t key.Event, c config.Event) {
	switch t.Direction {
	case key.DirPress:
		e.bindings.Press(t.Code)
	case key.DirRelease:
		e.bindings.Release(t.Code)
	}
}

func (e *Engine) Draw(c config.Event) {
	since := time.Now().Sub(e.started)

	// Handle key presses
	var lineRotate float32
	var camDelta mgl.Vec3
	if e.bindings.Pressed(KeyLineLeft) {
		lineRotate -= turnSpeed
	}
	if e.bindings.Pressed(KeyLineRight) {
		lineRotate += turnSpeed
	}
	if e.bindings.Pressed(KeyCamForward) {
		camDelta[2] -= moveSpeed
	}
	if e.bindings.Pressed(KeyCamReverse) {
		camDelta[2] += moveSpeed
	}
	if e.bindings.Pressed(KeyCamLeft) {
		camDelta[0] -= moveSpeed
	}
	if e.bindings.Pressed(KeyCamRight) {
		camDelta[0] += moveSpeed
	}
	if e.bindings.Pressed(KeyCamUp) {
		e.camera.MoveTo(e.camera.Position().Add(mgl.Vec3{0, moveSpeed, 0}))
	}
	if e.bindings.Pressed(KeyCamDown) {
		e.camera.MoveTo(e.camera.Position().Add(mgl.Vec3{0, -moveSpeed, 0}))
	}
	if camDelta[0]+camDelta[1]+camDelta[2] != 0 {
		e.following = false
		e.camera.Move(camDelta)
	} else if e.following {
		e.camera.Lerp(e.line.position.Add(e.followOffset), e.line.position, 0.1)
	}

	gl.ClearColor(0, 0, 0, 1)
	//gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Enable(gl.DEPTH_TEST)

	//gl.Disable(gl.CULL_FACE)
	//gl.DepthFunc(gl.LESS)
	//gl.SampleCoverage(4.0, false)

	// Spinny!
	//rotation := mgl.HomogRotate3D(float32(since.Seconds()), AxisFront)
	//e.scene.transform = &rotation

	if !e.paused {
		e.line.Tick(since, lineRotate)
	}
	e.scene.Draw(e.camera)

}

func fail(code int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(code)
}

func main() {
	log.SetOutput(os.Stdout)

	camera := NewQuatCamera()
	engine := Engine{
		camera:   camera,
		bindings: DefaultBindings(),
		scene: &Scene{
			ambientColor: mgl.Vec3{0.5, 0.5, 0.5},
		},
	}

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
			case key.Event:
				engine.Press(e, c)
			}
		}
	})
}
