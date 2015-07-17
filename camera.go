package main

import (
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	View() mgl.Mat4
	Projection() mgl.Mat4
}

type EulerCamera struct {
	projection      mgl.Mat4
	eye, center, up mgl.Vec3
}

// Perspective computes the projection matrix and saves it
func (c *EulerCamera) Perspective(fovy, aspect, near, far float32) {
	c.projection = mgl.Perspective(0.785, aspect, 0.1, 100.0)
}

func (c *EulerCamera) Move(position mgl.Vec3) {
	c.eye = position
}

func (c *EulerCamera) Pan(target mgl.Vec3, up mgl.Vec3) {
	c.center = target
	c.up = up
}

func (c *EulerCamera) View() mgl.Mat4 {
	return mgl.LookAtV(c.eye, c.center, c.up)
}

func (c *EulerCamera) Projection() mgl.Mat4 {
	return c.projection
}

func (c *EulerCamera) String() string {
	return fmt.Sprintf("eye:    %+v\ncenter: %+v\nup:     %+v\n", c.eye, c.center, c.up)
}
