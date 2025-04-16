package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func main() {
	// 1. Загружаем изображение
	imgPath := "input.jpg"
	img, err := loadImage(imgPath)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	// 2. Добавляем текст
	memeText := "Sample Meme Text"
	result := addTextToImage(img, memeText)

	// 3. Сохраняем результат
	outputPath := "output.jpg"
	err = saveImage(outputPath, result)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Meme generated successfully!")
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

func addTextToImage(img image.Image, text string) image.Image {
	// Создаем новое изображение с тем же размером
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)

	// Настройки текста
	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(image.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(20 * 64), Y: fixed.Int26_6((bounds.Max.Y - 20) * 64)},
	}

	// Рисуем текст
	d.DrawString(text)

	return rgba
}

func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch filepath.Ext(path) {
	case ".png":
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	}
}
