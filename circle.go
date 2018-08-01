package main

import (
	"math"
)

// Check whether a leaf node has a circle event and add to event queue if true
func checkCircleEvent(leafNode *node, sweepline int, eventQueue *PriorityQueue) *Item {
	if leafNode.previous == nil || leafNode.next == nil {
		return nil
	}

	leftSite := leafNode.previous.arcSite
	rightSite := leafNode.next.arcSite
	middleSite := leafNode.arcSite

	if leftSite == rightSite {
		return nil
	}

	// Let the center be (a, b) -> each point is equal distance to (a, b) since all lie on the circumference.
	// Given points (x1, y1), (x2, y2), (x3, y3) equate square of distance:
	// (a - x1)^2 + (b - y1)^2 = (a - x2)^2 + (b - y2)^2 = (a - x3)^2 + (b - y3)^2
	// Note - we need to find the constants in the expanded form below for each point.
	// (a - x)^2 + (b - y)^2 --> a^2 + b^2 - 2 a x + x^2 - 2 b y + y^2
	// We can use this to generate two linear equations and solve for (a, b)

	// Difference between first and second, and second and third - These are the constants for linear equations
	// x^2
	leftXSquaredDiff := (float64(leftSite.x) * float64(leftSite.x)) - (float64(middleSite.x) * float64(middleSite.x))
	rightXSquaredDiff := (float64(middleSite.x) * float64(middleSite.x)) - (float64(rightSite.x) * float64(rightSite.x))

	// y^2
	leftYSquaredDiff := (float64(leftSite.y) * float64(leftSite.y)) - (float64(middleSite.y) * float64(middleSite.y))
	rightYSquareDiff := (float64(middleSite.y) * float64(middleSite.y)) - (float64(rightSite.y) * float64(rightSite.y))

	// -2x -> (this is the a term constant)
	leftXLinearDiff := (2.0 * float64(leftSite.x)) - (2.0 * float64(middleSite.x))
	rightXLinearDiff := (2.0 * float64(middleSite.x)) - (2.0 * float64(rightSite.x))

	// -2y -> (this is the b term constant)
	leftYLinearDiff := (2.0 * float64(leftSite.y)) - (2.0 * float64(middleSite.y))
	rightYLinearDiff := (2.0 * float64(middleSite.y)) - (2.0 * float64(rightSite.y))

	// x^2 + y^2 --> (let this be k)
	constantsLeft := leftXSquaredDiff + leftYSquaredDiff
	constantsRight := rightXSquaredDiff + rightYSquareDiff

	// We now have the constants of two linear equations of the form k - 2ax - 2by = 0

	// (I think) equation for b - from substituting linear equation in other linear equation
	//      2x1.k2 - 2x2.k1
	// b = -------------------   (note: the '.' is multiplication)
	//      2y1.2x2 - 2y2.2x1
	b := ((-1.0 * leftXLinearDiff * constantsRight) - (-1.0 * rightXLinearDiff * constantsLeft)) / ((leftYLinearDiff * rightXLinearDiff) - (rightYLinearDiff * leftXLinearDiff))

	// now substitute b back into equation a = (k1 - 2y1.b) / 2x1
	a := (constantsLeft - (leftYLinearDiff * b)) / leftXLinearDiff

	// This gives us the circle center (a, b)
	// fmt.Println("(a, b) value is --> (", a, ", ", b, ")")

	// Calculate the radius as distance from center to any of the sites (we choose left site)
	radius := math.Sqrt(math.Pow((float64(leftSite.x)-a), 2) + math.Pow((float64(leftSite.y)-b), 2))

	// Check that the bottom of circle lies below the sweepline
	bottomOfCircle := b - radius
	if bottomOfCircle > float64(sweepline) {
		return nil
	}

	// The circle event is valid - create and add to event queue
	circleCenter := site{x: int(a), y: int(b)} // TODO - this feels like bad design - technically not a site
	circleEvent := &Item{
		value:    Event{eventType: "circle", location: circleCenter, leafNode: leafNode},
		priority: int(bottomOfCircle),
		index:    eventQueue.Len(),
	}
	eventQueue.Push(circleEvent)

	return circleEvent
}
