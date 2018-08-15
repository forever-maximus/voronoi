package main

import (
	"container/heap"

	"github.com/fogleman/gg"
)

type site struct {
	x, y float64
}

type breakpoint struct {
	leftSite  *site
	rightSite *site
}

type boundingBox struct {
	height, width float64
}

// Event - represents an event used in Fortune's algorithm for generating voronoi diagram
type Event struct {
	eventType string
	location  site
	leafNode  *node
}

func main() {
	// Some input sites
	// siteList := []site{
	// 	site{x: 40, y: 120},
	// 	site{x: 70, y: 150},
	// 	site{x: 120, y: 70},
	// 	site{x: 260, y: 170},
	// }

	siteList := []site{
		site{x: 40, y: 120},
		site{x: 70, y: 150},
		site{x: 120, y: 70},
		site{x: 260, y: 170},
		site{x: 176, y: 220},
		site{x: 246, y: 110},
		site{x: 430, y: 450},
		site{x: 200, y: 400},
		site{x: 400, y: 100},
	}

	// Create a priority queue, put the items in it
	pq := make(PriorityQueue, len(siteList))
	for i, coordinates := range siteList {
		pq[i] = &Item{
			value:    Event{eventType: "site", location: coordinates},
			priority: coordinates.y,
			index:    i,
		}
	}
	heap.Init(&pq)

	fortunesAlgorithm(&pq, siteList)
}

func fortunesAlgorithm(eventQueue *PriorityQueue, siteList []site) {
	beachline := redblacktree{root: nil}
	dcel := doublyConnectedEdgeList{vertices: nil, edges: nil}
	counter := 1
	for eventQueue.Len() > 0 {
		item := heap.Pop(eventQueue).(*Item)
		if item.value.eventType == "site" {
			// Site event
			beachline.insert(counter, &item.value.location, eventQueue, &dcel)
		} else {
			// Circle event
			beachline.removeArc(item.value.leafNode, eventQueue, &item.value.location, &dcel,
				item.priority, counter)
		}
		counter++
	}

	//beachline.inorderTraversal()

	// Add bounding box and connect half infinite edges to it
	boundingBox := boundingBox{height: 500, width: 500}
	connectEdgesToBoundary(beachline.root, boundingBox, &dcel)

	// Draw voronoi
	drawVoronoi(boundingBox, &dcel, siteList)
}

func drawVoronoi(boundingBox boundingBox, dcel *doublyConnectedEdgeList, siteList []site) {
	// The drawing module sets the top left as (0, 0) and bottom right as (boundary.width, boundary.height).
	// i.e. flipped the y axis direction from what was expected - ignore for now
	voronoi := gg.NewContext(int(boundingBox.width), int(boundingBox.height))
	voronoi.SetRGB(1, 1, 1)
	voronoi.Clear()
	voronoi.SetLineWidth(3)

	for _, halfEdge := range dcel.edges {
		voronoi.SetRGB(0.3, 0.7, 0.8)
		voronoi.DrawLine(halfEdge.originVertex.x, boundingBox.height-halfEdge.originVertex.y,
			halfEdge.twinEdge.originVertex.x, boundingBox.height-halfEdge.twinEdge.originVertex.y)
		voronoi.Stroke()
	}

	for _, site := range siteList {
		voronoi.SetRGB(0.9, 0.5, 0.6)
		voronoi.DrawPoint(site.x, boundingBox.height-site.y, 2.0)
		voronoi.Stroke()
	}
	voronoi.SavePNG("testingVoronoi.png")
}
