package main

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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/OutOfContainment/ReallyBadDictaphoneEmulator/sound"
)

const (
	idleState      = "idle"
	recordingState = "recording"
	pauseState     = "pause"
)

var (
	state = idleState

	currentTrack        int
	currentTracksAmount int
)

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

func skeleton(RBDE fyne.App, win fyne.Window, sound *sound.Sound) *fyne.Container {
	tracksInterface := canvas.NewText(fmt.Sprintf("%d / %d", currentTrack, currentTracksAmount), color.White)
	tracksInterface.TextSize = 35

	tracksIcon := container.NewGridWithColumns(2, widget.NewLabel(""), widget.NewIcon(theme.StorageIcon()))

	// add image
	logopath, err := filepath.Abs("images/Logo.png")
	wallpaperpath, err := filepath.Abs("images/wallpaper.jpg")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(logopath)

	logo := canvas.NewImageFromFile(logopath)
	wallpaper := canvas.NewImageFromFile(wallpaperpath)
	logo.FillMode = canvas.ImageFillOriginal
	title := container.New(layout.NewCenterLayout(), logo)

	// Clock timer at the top
	currTime := widget.NewLabel("")
	go clockUpdate_go(currTime)
	bar := container.New(
		layout.NewCenterLayout(),
		canvas.NewRectangle(color.NRGBA{R: 0, G: 0, B: 0, A: 55}),
		currTime,
		widget.NewSeparator(),
	)

	screen := container.NewPadded(
		//		canvas.NewRectangle(color.NRGBA{R: 127, G: 20, B: 60, A: 155}),
		wallpaper,
		container.NewBorder(
			bar,
			nil,
			nil,
			container.NewPadded(container.NewGridWithColumns(
				2,
				tracksIcon,
				tracksInterface,
			)),
			nil,
		),
	)

	buttons := container.NewCenter(
		container.NewVBox(
			container.NewHBox(
				widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() {
					log.Println("Previous Media Button Pressed")
					if currentTrack > 1 && state == idleState {
						currentTrack--
						screenCountUpdate(currentTrack, currentTracksAmount, tracksInterface)
					}
				}),
				widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {
					log.Println("Stop Button Pressed")
					sound.Stop()
					state = idleState
				}),
				widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
					log.Println("Play Button Pressed")
					sound.Play(currentTrack)
				}),
				widget.NewButtonWithIcon("", theme.MediaRecordIcon(), func() {
					log.Println("Record Button Pressed")
					sound.Record()
					if state == idleState {
						if currentTrack <= 10 {
							if currentTrack == currentTracksAmount {
								currentTrack++
							}
							currentTracksAmount++
							screenCountUpdate(currentTrack, currentTracksAmount, tracksInterface)
							log.Println(currentTrack, " / ", currentTracksAmount)
						}
					}
					state = recordingState
				}),
				widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
					log.Println("Pause Button Pressed")
					sound.Pause()
					state = pauseState
				}),
				widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() {
					log.Println("Previous Next Button Pressed")
					if currentTrack < currentTracksAmount && state == idleState {
						currentTrack++
						screenCountUpdate(currentTrack, currentTracksAmount, tracksInterface)
					}
				}),
			),
			/*
				container.NewHBox(
					widget.NewButton("Menu", func() {
						log.Println("Menu Button Pressed")
					}),
					layout.NewSpacer(),
					widget.NewButton("Delete", func() {
						log.Println("Delete Button Pressed")
					}),
				),
			*/
		),
	)

	// GUI layout
	skeleton := container.NewGridWithRows(
		2,
		screen,
		container.NewVBox(
			title,
			buttons,
		),
	)
	return skeleton
}
