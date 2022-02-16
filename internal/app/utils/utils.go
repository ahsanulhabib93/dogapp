package utils

const (
	DEFAULT_PER_PAGE = uint64(20)
)

func Int64Min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
