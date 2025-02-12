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
