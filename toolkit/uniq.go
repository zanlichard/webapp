package toolkit

import "sort"

func Uniq(data sort.Interface) (size int) {
	p, l := 0, data.Len()
	if l <= 1 {
		return l
	}
	for i := 1; i < l; i++ {
		if !data.Less(p, i) {
			continue
		}
		p++
		if p < i {
			data.Swap(p, i)
		}
	}
	return p + 1
}
