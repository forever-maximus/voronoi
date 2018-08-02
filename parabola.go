package main

import "math"

// Return the x coordinate of the intersection between two parabolas given their directrix and foci.
// Since the beachline is x-monotone the breakpoints can be differentiated by which parabola is left
// of the breakpoint and which is right (i.e. breakpoint (a,b) is not the same as (b,a) breakpoint).
func getBreakpointXCoordinate(focusPair *breakpoint, directrix float64) float64 {
	// Get the coefficients for each parabola
	a1, b1, c1 := getCoefficients(focusPair.leftSite, directrix)
	a2, b2, c2 := getCoefficients(focusPair.rightSite, directrix)

	aDiff := a1 - a2
	bDiff := b1 - b2
	cDiff := c1 - c2

	discriminant := (bDiff * bDiff) - (4 * aDiff * cDiff)

	// Quadratic formula
	intersection1 := ((-1 * bDiff) + math.Sqrt(discriminant)) / (2 * aDiff)
	intersection2 := ((-1 * bDiff) - math.Sqrt(discriminant)) / (2 * aDiff)

	// Problem: Given a focus pair (a,b), we get two intersections returned (x1 and x2), how do we
	//          know which one is the intersection for (a,b) (i.e. is x1 (a,b) or (b,a))
	// Idea:
	// 1. Determine which parabola sits on top of the other one by checking which has larger focal y coordinate.
	//    The one with larger focal y will be the base parabola with the other building off it
	// 2. If the base parabola is the left focus in the pair (eg. pair = (base, child)) then return the smallest
	//    x intercept, if the base is the right focus (eg. pair = (child, base)) then return the largest x intercept.
	if focusPair.leftSite.y > focusPair.rightSite.y {
		return math.Min(intersection1, intersection2)
	}
	return math.Max(intersection1, intersection2)
}

// Return the coefficients in the parabola standard form equation ax^2 + bx + c = 0
// Note - I got the equations to determine the coefficients from math stackexchange post
func getCoefficients(site *site, directrix float64) (float64, float64, float64) {
	// Double the distance from the focus to the directrix
	dp := 2.0 * (site.y - directrix)

	a := 1.0 / dp
	b := (-2.0 * site.x) / dp
	c := directrix + (dp / 4.0) + ((site.x * site.x) / dp)

	return a, b, c
}
