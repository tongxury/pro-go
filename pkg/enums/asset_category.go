package enums

type AssetCategory string

const (
	AssetCategory_File AssetCategory = "file"
)

func (t AssetCategory) Values() []string {
	return []string{
		AssetCategory_File.String(),
	}
}

func (t AssetCategory) String() string {
	return string(t)
}
