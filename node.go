package registry

import (
	"strings"
)

type Node struct {
	Name       string
	RouterName string
	Children   []*Node
	IsEnd      bool
	Data       *[]Address
}

var defaultSplit = "/"

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Children: make([]*Node, 0),
	}
}

func (t *Node) Put(path string, data *[]Address) {
	root := t
	paths := strings.Split(path, defaultSplit)
	routerName := ""
	prefix := strings.HasPrefix(path, defaultSplit)
	for index, name := range paths {
		if prefix && index == 0 {
			continue
		}
		children := t.Children
		isMatch := false
		for _, childNode := range children {
			// 如果user匹配到了，下一次就开始判断get
			if childNode.Name == name {
				isMatch = true
				routerName += defaultSplit + childNode.Name
				childNode.RouterName = routerName
				t = childNode
				break
			}
		}
		if !isMatch {
			isEnd := false
			if index == len(paths)-1 {
				isEnd = true
			}
			// 没有匹配到，那么这是一个新的路径，创建一个节点对象
			childNode := &Node{Name: name, Children: make([]*Node, 0), IsEnd: isEnd}
			if isEnd {
				childNode.Data = data
			}
			routerName += defaultSplit + childNode.Name
			childNode.RouterName = routerName
			children = append(children, childNode)
			t.Children = children
			t = childNode
		}
	}
	t = root
}

func (t *Node) Get(path string) *Node {
	paths := strings.Split(path, defaultSplit)
	prefix := strings.HasPrefix(path, defaultSplit)
	residuePath := path
	for index, name := range paths {
		if prefix && index == 0 {
			continue
		}
		children := t.Children
		for _, childNode := range children {
			if childNode.Name == name {

				residuePath = strings.TrimPrefix(residuePath, defaultSplit+name)

				t = childNode
				if index == len(paths)-1 {
					return childNode
				}
				break
			}
		}
	}
	return nil
}

func (t *Node) delete(path string) {
	paths := strings.Split(path, defaultSplit)
	prefix := strings.HasPrefix(path, defaultSplit)

	residuePath := path
	for index, name := range paths {
		if prefix && index == 0 {
			continue
		}
		children := t.Children
		for _, childNode := range children {
			if childNode.Name == name {

				residuePath = strings.TrimPrefix(residuePath, defaultSplit+name)

				t = childNode
				if index == len(paths)-1 {
					childNode = nil
				}
				break
			}
		}
	}
}

func Iterator(node *Node, fun func(key string, Data *[]Address)) {
	for _, children := range node.Children {
		if children.IsEnd {
			fun(children.RouterName, children.Data)
		} else {
			Iterator(children, fun)
		}
	}
}
