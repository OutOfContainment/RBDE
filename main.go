package main

// Functions that end with "_go" are meant to be Goroutines

import (
	"database/sql"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/OutOfContainment/RBDE/gui"
	"github.com/OutOfContainment/RBDE/sound"
	_ "github.com/mattn/go-sqlite3"
)

var tracksAmount int

func main() {
	log.Println("Start")

	// set font
	if _, err := os.Stat("font.ttf"); err == nil {
		os.Setenv("FYNE_FONT", "font.ttf")
	} else {
		log.Println("'font.ttf' not found; using default font.")
	}

	// open existing database || create database
	records := getDatabase()
	defer records.Close()

	sound := sound.NewSound(records)

	// Initialise window
	RBDE := app.New()
	win := RBDE.NewWindow("DiEmu")
	win.Resize(fyne.NewSize(240, 400))

	win.SetContent(gui.Skeleton(RBDE, win, sound, tracksAmount))

	// Open window
	defer log.Println("Goodbye.")
	defer win.ShowAndRun()
	defer log.Println("Open window ##")
}

func getDatabase() *sql.DB {
	if _, err := os.Stat("records.db"); err == nil {
		records, err := sql.Open("sqlite3", "./records.db")
		if err != nil {
			log.Fatal("Could not open existing database ", err)
		}

		tracksAmount = getTracksAmount(records)
		log.Println(tracksAmount, "tracks in opened database")
		return records

	} else {
		createDB()
		records, err := sql.Open("sqlite3", "./records.db")
		if err != nil {
			log.Fatal("Could not open database ", err)
		}

		createTable(records)
		return records
	}
}

func createDB() {
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

func getTracksAmount(db *sql.DB) int {
	getTracksAmountQuery := "SELECT COUNT(id) FROM record"
	getTracksAmountStatement, err := db.Prepare(getTracksAmountQuery)
	if err != nil {
		log.Fatal("ErrorSQLOpen", err)
	}

	getTracksAmountStatement.QueryRow().Scan(&tracksAmount)

	return tracksAmount
}
