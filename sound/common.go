package sound

const GetAllRecordsQuery = "SELECT * FROM record"

type Record struct {
	Id          int
	Name        string
	SampleCount uint32
	RawData     []byte
}
