package collision

import (
	"errors"
	"fmt"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type Tracker interface {
	Update() error
}

type Collider interface {
	Track(*[]mgl.Vec3) Tracker
	Reset()
	String() string
}

var CollisionBoundary = errors.New("collision with boundary")

type CollisionSegment struct {
	x0, y0, x1, y1 float32
}

func (seg *CollisionSegment) Error() string {
	return fmt.Sprintf("collision with segment: %v,%v -> %v,%v", seg.x0, seg.y0, seg.x1, seg.y1)
}

type Boundary struct {
	X1, Y1, X2, Y2 float32
}
