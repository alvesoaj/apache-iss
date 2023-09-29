package main

import (
	"fmt"
	"log"
	"os"

	"apache-instruction-set-simulator/extras"
	"apache-instruction-set-simulator/machines"
	"apache-instruction-set-simulator/utils"
)

func main() {
	var programName string = os.Args[1]
	if programName == "" {
		log.Fatal("programName param was not provided")
	}
	var cycles int = int(utils.CastStringToUint8(os.Args[2], 10))
	if cycles == 0 {
		log.Fatal("cycles param was not provided")
	}

	fmt.Println("process started")
	var memory *extras.Memory16x8bits = extras.NewMemory16x8bits()
	memory.LoadProgram(programName)
	var machine *machines.Apache8bits = machines.NewApache8bits(memory)
	machine.Run(cycles)
	fmt.Println("process finished")
}
