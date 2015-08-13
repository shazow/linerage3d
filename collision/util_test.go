package collision

import "testing"

func TestBoxCollision(t *testing.T) {
	tests := []struct {
		result bool

		a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32
	}{
		{true, 0, 0, 1, 1, 0, 0, 1, 1},
		{true, 0, 0, 1, 1, 1, 1, 0, 0},
		{true, 0, 0, 1, 1, 1, 0, 0, 1},
		{false, 0, 0, 1, 1, 2, 2, 3, 3},
	}

	for i, test := range tests {
		r := IsBoxCollision(test.a1_x, test.a1_y, test.a2_x, test.a2_y, test.b1_x, test.b1_y, test.b2_x, test.b2_y)
		if r != test.result {
			t.Errorf("IsBoxCollision test #%d failed: %v", i, test)
		}
	}
}

func TestIsCollision(t *testing.T) {
	tests := []struct {
		result bool

		a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32
	}{
		// Test collision is [a, b)
		{true, 0, 0, 1, 1, 1, 1, 0, 0},      // a -> -a
		{true, 0, 0, 1, 1, 0, 0, 1, 1},      // a -> a
		{true, 0, 0, 1, 1, 1, 1, 0.5, 0.5},  // a -> b intersect
		{false, 0, 0, 1, 1, 1, 1, 2, 2},     // a -> b
		{true, 2, 2, 1, 1, 1, 1, 1.5, 1.5},  // b -> a intersect
		{false, 2, 2, 1, 1, 1, 1, 0.5, 0.5}, // b -> a
		{false, 1, 0, 2, 0, 2, 0, 3, 0},     // a -> b
		{false, 0, 0, 1, 1, 1, 1, 1.5, 1.5}, // a -> b, non-collinear
		{true, 0, 0, 4, 4, 1, 1, 2, 2},      // a -> d, b -> c, collinear contained
		{true, 0, 0, 4, 4, 2, 2, 1, 1},      // a -> d, c -> b, collinear contained
		{false, 0, 1, 0, 4, 1, 3, 1, 2},     // a -> d, b -> c, parallel offset
		{false, 0, 1, 0, 4, 1, 2, 1, 3},     // a -> d, c -> b, parallel offset
		{true, 1, 0, 2, 0, 3, 0, 2, 0},      // a -> b <- c

		{true, 0, 1, 0, 4, 0, 1, 1, 3}, // a -> b, a->c, non-collinear
		{true, 0, 1, 0, 4, 0, 1, 0, 3},
		{false, 0, 1, 0, 4, 0, 4, 1, 3},
		{true, 0, 1, 0, 4, 0, 4, 0, 3},

		{true, 3, 1, 3, 2, 2.5, 1.5, 3.5, 1.5},
		{true, 3, 1, 3, 2, 2.5, 1.5, 3.5, 1},
		{true, 3, 1, 3, 2, 2.5, 1, 4, 1},
		{true, 1, 1, 1, 3, 0, 1, 4, 1},
		{true, 1, 0, 1, 1, 0, 1, 4, 1},
		{false, 0, 0, 1, 1, 2, 2, 3, 3},
		{false, 2, 0, 3, 0, 3, 0, 3, 1},
		{false, 2, 1, 3, 1, 1, 0, 2, 0}, // collinear disjoint vertically
		//{false, 1.37, 1.39, 1.34, 1.35, 1.31, 1.31, 1.34, 1.35},
	}

	for i, test := range tests {
		r := IsCollision2D(test.a1_x, test.a1_y, test.a2_x, test.a2_y, test.b1_x, test.b1_y, test.b2_x, test.b2_y)
		if r != test.result {
			t.Errorf("IsCollision2D test #%d failed: %v", i, test)
		}
	}
}
