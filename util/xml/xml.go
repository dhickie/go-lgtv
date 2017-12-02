package xml

import (
	"bytes"
	"encoding/xml"
	"errors"
)

// ErrNodeNotFound indicates that the specified node could not be found
var ErrNodeNotFound = errors.New("Node not found in xml")

// Node represents a node in some XML
type Node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []Node `xml:",any"`
}

// FindXMLNode finds the node with the specified name in the specified input XML string
func FindXMLNode(input, name string) (Node, error) {
	buf := bytes.NewBufferString(input)
	dec := xml.NewDecoder(buf)

	var n Node
	err := dec.Decode(&n)
	if err != nil {
		return Node{}, err
	}

	var m Node
	found := false
	walk([]Node{n}, func(n Node) bool {
		if n.XMLName.Local == name {
			found = true
			m = n
			return false
		}
		return true
	})

	if found {
		return m, nil
	}
	return Node{}, ErrNodeNotFound
}

func walk(nodes []Node, f func(Node) bool) {
	for _, n := range nodes {
		if f(n) {
			walk(n.Nodes, f)
		}
	}
}
