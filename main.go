package main

// Functions that end with "_go" are meant to be Goroutines

import (
	"database/sql"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/OutOfContainment/ReallyBadDictaphoneEmulator/sound"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Start")

	createDB()
	records, err := sql.Open("sqlite3", "./records.db")
	if err != nil {
		log.Fatal("Could not open database ", err)
	}
	defer records.Close()
	createTable(records)

	sound := sound.NewSound(records)

	// Initialise window
	RBDE := app.New()
	win := RBDE.NewWindow("DiEmu")
	win.Resize(fyne.NewSize(250, 500))

	win.SetContent(skeleton(RBDE, win, sound))

	// Open window
	defer log.Println("Goodbye.")
	defer win.ShowAndRun()
	defer log.Println("Open window ##")
}

func createDB() {
	os.Remove("records.db")

	log.Println("Creating records.db ..")
	file, err := os.Create("records.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("records.db created")
}

func createTable(records *sql.DB) {
	createRecordsTableSQL := `CREATE TABLE IF NOT EXISTS record (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		sample_count INTEGER,
		wav_data BLOB
	);`

	log.Println("Create records table ..")
	statement, err := records.Prepare(createRecordsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Records table created")
}
