package repository

// minHeapNode is an object for a min-heap node.
type minHeapNode struct {
	idx        int // min-heap index
	dist       int // distance from src
	stationIdx int // station index
}

type minHeap []*minHeapNode

func (pq minHeap) Len() int {
	return len(pq)
}

func (pq minHeap) Less(i, j int) bool {
	return pq[i].dist < pq[j].dist
}

func (pq minHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idx = i
	pq[j].idx = j
}

func (pq *minHeap) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

func (pq *minHeap) Push(x interface{}) {
	item := x.(*minHeapNode)
	*pq = append(*pq, item)
}
