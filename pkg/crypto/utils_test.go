package crypto

import (
	"testing"
)

const pubRSAKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7ioZMaThxys+nATwjs5m
7CoBdnfV0mz4TXMV6xoajeKPG0embTO6oVyIboSQNiy/yszBr/5VKoLaoGYhQ56x
DFzc/dc1m1LngdEZhHQJ2lQihYPtEMS+Kr2LgL20vnWMXrEzE9sgeV1CXx/fCLw4
gkqbAkUQim9eUAvDI4GEA2GD2M5GlLdx0UYX7BjvckO0N9g+MD72WmE64HWm+3Es
k8zMUMDeX2MvtiG78ltKSj3RYDyjwjI8E2bETvJ9Dy4jSXo9clxUvLQxVOf/UZMk
JAR8t5MC4kHFYHwJcjWRdrPltmfjHu7TK6jjU4Am3vnRvBr37KQQ2KMXFeN10unh
pQIDAQAB
-----END PUBLIC KEY-----
`

const privRSAKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEA7ioZMaThxys+nATwjs5m7CoBdnfV0mz4TXMV6xoajeKPG0em
bTO6oVyIboSQNiy/yszBr/5VKoLaoGYhQ56xDFzc/dc1m1LngdEZhHQJ2lQihYPt
EMS+Kr2LgL20vnWMXrEzE9sgeV1CXx/fCLw4gkqbAkUQim9eUAvDI4GEA2GD2M5G
lLdx0UYX7BjvckO0N9g+MD72WmE64HWm+3Esk8zMUMDeX2MvtiG78ltKSj3RYDyj
wjI8E2bETvJ9Dy4jSXo9clxUvLQxVOf/UZMkJAR8t5MC4kHFYHwJcjWRdrPltmfj
Hu7TK6jjU4Am3vnRvBr37KQQ2KMXFeN10unhpQIDAQABAoIBACQFSfowLdWpvLZs
KNXwcbCWSdJZHYXN5WARX+dG820yLuK5W3p4sGlnTVspwYXwDrHldgXgOZFMaTSJ
Pc60WaK9CM97lSgAyfLgZTObOUJEJ1R6N4ipuPlN4aN/Da8gqDJKKqd+JNM2P6uT
bnArx4AtOHSHbZECdwk0PjdIh2bbNhmOpTKKbHQxGndsNywI256SN+IpzhzVOPMP
isB6jrqgXbFHzhN7kBSMdq3vS28flzHKjyPabgztFMy8Rur7jW5YWUmxCsqKdy29
vFji4qH9mgHt1ZzEFTwlKt4+Bq1IODQZQfPoXpujHaDEMqE78yLGi1xJFdn1o0op
J88axgECgYEA+32zT9GWXUnQAhRjGPIJJYVcW46ZavXyQR1CyHIVDPjGlM5rJy9i
+sWERL6FIQRbLfwL5DFXGKpuK5bsNbtWZi7nS2evmtKDAHtCSpbsfE5bq/z/m/Nj
2XBT8dUoZ/P0xVTfCA6JkmK5nIpLfQ+SbNPpblzbOYgtvpiqxQh8rBECgYEA8m87
OoCboXqB5a/d9UfSOWDs8gO/KHUVXE/knhI162dkohTkTHzSCfXerHKzafxAR53/
jJJ36ZLVbUg3XAIqKy6sXnatBhB8xMiT4gObaEbhU1LbhhC1X26LgEfmz4vdxVuB
QtwPpeyilqWexwJ5q5p2Lsh9Xsb3lhELmND4wFUCgYAEDWJ2RspFdosDfZCbNksv
b5atYv7V2mCs6+vHjw8Hxnpsq2bOmtTddZFMCkXa4lcVxpnqc2ET5KshyKzFsN8T
hm3zqRgLRpkVyOaojQYCesC/ZLQ5rxJMzqKLowOjqSqog1WUq6dL1ItpGlFdEoMp
fcClJpnhs3AJQix+QETCAQKBgBuDwEcdfYxQKRn2YcyKwDM+6uV0w2dGEoyNjLbb
/j6fV26FzHtZ10TGIOWVhwNKW8lFB1He9bkOryZeAdpxbHPGMk3uTijYCjETSqVm
H2cwVDZuuvd2Qf94vmBqyKlZiGvzvLHn4+bC+pj6ZxDTGRf+ydb5bjEph8QCXzyS
ywiZAoGAWbOBt6wLOT9jLEufbL7qejFcpzaJPMoLJUE53q/UA7emI2YPnl8EORIt
MkEnDUhdmbrETQy1Kgzl8wRylckjiLUiHr/lIbkXw2y2VFQ5pQKPCoA9nttRQ68p
bDGcp73uUkxlalr7w6eYT+rWDVKQPc8m70ow0OzCP2Lir1u28zk=
-----END RSA PRIVATE KEY-----
`

func TestEncryptMessage(t *testing.T) {
	message := "This is a test message to encrypt.."
	byteMsg := []byte(message)
	cypherKey, cypherContent, err := EncryptMessage(byteMsg, []byte(pubRSAKey))
	if err != nil {
		t.Error(err.Error())
	}
	if cypherKey == nil || cypherContent == nil {
		t.Error("error to encrypt message")
	}

}

func TestDecryptMessage(t *testing.T) {
	message := "This is a test message to encrypt.."
	byteMsg := []byte(message)
	cypherKey, cypherContent, err := EncryptMessage(byteMsg, []byte(pubRSAKey))

	decodedByteMsg, err := DecryptMessage(cypherContent, cypherKey, []byte(privRSAKey))
	if err != nil {
		t.Error(err.Error())
	}
	if len(byteMsg) != len(decodedByteMsg) {
		t.Error("error to decrypt message")
	}
	for i, v := range byteMsg {
		if v != decodedByteMsg[i] {
			t.Error("error to decrypt message")
		}
	}
}
