package scenegraph

import (
	"go-game-engine/entity"
	"go-game-engine/types"
)

// Transform describes a transform of an object
type Transform struct {
	position types.Vector3
	scale    types.Vector3
}

// Node represents a node in the scene graph
type Node struct {
	Parent    *Node
	Entity    *entity.Entity
	Children  []*Node
	Transform Transform
}

var scene *Node

// CreateNode creates a node
func CreateNode(parent *Node) *Node {
	nd := Node{
		Parent:    parent,
		Entity:    nil,
		Children:  make([]*Node, 0),
		Transform: Transform{},
	}

	return &nd
}

// CreateScene initialises the scene
func CreateScene() *Node {
	scene = CreateNode(nil)

	return scene
}
