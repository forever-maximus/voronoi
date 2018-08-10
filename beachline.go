// Note - Currently this is just a standard Binary Search Tree Implementation
// The plan is to eventually move to a redblack tree
package main

import (
	"container/heap"
	"fmt"
)

const (
	red   bool = true
	black bool = false
)

// Internal Nodes - These will have a breakpoint but arcSite will be nil
// Leaf Nodes - These have an arcSite but no breakpoint
type node struct {
	left, right, parent, next, previous *node
	colour                              bool
	breakpoint                          *breakpoint
	arcSite                             *site
	key                                 int
	circleEvent                         *Item
	halfEdge                            *halfEdge
}

type redblacktree struct {
	root *node
}

func (rbtree *redblacktree) insert(newKey int, newSite *site, eventQueue *PriorityQueue, dcel *doublyConnectedEdgeList) {
	if rbtree.root == nil {
		rbtree.root = &node{key: newKey, colour: black, arcSite: newSite}
	} else {
		rbtree.root = rbtree.root.insert(rbtree.root, newKey, newSite, eventQueue, dcel)
	}
}

// Insert finds the arc on the beachline above the new site (this is the leaf node found) and replaces it with a subtree
// consisting of 2 internal nodes (breakpoints between arcs) and 3 leaf nodes (arcs on beachline).
//                                                x
// (leaf node found)                            /   \       | x = internal node (breakpoint)
//         o  ---------------------->          x     o      | o = leaf node (arc)
//                   (transform)             /  \
//                                          o    o
func (n *node) insert(currentNode *node, newKey int, newSite *site, eventQueue *PriorityQueue,
	dcel *doublyConnectedEdgeList) *node {
	// Check if this is a leaf node
	if currentNode.breakpoint == nil {

		if currentNode.circleEvent != nil {
			// Remove circle event from event queue as it is a false alarm
			heap.Remove(eventQueue, currentNode.circleEvent.index)
		}

		// Define the breakpoints that will be used in the two new internal nodes
		leftBreakpoint := breakpoint{
			leftSite:  currentNode.arcSite,
			rightSite: newSite,
		}
		rightBreakpoint := breakpoint{
			leftSite:  newSite,
			rightSite: currentNode.arcSite,
		}

		// The 3 leaf nodes that represent the arcs
		leftLeafNode := node{
			arcSite:  currentNode.arcSite,
			previous: currentNode.previous,
			key:      currentNode.key,
		}
		middleLeafNode := node{arcSite: newSite,
			previous: &leftLeafNode,
			key:      newKey,
		}
		rightLeafNode := node{
			arcSite:  currentNode.arcSite,
			next:     currentNode.next,
			previous: &middleLeafNode,
			key:      currentNode.key,
		}
		middleLeafNode.next = &rightLeafNode
		leftLeafNode.next = &middleLeafNode

		// Create and add half-edges to dcel structure
		leftHalfEdge := dcel.addIsolatedEdge()
		rightHalfEdge := leftHalfEdge.twinEdge

		// The 2 internal nodes which represent each edge being traced out
		leftInternalNode := node{
			left:       &leftLeafNode,
			right:      &middleLeafNode,
			breakpoint: &leftBreakpoint,
			halfEdge:   leftHalfEdge,
		}
		rightInternalNode := node{
			left:       &leftInternalNode,
			right:      &rightLeafNode,
			breakpoint: &rightBreakpoint,
			halfEdge:   rightHalfEdge,
		}

		// Set parent nodes
		leftInternalNode.parent = &rightInternalNode
		rightLeafNode.parent = &rightInternalNode
		middleLeafNode.parent = &leftInternalNode
		leftLeafNode.parent = &leftInternalNode

		// Check for circle event (i.e. check for unique triples of sites on beachline (a,b,c))
		leftLeafNode.circleEvent = checkCircleEvent(&leftLeafNode, newSite.y, eventQueue)
		rightLeafNode.circleEvent = checkCircleEvent(&rightLeafNode, newSite.y, eventQueue)

		return &rightInternalNode
	}

	// The directrix will be at the same y coordinate as the new site being added
	breakpointXCoordinate := getBreakpointXCoordinate(currentNode.breakpoint, newSite.y)

	if newSite.x < breakpointXCoordinate {
		currentNode.left = currentNode.insert(currentNode.left, newKey, newSite, eventQueue, dcel)
		currentNode.left.parent = currentNode
	} else if newSite.x > breakpointXCoordinate {
		currentNode.right = currentNode.insert(currentNode.right, newKey, newSite, eventQueue, dcel)
		currentNode.right.parent = currentNode
	}

	return currentNode
}

func (rbtree *redblacktree) removeArc(leafNode *node, eventQueue *PriorityQueue, circleCenter *site,
	dcel *doublyConnectedEdgeList, sweepline float64) {
	if leafNode.parent == nil {
		// This happens if only 1 node is in the tree and you call remove - this should never happen
		rbtree.root = nil
	}

	leftLeafNode := leafNode.previous
	rightLeafNode := leafNode.next

	// Check if left or right have circle events - these won't be valid once leaf node has been removed
	if leftLeafNode.circleEvent != nil {
		heap.Remove(eventQueue, leftLeafNode.circleEvent.index)
	}
	if rightLeafNode.circleEvent != nil {
		heap.Remove(eventQueue, rightLeafNode.circleEvent.index)
	}

	// Get the ancestors of the leafnode required
	siblingNode := getSibling(leafNode)
	parentNode := leafNode.parent
	grandparent := parentNode.parent
	greatGrandParent := grandparent.parent

	// Check whether leafnode is in right or left subtree of the grandparent - then replace leaf+parent with sibling
	if grandparent.left == parentNode {
		grandparent.left = siblingNode
	} else {
		grandparent.right = siblingNode
	}

	// Check if leaf node is left or right child
	alteredInternalNode := &node{}
	if parentNode.left == leafNode {
		// Leaf node is a left child
		// 1. Get inorder predecessor of the leaf node (this will always be an internal node)
		alteredInternalNode = inorderPredecessorOfLeafNode(leafNode)

		// 2. Get minimum leaf node in sibling subtree (if sibling is a leaf it will be sibling node)
		minimumInSiblingTree := getMinimumLeafNode(siblingNode)

		// 3. Replace the right site in the inorder predecessor's breakpoint (returned in step 1) with
		//    the minimum leaf in sibling subtree (step 2)
		alteredInternalNode.right = minimumInSiblingTree
	} else {
		// Leaf node is right child
		// 1. Get inorder successor of the leaf node (this will always be an internal node)
		alteredInternalNode = inorderSuccessorOfLeafNode(leafNode)

		// 2. Get maximum leaf node in sibling subtree (if sibling is a leaf it will be sibling node)
		maximumInSiblingTree := getMaximumLeafNode(siblingNode)

		// 3. Replace the left site in the inorder successor's breakpoint (returned in step 1) with
		//    the maximum leaf in sibling subtree (step 2)
		alteredInternalNode.left = maximumInSiblingTree
	}

	// Work out the replacement site that will go in the great grandparent node breakpoint
	arcSiteReplacement := &site{}
	if parentNode.breakpoint.leftSite != leafNode.arcSite {
		arcSiteReplacement = parentNode.breakpoint.leftSite
	} else {
		arcSiteReplacement = parentNode.breakpoint.rightSite
	}

	// Replace the removed arc from the breakpoint in the great grandparent
	isLeafParentLeft := true
	if greatGrandParent.breakpoint.leftSite == leafNode.arcSite {
		greatGrandParent.breakpoint.leftSite = arcSiteReplacement
	} else {
		greatGrandParent.breakpoint.rightSite = arcSiteReplacement
		isLeafParentLeft = false
	}

	// Assign which halfedge is inbound from left and which is from right
	leftHalfEdge, rightHalfEdge := &halfEdge{}, &halfEdge{}
	if isLeafParentLeft == true {
		leftHalfEdge = parentNode.halfEdge
		rightHalfEdge = greatGrandParent.halfEdge
	} else {
		leftHalfEdge = greatGrandParent.halfEdge
		rightHalfEdge = parentNode.halfEdge
	}

	// Create new halfedge pair
	newHalfEdge := dcel.addIsolatedEdge()

	// Add circle center as a new vertex of the voronoi diagram
	voronoiVertex := dcel.addIsolatedVertex(circleCenter.x, circleCenter.y)

	// Connect halfedges to the vertex
	leftHalfEdge.originVertex = voronoiVertex
	rightHalfEdge.originVertex = voronoiVertex
	newHalfEdge.originVertex = voronoiVertex

	// Connect halfedges with each other
	leftHalfEdge.twinEdge.nextEdge = newHalfEdge
	newHalfEdge.twinEdge.nextEdge = rightHalfEdge
	rightHalfEdge.twinEdge.nextEdge = leftHalfEdge

	// Fix the next and previous nodes of the nodes left and right of the node which was just removed
	leftLeafNode.next = rightLeafNode
	rightLeafNode.previous = leftLeafNode

	// Update the halfedge record in the great grandparent since the breakpoint has changed
	greatGrandParent.halfEdge = newHalfEdge.twinEdge

	// Check for new circle events now that the leaf has been removed from the beachline
	leftLeafNode.circleEvent = checkCircleEvent(leftLeafNode, sweepline, eventQueue)
	rightLeafNode.circleEvent = checkCircleEvent(rightLeafNode, sweepline, eventQueue)
}

// Assumption - due to the nature of the algorithm a node will always have a sibling unless root node
func getSibling(node *node) *node {
	if node.parent == nil {
		return nil // this will happen if the input node is the root node
	}
	if node.parent.right == node {
		return node.parent.left
	}
	return node.parent.right
}

// Return the inorder successor of a leaf node
func inorderSuccessorOfLeafNode(currentNode *node) *node {
	for currentNode.parent != nil {
		if currentNode.parent.left == currentNode {
			return currentNode.parent
		}
		currentNode = currentNode.parent
	}
	// The input leaf node was either on end of beachline (right-most leaf) or root node (tree size = 1)
	return nil
}

// Return the inorder predecessor of a leaf node
func inorderPredecessorOfLeafNode(currentNode *node) *node {
	for currentNode.parent != nil {
		if currentNode.parent.right == currentNode {
			return currentNode.parent
		}
		currentNode = currentNode.parent
	}
	// The input leaf node was either on start of beachline (left-most leaf) or root node (tree size = 1)
	return nil
}

// Return the minimum leaf node in a subtree
func getMinimumLeafNode(currentNode *node) *node {
	if currentNode.breakpoint != nil {
		if currentNode.left != nil {
			return getMinimumLeafNode(currentNode.left)
		}
		return getMinimumLeafNode(currentNode.right)
	}
	return currentNode
}

// Return the maximum leaf node in a subtree
func getMaximumLeafNode(currentNode *node) *node {
	if currentNode.breakpoint != nil {
		if currentNode.right != nil {
			return getMaximumLeafNode(currentNode.right)
		}
		return getMaximumLeafNode(currentNode.left)
	}
	return currentNode
}

func (rbtree *redblacktree) inorderTraversal() {
	rbtree.root.inorderTraversal(rbtree.root)
}

func (n *node) inorderTraversal(currentNode *node) {
	if currentNode != nil {
		currentNode.inorderTraversal(currentNode.left)
		//fmt.Printf("%d ", currentNode.key)
		fmt.Println(currentNode)
		currentNode.inorderTraversal(currentNode.right)
	}
}
