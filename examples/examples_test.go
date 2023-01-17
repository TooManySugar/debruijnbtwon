package debruijnbtwon_examples_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
	"math/rand"

	db2n "debruijnbtwon"
)

func CountB25Example() int {
	count := 0
	countingDiscard := func(dbs uint64) bool {
		count++
		return false
	}

	db2n.FindDeBruijnSeqK2N(5, countingDiscard)

	fmt.Println("found:", count)

	return count
}

func WriteB24Example(w io.Writer) {
	printVisitor := func(dbs uint64) bool {
		fmt.Fprintf(w, "%3X\n", dbs)
		return false
	}

	db2n.FindDeBruijnSeqK2N(4, printVisitor)
}

func WriteB24SubsequencesExample(w io.Writer) {

	// use bytes.Buffer to minimize writes to actual writer
	bb := bytes.Buffer{}

	n := uint64(4)

	// count of subsequences
	sc := 1 << n
	shift := sc - int(n)

	count := 0
	printN4SubsequenceVisitor := func(dbs uint64) bool {
		count++
		fmt.Fprintf(&bb, "0x%04x:", dbs)

		for j := 0; j < sc; j++ {
			// those casts required to make multiplication overflow
			// instead mask 01..1 where 1..1 with lenght n might be used
			// checkout debruijnbtwon_test.go magicTestNFactory how to do so
			testn := uint16(1) << j
			v := testn * uint16(dbs) >> shift
			fmt.Fprintf(&bb, " %2d", v)
		}
		fmt.Fprintln(&bb)

		w.Write(bb.Bytes())
		bb.Reset()

		return false
	}

	db2n.FindDeBruijnSeqK2N(n, printN4SubsequenceVisitor)
}

// Part of math.rand Source
type int63Source interface {
	Int63() int64
}

type bufferedRandBitSource struct {
	source int63Source
	offset uint8
	buf int64
}

func (s *bufferedRandBitSource) Bit() bool {
	if s.offset >= 64 {
		s.offset = 0
		s.buf = s.source.Int63()
	}
	bit := (s.buf >> s.offset&1) == 1
	s.offset++
	return bit
}

func Get1RandomB26Example() uint64 {
	src := rand.NewSource(time.Now().UnixNano())

	s := &bufferedRandBitSource{
		offset: 64,
		source: src,
	}
	var res uint64
	onRandFound := func(dbs uint64) bool {
		res = dbs
		return true
	}

	db2n.RandFindDeBruijnSeqK2N(6, s, onRandFound)

	return res
}

func CreateLog2LookupTableExample(dbs uint64) []byte {
	invert := make([]byte, 64)

	for j := 0; j < 64; j++ {
		v := uint64(1) << j * dbs >> 58
		invert[v] = byte(j)
	}

	return invert
}

// Examples running tests
// You can execute example using:
//     $ go test . -v -run <example name>
// Or run them all by:
//     $ go test . -v
////////////////////////////////////////////////////////////////////////////////

func TestCountB25Example(t *testing.T) {
	res := CountB25Example()
	ref := 2048

	if res != ref {
		t.Fatalf("\nExpected to count:\n%d\ngot:\n%d", ref, res)
	}
}

type ValidationWriter struct {
	w io.Writer
	r io.Reader

	t *testing.T
}

func (vw *ValidationWriter) Write(b []byte) (int, error) {
	wi, we := vw.w.Write(b)

	buf := make([]byte, len(b))

	_, err := vw.r.Read(buf)
	if err != nil {
		if err == io.EOF {
			vw.t.Fatalf("\nUnnexpected additional output:\n%#v", string(b))
		}
		vw.t.Fatal(err.Error())
	}

	if bytes.Compare(b, buf) != 0 {
		vw.t.Fatalf("\nExpected:\n%#v\ngot\n%#v", string(buf), string(b))
	}

	return wi, we
}

func utilTestWritingFuncOutput(t *testing.T,
	f func(io.Writer),
	outputRefPath string) {

	refFile, err := os.Open(outputRefPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer refFile.Close()

	vw := ValidationWriter{
		w: os.Stdout,
		r: refFile,
		t: t,
	}

	f(&vw)

	leftovers, err := io.ReadAll(refFile)
	if len(leftovers) == 0 {
		return
	}

	vw.t.Fatalf("Expected to print additional data:\n%#v", string(leftovers))
}

func utilUpdateWritingFuncOutput(f func(io.Writer), outputRefPath string) {
	refFile, err := os.Create(outputRefPath)
	if err != nil {
		panic(err.Error())
	}
	defer refFile.Close()

	f(refFile)
}

func TestWriteB24Example(t *testing.T) {
	f := WriteB24Example
	refPath := "./test/WriteB24ExampleOutput.bin"

	// utilUpdateWritingFuncOutput(f, refPath)
	utilTestWritingFuncOutput(t, f, refPath)
}

func TestWriteB24SubsequencesExample(t *testing.T) {
	f := WriteB24SubsequencesExample
	refPath := "./test/WriteB24SubsequencesExampleOutput.bin"

	// utilUpdateWritingFuncOutput(f, refPath)
	utilTestWritingFuncOutput(t, f, refPath)
}

func TestGet1RandomB26Example(t *testing.T) {
	fmt.Printf("random B(2,6) sequence: 0x%016X\n", Get1RandomB26Example())
}

func TestCreateLog2LookupTableExample(t *testing.T) {
	randomDBSeq := Get1RandomB26Example()
	lookupT := CreateLog2LookupTableExample(randomDBSeq)

	log2fromPow2n := func (pow2n uint64) uint8 {
		return lookupT[pow2n * randomDBSeq >> 58]
	}

	for j := 0; j < 64; j++ {
		if log2fromPow2n(uint64(1) << j) != byte(j) {
			t.Fatal("invalid table generated")
		}
	}

	fmt.Printf("magic:\n0x%016X\n", randomDBSeq)

	fmt.Printf("lookup table:\n[64]uint8{\n\t%2d", lookupT[0])
	for j := 1; j < 64; j++ {
		fmt.Printf(",")
		if j % 16 == 0 {
			fmt.Printf("\n\t")
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%2d", lookupT[j])
	}
	fmt.Printf("}\n")
}
