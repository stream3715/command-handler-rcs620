package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"go.bug.st/serial"
)

var sc = bufio.NewScanner(os.Stdin)

func main() {
	portName, err := getPortName()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	go readCommand(port)

	for {
		s := nextLine()
		switch s {
		case "0":
			sendCommand(port, []byte{0xd4, 0x02})

		case "1":
			sendCommand(port, []byte{0xd4, 0x01})

		case "r":
			sendCommand(port, []byte{0xD4, 0x18, 0x01})
			time.Sleep(time.Millisecond * 10)
			sendAck(port)
		case "q":
			os.Exit(0)
		}
	}
}

func nextLine() string {
	sc.Scan()
	return sc.Text()
}
