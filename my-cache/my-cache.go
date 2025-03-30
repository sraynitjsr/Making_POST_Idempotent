package mycache

import (
	"container/list"
)

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	head     *list.List
}

type Node struct {
	key   string
	value string
}

const CAPACITY = 5

var myCache = LRUCache{
	capacity: CAPACITY,
	cache:    make(map[string]*list.Element),
	head:     list.New(),
}

func Get(idempotencyKey string) string {
	if element, present := myCache.cache[idempotencyKey]; present {
		myCache.head.MoveToFront(element)
		return element.Value.(*Node).value
	}
	return ""
}

func Put(idempotencyKey, orderID string) {
	if element, present := myCache.cache[idempotencyKey]; present {
		element.Value.(*Node).value = orderID
		myCache.head.MoveToFront(element)
		return
	}

	if myCache.head.Len() >= myCache.capacity {
		oldestElement := myCache.head.Back()
		delete(myCache.cache, oldestElement.Value.(*Node).key)
		myCache.head.Remove(oldestElement)
	}

	newNode := &Node{key: idempotencyKey, value: orderID}
	newElement := myCache.head.PushFront(newNode)
	myCache.cache[idempotencyKey] = newElement
}
