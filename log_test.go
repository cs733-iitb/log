package log

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	cleanup()
	m.Run()
	cleanup()
}

func TestAppend(t *testing.T) {
	append(t, 0)
	append(t, 50)
}

// append 50 records of the form (i ->  "log item i = <i>")
func append(t *testing.T, start int) {
	lg := mkLog(t)
	defer lg.Close()

	for i := start; i < start+50; i++ {
		err := lg.Append([]byte(fmt.Sprintf("log item i = %d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Depends on results of TestAppend
func TestGet(t *testing.T) {
	lg := mkLog(t)
	defer lg.Close()

	checkIndex(t, lg, 99)
	for i := 0; i < 100; i++ {
		checkGet(t, lg, i)
	}
	lg.Close()
}

func checkGet(t *testing.T, lg *Log, i int) {
	data, err := lg.Get(int64(i))
	if err != nil {
		t.Fatal(err)
	}
	expected := fmt.Sprintf("log item i = %d", i)
	if expected != string(data) {
		t.Fatalf("Expected '%s', got '%s'", expected, string(data))
	}
}

// Depends on TestAppend, which should have inserted 100 records
func TestTruncate(t *testing.T) {
	lg := mkLog(t)
	err := lg.TruncateToEnd(50)
	lg.Close()
	if err != nil {
		t.Fatal(err)
	}

	lg = mkLog(t)
	checkIndex(t, lg, 49)
	lg.Close()
}

func checkIndex(t *testing.T, lg *Log, expected int) {
	i := lg.GetLastIndex()
	if i != int64(expected) {
		t.Fatal("Expected last index to be ", expected, " got ", i)
	}
}

func mkLog(t *testing.T) *Log {
	lg, err := Open("./logtest")
	if err != nil {
		t.Fatal(err)
	}
	return lg
}

func cleanup() {
	os.RemoveAll("./logtest")
}
