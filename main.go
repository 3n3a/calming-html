package main

import (
	"fmt"
	"html/template"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

    // Get List of Images in Folder X
    files, err := ioutil.ReadDir(DIRECTORY)
    if err != nil {
        log.Fatal(err)
    }

    var images []ImageTemplateElement
    for i, file := range files {
        if !file.IsDir() {

            // Get Image Size
            i_width, i_height := getImageSize(DIRECTORY, file.Name())

            // Define Loading Priority
            loadingPriority := "eager"
            if i+1 > 4 {
                loadingPriority = "lazy"
            }

            images = append(images, ImageTemplateElement{
                Index: i+1,
                Filename: fmt.Sprintf("images/%s", file.Name()),
                LoadingPriority: loadingPriority,
                Width: i_width,
                Height: i_height,
            })
        }
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


func getImageSize(basePath string, imageName string) (int, int) {
    if reader, err := os.Open(filepath.Join(basePath, imageName)); err == nil {
        defer reader.Close()
        fmt.Println(reader.Name())
        image, err := jpeg.DecodeConfig(reader)
        if err != nil {
            log.Fatal(err)
        }
        return image.Width, image.Height
    } else {
        log.Fatal(err)
    }
    return 0,0
}