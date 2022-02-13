package utils

func Int64Min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
