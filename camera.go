package main

import (
	"fmt"
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

const halfPi = math.Pi / 2.0

type Camera interface {
	View() mgl.Mat4
	Projection() mgl.Mat4
}

type EulerCamera struct {
	projection mgl.Mat4
	eye        mgl.Vec3

	yaw   float64
	pitch float64

	center mgl.Vec3
	up     mgl.Vec3
	right  mgl.Vec3
}

// Perspective computes the projection matrix and saves it
func (c *EulerCamera) Perspective(fovy, aspect, near, far float32) {
	c.projection = mgl.Perspective(fovy, aspect, near, far)
}

func (c *EulerCamera) updateVectors() {
	// Borrowed from https://github.com/mmchugh/planetary/blob/master/src/helpers/camera.cpp
	c.center = mgl.Vec3{
		float32(math.Cos(c.pitch) * math.Sin(c.yaw)),
		float32(math.Sin(c.pitch)),
		float32(math.Cos(c.pitch) * math.Cos(c.yaw)),
	}

	c.right = mgl.Vec3{
		float32(math.Sin(c.yaw - halfPi)),
		0,
		float32(math.Cos(c.pitch - halfPi)),
	}

	c.up = c.right.Cross(c.center)
}

// Rotate adjusts the direction vectors by a delta vector of {pitch, yaw, roll}.
// Roll is ignored for now.
func (c *EulerCamera) Rotate(delta mgl.Vec3) {
	c.pitch += float64(delta[0])
	c.yaw += float64(delta[1])

	// Limit vertical rotation to avoid gimbal lock
	if c.yaw > halfPi {
		c.yaw = halfPi
	} else if c.yaw < -halfPi {
		c.yaw = -halfPi
	}

	c.updateVectors()
}

// Move adjusts the position of the camera by a delta vector relative to the camera is facing.
func (c *EulerCamera) Move(delta mgl.Vec3) {
	c.eye = c.eye.Add(c.right.Mul(delta.X())).Add(c.up.Mul(delta.Y())).Add(c.center.Mul(delta.Z()))
}

// View returns the transform matrix from world space into camera space
func (c *EulerCamera) View() mgl.Mat4 {
	return mgl.LookAtV(c.eye, c.center, c.up)
}

// Projection returns the projection matrix for the camera perspective
func (c *EulerCamera) Projection() mgl.Mat4 {
	return c.projection
}

// String returns a string representation of the camera for debugging.
func (c *EulerCamera) String() string {
	return fmt.Sprintf("eye:    %+v\ncenter: %+v\nup:     %+v\n", c.eye, c.center, c.up)
}
