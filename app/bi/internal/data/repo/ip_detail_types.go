package repo

type IpDetail struct {
	IpRangeStart       string  `ch:"ip_range_start"`
	IpRangeEnd         string  `ch:"ip_range_end"`
	IpRangeStartNumber int64   `ch:"ip_range_start_num"`
	IpRangeEndNumber   int64   `ch:"ip_range_end_num"`
	CountryCode        string  `ch:"country_code"`
	State1             string  `ch:"state1"`
	State2             string  `ch:"state2"`
	City               string  `ch:"city"`
	Postcode           string  `ch:"postcode"`
	Latitude           float64 `ch:"latitude"`
	Longitude          float64 `ch:"longitude"`
	Timezone           string  `ch:"timezone"`
	//Verify      int8      `ck:"verify"`
}

type IpDetails []IpDetail
