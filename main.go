package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/JonCSykes/DragonTable/mapFile"
	"github.com/JonCSykes/DragonTable/widgetExt"
	"github.com/lxn/win"
)

var ScreenHeight int
var ScreenWidth int
var mainWindow fyne.Window
var mapFiles []*mapFile.MapFile
var mainContent *fyne.Container

func main() {

	GetScreenResolution()

	myApp := app.New()
	mainWindow = myApp.NewWindow("Dragon Table")

	BuildUI()

	mainWindow.SetContent(mainContent)

	mainWindow.SetPadded(true)
	//mainWindow.SetFullScreen(true)
	mainWindow.ShowAndRun()
}

func GetScreenResolution() {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	ScreenWidth = int(win.GetDeviceCaps(hDC, win.HORZRES))
	ScreenHeight = int(win.GetDeviceCaps(hDC, win.VERTRES))

	fmt.Println(strconv.Itoa(ScreenWidth) + " x " + strconv.Itoa(ScreenHeight))
}

func BuildUI() {

	mapList := BuildNavList()
	touchControlButton := BuildNavButtons()
	images := BuildImageDisplay()
	content := container.NewWithoutLayout()

	for _, image := range images {
		content.Add(image)
	}

	content.Add(mapList)
	content.Add(touchControlButton)

	mainContent = content
}

func BuildImageDisplay() []*canvas.Image {

	if len(mapFiles) == 0 {
		mapFiles = mapFile.GetMaps()
	}

	var images []*canvas.Image
	for i, mapFile := range mapFiles {
		if i > 1 {
			mapFile.Image.Hide()
		}

		mapFile.Image.Resize(fyne.NewSize(float32(ScreenWidth), float32(ScreenHeight)))
		mapFile.Image.FillMode = canvas.ImageFillContain
		mapFile.Image.Move(fyne.Position{X: 0, Y: 0})

		images = append(images, mapFile.Image)
	}
	return images
}

func BuildNavList() *widget.List {

	if len(mapFiles) == 0 {
		mapFiles = mapFile.GetMaps()
	}

	mapList := widget.NewList(
		func() int {
			return len(mapFiles)
		},
		func() fyne.CanvasObject {
			return widgetExt.NewImageButton("", nil, nil)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			o.(*widgetExt.ImageButton).SetText(mapFiles[i].FileName)
			o.(*widgetExt.ImageButton).OnTapped = func() {
				fmt.Println("Clicked : " + mapFiles[i].FileName)
				fileName := mapFiles[i].FileName + "." + mapFiles[i].Extension

				HideMaps()
				ShowMap(fileName)
			}
			o.(*widgetExt.ImageButton).Resize(fyne.Size{Width: 400, Height: 50})
			o.(*widgetExt.ImageButton).SetImage(mapFiles[i].ThumbResource)
		})

	mapList.Resize(fyne.NewSize(400, 1000))
	mapList.Move(fyne.Position{X: float32(ScreenWidth) - 410, Y: 10})

	return mapList
}

func HideMaps() {
	for _, child := range mainContent.Objects {
		switch x := child.(type) {
		case *canvas.Image:
			x.Hide()
		}
	}
}

func ShowMap(fileName string) {
	for _, child := range mainContent.Objects {
		switch x := child.(type) {
		case *canvas.Image:

			if fileName == x.Resource.Name() {
				x.Show()
			}
		}
	}
}

func BuildNavButtons() *widget.Button {

	var touchControlButton *widget.Button

	disabledTouchIcon, enableTouchError := fyne.LoadResourceFromPath("./resources/icons/hand-point-up-regular.svg")
	if enableTouchError != nil {
		fmt.Println(enableTouchError)
	}

	enabledTouchIcon, enableTouchError := fyne.LoadResourceFromPath("./resources/icons/hand-point-up-solid.svg")
	if enableTouchError != nil {
		fmt.Println(enableTouchError)
	}

	touchControlButton = widget.NewButtonWithIcon("", enabledTouchIcon, func() {

		if touchControlButton.Icon == disabledTouchIcon {
			touchControlButton.Importance = widget.HighImportance
			touchControlButton.SetIcon(enabledTouchIcon)
		} else {
			touchControlButton.Importance = widget.LowImportance
			touchControlButton.SetIcon(disabledTouchIcon)
		}
	})

	touchControlButton.Importance = widget.HighImportance
	touchControlButton.Resize(fyne.NewSize(100, 100))
	touchControlButton.Move(fyne.Position{X: float32(ScreenWidth) - 110, Y: float32(ScreenHeight) - 110})

	return touchControlButton
}
