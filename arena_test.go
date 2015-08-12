package main

import (
	"image"
	"reflect"
	"testing"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func TestArenaSegments(t *testing.T) {
	arena := newArena(image.Rect(-10, -10, 10, 10), nil)
	segments := []mgl.Vec3{
		mgl.Vec3{3, 0, 1},
		mgl.Vec3{3, 0, 2},
		mgl.Vec3{2.5, 0, 1.5},
		mgl.Vec3{3.5, 0, 1.5},
	}

	offset := 2
	var cs *cellSegment
	var err error

	if cs, err = arena.Add(cs, segments[:offset]); err != nil {
		t.Errorf("Unexpected collision for segment: %v", segments[offset-2:offset])
		t.Log(dumpArena(arena))
	}

	if got, want := arena.grid[arena.index(3, 2)].Len(), 2; want != got {
		t.Errorf("[3,2] segment length: got %v; want %v", got, want)
	}

	if got, want := arena.grid[arena.index(3, 1)].Len(), 2; want != got {
		t.Errorf("[3,1] segment length: got %v; want %v", got, want)
	}

	offset = 4
	if cs, err = arena.Add(cs, segments[:offset]); err == nil {
		t.Errorf("Missed collision for segment: %v", segments[offset-2:offset])
		t.Log(dumpArena(arena))
	}
}

func TestArenaExtend(t *testing.T) {
	arena := newArena(image.Rect(-10, -10, 10, 10), nil)
	segments := []mgl.Vec3{
		mgl.Vec3{3, 0, 1},
		mgl.Vec3{3, 0, 2},
		mgl.Vec3{3, 0, 3},
	}

	offset := 2
	var cs *cellSegment
	var err error

	if cs, err = arena.Add(cs, segments[:offset]); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	segments[len(segments)-1][2] = 5
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}
	oldIdx := cs.cellIdx

	segments[len(segments)-1][2] = 7
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	if cs.cellIdx == oldIdx {
		t.Errorf("Failed to update cell segment")
	}

	// Turn
	segments = append(segments, mgl.Vec3{3.5, 0, 7})
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	// No-op
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	// Turn
	segments = append(segments, mgl.Vec3{4.5, 0, 7.5})
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	// Extend
	segments[len(segments)-1][2] = 9
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	// No-op
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Unexpected collision for segment: %v", err)
	}

	segments[len(segments)-1][2] = 10
	if cs, err = arena.Add(cs, segments); err == nil {
		t.Errorf("Missed collision with boundary.")
	}
}

func TestArenaFail(t *testing.T) {
	arena := newArena(image.Rect(-10, -10, 10, 10), nil)
	segments := []mgl.Vec3{
		mgl.Vec3{0, 0, 0},
		mgl.Vec3{1, 0, 2},
		mgl.Vec3{0, 0, 3},
		mgl.Vec3{-1, 0, 2},
		mgl.Vec3{0, 0, 1},
	}

	var cs *cellSegment
	var err error
	for i := 2; i <= len(segments); i++ {
		t.Logf("Adding %v", segments[i-2:i])
		cs, err = arena.Add(cs, segments[:i])
		if err != nil {
			t.Errorf("Arena collision on segment %v: %s", segments[i-2:i], err)
		}
	}

	// Extend
	segments[len(segments)-1] = mgl.Vec3{0.4, 0, 0.6}
	if cs, err = arena.Add(cs, segments); err == nil {
		t.Errorf("Missed collision for segment: %v", segments[len(segments)-2:])
	}
}

func TestArena(t *testing.T) {
	arena := newArena(image.Rect(-10, -10, 10, 10), nil)

	if !reflect.DeepEqual(arena.bounds, f32bounds{X1: -10, Y1: -10, X2: 10, Y2: 10}) {
		t.Errorf("incorrect arena.bounds: %v", arena.bounds)
	}

	if got, want := arena.width, 20; want != got {
		t.Errorf("got %v; want %v", got, want)
	}

	if got, want := arena.index(-10, -10), 0; want != got {
		t.Errorf("got %v; want %v", got, want)
	}

	if got, want := arena.index(-9, -9), 21; want != got {
		t.Errorf("got %v; want %v", got, want)
	}

	segments := []mgl.Vec3{
		mgl.Vec3{0, 0, 0},
		mgl.Vec3{1, 0, 0},
		mgl.Vec3{1.1, 0, 0},
		mgl.Vec3{1.2, 0, 0},
		mgl.Vec3{2, 0, 0},
		mgl.Vec3{3, 0, 0},
		mgl.Vec3{3, 0, 1},
		mgl.Vec3{3, 0, 2},
		mgl.Vec3{2, 0, 2},
		mgl.Vec3{2.5, 0, 1.5},
	}

	var cs *cellSegment
	var err error
	for i := 2; i <= len(segments); i++ {
		t.Logf("Adding %v", segments[i-2:i])
		cs, err = arena.Add(cs, segments[:i])
		if err != nil {
			t.Errorf("Arena collision on segment %v: %s", segments[i-2:i], err)
		}
	}

	// Add collision segment
	segments = append(segments, mgl.Vec3{3.5, 0, 1.5})
	t.Logf("Adding %v", segments[len(segments)-2:])
	if cs, err = arena.Add(cs, segments); err == nil {
		t.Errorf("Missed collision for segment: %v", segments[len(segments)-2:])
		t.Log(dumpArena(arena))
	}

	// Add collision segment
	segments = append(segments, mgl.Vec3{2, 0, 1})
	t.Logf("Adding %v", segments[len(segments)-2:])
	if cs, err = arena.Add(cs, segments); err == nil {
		t.Errorf("Missed collision for segment: %v", segments[len(segments)-2:])
	}

	// Add non-collision segment
	segments = append(segments, mgl.Vec3{2.5, 0, 1})
	t.Logf("Adding %v", segments[len(segments)-2:])
	if cs, err = arena.Add(cs, segments); err != nil {
		t.Errorf("Collision for segment: %v", err)
	}

	// Add collision segment
	segments = append(segments, mgl.Vec3{4, 0, 1})
	t.Logf("Adding %v", segments[len(segments)-2:])
	if cs, err = arena.Add(cs, segments); err == nil {
		t.Errorf("Missed collision for segment: %v", segments[len(segments)-2:])
	}
}
