package main

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Vector interface {
	Position() mgl.Vec3
	Direction() mgl.Vec3
}

type World interface {
	Reset()
	Tick(time.Duration) error
	Focus() Vector
}
