package main

import (
	"math"
)

func connectEdgesToBoundary(currentNode *node, boundingBox boundingBox, dcel *doublyConnectedEdgeList) {
	// Each internal node left in the beachline represents a half infinite edge
	// Need to connect each of these to the bounding box
	if currentNode.breakpoint != nil {
		// Steps:
		// 1. Find the midpoint between the left and right site in breakpoint
		// 2. Use this point and the nodes half edge twin vertex point to determine the equation of the line
		// 3. Calculate where the line intercepts the bounding box (there will be two points)
		// 4. Determine which of the points is the closest and set the nodes half edge vertex as this point

		vertex := currentNode.halfEdge.twinEdge.originVertex

		// Initialise boundary vertices
		vertexBoundingA := getVertex()
		vertexBoundingA.x = -1.0
		vertexBoundingB := getVertex()

		// Consider ignoring if vertex lies outside bounding box - maybe add dummy vertex
		if vertex.x < 0 || vertex.x > boundingBox.width || vertex.y < 0 || vertex.y > boundingBox.height {
			fakeVertex := dcel.addIsolatedVertex(vertex.x, vertex.y)
			currentNode.halfEdge.originVertex = fakeVertex
		} else {

			xMidpoint := (currentNode.breakpoint.leftSite.x + currentNode.breakpoint.rightSite.x) / 2
			yMidpoint := (currentNode.breakpoint.leftSite.y + currentNode.breakpoint.rightSite.y) / 2

			gradient := (yMidpoint - vertex.y) / (xMidpoint - vertex.x)
			b := yMidpoint - (gradient * xMidpoint)
			bottomBoundInterceptX := (-1.0 * b) / gradient
			rightBoundInterceptY := (gradient * boundingBox.width) + b
			topBoundInterceptX := (boundingBox.height - b) / gradient

			// TODO - handle case of parallel line

			if b >= 0 && b <= boundingBox.height {
				vertexBoundingA.x = 0
				vertexBoundingA.y = b
			}
			if bottomBoundInterceptX >= 0 && bottomBoundInterceptX <= boundingBox.width {
				vertexBoundingB.x = bottomBoundInterceptX
				vertexBoundingB.y = 0
			}
			if rightBoundInterceptY >= 0 && rightBoundInterceptY <= boundingBox.height {
				if vertexBoundingA.x < 0 {
					vertexBoundingA.x = boundingBox.width
					vertexBoundingA.y = rightBoundInterceptY
				} else {
					vertexBoundingB.x = boundingBox.width
					vertexBoundingB.y = rightBoundInterceptY
				}
			}
			if topBoundInterceptX >= 0 && topBoundInterceptX <= boundingBox.width {
				if vertexBoundingA.x < 0 {
					vertexBoundingA.x = topBoundInterceptX
					vertexBoundingA.y = boundingBox.height
				} else {
					vertexBoundingB.x = topBoundInterceptX
					vertexBoundingB.y = boundingBox.height
				}
			}

			// Distance to A from vertex and distance to A from midpoint
			distanceVertexToA := math.Sqrt(math.Pow((vertex.x-vertexBoundingA.x), 2) + math.Pow((vertex.y-vertexBoundingA.y), 2))
			distanceMidpointToA := math.Sqrt(math.Pow((xMidpoint-vertexBoundingA.x), 2) +
				math.Pow((yMidpoint-vertexBoundingA.y), 2))

			// Add the boundary vertex which is closer to midpoint than node vertex to the dcel
			// and connect the halfedge on the node to it
			newVertex := getVertexPointer()
			if distanceMidpointToA < distanceVertexToA {
				newVertex = dcel.addIsolatedVertex(vertexBoundingA.x, vertexBoundingA.y)
			} else {
				newVertex = dcel.addIsolatedVertex(vertexBoundingB.x, vertexBoundingB.y)
			}
			currentNode.halfEdge.originVertex = newVertex
		}

		// Run on child nodes
		if currentNode.left != nil {
			connectEdgesToBoundary(currentNode.left, boundingBox, dcel)
		}
		if currentNode.right != nil {
			connectEdgesToBoundary(currentNode.right, boundingBox, dcel)
		}
	}
}
