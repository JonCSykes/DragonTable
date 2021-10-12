package main

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
var CurrentMap *canvas.Image
var CurrentMapSize fyne.Size
var MapControl *container.Scroll
var ZoomControl *fyne.Container

func main() {

	GetScreenResolution()

	myApp := app.New()
	mainWindow := myApp.NewWindow("Dragon Table - v0.1")

	BuildUI()

	mainWindow.SetContent(mainContent)

	mainWindow.SetPadded(true)
	mainWindow.SetFullScreen(true)
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

	content := container.NewWithoutLayout()

	wallpaper := BuildWallpaper()
	mapList := BuildNavList()
	navButtons := BuildNavButtons()
	lines := DrawGrid()

	InitCurrentMap()
	BuildZoomControls()

	MapControl = container.NewScroll(CurrentMap)
	MapControl.Resize(fyne.NewSize(float32(ScreenWidth), float32(ScreenHeight)))
	MapControl.Move(fyne.Position{X: 0, Y: 0})

	content.Add(wallpaper)
	content.Add(MapControl)

	for _, line := range lines {
		content.Add(line)
	}

	for _, navButton := range navButtons {
		content.Add(navButton)
	}

	content.Add(mapList)

	if ZoomControl != nil {
		content.Add(ZoomControl)
	}

	mainContent = content
}

func BuildWallpaper() *canvas.Image {
	dragonTableWallpaperResource, imageError := fyne.LoadResourceFromPath(DragonTableWallpaperPath)
	if imageError != nil {
		fmt.Println(imageError)
	}

	dragonTableImage := canvas.NewImageFromResource(dragonTableWallpaperResource)
	dragonTableImage.Resize(fyne.NewSize(float32(ScreenWidth), float32(ScreenHeight)))
	dragonTableImage.FillMode = canvas.ImageFillContain
	dragonTableImage.Move(fyne.Position{X: 0, Y: 0})

	return dragonTableImage
}

func InitCurrentMap() {
	if len(mapFiles) == 0 {
		mapFiles = mapFile.GetMaps()
	}

	CurrentMap = mapFiles[0].Image
	CurrentMap.FillMode = canvas.ImageFillStretch
	CurrentMapSize = CurrentMap.Size()

	CurrentMap.Hide()
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
				fmt.Println("Width : " + strconv.Itoa(mapFiles[i].Width))
				fmt.Println("Height : " + strconv.Itoa(mapFiles[i].Height))

				if !CurrentMap.Hidden && CurrentMap.Resource.Name() == mapFiles[i].FileName+"."+mapFiles[i].Extension {
					HideCurrentMap()
				} else {
					SetCurrentMap(mapFiles[i].Image)
					ShowCurrentMap()
				}
			}
			o.(*widgetExt.ImageButton).Resize(fyne.Size{Width: 200, Height: 50})
			o.(*widgetExt.ImageButton).SetImage(mapFiles[i].ThumbResource)
		})

	mapList.Resize(fyne.NewSize(200, 1000))
	mapList.Move(fyne.Position{X: float32(ScreenWidth) - 210, Y: 80})

	return mapList
}

func HideCurrentMap() {
	if CurrentMap != nil {
		CurrentMap.Hide()
		ZoomControl.Hide()
	}
}

func ShowCurrentMap() {
	if CurrentMap != nil {
		CurrentMap.Show()
		ZoomControl.Show()
	}
}

func SetCurrentMap(image *canvas.Image) {
	CurrentMap = image
	CurrentMap.FillMode = canvas.ImageFillStretch
	CurrentMap.Move(fyne.Position{X: 0, Y: 0})
	CurrentMapSize = CurrentMap.Size()

	fmt.Println(image.Size().Width, image.Size().Height)
	fmt.Println(CurrentMap.Size().Width, CurrentMap.Size().Height)
	fmt.Println(CurrentMapSize.Width, CurrentMapSize.Height)

	MapControl.Content = CurrentMap
	MapControl.Refresh()

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
			MapControl.Scrolled(&fyne.ScrollEvent{PointEvent: fyne.PointEvent{AbsolutePosition: fyne.NewPos(0, 0), Position: fyne.NewPos(0, 0)}, Scrolled: fyne.NewDelta(45, 45)})
		} else {
			touchControlButton.Importance = widget.MediumImportance
			touchControlButton.SetIcon(disabledTouchIcon)
			MapControl.Scrolled(&fyne.ScrollEvent{PointEvent: fyne.PointEvent{AbsolutePosition: fyne.NewPos(0, 0), Position: fyne.NewPos(0, 0)}, Scrolled: fyne.NewDelta(245, 245)})
		}
	})

	touchControlButton.Importance = widget.HighImportance
	touchControlButton.Resize(fyne.NewSize(50, 50))
	touchControlButton.Move(fyne.Position{X: float32(ScreenWidth) - 130, Y: 10})

	hamburgerButton = widget.NewButtonWithIcon("", hamburger, func() {

		for _, child := range mainContent.Objects {
			switch x := child.(type) {
			case *widget.List:

				if x.Hidden {
					x.Show()
					hamburgerButton.Importance = widget.HighImportance
				} else {
					x.Hide()
					hamburgerButton.Importance = widget.MediumImportance
				}
			}
		}
	})

	hamburgerButton.Importance = widget.HighImportance
	hamburgerButton.Resize(fyne.NewSize(50, 50))
	hamburgerButton.Move(fyne.Position{X: float32(ScreenWidth) - 70, Y: 10})

	gridButton = widget.NewButtonWithIcon("", gridIcon, func() {

		for _, child := range mainContent.Objects {
			switch x := child.(type) {
			case *canvas.Line:
				if x.Hidden {
					x.Show()
					gridButton.Importance = widget.HighImportance
				} else {
					x.Hide()
					gridButton.Importance = widget.MediumImportance
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

func BuildZoomControls() {

	f := 1.0
	data := binding.BindFloat(&f)
	label := widget.NewLabelWithData(binding.FloatToStringWithFormat(data, "Zoom: %0f"))
	slide := widget.NewSliderWithData(1, 2, data)
	slide.Step = 0.1
	slide.Resize(fyne.NewSize(600, 50))
	slide.OnChanged = func(value float64) {
		fmt.Println("Zoom Changed " + fmt.Sprintf("%f", value))
		if CurrentMap != nil {
			newWidth := CurrentMapSize.Width * float32(value)
			newHeight := CurrentMapSize.Height * float32(value)
			CurrentMap.Resize(fyne.NewSize(newWidth, newHeight))
			CurrentMap.SetMinSize(fyne.NewSize(newWidth, newHeight))
			fmt.Println(CurrentMapSize.Width, CurrentMapSize.Height)
			fmt.Println(newWidth, newHeight)
		}
	}
	ZoomControl = container.NewVBox(label, slide)
	ZoomControl.Move(fyne.NewPos(float32(ScreenWidth-200), float32(ScreenHeight-200)))
	ZoomControl.Hide()

}
