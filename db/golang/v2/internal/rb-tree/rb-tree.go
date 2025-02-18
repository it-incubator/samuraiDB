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
	// кейс когда у target нет детей
	if target.left == nil && target.right == nil {
		if target == t.root {
			t.root = nil
		} else {
			if target.parent.left == target {
				target.parent.left = nil
			} else {
				target.parent.right = nil
			}
		}
		return
	}

	// кейс когда у target один ребенок
	if target.left == nil || target.right == nil {
		nodeToFix = t.deleteSingleChild(target)
	} else {
		// кейс, с двумя детьми
		nodeToFix = t.deleteTwoChildren(target, &originalColor)
	}

	if originalColor == BLACK && nodeToFix != nil {
		t.fixDelete(nodeToFix)
	}
}

func (t *RBTree) fixDelete(currentNode *Node) {
	for currentNode != nil && currentNode != t.root && currentNode.color == BLACK {
		if currentNode == currentNode.parent.left {
			currentNode = t.fixDeleteLeft(currentNode)
		} else {
			currentNode = t.fixDeleteRight(currentNode)
		}
	}
	if currentNode != nil {
		currentNode.color = BLACK
	}
}

func (t *RBTree) fixDeleteLeft(currentNode *Node) *Node {
	parent := currentNode.parent
	sibling := parent.right
	if sibling == nil {
		return parent
	}
	// брат красный
	if sibling.color == RED {
		sibling.color = BLACK
		parent.color = RED
		t.rotateLeft(parent)
		sibling = parent.right
	}
	// оба ребёнка брата чёрные
	if (sibling.left == nil || sibling.left.color == BLACK) &&
		(sibling.right == nil || sibling.right.color == BLACK) {
		sibling.color = RED
		return parent
	}
	// правый ребёнок брата чёрный, а левый красный, делаем правый поворот на брате
	if sibling.right == nil || sibling.right.color == BLACK {
		if sibling.left != nil {
			sibling.left.color = BLACK
		}
		sibling.color = RED
		t.rotateRight(sibling)
		sibling = parent.right
	}
	// правый ребёнок брата красный, делаем левый поворот на родителе
	sibling.color = parent.color
	parent.color = BLACK
	if sibling.right != nil {
		sibling.right.color = BLACK
	}
	t.rotateLeft(parent)
	return t.root
}

func (t *RBTree) fixDeleteRight(currentNode *Node) *Node {
	parent := currentNode.parent
	sibling := parent.left
	if sibling == nil {
		return parent
	}
	if sibling.color == RED {
		sibling.color = BLACK
		parent.color = RED
		t.rotateRight(parent)
		sibling = parent.left
	}
	if (sibling.left == nil || sibling.left.color == BLACK) &&
		(sibling.right == nil || sibling.right.color == BLACK) {
		sibling.color = RED
		return parent
	}
	if sibling.left == nil || sibling.left.color == BLACK {
		if sibling.right != nil {
			sibling.right.color = BLACK
		}
		sibling.color = RED
		t.rotateLeft(sibling)
		sibling = parent.left
	}
	sibling.color = parent.color
	parent.color = BLACK
	if sibling.left != nil {
		sibling.left.color = BLACK
	}
	t.rotateRight(parent)
	return t.root
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
// он находит наследника -> заменяет target наследником -> возвращает узел
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

func (t *RBTree) transplant(target, replacement *Node) {
	// кейс когда target корень
	if target.parent == nil {
		t.root = replacement
		// определяем в каком узле происходит замена
	} else if target.parent.left != target {
		target.parent.right = replacement
	} else {
		target.parent.left = replacement
	}
	// обмен родителями
	if replacement != nil {
		replacement.parent = target.parent
	}
}

func (t *RBTree) deleteTwoChildrenOld(target *Node, originalColor *string) *Node {
	predecessor := target.left
	for predecessor.right != nil {
		predecessor = predecessor.right
	}
	*originalColor = predecessor.color

	var replacement *Node
	if predecessor.left == nil {
		replacement = predecessor
	} else {
		replacement = predecessor.left
	}

	// Если старый узел не является непосредственным левым ребёнком target,
	// то перемещаем его левое поддерево на его место
	if predecessor.parent != target {
		t.transplant(predecessor, predecessor.left)
		predecessor.left = target.left
		predecessor.left.parent = predecessor
	}

	// Заменяем target предшественником.
	t.transplant(target, predecessor)
	predecessor.right = target.right
	predecessor.right.parent = predecessor
	predecessor.color = target.color

	if replacement != nil && *originalColor == BLACK {
		replacement.color = BLACK
	}

	return replacement
}

// processBlackSibling проверяет, что оба ребенка узла sibling отсутствуют или черные.
// Если условие выполнено, функция устанавливает sibling в красный и возвращает true,
// а также возвращает родителя (newCurrent) в качестве нового currentNode для балансировки.
func processBlackSibling(sibling, parent *Node, currentNode *Node) (handled bool, newCurrent *Node) {

	if sibling == nil {
		return false, currentNode
	}

	if (sibling.left == nil || sibling.left.color == BLACK) &&
		(sibling.right == nil || sibling.right.color == BLACK) {
		sibling.color = RED
		return true, parent
	}
	return false, currentNode
}
