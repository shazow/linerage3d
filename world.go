package main

import (
	"time"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type World interface {
	Tick(time.Duration)
	Focus() mgl.Vec3
}
