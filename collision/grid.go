package collision

import (
	"bytes"
	"fmt"
	"image"

	mgl "github.com/go-gl/mathgl/mgl32"
)

func GridCollider(bounds image.Rectangle) Collider {
	size := bounds.Size()

	grid := &gridCollider{
		bounds: Boundary{float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Max.X), float32(bounds.Max.Y)},
		width:  size.X,
		height: size.Y,
	}
	grid.Reset()
	return grid
}

type gridCollider struct {
	bounds Boundary
	width  int
	height int
	grid   []gridCell
}

func (grid *gridCollider) String() string {
	w := &bytes.Buffer{}
	for idx, cell := range grid.grid {
		if len(cell) == 0 {
			continue
		}

		x := idx%grid.width + int(grid.bounds.X1)
		y := (idx-x)/grid.height + int(grid.bounds.Y1)

		fmt.Fprintf(w, "[%d,%d] %+v\n", x, y, cell)
	}
	return w.String()
}

func (grid *gridCollider) Reset() {
	grid.grid = make([]gridCell, grid.width*grid.height, grid.width*grid.height)
}

func (grid *gridCollider) Track(segment *[]mgl.Vec3) Tracker {
	return &gridTracker{
		grid:    grid,
		segment: segment,
	}
}

func (grid *gridCollider) index(x, y float32) int {
	return int(x-grid.bounds.X1) + int(y-grid.bounds.Y1)*grid.width
}

type gridTracker struct {
	grid    *gridCollider
	cs      *cellSegment
	segment *[]mgl.Vec3
}

func (tracker *gridTracker) Update() error {
	cs := tracker.cs
	segment := *tracker.segment
	grid := tracker.grid

	if len(segment) < 2 {
		return nil
	}

	offset := len(segment) - 2

	var vec mgl.Vec3 = segment[offset+1]
	x1, y1 := vec[0], vec[2]

	// Check boundary
	if grid.bounds.X1 >= x1 || x1 >= grid.bounds.X2 ||
		grid.bounds.Y1 >= y1 || y1 >= grid.bounds.Y2 {
		return CollisionBoundary
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

	tracker.cs = cs
	return err
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
