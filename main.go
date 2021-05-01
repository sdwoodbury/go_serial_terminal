package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	serial "github.com/albenik/go-serial"
)

func readWrite(port serial.Port, terminalReaderChan chan []byte, terminalWriterChan chan []byte, shouldQuit chan bool) {
	buff := make([]byte, 128)
	readSerial := func() {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n > 0 {
			terminalWriterChan <- buff[:n]
		}
	}

	for {
		select {
		case <-shouldQuit:
			port.Close()
			return
		case incoming := <-terminalReaderChan:
			port.Write(incoming)
		default:
			readSerial()
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func terminalReader(terminalReaderChan chan []byte, shouldQuit chan bool) {
	buf := make([]byte, 128)
	for {
		select {
		case <-shouldQuit:
			return
		default:
			n, _ := os.Stdin.Read(buf)
			if n == 0 {
				time.Sleep(1 * time.Millisecond)
			} else {
				terminalReaderChan <- buf[:n]
			}
		}
	}
}

func terminalWriter(terminalWriterChan chan []byte, shouldQuit chan bool) {
	for {
		select {
		case <-shouldQuit:
			return
		case incoming := <-terminalWriterChan:
			os.Stdout.Write(incoming)
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func getPort(comPort string, baudRate int) serial.Port {
	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
	}
	port, err := serial.Open(comPort, mode)
	if err != nil {
		log.Fatal(err)
	}

	return port
}

func main() {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	baudRate, _ := strconv.ParseInt(os.Args[2], 10, 32)
	comPort := getPort(os.Args[1], int(baudRate))

	shouldQuit := make(chan bool, 3)
	terminalReaderChan := make(chan []byte, 128)
	terminalWriterChan := make(chan []byte, 128)

	go readWrite(comPort, terminalReaderChan, terminalWriterChan, shouldQuit)
	go terminalWriter(terminalWriterChan, shouldQuit)
	go terminalReader(terminalReaderChan, shouldQuit)

	<-signalChan

	log.Print("os.Interrupt - shutting down...\n")

	shouldQuit <- true
	shouldQuit <- true
	shouldQuit <- true

	os.Exit(0)
}
