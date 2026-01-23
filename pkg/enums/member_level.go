package enums

type MemberLevel string

const (
	MemberLevel_Free    MemberLevel = "free"
	MemberLevel_Basic   MemberLevel = "basic"
	MemberLevel_Pro     MemberLevel = "pro"
	MemberLevel_ProPlus MemberLevel = "pro_plus"
)

func (t MemberLevel) Values() []string {
	return []string{
		MemberLevel_Free.String(),
		MemberLevel_Basic.String(),
		MemberLevel_Pro.String(),
		MemberLevel_ProPlus.String()}
}

func (t MemberLevel) Name() string {
	return map[MemberLevel]string{
		MemberLevel_Free:    "Free",
		MemberLevel_Basic:   "Basic",
		MemberLevel_Pro:     "Pro",
		MemberLevel_ProPlus: "Pro+",
	}[t]
}

func (t MemberLevel) Color() string {
	return map[MemberLevel]string{
		MemberLevel_Free:    "",
		MemberLevel_Basic:   "#69b1ff",
		MemberLevel_Pro:     "#003eb3",
		MemberLevel_ProPlus: "#ff4d4f",
	}[t]
}

func MemberLevelByCode(code int) MemberLevel {
	return map[int]MemberLevel{
		1: MemberLevel_Free,
		2: MemberLevel_Basic,
		3: MemberLevel_Pro,
		4: MemberLevel_ProPlus,
	}[code]
}

func (t MemberLevel) String() string {
	return string(t)
}

func (t MemberLevel) Sort() int {
	for i, x := range t.Values() {
		if t.String() == x {
			return i
		}
	}
	return 0
}
