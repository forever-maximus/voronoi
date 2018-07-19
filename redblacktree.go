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

type node struct {
	left, right *node
	colour      bool
	sites       []site
	key         int
	circleEvent *Item
	edge        *halfEdge
}

type redblacktree struct {
	root *node
}

func (rbtree *redblacktree) insert(newKey int) {
	if rbtree.root == nil {
		rbtree.root = &node{colour: black, key: newKey}
	} else {
		rbtree.root.insert(rbtree.root, newKey)
	}
}

func (n *node) insert(currentNode *node, newKey int) *node {
	if currentNode == nil {
		return &node{colour: red, key: newKey}
	}

	if newKey < currentNode.key {
		currentNode.left = currentNode.insert(currentNode.left, newKey)
	} else if newKey > currentNode.key {
		currentNode.right = currentNode.insert(currentNode.right, newKey)
	}

	return currentNode
}

func (rbtree *redblacktree) search(keyToFind int) *node {
	if rbtree.root == nil || rbtree.root.key == keyToFind {
		return rbtree.root
	}
	return rbtree.root.search(rbtree.root, keyToFind)
}

func (n *node) search(currentNode *node, keyToFind int) *node {
	if currentNode == nil || currentNode.key == keyToFind {
		return currentNode
	}
	if keyToFind > currentNode.key {
		return currentNode.search(currentNode.right, keyToFind)
	}
	return currentNode.search(currentNode.left, keyToFind)
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
