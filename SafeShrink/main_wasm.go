//go:build js && wasm

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
	"syscall/js"

	"github.com/nfnt/resize"
)

type CompressResult struct {
	Success          bool    `json:"success"`
	CompressedData   string  `json:"compressedData"`
	OriginalSize     int64   `json:"originalSize"`
	CompressedSize   int64   `json:"compressedSize"`
	CompressionRatio float64 `json:"compressionRatio"`
	Message          string  `json:"message,omitempty"`
}

type CompressOptions struct {
	Quality    int  `json:"quality"`
	MaxWidth   uint `json:"maxWidth,omitempty"`
	MaxHeight  uint `json:"maxHeight,omitempty"`
	KeepAspect bool `json:"keepAspect"`
}

func main() {
	c := make(chan struct{}, 0)

	// Register the compressImage function
	js.Global().Set("compressImage", js.FuncOf(compressImageWrapper))
	js.Global().Set("getSupportedFormats", js.FuncOf(getSupportedFormats))

	fmt.Println("SafeShrink WASM module loaded successfully!")
	<-c
}

func compressImageWrapper(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{
			"success": false,
			"message": "Invalid arguments: expected base64Data and options",
		}
	}

	// Parse arguments
	base64Data := args[0].String()
	optsJS := args[1]

	// Parse options
	opts := CompressOptions{
		Quality:    85,
		KeepAspect: true,
	}

	if optsJS.Get("quality").Type() != js.TypeUndefined {
		opts.Quality = optsJS.Get("quality").Int()
	}
	if optsJS.Get("maxWidth").Type() != js.TypeUndefined {
		opts.MaxWidth = uint(optsJS.Get("maxWidth").Int())
	}
	if optsJS.Get("maxHeight").Type() != js.TypeUndefined {
		opts.MaxHeight = uint(optsJS.Get("maxHeight").Int())
	}
	if optsJS.Get("keepAspect").Type() != js.TypeUndefined {
		opts.KeepAspect = optsJS.Get("keepAspect").Bool()
	}

	// Remove data URL prefix if present
	if strings.HasPrefix(base64Data, "data:") {
		if idx := strings.Index(base64Data, ","); idx != -1 {
			base64Data = base64Data[idx+1:]
		}
	}

	// Decode base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Failed to decode base64 data: %v", err),
		}
	}

	// Compress the image
	compressedData, err := compressImage(imageData, opts)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Compression failed: %v", err),
		}
	}

	// Calculate compression ratio
	originalSize := int64(len(imageData))
	compressedSize := int64(len(compressedData))
	ratio := (1.0 - float64(compressedSize)/float64(originalSize)) * 100

	// Encode compressed data to base64
	compressedBase64 := base64.StdEncoding.EncodeToString(compressedData)

	return map[string]interface{}{
		"success":          true,
		"compressedData":   compressedBase64,
		"originalSize":     originalSize,
		"compressedSize":   compressedSize,
		"compressionRatio": ratio,
	}
}

func getSupportedFormats(this js.Value, args []js.Value) interface{} {
	return map[string]interface{}{
		"supported": []string{"jpeg", "jpg", "png"},
		"output":    []string{"jpeg", "png"},
	}
}

func compressImage(data []byte, opts CompressOptions) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Resize if needed
	if opts.MaxWidth > 0 || opts.MaxHeight > 0 {
		bounds := img.Bounds()
		width := uint(bounds.Dx())
		height := uint(bounds.Dy())

		newWidth := width
		newHeight := height

		// Calculate new dimensions
		if opts.MaxWidth > 0 && width > opts.MaxWidth {
			newWidth = opts.MaxWidth
			if opts.KeepAspect {
				newHeight = uint(float64(height) * float64(newWidth) / float64(width))
			}
		}

		if opts.MaxHeight > 0 && newHeight > opts.MaxHeight {
			newHeight = opts.MaxHeight
			if opts.KeepAspect {
				newWidth = uint(float64(newWidth) * float64(newHeight) / float64(height))
			}
		}

		// Resize if dimensions changed
		if newWidth != width || newHeight != height {
			img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
		}
	}

	// Encode with compression
	var buf bytes.Buffer

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: opts.Quality})
	case "png":
		// Convert PNG to JPEG for better compression if quality < 100
		if opts.Quality < 100 {
			err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: opts.Quality})
		} else {
			err = png.Encode(&buf, img)
		}
	default:
		// Default to JPEG for unsupported formats
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: opts.Quality})
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	// If compressed is larger, try with lower quality or return original
	compressedData := buf.Bytes()
	if len(compressedData) >= len(data) && opts.Quality > 50 {
		buf.Reset()
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: opts.Quality - 20})
		if err == nil {
			newCompressed := buf.Bytes()
			if len(newCompressed) < len(data) {
				return newCompressed, nil
			}
		}
		return data, nil
	}

	return compressedData, nil
}
