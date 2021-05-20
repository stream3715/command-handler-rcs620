package lib

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func GetPortName() (string, error) {
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

func SendCommand(port serial.Port, command []byte) (int, error) {
	builded := buildCommand(command)
	fmt.Println(hex.EncodeToString(command))
	n, error := port.Write(builded)
	if error != nil {
		return -1, error
	}
	return n, nil
}

func SendAck(port serial.Port) (int, error) {
	builded := []byte{0x00, 0x00, 0xFF, 0x00, 0xFF, 0x00}
	n, error := port.Write(builded)
	if error != nil {
		return -1, error
	}
	return n, nil
}

func ReadCommand(port serial.Port) {
	pool := []byte{}
	for {
		buff := make([]byte, 300)
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
		}

		received := buff[:n]
		pool = append(pool, received...)

		for len(pool) > 5 {
			var splitted []byte
			var err error
			splitted, pool, err = popFirstResponse(pool)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(hex.EncodeToString(splitted))
		}
	}
}
