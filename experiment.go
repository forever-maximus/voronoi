package main

import (
	"container/heap"
	"fmt"
)

type site struct {
	x, y int
}

type breakpoint struct {
	leftSite  *site
	rightSite *site
}

// Event - represents an event used in Fortune's algorithm for generating voronoi diagram
type Event struct {
	eventType string
	location  site
}

func main() {
	// Some input sites
	siteList := []site{
		site{x: 60, y: 100},
		site{x: 30, y: 40},
		site{x: 140, y: 160},
		site{x: 300, y: 220},
		site{x: 210, y: 290},
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, len(siteList))
	for i, coordinates := range siteList {
		pq[i] = &Item{
			value:    Event{eventType: "site", location: coordinates},
			priority: coordinates.y,
			index:    i,
		}
	}
	heap.Init(&pq)

	// Testing BST implementation
	beachline := redblacktree{root: nil}
	testSite := site{x: 40, y: 120}
	beachline.insert(10, &testSite)
	beachline.inorderTraversal()
	fmt.Println()

	secondTestSite := site{x: 70, y: 100}
	beachline.insert(12, &secondTestSite)
	beachline.inorderTraversal()
	fmt.Println()

	thirdTestSite := site{x: 120, y: 70}
	beachline.insert(15, &thirdTestSite)
	beachline.inorderTraversal()

	//fortunesAlgorithm(&pq)
}

func fortunesAlgorithm(eventQueue *PriorityQueue) {
	for eventQueue.Len() > 0 {
		item := heap.Pop(eventQueue).(*Item)
		if item.value.eventType == "site" {
			// Site event

		} else {
			// Circle event
		}
		fmt.Printf("%.2d:%+v ", item.priority, item.value)
	}
}
