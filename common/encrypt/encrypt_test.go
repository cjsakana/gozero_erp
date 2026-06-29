package encrypt

import (
	"testing"
)

func TestEncryptMobile(t *testing.T) {
	mobile := "13800138000"
	encryptedMobile, err := EncMobile(mobile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(encryptedMobile)
	decryptedMobile, err := DecMobile(encryptedMobile)
	if err != nil {
		t.Fatal(err)
	}
	if mobile != decryptedMobile {
		t.Fatalf("expected %s, but got %s", mobile, decryptedMobile)
	}
	t.Log(decryptedMobile)
}

func TestEncryptPassword(t *testing.T) {
	password := "123456"
	encryptedPassword, err := HashPassword(password)
	t.Log(encryptedPassword)
	if err != nil {
		t.Fatal(err)
	}
	ok := CheckPassword(password, encryptedPassword)
	if !ok {
		t.Fatal(ok)
	} else {
		t.Log(ok)
	}
}
func TestE(t *testing.T) {
	password := "123456"
	encryptedPassword := "$2a$10$nf0UEOnlNS2gYtdxz3XtEu3pyNuBAfNDGIH7qkVMUH0gFwgH2rkcq"
	ok := CheckPassword(password, encryptedPassword)
	if !ok {
		t.Fatal(ok)
	} else {
		t.Log(ok)
	}
}

func TestEncryptIDCard(t *testing.T) {
	idCard := "110105199001011234"
	encryptIDCard, err := EncIDCard(idCard)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(encryptIDCard)
	oriCard, err := DecIDCard(encryptIDCard)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(oriCard)
}

func TestEncryptAccount(t *testing.T) {
	account := "6222025836914728"
	encryptAccount, err := EncAccount(account)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(encryptAccount)
	oriCard, err := DecAccount(encryptAccount)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(oriCard)
}
