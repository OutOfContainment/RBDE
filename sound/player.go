package sound

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gen2brain/malgo"
)

const getRecordByIdQuery = `SELECT id, name, sample_count, wav_data FROM 
( SELECT ROW_NUMBER() OVER (ORDER BY id ASC) AS row_id, * FROM record ORDER BY row_id )
WHERE row_id = ?`

type player struct {
	stopChan chan struct{}
	paused   bool
	db       *sql.DB
}

func newPlayer(stopChan chan struct{}, db *sql.DB) *player {
	return &player{stopChan: stopChan, db: db}
}

func (p *player) Play(recordId int) error {
	if p.paused {
		p.paused = false
		return nil
	}

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		//
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	// TODO: query record by id or name
	getRecordByIdStatement, err := p.db.Prepare(getRecordByIdQuery)
	if err != nil {
		return err
	}

	record := Record{}
	row := getRecordByIdStatement.QueryRow(recordId)
	row.Scan(&record.Id, &record.Name, &record.SampleCount, &record.RawData)

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Duplex)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = 44100
	deviceConfig.Alsa.NoMMap = 1

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	var playbackSampleCount uint32
	done := make(chan struct{})

	onSendFrames := func(pSample, nil []byte, framecount uint32) {
		if p.paused {
			return
		}

		samplesToRead := framecount * deviceConfig.Playback.Channels * sizeInBytes
		if samplesToRead > record.SampleCount-playbackSampleCount {
			samplesToRead = record.SampleCount - playbackSampleCount
		}

		copy(pSample, record.RawData[playbackSampleCount:playbackSampleCount+samplesToRead])

		playbackSampleCount += samplesToRead

		if playbackSampleCount == uint32(len(record.RawData)) {
			done <- struct{}{}
		}
	}

	fmt.Println("Playing...")
	playbackCallbacks := malgo.DeviceCallbacks{
		Data: onSendFrames,
	}

	device, err := malgo.InitDevice(ctx.Context, deviceConfig, playbackCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	select {
	case <-done:
		log.Println("Record playing has been fully played.")
		break
	case <-p.stopChan:
		log.Println("Record playing has been stopped.")
		break
	}

	device.Uninit()

	return nil
}

func (p *player) Pause() {
	p.paused = true
}
