package utilities

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	keyLength := 32
	random, err := GenerateRandomString(keyLength)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, len(random), keyLength)
}

func TestHash(t *testing.T) {
	data := "123456"
	generated := "$argon2id$v=19$m=65536,t=1,p=6$ZYw8nPdCJwLql+boKoEO4w$VvtqLRPkn583kOik8YW04RC0G3b0LOkoBCdteVzOuNs"
	_, err := Hash(data)
	if err != nil {
		t.Error(err)
		return
	}

	verified, err := VerifyHash(data, generated)

	assert.Equal(t, verified, true)
}

func TestBase64EncodingDecoding(t *testing.T) {
	text := "test.123456"
	encoded := EncodeBase64([]byte(text))
	decoded, err := DecodeBase64(encoded)

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, text, string(decoded))
}

func TestBase64URLEncodingDecoding(t *testing.T) {
	text := "sample.url"
	encoded := EncodeBase64URL([]byte(text))
	decoded, err := DecodeBase64URL(encoded)

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, text, string(decoded))
}

func TestAESEncryptDecrypt(t *testing.T) {
	passphrase := "secret-123456"
	salt := "salt-123456"
	message := "my secret message"

	aesEncryptedKey, errEncryptingKey := GenerateAESKey([]byte(passphrase), []byte(salt))
	if errEncryptingKey != nil {
		t.Error(errEncryptingKey)
		return
	}

	encryptedMessage, errEncrypting := EncryptAES(aesEncryptedKey, []byte(message))
	if errEncrypting != nil {
		t.Error(errEncrypting)
		return
	}

	decryptedMessage, errDecrypting := DecryptAES(aesEncryptedKey, encryptedMessage)
	if errDecrypting != nil {
		t.Error(errDecrypting)
		return
	}

	assert.Equal(t, message, string(decryptedMessage))
}
