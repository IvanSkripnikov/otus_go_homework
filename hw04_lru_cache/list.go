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
		list.FirstElement = &val
		list.LastElement = &val
	} else {
		firstElement := list.FirstElement
		firstElement.Prev = &val
		val.Next = firstElement
		list.FirstElement = &val
	}
	list.Count++

	return &val
}

func (list *list) PushBack(v interface{}) *ListItem {
	val := ListItem{Value: v}
	if list.Count == 0 {
		list.FirstElement = &val
		list.LastElement = &val
	} else {
		lastElement := list.LastElement
		lastElement.Next = &val
		val.Prev = lastElement
		list.LastElement = &val
	}
	list.Count++

	return &val
}

func (list *list) Remove(i *ListItem) {
	currentElement := list.FirstElement
	for index := 0; index < list.Count; index++ {
		if i == currentElement {
			if index == 0 {
				nextElement := currentElement.Next
				nextElement.Prev = nil
				list.FirstElement = nextElement
			} else if index == list.Count-1 {
				prevElement := currentElement.Prev
				prevElement.Next = nil
				list.LastElement = prevElement
			} else {
				nextElement := currentElement.Next
				prevElement := currentElement.Prev
				i.Next.Prev = prevElement
				i.Prev.Next = nextElement
			}
			list.Count--
			break
		} else {
			currentElement = currentElement.Next
		}
	}
}

func (list *list) MoveToFront(i *ListItem) {
	currentElement := list.FirstElement
	for index := 0; index < list.Count; index++ {
		if i == currentElement {
			if index == 0 {
				break
			} else if index == list.Count-1 {
				prevElement := i.Prev
				prevElement.Next = nil
				list.LastElement = prevElement

				firstElement := list.FirstElement
				firstElement.Prev = i
				i.Prev = nil
				i.Next = firstElement
				list.FirstElement = i
			} else {
				nextElement := i.Next
				prevElement := i.Prev
				i.Prev = prevElement
				i.Next = nextElement

				firstElement := list.FirstElement
				firstElement.Prev = i
				i.Prev = nil
				i.Next = firstElement
				list.FirstElement = i
			}
			break
		} else {
			currentElement = currentElement.Next
		}
	}
}

func NewList() List {
	return new(list)
}
