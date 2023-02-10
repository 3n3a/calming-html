package main

import (
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/webp"

	"github.com/nickalie/go-webpbin"
)

type ImageTemplateElement struct {
    Index int
    Filename string
    LoadingPriority string
    Width int
    Height int
}

type CalmingImagesData struct {
    PageTitle string
    Images []ImageTemplateElement
}

func main() {
    DIRECTORY := "public/images"
    LAYOUT_NAME := "layout.html"
    OUT_PATH := "public/index.html"
    IMAGE_WEB_BASE_PATH := "https://ik.imagekit.io/njhey0rxzni"

    err := tranformAllImagesToWebp(DIRECTORY)
    if err != nil {
        log.Fatal(err)
    }

    images, err := getImagesForSite(DIRECTORY, IMAGE_WEB_BASE_PATH, "webp")
    if err != nil {
        log.Fatal(err)
    }

    // Fill out Template with List of Images
    data := CalmingImagesData{
        PageTitle: "Calming Images",
        Images: images,
    }

    t := template.Must(template.ParseFiles(LAYOUT_NAME))
    out, err := os.Create(OUT_PATH)
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()
    
    // Save filled out TEmplate
    err = t.Execute(out, data)
    if err != nil {
        log.Fatal(err)
    }
    
}

func getImagesForSite(folderPath string, imageBasePathWeb string, targetExtension string) ([]ImageTemplateElement, error) {
    // Get List of Images in Folder X
    files, err := ioutil.ReadDir(folderPath)
    if err != nil {
        log.Fatal(err)
    }

    var images []ImageTemplateElement
    for i, file := range files {
        if !file.IsDir() && (getFileExtension(file.Name()) == targetExtension) {            
            // Get Image Size
            i_width, i_height := getImageSize(folderPath, file.Name(), targetExtension)

            // Define Loading Priority
            loadingPriority := "eager"
            if i+1 > 4 {
                loadingPriority = "lazy"
            }

            images = append(images, ImageTemplateElement{
                Index: i+1,
                Filename: fmt.Sprintf("%s/%s", imageBasePathWeb, file.Name()),
                LoadingPriority: loadingPriority,
                Width: i_width,
                Height: i_height,
            })
        }
    }
    return images, nil
}

func tranformAllImagesToWebp(folderPath string) error {
    files, err := ioutil.ReadDir(folderPath)
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        if !file.IsDir() && (getFileExtension(file.Name()) != "webp") {
            // Transform to webp
            newName, err := transformToWebp(folderPath, file.Name())
            if err != nil {
                log.Fatal(err)
            }
            fmt.Println("Transformed", newName)
        }
    }
    return nil
}

func transformToWebp(basePath string, imageName string) (string, error) {
    rawName := getFilenameNoExtension(imageName)
    newName := fmt.Sprintf("%s.webp", rawName)

    err := webpbin.NewCWebP().
        Quality(60).
        InputFile(filepath.Join(basePath, imageName)).
        OutputFile(filepath.Join(basePath, newName)).
        Run()

    return newName, err
}

func getFileExtension(fileName string) string {
    s := strings.Split(fileName, ".")
    return s[len(s)-1]
}

func getFilenameNoExtension(fileName string) string {
    s := strings.Split(fileName, ".")
    rawName := s[0]
    return rawName
}

func getImageSize(basePath string, imageName string, imageType string) (int, int) {
    if reader, err := os.Open(filepath.Join(basePath, imageName)); err == nil {
        defer reader.Close()

        var image image.Config
        var err error
        switch imageType {
        case "jpeg":
            image, err = jpeg.DecodeConfig(reader)
        case "webp":
            image, err = webp.DecodeConfig(reader)
        }
        if err != nil {
            log.Fatal(err)
        }
        return image.Width, image.Height
    } else {
        log.Fatal(err)
    }
    return 0,0
}