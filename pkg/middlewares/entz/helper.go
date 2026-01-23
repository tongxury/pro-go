package entz

func VerifyPageSize(page, size int64) (int, int, int, bool) {
	if size == 0 {
		return 0, 0, 0, false
	}

	if size > 500 {
		size = 500
	}

	if page < 0 {
		page = 1
	}

	return int((page - 1) * size), int(size), int(page), true

}
