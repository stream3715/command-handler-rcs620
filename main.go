package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func getPortName() (string, error) {
	ports, error := enumerator.GetDetailedPortsList()
	if error != nil {
		return "", error
	}
	for _, port := range ports {
		if port.IsUSB && port.VID == "0403" && port.PID == "6001" {
			return port.Name, nil
		}
	}
	return "", errors.New("SERIAL_NOT_CONNECTED")
}

func main() {
	portName, err := getPortName()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(port)
	go sendCommand(port, []byte{0xd4, 0x02})
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func sendCommand(port serial.Port, command []byte) {
	send := []byte{0x00, 0x00, 0xff}
	len := len(command)
	lcs := 0x100 - len
	send = append(send, byte(len), byte(lcs))
	send = append(send, command...)
	dlen := 0
	for _, v := range command {
		dlen += int(v)
	}
	dcs := 0x100 - (dlen % 0x100)
	send = append(send, byte(dcs), 0x00)
	fmt.Println(send)
}
