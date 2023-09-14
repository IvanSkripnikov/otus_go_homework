package main

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	Count        int
	FirstElement *ListItem
	LastElement  *ListItem
}

func (list list) Len() int {
	return list.Count
}

func (list list) Front() *ListItem {
	return list.FirstElement
}

func (list list) Back() *ListItem {
	return list.LastElement
}

func (list *list) PushFront(v interface{}) *ListItem {
	val := ListItem{Value: v}
	if list.Count == 0 {
		val.Prev = nil
		val.Next = nil
		list.FirstElement = &val
		list.LastElement = &val
	} else {
		firstElement := list.FirstElement
		val.Prev = nil
		val.Next = firstElement
		firstElement.Prev = &val
		list.FirstElement = &val
	}
	list.Count++

	return &val
}

func (list *list) PushBack(v interface{}) *ListItem {
	val := ListItem{Value: v}
	if list.Count == 0 {
		val.Prev = nil
		val.Next = nil
		list.FirstElement = &val
		list.LastElement = &val
	} else {
		lastElement := list.LastElement
		val.Prev = lastElement
		val.Next = nil
		lastElement.Next = &val
		list.LastElement = &val
	}
	list.Count++

	return &val
}

func (list *list) Remove(i *ListItem) {
	if list.Count == 0 {
		return
	}
	if list.Count == 1 {
		list.FirstElement = nil
		list.LastElement = nil
	} else {
		if i.Prev == nil {
			nextElement := i.Next
			nextElement.Prev = nil
			list.FirstElement = nextElement
		} else if i.Next == nil {
			prevElement := i.Prev
			prevElement.Next = nil
			list.LastElement = prevElement
		} else {
			nextElement := i.Next
			prevElement := i.Prev
			i.Next.Prev = prevElement
			i.Prev.Next = nextElement
		}
	}
	list.Count--
}

func (list *list) MoveToFront(i *ListItem) {
	// проверка на непустой список
	if list.Count == 0 {
		return
	}
	// проверка на то, что это первый элемент
	if i.Prev == nil {
		return
	}

	if i.Next == nil {
		prevElement := i.Prev
		prevElement.Next = nil
		list.LastElement = prevElement

		list.setCurrentElementFirst(i)
	} else {
		nextElement := i.Next
		prevElement := i.Prev
		i.Next.Prev = prevElement
		i.Prev.Next = nextElement

		list.setCurrentElementFirst(i)
	}
}

func NewList() List {
	return new(list)
}

func (list *list) setCurrentElementFirst(i *ListItem) {
	firstElement := list.FirstElement
	i.Prev = nil
	i.Next = firstElement
	firstElement.Prev = i
	list.FirstElement = i
}
