package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

// TODO: ...
// Ref: https://github.com/krux02/turnt-octo-wallhack/blob/master/particles/ParticleSystem.go

type Emitter interface {
	Shape
	Tick(time.Duration)
	MoveTo(mgl.Vec3)
}

func RandomParticle(origin mgl.Vec3, force float32) *particle {
	return &particle{
		position: origin,
		velocity: mgl.Vec3{
			(0.5 - rand.Float32()) * force,
			rand.Float32() * force,
			(0.5 - rand.Float32()) * force,
		},
	}

}

const particleLen = 9 * 3

type particle struct {
	position mgl.Vec3
	velocity mgl.Vec3
}

func (p *particle) Vertices() []float32 {
	var size float32 = 0.1
	pos := p.position
	a := pos.Add(mgl.Vec3{-size / 2, -size * 2, 0})
	b := pos.Add(mgl.Vec3{size / 2, -size, 0})
	return []float32{
		pos[0], pos[1], pos[2], // Top
		pos[0] - size, pos[1] - size, pos[2], // Bottom left
		pos[0] + size, pos[1] - size, pos[2], // Bottom right

		// Arrow handle
		b[0], b[1], b[2], // Top Right
		a[0], b[1], a[2], // Top Left
		a[0], a[1], a[2], // Bottom Left
		a[0], a[1], a[2], // Bottom Left
		b[0], b[1], b[2], // Top Right
		b[0], a[1], b[2], // Bottom Right
	}
}

func (p *particle) Tick(force mgl.Vec3) {
	p.velocity = force.Add(p.velocity)
	p.position = p.position.Add(p.velocity)
}

var particleForce float32 = 0.06
var gravityForce = mgl.Vec3{0, -0.1, 0}

func ParticleEmitter(origin mgl.Vec3, num int, rate float32) Emitter {
	bufSize := num * particleLen * vecSize
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferInit(gl.ARRAY_BUFFER, bufSize, gl.DYNAMIC_DRAW)

	return &particleEmitter{
		VBO:       vbo,
		origin:    origin,
		rate:      rate,
		particles: make([]*particle, 0, num),
		num:       num,
	}
}

type particleEmitter struct {
	VBO       gl.Buffer
	origin    mgl.Vec3
	rate      float32
	particles []*particle
	num       int
	lastTick  time.Duration
}

func (emitter *particleEmitter) MoveTo(pos mgl.Vec3) {
	emitter.origin = pos
}

func (emitter *particleEmitter) Tick(since time.Duration) {
	interval := float32((since - emitter.lastTick).Seconds())
	emitter.lastTick = since

	// Randomize emitting
	n := int(emitter.rate * interval * rand.Float32())

	extra := len(emitter.particles) + 1 + n - emitter.num
	if extra > 0 {
		// Pop oldest particles
		emitter.particles = emitter.particles[extra:]
	}

	for i := 0; i <= n; i++ {
		p := RandomParticle(emitter.origin, particleForce)
		emitter.particles = append(emitter.particles, p)
	}

	f := gravityForce.Mul(interval)
	for _, particle := range emitter.particles {
		particle.Tick(f)
	}

	emitter.Buffer()
}

func (emitter *particleEmitter) Buffer() {
	data := emitter.Bytes()
	if len(data) > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, emitter.VBO)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, data)
	}
}

func (emitter *particleEmitter) Len() int {
	return len(emitter.particles)
}

func (emitter *particleEmitter) Bytes() []byte {
	buf := bytes.Buffer{}
	for _, particle := range emitter.particles {
		binary.Write(&buf, binary.LittleEndian, particle.Vertices())
	}
	return buf.Bytes()
}

func (emitter *particleEmitter) Stride() int {
	return vecSize * vertexDim
}

func (emitter *particleEmitter) Draw(shader Shader, camera Camera) {
	gl.BindBuffer(gl.ARRAY_BUFFER, emitter.VBO)

	gl.EnableVertexAttribArray(shader.Attrib("vertCoord"))
	gl.VertexAttribPointer(shader.Attrib("vertCoord"), vertexDim, gl.FLOAT, false, emitter.Stride(), 0)

	gl.DrawArrays(gl.TRIANGLES, 0, emitter.Len()*particleLen/vertexDim)

	gl.DisableVertexAttribArray(shader.Attrib("vertCoord"))
}

func (emitter *particleEmitter) Close() {
	gl.DeleteBuffer(emitter.VBO)
}
