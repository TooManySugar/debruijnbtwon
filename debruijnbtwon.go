package debruijnbtwon

import (
	"fmt"
)

type treeSearcher struct {
	// this 4 values are all constant for n
	width uint64
	uMasks []uint64
	sMasks []uint64
	numMask uint64

	// set of seen before subsequences
	seen []bool

	// function to handle found De Bruijn sequance
	// on return == true search stops
	onFound func(uint64) bool

	// stop search flag
	stop bool
}

// Recursive search on non-objective binary tree (there is no b-tree only value)
// visiting only non visited 'nodes' (there is no nodes only ownNum)
// offset is distance from 'root' to 'leaf' at current 'node', 0 => is leaf
func (ts *treeSearcher) step(value uint64, ownNum uint64, offset uint64) {
	if offset == 0 {
		if value & 1 == 0 {
			return
		}

		ts.stop = ts.onFound(value)
		return
	}

	value0 := value & ts.uMasks[offset]

	var nextNum uint64

	// ownNum
	// vvvvvv
	// _xxxxx0
	//  ^^^^^^
	//  nextNum
	nextNum = (ownNum << 1) & ts.numMask
	if !ts.seen[nextNum] {
		ts.seen[nextNum] = true
		ts.step(value0, nextNum, offset - 1)
		ts.seen[nextNum] = false
	}

	if ts.stop {
		return
	}

	// ownNum
	// vvvvvv
	// _xxxxx1
	//  ^^^^^^
	//  nextNum
	nextNum |= 1
	if !ts.seen[nextNum] {
		ts.seen[nextNum] = true
		ts.step(value0 | ts.sMasks[offset], nextNum, offset - 1)
		ts.seen[nextNum] = false
	}
}

func (ts *treeSearcher) initMasks() {

	var sz uint64 = 1 << ts.width + 1
	ts.uMasks = make([]uint64, sz)
	ts.sMasks = make([]uint64, sz)

	for i := uint64(0); i < sz; i++ {
		//      seq
		//   xxvvvvvv..
		// & 1111111100 <- this masks
		//   xxvvvvvv00
		ts.uMasks[i] = ^((uint64(1) << i) - 1)
		//      seq
		//   xxvvvvvv..
		// | 0000000010 <- this masks
		//   xxvvvvvv1.
		ts.sMasks[i] = uint64(1) << (i - 1)
	}
}

type ErrorOutOfRange struct {
	value uint64
}

func (e ErrorOutOfRange) Error() string {
	return fmt.Sprintf("n must be in range [1, 6] got: %d", e.value)
}

// Search on non-objective binary tree to find B(2,n) De Bruijn Sequences
//     n - lenght of all possible subsequences
//     onFound - function to handle found De Bruijn sequence
//               accepts sequence as unsigned integer without trailing zeroes
//               for given n amount of trailing zeroes is (n - 1)
//               return value defines if search stops (true) or not (false)
//
// returns ErrorOutOfRange if n not in range [1, 6] otherwise returns nil
//
func FindDeBruijnSeqK2N(n uint64, onFound func(uint64) bool) error {
	if n > 6 || n < 1 {
		return ErrorOutOfRange{n}
	}

	ts := treeSearcher {
		width: n,

		seen: make([]bool, 1 << n),
		onFound: onFound,
		stop: false,

		numMask: uint64(1) << n - 1,
	}
	ts.initMasks()

	// starting 'node'
	ts.seen[0] = true

	ts.step(uint64(0), 0, 1 << n - n)
	return nil
}
