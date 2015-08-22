package collision

import (
	"image"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// LinearCollider is O(N) collision checker. Inefficient but simple and accurate.
func LinearCollider(bounds image.Rectangle) Collider {
	size := bounds.Size()

	collider := &linearCollider{
		bounds: Boundary{float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y)},
		width:  size.X,
		height: size.Y,
	}
	collider.Reset()
	return collider
}

type linearCollider struct {
	bounds   Boundary
	width    int
	height   int
	segments []*[]mgl.Vec3
}

func (collider *linearCollider) Track(segment *[]mgl.Vec3) Tracker {
	collider.segments = append(collider.segments, segment)
	return &linearTracker{
		collider: collider,
		segment:  segment,
	}
}

func (collider *linearCollider) Reset() {
	collider.segments = []*[]mgl.Vec3{}
}

func (collider *linearCollider) String() string {
	// TODO:
	return "<linearCollider>"
}

type linearTracker struct {
	collider *linearCollider
	segment  *[]mgl.Vec3
}

func (tracker *linearTracker) Update() error {
	collider := tracker.collider
	segment := *tracker.segment
	n := len(segment)
	if n < 2 {
		// Not long enough
		return nil
	}

	var vec mgl.Vec3 = segment[n-1]
	x1, y1 := vec[0], vec[2]

	// Check boundary
	if collider.bounds.X1 >= x1 || x1 >= collider.bounds.X2 ||
		collider.bounds.Y1 >= y1 || y1 >= collider.bounds.Y2 {
		return CollisionBoundary
	}

	vec = segment[n-2]
	x0, y0 := vec[0], vec[2]

	for _, segment_ref := range collider.segments {
		segment = *segment_ref
		m := len(segment)
		if segment_ref == tracker.segment {
			// Don't compare the last element for the current line
			m -= 1
		}
		for i := 1; i < m; i += 1 {
			vec = segment[i-1]
			seg_x0, seg_y0 := vec[0], vec[2]

			vec = segment[i]
			seg_x1, seg_y1 := vec[0], vec[2]

			if IsCollision2D(seg_x0, seg_y0, seg_x1, seg_y1, x0, y0, x1, y1) {
				return &CollisionSegment{seg_x0, seg_y0, seg_x1, seg_y1}
			}
		}
	}
	return nil
}
