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

	for current != nil {
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
	if currentNode.parent.color == BLACK {
		return
	}

	if currentNode.parent.color == RED {
		// если родитель красный, нужно проверить его "дядю"
		parent := currentNode.parent

		if parent == nil || parent.parent == nil {
			return
		}

		grandParent := parent.parent
		// нужно найти "дядю"
		var uncle *Node
		if grandParent.left == parent { // сравнивал сначала ключи, но кажется проще сравнивать ссылки
			uncle = grandParent.right // если родитель текущей ноды слева, то дядя справа
		} else {
			uncle = grandParent.left
		}
		// если дядя красный, перекрашиваем и проверяем дерево выше, возможно там нужно делать так же изменения
		if uncle != nil && uncle.color == RED {
			parent.color = BLACK
			uncle.color = BLACK
			grandParent.color = RED
			t.fixInsert(grandParent)
			return
		}
		//здесь кажется нужно делать поворот(как определить какой?), похоже,
		//что в этом случае измениться черная высота (хотя я ее никак не отслеживаю)
		if uncle == nil || uncle.color == BLACK {

		}

	}
}

// хз пока как их реализовать
func (t *RBTree) rotateLeft() {

}

func (t *RBTree) rotateRight() {

}
