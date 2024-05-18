package custom

import (
	"container/list"
)


type RoundRobin struct {
	queue              *list.List // use for round robin
}

func NewRoundRobin() *RoundRobin{
	return &RoundRobin{
		queue: list.New(),
	}
}

func (rr *RoundRobin) Add(T any) {
	rr.queue.PushBack(T)
}

func (rr *RoundRobin) RemoveNode(removeNode any) {
	for e := rr.queue.Front(); e != nil; e = e.Next() {
		if e.Value == removeNode {
			rr.queue.Remove(e)
			break
		}
	} 
}

func (rr *RoundRobin) GetCurNode() any{
	node := rr.queue.Front() // iterator
	val := node.Value; // value
	rr.queue.Remove(node)
	rr.queue.PushBack(val)
	return val
}
