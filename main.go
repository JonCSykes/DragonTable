package main

/*
#cgo CFLAGS: -I../lib
#cgo LDFLAGS: -L. -lEloMtApi
#include <EloInterface.h>
*/
import "C"
import (
	"fmt"
	"image/color"
	"math"
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

const ScreenDimensionWidth int = 30
const ScreenDimensionHeight int = 16
const ZoomSliderWidth float32 = 150
const ZoomSliderHeight float32 = 50
const ZoomSliderXOffset int = 200
const ZoomSliderYOffset int = 150

const DragonTableWallpaperPath string = "./resources/images/dragontable.jpg"

var ScreenHeight int
var ScreenWidth int
var TouchEnabled bool

var MainWindow fyne.Window
var mapFiles []*mapFile.MapFile
var mainContent *fyne.Container
var CurrentMap *canvas.Image
var CurrentMapSize fyne.Size
var MapContent *fyne.Container
var MapControl *container.Scroll
var ZoomControl *fyne.Container
var ZoomSlider *widget.Slider

type DeltaXY struct {
	TX int64
	TY int64
	DX int64
	DY int64
}

func main() {

	deltaChan := make(chan DeltaXY)
	TouchEnabled = true

	go streamTouchInput(deltaChan)
	go triggerScrolledEvent(deltaChan)

	GetScreenResolution()

	myApp := app.New()
	MainWindow = myApp.NewWindow("Dragon Table - v0.1")

	BuildUI()

	MainWindow.SetContent(mainContent)

	MainWindow.SetPadded(true)
	MainWindow.SetFullScreen(true)
	MainWindow.ShowAndRun()
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
	InitCurrentMap()

	BuildZoomControls()
	lines := DrawGrid()
	MapContent = container.NewWithoutLayout()
	MapContent.Add(CurrentMap)
	for _, line := range lines {
		MapContent.Add(line)
	}

	MapControl = container.NewScroll(MapContent)
	MapControl.Resize(fyne.NewSize(float32(ScreenWidth), float32(ScreenHeight)))
	MapControl.Move(fyne.Position{X: -2, Y: -2})

	content.Add(wallpaper)
	content.Add(MapControl)

	for _, navButton := range navButtons {
		content.Add(navButton)
	}

	mapList.Refresh()
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
	dragonTableImage.Move(fyne.Position{X: -4, Y: -4})

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
			o.(*widgetExt.ImageButton).Resize(fyne.Size{Width: 250, Height: 50})
			o.(*widgetExt.ImageButton).SetImage(mapFiles[i].ThumbResource)
		})

	mapList.Resize(fyne.NewSize(250, 1000))
	mapList.Move(fyne.Position{X: float32(ScreenWidth) - 260, Y: 80})

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
		SetZoomSliderRange()
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

	lines := DrawGrid()
	MapContent = container.NewWithoutLayout()
	MapContent.Add(CurrentMap)
	for _, line := range lines {
		MapContent.Add(line)
	}

	MapControl.Content = MapContent
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

	syncIcon, syncError := fyne.LoadResourceFromPath("./resources/icons/sync-alt-solid.svg")
	if syncError != nil {
		fmt.Println(syncError)
	}

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

	syncButton := widget.NewButtonWithIcon("", syncIcon, func() {
		fmt.Println("Refreshing Main Content")
		mapFiles = nil
		BuildUI()
		MainWindow.SetContent(mainContent)
	})

	syncButton.Importance = widget.HighImportance
	syncButton.Resize(fyne.NewSize(50, 50))
	syncButton.Move(fyne.Position{X: float32(ScreenWidth) - 130, Y: 10})

	touchControlButton = widget.NewButtonWithIcon("", enabledTouchIcon, func() {

		if touchControlButton.Icon == disabledTouchIcon {
			touchControlButton.Importance = widget.HighImportance
			touchControlButton.SetIcon(enabledTouchIcon)
			TouchEnabled = true
		} else {
			touchControlButton.Importance = widget.MediumImportance
			touchControlButton.SetIcon(disabledTouchIcon)
			TouchEnabled = false
		}
	})

	touchControlButton.Importance = widget.HighImportance
	touchControlButton.Resize(fyne.NewSize(50, 50))
	touchControlButton.Move(fyne.Position{X: float32(ScreenWidth) - 190, Y: 10})

	gridButton = widget.NewButtonWithIcon("", gridIcon, func() {
		for _, child := range MapContent.Objects {
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
	gridButton.Move(fyne.Position{X: float32(ScreenWidth) - 240, Y: 10})

	navButtons = append(navButtons, touchControlButton, hamburgerButton, gridButton, syncButton)

	return navButtons
}

func DrawGrid() []*canvas.Line {

	var lines []*canvas.Line
	screenGridOffset := 5

	vLineSpace := ScreenWidth / ScreenDimensionWidth
	hLineSpace := ScreenHeight / ScreenDimensionHeight

	fmt.Println("Number of lines across: ", int(CurrentMapSize.Width/float32(vLineSpace)))
	fmt.Println("Number of lines down: ", int(CurrentMapSize.Height/float32(hLineSpace)))

	for i := 0; i < int(CurrentMapSize.Width*float32(ZoomSlider.Max)/float32(vLineSpace)); i++ {
		line := canvas.NewLine(color.RGBA{R: 56, G: 56, B: 56, A: 255})
		line.StrokeWidth = 1
		line.Position1 = fyne.NewPos(float32(vLineSpace*i), float32(0-screenGridOffset))
		line.Position2 = fyne.NewPos(float32(vLineSpace*i), CurrentMapSize.Height*float32(ZoomSlider.Max)+float32(screenGridOffset))
		line.Hide()

		lines = append(lines, line)
	}

	for i := 0; i < int(CurrentMapSize.Height*float32(ZoomSlider.Max)/float32(hLineSpace)); i++ {
		line := canvas.NewLine(color.RGBA{R: 56, G: 56, B: 56, A: 255})
		line.StrokeWidth = 1
		line.Position1 = fyne.NewPos(float32(0-screenGridOffset), float32(hLineSpace*i))
		line.Position2 = fyne.NewPos(CurrentMapSize.Width*float32(ZoomSlider.Max)+float32(screenGridOffset), float32(hLineSpace*i))
		line.Hide()

		lines = append(lines, line)
	}

	return lines
}

func BuildZoomControls() {

	f := 1.0
	data := binding.BindFloat(&f)
	ZoomSlider = widget.NewSliderWithData(0.1, 2, data)
	ZoomSlider.Step = 0.1
	ZoomSlider.Resize(fyne.NewSize(ZoomSliderWidth, ZoomSliderHeight))
	ZoomSlider.OnChanged = func(value float64) {
		fmt.Println("Zoom Changed " + fmt.Sprintf("%f", value))
		if CurrentMap != nil {
			newWidth := float32(math.Abs(float64(CurrentMapSize.Width) * value))
			newHeight := float32(math.Abs(float64(CurrentMapSize.Height) * value))

			CurrentMap.Resize(fyne.NewSize(newWidth, newHeight))
			CurrentMap.SetMinSize(fyne.NewSize(newWidth, newHeight))
			fmt.Println(CurrentMapSize.Width, CurrentMapSize.Height)
			fmt.Println(newWidth, newHeight)
			MapControl.Refresh()
		}
	}
	ZoomControl = container.NewWithoutLayout(ZoomSlider)
	ZoomControl.Resize(fyne.NewSize(ZoomSliderWidth, ZoomSliderHeight))
	ZoomControl.Move(fyne.NewPos(float32(ScreenWidth-ZoomSliderXOffset), float32(ScreenHeight-ZoomSliderYOffset)))
	ZoomControl.Hide()

}

func SetZoomSliderRange() {

	heightRatio := float32(ScreenHeight) / CurrentMapSize.Height
	widthRatio := float32(ScreenWidth) / CurrentMapSize.Width

	if heightRatio > widthRatio {
		ZoomSlider.Min = float64(math.Round(float64(heightRatio)*100) / 100)
		ZoomSlider.Max = 2
		ZoomSlider.Value = 1
		fmt.Println(ZoomSlider.Min)
	} else {
		ZoomSlider.Min = float64(math.Round(float64(widthRatio)*100) / 100)
		ZoomSlider.Max = 2
		ZoomSlider.Value = 1
		fmt.Println(ZoomSlider.Min)
	}
	ZoomControl.Refresh()
}

func triggerScrolledEvent(deltaChan chan DeltaXY) {

	for {
		delta := <-deltaChan
		if MapControl != nil && TouchEnabled {
			//fmt.Println(delta.DX*10, delta.DY*10)
			MapControl.Scrolled(&fyne.ScrollEvent{PointEvent: fyne.PointEvent{AbsolutePosition: fyne.NewPos(0, 0), Position: fyne.NewPos(0, 0)}, Scrolled: fyne.NewDelta(float32(delta.DX/2), float32(delta.DY/2))})
		}
	}
}

func streamTouchInput(deltaChan chan DeltaXY) {

	var x, px, y, py, z C.int
	var status C.TOUCH_STATUS
	var is C.bool
	is = false

	for {
		C.EloGetTouchPacket(C.int(0), &x, &y, &z, &status, is)

		if status == 2 && px > 0 && py > 0 {
			deltaChan <- DeltaXY{TX: int64(x), TY: int64(y), DX: int64(x - px), DY: int64(y - py)}
		}

		px = x
		py = y
	}

}
