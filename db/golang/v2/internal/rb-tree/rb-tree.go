package rb_tree

const RED = "Red"
const BLACK = "Black"

type Color struct {
	Red   string
	Black string
}
type Node struct {
	key    int
	color  string
	left   *Node
	right  *Node
	parent *Node
}

type RBTree struct {
	root *Node
}

type Iterator struct {
	tree            *RBTree
	currentNode     *Node
	currentPosition int // не уверен, что нужно знать
}

func New() *RBTree {
	return &RBTree{}
}

func (t *RBTree) InsertTree(key int) *Node {
	current := t.root
	// empty tree
	if t.root == nil {
		t.root = &Node{
			key:    key,
			color:  BLACK,
			left:   nil,
			right:  nil,
			parent: nil,
		}
		return t.root
	}

	var parent *Node
	for {
		if current == nil {
			//create current node
			currentNode := &Node{
				key:    key,
				color:  RED,
				left:   nil,
				right:  nil,
				parent: parent,
			}
			if key >= parent.key {
				parent.right = currentNode
			}

			if key <= parent.key {
				parent.left = currentNode
			}
			//похоже в этот момент нужно будет проводить балансировку
			t.fixInsert(currentNode)

			return currentNode
		}

		if key >= current.key {
			parent = current
			current = current.right
		}

		if key <= current.key {
			parent = current
			current = current.left
		}
	}
}

func (t *RBTree) find(key int) bool {
	current := t.root

	for current == nil {

		if current.key == key {
			return true
		}

		if key > current.key {
			current = current.right
		} else {
			current = current.left
		}

	}
	return false
}

func (t *RBTree) fixInsert(currentNode *Node) {
}
