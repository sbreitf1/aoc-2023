package helper

import (
	"container/heap"

	"golang.org/x/exp/constraints"
)

// inspired by https://pkg.go.dev/container/heap#example-package-PriorityQueue

type PriorityQueue[P constraints.Ordered, T any] struct {
	items priorityQueueItemList[P, T]
}

func MakePriorityQueue[P constraints.Ordered, T any]() *PriorityQueue[P, T] {
	return &PriorityQueue[P, T]{}
}

func (pq *PriorityQueue[P, T]) Push(priority P, obj T) {
	heap.Push(&pq.items, &priorityQueueItem[P, T]{Object: obj, Priority: priority})
}

func (pq *PriorityQueue[P, T]) Pop() (T, P) {
	item := heap.Pop(&pq.items).(*priorityQueueItem[P, T])
	return item.Object, item.Priority
}

func (pq *PriorityQueue[P, T]) Len() int {
	return pq.items.Len()
}

type priorityQueueItem[P constraints.Ordered, T any] struct {
	Object   T
	Priority P
	Index    int
}

type priorityQueueItemList[P constraints.Ordered, T any] []*priorityQueueItem[P, T]

func (pq priorityQueueItemList[P, T]) Len() int { return len(pq) }

func (pq priorityQueueItemList[P, T]) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq priorityQueueItemList[P, T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityQueueItemList[P, T]) Push(x any) {
	n := len(*pq)
	item := x.(*priorityQueueItem[P, T])
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueueItemList[P, T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
