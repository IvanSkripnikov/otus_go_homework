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

type Vault struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (elem *lruCache) Set(key Key, value interface{}) bool {
	item, ok := elem.items[key]

	if ok {
		if vaultItem, vaultOk := item.Value.(*Vault); vaultOk {
			vaultItem.Value = value
			elem.queue.MoveToFront(item)

			return true
		}
	} else {
		elem.items[key] = elem.queue.PushFront(&Vault{Key: key, Value: value})

		if elem.queue.Len() > elem.capacity {
			lastItem := elem.queue.Back()
			if vaultItem, vaultOk := lastItem.Value.(*Vault); vaultOk {
				delete(elem.items, vaultItem.Key)
				elem.queue.Remove(lastItem)
			}
		}
	}

	return false
}

func (elem *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := elem.items[key]
	if !ok {
		return nil, false
	}

	elem.queue.MoveToFront(item)
	if vaultItem, vaultOk := item.Value.(*Vault); vaultOk {
		return vaultItem.Value, ok
	}

	return nil, false
}

func (elem *lruCache) Clear() {
	elem.queue = NewList()
	elem.items = make(map[Key]*ListItem, elem.capacity)
}
