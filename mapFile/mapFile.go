package mapFile

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/gxcbuf/graphics-go/graphics"
)

// MapFile :
type MapFile struct {
	Path              string
	FileName          string
	Extension         string
	FullPath          string
	FullThumbnailPath string
	Image             *canvas.Image
	ImageResource     fyne.Resource
	ThumbResource     fyne.Resource
}

const MapPath string = "./resources/maps"

func InitMapFile(fullPath string) *MapFile {

	lastSlash := strings.LastIndex(fullPath, "/")
	path := fullPath[:lastSlash+1]
	fullFileName := fullPath[lastSlash+1:]
	extension := fullFileName[strings.LastIndex(fullFileName, ".")+1:]
	fileName := fullFileName[:strings.LastIndex(fullFileName, ".")]

	imageResource, imageError := fyne.LoadResourceFromPath(fullPath)
	if imageError != nil {
		fmt.Println(imageError)
	}

	image := canvas.NewImageFromResource(imageResource)

	newMapFile := MapFile{FileName: fileName, Path: path, Extension: extension, FullPath: fullPath, ImageResource: imageResource, Image: image}

	newMapFile.GenerateThumb()

	return &newMapFile
}

func (mapFile *MapFile) GenerateThumb() {

	imagePath, _ := os.Open(mapFile.FullPath)
	defer imagePath.Close()
	srcImage, _, _ := image.Decode(imagePath)

	dstImage := image.NewRGBA(image.Rect(0, 0, 200, 50))

	graphics.Thumbnail(dstImage, srcImage)

	fullThumbnailPath := mapFile.Path + mapFile.FileName + "_thumb." + mapFile.Extension

	newImage, _ := os.Create(fullThumbnailPath)
	defer newImage.Close()
	jpeg.Encode(newImage, dstImage, &jpeg.Options{Quality: jpeg.DefaultQuality})

	mapThumb, thumbError := fyne.LoadResourceFromPath(fullThumbnailPath)
	if thumbError != nil {
		fmt.Println(thumbError)
	}

	mapFile.FullThumbnailPath = fullThumbnailPath
	mapFile.ThumbResource = mapThumb
}

func GetMaps() []*MapFile {

	var mapFiles []*MapFile
	files, err := ioutil.ReadDir(MapPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), "_thumb") {
			mapFile := InitMapFile(MapPath + "/" + file.Name())
			mapFiles = append(mapFiles, mapFile)
		}
	}

	return mapFiles
}
