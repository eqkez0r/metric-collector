package encrypting

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"os"
	"path"
)

const (
	privateKeyName = "private.pem"
	publicKeyName  = "public.pem"
	KeyLengthBits  = 4096
)

func GenerateIfNotExist(keysPath string) error {
	privateKeyPath := path.Join(keysPath, privateKeyName)
	publicKeyPath := path.Join(keysPath, publicKeyName)

	privateKeyFile, err := os.OpenFile(privateKeyPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	publicKeyFile, err := os.OpenFile(publicKeyPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	privateKey, err := rsa.GenerateKey(rand.Reader, KeyLengthBits)
	if err != nil {
		return err
	}

	var privateKeyPEM, publicKeyPEM bytes.Buffer

	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return err
	}

	err = pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	if err != nil {
		return err
	}

	_, err = privateKeyFile.Write(privateKeyPEM.Bytes())
	if err != nil {
		return err
	}

	_, err = publicKeyFile.Write(publicKeyPEM.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func GetPrivateKey(keysPath string) (*rsa.PrivateKey, error) {
	privateKeyPath := path.Join(keysPath, privateKeyName)

	privateKeyFile, err := os.Open(privateKeyPath)
	if err != nil {
		return nil, err
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := io.ReadAll(privateKeyFile)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyBytes)

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GetPublicKey(keysPath string) (*rsa.PublicKey, error) {
	publicKeyPath := path.Join(keysPath, publicKeyName)

	publicKeyFile, err := os.Open(publicKeyPath)
	if err != nil {
		return nil, err
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := io.ReadAll(publicKeyFile)
	if err != nil {
		return nil, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyBytes)

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func Encrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	var encryptedData bytes.Buffer

	chunkSize := key.Size()

	chunks := len(data) / chunkSize

	for i := 0; i <= chunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(data) {
			end = len(data)
		}
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, key, data[start:end])
		if err != nil {
			return nil, err
		}
		_, err = encryptedData.Write(encryptedChunk)
		if err != nil {
			return nil, err
		}
	}

	return encryptedData.Bytes(), nil
}

func Decrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	chunkSize := key.Size()

	var decryptedData bytes.Buffer

	chunks := len(data) / chunkSize
	log.Println(len(data), chunkSize, chunks)
	for i := 0; i < chunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(data) {
			end = len(data)
		}
		decryptedChunk, err := rsa.DecryptPKCS1v15(rand.Reader, key, data[start:end])
		if err != nil {
			return nil, err
		}
		_, err = decryptedData.Write(decryptedChunk)
		if err != nil {
			return nil, err
		}
	}

	return decryptedData.Bytes(), nil
}
