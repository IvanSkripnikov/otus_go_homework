package main

import "fmt"

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
	List
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

func (list list) PushFront(v interface{}) *ListItem {
	val := ListItem{Value: v}
	list.FirstElement = &val
	list.Count++
	if list.Count == 0 {
		list.LastElement = &val
	}

	return &val
}

func (list list) PushBack(v interface{}) *ListItem {
	val := ListItem{Value: v}
	list.LastElement = &val
	list.Count++
	if list.Count == 0 {
		list.FirstElement = &val
	}

	return &val
}

func (list list) Remove(i *ListItem) {
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
				nextElement.Prev = prevElement
			}
			list.Count--
			break
		} else {
			currentElement = currentElement.Next
		}
	}
}

func (list list) MoveToFront(i *ListItem) {
	currentElement := list.FirstElement
	for index := 0; index < list.Count; index++ {
		if i == currentElement {
			if index == 0 {

			} else if index == list.Count-1 {
				prevElement := currentElement.Prev
				prevElement.Next = nil
				list.LastElement = prevElement
			} else {
				nextElement := currentElement.Next
				prevElement := currentElement.Prev
				nextElement.Prev = prevElement
			}
			break
		} else {
			currentElement = currentElement.Next
		}
	}

	firstElement := list.FirstElement
	i.Prev = nil
	i.Next = firstElement
	list.FirstElement = i
}

func NewList() List {
	return new(list)
}

func main() {
	fmt.Println("working")
}
