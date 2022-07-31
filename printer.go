package gongoff

import (
	"bufio"
	"errors"
	"fmt"
	"go.bug.st/serial"
	"net"
)

type Printer interface {
	Open() error
	IsOpen() bool
	PrintDocument(Document) error
	PrintCommands([]Command) error
	Close() error

	flush() error
}

type GenericPrinter struct {
	dst *bufio.Writer
}

func (p *GenericPrinter) IsOpen() bool {
	return p.dst != nil
}

func (p *GenericPrinter) flush() error {
	return p.dst.Flush()
}

func (p *GenericPrinter) PrintDocument(doc Document) error {
	return p.PrintCommands(doc.get())
}

// PrintCommands prints the given commands to the printer.
// It allows for more flexibility than PrintDocument
func (p *GenericPrinter) PrintCommands(commands []Command) error {
	for _, command := range commands {
		commandString, err := command.get()
		if err != nil {
			return err
		}
		_, err = p.dst.Write([]byte(commandString))
		if err != nil {
			return err
		}
	}
	err := p.flush()
	if err != nil {
		return err
	}
	return nil
}

type NetworkPrinter struct {
	GenericPrinter
	socket *net.Conn
	ip     string
	port   int
}

func NewNetworkPrinter(ip string, port int) *NetworkPrinter {
	return &NetworkPrinter{ip: ip, port: port}
}

func (p *NetworkPrinter) Open() error {

	socket, err := net.Dial("tcp", fmt.Sprintf("%s:%d", p.ip, p.port))
	if err != nil {
		return err
	}
	p.socket = &socket
	p.dst = bufio.NewWriter(socket)
	return nil

}

func (p *NetworkPrinter) Close() error {
	if p.socket != nil {
		sP := *p.socket
		err := sP.Close()
		if err != nil {
			return err
		}
		p.dst = nil
		p.socket = nil
		return nil
	} else {
		return nil
	}
}

type SerialPrinter struct {
	GenericPrinter
	serialPort *serial.Port
	port       string
	baudRate   int
}

func NewSerialPrinter(port string) *SerialPrinter {
	return &SerialPrinter{port: port, baudRate: 9600}
}

func (p *SerialPrinter) Open() error {

	ports, err := serial.GetPortsList()
	if err != nil {
		return err
	}
	if len(ports) == 0 {
		return errors.New("no serial ports found")
	}

	for _, port := range ports {
		if port == p.port {
			serialPort, err := serial.Open(port, &serial.Mode{BaudRate: p.baudRate})
			if err != nil {
				return err
			}
			p.serialPort = &serialPort
			p.dst = bufio.NewWriter(serialPort)
			return nil
		}
	}

	return errors.New("chosen serial port not found")
}

func (p *SerialPrinter) Close() error {
	if p.serialPort != nil {
		sP := *p.serialPort
		err := sP.Close()
		if err != nil {
			return err
		}
		p.dst = nil
		p.serialPort = nil
		return nil
	} else {
		return nil
	}
}
