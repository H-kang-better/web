package msgo

import "testing"

func TestTreeNode(t *testing.T) {
	root := &treeNode{
		name:     "/",
		children: make([]*treeNode, 0),
	}
	root.Put("/user/get/:id")
	root.Put("/user/create/hello")
	root.Get("/user/get/1")
}
