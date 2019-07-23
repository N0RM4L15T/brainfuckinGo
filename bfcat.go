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
	LoopBegin
	LoopEnd
)

const (
	pointerMax = 65535
	byteMax    = 127
)

type bfProgram struct {
	memories  []byte
	pointer   uint16
	loopStack []int
	reader    *bufio.Reader
	command   []Bfunc
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
			bf.Interpret(char)
		}

		bf.Excute()
	}

	fmt.Printf("\nAll File Has Been Read Successfully!\n")
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
	fmt.Printf("%s", string(bf.memories[bf.pointer]))
}

func (bf *bfProgram) scanV() {

	var input byte
	fmt.Scanf("%v", &input)

	if input > 127 {
		input = 0
	}

	bf.memories[bf.pointer] = input

}

func (bf *bfProgram) startL(i int) int {

	if bf.memories[bf.pointer] != 0 {

		bf.loopStack = append(bf.loopStack, i)
		return i

	} else {

		for n := i; n <= len(bf.command); n++ {
			if bf.command[n] == LoopEnd {
				return n
			}
		}

		fmt.Println("Can't find end of loop")
		return i

	}

}

func (bf *bfProgram) finishL(i int) int {

	if bf.memories[bf.pointer] != 0 {

		return bf.loopStack[len(bf.loopStack)-1]

	} else {

		bf.loopStack = bf.loopStack[:len(bf.loopStack)-1]
		return i

	}

}

func (bf *bfProgram) Interpret(char byte) {

	insert := func(command Bfunc) {
		bf.command = append(bf.command, command)
	}

	switch char {
	case '>':
		insert(PointerIncrease)
	case '<':
		insert(PointerDecrease)
	case '+':
		insert(ValueIncrease)
	case '-':
		insert(ValueDecrease)
	case '.':
		insert(ValuePrint)
	case ',':
		insert(ValueScan)
	case '[':
		insert(LoopBegin)
	case ']':
		insert(LoopEnd)
	}

}

func (bf *bfProgram) Excute() {

	for i := 0; i < len(bf.command); i++ {

		switch bf.command[i] {
		case PointerIncrease:
			bf.increaseP()
		case PointerDecrease:
			bf.decreaseP()
		case ValueIncrease:
			bf.increaseV()
		case ValueDecrease:
			bf.decreaseV()
		case ValuePrint:
			bf.printV()
		case ValueScan:
			bf.scanV()
		case LoopBegin:
			i = bf.startL(i)
		case LoopEnd:
			i = bf.finishL(i)
		}

	}

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
