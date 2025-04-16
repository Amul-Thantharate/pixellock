# 🔐 PixelLock - Image Encryption & Steganography Tool

A powerful CLI tool for encrypting images and hiding secret messages using steganography, built with Go. PixelLock provides enterprise-grade security for your sensitive images and allows for covert communication through steganographic techniques.

![Version](https://img.shields.io/badge/version-v1.0.0-blue)
![Go Version](https://img.shields.io/badge/go-1.24-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

## 🌟 Features

- 🔒 **Image Encryption**: Secure your images using AES-256 GCM encryption, a highly secure authenticated encryption mode that provides both confidentiality and data authenticity
- 📁 **Batch Processing**: Process multiple images in directories recursively, making it easy to secure entire collections of sensitive images at once
- 💌 **Steganography**: Hide and reveal secret messages within images without visible changes, using advanced LSB (Least Significant Bit) techniques
- 🎨 **Multiple Format Support**: Works with PNG, JPEG, and other common image formats, maintaining compatibility with your existing workflows
- 🔑 **Key Management**: Generate and manage encryption keys securely with built-in tools for key generation, storage, and retrieval

## 📋 Prerequisites

- Go 1.24 or higher
- Make (optional, for using Makefile commands)
- Docker (optional, for containerized usage)

## 🚀 Installation

### Using Go

```bash
# Clone the repository
git clone https://github.com/Amul-Thantharate/pixellock.git

# Navigate to the project directory
cd pixellock

# Install dependencies
make install-deps

# Build the project
make build
```

### Using Docker

```bash
# Build Docker image
make docker-build

# Run using Docker
make docker-run
```

## 💡 Usage

### Encrypt Images

PixelLock uses AES-256 in GCM mode, providing authenticated encryption that ensures both confidentiality and integrity of your images.

```bash
# Encrypt a single image
pixellock encrypt -i input.png -o encrypted.enc -k <base64-key>

# Encrypt a directory of images recursively
pixellock encrypt -i images/ -o encrypted/ -r
```

### Decrypt Images

Decrypt your images using the same key that was used for encryption. The authentication feature of GCM ensures that tampered files will be detected during decryption.

```bash
# Decrypt a single image
pixellock decrypt -i encrypted.enc -o decrypted.png -k <base64-key>

# Decrypt a directory of images
pixellock decrypt -i encrypted/ -o decrypted/ -r
```

### Steganography

The steganography feature uses sophisticated algorithms to embed data within the least significant bits of image pixels, making the changes imperceptible to the human eye and resistant to statistical analysis.

```bash
# Hide a message in an image
pixellock stego hide -i input.png -o output.png -m "Secret message"

# Reveal a hidden message
pixellock stego reveal -i output.png
```

### Generate Encryption Key

PixelLock's key generation uses a cryptographically secure random number generator to create high-entropy keys suitable for AES-256 encryption.

```bash
# Generate and display a new key
pixellock keygen

# Generate and save key to file
pixellock keygen --output mykey.key
```

## 🛠 Available Commands

- `encrypt` (aliases: `e`): Encrypt images using AES-256 GCM for maximum security
- `decrypt` (aliases: `d`): Decrypt previously encrypted images with authentication
- `keygen`: Generate cryptographically secure encryption keys of appropriate length
- `stego`: Steganography operations for covert communication
  - `hide`: Hide messages in images using advanced LSB techniques
  - `reveal`: Extract hidden messages without damaging the carrier image

## 🔧 Makefile Commands

- `make build`: Build the application with optimized settings
- `make run`: Build and run the application with default parameters
- `make test`: Run comprehensive test suite including unit and integration tests
- `make clean`: Clean build artifacts and temporary files
- `make format`: Format code according to Go best practices
- `make docker-build`: Build Docker image with minimal footprint
- `make docker-run`: Run in Docker container with appropriate volume mounts
- `make help`: Show all available commands with detailed descriptions

## 🔐 Security Notes

- Always store encryption keys securely in a password manager or hardware security module
- Use environment variable `IMAGE_ENCRYPTION_KEY` for automated processes to avoid key exposure in command history
- Back up your encryption keys - lost keys mean unrecoverable images with no backdoor recovery options
- Default encryption uses AES-256 GCM, providing 256-bit security with authenticated encryption
- The tool implements secure memory handling to minimize the risk of key exposure through memory dumps

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

Amul Thantharate

## 🙏 Acknowledgments

- [urfave/cli](https://github.com/urfave/cli) for CLI framework
- [gookit/color](https://github.com/gookit/color) for terminal colors
- Cryptography community for best practices and security standards
- Open source contributors who provided feedback and suggestions
