package encryptor

import (
	"bytes"
	"github.com/alexmullins/zip"
	"io"
	"log"
	"os"
	"testing"
)

func TestName1(t *testing.T) {
	s := "1d49c53c684343ef1c09acd40cd7091d63a1b3447e7aa8fe6ec962d8840f2c09"

	en, _ := DESCBC("Ur2mCula").Encrypt(s)
	println(en)

	decrypt, _ := DESCBC("Ur2mCula").Decrypt("68e33e181f757e22d33a2b117cc70840a46c98283b7669a198f9078fe4d15dc0605ef17d9d34c4afe815c5865a9316f4230e6e7ea101dd8ef79e9700590dcc02c7a7b507d5e45489")
	println(decrypt)

}

func TestA(t *testing.T) {
	contents := []byte("Hello World")
	fzip, err := os.Create(`./test.zip`)
	if err != nil {
		log.Fatalln(err)
	}

	zipw := zip.NewWriter(fzip)
	defer zipw.Close()
	w, err := zipw.Encrypt(`test.txt`, `golang`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(w, bytes.NewReader(contents))
	if err != nil {
		log.Fatal(err)
	}
	zipw.Flush()
}
