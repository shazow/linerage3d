package collision

// IsBoundingBox returns true if a box intercepts b box.
func IsBoxCollision(a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32) bool {
	return a1_x <= b2_x && a2_x >= b1_x && a1_y <= b2_y && a2_y >= b1_y
}

// IsCollision returns true if segment a1->a2 intersects segment b1->b2.
// Collisions are checked [a,b). That is, a->b->c will not collide, but
// a->b,a->c will collide.
func IsCollision2D(a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32) bool {
	// Partly based on https://stackoverflow.com/questions/563198/
	s1_x := a2_x - a1_x
	s1_y := a2_y - a1_y
	s2_x := b2_x - b1_x
	s2_y := b2_y - b1_y

	connected := a2_x == b1_x && a2_y == b1_y
	denom := s1_x*s2_y - s2_x*s1_y
	if denom == 0 {
		// Collinear

		if connected {
			// Pointing away?
			return (s1_x*s2_x < 0) || (s1_y*s2_y < 0)
		}

		// Any of the wrong points connected? (Head-on connected)
		if a2_x == b2_x && a2_y == b2_y {
			return true
		}

		// Basically box collision
		return a1_x <= b2_x && a2_x >= b1_x && a1_y <= b2_y && a2_y >= b1_y
	}

	if connected {
		// Connected but not collinear
		return false
	}

	denomPositive := denom > 0

	s3_x := a1_x - b1_x
	s3_y := a1_y - b1_y

	s_numer := s1_x*s3_y - s1_y*s3_x
	if (s_numer <= 0) == denomPositive {
		return false
	}

	t_numer := s2_x*s3_y - s2_y*s3_x
	if (t_numer <= 0) == denomPositive {
		return false
	}

	if ((s_numer >= denom) == denomPositive) || ((t_numer >= denom) == denomPositive) {
		return false
	}

	/*
		// Intersecting point
		t := t_numer / denom
		i_x := a1_x + (t * s1_x)
		i_y := a1_x + (t * s1_y)
		fmt.Printf("collided at %v,%v when comparing %v\n", i_x, i_y, []float32{a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y})

		fmt.Println(s_numer, t_numer, denom, "numer", a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y)
	*/
	return true
}
