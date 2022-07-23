package gongoff

import (
	"fmt"
	"testing"
)

// TestCommandProduct tests the creation of a CommandProduct command.
func TestCommandProduct(t *testing.T) {
	productName := "BREAD"
	quantity := 2
	department := 3
	commandProduct := NewCommandProduct(&productName, 750, &quantity, &department)
	command, err := commandProduct.get()
	if err != nil {
		panic(err)
	}
	if command != "\"BREAD\"2*750H3R" {
		t.Errorf("Expected \"BREAD\"2*750H3R, got %s", command)
	}

	commandProductDefaults := NewCommandProduct(nil, 750, nil, nil)
	commandDefaults, err := commandProductDefaults.get()
	if err != nil {
		panic(err)
	}
	if commandDefaults != "750HR" {
		t.Errorf("Expected 750HR, got %s", commandDefaults)
	}

	fmt.Println("Completed testCommandProduct")
}

// TestCommandMessage tests the creation of a CommandMessage command.
func TestCommandMessage(t *testing.T) {
	commandMessage := NewCommandMessage("Hello World!")
	command, err := commandMessage.get()
	if err != nil {
		panic(err)
	}
	if command != "\"Hello World!                           \"@40F" {
		t.Errorf("Expected \"Hello World!                           \"@40F, got %s", command)
	}

	commandMessageLong := NewCommandMessage("Hello World!                             truncate this")
	command, err = commandMessageLong.get()
	if err != nil {
		panic(err)
	}
	if command != "\"Hello World!                             trunc\"@40F" {
		t.Errorf("Expected \"Hello World!                             trunc\"@40F, got %s", command)
	}

	fmt.Println("Completed testCommandMessage")
}
