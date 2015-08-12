package main

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type World interface {
	Reset()
	Tick(time.Duration) error
	Focus() mgl.Vec3
}
