package vm

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Runtime struct {
}

type VM struct {
	Stdout           io.Writer
	filename         string
	debug            bool
	content          *codeReader
	code             PyObject
	interned_strings []PyObject

	mainframe *PyFrame
	freevars  []PyObject // TODO: better place?! length!? is usage right?!
}

func (vm *VM) log(msg string) {
	log.Println(fmt.Sprintf("[VM] %s", msg))
}

func (vm *VM) parse() error {
	log.Println("Parsing...")
	if magic, _ := vm.content.readDWord(); magic != 168686339 {
		log.Fatal("No valid compiled python file (invalid magic)")
	}

	timestamp, _ := vm.content.readDWord()
	t := time.Unix(int64(timestamp), 0)
	log.Printf("File created: %s (timestamp: %d)\n", t, timestamp)

	vm.interned_strings = make([]PyObject, 0, 5000) // TODO: Wahllose Kapazität besser bestimmen!
	vm.code = vm.readObject()

	return nil
}

func (vm *VM) Filename() *string {
	return vm.code.(*PyCode).filename
}

func (vm *VM) Name() *string {
	return vm.code.(*PyCode).name
}

func (vm *VM) Run() error {
	vm.mainframe = NewPyFrame(uint64(vm.code.(*PyCode).stacksize))

	if retval, err := vm.code.(*PyCode).eval(vm.mainframe); err != nil {
		return err
	} else {
		vm.log(fmt.Sprintf("Returning value: %v (%T)", *retval.asString(), retval))
	}

	return nil
}

var debugMode bool = false

func NewVM(filename string, debug bool) (*VM, error) {
	debugMode = debug

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	vm := &VM{
		Stdout:   os.Stdout,
		content:  NewCodeReader(content),
		filename: filename,
		debug:    debug,
		freevars: make([]PyObject, 1000, 1000),
	}

	if err := vm.parse(); err != nil {
		return nil, err
	}

	return vm, nil
}
