package enums

type FileCategory string

const (
	FileCategory_Doc   FileCategory = "doc"
	FileCategory_Pdf   FileCategory = "pdf"
	FileCategory_Image FileCategory = "image"
)

func (t FileCategory) Values() []string {
	return []string{
		FileCategory_Doc.String(),
		FileCategory_Pdf.String(),
		FileCategory_Image.String(),
	}
}

func (t FileCategory) String() string {
	return string(t)
}
