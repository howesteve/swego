package lichdata

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/philhofer/fwd"
)

const testData = "\xbflength of input data expected"

func testReader() byteReader {
	return strings.NewReader("30<" + testData + ">")
}

type basicReader struct {
	r io.Reader
}

func (br *basicReader) Read(buf []byte) (int, error) {
	return br.r.Read(buf)
}

func TestReader(t *testing.T) {
	readers := []io.Reader{
		testReader(),
		&basicReader{testReader()},
	}

	for _, r := range readers {
		name := reflect.TypeOf(r).Elem().Name()
		t.Run(name, func(t *testing.T) {
			data, err := ReadFrom(r)
			if err != nil {
				t.Errorf("err = %v, want: nil", err)
			}

			got := string(data)
			want := testData
			if got != want {
				t.Errorf("data = %q, want: %q", got, want)
			}
		})
	}
}

func BenchmarkReaderBufio(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readData(bufio.NewReader(testReader()))
	}
}

func BenchmarkReaderFwd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readData(fwd.NewReader(testReader()))
	}
}

func TestWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)

	n, err := io.WriteString(w, testData)
	if err != nil {
		t.Errorf("err = %v, want: nil", err)
	}

	nstr := strconv.Itoa(n)
	want := len(testData) + len(nstr) + 2
	if n != want {
		t.Errorf("n = %d, want: %d", n, want)
	}
}

type discard struct{}

func (discard) Write([]byte) (int, error) { return 0, nil }
func (discard) WriteByte(byte) error      { return nil }
func (discard) Flush() error              { return nil }

var _ BufWriter = discard{}

func BenchmarkWriterBufio(b *testing.B) {
	d := discard{}
	data := []byte(testData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writeData(bufio.NewWriter(d), data)
	}
}

func BenchmarkWriterFwd(b *testing.B) {
	d := discard{}
	data := []byte(testData)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		writeData(fwd.NewWriter(d), data)
	}
}
