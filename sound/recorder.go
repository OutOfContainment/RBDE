package sound

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gen2brain/malgo"
)

const insertRecordQuery = "INSERT INTO record (name, sample_count ,wav_data) VALUES (?, ?, ?)"

type recorder struct {
	stopChan chan struct{}
	paused   bool
	db       *sql.DB
}

func newRecorder(stopChan chan struct{}, db *sql.DB) *recorder {
	return &recorder{stopChan: stopChan, db: db}
}

func (r *recorder) Record(name string) error {
	if r.paused {
		r.paused = false
		return nil
	}
	log.Println("Recording started.")

	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("LOG <%v>\n", message)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Duplex)
	deviceConfig.Capture.Format = malgo.FormatS16
	deviceConfig.Capture.Channels = 1
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = 1
	deviceConfig.SampleRate = 44100
	deviceConfig.Alsa.NoMMap = 1

	// var playbackSampleCount uint32
	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Capture.Format))
	onRecvFrames := func(pSample2, pSample []byte, framecount uint32) {
		if r.paused {
			return
		}

		sampleCount := framecount * deviceConfig.Capture.Channels * sizeInBytes

		newCapturedSampleCount := capturedSampleCount + sampleCount

		pCapturedSamples = append(pCapturedSamples, pSample...)

		capturedSampleCount = newCapturedSampleCount

	}

	fmt.Println("Recording...")
	captureCallbacks := malgo.DeviceCallbacks{
		Data: onRecvFrames,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, captureCallbacks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = device.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		<-r.stopChan
		log.Println("Record stopped.")
		device.Uninit()
		wg.Done()
	}()

	wg.Wait()

	// just for the sake of it, to keep things JUST a bit cleannnnnn
	record := Record{
		Name:        name,
		SampleCount: capturedSampleCount,
		RawData:     pCapturedSamples,
	}

	insertNewRecordStatement, err := r.db.Prepare(insertRecordQuery)
	if err != nil {
		log.Println("statement ", err)
		return err
	}

	_, err = insertNewRecordStatement.Exec(record.Name, record.SampleCount, record.RawData)
	if err != nil {
		return err
	}

	log.Println(record.Name, " ", record.SampleCount, " ", len(record.RawData))

	return nil
}

func (r *recorder) Pause() {
	r.paused = true
}
