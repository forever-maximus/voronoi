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
	if currentNode.key < keyToFind {
		return currentNode.search(currentNode.right, keyToFind)
	}
	return currentNode.search(currentNode.left, keyToFind)
}

func (rbtree *redblacktree) inorderTraversal() {
	rbtree.root.inorderTraversal(rbtree.root)
}

func (n *node) inorderTraversal(currentNode *node) {
	if currentNode != nil {
		currentNode.inorderTraversal(currentNode.left)
		fmt.Println(currentNode.key)
		currentNode.inorderTraversal(currentNode.right)
	}
}
