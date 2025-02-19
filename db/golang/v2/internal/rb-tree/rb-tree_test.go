package rb_tree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRBTree_InsertTree(t *testing.T) {
	tree := New[int, string]()

	t.Run("Создаем корень", func(t *testing.T) {
		node10 := tree.InsertTree(10, "data 10")

		t.Log("Структура дерева после вставки узла 10: \n")
		tree.Print()

		require.Equal(t, BLACK, node10.color, "Корень должен быть черным", node10.color)
	})

	t.Run("Вставка ребенка вправо", func(t *testing.T) {
		node20 := tree.InsertTree(20, "data 20")

		t.Log("Структура дерева после вставки узла 20: \n")
		tree.Print()

		require.NotNil(t, node20)
		require.Equal(t, 20, node20.key, "Узел должен иметь ключ 20")
		require.Equal(t, RED, node20.color, "Новый узел должен быть красного цвета", node20.color)
		require.Equal(t, tree.root, node20.parent, "Родитель узла с ключом 20 должен быть корнем")
		require.Equal(t, node20, tree.root.right, "Узел с ключом 20 должен быть правым потомком корня")
	})

	t.Run("Вставка ребенка влево", func(t *testing.T) {
		node3 := tree.InsertTree(3, "data-3")

		t.Log("Структура дерева после вставки узла 3: \n")
		tree.Print()

		require.NotNil(t, node3)
		require.Equal(t, 3, node3.key, "Узел должен иметь ключ 3")
		require.Equal(t, RED, node3.color, "Новый узел должен быть красного цвета", node3.color)
		require.Equal(t, tree.root, node3.parent, "Родитель узла с ключом 3 должен быть корнем")
		require.Equal(t, node3, tree.root.left, "Узел с ключом 3 должен быть левым потомком корня")
	})

	t.Run("Вставляем дополнительные узлы и проверяем балансировку", func(t *testing.T) {
		node30 := tree.InsertTree(30, "data 10")

		t.Log("Структура дерева после вставки узла 30: \n")
		tree.Print()

		require.NotNil(t, node30)
		require.Equal(t, 30, node30.key)
		require.Equal(t, BLACK, node30.parent.color, "Родитель узла 30 должен быть черным", node30.parent.color)
		require.Equal(t, BLACK, tree.root.color, "Корень должен быть черным", tree.root.color)

		node40 := tree.InsertTree(40, "data 10")

		t.Log("Структура дерева после вставки узла 40: \n")
		tree.Print()

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
		t.Log("Структура дерева после балансировкой: \n")
		tree.Print()

		require.Equal(t, BLACK, root.color, "Корень должен оставаться черным")
		require.Equal(t, RED, childLeft.color, "Дочерний узел должен оставаться красного цвета")
		require.Equal(t, root, childLeft.parent, "Родитель дочернего узла не должен измениться")
		require.Equal(t, childLeft, root.left, "Дочерний узел должен быть левым потомком корня")
		require.Nil(t, root.right, "Правый узел должен быть nil")
	})

	t.Run("Должны перекраситься ноды (родитель и дядя красные)", func(t *testing.T) {
		tree, root, childLeft, childRight, newNode := createRecoloringTree()

		// На этом этапе родитель (childRight) и дядя (childLeft) красные
		t.Log("Структура дерева до балансировки: \n")
		tree.Print()
		tree.fixInsert(newNode)
		t.Log("Структура дерева после балансировкой: \n")
		tree.Print()

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
		t.Log("Структура дерева после балансировкой: \n")
		tree.Print()

		tree.fixInsert(newNode)
		t.Log("Структура дерева после балансировкой: \n")
		tree.Print()
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
		t.Log("Структура дерева перед балансировкой: \n")
		tree.Print()

		tree.fixInsert(newNode)
		t.Log("Структура дерева после балансировкой: \n")
		tree.Print()

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
		tree := New[int, string]()
		res := tree.findNode(1)

		require.Nil(t, res, "Дерево пустое")
	})

	t.Run("Должен найти ноду по ключу", func(t *testing.T) {
		tree, root, childLeft := createSimpleTree()
		res := tree.findNode(childLeft.key)
		t.Log("Структура дерева \n")
		tree.Print()

		require.NotNil(t, root, "У дерева есть корень")
		require.NotNil(t, root.left, "У дерева есть левый ребенок")
		require.Equal(t, res.key, childLeft.key, "Функция должна вернуть ноду с искомым ключам")
	})

	t.Run("Должен вернуть nil если такой ноды нет", func(t *testing.T) {
		tree, root, _ := createSimpleTree()
		res := tree.findNode(999)
		t.Log("Структура дерева \n")
		tree.Print()

		require.NotNil(t, root, "У дерева есть корень")
		require.NotNil(t, root.left, "У дерева есть левый ребенок")
		require.Nil(t, res, "Функция должна вернуть nil ")
	})

	t.Run("Проверяем корректный поиск если ключ является ключом корневой ноды", func(t *testing.T) {
		tree, root, _ := createSimpleTree()
		res := tree.findNode(root.key)
		t.Log("Структура дерева \n")
		tree.Print()

		require.NotNil(t, root, "У дерева есть корень")
		require.NotNil(t, root.left, "У дерева есть левый ребенок")
		t.Log("Искомый ключ ->", res.key)
		require.Equal(t, res.key, root.key, "Функция должна вернуть ноду с искомым ключам")
	})
}

func TestRBTree_Delete(t *testing.T) {
	t.Run("Удаление из пустого дерева", func(t *testing.T) {
		tree := New[int, string]()
		t.Log("Структура дерева \n")
		tree.Print()
		tree.Delete(11)
		require.Nil(t, tree.root, "Дерево должно остаться пустым")
	})

	t.Run("Удаление узла без детей", func(t *testing.T) {
		tree := createLongTree()
		t.Log("Структура дерева до удаления\n")
		tree.Print()
		parent := tree.findNode(20).parent
		tree.Delete(20)
		node := tree.findNode(20)
		require.Nil(t, node, "Удаляемая нода не должна существовать в дереве")
		require.Equal(t, BLACK, parent.color, "Цвет родителя должен быть черным")
		require.NoError(t, checkRBInvariants(tree), "Инварианты RB-дерева нарушены после удаления %d")

	})

	t.Run("Удаление узла с одним ребенком", func(t *testing.T) {
		tree := createLongTree()
		// вставляем узел, который станет родителем для дальнейшей проверки
		newNode := tree.InsertTree(10, "data-10")
		require.Equal(t, RED, newNode.color, "Новый узел должен быть красного цвета")
		parent := newNode.parent
		require.Equal(t, BLACK, parent.color, "Родитель нового узла должен быть черного цвета")
		t.Log("Структура дерева до удаления\n")
		tree.Print()
		// Удаляем узел, у которого только один ребенок
		tree.Delete(9)
		node := tree.findNode(9)
		t.Log("Структура дерева после удаления\n")
		tree.Print()
		require.Nil(t, node, "Удаляемая нода не должна существовать в дереве")
		// проверяем что балансировка отработала
		require.Equal(t, RED, newNode.color, "После балансировки остаться красным")
		require.Equal(t, BLACK, newNode.parent.color, "Родитель узла должен стать черным")
		require.Equal(t, RED, newNode.parent.left.color, "После балансировки остаться красным")
		require.NoError(t, checkRBInvariants(tree), "Инварианты RB-дерева нарушены после удаления")

	})

	t.Run("Удаление узла с двумя детьми", func(t *testing.T) {
		tree := createLongTree()
		t.Log("Структура дерева до удаления\n")
		tree.Print()
		tree.Delete(15)
		t.Log("Структура дерева после удаления\n")
		tree.Print()
		node := tree.findNode(15)
		require.Nil(t, node, "Удаляемый узел должен отсутствовать в дереве")
		require.NoError(t, checkRBInvariants(tree), "Инварианты RB-дерева нарушены после удаления")

	})

	t.Run("Удаление корневого узла", func(t *testing.T) {
		tree := createLongTree()
		originalRootKey := tree.root.key
		tree.Delete(originalRootKey)
		t.Log("Структура дерева после удаления корневого узла: \n")
		tree.Print()
		require.NotEqual(t, originalRootKey, tree.root.key, "Новый корень должен отличаться от удаленного")
		require.Equal(t, BLACK, tree.root.color, "Новый корень должен быть черным")
		require.NoError(t, checkRBInvariants(tree), "Инварианты RB-дерева нарушены")
	})

	t.Run("Последовательное удаление узлов", func(t *testing.T) {
		tree := createLongTree()
		//удаляем несколько узлов по одному.
		for _, key := range []int{1, 5, 9, 13, 17} {
			tree.Delete(key)
			t.Log("После удаления \n")
			tree.Print()
			require.Nil(t, tree.findNode(key), "Удаленного узла нет в дереве %d", key)
			require.NoError(t, checkRBInvariants(tree), "Инварианты RB-дерева нарушены после удаления %d", key)
		}
	})
}

func BenchmarkRBTree_InsertTree(b *testing.B) {
	b.ReportAllocs()
	tree := New[int, string]()
	for i := 0; i < b.N; i++ {
		tree.InsertTree(i, fmt.Sprintf("data-%d", i))
	}
}

func BenchmarkRBTree_Find(b *testing.B) {
	tree := New[int, string]()
	for i := 0; i < 10000; i++ {
		tree.InsertTree(i, fmt.Sprintf("data %d", i))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		// Проводим поиск случайного ключа
		_ = tree.findNode(i % 10000)
	}
}

// В тесте при измерениях аллоцируется на каждую итерацию новое дерево,
// поэтому цифра 10001 allocs/op это нормально
func BenchmarkRBTree_Delete_Only(b *testing.B) {
	// Предварительно создаем дерево с 10000 узлов
	tree := New[int, string]()
	for j := 0; j < 10000; j++ {
		tree.InsertTree(j, fmt.Sprintf("data %d", j))
	}
	keys := make([]int, 10000)
	for i := 0; i < 10000; i++ {
		keys[i] = i
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		localTree := cloneTree(tree)
		for _, key := range keys {
			localTree.Delete(key)
		}
	}
}

// Более чистый бенчмарк для удаления без аллокаций на каждой итерации
func BenchmarkRBTree_Delete_Sequential(b *testing.B) {
	numNodes := 10000
	tree := New[int, string]()
	for j := 0; j < numNodes; j++ {
		tree.InsertTree(j, fmt.Sprintf("data %d", j))
	}

	keys := make([]int, numNodes-1)
	for i := 0; i < numNodes-1; i++ {
		keys[i] = i
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if len(keys) == 0 {
			break
		}

		key := keys[0]
		keys = keys[1:]
		tree.Delete(key)
	}
}

// Более чистый бенчмарк для удаления, но в обратном порядке, без аллокаций на каждой итерации
func BenchmarkRBTree_Delete_Reverse(b *testing.B) {
	numNodes := 10000
	tree := New[int, string]()
	for j := 0; j < numNodes; j++ {
		tree.InsertTree(j, fmt.Sprintf("data %d", j))
	}

	keys := make([]int, numNodes-1)
	for i := 0; i < numNodes-1; i++ {
		keys[i] = numNodes - 1 - i
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if len(keys) == 0 {
			break
		}

		key := keys[0]
		keys = keys[1:] // Убираем ключ из списка
		tree.Delete(key)
	}
}
