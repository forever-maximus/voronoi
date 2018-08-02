package main

import (
	"container/heap"
	"fmt"
)

type site struct {
	x, y float64
}

type breakpoint struct {
	leftSite  *site
	rightSite *site
}

// Event - represents an event used in Fortune's algorithm for generating voronoi diagram
type Event struct {
	eventType string
	location  site
	leafNode  *node
}

func main() {
	// Some input sites
	siteList := []site{
		site{x: 40, y: 120},
		site{x: 70, y: 150},
		site{x: 120, y: 70},
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

	fortunesAlgorithm(&pq)
}

func fortunesAlgorithm(eventQueue *PriorityQueue) {
	beachline := redblacktree{root: nil}
	counter := 1
	for eventQueue.Len() > 0 {
		item := heap.Pop(eventQueue).(*Item)
		if item.value.eventType == "site" {
			// Site event
			beachline.insert(counter, &item.value.location, eventQueue)
			counter++
		} else {
			// Circle event
			fmt.Println("Circle event detected!")
		}
	}
	beachline.inorderTraversal()
}
