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

func (lruCache *lruCache) Set(key Key, value interface{}) bool {
	item, ok := lruCache.items[key]

	if ok {
		if vaultItem, vaultOk := item.Value.(*Vault); vaultOk {
			vaultItem.Value = value
			lruCache.queue.MoveToFront(item)

			return true
		}
	} else {
		vault := Vault{Key: key, Value: value}
		lruCache.items[key] = lruCache.queue.PushFront(&vault)
		if lruCache.queue.Len() > lruCache.capacity {
			RemoveOverflowElements(lruCache)
		}
	}

	return false
}

func (lruCache *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := lruCache.items[key]
	if !ok {
		return nil, false
	}
	lruCache.queue.MoveToFront(item)
	if vaultItem, vaultOk := item.Value.(*Vault); vaultOk {
		return vaultItem.Value, ok
	}

	return nil, false
}

func (lruCache *lruCache) Clear() {
	lruCache.queue = NewList()
	lruCache.items = make(map[Key]*ListItem, lruCache.capacity)
}

func RemoveOverflowElements(lruCache *lruCache) {
	lastItem := lruCache.queue.Back()
	if vaultItem, vaultOk := lastItem.Value.(*Vault); vaultOk {
		delete(lruCache.items, vaultItem.Key)
		lruCache.queue.Remove(lastItem)
	}
}
