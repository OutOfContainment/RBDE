package gui

import (
	"fyne.io/fyne/v2/theme"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/OutOfContainment/RBDE/sound"
)

const (
	idleState      = "idle"
	recordingState = "recording"
	pauseState     = "pause"

	tracksLimit = 10
)

var state string = idleState

func buttons(
	RBDE fyne.App,
	win fyne.Window,
	sound *sound.Sound,
	currentTracksAmount int,
	tracksInterface *canvas.Text) *fyne.Container {

	playButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(),
		func() {
			log.Println("Play Button Pressed")
			sound.Play(currentTrack)
		})

	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(),
		func() {
			log.Println("Stop Button Pressed")
			sound.Stop()
			state = idleState
		})

	recordButton := widget.NewButtonWithIcon("", theme.MediaRecordIcon(),
		func() {
			log.Println("Record Button Pressed")
			sound.Record()
			if state == idleState {
				if currentTrack <= tracksLimit {
					if currentTrack == currentTracksAmount {
						currentTrack++
					}
					currentTracksAmount++
					screenCountUpdate(currentTrack, currentTracksAmount, tracksInterface)
					log.Println(currentTrack, " / ", currentTracksAmount)
				}
			}
			state = recordingState
		})

	pauseButton := widget.NewButtonWithIcon("", theme.MediaPauseIcon(),
		func() {
			log.Println("Pause Button Pressed")
			sound.Pause()
			state = pauseState
		})

	prevTrackButton := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(),
		func() {
			log.Println("Previous media button pressed")
			if currentTrack > 1 && state == idleState {
				currentTrack--
				screenCountUpdate(currentTrack, currentTracksAmount, tracksInterface)
			}
		})

	nextTrackButton := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(),
		func() {
			log.Println("Next mexia button pressed")
			if currentTrack < currentTracksAmount && state == idleState {
				currentTrack++
				screenCountUpdate(currentTrack, currentTracksAmount, tracksInterface)
			}
		})

	deleteButton := widget.NewButtonWithIcon("Delete all media", theme.DeleteIcon(),
		func() {
			log.Println("Delete button pressed")
			sound.ClearRecords()
			currentTrack, currentTracksAmount = 0, 0
			screenCountUpdate(currentTrack, currentTracksAmount, tracksInterface)
		})
	/*
		menuButton := widget.NewButton("Menu",
			func() {
				log.Println("Menu Button Pressed")
			})

	*/

	buttonLayout := container.NewCenter(container.NewVBox(container.NewHBox(
		prevTrackButton,
		stopButton,
		playButton,
		recordButton,
		pauseButton,
		nextTrackButton,
	),
		container.NewPadded( /* menuButton, layout.NewSpacer(),*/ deleteButton)))
	return buttonLayout
}
