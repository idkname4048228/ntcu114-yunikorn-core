package heap

import (
	"container/heap"
)

// 定義一個座標及其附加信息的結構體
type CoordinateInfo struct {
	Coordinates []float64 // 任意維度的座標
	Value       float64   // 要比較的值，堆會根據這個值來排序
	Index       int       // 在堆中的索引
}

// 定義一個座標最小堆的類型
type CoordinateHeap []*CoordinateInfo

// Len 代表堆的長度
func (h CoordinateHeap) Len() int { return len(h) }

// Less 用於比較兩個座標的信息
func (h CoordinateHeap) Less(i, j int) bool {
	return h[i].Value < h[j].Value
}

// Swap 用於交換堆中的兩個元素
func (h CoordinateHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

// Push 用於向堆中插入元素
func (h *CoordinateHeap) Push(x interface{}) {
	n := len(*h)
	coord := x.(*CoordinateInfo)
	coord.Index = n
	*h = append(*h, coord)
}

// Pop 用於從堆中移除最小的元素
func (h *CoordinateHeap) Pop() interface{} {
	old := *h
	n := len(old)
	coord := old[n-1]
	old[n-1] = nil  // 避免內存泄漏
	coord.Index = -1 // 安全地移除
	*h = old[0 : n-1]
	return coord
}

// 更新座標上的資訊
func (h *CoordinateHeap) Update(coord *CoordinateInfo, newValue float64) {
	coord.Value = newValue
	heap.Fix(h, coord.Index)
}

// 查找座標是否存在
func (h CoordinateHeap) Find(coords []float64) *CoordinateInfo {
	for _, coord := range h {
		equal := true
		for i := range coords {
			if coord.Coordinates[i] != coords[i] {
				equal = false
				break
			}
		}
		if equal {
			return coord
		}
	}
	return nil
}