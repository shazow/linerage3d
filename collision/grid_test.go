package collision

import (
	"image"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func TestGrid(t *testing.T) {
	collider := GridCollider(image.Rect(-10, -10, 10, 10))

	segments := []mgl.Vec3{}
	tracker := collider.Track(&segments)

	tests := []struct {
		collides bool
		vec      mgl.Vec3
	}{
		{false, mgl.Vec3{0, 0, 0}},
		{false, mgl.Vec3{1, 0, 0}},
		{false, mgl.Vec3{1.1, 0, 0}},
		{false, mgl.Vec3{1.2, 0, 0}},
		{false, mgl.Vec3{2, 0, 0}},
		{false, mgl.Vec3{3, 0, 0}},
		{false, mgl.Vec3{3, 0, 1}},
		{false, mgl.Vec3{3, 0, 2}},
		{false, mgl.Vec3{2, 0, 2}},
		{false, mgl.Vec3{2.5, 0, 1.5}},

		{true, mgl.Vec3{3.5, 0, 1.5}},
		{true, mgl.Vec3{2, 0, 1}},
		{false, mgl.Vec3{2.5, 0, 1}},
		{true, mgl.Vec3{4, 0, 1}},
	}

	var err error
	for _, test := range tests {
		t.Logf("Adding %v", test.vec)
		segments = append(segments, test.vec)
		err = tracker.Update()
		if (err == nil) == test.collides {
			t.Errorf("expected collision=%v; got %s", test.collides, err)
		}
	}
}

func TestGridExtend2(t *testing.T) {
	collider := GridCollider(image.Rect(-10, -10, 10, 10))
	segments := []mgl.Vec3{}
	tracker := collider.Track(&segments)

	tests := []struct {
		collides bool
		vec      mgl.Vec3
		extend   bool
	}{
		{false, mgl.Vec3{0, 0, 0}, false},
		{false, mgl.Vec3{1, 0, 2}, false},
		{false, mgl.Vec3{0, 0, 3}, false},
		{false, mgl.Vec3{-1, 0, 2}, false},
		{false, mgl.Vec3{0, 0, 1}, false},
		{true, mgl.Vec3{0.4, 0, 0.6}, true},
	}

	var err error
	for _, test := range tests {
		t.Logf("Adding %v", test.vec)
		if test.extend {
			segments[len(segments)-1] = test.vec
		} else {
			segments = append(segments, test.vec)
		}
		err = tracker.Update()
		if (err == nil) == test.collides {
			t.Errorf("expected collision=%v; got %s", test.collides, err)
		}
	}
}

func TestGridExtend(t *testing.T) {
	collider := GridCollider(image.Rect(-10, -10, 10, 10))
	segments := []mgl.Vec3{}
	tracker := collider.Track(&segments)

	tests := []struct {
		collides bool
		vec      mgl.Vec3
		extend   bool
	}{
		{false, mgl.Vec3{3, 0, 3}, false},
		{false, mgl.Vec3{3, 0, 5}, true},
		{false, mgl.Vec3{3, 0, 7}, true},
		{false, mgl.Vec3{3.5, 0, 7}, false},   // turn
		{false, mgl.Vec3{3.5, 0, 7}, true},    // no-op
		{false, mgl.Vec3{4.5, 0, 7.5}, false}, // turn
		{false, mgl.Vec3{4.5, 0, 9}, true},    // extend
		{false, mgl.Vec3{4.5, 0, 9}, true},    // no-op
		{true, mgl.Vec3{4.5, 0, 10}, true},    // extend, boundary collision
	}

	var err error
	for _, test := range tests {
		t.Logf("Adding %v", test.vec)
		if test.extend {
			segments[len(segments)-1] = test.vec
		} else {
			segments = append(segments, test.vec)
		}
		err = tracker.Update()
		if (err == nil) == test.collides {
			t.Errorf("expected collision=%v; got %s", test.collides, err)
		}
	}
}
