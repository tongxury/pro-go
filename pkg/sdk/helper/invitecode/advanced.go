package invitecode

var defaultGenerator, _ = NewGenerator(6)

func IdToCode(id int64) string {
	encode, _ := defaultGenerator.Encode(uint64(id))
	return encode
}

func CodeToId(code string) int64 {
	return int64(defaultGenerator.Decode(code))
}
