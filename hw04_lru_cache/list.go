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
	// обрабатываем случай пустого списка
	if list.Count == 0 {
		return
	}

	// обрабатываем случай одного элемента в списке
	if list.Count == 1 {
		list.FirstElement = nil
		list.LastElement = nil
		list.Count--
		return
	}

	// обрабатываем случай удаления первого элемента
	if i.Prev == nil {
		nextElement := i.Next
		nextElement.Prev = nil
		list.FirstElement = nextElement
		list.Count--
		return
	}

	// обрабатываем случай удаления последнего элемента
	if i.Next == nil {
		prevElement := i.Prev
		prevElement.Next = nil
		list.LastElement = prevElement
		list.Count--
		return
	}

	list.LinkNeighboringForElement(i)
	list.Count--
}

func (list *list) MoveToFront(i *ListItem) {
	// проверка на непустой список
	if list.Count == 0 {
		return
	}

	// обрабатываем случай первого элемента
	if i.Prev == nil {
		return
	}

	// обрабатываем случай последнего элемента
	if i.Next == nil {
		prevElement := i.Prev
		prevElement.Next = nil
		list.LastElement = prevElement

		list.SetCurrentElementFirst(i)
	} else {
		list.LinkNeighboringForElement(i)
		list.SetCurrentElementFirst(i)
	}
}

func NewList() List {
	return new(list)
}

func (list *list) SetCurrentElementFirst(i *ListItem) {
	firstElement := list.FirstElement
	i.Prev = nil
	i.Next = firstElement
	firstElement.Prev = i
	list.FirstElement = i
}

func (list *list) LinkNeighboringForElement(i *ListItem) {
	nextElement := i.Next
	prevElement := i.Prev
	i.Next.Prev = prevElement
	i.Prev.Next = nextElement
}
