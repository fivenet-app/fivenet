package utils

func SliceUint64ToInt64(in []uint64) (out []int64) {
	for _, i := range in {
		out = append(out, int64(i))
	}
	return
}

func SliceInt32ToInt(in []int32) (out []int) {
	for _, i := range in {
		out = append(out, int(i))
	}
	return
}
