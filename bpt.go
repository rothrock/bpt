package bpt

import (
	"fmt"
	"sort"
)

type Record struct {
	key   string
	value string
}

type ByKey []Record

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].key < a[j].key }

// There are 4 types of nodes in a B+Tree.
// Initially the root serves as a leaf until it splits.
// Subsequent leaf splits will eventually spawn internal
// nodes between the root and the leaf.
type Node struct {
	whatAmI   string // rootleaf, root, internal, leaf
	records   []Record
	childPtrs []*Node
	next      *Node
}

type Tree struct {
	root *Node
}

func NewBPT() *Tree {
	newRoot := new(Node)
	newRoot.whatAmI = "rootleaf"
	newTree := Tree{root: newRoot}
	return &newTree
}

// Searches the tree n for the record associated with the given key.
func (t *Tree) Find(key string) (Record, bool, error) {

	n := t.root
	leaf, _, idx, found, err := findLeaf(n, n, key)
	if err != nil {
		return Record{key: "", value: ""}, false, err
	}
	if found {
		return leaf.records[idx], true, err
	}
	return Record{key: "", value: ""}, false, err
}

// Inserts a record on the leaf of a tree
// rooted at n. If the root splits, the returned
// node will be the root of the new tree.
func (t *Tree) Insert(r Record) error {

	n := t.root
	leaf, parent, _, found, err := findLeaf(n, n, r.key)
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("Key collision")
	}

	// Add the record
	leaf.records = append(leaf.records, r)
	sort.Sort(ByKey(leaf.records))

	if len(leaf.records) == 6 {
		// split needed
		right := new(Node)
		right.whatAmI = "leaf"
		right.records = make([]Record, 3)
		copy(right.records, leaf.records[3:])
		leaf.records = leaf.records[:3]
		leaf.next = right
		t.root, err = insertNode(n, parent, right.records[0].key, right)
	}
	return nil
}

// Inserts the key and associated node pointer on the target.
func insertKeyAndPtr(target *Node, key string, n *Node) {
	var i int
	var r Record
	for i, r = range target.records {
		if key < r.key {
			break
		}
	}

	rtmp := make([]Record, 0)
	ctmp := make([]*Node, 0)

	switch {
	case key > r.key:
		target.records = append(target.records, Record{key: key, value: ""})
		target.childPtrs = append(target.childPtrs, n)
	case i == 0:
		rtmp = make([]Record, 1)
		ctmp = make([]*Node, 1)
		rtmp[0] = Record{key: key, value: ""}
		rtmp = append(rtmp, target.records...)

		ctmp[0] = target.childPtrs[0]
		ctmp = append(ctmp, n)
		ctmp = append(ctmp, target.childPtrs[1:]...)
		target.records = rtmp
		target.childPtrs = ctmp
	default:
		rtmp = append(rtmp, target.records[:i]...)
		rtmp = append(rtmp, Record{key: key, value: ""})
		rtmp = append(rtmp, target.records[i:]...)

		i++ // Since we have more ptrs than records.
		ctmp = append(ctmp, target.childPtrs[:i]...)
		ctmp = append(ctmp, n)
		ctmp = append(ctmp, target.childPtrs[i:]...)
		target.records = rtmp
		target.childPtrs = ctmp

	}
}

// Splits a root node returning a new root with
// 'target' and 'right' as child pointers.
// This function does not apply to a rootleaf node.
func splitRoot(target *Node, right *Node, splitKey string) *Node {
	target.whatAmI = "internal"
	root := new(Node)
	root.whatAmI = "root"
	root.childPtrs = make([]*Node, 2)
	root.childPtrs[0] = target
	root.childPtrs[1] = right
	root.records = make([]Record, 1)
	root.records[0] = Record{key: splitKey, value: ""}
	return root
}

//Splits and internal node. The function updates the target
// and returns the newly created node.
func splitNode(root *Node, target *Node, key string) *Node {

	right := new(Node)
	right.whatAmI = "internal"
	right.records = make([]Record, 3)
	right.childPtrs = make([]*Node, 4)

	copy(right.records, target.records[3:])
	target.records = target.records[:2]

	copy(right.childPtrs, target.childPtrs[3:])
	target.childPtrs = target.childPtrs[:3]

	return right
}

// Inserts the key and related node pointer (n) on the target node.
// Splits if necessary, and returns a new root if needed.
// This function is not for inserting records on a leaf.
// The target should be one of root, rootleaf, or internal.
func insertNode(root *Node, target *Node, key string, n *Node) (*Node, error) {
	if target.whatAmI == "rootleaf" {
		newRoot := new(Node)
		target.whatAmI = "leaf"
		newRoot.whatAmI = "root"
		newRoot.records = append(newRoot.records, Record{key: key, value: ""})
		newRoot.childPtrs = append(newRoot.childPtrs, root, n)
		return newRoot, nil
	}

	insertKeyAndPtr(target, key, n)

	// Nodes need to split once they get to 6 entries.
	if len(target.records) == 6 { //split needed

		splitKey := target.records[2].key
		destinationNode, err := findParent(root, target)
		if err != nil {
			return nil, fmt.Errorf("Can't split because...", err)
		}

		right := splitNode(root, target, splitKey)

		// Need to make a new root node.
		if target.whatAmI == "root" {
			return splitRoot(target, right, splitKey), nil
		}

		root, err = insertNode(root, destinationNode, splitKey, right)
		if err != nil {
			return nil, fmt.Errorf("Recursive call to insertNode failed because of ", err)
		}
	}

	return root, nil
}

// Attempts to find the node that is the parent of the target
// node. Search starts from the given root.
func findParent(root *Node, target *Node) (*Node, error) {

	if root == target {
		return root, nil
	}

	if root.whatAmI == "leaf" {
		err := fmt.Errorf("Can't find parent starting from a leaf")
		return nil, err
	}

	key := target.records[0].key
	for i, r := range root.records {
		if key < r.key {
			if root.childPtrs[i] == target {
				return root, nil
			} else {
				return findParent(root.childPtrs[i], target)
			}
		}
	}
	return findParent(root.childPtrs[len(root.records)], target)
}

// Finds the relevant leaf node for the given key.
// Returns the leaf node, its parent, the key index,
// true or false if the key was actually found, and possible errors.
// This function is used to test where a record should end up
// or to actually find the record associated with the given key.
func findLeaf(n *Node, parent *Node, key string) (*Node, *Node, int, bool, error) {

	if n == nil {
		return n, parent, 0, false, fmt.Errorf("n is nil")
	}
	if len(n.childPtrs) == 0 { //we're on a leaf,
		for i, r := range n.records {
			if key == r.key {
				return n, parent, i, true, nil
			}
		}
		return n, parent, 0, false, nil
	}

	for i, r := range n.records {
		if key < r.key {
			return findLeaf(n.childPtrs[i], n, key)
		}
		if key == r.key {
			return findLeaf(n.childPtrs[i+1], n, key)
		}
	}

	return findLeaf(n.childPtrs[len(n.records)], n, key)
}
