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
	edge                                *halfEdge
}

type redblacktree struct {
	root *node
}

func (rbtree *redblacktree) insert(newKey int, newSite *site, eventQueue *PriorityQueue) {
	if rbtree.root == nil {
		rbtree.root = &node{key: newKey, colour: black, arcSite: newSite}
	} else {
		rbtree.root = rbtree.root.insert(rbtree.root, newKey, newSite, eventQueue)
	}
}

func (n *node) insert(currentNode *node, newKey int, newSite *site, eventQueue *PriorityQueue) *node {
	// Check if this is a leaf node
	if currentNode.breakpoint == nil {

		if currentNode.circleEvent != nil {
			// Remove circle event from event queue as it is a false alarm
			heap.Remove(eventQueue, currentNode.circleEvent.index)
		}

		// Define the breakpoints that will be used in the two new internal nodes
		leftBreakpoint := breakpoint{leftSite: currentNode.arcSite, rightSite: newSite}
		rightBreakpoint := breakpoint{leftSite: newSite, rightSite: currentNode.arcSite}

		// The 3 leaf nodes that represent the arcs
		leftLeafNode := node{arcSite: currentNode.arcSite, previous: currentNode.previous, key: currentNode.key}
		middleLeafNode := node{arcSite: newSite, previous: &leftLeafNode, key: newKey}
		rightLeafNode := node{arcSite: currentNode.arcSite, next: currentNode.next, previous: &middleLeafNode, key: currentNode.key}
		middleLeafNode.next = &rightLeafNode
		leftLeafNode.next = &middleLeafNode

		leftInternalNode := node{left: &leftLeafNode, right: &middleLeafNode, breakpoint: &leftBreakpoint}
		rightInternalNode := node{left: &leftInternalNode, right: &rightLeafNode, breakpoint: &rightBreakpoint}

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

	if float64(newSite.x) < breakpointXCoordinate {
		currentNode.left = currentNode.insert(currentNode.left, newKey, newSite, eventQueue)
		currentNode.left.parent = currentNode
	} else if float64(newSite.x) > breakpointXCoordinate {
		currentNode.right = currentNode.insert(currentNode.right, newKey, newSite, eventQueue)
		currentNode.right.parent = currentNode
	}

	return currentNode
}

func (rbtree *redblacktree) delete(keyToRemove int) *node {
	if rbtree.root == nil {
		return nil
	}
	return rbtree.root.delete(rbtree.root, keyToRemove)
}

func (n *node) delete(currentNode *node, keyToRemove int) *node {
	// Base case
	if currentNode == nil {
		return currentNode
	}

	if keyToRemove < currentNode.key {
		currentNode.left = currentNode.delete(currentNode.left, keyToRemove)

	} else if keyToRemove > currentNode.key {
		currentNode.right = currentNode.delete(currentNode.right, keyToRemove)

	} else {

		// If node only has one child or no children
		if currentNode.left == nil {
			return currentNode.right
		} else if currentNode.right == nil {
			return currentNode.left
		}

		// handle case where the node to delete has two children
		// Get smallest node in right subtree (ie. inorder successor to current node)
		inorderSuccessor := currentNode.minValueNode(currentNode.right)
		// Copy it's value to the current node
		currentNode.key = inorderSuccessor.key
		// Delete the inorder successor
		currentNode.right = currentNode.delete(currentNode.right, inorderSuccessor.key)
	}

	return currentNode
}

// Find the node with minimum value given a binary search tree
func (n *node) minValueNode(currentNode *node) *node {
	if currentNode.left != nil {
		return currentNode.minValueNode(currentNode.left)
	}
	return currentNode
}

// // Return the minimum value leaf node in the subtree
// func minValueLeafNode(currentNode *node) *node {
// 	if currentNode.breakpoint != nil {
// 		if currentNode.left != nil {
// 			return minValueLeafNode(currentNode.left)
// 		} else if currentNode.right != nil {
// 			return minValueLeafNode(currentNode.right)
// 		}
// 		// nothing should get here - if it does it is an error - all internal nodes should have at least one child
// 	}
// 	return currentNode
// }

// // Return the maximum value leaf node in the subtree
// func maxValueLeafNode(currentNode *node) *node {
// 	if currentNode.breakpoint != nil {
// 		if currentNode.right != nil {
// 			return maxValueLeafNode(currentNode.right)
// 		} else if currentNode.left != nil {
// 			return maxValueLeafNode(currentNode.left)
// 		}
// 		// error - all internal nodes should have at least one child
// 	}
// 	return currentNode
// }

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
