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

// TestCommandPayment tests the creation of a CommandPayment command.
func TestCommandPayment(t *testing.T) {
	testDesc := "test"
	testAmount := 100
	commandPayment, err := NewCommandPayment(TerminatorTypePaymentCards, &testAmount, &testDesc)
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	command, err := commandPayment.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "100H\"test\"3T" {
		t.Errorf("Expected 100H\"test\"3T, got %s", command)
	}

	_, err = NewCommandPayment(TerminatorTypeSold, nil, nil)
	if err == nil {
		t.Errorf("Expected error != nil, got nil")
	}

	fmt.Println("Completed testCommandPayment")
}

func TestNewCommandCustomerIdentifier(t *testing.T) {
	commandCustomerIdentifier, err := NewCommandCustomerIdentifier("RSSMRA00A01F205F")
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	command, err := commandCustomerIdentifier.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "\"RSSMRA00A01F205F\"@39F" {
		t.Errorf("Expected \"RSSMRA00A01F205F\"@39F, got %s", command)
	}

	_, err = NewCommandCustomerIdentifier("test")
	if err == nil {
		t.Errorf("Expected error != nil, got nil")
	}

	fmt.Println("Completed testNewCommandCustomerIdentifier")
}
