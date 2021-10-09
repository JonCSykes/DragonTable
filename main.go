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

func main() {

	GetScreenResolution()

	myApp := app.New()
	myWindow := myApp.NewWindow("Dragon Table")

	myWindow.SetPadded(true)
	myWindow.SetFullScreen(true)
	myWindow.SetContent(BuildUI())
	myWindow.ShowAndRun()
}

func GetScreenResolution() {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	ScreenWidth = int(win.GetDeviceCaps(hDC, win.HORZRES))
	ScreenHeight = int(win.GetDeviceCaps(hDC, win.VERTRES))

	fmt.Println(strconv.Itoa(ScreenWidth) + " x " + strconv.Itoa(ScreenHeight))
}

func BuildUI() *fyne.Container {

	image := BuildImageDisplay()
	mapList := BuildNavList()
	touchControlButton := BuildNavButtons()

	content := container.NewWithoutLayout(image, mapList, touchControlButton)

	return content
}

func BuildImageDisplay() *canvas.Image {
	image := canvas.NewImageFromFile("./resources/maps/campaign_khorvaire.jpg")
	image.Resize(fyne.NewSize(float32(ScreenWidth), float32(ScreenHeight)))
	image.FillMode = canvas.ImageFillStretch
	image.Move(fyne.Position{X: 0, Y: 0})

	return image
}

func BuildNavList() *widget.List {

	mapFiles := mapFile.GetMaps()

	mapList := widget.NewList(
		func() int {
			return len(mapFiles)
		},
		func() fyne.CanvasObject {
			return widgetExt.NewImageButton("test", nil, nil)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			mapThumb, thumbError := fyne.LoadResourceFromPath(mapFiles[i].FullThumbnailPath)
			if thumbError != nil {
				fmt.Println(thumbError)
			}
			//o.(*widget.Button).SetText(mapFiles[i].FileName)
			o.(*widgetExt.ImageButton).Resize(fyne.Size{Width: 400, Height: 100})
			o.(*widgetExt.ImageButton).SetImage(mapThumb)
		})

	mapList.Resize(fyne.NewSize(400, 1000))
	mapList.Move(fyne.Position{X: float32(ScreenWidth) - 410, Y: 10})

	return mapList
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
