package cryptox

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateRandomKey(t *testing.T) {
	key, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}
	if len(key) != KeySize {
		t.Errorf("Generated key has incorrect size: got %d, want %d", len(key), KeySize)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	plaintext := []byte("This is a test message.")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text does not match plaintext: got %s, want %s", string(decrypted), string(plaintext))
	}
}

func createImageFile(t *testing.T, filename string) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	colorRed := color.RGBA{255, 0, 0, 255}
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.SetRGBA(x, y, colorRed)
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
}

func TestIsImageFile(t *testing.T) {
	tempDir := t.TempDir()

	pngFile := filepath.Join(tempDir, "test.png")
	createImageFile(t, pngFile)
	if !isImageFile(pngFile) {
		t.Errorf("isImageFile should return true for PNG file")
	}

	txtFile := filepath.Join(tempDir, "test.txt")
	err := ioutil.WriteFile(txtFile, []byte("not an image"), 0644)
	if err != nil {
		t.Fatalf("failed to create file")
	}
	if isImageFile(txtFile) {
		t.Errorf("isImageFile should return false for a non-image file")
	}
}
