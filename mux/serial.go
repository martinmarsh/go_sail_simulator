/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package mux

import (
	"fmt"
	"go_sail_simulator/buffer"
	"strconv"
	"time"
	"go.bug.st/serial"
)

func serialProcess(name string, config map[string][]string, channels *map[string](chan string)) {
	fmt.Println("started navmux serial " + name)
	baud, err := strconv.ParseInt(config["baud"][0], 10, 64)
	if err != nil {
		baud = 4800
	}
	mode := &serial.Mode{
		BaudRate: int(baud),
	}
	portName := config["name"][0]
	port, err := serial.Open(portName, mode)
	if err != nil {
		fmt.Println("no serial port " + portName)
	} else {
		if len(config["outputs"]) > 0 {
			fmt.Println("Open read serial port " + portName)
			go serialReader(name, port, config["outputs"], channels)

		}
		if len(config["input"]) > 0 {
			fmt.Println("Open write serial port " + portName)
			go serialWriter(name, port, config["input"], channels)
		}

	}

}

func serialReader(name string, port serial.Port, outputs []string, channels *map[string](chan string)) {
	buff := make([]byte, 50)
	cb := buffer.MakeByteBuffer(300, 100)
	time.Sleep(100 * time.Millisecond)
	for {
		n, err := port.Read(buff)
		if err != nil {
			fmt.Println("FATAL Error on port " + name)
			time.Sleep(5 * time.Second)
		}
		if n == 0 {
			fmt.Println("\nEOF on read of " + name)
			time.Sleep(5 * time.Second)
		} else {
			for i := 0; i < n; i++ {
				if buff[i] != 10 {
					cb.Write_byte(buff[i])
				}
			}
		}
		for {
			str := cb.ReadString()
			if len(str) == 0 {
				break
			}
			//fmt.Printf("\n\nSending from serial '%s'\n", str)
			for _, out := range outputs {
				(*channels)[out] <- str
			}

		}
	}
}

func serialWriter(name string, port serial.Port, input []string, channels *map[string](chan string)) {
	time.Sleep(100 * time.Millisecond)
	for {
		for _, in := range input {
			//fmt.Printf("\nWaiting for input to serial '%s'\n", in)
			str := <-(*channels)[in]
			str += "\r\n"
			//fmt.Printf("\nWriting to serial '%s'\n", str)
			_, err := port.Write([]byte(str))
			if err != nil {
				fmt.Println("FATAL Error on port" + name)
				time.Sleep(time.Minute)
			}

		}
	}

}
