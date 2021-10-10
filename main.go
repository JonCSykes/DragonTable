package main

import (
	"fmt"
	"image/color"
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

const ScreenDimensionWidth int = 23
const ScreenDimensionHeight int = 13

const DragonTableWallpaperPath string = "./resources/images/dragontable.jpg"

var ScreenHeight int
var ScreenWidth int

var mapFiles []*mapFile.MapFile
var mainContent *fyne.Container

func main() {

	GetScreenResolution()

	myApp := app.New()
	mainWindow := myApp.NewWindow("Dragon Table - v0.1")

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
	navButtons := BuildNavButtons()
	images := BuildImageDisplay()
	content := container.NewWithoutLayout()
	lines := DrawGrid()

	for _, image := range images {
		content.Add(image)
	}

	for _, line := range lines {
		content.Add(line)
	}

	for _, navButton := range navButtons {
		content.Add(navButton)
	}

	content.Add(mapList)

	mainContent = content

	fmt.Println("Execute BuildUI")
	fmt.Println(len(mainContent.Objects))
}

func BuildImageDisplay() []*canvas.Image {

	fmt.Println("Execute BuildImageDisplay")
	if len(mapFiles) == 0 {
		fmt.Println("Execute GetMaps in BuildImageDisplay")
		mapFiles = mapFile.GetMaps()
	}

	var images []*canvas.Image

	dragonTableWallpaperResource, imageError := fyne.LoadResourceFromPath(DragonTableWallpaperPath)
	if imageError != nil {
		fmt.Println(imageError)
	}

	dragonTableImage := canvas.NewImageFromResource(dragonTableWallpaperResource)
	dragonTableImage.Resize(fyne.NewSize(float32(ScreenWidth), float32(ScreenHeight)))
	dragonTableImage.FillMode = canvas.ImageFillContain
	dragonTableImage.Move(fyne.Position{X: 0, Y: 0})

	images = append(images, dragonTableImage)

	for _, mapFile := range mapFiles {

		mapFile.Image.Hide()

		mapFile.Image.Resize(fyne.NewSize(float32(ScreenWidth), float32(ScreenHeight)))
		mapFile.Image.FillMode = canvas.ImageFillContain
		mapFile.Image.Move(fyne.Position{X: 0, Y: 0})

		images = append(images, mapFile.Image)
	}
	return images
}

func BuildNavList() *widget.List {

	fmt.Println("Execute BuildNavList")

	if len(mapFiles) == 0 {
		fmt.Println("Execute GetMaps in BuildNavList")
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
			o.(*widgetExt.ImageButton).Resize(fyne.Size{Width: 200, Height: 50})
			o.(*widgetExt.ImageButton).SetImage(mapFiles[i].ThumbResource)
		})

	mapList.Resize(fyne.NewSize(200, 1000))
	mapList.Move(fyne.Position{X: float32(ScreenWidth) - 210, Y: 80})

	return mapList
}

func HideMaps() {

	fmt.Println("Execute HideMaps")
	fmt.Println(len(mainContent.Objects))

	for _, child := range mainContent.Objects {
		switch x := child.(type) {
		case *canvas.Image:
			if x.Resource.Name() != "dragontable.jpg" {
				x.Hide()
			}
		}
	}
}

func ShowMap(fileName string) {

	fmt.Println("Execute ShowMap")
	fmt.Println(len(mainContent.Objects))

	for _, child := range mainContent.Objects {
		switch x := child.(type) {
		case *canvas.Image:

			if fileName == x.Resource.Name() {
				x.Show()
			}
		}
	}
}

func BuildNavButtons() []*widget.Button {

	var navButtons []*widget.Button

	var touchControlButton, hamburgerButton, gridButton *widget.Button

	hamburger, hamburgerError := fyne.LoadResourceFromPath("./resources/icons/bars-solid.svg")
	if hamburgerError != nil {
		fmt.Println(hamburgerError)
	}

	disabledTouchIcon, enableTouchError := fyne.LoadResourceFromPath("./resources/icons/hand-point-up-regular.svg")
	if enableTouchError != nil {
		fmt.Println(enableTouchError)
	}

	enabledTouchIcon, enableTouchError := fyne.LoadResourceFromPath("./resources/icons/hand-point-up-solid.svg")
	if enableTouchError != nil {
		fmt.Println(enableTouchError)
	}

	gridIcon, gridError := fyne.LoadResourceFromPath("./resources/icons/border-all-solid.svg")
	if gridError != nil {
		fmt.Println(gridError)
	}

	touchControlButton = widget.NewButtonWithIcon("", enabledTouchIcon, func() {

		if touchControlButton.Icon == disabledTouchIcon {
			touchControlButton.Importance = widget.HighImportance
			touchControlButton.SetIcon(enabledTouchIcon)
		} else {
			touchControlButton.Importance = widget.MediumImportance
			touchControlButton.SetIcon(disabledTouchIcon)
		}
	})

	touchControlButton.Importance = widget.HighImportance
	touchControlButton.Resize(fyne.NewSize(50, 50))
	touchControlButton.Move(fyne.Position{X: float32(ScreenWidth) - 130, Y: 10})

	hamburgerButton = widget.NewButtonWithIcon("", hamburger, func() {
		if hamburgerButton.Importance == widget.HighImportance {
			hamburgerButton.Importance = widget.MediumImportance
		} else {
			hamburgerButton.Importance = widget.HighImportance
		}

		for _, child := range mainContent.Objects {
			switch x := child.(type) {
			case *widget.List:

				if x.Hidden {
					x.Show()
				} else {
					x.Hide()
				}
			}
		}
	})

	hamburgerButton.Importance = widget.HighImportance
	hamburgerButton.Resize(fyne.NewSize(50, 50))
	hamburgerButton.Move(fyne.Position{X: float32(ScreenWidth) - 70, Y: 10})

	gridButton = widget.NewButtonWithIcon("", gridIcon, func() {

		if gridButton.Importance == widget.HighImportance {
			gridButton.Importance = widget.MediumImportance
		} else {
			gridButton.Importance = widget.HighImportance
		}

		for _, child := range mainContent.Objects {
			switch x := child.(type) {
			case *canvas.Line:
				if x.Hidden {
					x.Show()
				} else {
					x.Hide()
				}
			}
		}
	})

	gridButton.Importance = widget.MediumImportance
	gridButton.Resize(fyne.NewSize(50, 50))
	gridButton.Move(fyne.Position{X: float32(ScreenWidth) - 190, Y: 10})

	navButtons = append(navButtons, touchControlButton, hamburgerButton, gridButton)

	return navButtons
}

func DrawGrid() []*canvas.Line {

	var lines []*canvas.Line
	screenGridOffset := 5

	vLineSpace := ScreenWidth / ScreenDimensionWidth
	hLineSpace := ScreenHeight / ScreenDimensionHeight

	for i := 0; i < int(ScreenDimensionWidth); i++ {
		line := canvas.NewLine(color.RGBA{R: 56, G: 56, B: 56, A: 255})
		line.StrokeWidth = 1
		line.Position1 = fyne.NewPos(float32(vLineSpace*i), float32(0-screenGridOffset))
		line.Position2 = fyne.NewPos(float32(vLineSpace*i), float32(ScreenHeight+screenGridOffset))
		line.Hide()

		lines = append(lines, line)
	}

	for i := 0; i < int(ScreenDimensionHeight); i++ {
		line := canvas.NewLine(color.RGBA{R: 56, G: 56, B: 56, A: 255})
		line.StrokeWidth = 1
		line.Position1 = fyne.NewPos(float32(0-screenGridOffset), float32(hLineSpace*i))
		line.Position2 = fyne.NewPos(float32(ScreenWidth+screenGridOffset), float32(hLineSpace*i))
		line.Hide()

		lines = append(lines, line)
	}

	return lines
}
