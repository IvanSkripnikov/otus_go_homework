package main

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (elem *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := elem.items[key]; !ok {
		elem.items[key] = elem.queue.PushFront(value)

		if elem.queue.Len() > elem.capacity {
			lastItem := elem.queue.Back()
			val := lastItem.Value
			for k, v := range elem.items {
				if v.Value == val {
					delete(elem.items, k)
					elem.queue.Remove(lastItem)
				}
			}
		}
	} else {
		item.Value = value
		elem.queue.MoveToFront(item)
		return true
	}

	return false
}

func (elem *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := elem.items[key]; ok {
		elem.queue.MoveToFront(item)
		return item.Value, ok
	}

	return nil, false
}

func (elem *lruCache) Clear() {
	elem.queue = NewList()
	elem.items = make(map[Key]*ListItem, elem.capacity)
}
