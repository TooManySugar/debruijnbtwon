package debruijnbtwon_test

import (
	"fmt"
	"testing"
	"debruijnbtwon"
)

func utilTestFindB2NMustErr(t* testing.T, n uint64, expectedErrText string) {
	err := debruijnbtwon.FindDeBruijnSeqK2N(n, nil)
	if err == nil {
		t.Fatal("Expected error got nil")
	}
	if err.Error() != expectedErrText {
		t.Fatalf("Expected error:\n%#v\ngot error:\n%#v",
		         expectedErrText,
		         err.Error())
	}
	fmt.Printf("    %s\n", err.Error())
}

func magicTestNFactory(t *testing.T, n uint64) (func (magic uint64) bool,
                                                *uint64,
                                                error) {

	if n < 1 || n > 6 {
		return nil, nil, fmt.Errorf("n is out of range [1, 6]")
	}

	// pow2n at the same time is:
	//    maximum value of subsequence + 1
	//    total ammount of subsequences
	pow2n := uint64(1 << n)
	shift := pow2n - n
	mask  := pow2n - 1

	notMagicFormat :=
		fmt.Sprintf("\n0x%%0%dX\n0b%%0%db\nis not a De Bruijn sequence",
		            pow2n / 4,
		            pow2n)

	notMagic := func (magic uint64) {
		t.Fatalf(notMagicFormat, magic, magic)
	}

	count := uint64(0)

	magicTestN := func (magic uint64) bool {

		// If first bit is 0 where is no magic. sad
		if magic & 1 == 0 {
			notMagic(magic)
		}

		set := make([]bool, pow2n)
		var testn, v uint64

		for j := uint64(0); j < pow2n; j++ {
			testn = uint64(1) << j
			// Multiplication magic by testn is mostly left bit shift.
			// Doing so and shifting result back to right by 2^n - n (which what
			// shift variable is) we acquire subsequence at bits 0 to n
			// counting from right (big endian).
			// After which we apply mask to get uint64 with first 6 (or less)
			// significant bits which is a subsequence.
			v = (testn * magic >> shift) & mask

			if set[v] {
				notMagic(magic)
			}
			set[v] = true
		}

		count++
		return false
	}

	return magicTestN, &count, nil
}

func utilTestFindB2N(t *testing.T, n uint64) {
	magicTestN, count, err := magicTestNFactory(t, n)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = debruijnbtwon.FindDeBruijnSeqK2N(n, magicTestN)
	if err != nil {
		t.Fatal(err.Error())
	}

	// cardinality of de Bruijn B(2, N) sequences:
	//     2^(2^(n-1) - n)
	mustCount := uint64(1) << ((1 << (n - 1)) - n)
	if *count != mustCount {
		t.Fatalf("Found %d instead of %d sequences", *count, mustCount)
	}
	fmt.Printf("    %d sequences found\n", *count)
}

func TestFindB20(t *testing.T) {
	utilTestFindB2NMustErr(
		t,
		0,
		"n must be in range [1, 6] got: 0")
}

func TestFindB27(t *testing.T) {
	utilTestFindB2NMustErr(
		t,
		7,
		"n must be in range [1, 6] got: 7")
}

func TestFindB21(t *testing.T) {
	utilTestFindB2N(t, 1)
}

func TestFindB22(t *testing.T) {
	utilTestFindB2N(t, 2)
}

func TestFindB23(t *testing.T) {
	utilTestFindB2N(t, 3)
}

func TestFindB24(t *testing.T) {
	utilTestFindB2N(t, 4)
}

func TestFindB25(t *testing.T) {
	utilTestFindB2N(t, 5)
}

// Takes 44s at my Ryzen 5 3500U laptop
func TestFindB26(t *testing.T) {
	utilTestFindB2N(t, 6)
}
