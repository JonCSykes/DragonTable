package mapFile

import (
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gxcbuf/graphics-go/graphics"
)

// MapFile :
type MapFile struct {
	Path              string
	FileName          string
	Extension         string
	FullPath          string
	FullThumbnailPath string
}

const MapPath string = "./resources/maps"

func InitMapFile(fullPath string) *MapFile {

	lastSlash := strings.LastIndex(fullPath, "/")
	path := fullPath[:lastSlash+1]
	fullFileName := fullPath[lastSlash+1:]
	extension := fullFileName[strings.LastIndex(fullFileName, ".")+1:]
	fileName := fullFileName[:strings.LastIndex(fullFileName, ".")]

	newMapFile := MapFile{FileName: fileName, Path: path, Extension: extension, FullPath: fullPath}

	newMapFile.GenerateThumb()

	return &newMapFile
}

func (mapFile *MapFile) GenerateThumb() {

	imagePath, _ := os.Open(mapFile.FullPath)
	defer imagePath.Close()
	srcImage, _, _ := image.Decode(imagePath)

	dstImage := image.NewRGBA(image.Rect(0, 0, 400, 100))

	graphics.Thumbnail(dstImage, srcImage)

	fullThumbnailPath := mapFile.Path + mapFile.FileName + "_thumb." + mapFile.Extension

	newImage, _ := os.Create(fullThumbnailPath)
	defer newImage.Close()
	jpeg.Encode(newImage, dstImage, &jpeg.Options{Quality: jpeg.DefaultQuality})

	mapFile.FullThumbnailPath = fullThumbnailPath
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
