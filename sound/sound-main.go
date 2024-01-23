package sound

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	idleState = "idle"

	recordingState = "recording"
	recPausedState = "paused recording"

	playingState    = "playing"
	playPausedState = "paused playing"
)

type Sound struct {
	recorder *recorder
	player   *player

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

func (s *Sound) Record() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != idleState && s.state != recPausedState {
		return nil
	}

	s.state = recordingState

	var err error
	go func() {
		recordName := fmt.Sprintf("%d_record", time.Now().UnixMilli())
		if recordErr := s.recorder.Record(recordName); err != nil {
			err = recordErr
		}
	}()
	return err
}

func (s *Sound) Play(recordId int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state != idleState && s.state != playPausedState {
		return
	}

	s.state = playingState

	go func() {
		s.player.Play(recordId)
		s.state = idleState
	}()
}

func (s *Sound) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.state {
	case idleState:
		log.Println("Nothing to stop.")
	case playingState:
		go func() {
			s.stopChans[playingState] <- struct{}{}
		}()
	default:
		go func() {
			s.stopChans[recordingState] <- struct{}{}
		}()
	}

	s.state = idleState
}

func (s *Sound) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch s.state {
	case idleState:
		log.Println("Nothing to pause.")
	case recordingState:
		s.recorder.Pause()
		s.state = recPausedState
	case playingState:
		s.player.Pause()
		s.state = playPausedState
	default:
		log.Println("Unknown state.")
	}
}
