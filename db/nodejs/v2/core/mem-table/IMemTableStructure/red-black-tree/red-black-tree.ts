import {IMemTableStructure} from "../i-mem-table-structure";

export enum Color {
    RED = "RED",
    BLACK = "BLACK"
}

export interface Comparable {
    compare(other: this): number;
}

export class TreeNode<TKey, TValue> {
    key: TKey;
    value: TValue;
    color: Color;
    left: TreeNode<TKey, TValue> | null;
    right: TreeNode<TKey, TValue> | null;
    parent: TreeNode<TKey, TValue> | null;

    constructor(key: TKey, value: TValue, color: Color = Color.RED, parent: TreeNode<TKey, TValue> | null = null) {
        this.key = key;
        this.value = value;
        this.color = color;
        this.left = null;
        this.right = null;
        this.parent = parent;
    }
}

export class RedBlackTree<TKey, TValue> implements IMemTableStructure<TKey, TValue> {
    private root: TreeNode<TKey, TValue> | null = null;

    private rotateLeft(node: TreeNode<TKey, TValue>): void {
        const rightChild = node.right;
        if (!rightChild) return;

        node.right = rightChild.left;
        if (rightChild.left) rightChild.left.parent = node;

        rightChild.parent = node.parent;
        if (!node.parent) {
            this.root = rightChild;
        } else if (node === node.parent.left) {
            node.parent.left = rightChild;
        } else {
            node.parent.right = rightChild;
        }

        rightChild.left = node;
        node.parent = rightChild;
    }

    private rotateRight(node: TreeNode<TKey, TValue>): void {
        const leftChild = node.left;
        if (!leftChild) return;

        node.left = leftChild.right;
        if (leftChild.right) leftChild.right.parent = node;

        leftChild.parent = node.parent;
        if (!node.parent) {
            this.root = leftChild;
        } else if (node === node.parent.right) {
            node.parent.right = leftChild;
        } else {
            node.parent.left = leftChild;
        }

        leftChild.right = node;
        node.parent = leftChild;
    }

    private fixInsertion(node: TreeNode<TKey, TValue>): void {
        while (node.parent && node.parent.color === Color.RED) {
            const grandparent = node.parent.parent;
            if (!grandparent) break;

            if (node.parent === grandparent.left) {
                const uncle = grandparent.right;
                if (uncle && uncle.color === Color.RED) {
                    node.parent.color = Color.BLACK;
                    uncle.color = Color.BLACK;
                    grandparent.color = Color.RED;
                    node = grandparent;
                } else {
                    if (node === node.parent.right) {
                        node = node.parent;
                        this.rotateLeft(node);
                    }
                    node.parent!.color = Color.BLACK;
                    grandparent.color = Color.RED;
                    this.rotateRight(grandparent);
                }
            } else {
                const uncle = grandparent.left;
                if (uncle && uncle.color === Color.RED) {
                    node.parent.color = Color.BLACK;
                    uncle.color = Color.BLACK;
                    grandparent.color = Color.RED;
                    node = grandparent;
                } else {
                    if (node === node.parent.left) {
                        node = node.parent;
                        this.rotateRight(node);
                    }
                    node.parent!.color = Color.BLACK;
                    grandparent.color = Color.RED;
                    this.rotateLeft(grandparent);
                }
            }
        }
        this.root!.color = Color.BLACK;
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
            if (key < current.key) {
                current = current.left;
            } else if (key > current.key) {
                current = current.right;
            } else {
                // If key already exists, update value
                current.value = value;
                return;
            }
        }

        const newNode = new TreeNode(key, value, Color.RED, parent);
        if (key < parent!.key) {
            parent!.left = newNode;
        } else {
            parent!.right = newNode;
        }

        this.fixInsertion(newNode);

        this.print()
    }

    delete(key: TKey): void {
        const node = this.findNode(key);
        if (!node) return;

        let y = node;
        let yOriginalColor = y.color;
        let x: TreeNode<TKey, TValue> | null;

        if (!node.left) {
            x = node.right;
            this.transplant(node, node.right);
        } else if (!node.right) {
            x = node.left;
            this.transplant(node, node.left);
        } else {
            y = this.minimum(node.right);
            yOriginalColor = y.color;
            x = y.right;

            if (y.parent === node) {
                if (x) x.parent = y;
            } else {
                this.transplant(y, y.right);
                y.right = node.right;
                if (y.right) y.right.parent = y;
            }

            this.transplant(node, y);
            y.left = node.left;
            if (y.left) y.left.parent = y;
            y.color = node.color;
        }

        if (yOriginalColor === Color.BLACK && x) {
            this.fixDeletion(x);
        }
    }

    private transplant(u: TreeNode<TKey, TValue>, v: TreeNode<TKey, TValue> | null): void {
        if (!u.parent) {
            this.root = v;
        } else if (u === u.parent.left) {
            u.parent.left = v;
        } else {
            u.parent.right = v;
        }
        if (v) v.parent = u.parent;
    }

    private fixDeletion(x: TreeNode<TKey, TValue>): void {
        while (x !== this.root && x.color === Color.BLACK) {
            if (x === x.parent!.left) {
                let w = x.parent!.right!;
                if (w.color === Color.RED) {
                    w.color = Color.BLACK;
                    x.parent!.color = Color.RED;
                    this.rotateLeft(x.parent!);
                    w = x.parent!.right!;
                }
                if ((!w.left || w.left.color === Color.BLACK) &&
                    (!w.right || w.right.color === Color.BLACK)) {
                    w.color = Color.RED;
                    x = x.parent!;
                } else {
                    if (!w.right || w.right.color === Color.BLACK) {
                        if (w.left) w.left.color = Color.BLACK;
                        w.color = Color.RED;
                        this.rotateRight(w);
                        w = x.parent!.right!;
                    }
                    w.color = x.parent!.color;
                    x.parent!.color = Color.BLACK;
                    if (w.right) w.right.color = Color.BLACK;
                    this.rotateLeft(x.parent!);
                    x = this.root!;
                }
            } else {
                let w = x.parent!.left!;
                if (w.color === Color.RED) {
                    w.color = Color.BLACK;
                    x.parent!.color = Color.RED;
                    this.rotateRight(x.parent!);
                    w = x.parent!.left!;
                }
                if ((!w.right || w.right.color === Color.BLACK) &&
                    (!w.left || w.left.color === Color.BLACK)) {
                    w.color = Color.RED;
                    x = x.parent!;
                } else {
                    if (!w.left || w.left.color === Color.BLACK) {
                        if (w.right) w.right.color = Color.BLACK;
                        w.color = Color.RED;
                        this.rotateLeft(w);
                        w = x.parent!.left!;
                    }
                    w.color = x.parent!.color;
                    x.parent!.color = Color.BLACK;
                    if (w.left) w.left.color = Color.BLACK;
                    this.rotateRight(x.parent!);
                    x = this.root!;
                }
            }
        }
        x.color = Color.BLACK;
    }

    private minimum(node: TreeNode<TKey, TValue>): TreeNode<TKey, TValue> {
        let current = node;
        while (current.left) current = current.left;
        return current;
    }

    public findNode(key: TKey): TreeNode<TKey, TValue> | null {
        let current = this.root;
        while (current) {
            if (key < current.key) {
                current = current.left;
            } else if (key > current.key) {
                current = current.right;
            } else {
                return current;
            }
        }
        return null;
    }

    find(key: TKey): TValue | null {
        const node = this.findNode(key);
        return node ? node.value : null;
    }

    print(): void {
        if (!this.root) {
            console.log("Empty tree");
            return;
        }

        const printNode = (node: TreeNode<TKey, TValue>, prefix: string, isLeft: boolean): void => {
            if (!node) return;

            console.log(
                prefix +
                (isLeft ? "├── " : "└── ") +
                `${node.key}(${node.color})`
            );

            // Рекурсивный вызов для левого и правого поддерева
            if (node.left) {
                printNode(node.left, prefix + (isLeft ? "│   " : "    "), true);
            }
            if (node.right) {
                printNode(node.right, prefix + (isLeft ? "│   " : "    "), false);
            }
        };

        // Начинаем с корня
        console.log(`${this.root.key}(${this.root.color})`);
        if (this.root.left) {
            printNode(this.root.left, "", true);
        }
        if (this.root.right) {
            printNode(this.root.right, "", false);
        }
    }
}
