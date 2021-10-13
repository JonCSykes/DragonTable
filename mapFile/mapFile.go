package mapFile

import (
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	Height            int
	Width             int
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

	dstImage := image.NewRGBA(image.Rect(0, 0, 250, 50))

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

			fileReader, err := os.Open(MapPath + "/" + file.Name())
			if err != nil {
				fmt.Println("Impossible to open the file:", err)
			} else {
				defer fileReader.Close()
				imageConfig, _, err := image.DecodeConfig(fileReader)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %v\n", file.Name(), err)
					continue
				}

				mapFile.Width = imageConfig.Width
				mapFile.Height = imageConfig.Height

				mapFile.Image.Resize(fyne.NewSize(float32(imageConfig.Width), float32(imageConfig.Height)))
				mapFile.Image.SetMinSize(fyne.NewSize(float32(imageConfig.Width), float32(imageConfig.Height)))

				fmt.Println("Created " + mapFile.FileName + " : " + strconv.Itoa(mapFile.Width) + "x" + strconv.Itoa(mapFile.Height))
			}

			mapFiles = append(mapFiles, mapFile)
		}
	}

	return mapFiles
}
