package toolkit

type UInt64Slice []uint64

func (a UInt64Slice) Len() int {
	return len(a)
}

func (a UInt64Slice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a UInt64Slice) Less(i, j int) bool {
	return a[i] < a[j]
}
