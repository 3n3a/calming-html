package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

type Filepath struct {
    Filename string
}

type CalmingImagesData struct {
    PageTitle string
    Images []Filepath
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

    var images []Filepath
    for _, file := range files {
        if !file.IsDir() {
            images = append(images, Filepath{Filename: fmt.Sprintf("images/%s", file.Name())})
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
