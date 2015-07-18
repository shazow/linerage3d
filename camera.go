package main

import (
	"fmt"
	"math"

	mgl "github.com/go-gl/mathgl/mgl32"
)

const halfPi = math.Pi / 2.0

var worldUp = mgl.Vec3{0, 1, 0}

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
	// TODO: Read up on http://learnopengl.com/#!Getting-started/Camera
	// Borrowed from https://github.com/mmchugh/planetary/blob/master/src/helpers/camera.cpp
	c.center = mgl.Vec3{
		float32(math.Cos(c.pitch) * math.Cos(c.yaw)),
		float32(math.Sin(c.pitch)),
		float32(math.Cos(c.pitch) * math.Sin(c.yaw)),
	}.Normalize()

	c.right = c.center.Cross(worldUp).Normalize()
	c.up = c.right.Cross(c.center)
}

// Rotate adjusts the direction vectors by a delta vector of {pitch, yaw, roll}.
// Roll is ignored for now.
func (c *EulerCamera) Rotate(delta mgl.Vec3) {
	c.yaw += float64(delta.X())
	c.pitch += float64(delta.Y())

	// Limit vertical rotation to avoid gimbal lock
	if c.pitch > halfPi {
		c.pitch = halfPi
	} else if c.pitch < -halfPi {
		c.pitch = -halfPi
	}

	c.updateVectors()
}

// RotateTo adjusts the yaw and pitch to face a point.
func (c *EulerCamera) RotateTo(point mgl.Vec3) {
	// TODO: https://math.stackexchange.com/questions/470112/calculate-camera-pitch-yaw-to-face-point
}

// Move adjusts the position of the camera by a delta vector relative to the camera is facing.
func (c *EulerCamera) Move(delta mgl.Vec3) {
	c.eye = c.eye.Add(c.right.Mul(delta.X())).Add(c.up.Mul(delta.Y())).Add(c.center.Mul(delta.Z()))
}

// View returns the transform matrix from world space into camera space
func (c *EulerCamera) View() mgl.Mat4 {
	return mgl.LookAtV(c.eye, c.eye.Add(c.center), c.up)
}

// Projection returns the projection matrix for the camera perspective
func (c *EulerCamera) Projection() mgl.Mat4 {
	return c.projection
}

// String returns a string representation of the camera for debugging.
func (c *EulerCamera) String() string {
	return fmt.Sprintf(`	eye:%v
	center: %v
	up:     %v
	pitch, yaw: %v, %v`+"\n", c.eye, c.center, c.up, c.pitch, c.yaw)
}
