package wkv_a

import (
	"errors"
	"sync"
)

var Kvt *Boot

type Boot struct {
	root *node
	lock sync.RWMutex
}

type node struct {
	value *list
	left  *node
	right *node
}

type list struct {
	index uint64
	key   string
	value string
	next  *list
}

func (b *Boot) Set(key string, content string) error {

	b.lock.Lock()
	defer b.lock.Unlock()

	index, err := tranUint(key)
	if err != nil {
		return err
	}

	newNode := &node{&list{index: index, value: content, key: key}, nil, nil}
	if b.root == nil {
		b.root = newNode
	} else {
		insert(b.root, newNode)
	}
	return nil
}

func insert(root, newNode *node) {

	if newNode.value.index < root.value.index {

		if root.left == nil {
			root.left = newNode
		} else {
			insert(root.left, newNode)
		}

	} else if newNode.value.index > root.value.index {

		if root.right == nil {
			root.right = newNode
		} else {
			insert(root.right, newNode)
		}

	} else {
		matchList(root, newNode)
	}
}

func matchList(node *node, newNode *node) {
	if node.value.key == newNode.value.key {
		node.value.value = newNode.value.value
	} else {
		if node.value.next == nil {
			node.value.next = newNode.value
		} else {
			cycleSetList(node.value, newNode.value)
		}
	}
}

func cycleSetList(fat *list, chil *list) {
	if fat.next == nil {
		fat.next = chil
		return
	}
	if fat.key == chil.key {
		fat.value = chil.value
		return
	}
	cycleSetList(fat.next, chil)
}

func (b *Boot) Get(key string) (string, error) {

	b.lock.RLock()
	defer b.lock.RUnlock()

	index, err := tranUint(key)
	if err != nil {
		return "", err
	}

	val, exist := search(b.root, index, key)
	if !exist {
		return "", errors.New("not fount")
	}

	return val, nil
}

func search(root *node, index uint64, key string) (string, bool) {
	if root == nil {
		return "", false
	}
	if index < root.value.index {
		return search(root.left, index, key)
	} else if index > root.value.index {
		return search(root.right, index, key)
	}
	return searchFormList(root.value, key)
}

func searchFormList(l *list, key string) (string, bool) {
	if l.next == nil && l.key != key {
		return "", false
	}
	if l.key == key {
		return l.value, true
	}
	return searchFormList(l.next, key)
}

func (b *Boot) Del(key string) error {

	index, err := tranUint(key)
	if err != nil {
		return err
	}
	_, res := remove(b.root, index, key)
	if !res {
		return errors.New("key not exist")
	}
	return nil
}

func remove(node *node, index uint64, key string) (*node, bool) {
	if node == nil {
		return nil, true
	}
	var exist bool
	if index < node.value.index {
		node.left, exist = remove(node.left, index, key)
		return node, exist
	}
	if index > node.value.index {
		node.right, exist = remove(node.right, index, key)
		return node, exist
	}

	exist = true

	var clearNode bool
	clearNode = removeFromList(node.value, node.value, key, clearNode)

	if clearNode {
		if node.left == nil && node.right == nil {
			node = nil
			return node, exist
		}

		if node.left == nil {
			node = node.right
			return node, exist
		}
		if node.right == nil {
			node = node.left
			return node, exist
		}

		rightMin, _ := min(node.right)
		node.value = rightMin
		node.right, _ = remove(node.right, rightMin.index, key)
	}

	return node, exist
}

func removeFromList(fat, now *list, key string, left bool) bool {
	if now.next != nil {
		left = true
	}
	if now == nil {
		return left
	}
	if now.key == key {
		fat.next = now.next
		return left
	}
	fat = now
	now = now.next
	return removeFromList(fat, now, key, left)
}

func min(node *node) (*list, bool) {
	if node == nil {
		return nil, false
	}
	n := node
	for {
		if n.left == nil {
			return n.value, true
		}
		n = n.left
	}
}
