package rb_tree

import (
	"fmt"
)

func createSimpleTree() (tree *RBTree[int, string], root, childLeft *Node[int, string]) {
	tree = New[int, string]()

	root = &Node[int, string]{
		key:   10,
		data:  "data 10",
		color: BLACK,
	}
	childLeft = &Node[int, string]{
		key:   5,
		data:  "data 5",
		color: RED,
	}
	root.left = childLeft
	childLeft.parent = root
	tree.root = root
	return
}

func createRecoloringTree() (tree *RBTree[int, string], root, childLeft, childRight, newNode *Node[int, string]) {
	tree = New[int, string]()
	root = &Node[int, string]{
		key:   10,
		data:  "data-10",
		color: BLACK,
	}
	childLeft = &Node[int, string]{
		key:   5,
		data:  "data-5",
		color: RED,
	}
	childRight = &Node[int, string]{
		key:   15,
		data:  "data-15",
		color: RED,
	}
	root.left = childLeft
	root.right = childRight
	childLeft.parent = root
	childRight.parent = root
	tree.root = root

	newNode = &Node[int, string]{
		key:    20,
		data:   "data-20",
		color:  RED,
		parent: childRight,
	}
	childRight.right = newNode
	return
}

func createLeftRotateTree() (tree *RBTree[int, string], root, parent, newNode *Node[int, string]) {
	// Родитель – левый ребёнок, новая нода вставляется как правый ребёнок родителя
	tree = New[int, string]()
	root = &Node[int, string]{
		key:   10,
		data:  "data-10",
		color: BLACK,
	}
	parent = &Node[int, string]{
		key:   5,
		data:  "data-5",
		color: RED,
	}
	root.left = parent
	parent.parent = root
	tree.root = root

	// newNode вставляется как правый ребёнок родителя
	newNode = &Node[int, string]{
		key:    7,
		data:   "data-7",
		color:  RED,
		parent: parent,
	}
	parent.right = newNode
	return
}

func createRightRotateTree() (tree *RBTree[int, string], root, parent, newNode *Node[int, string]) {
	//Родитель – правый ребёнок, новая нода вставляется как левый ребёнок родителя
	tree = New[int, string]()
	root = &Node[int, string]{
		key:   10,
		data:  "data-10",
		color: BLACK,
	}
	parent = &Node[int, string]{
		key:   15,
		data:  "data-15",
		color: RED,
	}
	root.right = parent
	parent.parent = root
	tree.root = root

	// newNode вставляется как левый ребёнок родителя.
	newNode = &Node[int, string]{
		key:    13,
		data:   "data-13",
		color:  RED,
		parent: parent,
	}
	parent.left = newNode
	return
}

func createLongTree() *RBTree[int, string] {
	arr := []int{11, 1, 12, 2, 13, 3, 14, 4, 15, 5, 16, 6, 17, 7, 18, 8, 19, 9, 20}
	tree := New[int, string]()
	for i := 0; i < len(arr); i++ {
		tree.InsertTree(arr[i], fmt.Sprintf("data-%d", i))
	}
	return tree
}

func checkRBInvariants(tree *RBTree[int, string]) error {
	if tree.root == nil {
		return nil
	}
	// 1) Корень должен быть чёрным
	if tree.root.color != BLACK {
		return fmt.Errorf("корень не чёрный")
	}
	// 2) Проверить, что нет двух подряд красных узлов
	if err := checkNoConsecutiveReds(tree.root); err != nil {
		return err
	}
	// 3) Проверить равную "чёрную высоту" по всем путям
	blackHeight := -1
	if err := checkBlackHeight(tree.root, 0, &blackHeight); err != nil {
		return err
	}
	return nil
}

// checkNoConsecutiveReds проверяет, что нет двух подряд красных узлов
func checkNoConsecutiveReds(node *Node[int, string]) error {
	if node == nil {
		return nil
	}
	if node.color == RED {
		if (node.left != nil && node.left.color == RED) ||
			(node.right != nil && node.right.color == RED) {
			return fmt.Errorf("найдены два подряд красных узла: %v", node.key)
		}
	}
	if err := checkNoConsecutiveReds(node.left); err != nil {
		return err
	}
	if err := checkNoConsecutiveReds(node.right); err != nil {
		return err
	}
	return nil
}

// checkBlackHeight проверяет, что все пути от корня до nil имеют одинаковую "чёрную высоту".
func checkBlackHeight(node *Node[int, string], currentBlackCount int, reference *int) error {
	if node == nil {
		// дошли до nil-узла
		if *reference == -1 {
			*reference = currentBlackCount
		} else if currentBlackCount != *reference {
			return fmt.Errorf("нарушена равная черная высота: %d != %d", currentBlackCount, *reference)
		}
		return nil
	}
	// если узел чёрный, увеличиваем счётчик
	nextBlackCount := currentBlackCount
	if node.color == BLACK {
		nextBlackCount++
	}
	// рекурсивно проверяем левую и правую ветви
	if err := checkBlackHeight(node.left, nextBlackCount, reference); err != nil {
		return err
	}
	if err := checkBlackHeight(node.right, nextBlackCount, reference); err != nil {
		return err
	}
	return nil
}
