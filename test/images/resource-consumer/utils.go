package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

const (
	consumeCPUBinary = "./consume-cpu/consume-cpu"
	consumeMemBinary = "stress"
)

// ConsumeCPU consumes a given number of millicores for the specified duration.
func ConsumeCPU(millicores int, durationSec int) {
	log.Printf("ConsumeCPU millicores: %v, durationSec: %v", millicores, durationSec)
	// creating new consume cpu process
	arg1 := fmt.Sprintf("-millicores=%d", millicores)
	arg2 := fmt.Sprintf("-duration-sec=%d", durationSec)
	consumeCPU := exec.Command(consumeCPUBinary, arg1, arg2)
	consumeCPU.Run()
}

// ConsumeMem consumes a given number of megabytes for the specified duration.
func ConsumeMem(megabytes int, durationSec int) {
	log.Printf("ConsumeMem megabytes: %v, durationSec: %v", megabytes, durationSec)
	megabytesString := strconv.Itoa(megabytes) + "M"
	durationSecString := strconv.Itoa(durationSec)
	// creating new consume memory process
	consumeMem := exec.Command(consumeMemBinary, "-m", "1", "--vm-bytes", megabytesString, "--vm-hang", "0", "-t", durationSecString)
	consumeMem.Run()
}

// GetCurrentStatus prints out a no-op.
func GetCurrentStatus() {
	log.Printf("GetCurrentStatus")
	// not implemented
}
