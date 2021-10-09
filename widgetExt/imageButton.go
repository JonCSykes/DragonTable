package widgetExt

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ImageButton widget has a text label and a background image that triggers an event func when clicked
type ImageButton struct {
	widget.DisableableWidget
	Text          string
	Image         fyne.Resource
	Importance    widget.ButtonImportance
	Alignment     widget.ButtonAlign
	IconPlacement widget.ButtonIconPlacement

	OnTapped func() `json:"-"`

	hovered, focused bool
	tapAnim          *fyne.Animation
}

type imageButtonRenderer struct {
	fyne.WidgetRenderer

	objects     []fyne.CanvasObject
	image       *canvas.Image
	label       *widget.RichText
	background  *canvas.Rectangle
	tapBG       *canvas.Rectangle
	imageButton *ImageButton
	layout      fyne.Layout
}

func NewImageButton(text string, image fyne.Resource, tapped func()) *ImageButton {
	imageButton := &ImageButton{
		Text:     text,
		Image:    image,
		OnTapped: tapped,
	}
	imageButton.ExtendBaseWidget(imageButton)

	return imageButton
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (imageButton *ImageButton) CreateRenderer() fyne.WidgetRenderer {

	imageButton.ExtendBaseWidget(imageButton)

	seg := &widget.TextSegment{Text: imageButton.Text, Style: widget.RichTextStyleStrong}
	seg.Style.Alignment = fyne.TextAlignCenter

	text := widget.NewRichText(seg)
	//text.inset = fyne.NewSize(theme.Padding()*2, theme.Padding()*2)

	background := canvas.NewRectangle(theme.ButtonColor())

	tapBG := canvas.NewRectangle(color.Transparent)

	imageButton.tapAnim = newButtonTapAnimation(tapBG, imageButton)
	imageButton.tapAnim.Curve = fyne.AnimationEaseOut

	renderer := &imageButtonRenderer{

		background:  background,
		tapBG:       tapBG,
		imageButton: imageButton,
		label:       text,
		layout:      layout.NewHBoxLayout(),
	}
	renderer.updateImageAndText()
	renderer.applyTheme()
	return renderer
}

// Cursor returns the cursor type of this widget
func (imageButton *ImageButton) Cursor() desktop.Cursor {
	return desktop.DefaultCursor
}

// FocusGained is a hook called by the focus handling logic after this object gained the focus.
func (imageButton *ImageButton) FocusGained() {
	imageButton.focused = true
	imageButton.Refresh()
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (imageButton *ImageButton) FocusLost() {
	imageButton.focused = false
	imageButton.Refresh()
}

// MinSize returns the size that this widget should not shrink below
func (imageButton *ImageButton) MinSize() fyne.Size {
	imageButton.ExtendBaseWidget(imageButton)
	return imageButton.BaseWidget.MinSize()
}

// MouseIn is called when a desktop pointer enters the widget
func (imageButton *ImageButton) MouseIn(*desktop.MouseEvent) {
	imageButton.hovered = true
	imageButton.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (imageButton *ImageButton) MouseMoved(*desktop.MouseEvent) {
}

// MouseOut is called when a desktop pointer exits the widget
func (imageButton *ImageButton) MouseOut() {
	imageButton.hovered = false
	imageButton.Refresh()
}

// SetImage updates the image on the image button - pass nil to hide the image.
func (imageButton *ImageButton) SetImage(image fyne.Resource) {
	imageButton.Image = image

	imageButton.Refresh()
}

// SetText allows the button label to be changed
func (imageButton *ImageButton) SetText(text string) {
	imageButton.Text = text

	imageButton.Refresh()
}

// Tapped is called when a pointer tapped event is captured and triggers any tap handler
func (imageButton *ImageButton) Tapped(*fyne.PointEvent) {
	if imageButton.Disabled() {
		return
	}

	imageButton.tapAnimation()
	imageButton.Refresh()

	if imageButton.OnTapped != nil {
		imageButton.OnTapped()
	}
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (imageButton *ImageButton) TypedRune(rune) {
}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (imageButton *ImageButton) TypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeySpace {
		imageButton.Tapped(nil)
	}
}

func (imageButton *ImageButton) tapAnimation() {
	if imageButton.tapAnim == nil {
		return
	}
	imageButton.tapAnim.Stop()
	imageButton.tapAnim.Start()
}

func (r *imageButtonRenderer) Destroy() {
}

func (r *imageButtonRenderer) Refresh() {

	r.label.Segments[0].(*widget.TextSegment).Text = r.imageButton.Text
	r.updateImageAndText()
	r.applyTheme()
	r.background.Refresh()
	r.Layout(r.imageButton.Size())
	canvas.Refresh(r.imageButton)
}

func (renderer *imageButtonRenderer) updateImageAndText() {
	if renderer.imageButton.Image != nil && renderer.imageButton.Visible() {
		if renderer.image == nil {
			renderer.image = canvas.NewImageFromResource(renderer.imageButton.Image)
			renderer.image.Translucency = float64(0.5)
			renderer.image.Resize(renderer.imageButton.Size())
			renderer.image.FillMode = canvas.ImageFillOriginal
			renderer.SetObjects([]fyne.CanvasObject{renderer.background, renderer.tapBG, renderer.image, renderer.label})
		}
		if renderer.imageButton.Disabled() {
			renderer.image.Resource = theme.NewDisabledResource(renderer.imageButton.Image)
		} else {
			renderer.image.Resource = renderer.imageButton.Image
		}
		renderer.image.Refresh()
		renderer.image.Show()
	} else if renderer.image != nil {
		renderer.image.Hide()
	}
	if renderer.imageButton.Text == "" {
		renderer.label.Hide()
	} else {
		renderer.label.Show()
	}
	renderer.label.Refresh()
}

func (renderer *imageButtonRenderer) SetObjects(objects []fyne.CanvasObject) {
	renderer.objects = objects
}

func (renderer *imageButtonRenderer) Objects() []fyne.CanvasObject {
	return renderer.objects
}

func (renderer *imageButtonRenderer) applyTheme() {
	renderer.background.FillColor = renderer.buttonColor()
	renderer.label.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNameForeground
	if renderer.imageButton.Disabled() {
		renderer.label.Segments[0].(*widget.TextSegment).Style.ColorName = theme.ColorNameDisabled
	}
	renderer.label.Refresh()
	if renderer.image != nil && renderer.image.Resource != nil {
		switch res := renderer.image.Resource.(type) {
		case *theme.ThemedResource:

			renderer.image.Resource = theme.NewInvertedThemedResource(res)
			renderer.image.Refresh()

		case *theme.InvertedThemedResource:

			renderer.image.Resource = res.Original()
			renderer.image.Refresh()

		}
	}
}

func (renderer *imageButtonRenderer) buttonColor() color.Color {
	switch {
	case renderer.imageButton.Disabled():
		return theme.DisabledButtonColor()
	case renderer.imageButton.focused:
		return blendColor(theme.ButtonColor(), theme.FocusColor())
	case renderer.imageButton.hovered:
		bg := theme.ButtonColor()
		return blendColor(bg, theme.HoverColor())
	default:
		return theme.ButtonColor()
	}
}

// Layout the components of the button widget
func (r *imageButtonRenderer) Layout(size fyne.Size) {
	var inset fyne.Position
	bgSize := size

	r.background.Move(inset)
	r.background.Resize(bgSize)

	hasImage := r.image != nil
	hasLabel := r.label.Segments[0].(*widget.TextSegment).Text != ""
	if !hasImage && !hasLabel {
		// Nothing to layout
		return
	}
	imageSize := bgSize
	labelSize := r.label.MinSize()
	padding := r.padding()
	if hasLabel {
		if hasImage {
			// Both
			var objects []fyne.CanvasObject

			objects = append(objects, r.image, r.label)

			r.image.SetMinSize(imageSize)
			min := r.layout.MinSize(objects)
			r.layout.Layout(objects, min)

			r.label.Move(alignedPosition(r.imageButton.Alignment, padding, labelSize, size))
			r.image.Move(r.image.Position())
		} else {
			// Label Only
			r.label.Move(alignedPosition(r.imageButton.Alignment, padding, labelSize, size))
			r.label.Resize(labelSize)
		}
	} else {
		// Image Only
		r.image.Move(alignedPosition(r.imageButton.Alignment, padding, imageSize, size))
		r.image.Resize(imageSize)
	}
}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (r *imageButtonRenderer) MinSize() (size fyne.Size) {
	hasIcon := r.image != nil
	hasLabel := r.label.Segments[0].(*widget.TextSegment).Text != ""
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	labelSize := r.label.MinSize()
	if hasLabel {
		size.Width = labelSize.Width
	}
	if hasIcon {
		if hasLabel {
			size.Width += theme.Padding()
		}
		size.Width += iconSize.Width
	}
	size.Height = fyne.Max(labelSize.Height, iconSize.Height)
	size = size.Add(r.padding())
	return
}

func (r *imageButtonRenderer) padding() fyne.Size {
	if r.imageButton.Text == "" {
		return fyne.NewSize(theme.Padding()*4, theme.Padding()*4)
	}
	return fyne.NewSize(theme.Padding()*6, theme.Padding()*4)
}

func newButtonTapAnimation(bg *canvas.Rectangle, w fyne.Widget) *fyne.Animation {
	return fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
		mid := (w.Size().Width - theme.Padding()) / 2
		size := mid * done
		bg.Resize(fyne.NewSize(size*2, w.Size().Height-theme.Padding()))
		bg.Move(fyne.NewPos(mid-size, theme.Padding()/2))

		r, g, bb, a := ToNRGBA(theme.PressedColor())
		aa := uint8(a)
		fade := aa - uint8(float32(aa)*done)
		bg.FillColor = &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(bb), A: fade}
		canvas.Refresh(bg)
	})
}

func alignedPosition(align widget.ButtonAlign, padding, objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	switch align {
	case widget.ButtonAlignCenter:
		pos.X = (layoutSize.Width - objectSize.Width) / 2
	case widget.ButtonAlignLeading:
		pos.X = padding.Width / 2
	case widget.ButtonAlignTrailing:
		pos.X = layoutSize.Width - objectSize.Width - padding.Width/2
	}
	return
}

func blendColor(under, over color.Color) color.Color {
	// This alpha blends with the over operator, and accounts for RGBA() returning alpha-premultiplied values
	dstR, dstG, dstB, dstA := under.RGBA()
	srcR, srcG, srcB, srcA := over.RGBA()

	srcAlpha := float32(srcA) / 0xFFFF
	dstAlpha := float32(dstA) / 0xFFFF

	outAlpha := srcAlpha + dstAlpha*(1-srcAlpha)
	outR := srcR + uint32(float32(dstR)*(1-srcAlpha))
	outG := srcG + uint32(float32(dstG)*(1-srcAlpha))
	outB := srcB + uint32(float32(dstB)*(1-srcAlpha))
	// We create an RGBA64 here because the color components are already alpha-premultiplied 16-bit values (they're just stored in uint32s).
	return color.RGBA64{R: uint16(outR), G: uint16(outG), B: uint16(outB), A: uint16(outAlpha * 0xFFFF)}

}
