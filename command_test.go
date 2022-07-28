package gongoff

import (
	"fmt"
	"testing"
	"time"
)

// TestCommandProduct tests the creation of a CommandProduct command.
func TestCommandProduct(t *testing.T) {
	productName := "BREAD"
	quantity := 2
	department := 3
	commandProduct := NewCommandProduct(750, &productName, &quantity, &department)
	command, err := commandProduct.get()
	if err != nil {
		panic(err)
	}
	if command != "\"BREAD\"2*750H3R" {
		t.Errorf("Expected \"BREAD\"2*750H3R, got %s", command)
	}

	commandProductDefaults := NewCommandProduct(750, nil, nil, nil)
	commandDefaults, err := commandProductDefaults.get()
	if err != nil {
		panic(err)
	}
	if commandDefaults != "750HR" {
		t.Errorf("Expected 750HR, got %s", commandDefaults)
	}

	fmt.Println("Completed testCommandProduct")
}

// TestCommandTrailer tests the creation of a CommandMessage command.
func TestCommandTrailer(t *testing.T) {
	commandMessage := NewCommandTrailer("Hello World!")
	command, err := commandMessage.get()
	if err != nil {
		panic(err)
	}
	if command != "\"Hello World!                           \"@40F" {
		t.Errorf("Expected \"Hello World!                           \"@40F, got %s", command)
	}

	commandMessageLong := NewCommandTrailer("Hello World!                             truncate this")
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

func TestCommandCustomerIdentifier(t *testing.T) {
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

func TestCommandDiscountPercentage(t *testing.T) {
	commandDiscountPercentage := NewCommandDiscountPercentage(50.12)
	command, err := commandDiscountPercentage.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "50.121M" {
		t.Errorf("Expected 50.121M, got %s", command)
	}

	fmt.Println("Completed testCommandDiscountPercentage")
}

func TestCommandDiscountAmount(t *testing.T) {
	commandDiscountAmount := NewCommandDiscountAmount(1126)
	command, err := commandDiscountAmount.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "1126H3M" {
		t.Errorf("Expected 1126H3M, got %s", command)
	}

	fmt.Println("Completed testCommandDiscountAmount")
}

func TestCommandBarcode(t *testing.T) {
	commandBarcode, err := NewCommandBarcode("1234567890123")
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	command, err := commandBarcode.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "\"1234567890123\"1Z" {
		t.Errorf("Expected \"1234567890123\"1Z, got %s", command)
	}

	_, err = NewCommandBarcode("test")
	if err == nil {
		t.Errorf("Expected error != nil, got nil")
	}

	fmt.Println("Completed testCommandBarcode")
}

func TestCommandOpenDocumentCommercialReturn(t *testing.T) {
	testDate := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	documentId := NewDocumentId(12, 23, testDate, nil)
	commandOpenDocumentCommercialReturn := NewCommandOpenDocumentCommercialReturn(*documentId)
	command, err := commandOpenDocumentCommercialReturn.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "\"0012-0023-01-01-18\"104M" {
		t.Errorf("Expected \"0012-0023-01-01-18\"104M, got %s", command)
	}

	fmt.Println("Completed testCommandOpenDocumentCommercialReturn")
}

func TestCommandOpenInvoice(t *testing.T) {

	invNum := 123
	commandOpenInvoice := NewCommandOpenInvoice(&invNum)
	command, err := commandOpenInvoice.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "\"00123\"101M" {
		t.Errorf("Expected \"00123\"101M, got %s", command)
	}

	commandOpenInvoiceZero := NewCommandOpenInvoice(nil)
	commandZero, err := commandOpenInvoiceZero.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if commandZero != "\"00000\"101M" {
		t.Errorf("Expected \"00000\"101M, got %s", command)
	}

	fmt.Println("Completed testCommandOpenInvoice")
}

func TestCommandOpenInvoiceCommercialDocument(t *testing.T) {

	invNum := 123
	commandOpenInvoice := NewCommandOpenInvoiceCommercialDocument(&invNum)
	command, err := commandOpenInvoice.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "\"00123\"111M" {
		t.Errorf("Expected \"00123\"111M, got %s", command)
	}

	commandOpenInvoiceZero := NewCommandOpenInvoiceCommercialDocument(nil)
	commandZero, err := commandOpenInvoiceZero.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if commandZero != "\"00000\"111M" {
		t.Errorf("Expected \"00000\"111M, got %s", command)
	}

	fmt.Println("Completed testCommandOpenInvoice")
}

func TestCommandInvoiceDetails(t *testing.T) {

	commandInvoiceDetails := NewCommandInvoiceDetails("Mario Rossi")
	command, err := commandInvoiceDetails.get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if command != "\"Mario Rossi                             \"@38F" {
		t.Errorf("Expected \"Mario Rossi                             \"@38F, got %s", command)
	}

	fmt.Println("Completed testCommandInvoiceDetails")
}
