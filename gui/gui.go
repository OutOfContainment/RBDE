package gui

import (
	"fmt"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"log"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/OutOfContainment/RBDE/sound"
)

var currentTrack int

func Skeleton(
	RBDE fyne.App,
	win fyne.Window,
	sound *sound.Sound,
	currentTracksAmount int) *fyne.Container {

	tracksInterface := canvas.NewText(fmt.Sprintf("%d / %d",
		currentTrack, currentTracksAmount), color.White)
	tracksInterface.TextSize = 35

	tracksIcon := container.NewGridWithColumns(2,
		widget.NewLabel(""), widget.NewIcon(theme.StorageIcon()))

	// add images
	wallpaperpath, err := filepath.Abs("images/wallpaper.jpg")
	if err != nil {
		log.Fatal("Error getting wallpaper.png file", err)
	} else {
		log.Println(wallpaperpath)
	}
	logopath, err := filepath.Abs("images/Logo.png")
	if err != nil {
		log.Fatal("Error getting Logo.png file", err)
	} else {
		log.Println(logopath)
	}

	wallpaper := canvas.NewImageFromFile(wallpaperpath)

	logo := canvas.NewImageFromFile(logopath)
	logo.FillMode = canvas.ImageFillOriginal
	logoLayout := container.NewCenter(logo)

	// Clock timer at the top
	currTime := widget.NewLabel("")
	go clockUpdate_go(currTime)
	bar := container.NewCenter(currTime)

	screen := container.NewPadded(
		//		canvas.NewRectangle(color.NRGBA{R: 127, G: 20, B: 60, A: 155}),
		wallpaper,
		container.NewBorder(
			bar, // top
			nil, // bottom
			nil, // left
			container.NewPadded( // right
				container.NewGridWithColumns(2, tracksIcon, tracksInterface)),
			nil, //center
		),
	)

	// GUI layout
	Skeleton := container.NewGridWithRows(2, // grid amount
		screen,
		container.NewGridWithRows(2,
			logoLayout,
			buttons(RBDE, win, sound, currentTracksAmount, tracksInterface),
		),
	)
	return Skeleton
}

func clockUpdate_go(currTime *widget.Label) {
	for range time.Tick(time.Second) {
		time := time.Now().Format(time.TimeOnly)
		currTime.SetText(time)
	}

}

func screenCountUpdate(currentTrack, currentTracksAmount int, tracksInterface *canvas.Text) {
	tracksInterface.Text = fmt.Sprintf("%d / %d", currentTrack, currentTracksAmount)
	tracksInterface.Refresh()
}
