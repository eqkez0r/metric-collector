package encrypting

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testData = "asdasdasdasdasdafdsgsdfljhsdt;lkfgjserl;kgjtl"

func TestEncrypting(t *testing.T) {
	path := os.Getenv("CRYPTO_KEY")
	t.Log("path", path)
	prkey, err := GetPrivateKey(path)
	if err != nil {
		t.Fatal(err)
	}
	pukey, err := GetPublicKey(path)
	if err != nil {
		t.Fatal(err)
	}

	b := prkey.PublicKey.Equal(pukey)
	if !b {
		t.Fatal("invalid public key")
	}

	encryptedData, err := Encrypt(pukey, []byte(testData))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("encryptedData", encryptedData)

	decryptedData, err := Decrypt(prkey, encryptedData)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(decryptedData), testData)
}
