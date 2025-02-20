package rb_tree

import (
	"cmp"
	"fmt"
)

func createNewTestNode[K cmp.Ordered, D any](key K, data D, color Color, parent *Node[K, D]) *Node[K, D] {
	return &Node[K, D]{
		key:    key,
		data:   data,
		color:  color,
		parent: parent,
	}
}

func buildTestTree[K cmp.Ordered, D any](setup func(tree *RBTree[K, D])) *RBTree[K, D] {
	tree := New[K, D]()
	setup(tree)
	return tree
}

func createSimpleTree() (tree *RBTree[int, string], root, childLeft *Node[int, string]) {
	tree = buildTestTree[int, string](func(tree *RBTree[int, string]) {
		// Создаём корневой узел без родителя.
		root = createNewTestNode(10, "data 10", BLACK, nil)
		// Создаём левый дочерний узел с родителем root.
		childLeft = createNewTestNode(5, "data 5", RED, root)
		// Устанавливаем связь.
		root.left = childLeft
		tree.root = root
	})
	return
}

func createRecoloringTree() (tree *RBTree[int, string], root, childLeft, childRight, node *Node[int, string]) {
	tree = buildTestTree[int, string](func(tree *RBTree[int, string]) {
		root = createNewTestNode(10, "data-10", BLACK, nil)
		childLeft = createNewTestNode(5, "data-5", RED, root)
		childRight = createNewTestNode(15, "data-15", RED, root)
		root.left = childLeft
		root.right = childRight
		childLeft.parent = root
		childRight.parent = root
		tree.root = root

		node = createNewTestNode(20, "data-20", RED, childRight)
		childRight.right = node
	})
	return
}

func createLeftRotateTree() (tree *RBTree[int, string], root, parent, node *Node[int, string]) {
	tree = buildTestTree[int, string](func(tree *RBTree[int, string]) {
		root = createNewTestNode(10, "data-10", BLACK, nil)
		// Родитель здесь — левый ребёнок корня.
		parent = createNewTestNode(5, "data-5", RED, root)
		root.left = parent
		parent.parent = root
		tree.root = root

		// createNewTestNode вставляется как правый ребёнок родителя.
		node = createNewTestNode(7, "data-7", RED, parent)
		parent.right = node
	})
	return
}

func createRightRotateTree() (tree *RBTree[int, string], root, parent, node *Node[int, string]) {
	tree = buildTestTree[int, string](func(tree *RBTree[int, string]) {
		root = createNewTestNode(10, "data-10", BLACK, nil)
		// Родитель здесь — правый ребёнок корня.
		parent = createNewTestNode(15, "data-15", RED, root)
		root.right = parent
		parent.parent = root
		tree.root = root

		// createNewTestNode вставляется как левый ребёнок родителя.
		node = createNewTestNode(13, "data-13", RED, parent)
		parent.left = node
	})
	return
}

func createLongTree() *RBTree[int, string] {
	tree := New[int, string]()
	arr := []int{11, 1, 12, 2, 13, 3, 14, 4, 15, 5, 16, 6, 17, 7, 18, 8, 19, 9, 20}
	for i, key := range arr {
		tree.InsertTree(key, fmt.Sprintf("data-%d", i))
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

// cloneTree и cloneNode нужны для бенчмарков, чтобы убрать подсчет в бенчмаркох алокаций на создание дерева,
// а считать только алокации при удалении
func cloneTree[K cmp.Ordered, V any](tree *RBTree[K, V]) *RBTree[K, V] {
	newTree := New[K, V]()
	newTree.root = cloneNode(tree.root, nil)
	return newTree
}

// cloneNode создает глубокую копию узла и всех его поддеревьев.
// parent — родительский узел для вновь созданного клона.
func cloneNode[K cmp.Ordered, V any](node *Node[K, V], parent *Node[K, V]) *Node[K, V] {
	if node == nil {
		return nil
	}
	newNode := &Node[K, V]{
		key:    node.key,
		color:  node.color,
		parent: parent,
	}
	newNode.left = cloneNode(node.left, newNode)
	newNode.right = cloneNode(node.right, newNode)
	return newNode
}
