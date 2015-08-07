package main

import (
	"errors"
	"image"
	"log"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

func NewArenaNode(bounds image.Rectangle, shader Shader) *arena {
	shape := NewStaticShape()
	shape.vertices = floorVertices
	shape.Buffer()

	return newArena(bounds, &Node{
		Shape:  shape,
		shader: shader,
	})
}

func newArena(bounds image.Rectangle, node *Node) *arena {
	size := bounds.Size()

	return &arena{
		Node:   node,
		bounds: f32bounds{float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y)},
		width:  size.X,
		grid:   make([][]cellSegment, size.X*size.Y, size.X*size.Y),
	}
}

type cellSegment struct {
	cellIdx    int
	segmentIdx int
	offset     int
	segment    []mgl.Vec3
}

type f32bounds struct {
	X1, Y1, X2, Y2 float32
}

type arena struct {
	*Node

	bounds f32bounds
	width  int
	grid   [][]cellSegment
}

func (grid *arena) index(x, y float32) int {
	return int(x-grid.bounds.X1) + int(y-grid.bounds.Y1)*grid.width
}

var ErrBadSegment = errors.New("insufficient segment length")
var CollisionBoundary = errors.New("collision with boundary")
var CollisionSegment = errors.New("collision with segment")

func (grid *arena) Add(cs *cellSegment, segment []mgl.Vec3) (*cellSegment, error) {
	if len(segment) < 2 {
		return nil, ErrBadSegment
	}

	offset := len(segment) - 1

	var vec mgl.Vec3 = segment[offset]
	x1, y1 := vec[0], vec[2]

	// Check boundary
	if grid.bounds.X1 >= x1 || x1 >= grid.bounds.X2 ||
		grid.bounds.Y1 >= y1 || y1 >= grid.bounds.Y2 {
		return nil, CollisionBoundary
	}

	idx := grid.index(x1, y1)
	cell := &grid.grid[idx]

	// Check segment collision
	vec = segment[offset-1]
	x0, y0 := vec[0], vec[2]

	// Check cell segments
	for _, cellSegment := range *cell {
		n := len(cellSegment.segment)
		if cellSegment.cellIdx == idx && n > 0 {
			// Same segment, skip comparing the last element
			n -= 1
		}
		for i := 1; i < n; i += 2 {
			vec = cellSegment.segment[i-1]
			seg_x0, seg_y0 := vec[0], vec[2]

			vec = cellSegment.segment[i]
			seg_x1, seg_y1 := vec[0], vec[2]

			if IsCollision2D(x0, y0, x1, y1, seg_x0, seg_y0, seg_x1, seg_y1) {
				log.Printf("Collision: %v against %v", []float32{x0, y0, x1, y1}, []float32{seg_x0, seg_y0, seg_x1, seg_y1})
				return nil, CollisionSegment
			}
		}
	}

	if cs == nil || cs.cellIdx != idx {
		if cs != nil {
			// Append final piece to the original cellSegment
			cs.segment = segment[cs.offset:]
		}
		*cell = append(*cell, cellSegment{
			cellIdx:    idx,
			segmentIdx: len(*cell),
			offset:     offset,
		})
		cs = &grid.grid[idx][len(*cell)-1]
	}
	cs.segment = segment[cs.offset:]
	return cs, nil
}

func (grid *arena) IsCollision(segment []mgl.Vec3) bool {
	if len(segment) < 2 {
		return false
	}

	var vec mgl.Vec3 = segment[len(segment)-1]
	a1_x, a1_y := vec[0], vec[2]
	if grid.bounds.X1 >= a1_x || a1_x >= grid.bounds.X2 ||
		grid.bounds.Y1 >= a1_y || a1_y >= grid.bounds.Y2 {
		return true
	}

	vec = segment[len(segment)-2]
	a2_x, a2_y := vec[0], vec[2]
	idx := grid.index(a1_x, a1_y)

	// TODO: Check cellSegments of the other idx?
	for _, cellSegment := range grid.grid[idx] {
		for i := 1; i < len(cellSegment.segment); i += 2 {
			vec = cellSegment.segment[i-1]
			b1_x, b1_y := vec[0], vec[2]

			vec = cellSegment.segment[i]
			b2_x, b2_y := vec[0], vec[2]

			if IsCollision2D(a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y) {
				log.Printf("Collision: %v", []float32{a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y})
				return true
			}
		}
	}

	return false
}

func (shape *arena) Draw(camera Camera) {
	shader := shape.shader
	gl.Uniform3fv(shader.Uniform("material.ambient"), []float32{0.2, 0.2, 0.2})
	gl.Uniform3fv(shader.Uniform("lights[0].color"), []float32{0.2, 0.1, 0.1})
	gl.Uniform1f(shader.Uniform("lights[0].intensity"), 0.3)

	shape.Node.Draw(camera)
}

var arenaVertices = []float32{
	// Floor
	-100, 0, -100,
	100, 0, -100,
	100, 0, 100,
	100, 0, 100,
	-100, 0, -100,
	-100, 0, 100,
}
