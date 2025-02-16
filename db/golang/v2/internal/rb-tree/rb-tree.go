package rb_tree

const RED = "Red"
const BLACK = "Black"

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

			if key > parent.key {
				parent.right = currentNode
			}

			if key < parent.key {
				parent.left = currentNode
			}
			//похоже в этот момент нужно будет проводить балансировку
			t.fixInsert(currentNode)

			return currentNode
		}

		if key > current.key {
			parent = current
			current = current.right
		} else {
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
	if currentNode.parent != nil && currentNode.parent.color == BLACK {
		return
	}

	if currentNode.parent != nil && currentNode.parent.color == RED {
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
		// если дядя красный, перекрашиваем и проверяем дерево выше, возможно там нужно делать также изменения
		if uncle != nil && uncle.color == RED {
			parent.color = BLACK
			uncle.color = BLACK
			grandParent.color = RED
			t.fixInsert(grandParent)
			return
		} else {
			if parent == grandParent.left {
				if currentNode == parent.right { //если текущая нода в правой ветке, то делаем левый поворот
					t.rotateLeft(parent)
					currentNode = parent
					parent = currentNode.parent
				}
				t.rotateRight(grandParent) //иначе правый
			} else {
				if currentNode == parent.left { //если текущая нода в левой ветке, то делаем правый поворот
					t.rotateRight(parent)
					currentNode = parent
					parent = currentNode.parent
				}
				t.rotateLeft(grandParent)
			}
		}
		// случай когда не нужны повороты
		parent.color = BLACK
		grandParent.color = RED
	}
	t.root.color = BLACK
}

func (t *RBTree) rotateLeft(node *Node) {
	rightChild := node.right
	if rightChild == nil {
		return
	}

	node.right = rightChild.left
	if rightChild.left != nil {
		rightChild.left.parent = node
	}

	rightChild.parent = node.parent
	if node.parent == nil {
		t.root = rightChild
	} else if node == node.parent.left {
		node.parent.left = rightChild
	} else {
		node.parent.right = rightChild
	}

	rightChild.left = node
	node.parent = rightChild
}

func (t *RBTree) rotateRight(node *Node) {
	leftChild := node.left
	if leftChild == nil {
		return
	}

	node.left = leftChild.right
	if leftChild.right != nil {
		leftChild.right.parent = node
	}

	leftChild.parent = node.parent
	if node.parent == nil {
		t.root = leftChild
	} else if node == node.parent.right {
		node.parent.right = leftChild
	} else {
		node.parent.left = leftChild
	}

	leftChild.right = node
	node.parent = leftChild
}

func (t *RBTree) findNode(key int) *Node {
	current := t.root

	for current != nil {
		if current.key == key {
			return current
		}
		if key > current.key {
			current = current.right
		} else {
			current = current.left
		}
	}
	return nil
}

func (t *RBTree) Delete(key int) {
	target := t.findNode(key)
	if target == nil {
		return
	}

	originalColor := target.color
	var nodeToFix *Node

	// кейс когда у target один ребенок
	if target.left == nil || target.right == nil {
		nodeToFix = t.deleteSingleChild(target)
	} else {
		// кейс, с двумя детьми
		nodeToFix = t.deleteTwoChildren(target, &originalColor)
	}

	if originalColor == BLACK {
		t.fixDelete(nodeToFix)
	}
}

func (t *RBTree) deleteSingleChild(target *Node) *Node {
	var child *Node
	if target.left != nil {
		child = target.left
	} else {
		child = target.right
	}
	t.transplant(target, child)
	return child
}

// deleteTwoChildren обрабатывает случай, когда у узла target два ребенка
// он находит наследника (минимальный узел правого поддерева) -> заменяет target наследником -> возвращает узел
// ‼️ для него может потребоваться балансировка и при этом originalColor нужно обновить цветом наследника.
func (t *RBTree) deleteTwoChildren(target *Node, originalColor *string) *Node {
	// ищем наследника
	successor := target.right
	for successor.left != nil {
		successor = successor.left
	}

	*originalColor = successor.color
	replaceNode := successor.right
	//todo думаю стоит написать комментарии что происходит в коде ниже
	if successor.parent != target {
		t.transplant(successor, successor.right)
		successor.right = target.right
		if successor.right != nil {
			successor.right.parent = successor
		}
	}

	t.transplant(target, successor)
	successor.left = target.left
	if successor.left != nil {
		successor.left.parent = successor
	}
	successor.color = target.color

	return replaceNode
}

func (t *RBTree) fixDelete(node *Node) {

}

func (t *RBTree) transplant(node *Node, node2 *Node) {

}
