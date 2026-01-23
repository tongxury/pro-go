package data

type File struct {
	Name string
	Body []byte
	URL  string
	Md5  string
}

type Files []File
