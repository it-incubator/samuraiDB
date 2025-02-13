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

		t.Log("Структура дерева после вставки узла 10:", treeToString(tree.root, ""))

		require.Equal(t, BLACK, node10.color, "Корень должен быть черным", node10.color)
	})

	t.Run("Вставка ребенка вправо", func(t *testing.T) {
		node20 := tree.InsertTree(20)

		t.Log("Структура дерева после вставки узла 20:", treeToString(tree.root, ""))

		require.NotNil(t, node20)
		require.Equal(t, 20, node20.key, "Узел должен иметь ключ 20")
		require.Equal(t, RED, node20.color, "Новый узел должен быть красного цвета", node20.color)
		require.Equal(t, tree.root, node20.parent, "Родитель узла с ключом 20 должен быть корнем")
		require.Equal(t, node20, tree.root.right, "Узел с ключом 20 должен быть правым потомком корня")
	})

	t.Run("Вставка ребенка влево", func(t *testing.T) {
		node3 := tree.InsertTree(3)

		t.Log("Структура дерева после вставки узла 3:", treeToString(tree.root, ""))

		require.NotNil(t, node3)
		require.Equal(t, 3, node3.key, "Узел должен иметь ключ 3")
		require.Equal(t, RED, node3.color, "Новый узел должен быть красного цвета", node3.color)
		require.Equal(t, tree.root, node3.parent, "Родитель узла с ключом 3 должен быть корнем")
		require.Equal(t, node3, tree.root.left, "Узел с ключом 3 должен быть левым потомком корня")
	})

	t.Run("Вставляем дополнительные узлы и проверяем балансировку", func(t *testing.T) {
		node30 := tree.InsertTree(30)

		t.Log("Структура дерева после вставки узла 30:", treeToString(tree.root, ""))

		require.NotNil(t, node30)
		require.Equal(t, 30, node30.key)
		require.Equal(t, BLACK, node30.parent.color, "Родитель узла 30 должен быть черным", node30.parent.color)
		require.Equal(t, BLACK, tree.root.color, "Корень должен быть черным", tree.root.color)

		node40 := tree.InsertTree(40)

		t.Log("Структура дерева после вставки узла 40:", treeToString(tree.root, ""))

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
	result := indent + fmt.Sprintf("Key: %d, Color: %s\n", node.key, node.color)
	result += treeToString(node.left, indent+"  ")
	result += treeToString(node.right, indent+"  ")
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
