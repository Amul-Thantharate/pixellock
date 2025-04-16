# 🎯 PixelLock Demo Guide

This guide demonstrates common use cases and examples for the PixelLock CLI tool.

## 📝 Basic Examples

### 1. Key Generation
```bash
# Generate and display a new key
bin/image-encryption keygen

# Generate and save key to file
bin/image-encryption keygen --output mykey.key
```

### 2. Single Image Encryption
```bash
# Basic image encryption
bin/image-encryption encrypt --input myimage.png --output encrypted_image.enc

# Encrypt using existing key
bin/image-encryption encrypt --input myimage.png --output encrypted_image.enc --key "YOUR_BASE64_KEY"

# Encrypt with key from environment variable
export IMAGE_ENCRYPTION_KEY="YOUR_BASE64_KEY"
bin/image-encryption encrypt --input myimage.png --output encrypted_image.enc
```

### 3. Single Image Decryption
```bash
# Basic decryption (PNG output)
bin/image-encryption decrypt --input encrypted_image.enc --output decrypted.png --key "YOUR_BASE64_KEY"

# Decrypt to JPEG format
bin/image-encryption decrypt --input encrypted_image.enc --output decrypted.jpg --key "YOUR_BASE64_KEY" --output-format jpg
```

## 🗂️ Batch Processing Examples

### 1. Directory Encryption
```bash
# Encrypt directory (non-recursive)
bin/image-encryption encrypt --input images_dir --output encrypted_images

# Encrypt directory recursively
bin/image-encryption encrypt --input images_dir --output encrypted_images --recursive

# Encrypt with overwrite option
bin/image-encryption encrypt --input images_dir --output encrypted_images --recursive --overwrite
```

### 2. Directory Decryption
```bash
# Decrypt directory
bin/image-encryption decrypt --input encrypted_images --output decrypted_images --key "YOUR_BASE64_KEY"

# Decrypt with custom extension
bin/image-encryption decrypt --input encrypted_files --output decrypted_files --key "YOUR_BASE64_KEY" --encrypted-ext .xyz
```

## 🕵️ Steganography Examples

### 1. Hide Messages
```bash
# Hide message (PNG output)
bin/image-encryption stego hide --input original.png --output hidden.png --message "Secret message"

# Hide message (JPEG output)
bin/image-encryption stego hide --input original.png --output hidden.jpg --message "Secret message" --output-format jpg
```

### 2. Reveal Messages
```bash
# Extract hidden message
bin/image-encryption stego reveal --input hidden.png
```

## 🔐 Security Best Practices

1. **Key Management**
```bash
# Generate and save key securely
bin/image-encryption keygen --output secure.key
chmod 600 secure.key

# Use environment variable
export IMAGE_ENCRYPTION_KEY=$(cat secure.key)
```

2. **Batch Processing with Key File**
```bash
# Read key from file and process directory
KEY=$(cat secure.key)
bin/image-encryption encrypt --input photos/ --output encrypted/ --key "$KEY" --recursive
```

## 🐳 Docker Examples

```bash
# Build and run with Docker
docker build -t image-encryption .
docker run -v $(pwd):/data image-encryption encrypt --input /data/photo.png --output /data/encrypted.enc

# Interactive mode
docker run -it -v $(pwd):/data image-encryption
```

## 📊 Example Results

```bash
# Key Generation Output
$ bin/image-encryption keygen
Generated Key (base64 encoded): dGhpcyBpcyBhbiBleGFtcGxlIGtleSBmb3IgZGVtbyBwdXJwb3NlcyE=

# Encryption Output
$ bin/image-encryption encrypt --input photo.png --output encrypted.enc
Image encrypted and saved to: encrypted.enc

# Decryption Output
$ bin/image-encryption decrypt --input encrypted.enc --output decrypted.png --key "YOUR_KEY"
Image decrypted and saved to: decrypted.png

# Steganography Output
$ bin/image-encryption stego reveal --input hidden.png
Hidden Message: Your secret message here
```

## 🚀 Quick Start Script

```bash
#!/bin/bash
# Quick demo script

# Generate a key
echo "Generating encryption key..."
KEY=$(bin/image-encryption keygen | grep "Generated Key" | cut -d':' -f2 | tr -d ' ')

# Create test image
convert -size 100x100 xc:white test.png

# Encrypt
echo "Encrypting test image..."
bin/image-encryption encrypt --input test.png --output encrypted.enc --key "$KEY"

# Decrypt
echo "Decrypting test image..."
bin/image-encryption decrypt --input encrypted.enc --output decrypted.png --key "$KEY"

# Hide message
echo "Testing steganography..."
bin/image-encryption stego hide --input test.png --output hidden.png --message "Hello, World!"
bin/image-encryption stego reveal --input hidden.png
