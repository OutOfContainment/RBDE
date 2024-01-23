package sound

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	idleState      = "idle"
	recordingState = "recording"
	playingState   = "playing"
)

type Sound struct {
	recorder *recorder
	player   *player

	lastRecordName string

	stopChans map[string]chan struct{}
	mu        sync.Mutex

	state string
}

func NewSound(db *sql.DB) *Sound {
	stopChans := make(map[string]chan struct{}, 2)
	stopChans[recordingState] = make(chan struct{})
	stopChans[playingState] = make(chan struct{})

	return &Sound{
		recorder: newRecorder(stopChans[recordingState], db),
		player:   newPlayer(stopChans[playingState], db),

		stopChans: stopChans,
		mu:        sync.Mutex{},

		state: idleState,
	}
}

func (d *Sound) Record() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.state == recordingState {
		return nil
	}

	d.state = recordingState

	var err error
	go func() {
		recordName := fmt.Sprintf("%d_record", time.Now().UnixMilli())
		if recordErr := d.recorder.Record(recordName); err != nil {
			err = recordErr
		} else {
			d.lastRecordName = recordName
		}
	}()

	return err
}

func (d *Sound) Play() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.state == playingState {
		return
	}

	d.state = playingState

	go d.player.Play(d.lastRecordName)
}

func (d *Sound) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch d.state {
	case idleState:
		log.Println("Nothing to stop.")
	case recordingState:
		go func() {
			d.stopChans[recordingState] <- struct{}{}
		}()
	case playingState:
		go func() {
			d.stopChans[playingState] <- struct{}{}
		}()
	default:
		log.Println("Unknown state.")
	}

	d.state = idleState
}

func (d *Sound) Pause() {
	d.mu.Lock()
	defer d.mu.Unlock()
	switch d.state {
	case idleState:
		log.Println("Nothing to pause.")
	case recordingState:
		go func() {
			d.recorder.Pause()
		}()
	case playingState:
		go func() {
			d.player.Pause()
		}()
	default:
		log.Println("Unknown state.")
	}
}
