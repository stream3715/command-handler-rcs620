package main

import "errors"

func buildCommand(command []byte) []byte {
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
	return send
}

func popFirstResponse(pool []byte) ([]byte, []byte, error) {
	first := []byte{}
	var length, lcs, dcs byte
	var commandbody []byte
	var err error

	if pool[0] == 0x00 && pool[1] == 0x00 && pool[2] == 0xff {
		length = pool[3]
		lcs = pool[4]
		if length != 0 {
			commandbody = pool[5 : length+5]
			dcs = pool[length+5]
			for _, v := range commandbody {
				dcs += v
			}
			if (length+lcs != 0x00) || dcs != 0x00 {
				err = errors.New("INVALID_CHECKSUM")
			}
			if pool[length+6] != 0x00 {
				err = errors.New("INVALID_SUFFIX")
			} else {
				if len(pool) == int(length)+7 {
					first = commandbody
					pool = []byte{}
				} else {
					first = commandbody
					pool = pool[length+6:]
				}
			}
		} else if length == 0 {
			if len(pool) == 6 {
				first = []byte{}
				pool = []byte{}
			} else {
				first = []byte{}
				pool = pool[6:]
			}
		}
	} else {
		err = errors.New("INVALID_PREFIX")
	}
	if err != nil {
		nextIndex := findNextPreambleIndex(pool)
		first = pool[:nextIndex-1]
		pool = pool[nextIndex-1:]
	}
	return first, pool, err
}

func findNextPreambleIndex(pool []byte) int {
	for i := 1; i < len(pool)-2; i++ {
		if pool[i] == 0x00 && pool[i+1] == 0x00 || pool[i+2] == 0xff {
			return i
		}
	}
	return 0
}
