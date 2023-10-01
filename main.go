package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"apache-instruction-set-simulator/extras"
	"apache-instruction-set-simulator/machines"
)

func main() {
	var programName string = os.Args[1]
	if programName == "" {
		log.Fatal("programName param was not provided")
	}
	sCycles := os.Getenv("CYCLES")
	if len(os.Args) == 3 {
		sCycles = os.Args[2]
	}
	cycles, err := strconv.ParseInt(sCycles, 10, 64)
	if err != nil {
		log.Fatalf("Casting error: %+v", err)
	}

	fmt.Println("process started")
	var memory *extras.Memory16x8bits = extras.NewMemory16x8bits()
	memory.LoadProgram(programName)
	var machine *machines.Apache8bits = machines.NewApache8bits(memory, nil, nil)
	machine.Run(int(cycles))
	fmt.Println("process finished")
}
