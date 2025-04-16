package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	gookitcolor "github.com/gookit/color" // Renamed to avoid conflict
	"github.com/urfave/cli/v2"
)

// Constants
const (
	Version  = "v1.0.0"           // Updated Version
	Author   = "Amul Thantharate" // Tool Author
	KeySize  = 32                 // AES-256 key size (32 bytes)
	AsciiArt = `
   __    _   _ ____  ____
  / /   / \ | |  _ \|  _ \
 / /   / _ \| | |_) | |_) |
/ /___/ ___ \ |  __/|  __/
\____/_/   \_\_|   |_|
 Image Encryption Tool
`
	EncryptedExtension = ".enc"
	StegoMessageLimit  = 1000 // Maximum message length for steganography
)

// Helper Functions

// GenerateRandomKey generates a random AES key.
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, KeySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// Encrypt encrypts data using AES-256 GCM.
func Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to create nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-256 GCM.
func Decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open GCM: %w", err)
	}

	return plaintext, nil
}

// LoadImage loads an image from a file.
func LoadImage(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

// SaveImage saves an image to a file.  Supports PNG and JPEG.
func SaveImage(filename string, img image.Image, outputFormat string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer f.Close()

	switch strings.ToLower(outputFormat) {
	case "jpg", "jpeg":
		opt := &jpeg.Options{Quality: 90} // Adjust quality as needed (0-100)
		err = jpeg.Encode(f, img, opt)
		if err != nil {
			return fmt.Errorf("failed to encode image to JPEG: %w", err)
		}
	default: // Default to PNG
		err = png.Encode(f, img)
		if err != nil {
			return fmt.Errorf("failed to encode image to PNG: %w", err)
		}
	}
	return nil
}

// SaveImage saves an image to a file with default format PNG.  Supports PNG and JPEG.
func SaveImageDefault(filename string, img image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("failed to encode image to PNG: %w", err)
	}

	return nil
}

// ImageToBytes converts an image to a byte slice.
func ImageToBytes(img image.Image) ([]byte, error) {
	// Encode the image to PNG in memory
	buf := new(bytes.Buffer) // Import "bytes"
	err := png.Encode(buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image to bytes: %w", err)
	}

	return buf.Bytes(), nil
}

// BytesToImage converts a byte slice to an image.
func BytesToImage(data []byte) (image.Image, error) {
	r := bytes.NewReader(data) // Import "bytes"
	img, err := png.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bytes to image: %w", err)
	}
	return img, nil
}

func isImageFile(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		return false // Or log the error
	}
	defer f.Close()

	_, format, err := image.DecodeConfig(f)
	if err != nil {
		return false // Or log the error
	}

	// List of supported formats
	supportedFormats := []string{"jpeg", "jpg", "png", "gif", "bmp", "tiff"}
	for _, supportedFormat := range supportedFormats {
		if strings.ToLower(format) == supportedFormat {
			return true
		}
	}
	return false
}

// CLI Commands

// encryptCmd encrypts an image or a directory of images.
var encryptCmd = &cli.Command{
	Name:    "encrypt",
	Aliases: []string{"e"},
	Usage:   "Encrypt an image or a directory of images",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Value:    "",
			Usage:    "Input image file or directory",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "encrypted_output", // Default output directory/file prefix
			Usage:   "Output encrypted image file or directory",
		},
		&cli.StringFlag{
			Name:    "key",
			Aliases: []string{"k"},
			Value:   "",
			Usage:   "Encryption key (base64 encoded). If not provided, a new key will be generated and printed/saved.",
		},
		&cli.StringFlag{
			Name:  "keyfile",
			Usage: "File to save the generated key to (if no key is provided). If specified, the key will be saved here after generation.",
			Value: "", // Default empty string = do not save to file.
		},
		&cli.BoolFlag{
			Name:  "print-key",
			Usage: "Print the generated key (if no key is provided). If set to false (default), the key will only be printed to stdout once and NOT stored anywhere. If set to true, the key will be printed even if a key is provided.",
			Value: false,
		},
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "Recursively search subdirectories for images to encrypt.",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite existing files in the output directory without warning.",
			Value: false,
		},
	},
	Action: func(c *cli.Context) error {
		inputPath := c.String("input")
		outputPath := c.String("output")
		keyBase64 := c.String("key")
		keyFile := c.String("keyfile")
		printKey := c.Bool("print-key")
		recursive := c.Bool("recursive")
		overwrite := c.Bool("overwrite")

		// Get key
		var key []byte
		var err error

		// Check environment variable first
		if keyBase64 == "" {
			keyBase64 = os.Getenv("IMAGE_ENCRYPTION_KEY")
			if keyBase64 != "" {
				gookitcolor.Yellow.Println("Using key from environment variable IMAGE_ENCRYPTION_KEY")
			}
		}

		if keyBase64 == "" {
			// Generate a new key
			key, err = GenerateRandomKey()
			if err != nil {
				gookitcolor.Red.Println(fmt.Errorf("failed to generate key: %w", err))
				return err
			}

			keyBase64Encoded := base64.StdEncoding.EncodeToString(key)

			if keyFile != "" {
				// Save the key to a file
				err = ioutil.WriteFile(keyFile, []byte(keyBase64Encoded), 0600) // Permissions 0600: read/write for owner only
				if err != nil {
					gookitcolor.Red.Println(fmt.Errorf("failed to save key to file: %w", err))
					return err
				}
				gookitcolor.Green.Println("Generated Key (base64 encoded):", keyBase64Encoded)
				gookitcolor.Green.Println("Key saved to file:", keyFile)

			} else {
				if printKey {
					gookitcolor.Green.Println("Generated Key (base64 encoded):", keyBase64Encoded)
				} else {
					gookitcolor.Green.Println("Generated Key (base64 encoded):", keyBase64Encoded)
					gookitcolor.Yellow.Println("IMPORTANT: This key is only displayed once. Do NOT lose it! Save it somewhere secure.")
				}
			}

		} else {
			// Decode the key from base64
			key, err = base64.StdEncoding.DecodeString(keyBase64)
			if err != nil {
				gookitcolor.Red.Println(fmt.Errorf("failed to decode key: %w", err))
				return err
			}
			if len(key) != KeySize {
				gookitcolor.Red.Println("invalid key size: key must be %d bytes when base64 decoded", KeySize)
				return fmt.Errorf("invalid key size: key must be %d bytes when base64 decoded", KeySize)
			}
			if printKey {
				gookitcolor.Green.Println("Using provided Key (base64 encoded):", base64.StdEncoding.EncodeToString(key))
			}
		}

		// Check if the input is a file or a directory
		fileInfo, err := os.Stat(inputPath)
		if err != nil {
			log.Printf("failed to stat input path: %v", err) // Use log for errors early
			return err
		}

		if fileInfo.IsDir() {
			// Process directory
			return encryptDirectory(inputPath, outputPath, key, recursive, overwrite)
		} else {
			// Process single file
			return encryptFile(inputPath, outputPath, key, overwrite)
		}
	},
}

func encryptFile(inputFilename, outputFilename string, key []byte, overwrite bool) error {
	// Check if the output file exists and if overwriting is allowed
	if _, err := os.Stat(outputFilename); err == nil && !overwrite {
		// File exists and overwrite is not allowed
		gookitcolor.Yellow.Printf("Output file %s already exists.  Overwrite with --overwrite flag.\n", outputFilename)
		return nil
	}

	// Load image
	img, err := LoadImage(inputFilename)
	if err != nil {
		log.Printf("failed to load image: %v", err) // Use log for errors
		return err
	}

	// Convert image to bytes
	imgBytes, err := ImageToBytes(img)
	if err != nil {
		log.Printf("failed to convert image to bytes: %v", err) // Use log for errors
		return err
	}

	// Encrypt the image bytes
	ciphertext, err := Encrypt(key, imgBytes)
	if err != nil {
		log.Printf("failed to encrypt: %v", err) // Use log for errors
		return err
	}

	// Save the encrypted data to a new file
	err = os.MkdirAll(filepath.Dir(outputFilename), os.ModeDir|0755) // Ensure output directory exists
	if err != nil {
		log.Printf("failed to create output directory: %v", err) // Use log for errors
		return err
	}

	err = ioutil.WriteFile(outputFilename, ciphertext, 0644)
	if err != nil {
		log.Printf("failed to write encrypted data to file: %v", err) // Use log for errors
		return err
	}

	gookitcolor.Cyan.Println("Image encrypted and saved to:", outputFilename)
	return nil
}

func encryptDirectory(inputDir, outputDir string, key []byte, recursive bool, overwrite bool) error {
	var wg sync.WaitGroup
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate the error
		}

		if info.IsDir() && path != inputDir && !recursive {
			return filepath.SkipDir // Skip subdirectories if not recursive
		}

		if !info.IsDir() { // Only check files
			if isImageFile(path) { // Use the file path
				// Construct the output filename
				relPath, err := filepath.Rel(inputDir, path)
				if err != nil {
					log.Printf("failed to get relative path: %v", err)
					return err
				}

				outputFilename := filepath.Join(outputDir, relPath+EncryptedExtension) // Append .enc extension

				wg.Add(1)
				go func(p, o string) {
					defer wg.Done()
					err := encryptFile(p, o, key, overwrite)
					if err != nil {
						log.Printf("Error encrypting %s: %v\n", p, err)
					}
				}(path, outputFilename) // Encrypt each image file
			}
		}
		return nil
	})
	wg.Wait() // Wait for all goroutines to complete

	if err != nil {
		log.Printf("error walking the path %s: %v", inputDir, err)
		return err
	}

	return nil
}

// decryptCmd decrypts an image.
var decryptCmd = &cli.Command{
	Name:    "decrypt",
	Aliases: []string{"d"},
	Usage:   "Decrypt an image or a directory of images",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i"},
			Value:    "",
			Usage:    "Input encrypted image file or directory",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "decrypted_output",
			Usage:   "Output decrypted image file or directory",
		},
		&cli.StringFlag{
			Name:     "key",
			Aliases:  []string{"k"},
			Value:    "",
			Usage:    "Encryption key (base64 encoded)",
			Required: true, // Key is now required for decryption
		},
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "Recursively search subdirectories for encrypted images to decrypt.",
			Value:   false,
		},
		&cli.StringFlag{
			Name:  "encrypted-ext",
			Value: EncryptedExtension, // Default encrypted extension
			Usage: "The extension of encrypted files (e.g., .enc, .xyz)",
		},
		&cli.BoolFlag{
			Name:  "overwrite",
			Usage: "Overwrite existing files in the output directory without warning.",
			Value: false,
		},
		&cli.StringFlag{ // New flag for output format
			Name:  "output-format",
			Value: "png", // Default output format
			Usage: "Output image format (png, jpg, jpeg)",
		},
	},
	Action: func(c *cli.Context) error {
		inputPath := c.String("input")
		outputPath := c.String("output")
		keyBase64 := c.String("key")
		recursive := c.Bool("recursive")
		encryptedExt := c.String("encrypted-ext")
		overwrite := c.Bool("overwrite")
		outputFormat := c.String("output-format") // Retrieve output format

		// Decode the key from base64
		key, err := base64.StdEncoding.DecodeString(keyBase64)
		if err != nil {
			log.Printf("failed to decode key: %v", err)
			return err
		}

		if len(key) != KeySize {
			log.Printf("invalid key size: key must be %d bytes when base64 decoded", KeySize)
			return fmt.Errorf("invalid key size: key must be %d bytes when base64 decoded", KeySize)
		}

		// Check if the input is a file or a directory
		fileInfo, err := os.Stat(inputPath)
		if err != nil {
			log.Printf("failed to stat input path: %v", err)
			return err
		}

		if fileInfo.IsDir() {
			// Process directory
			return decryptDirectory(inputPath, outputPath, key, recursive, encryptedExt, overwrite, outputFormat)
		} else {
			// Process single file
			return decryptFile(inputPath, outputPath, key, overwrite, outputFormat)
		}
	},
}

func decryptFile(inputFilename, outputFilename string, key []byte, overwrite bool, outputFormat string) error {
	// Check if the output file exists and if overwriting is allowed
	if _, err := os.Stat(outputFilename); err == nil && !overwrite {
		// File exists and overwrite is not allowed
		gookitcolor.Yellow.Printf("Output file %s already exists.  Overwrite with --overwrite flag.\n", outputFilename)
		return nil
	}
	// Read the encrypted data from the file
	ciphertext, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		log.Printf("failed to read encrypted file: %v", err)
		return err
	}

	// Decrypt the data
	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		log.Printf("failed to decrypt: %v", err)
		return err
	}

	// Convert the decrypted bytes back to an image
	img, err := BytesToImage(plaintext)
	if err != nil {
		log.Printf("failed to convert decrypted bytes to image: %v", err)
		return err
	}

	// Save the decrypted image to a file
	err = os.MkdirAll(filepath.Dir(outputFilename), os.ModeDir|0755) // Ensure output directory exists
	if err != nil {
		log.Printf("failed to create output directory: %v", err)
		return err
	}

	err = SaveImage(outputFilename, img, outputFormat) // Pass outputFormat to SaveImage
	if err != nil {
		log.Printf("failed to save decrypted image: %v", err)
		return err
	}

	gookitcolor.Cyan.Println("Image decrypted and saved to:", outputFilename)
	return nil
}

func decryptDirectory(inputDir, outputDir string, key []byte, recursive bool, encryptedExt string, overwrite bool, outputFormat string) error {
	var wg sync.WaitGroup
	err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // Propagate the error
		}

		if info.IsDir() && path != inputDir && !recursive {
			return filepath.SkipDir // Skip subdirectories if not recursive
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), encryptedExt) { // Decrypt only .enc files
			// Construct the output filename
			relPath, err := filepath.Rel(inputDir, path)
			if err != nil {
				log.Printf("failed to get relative path: %v", err)
				return err
			}

			outputFilename := filepath.Join(outputDir, strings.TrimSuffix(relPath, encryptedExt)) // Remove .enc extension

			wg.Add(1)
			go func(p, o string) {
				defer wg.Done()
				err := decryptFile(p, o, key, overwrite, outputFormat) // Pass outputFormat
				if err != nil {
					log.Printf("Error decrypting %s: %v\n", p, err)
				}
			}(path, outputFilename) // Decrypt each image file
		}
		return nil
	})

	wg.Wait()
	if err != nil {
		log.Printf("error walking the path %s: %v", inputDir, err)
		return err
	}

	return nil
}

var keygenCmd = &cli.Command{
	Name:  "keygen",
	Usage: "Generate a new encryption key",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "output",
			Value: "",
			Usage: "File to save the generated key to",
		},
	},
	Action: func(c *cli.Context) error {
		keyFile := c.String("output")
		key, err := GenerateRandomKey()
		if err != nil {
			log.Printf("failed to generate key: %v", err)
			return err
		}

		keyBase64Encoded := base64.StdEncoding.EncodeToString(key)

		if keyFile != "" {
			// Save the key to a file
			err = ioutil.WriteFile(keyFile, []byte(keyBase64Encoded), 0600) // Permissions 0600: read/write for owner only
			if err != nil {
				log.Printf("failed to save key to file: %v", err)
				return err
			}
			gookitcolor.Green.Println("Generated Key (base64 encoded):", keyBase64Encoded)
			gookitcolor.Green.Println("Key saved to file:", keyFile)
		} else {
			gookitcolor.Green.Println("Generated Key (base64 encoded):", keyBase64Encoded)
		}

		return nil
	},
}

// steganographyCmd implements steganography features
var steganographyCmd = &cli.Command{
	Name:  "stego",
	Usage: "Hide or reveal a message within an image using steganography",
	Subcommands: []*cli.Command{
		{
			Name:  "hide",
			Usage: "Hide a message within an image",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "input",
					Aliases:  []string{"i"},
					Value:    "",
					Usage:    "Input image file",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "output",
					Aliases:  []string{"o"},
					Value:    "stego_output.png",
					Usage:    "Output stego image file",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "message",
					Aliases:  []string{"m"},
					Value:    "",
					Usage:    "Message to hide",
					Required: true,
				},
				&cli.StringFlag{ // New flag for output format
					Name:  "output-format",
					Value: "png", // Default output format
					Usage: "Output image format (png, jpg, jpeg)",
				},
			},
			Action: func(c *cli.Context) error {
				inputPath := c.String("input")
				outputPath := c.String("output")
				message := c.String("message")
				outputFormat := c.String("output-format")

				if len(message) > StegoMessageLimit {
					gookitcolor.Red.Println("Message too long. Max message length is", StegoMessageLimit, "characters.")
					return fmt.Errorf("message too long. Max message length is %d characters", StegoMessageLimit)
				}

				return hideMessage(inputPath, outputPath, message, outputFormat)
			},
		},
		{
			Name:  "reveal",
			Usage: "Reveal a hidden message from an image",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "input",
					Aliases:  []string{"i"},
					Value:    "",
					Usage:    "Input stego image file",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				inputPath := c.String("input")
				message, err := revealMessage(inputPath)
				if err != nil {
					gookitcolor.Red.Println(fmt.Errorf("failed to reveal message: %w", err))
					return err
				}
				gookitcolor.Green.Println("Hidden Message:", message)
				return nil
			},
		},
	},
}

// hideMessage hides a message within an image using LSB steganography
// hideMessage hides a message within an image using LSB steganography
func hideMessage(inputFilename, outputFilename, message string, outputFormat string) error {
	img, err := LoadImage(inputFilename)
	if err != nil {
		log.Printf("failed to load image: %v", err)
		return err
	}

	// Convert to RGBA
	b := img.Bounds()
	rgbaImg := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, b.Min, draw.Src)

	bounds := rgbaImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	message += "\x00" // Null terminate the message

	messageIndex := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if messageIndex < len(message) {
				r, g, b, a := rgbaImg.At(x, y).RGBA()
				// Modify the least significant bit of the red channel to store the message
				r = (r &^ 1) | uint32(message[messageIndex]>>7)
				g = (g &^ 1) | uint32((message[messageIndex]>>6)&1)
				b = (b &^ 1) | uint32((message[messageIndex]>>5)&1)
				a = (a &^ 1) | uint32((message[messageIndex]>>4)&1)
				rgbaImg.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
				messageIndex++
			}
		}
	}

	err = os.MkdirAll(filepath.Dir(outputFilename), os.ModeDir|0755) // Ensure output directory exists
	if err != nil {
		log.Printf("failed to create output directory: %v", err)
		return err
	}

	err = SaveImage(outputFilename, rgbaImg, outputFormat) // Save using the specified output format
	if err != nil {
		log.Printf("failed to encode stego image: %v", err)
		return err
	}
	gookitcolor.Cyan.Println("Message hidden and saved to:", outputFilename)
	return nil
}

// revealMessage reveals a hidden message from an image
func revealMessage(inputFilename string) (string, error) {
	img, err := LoadImage(inputFilename)
	if err != nil {
		log.Printf("failed to load image: %v", err)
		return "", err
	}

	// Convert to RGBA
	b := img.Bounds()
	rgbaImg := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(rgbaImg, rgbaImg.Bounds(), img, b.Min, draw.Src)

	bounds := rgbaImg.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var messageBits bytes.Buffer
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := rgbaImg.At(x, y).RGBA()
			messageBits.WriteByte(uint8(((r & 1) << 7) | ((g & 1) << 6) | ((b & 1) << 5) | ((a & 1) << 4))) // Extract LSB of each channel

		}
	}
	message := messageBits.String()
	message = strings.Split(message, "\x00")[0] // Read until null terminator
	return message, nil
}

// main function
func main() {
	cli.VersionFlag = &cli.BoolFlag{ //Add the version flag
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the version",
	}
	cli.HelpFlag = &cli.BoolFlag{ //Override the help flag so that it is at the end
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show help",
	}
	app := &cli.App{
		Name:    "image-encryption",
		Usage:   "Encrypt, decrypt, and hide messages within images using AES-256 GCM and steganography",
		Version: Version, //Set the version from the constant
		Authors: []*cli.Author{
			{
				Name:  Author, // Set the Author from constant
				Email: "",     // Can add an email here
			},
		},
		Commands: []*cli.Command{
			encryptCmd,
			decryptCmd,
			keygenCmd,
			steganographyCmd,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Value: false,
				Usage: "Enable verbose logging",
			},
			&cli.BoolFlag{
				Name:    "about",
				Aliases: []string{"a"},
				Usage:   "About this tool",
			},
		},
		Before: func(c *cli.Context) error {
			// Print AsciiArt on startup
			gookitcolor.HiBlue.Println(AsciiArt)

			if c.Bool("verbose") {
				log.SetFlags(log.LstdFlags | log.Lshortfile) // Enhanced logging
				log.Println("Verbose mode enabled")
			}

			if c.Bool("about") {
				fmt.Printf("Image Encryption Tool\n")
				fmt.Printf("Version: %s\n", Version)
				fmt.Printf("Author: %s\n", Author)
				fmt.Printf("Go Version: %s\n", runtime.Version())
				fmt.Printf("Operating System: %s %s\n", runtime.GOOS, runtime.GOARCH)
				os.Exit(0) // Exit after printing about information
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
