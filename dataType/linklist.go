package dataType

import (
	"fmt"
	//"main/types"
	"log"
	"sync"
)

type Node struct {
	ID    int
	Value int
	Prev  *Node
	Next  *Node
}

type DoublyLinkedList struct {
	Head    *Node
	Tail    *Node
	Size    int
	maxSize int
	nodeMap map[int]*Node // 用于快速查找节点的 map
	mu      sync.RWMutex  // 添加读写锁
}

func NewDoublyLinkedList(maxSize int) *DoublyLinkedList {
	return &DoublyLinkedList{
		maxSize: maxSize,
		nodeMap: make(map[int]*Node),
	}
}
func (list *DoublyLinkedList) Insert(id int, value int) {
	list.mu.Lock() // 加写锁
	defer list.mu.Unlock()
	newNode := &Node{ID: id, Value: value}
	// 查找是否已经存在相同的ID
	existingNode := list.findNodeById(id)
	log.Printf("%d=%d,add at link", newNode.ID, newNode.Value)
	if existingNode != nil {
		if value > existingNode.Value { // 如果新值更大，则更新
			list.removeNode(existingNode) // 删除旧节点
			list.insertNewNode(newNode)   //  插入新节点
		}
		return
	}
	list.insertNewNode(newNode)

}

func (list *DoublyLinkedList) insertNewNode(newNode *Node) {
	if list.Size == list.maxSize && newNode.Value <= list.Tail.Value { // 如果列表满了，并且新节点比尾部的小
		return // 直接返回
	}

	if list.Head == nil { // 链表为空
		list.Head = newNode
		list.Tail = newNode
	} else if newNode.Value >= list.Head.Value { // 如果新节点比头节点大，则插入到头节点
		newNode.Next = list.Head
		list.Head.Prev = newNode
		list.Head = newNode
	} else { // 如果新节点比头节点小，则需要遍历插入
		current := list.Head
		for current.Next != nil && newNode.Value < current.Next.Value { // 寻找插入的位置
			current = current.Next
		}

		newNode.Next = current.Next
		if current.Next != nil { // 如果不是添加到最后
			current.Next.Prev = newNode
		} else {
			list.Tail = newNode // 新节点是最后一个节点
		}

		newNode.Prev = current
		current.Next = newNode
	}
	list.nodeMap[newNode.ID] = newNode // 存储到map
	list.Size++

	if list.Size > list.maxSize { // 如果链表长度超过 20， 则移除最后的节点
		remove := list.removeLast()
		delete(list.nodeMap, remove.ID) // 从 map 中删除

	}
}

func (list *DoublyLinkedList) removeLast() *Node {
	var node = list.Tail
	if list.Size == 0 {
		return node
	}

	if list.Size == 1 {
		list.Head = nil
		list.Tail = nil

	} else {

		list.Tail = list.Tail.Prev
		list.Tail.Next = nil

	}
	list.Size--
	return node
}

func (list *DoublyLinkedList) removeNode(node *Node) {
	if node == list.Head && node == list.Tail { // 只有一个节点
		list.Head = nil
		list.Tail = nil
	} else if node == list.Head { // 删除头节点
		list.Head = node.Next
		list.Head.Prev = nil
	} else if node == list.Tail { // 删除尾节点
		list.Tail = node.Prev
		list.Tail.Next = nil
	} else { // 删除中间的节点
		node.Prev.Next = node.Next
		node.Next.Prev = node.Prev
	}
	delete(list.nodeMap, node.ID) // 从 map 中删除
	list.Size--
}
func (list *DoublyLinkedList) findNodeById(id int) *Node {
	node, ok := list.nodeMap[id]
	if ok {
		return node
	}
	return nil
}

func (list *DoublyLinkedList) GetTop(n int) []string {
	list.mu.RLock()
	defer list.mu.RUnlock()
	var result []string
	current := list.Head
	for current != nil && len(result) < n {
		result = append(result, fmt.Sprintf("%d=%d", current.ID, current.Value))
		current = current.Next
	}
	return result

}
