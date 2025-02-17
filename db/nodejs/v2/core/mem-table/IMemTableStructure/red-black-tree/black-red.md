# Красно-черное дерево (Red-Black Tree)

Реализация сбалансированного двоичного дерева поиска, которое автоматически поддерживает баланс при вставке и удалении элементов.

## Свойства красно-черного дерева

1. Каждый узел либо красный, либо черный
2. Корень всегда черный
3. Все листья (NULL) считаются черными
4. У красного узла оба потомка черные
5. Для каждого узла все простые пути от него до листьев содержат одинаковое число черных узлов

## API

### Конструктор

```typescript
constructor()
```

Создает пустое красно-черное дерево.

### Методы

#### insert(key: TKey, value: TValue): void

Вставляет новый элемент в дерево или обновляет значение существующего ключа.

```typescript
const tree = new RedBlackTree<number, string>();
tree.insert(1, "one");
```

#### delete(key: TKey): void

Удаляет элемент с указанным ключом из дерева.

```typescript
tree.delete(1);
```

#### find(key: TKey): TValue | null

Поиск значения по ключу. Возвращает null, если ключ не найден.

```typescript
const value = tree.find(1);
```

#### print(): void

Выводит визуальное представление дерева в консоль.

```typescript
tree.print();
```

Пример вывода:
```
10(BLACK)
├── 5(RED)
│   ├── 3(BLACK)
│   └── 7(BLACK)
└── 15(RED)
    ├── 13(BLACK)
    └── 17(BLACK)
```

Особенности реализации:

1. Если дерево пустое, выводится сообщение "Empty tree"
2. Для каждого узла выводится его ключ и цвет в скобках
3. Используются специальные символы для отображения связей между узлами:
    - `├──` для левого потомка
    - `└──` для правого потомка
    - `│   ` для вертикальных линий
4. Метод рекурсивно обходит дерево, поддерживая правильные отступы и связи

## Структура узла

```typescript
class TreeNode<TKey, TValue> {
    key: TKey;
    value: TValue;
    color: Color;
    left: TreeNode<TKey, TValue> | null;
    right: TreeNode<TKey, TValue> | null;
    parent: TreeNode<TKey, TValue> | null;
}
```

## Временная сложность

| Операция | В среднем | В худшем случае |
|----------|-----------|-----------------|
| Вставка  | O(log n)  | O(log n)       |
| Удаление | O(log n)  | O(log n)       |
| Поиск    | O(log n)  | O(log n)       |

## Пример использования

```typescript
const tree = new RedBlackTree<number, string>();

// Вставка элементов
tree.insert(10, "десять");
tree.insert(5, "пять");
tree.insert(15, "пятнадцать");

// Поиск элемента
const value = tree.find(5); // "пять"

// Визуализация дерева
tree.print();

// Удаление элемента
tree.delete(5);
```

## Примечания по реализации

- Дерево поддерживает обобщенные типы для ключей и значений
- Автоматически балансируется после вставки и удаления
- Реализует интерфейс IMemTableStructure
- Все операции гарантированно выполняются за O(log n) времени
- Дополнительное пространство O(n)

## Ограничения

- Ключи должны быть сравнимыми (поддерживать операторы <, >, =)
- Не поддерживает дубликаты ключей (при вставке существующего ключа значение обновляется)

## Использование памяти

Каждый узел дерева содержит:
- Ключ и значение
- Цвет (1 бит)
- Три ссылки (на родителя и двух потомков)

Общее использование памяти: O(n), где n - количество элементов в дереве.





Использование UUID в качестве ключа в красно-черном дереве возможно, но имеет свои особенности и потенциальные проблемы. Давайте рассмотрим это подробнее:

Текущая реализация использует сравнение ключей через операторы `<` и `>`:

```typescript
if (key < current.key) {
    current = current.left;
} else if (key > current.key) {
    current = current.right;
}
```

Это не будет корректно работать с UUID, так как прямое сравнение строк UUID не даст правильного результата для построения дерева.

Вот как можно исправить это:

1. Сначала определим интерфейс для сравнения:

```typescript
export interface Comparable<T> {
    compareTo(other: T): number;
}

export class UUID implements Comparable<UUID> {
    private value: string;

    constructor(value: string) {
        this.value = value;
    }

    compareTo(other: UUID): number {
        return this.value.localeCompare(other.value);
    }

    toString(): string {
        return this.value;
    }
}
```

2. Изменим класс RedBlackTree:

```typescript
export class RedBlackTree<TKey extends Comparable<TKey>, TValue> implements IMemTableStructure<TKey, TValue> {
    // ... остальной код ...

    private compare(a: TKey, b: TKey): number {
        return a.compareTo(b);
    }

    insert(key: TKey, value: TValue): void {
        if (!this.root) {
            this.root = new TreeNode(key, value, Color.BLACK);
            return;
        }

        let parent: TreeNode<TKey, TValue> | null = null;
        let current: TreeNode<TKey, TValue> | null = this.root;

        while (current) {
            parent = current;
            const compareResult = this.compare(key, current.key);
            if (compareResult < 0) {
                current = current.left;
            } else if (compareResult > 0) {
                current = current.right;
            } else {
                current.value = value;
                return;
            }
        }

        // ... остальной код метода ...
    }

    // Аналогично изменить все методы, где используется сравнение ключей
}
```

3. Пример использования:

```typescript
const tree = new RedBlackTree<UUID, string>();

tree.insert(new UUID("550e8400-e29b-41d4-a716-446655440000"), "value1");
tree.insert(new UUID("550e8400-e29b-41d4-a716-446655440001"), "value2");
```

4. Обновим документацию:

```markdown
## Работа с UUID

При использовании UUID в качестве ключей, необходимо обеспечить корректное сравнение:

```typescript
const tree = new RedBlackTree<UUID, string>();

// Правильно
tree.insert(new UUID("550e8400-e29b-41d4-a716-446655440000"), "value1");

// Неправильно
tree.insert("550e8400-e29b-41d4-a716-446655440000", "value1");
```



Рекомендации:
1. Если производительность критична, лучше использовать числовые идентификаторы
2. Если необходимо использовать UUID, реализуйте эффективное сравнение
3. Рассмотрите возможность использования упорядоченных UUID (ULID, TSID) для лучшей производительности дерева

В любом случае, важно добавить корректное сравнение ключей через интерфейс Comparable, чтобы гарантировать правильную работу дерева независимо от типа ключа.