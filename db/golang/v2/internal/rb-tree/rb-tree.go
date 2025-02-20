package rb_tree

import (
	"cmp"
	"fmt"
)

type Color string

const (
	RED   Color = "Red"
	BLACK Color = "Black"
)

type Node[K cmp.Ordered, D any] struct {
	key    K
	data   D
	color  Color
	left   *Node[K, D]
	right  *Node[K, D]
	parent *Node[K, D]
}

type RBTree[K cmp.Ordered, D any] struct {
	root *Node[K, D]
}

func New[K cmp.Ordered, D any]() *RBTree[K, D] {
	return &RBTree[K, D]{}
}

func (t *RBTree[K, D]) InsertTree(key K, data D) *Node[K, D] {
	current := t.root
	// empty tree
	if t.root == nil {
		t.root = &Node[K, D]{
			key:    key,
			data:   data,
			color:  BLACK,
			left:   nil,
			right:  nil,
			parent: nil,
		}
		return t.root
	}

	var parent *Node[K, D]
	for {
		if current == nil {
			//create current node
			currentNode := &Node[K, D]{
				key:    key,
				data:   data,
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

func (t *RBTree[K, D]) fixInsert(currentNode *Node[K, D]) {
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
		var uncle *Node[K, D]
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

func (t *RBTree[K, D]) rotateLeft(node *Node[K, D]) {
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

func (t *RBTree[K, D]) rotateRight(node *Node[K, D]) {
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

func (t *RBTree[K, D]) findNode(key K) *Node[K, D] {
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

// Delete удаляет узел с указанным ключом из красно-чёрного дерева
// сначала находит целевой узел, затем, если у него два потомка, заменяет его на предшественника,
// после чего корректирует дерево, чтобы сохранить свойства красно-чёрного дерева.
func (t *RBTree[K, D]) Delete(key K) {
	var childNode *Node[K, D]
	targetNode := t.findNode(key)
	if targetNode == nil {
		return
	}
	// если у узла два потомка, находим предшественника (это максимальный узел в левом поддереве)
	if targetNode.left != nil && targetNode.right != nil {
		predecessorNode := targetNode.left.findMaxNode()
		targetNode.key = predecessorNode.key
		targetNode = predecessorNode
	}
	// если у узла один ребенок
	if targetNode.left == nil || targetNode.right == nil {
		if targetNode.right == nil {
			childNode = targetNode.left
		} else {
			childNode = targetNode.right
		}
		// если удаляемый узел чёрный, требуется балансировка
		if targetNode.color == BLACK {
			targetNode.color = nodeColor(childNode)
			t.propagateFixup(targetNode)
		}
		t.replaceNode(targetNode, childNode)
		// если удалённый узел был корнем, новый узел (если он есть) должен быть чёрным
		if targetNode.parent == nil && childNode != nil {
			childNode.color = BLACK
		}
	}
}

// findMaxNode возвращает узел с максимальным ключом в поддереве,
// перемещаясь к самому правому узлу.
func (n *Node[K, D]) findMaxNode() *Node[K, D] {
	if n == nil {
		return nil
	}
	for n.right != nil {
		n = n.right
	}
	return n
}

// propagateFixup обрабатывает случай удаления, когда узел не является корнем
// если узел уже стал корневым, дальнейшая корректировка не требуется.
func (t *RBTree[K, D]) propagateFixup(deletedNode *Node[K, D]) {
	if deletedNode.parent == nil {
		return
	}
	t.adjustRedSibling(deletedNode)
}

// adjustRedSibling обрабатывает случай, когда "брат" (sibling) удалённого узла красный,
// в этом случае происходит перекраска и поворот для поднятия проблемы выше по дереву
func (t *RBTree[K, D]) adjustRedSibling(deletedNode *Node[K, D]) {
	sibling := deletedNode.findSibling()
	if nodeColor(sibling) == RED {
		deletedNode.parent.color = RED
		sibling.color = BLACK
		if deletedNode == deletedNode.parent.left {
			t.rotateLeft(deletedNode.parent)
		} else {
			t.rotateRight(deletedNode.parent)
		}
	}
	t.balanceWithBlackNodes(deletedNode)
}

// balanceWithBlackNodes обрабатывает случай, когда родитель, брат и оба ребёнка брата чёрные.
// В этом случае брат перекрашивается в красный, а алгоритм рекурсивно продолжается для родителя.
func (t *RBTree[K, D]) balanceWithBlackNodes(deletedNode *Node[K, D]) {
	sibling := deletedNode.findSibling()
	if nodeColor(deletedNode.parent) == BLACK &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.left) == BLACK &&
		nodeColor(sibling.right) == BLACK {
		sibling.color = RED
		t.propagateFixup(deletedNode.parent)
	} else {
		t.adjustRedParent(deletedNode)
	}
}

// adjustRedParent обрабатывает случай, когда родитель красный, а брат и его потомки — чёрные.
// Здесь происходит обмен цвета между родителем и братом.
func (t *RBTree[K, D]) adjustRedParent(deletedNode *Node[K, D]) {
	sibling := deletedNode.findSibling()
	if nodeColor(deletedNode.parent) == RED &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.left) == BLACK &&
		nodeColor(sibling.right) == BLACK {
		sibling.color = RED
		deletedNode.parent.color = BLACK
	} else {
		t.rotateSiblingForBalance(deletedNode)
	}
}

// rotateSiblingForBalance обрабатывает случай, когда брат чёрный, а один из его детей красный,
// что позволяет выполнить поворот и подготовить ситуацию для финальной балансировки
func (t *RBTree[K, D]) rotateSiblingForBalance(deletedNode *Node[K, D]) {
	sibling := deletedNode.findSibling()
	if deletedNode == deletedNode.parent.left &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.left) == RED &&
		nodeColor(sibling.right) == BLACK {
		sibling.color = RED
		sibling.left.color = BLACK
		t.rotateRight(sibling)
	} else if deletedNode == deletedNode.parent.right &&
		nodeColor(sibling) == BLACK &&
		nodeColor(sibling.right) == RED &&
		nodeColor(sibling.left) == BLACK {
		sibling.color = RED
		sibling.right.color = BLACK
		t.rotateLeft(sibling)
	}
	t.finalizeDeletionBalance(deletedNode)
}

// finalizeDeletionBalance выполняет окончательную корректировку, устанавливая цвета брата и родителя
// и выполняет поворот для восстановления свойств красно-чёрного дерева
func (t *RBTree[K, D]) finalizeDeletionBalance(deletedNode *Node[K, D]) {
	sibling := deletedNode.findSibling()
	sibling.color = nodeColor(deletedNode.parent)
	deletedNode.parent.color = BLACK
	if deletedNode == deletedNode.parent.left && nodeColor(sibling.right) == RED {
		sibling.right.color = BLACK
		t.rotateLeft(deletedNode.parent)
	} else if nodeColor(sibling.left) == RED {
		sibling.left.color = BLACK
		t.rotateRight(deletedNode.parent)
	}
}

func (n *Node[K, D]) findSibling() *Node[K, D] {
	if n == nil || n.parent == nil {
		return nil
	}
	if n == n.parent.left {
		return n.parent.right
	}
	return n.parent.left
}

func (t *RBTree[K, D]) replaceNode(oldNode, newNode *Node[K, D]) {
	if oldNode.parent == nil {
		t.root = newNode
	} else {
		if oldNode == oldNode.parent.left {
			oldNode.parent.left = newNode
		} else {
			oldNode.parent.right = newNode
		}
	}
	if newNode != nil {
		newNode.parent = oldNode.parent
	}
}

func (t *RBTree[K, D]) transplant(target, replacement *Node[K, D]) {
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

func nodeColor[K cmp.Ordered, D any](node *Node[K, D]) Color {
	if node == nil {
		return BLACK
	}
	return node.color
}

func (t *RBTree[K, D]) Print() {
	if t.root == nil {
		fmt.Println("[Empty tree]")
		return
	}
	t.printNode(t.root, "", true)
}

func (t *RBTree[K, D]) printNode(node *Node[K, D], prefix string, isTail bool) {
	if node == nil {
		return
	}

	// Добавляем указатели на детей
	pointers := "├── "
	if isTail {
		pointers = "└── "
	}

	// Формируем цветовую метку
	color := "⚫️"
	if node.color == RED {
		color = "🔴"
	}

	// Выводим текущий узел
	fmt.Printf("%s%s%v(%s)-%v\n", prefix, pointers, node.key, color, node.data)

	// Вычисляем новый префикс для детей
	newPrefix := prefix
	if isTail {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	// Рекурсивно выводим детей
	if node.left != nil || node.right != nil {
		t.printNode(node.right, newPrefix, false)
		t.printNode(node.left, newPrefix, true)
	}
}
