package main

import (
	"errors"
	"image"

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
		height: size.Y,
		grid:   make([]gridCell, size.X*size.Y, size.X*size.Y),
	}
}

type f32bounds struct {
	X1, Y1, X2, Y2 float32
}

type cellSegment struct {
	cellIdx    int
	segmentIdx int
	offset     int
	segment    []mgl.Vec3
}

type gridCell []cellSegment

func (cell *gridCell) Len() int {
	n := 0
	for _, cs := range *cell {
		n += len(cs.segment)
	}
	return n
}

func (cell *gridCell) IsCollision(x0, y0, x1, y1 float32) bool {
	var vec mgl.Vec3
	for _, cellSegment := range *cell {
		for i := 1; i < len(cellSegment.segment); i += 1 {
			vec = cellSegment.segment[i-1]
			seg_x0, seg_y0 := vec[0], vec[2]

			vec = cellSegment.segment[i]
			seg_x1, seg_y1 := vec[0], vec[2]

			if IsCollision2D(x0, y0, x1, y1, seg_x0, seg_y0, seg_x1, seg_y1) {
				return true
			}
		}
	}
	return false
}

type arena struct {
	*Node

	bounds f32bounds
	width  int
	height int
	grid   []gridCell
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
	if grid.bounds.X1 > x1 || x1 > grid.bounds.X2 ||
		grid.bounds.Y1 > y1 || y1 > grid.bounds.Y2 {
		return nil, CollisionBoundary
	}

	// Check segment collision
	vec = segment[offset-1]
	x0, y0 := vec[0], vec[2]

	checkIdx := []int{grid.index(x0, y0), grid.index(x1, y1)}

	if checkIdx[0] == checkIdx[1] {
		checkIdx = checkIdx[:1]
	}

	// Check cell segments
	var err error
	for _, idx := range checkIdx {
		if grid.grid[idx].IsCollision(x0, y0, x1, y1) {
			err = CollisionSegment
			break
		}
	}

	for _, idx := range checkIdx {
		cell := &grid.grid[idx]
		if cs == nil || cs.cellIdx != idx {
			*cell = append(*cell, cellSegment{
				cellIdx:    idx,
				segmentIdx: len(*cell),
				offset:     offset - 1,
			})
			cs = &grid.grid[idx][len(*cell)-1]
		}
		cs.segment = segment[cs.offset:]
	}
	return cs, err
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
