package pkg

import "sort"

// 思路
// 维护一个左端点小于等于当前 q 的区间集合
// 并且这个集合用区间长度排序
// 这个集合可以用最小堆表示 最小堆的最上层是长度最小的值
// 如果右端点比 q 小 那么这个区间也不可能容得下 下一个 q 可以直接删除
// 依次重复知道查找到结果，找到以后不要从最小堆删除

type interval []int

func (it interval) len() int {
	return it[1] - it[0] + 1
}

type heap struct {
	intervals []interval
	size      int
}

// add 往堆里面加一个
// 从下往上调整
func (i *heap) add(it interval) {
	i.intervals[i.size] = it
	i.size++

	cur := i.size - 1
	swapped := true

	for cur != 0 && swapped {
		swapped = i.adjustParent(cur-1/2) != -1
		cur = (cur - 1) / 2
	}
}

func (i *heap) adjustParent(p int) (next int) {
	next = p
	if p*2+1 < i.size && i.intervals[p*2+1].len() < i.intervals[p].len() {
		next = p*2 + 1
	}
	if p*2+2 < i.size && i.intervals[p*2+2].len() < i.intervals[p].len() {
		next = p*2 + 2
	}

	if next == p {
		return -1
	}

	tmp := i.intervals[p]
	i.intervals[p] = i.intervals[next]
	i.intervals[next] = tmp
	return next
}

func (i *heap) top() *interval {
	if i.size == 0 {
		return nil
	}
	return &i.intervals[0]
}

func (i *heap) alloc(size int) {
	i.intervals = make([]interval, size)
}

func (i *heap) rm() {
	if i.size == 0 {
		return
	}

	if i.size == 1 {
		i.size--
		return
	}

	i.intervals[0] = i.intervals[i.size-1]
	i.size--

	cur := 0

	for cur >= 0 {
		cur = i.adjustParent(cur)
	}
}

type query struct {
	q   int
	idx int
	res int
}

func minInterval(intervals [][]int, queries []int) []int {
	var heap heap

	heap.alloc(len(intervals))
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][1]-intervals[i][0] < intervals[j][1]-intervals[j][0]
	})

	qs := make([]query, len(queries))
	for i := range queries {
		qs[i].q = queries[i]
		qs[i].idx = i
	}

	sort.Slice(qs, func(i, j int) bool {
		return qs[i].q < qs[j].q
	})

	i := 0

	for j := 0; j < len(qs); j++ {
		for {
			if intervals[i][0] > qs[j].q {
				break
			}
			if intervals[i][1] <= qs[j].q {
				heap.add(intervals[i])
			}
			i++
		}

		for heap.size > 0 && (*heap.top())[1] < qs[j].q {
			heap.rm()
		}

		if heap.size == 0 {
			qs[j].res = -1
		} else {
			qs[j].res = heap.top().len()
		}
	}

	res := make([]int, len(queries))

	for j := range qs {
		res[qs[j].idx] = qs[j].res
	}
	return res
}
