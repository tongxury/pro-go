package cryptoz

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRsa(t *testing.T) {

	pub := `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwiy6LzMP4bhqbM/j/O4p
LeoJhyHedUFI1QmcU1lWr+8H1DlnLN+vtizrZKdQFlh+A7ipYChreK5BNfuQQxma
CIctYDErlGufG0lXuGa2Q0TXKMG7jeFmwNky6w06w/uuSQIoz3Mv0q39C4jjCST4
rkYcpePQ6JvQnjCx9lAlpl4k0IaQHAwYRXBr83TqkKpPKZJdmMixai0NdRN7C6Ph
iq6taBpFSNQQT7LQuJPX/aoDb30ufp/JJ4zIV+pg50Ub/z2azmpDvhVqa1S6fPMm
nZi07jaz8wmLO9Mt0BzdeG0eDYJBW4l3GPpO21B6bk0I/eTsFErCuIsDykcJZXSj
KQIDAQAB
-----END PUBLIC KEY-----
	`

	pri := `
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDCLLovMw/huGps
z+P87ikt6gmHId51QUjVCZxTWVav7wfUOWcs36+2LOtkp1AWWH4DuKlgKGt4rkE1
+5BDGZoIhy1gMSuUa58bSVe4ZrZDRNcowbuN4WbA2TLrDTrD+65JAijPcy/Srf0L
iOMJJPiuRhyl49Dom9CeMLH2UCWmXiTQhpAcDBhFcGvzdOqQqk8pkl2YyLFqLQ11
E3sLo+GKrq1oGkVI1BBPstC4k9f9qgNvfS5+n8knjMhX6mDnRRv/PZrOakO+FWpr
VLp88yadmLTuNrPzCYs70y3QHN14bR4NgkFbiXcY+k7bUHpuTQj95OwUSsK4iwPK
RwlldKMpAgMBAAECggEAIHwDGepc+dI0W8fbyHC/iuLgfS75XHxzhtB4yqjji8Nd
d++yhxtU9hFFwC0NhO+BBXZbP68Da8kbN8DCPbeGwW579N/E/quSjqoSdtMYIuDd
bgAbNH1FB3ZOwmwQLMFqQuSNl0cZ9REOiGN6OAlrYRxxpn8acA/BMvXKj/6Qjehf
e3g+2eAYwoX3ztZuiIs35bihdv38GghQdVRwB+NWbtxzhEfV/c70o5YAjE0nyoLc
ohEogoWB29XEyu+NELrMTvl+z7UPh+Z3PzeNIG+rAJM63oI8XLY3kMoja8aEZwwv
bh0NOaLZvrCvzUmbNy08jpU9pzR2YzPSta4MuTcPUQKBgQDjNIMKYs2jGE8tBozh
r4Qbk6aJbJGlQYW/1796BDuBjczHhcdhl/6Zt0NJJpj2A1f9V+Y9NIQZjI10sO4E
/+ExNSCIKGSxxy/IZdPOhE04XrLbBokFALe22/yZhnD2quGbIuXcwIt2hyVfLpCh
DKw6Z3JM7RJjKoQJ8qdB6ouhmwKBgQDayJH9yDZFS8ssUhiRotAUXoZTYz0LyOmO
fJ5Q4Bo9wNq+jApcUhuMBG8POL8bVMJ+RZehKacGVuuDFM/m9v2fqEMLJhJKNqPm
0MM5rXID6HQzVpuXP5d/OEqqkJuNkpkW5XZ9SvMW2v93hU2dB2WBohSQ2mYZs7U1
lLwFswvsiwKBgQCYyUxBpLWSMpuztI7yiVv2S3EXQsoibhBqNMRPYh89/MQzfBPa
3iJY7jMyMuFztkXqWLy8dd9LawgI6530ALpHo+lPhpJINqE8SrWHT9K50HzH6voj
QhtIvWB9QToftkPmVi5rJ5PhTfpkqmSZ2HLNB5mGf3n487M9GU8+dWIWdwKBgFLi
QUyfmM/P0vzLbTtfLu1IkiLtKadZSgIM+/0vqUFT2ortis9G2+DDnT9rBBtalQQ5
YSRRH1GrhDV4oPqi/5qIqD2FAtDSum0rEYq8RsFsQvlgCjnWgZJUxRSxC/0qWIzw
CV+WEVnLRZUGD006DB58RMZLtvpttmzCGCkgl5fZAoGAJqCTv9w1hxFU7SztPSl+
CPHXY97bmN7v81V4+1URNnimLFhla02igW6TIEDZEpsPBZtbhwz3MyjXes4RJCGb
sVQvquH0+U4b2jRp3tBdzLxRXYKpNougvIl3oxwQMo4m/pgvTBDZzL9zKWeinufY
M0ZrbKtiitLIIgcfc+h8yxg=
-----END PRIVATE KEY-----
	`

	txt := "text"

	cipherText, err1 := RSA().Encrypt([]byte(txt), pub)

	plainText, err2 := RSA().Decrypt(cipherText, pri)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, txt, string(plainText))

	hex, err3 := RSA().EncryptString(txt, pub)
	fmt.Println(hex)
	assert.Nil(t, err3)
	decrypt, err4 := RSA().DecryptString(hex, pri)
	assert.Nil(t, err4)
	assert.Equal(t, decrypt, string(plainText))

}
