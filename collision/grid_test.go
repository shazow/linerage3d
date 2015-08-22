package collision

import (
	"image"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type lineStep struct {
	collides bool
	vec      mgl.Vec3
	extend   bool
}

type stepSuite struct {
	name   string
	bounds image.Rectangle
	steps  []lineStep
}

var suites []stepSuite

func init() {
	bounds := image.Rect(-10, -10, 10, 10)
	suites = []stepSuite{
		{
			"Basic",
			bounds,
			[]lineStep{
				{false, mgl.Vec3{0, 0, 0}, false},
				{false, mgl.Vec3{1, 0, 0}, false},
				{false, mgl.Vec3{1.1, 0, 0}, false},
				{false, mgl.Vec3{1.2, 0, 0}, false},
				{false, mgl.Vec3{2, 0, 0}, false},
				{false, mgl.Vec3{3, 0, 0}, false},
				{false, mgl.Vec3{3, 0, 1}, false},
				{false, mgl.Vec3{3, 0, 2}, false},
				{false, mgl.Vec3{2, 0, 2}, false},
				{false, mgl.Vec3{2.5, 0, 1.5}, false},

				{true, mgl.Vec3{3.5, 0, 1.5}, false},
				{true, mgl.Vec3{2, 0, 1}, false},
				{false, mgl.Vec3{2.5, 0, 1}, false},
				{true, mgl.Vec3{4, 0, 1}, false},
			},
		},
		{
			"Extending",
			bounds,
			[]lineStep{
				{false, mgl.Vec3{3, 0, 3}, false},
				{false, mgl.Vec3{3, 0, 5}, true},
				{false, mgl.Vec3{3, 0, 7}, true},
				{false, mgl.Vec3{3.5, 0, 7}, false},   // turn
				{false, mgl.Vec3{3.5, 0, 7}, true},    // no-op
				{false, mgl.Vec3{4.5, 0, 7.5}, false}, // turn
				{false, mgl.Vec3{4.5, 0, 9}, true},    // extend
				{false, mgl.Vec3{4.5, 0, 9}, true},    // no-op
				{true, mgl.Vec3{4.5, 0, 10}, true},    // extend, boundary collision
			},
		},
		{
			"Extending 2",
			bounds,
			[]lineStep{
				{false, mgl.Vec3{0, 0, 0}, false},
				{false, mgl.Vec3{1, 0, 2}, false},
				{false, mgl.Vec3{0, 0, 3}, false},
				{false, mgl.Vec3{-1, 0, 2}, false},
				{false, mgl.Vec3{0, 0, 1}, false},
				{true, mgl.Vec3{0.4, 0, 0.6}, true},
			},
		},
	}
}

func colliderTester(t *testing.T, newCollider func(image.Rectangle) Collider) {
	for _, suite := range suites {
		t.Logf("Starting suite: %s", suite.name)

		collider := newCollider(suite.bounds)
		segments := []mgl.Vec3{}
		tracker := collider.Track(&segments)

		var err error
		for _, test := range suite.steps {
			t.Logf("adding %v", test.vec)
			segments = append(segments, test.vec)
			err = tracker.Update()
			if (err == nil) == test.collides {
				t.Errorf("expected collision=%v; got %s", test.collides, err)
			}
		}
	}
}

func TestGrid(t *testing.T) {
	colliderTester(t, GridCollider)
}

func TestLinear(t *testing.T) {
	colliderTester(t, LinearCollider)
}
