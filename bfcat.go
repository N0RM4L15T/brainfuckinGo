package main

import (
	"bufio"
	"fmt"
	"os"
)

type Bfunc uint8

const (
	PointerIncrease Bfunc = iota
	PointerDecrease
	ValueIncrease
	ValueDecrease
	ValuePrint
	ValueScan
	LoopStart
	LoopEnd
)

const (
	pointerMax = 65535
	byteMax    = 127
)

type bfProgram struct {
	memories []byte
	pointer  uint16
	stack    []int
	reader   *bufio.Reader
}

func InitBF(file *os.File) *bfProgram {
	bf := bfProgram{}
	bf.memories = make([]byte, pointerMax+1)
	bf.pointer = 0
	bf.reader = bufio.NewReader(file)
	return &bf
}

func main() {

	checkArgs()

	fileNames := os.Args[1:]

	for _, selected := range fileNames {

		file, err := os.Open(selected)
		checkErr(err)

		var char byte

		bf := InitBF(file)

		for err == nil {
			char, err = bf.reader.ReadByte()
			bf.Control(char)
		}
	}

	fmt.Printf("\nAll File Read Successfully!\n")
}

func (bf *bfProgram) increaseP() {
	if bf.pointer != pointerMax {
		bf.pointer++
	} else {
		bf.pointer = 0
	}
}

func (bf *bfProgram) decreaseP() {
	if bf.pointer != 0 {
		bf.pointer--
	} else {
		bf.pointer = pointerMax
	}
}

func (bf *bfProgram) increaseV() {
	if bf.memories[bf.pointer] != byteMax {
		bf.memories[bf.pointer]++
	} else {
		bf.memories[bf.pointer] = 0
	}
}

func (bf *bfProgram) decreaseV() {
	if bf.memories[bf.pointer] != 0 {
		bf.memories[bf.pointer]--
	} else {
		bf.memories[bf.pointer] = byteMax
	}
}

func (bf *bfProgram) printV() {
	fmt.Printf("%d ", bf.memories[bf.pointer])
}

func (bf *bfProgram) scanV() {

	var input byte
	fmt.Scanf("%v", &input)

	if input > 127 {
		input = 0
	}

	bf.memories[bf.pointer] = input

}

func (bf *bfProgram) loopStart() {
	bf.stack = append(bf.stack, 0)
}

func (bf *bfProgram) loopEnd() {

	if bf.memories[bf.pointer] != 0 {
		bf.loopBack()
	} else {
		bf.loopFinish()
	}
}

func (bf *bfProgram) loopRefresh() {

	for n := range bf.stack {
		bf.stack[n]++
	}

}

func (bf *bfProgram) loopBack() {

	for n := range bf.stack {
		bf.stack[n] -= bf.stack[len(bf.stack)-1]
	}

	for i := 0; i < bf.stack[len(bf.stack)-1]; i++ {
		bf.reader.UnreadByte()
	}

}

func (bf *bfProgram) loopFinish() {

	bf.stack[len(bf.stack)-1] = 0
	bf.stack = bf.stack[:len(bf.stack)-1]

}

func (bf *bfProgram) Control(char byte) {

	switch char {
	case '>':
		bf.increaseP()
	case '<':
		bf.decreaseP()
	case '+':
		bf.increaseV()
	case '-':
		bf.decreaseV()
	case '.':
		bf.printV()
	case ',':
		bf.scanV()
	case '[':
		bf.loopStart()
	case ']':
		bf.loopEnd()
		return
	}

	bf.loopRefresh()
}

func checkArgs() {

	if len(os.Args) == 1 {
		fmt.Println("Arguments should be more than one.")
		os.Exit(0)
	}

}

func checkErr(err error) {

	if err != nil {
		fmt.Println(err)
	}

}
