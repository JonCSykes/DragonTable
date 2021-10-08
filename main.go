package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/JonCSykes/DragonTable/mapFile"
	"github.com/lxn/win"
)

var touchControlButton *widget.Button

func main() {
	hDC := win.GetDC(0)
	defer win.ReleaseDC(0, hDC)
	screenWidth := int(win.GetDeviceCaps(hDC, win.HORZRES))
	screenHeight := int(win.GetDeviceCaps(hDC, win.VERTRES))

	myApp := app.New()
	myWindow := myApp.NewWindow("Dragon Table")

	disabledTouchIcon, enableTouchError := fyne.LoadResourceFromPath("./resources/icons/hand-point-up-regular.svg")
	if enableTouchError != nil {
		fmt.Println(enableTouchError)
	}

	enabledTouchIcon, enableTouchError := fyne.LoadResourceFromPath("./resources/icons/hand-point-up-solid.svg")
	if enableTouchError != nil {
		fmt.Println(enableTouchError)
	}

	mapFiles := mapFile.GetMaps()

	mapList := widget.NewList(
		func() int {
			return len(mapFiles)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Maps")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(mapFiles[i].FileName)
		})

	mapList.Resize(fyne.NewSize(500, 1000))
	mapList.Move(fyne.Position{X: float32(screenWidth) - 510, Y: 10})

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
	touchControlButton.Move(fyne.Position{X: float32(screenWidth) - 110, Y: float32(screenHeight) - 110})

	image := canvas.NewImageFromFile("./resources/maps/campaign_khorvaire.jpg")
	image.Resize(fyne.NewSize(float32(screenWidth), float32(screenHeight)))
	image.FillMode = canvas.ImageFillStretch
	image.Move(fyne.Position{X: 0, Y: 0})

	content := container.NewWithoutLayout(image, mapList, touchControlButton)

	myWindow.SetPadded(true)
	myWindow.SetFullScreen(true)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
