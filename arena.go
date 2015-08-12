package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"

	mgl "github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/gl"
)

var ErrBadSegment = errors.New("insufficient segment length")
var CollisionBoundary = errors.New("collision with boundary")

type CollisionSegment struct {
	x0, y0, x1, y1 float32
}

func (seg *CollisionSegment) Error() string {
	return fmt.Sprintf("collision with segment: %v,%v -> %v,%v", seg.x0, seg.y0, seg.x1, seg.y1)
}

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

	grid := &arena{
		Node:   node,
		bounds: f32bounds{float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y)},
		width:  size.X,
		height: size.Y,
	}
	grid.Reset()
	return grid
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

func (cell *gridCell) IsCollision(x0, y0, x1, y1 float32) error {
	var vec mgl.Vec3
	for _, cellSegment := range *cell {
		for i := 1; i < len(cellSegment.segment); i += 1 {
			vec = cellSegment.segment[i-1]
			seg_x0, seg_y0 := vec[0], vec[2]

			vec = cellSegment.segment[i]
			seg_x1, seg_y1 := vec[0], vec[2]

			if IsCollision2D(seg_x0, seg_y0, seg_x1, seg_y1, x0, y0, x1, y1) {
				if seg_x0 == x0 && seg_y0 == y0 && seg_x1 == x1 && seg_y1 == y1 {
					// FIXME: Temporary hack to avoid checking the same segment
					continue
				}
				return &CollisionSegment{seg_x0, seg_y0, seg_x1, seg_y1}
			}
		}
	}
	return nil
}

type arena struct {
	*Node

	bounds f32bounds
	width  int
	height int
	grid   []gridCell

	// Debug
	checked int
}

func (grid *arena) Reset() {
	grid.grid = make([]gridCell, grid.width*grid.height, grid.width*grid.height)
}

func (grid *arena) index(x, y float32) int {
	return int(x-grid.bounds.X1) + int(y-grid.bounds.Y1)*grid.width
}

func (grid *arena) Add(cs *cellSegment, segment []mgl.Vec3) (*cellSegment, error) {
	if len(segment) < 2 {
		return nil, ErrBadSegment
	}

	offset := len(segment) - 2

	var vec mgl.Vec3 = segment[offset+1]
	x1, y1 := vec[0], vec[2]

	// Check boundary
	if grid.bounds.X1 >= x1 || x1 >= grid.bounds.X2 ||
		grid.bounds.Y1 >= y1 || y1 >= grid.bounds.Y2 {
		return nil, CollisionBoundary
	}

	// Check segment collision
	vec = segment[offset]
	x0, y0 := vec[0], vec[2]

	dx, dy := x1-x0, y1-y0
	steps := (dx*dx + dy*dy)
	if steps > 1.0 {
		dx, dy = dx/steps, dy/steps
	}

	// Extending only works if the extension is collinearly forward. All hell
	// will break loose otherwise.
	extending := cs != nil && cs.offset == offset

	// Check cell segments
	var lastIdx int = -1
	var err error
	for ; steps >= 0; steps -= 1.0 {
		idx := grid.index(x1-dx*steps, y1-dy*steps)
		if idx == lastIdx {
			// FIXME: This is happening more than it should
			continue
		}
		if idx > len(grid.grid) && err == nil {
			err = CollisionBoundary
			break
		}
		lastIdx = idx

		if extending && false {
			// Skip until we reach the last cell we filled last
			if cs.cellIdx == idx {
				extending = false
			} else {
				continue
			}
		}

		if cs != nil && cs.cellIdx == idx && len(cs.segment) == len(segment)-offset+1 {
			// Already set
			continue
		}

		if err == nil {
			err = grid.grid[idx].IsCollision(x0, y0, x1, y1)
			grid.checked++
		}

		cell := &grid.grid[idx]
		if cs == nil || cs.cellIdx != idx {
			*cell = append(*cell, cellSegment{
				cellIdx:    idx,
				segmentIdx: len(*cell),
				offset:     offset,
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

func dumpArena(arena *arena) string {
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "arena.bounds: %+v\n", arena.bounds)
	fmt.Fprintf(w, "arena.grid:\n")

	for idx, cell := range arena.grid {
		if len(cell) == 0 {
			continue
		}

		x := idx%arena.width + int(arena.bounds.X1)
		y := (idx-x)/arena.height + int(arena.bounds.Y1)

		fmt.Fprintf(w, "[%d,%d] %+v\n", x, y, cell)
	}
	return w.String()
}
