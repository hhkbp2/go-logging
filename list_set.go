package logging

import (
	"container/list"
)

type ListSet struct {
	*list.List
}

func NewListSet() *ListSet {
	return &ListSet{
		List: list.New(),
	}
}

func (self *ListSet) SetAdd(i interface{}) {
	self.List.PushBack(i)
}

func (self *ListSet) SetRemove(i interface{}) bool {
	for e := self.List.Front(); e != nil; e = e.Next() {
		if e.Value == i {
			self.List.Remove(e)
			return true
		}
	}
	return false
}

func (self *ListSet) SetContains(i interface{}) bool {
	for e := self.List.Front(); e != nil; e = e.Next() {
		if e.Value == i {
			return true
		}
	}
	return false
}

func (self *ListSet) SetClone() *ListSet {
	newList := list.New()
	for e := self.List.Front(); e != nil; e = e.Next() {
		newList.PushBack(e.Value)
	}
	return &ListSet{
		List: newList,
	}
}
