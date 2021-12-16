package repository

type Pagination struct {
	Offset int64
	Limit  int64
}

func Paginate(pagination Pagination, sliceLength int64) (int64, int64) {
	start := pagination.Offset

	if start > sliceLength {
		start = sliceLength
	}

	end := start + pagination.Limit
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
