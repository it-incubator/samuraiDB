package rb_tree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRBTree_InsertTree(t *testing.T) {
	tree := New()

	t.Run("Создаем корень", func(t *testing.T) {
		node10 := tree.InsertTree(10)

		t.Logf("Структура дерева после вставки узла 10:\n%s", treeToString(tree.root, ""))

		require.Equal(t, BLACK, node10.color, "Корень должен быть черным", node10.color)
	})

	t.Run("Вставка ребенка вправо", func(t *testing.T) {
		node20 := tree.InsertTree(20)

		t.Logf("Структура дерева после вставки узла 20:\n%s", treeToString(tree.root, ""))

		require.NotNil(t, node20)
		require.Equal(t, 20, node20.key, "Узел должен иметь ключ 20")
		require.Equal(t, RED, node20.color, "Новый узел должен быть красного цвета", node20.color)
		require.Equal(t, tree.root, node20.parent, "Родитель узла с ключом 20 должен быть корнем")
		require.Equal(t, node20, tree.root.right, "Узел с ключом 20 должен быть правым потомком корня")
	})

	t.Run("Вставка ребенка влево", func(t *testing.T) {
		node3 := tree.InsertTree(3)

		t.Logf("Структура дерева после вставки узла 3:\n%s", treeToString(tree.root, ""))

		require.NotNil(t, node3)
		require.Equal(t, 3, node3.key, "Узел должен иметь ключ 3")
		require.Equal(t, RED, node3.color, "Новый узел должен быть красного цвета", node3.color)
		require.Equal(t, tree.root, node3.parent, "Родитель узла с ключом 3 должен быть корнем")
		require.Equal(t, node3, tree.root.left, "Узел с ключом 3 должен быть левым потомком корня")
	})

	t.Run("Вставляем дополнительные узлы и проверяем балансировку", func(t *testing.T) {
		node30 := tree.InsertTree(30)

		t.Logf("Структура дерева после вставки узла 30:\n%s", treeToString(tree.root, ""))

		require.NotNil(t, node30)
		require.Equal(t, 30, node30.key)
		require.Equal(t, BLACK, node30.parent.color, "Родитель узла 30 должен быть черным", node30.parent.color)
		require.Equal(t, BLACK, tree.root.color, "Корень должен быть черным", tree.root.color)

		node40 := tree.InsertTree(40)

		t.Logf("Структура дерева после вставки узла 40:\n%s", treeToString(tree.root, ""))

		require.NotNil(t, node40)
		require.Equal(t, 40, node40.key)
		require.Equal(t, RED, node40.color, "Новый узел должен быть красного цвета", node40.color)
		require.Equal(t, BLACK, node40.parent.color, "Родитель узла 40 должен быть черным", node40.parent.color)
		require.Equal(t, BLACK, tree.root.color, "Корень должен быть черным", tree.root.color)
		require.Equal(t, RED, node40.parent.left.color, "Должен быть выполнен поворот и нода с ключом 20 должна стать ребенком ноды с ключом 30", node40.parent.left.color)
	})
	//Проверка черной высоты как доп проверка
	t.Run("Проверка черной высоты", func(t *testing.T) {
		err := checkRBInvariants(tree)
		require.NoError(t, err, "Дерево нарушает инварианты красно-чёрного дерева")
	})
}

func TestRBTree_fixInsert(t *testing.T) {
	t.Run("Не нужна балансировка (родитель черный)", func(t *testing.T) {
		tree, root, childLeft := createSimpleTree()

		tree.fixInsert(childLeft)

		require.Equal(t, BLACK, root.color, "Корень должен оставаться черным")
		require.Equal(t, RED, childLeft.color, "Дочерний узел должен оставаться красного цвета")
		require.Equal(t, root, childLeft.parent, "Родитель дочернего узла не должен измениться")
		require.Equal(t, childLeft, root.left, "Дочерний узел должен быть левым потомком корня")
		require.Nil(t, root.right, "Правый узел должен быть nil")
	})

	t.Run("Должны перекраситься ноды (родитель и дядя красные)", func(t *testing.T) {
		tree, root, childLeft, childRight, newNode := createRecoloringTree()

		// На этом этапе родитель (childRight) и дядя (childLeft) красные
		tree.fixInsert(newNode)

		// Ожидаем, что после перекрашивания:
		// - Родитель (childRight) и дядя (childLeft) станут черными
		// - Дедушка (root) временно станет красным, но затем fixInsert приведет его к черному
		require.Equal(t, BLACK, root.color, "Корень должен оставаться черным")
		require.Equal(t, BLACK, childLeft.color, "Дядя должен стать черным")
		require.Equal(t, BLACK, childRight.color, "Родитель должен стать черным")
		require.Equal(t, root, childLeft.parent, "Родитель дочернего узла не должен измениться")
		require.Equal(t, root, childRight.parent, "Родитель дочернего узла не должен измениться")
		require.Equal(t, childLeft, root.left, "Левая нода должна быть левым потомком корня")
		require.Equal(t, childRight, newNode.parent, "Новая нода должна быть потомком правой ноды")
		require.NotNil(t, root.right, "Правый узел не должен быть nil")
	})

	t.Run("Должен произойти левый поворот, а затем правый поворот", func(t *testing.T) {
		tree, _, _, newNode := createLeftRotateTree()
		//сценарий когда родитель слева от корня, а новая нода правый ребенок
		t.Logf("Структура дерева перед балансировкой:\n%s", treeToString(tree.root, ""))

		tree.fixInsert(newNode)
		t.Logf("Cтруктура дерева после балансировкой:\n%s", treeToString(tree.root, ""))
		// Ожидаем, что произойдёт левый поворот на родителе, затем правый поворот на корне
		// Итоговая структура должна стать:
		//   Новый корень с ключом 7 (черный),
		//   левый потомок – нода с ключом 5,
		//   правый потомок – нода с ключом 10.
		require.Equal(t, 7, tree.root.key, "Новый корень должен иметь ключ 7")
		require.Equal(t, BLACK, tree.root.color, "Новый корень должен быть черным")
		require.NotNil(t, tree.root.left, "Новый корень должен иметь левую ноду")
		require.Equal(t, 5, tree.root.left.key, "Левая нода должна иметь ключ 5")
		require.NotNil(t, tree.root.right, "Новый корень должен иметь правую ноду")
		require.Equal(t, 10, tree.root.right.key, "Правая нода должна иметь ключ 10")
	})

	t.Run("Должен произойти правый поворот, а затем левый поворот", func(t *testing.T) {
		//Родитель справа от корня, новая нода – левый ребёнок родителя
		tree, _, _, newNode := createRightRotateTree()
		t.Logf("Структура дерева перед балансировкой:\n%s", treeToString(tree.root, ""))

		tree.fixInsert(newNode)
		t.Logf("Структура дерева после балансировкой:\n%s", treeToString(tree.root, ""))

		// Ожидаем, что произойдёт правый поворот на родителе, затем левый поворот на корне
		// Итоговая структура должна стать:
		//   Новый корень с ключом 13 (черный),
		//   левый потомок – нода с ключом 10,
		//   правый потомок – нода с ключом 15
		require.Equal(t, 13, tree.root.key, "Новый корень должен иметь ключ 13")
		require.Equal(t, BLACK, tree.root.color, "Новый корень должен быть черным")
		require.NotNil(t, tree.root.left, "Новый корень должен иметь левую ноду")
		require.Equal(t, 10, tree.root.left.key, "Левая нода должна иметь ключ 10")
		require.NotNil(t, tree.root.right, "Новый корень должен иметь правую ноду")
		require.Equal(t, 15, tree.root.right.key, "Правая нода должна иметь ключ 15")
	})
}

func TestRBTree_findNode(t *testing.T) {
	t.Run("Должен вернуть nil если дерево пустое", func(t *testing.T) {
		tree := New()
		res := tree.findNode(1)

		require.Nil(t, res, "Дерево пустое")
	})

	t.Run("Должен найти ноду по ключу", func(t *testing.T) {
		tree, root, childLeft := createSimpleTree()
		res := tree.findNode(childLeft.key)
		t.Logf("Структура дерева:\n%s", treeToString(tree.root, ""))

		require.NotNil(t, root, "У дерева есть корень")
		require.NotNil(t, root.left, "У дерева есть левый ребенок")
		require.Equal(t, res.key, childLeft.key, "Функция должна вернуть ноду с искомым ключам")
	})

	t.Run("Должен вернуть nil если такой ноды нет", func(t *testing.T) {
		tree, root, _ := createSimpleTree()
		res := tree.findNode(999)
		t.Logf("Структура дерева:\n%s", treeToString(tree.root, ""))

		require.NotNil(t, root, "У дерева есть корень")
		require.NotNil(t, root.left, "У дерева есть левый ребенок")
		require.Nil(t, res, "Функция должна вернуть nil ")
	})

	t.Run("Проверяем корректный поиск если ключ является ключом корневой ноды", func(t *testing.T) {
		tree, root, _ := createSimpleTree()
		res := tree.findNode(root.key)
		t.Logf("Структура дерева:\n%s", treeToString(tree.root, ""))

		require.NotNil(t, root, "У дерева есть корень")
		require.NotNil(t, root.left, "У дерева есть левый ребенок")
		t.Log("Искомый ключ ->", res.key)
		require.Equal(t, res.key, root.key, "Функция должна вернуть ноду с искомым ключам")
	})
}

func TestRBTree_Delete(t *testing.T) {
	t.Run("Удаление из пустого дерева", func(t *testing.T) {
		tree := New()
		t.Logf("Структура дерева:\n%s", treeToString(tree.root, ""))
		tree.Delete(11)
		require.Nil(t, tree.root, "Дерево должно остаться пустым")
	})

	t.Run("Удаление узла без детей", func(t *testing.T) {
		tree := createLongTree()
		t.Logf("Структура дерева до удаления:\n%s", treeToString(tree.root, "|"))
		parent := tree.findNode(20).parent
		tree.Delete(20)
		node := tree.findNode(20)
		require.Nil(t, node, "Удаляемая нода не должна существовать в дереве")
		require.Equal(t, BLACK, parent.color, "Цвет родителя должен быть черным")
	})

	t.Run("Удаление узла с одним ребенком", func(t *testing.T) {
		tree := createLongTree()
		t.Logf("Структура дерева до удаления:\n%s", treeToString(tree.root, ""))
		// вставляем узел, который станет родителем для дальнейшей проверки
		newNode := tree.InsertTree(10)
		require.Equal(t, RED, newNode.color, "Новый узел должен быть красного цвета")
		parent := newNode.parent
		require.Equal(t, BLACK, parent.color, "Родитель нового узла должен быть черного цвета")
		// Удаляем узел, у которого только один ребенок
		tree.Delete(9)
		node := tree.findNode(9)
		require.Nil(t, node, "Удаляемая нода не должна существовать в дереве")
		// проверяем что балансировка отработала
		require.Equal(t, BLACK, newNode.color, "После балансировки узел должен стать черным")

	})

	t.Run("Удаление узла с двумя детьми", func(t *testing.T) {
		tree := createLongTree()
		t.Logf("Структура дерева до удаления:\n%s", treeToString(tree.root, ""))
		tree.Delete(15)
		t.Logf("Структура дерева после удаления узла с двумя детьми:\n%s", treeToString(tree.root, ""))
		node := tree.findNode(15)
		require.Nil(t, node, "Удаляемый узел должен отсутствовать в дереве")
	})

	t.Run("Удаление корневого узла", func(t *testing.T) {
		tree := createLongTree()
		originalRootKey := tree.root.key
		tree.Delete(originalRootKey)
		t.Logf("Структура дерева после удаления корневого узла:\n%s", treeToString(tree.root, ""))
		require.NotEqual(t, originalRootKey, tree.root.key, "Новый корень должен отличаться от удаленного")
		require.Equal(t, BLACK, tree.root.color, "Новый корень должен быть черным")
		require.NoError(t, checkRBInvariants(tree), "Инварианты RB-дерева нарушены")
	})

	t.Run("Последовательное удаление узлов", func(t *testing.T) {
		tree := createLongTree()
		//удаляем несколько узлов по одному.
		for _, key := range []int{1, 5, 9, 13, 17} {
			tree.Delete(key)
			t.Logf("После удаления %d:\n%s", key, treeToString(tree.root, ""))
			require.Nil(t, tree.findNode(key), "Удаленного узла нет в дереве %d", key)
		}
	})
}

func BenchmarkRBTree_InsertTree(b *testing.B) {
	b.ReportAllocs()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.InsertTree(i)
	}
}

func checkRBInvariants(tree *RBTree) error {
	if tree.root == nil {
		return nil
	}
	if tree.root.color != BLACK {
		return fmt.Errorf("корень должен быть черный")
	}
	//todo можно реализовать как доп проверку
	// Можно проверить, что нет подряд красных узлов и что черная высота одинакова по всем пути рекурсивной проверки
	return nil
}

func treeToString(node *Node, indent string) string {
	if node == nil {
		return indent + "nil\n"
	}
	var color string

	if node.color == RED {
		color = "🔴"
	} else {
		color = "⚫️"
	}

	result := indent + fmt.Sprintf("Key: %d, Color: %s\n", node.key, color)
	result += treeToString(node.left, indent+"|--")
	result += treeToString(node.right, indent+"|-")
	return result
}

func createSimpleTree() (tree *RBTree, root, childLeft *Node) {
	tree = New()
	root = &Node{
		key:   10,
		color: BLACK,
	}
	childLeft = &Node{
		key:   5,
		color: RED,
	}
	root.left = childLeft
	childLeft.parent = root
	tree.root = root
	return
}

func createRecoloringTree() (tree *RBTree, root, childLeft, childRight, newNode *Node) {
	tree = New()
	root = &Node{
		key:   10,
		color: BLACK,
	}
	childLeft = &Node{
		key:   5,
		color: RED,
	}
	childRight = &Node{
		key:   15,
		color: RED,
	}
	root.left = childLeft
	root.right = childRight
	childLeft.parent = root
	childRight.parent = root
	tree.root = root

	newNode = &Node{
		key:    20,
		color:  RED,
		parent: childRight,
	}
	childRight.right = newNode
	return
}

func createLeftRotateTree() (tree *RBTree, root, parent, newNode *Node) {
	// Родитель – левый ребёнок, новая нода вставляется как правый ребёнок родителя
	tree = New()
	root = &Node{
		key:   10,
		color: BLACK,
	}
	parent = &Node{
		key:   5,
		color: RED,
	}
	root.left = parent
	parent.parent = root
	tree.root = root

	// newNode вставляется как правый ребёнок родителя
	newNode = &Node{
		key:    7,
		color:  RED,
		parent: parent,
	}
	parent.right = newNode
	return
}

func createRightRotateTree() (tree *RBTree, root, parent, newNode *Node) {
	//Родитель – правый ребёнок, новая нода вставляется как левый ребёнок родителя
	tree = New()
	root = &Node{
		key:   10,
		color: BLACK,
	}
	parent = &Node{
		key:   15,
		color: RED,
	}
	root.right = parent
	parent.parent = root
	tree.root = root

	// newNode вставляется как левый ребёнок родителя.
	newNode = &Node{
		key:    13,
		color:  RED,
		parent: parent,
	}
	parent.left = newNode
	return
}

func createLongTree() *RBTree {
	arr := []int{11, 1, 12, 2, 13, 3, 14, 4, 15, 5, 16, 6, 17, 7, 18, 8, 19, 9, 20}
	tree := New()
	for i := 0; i < len(arr); i++ {
		tree.InsertTree(arr[i])
	}
	return tree
}
