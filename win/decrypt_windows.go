//go:build windows

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

func decryptSecret(blob string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		raw, err = base64.RawStdEncoding.DecodeString(blob)
		if err != nil {
			return "", err
		}
	}
	if len(raw) > 3 && string(raw[:3]) == "v10" {
		raw = raw[3:]
	}
	if key, err := chromeKey(); err == nil && len(raw) > 12+16 {
		if plain, err := openGCM(key, raw); err == nil {
			return string(plain), nil
		}
	}
	plain, err := dpapi(raw)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func chromeKey() ([]byte, error) {
	path := filepath.Join(os.Getenv("APPDATA"), "Grok Bot", "Local State")
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc struct {
		OSCrypt struct {
			EncryptedKey string `json:"encrypted_key"`
		} `json:"os_crypt"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, err
	}
	raw, err := base64.StdEncoding.DecodeString(doc.OSCrypt.EncryptedKey)
	if err != nil || len(raw) < 6 || string(raw[:5]) != "DPAPI" {
		return nil, errors.New("没有系统加密密钥")
	}
	return dpapi(raw[5:])
}

func openGCM(key, payload []byte) ([]byte, error) {
	if len(key) != 32 || len(payload) < 12+16 {
		return nil, errors.New("密钥长度不对")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce, ct := payload[:gcm.NonceSize()], payload[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ct, nil)
}

func dpapi(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("空数据")
	}
	in := windows.DataBlob{Size: uint32(len(data)), Data: &data[0]}
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, 0, &out); err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return append([]byte(nil), unsafe.Slice(out.Data, out.Size)...), nil
}
