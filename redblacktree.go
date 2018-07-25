// Note - Currently this is just a standard Binary Search Tree Implementation
// The plan is to eventually move to a redblack tree
package main

import (
	"fmt"
)

const (
	red   bool = true
	black bool = false
)

// Internal Nodes - These will have a breakpoint but arcSite will be nil
// Leaf Nodes - These have an arcSite but no breakpoint
type node struct {
	left, right *node
	colour      bool
	breakpoint  *breakpoint
	arcSite     *site
	key         int
	circleEvent *Item
	edge        *halfEdge
}

type redblacktree struct {
	root *node
}

func (rbtree *redblacktree) insert(newKey int, newSite *site) {
	if rbtree.root == nil {
		rbtree.root = &node{colour: black, arcSite: newSite}
	} else {
		rbtree.root = rbtree.root.insert(rbtree.root, newKey, newSite)
	}
}

func (n *node) insert(currentNode *node, newKey int, newSite *site) *node {
	// Check if this is a leaf node
	if currentNode.breakpoint == nil {
		// TODO - check if currentNode has a circle event and if it does this needs to be removed
		// from the event queue as it is a false alarm

		// Define the breakpoints that will be used in the two new internal nodes
		leftBreakpoint := breakpoint{leftSite: currentNode.arcSite, rightSite: newSite}
		rightBreakpoint := breakpoint{leftSite: newSite, rightSite: currentNode.arcSite}

		// The 3 leaf nodes that represent the arcs
		leftLeafNode := node{arcSite: currentNode.arcSite}
		middleLeafNode := node{arcSite: newSite}
		rightLeafNode := node{arcSite: currentNode.arcSite}

		leftInternalNode := node{key: newKey, left: &leftLeafNode, right: &middleLeafNode, breakpoint: &leftBreakpoint}
		rightInternalNode := node{key: newKey, left: &leftInternalNode, right: &rightLeafNode, breakpoint: &rightBreakpoint}

		return &rightInternalNode
	}

	// The directrix will be at the same y coordinate as the new site being added
	breakpointXCoordinate := getBreakpointXCoordinate(currentNode.breakpoint, newSite.y)

	if float64(newSite.x) < breakpointXCoordinate {
		currentNode.left = currentNode.insert(currentNode.left, newKey, newSite)
	} else if float64(newSite.x) > breakpointXCoordinate {
		currentNode.right = currentNode.insert(currentNode.right, newKey, newSite)
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

func (rbtree *redblacktree) inorderTraversal() {
	rbtree.root.inorderTraversal(rbtree.root)
}

func (n *node) inorderTraversal(currentNode *node) {
	if currentNode != nil {
		currentNode.inorderTraversal(currentNode.left)
		fmt.Printf("%d ", currentNode.key)
		currentNode.inorderTraversal(currentNode.right)
	}
}
