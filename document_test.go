package gongoff

import (
	"fmt"
	"testing"
	"time"
)

func TestDocumentCommercial(t *testing.T) {

	testProduct := "BREAD"
	commandProduct := NewCommandProduct(750, &testProduct, nil, nil)
	commandPayment, err := NewCommandPayment(TerminatorTypePaymentCash, nil, nil)
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	commercialDoc := NewDocumentCommercial(
		[]CommandProduct{
			*commandProduct,
		},
		[]CommandPayment{
			*commandPayment,
		},
		nil,
		nil,
		nil,
		NewCommandTrailer("Hello World!"),
	)

	commands := commercialDoc.get()
	if len(commands) != 3 {
		t.Errorf("Expected 3 commands, got %d", len(commands))
	}

	fmt.Println("Completed testDocumentCommercial")
}

func TestDocumentManagement(t *testing.T) {

	managementDoc := NewDocumentManagement([]string{"test", "test2", "test3"})
	commands := managementDoc.get()
	if len(commands) != 5 {
		t.Errorf("Expected 5 commands, got %d", len(commands))
	}

	openCommand, err := commands[0].get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if openCommand != "j" {
		t.Errorf("Expected j, got %s", openCommand)
	}

	closeCommand, err := commands[4].get()
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	if closeCommand != "J" {
		t.Errorf("Expected J, got %s", closeCommand)
	}

	fmt.Println("Completed testDocumentManagement")
}

func TestDocumentCommercialReturn(t *testing.T) {

	documentId := NewDocumentId(1, 2, time.Now(), nil)
	commandOpenDocumentCommercialReturn := NewCommandOpenDocumentCommercialReturn(*documentId)
	documentCommercialReturn := NewDocumentCommercialReturn(*commandOpenDocumentCommercialReturn, nil, nil)
	commands := documentCommercialReturn.get()
	if len(commands) != 1 {
		t.Errorf("Expected 1 commands, got %d", len(commands))
	}

	fmt.Println("Completed testDocumentCommercialReturn")
}

func TestDocumentCommercialCancellation(t *testing.T) {

	documentId := NewDocumentId(1, 2, time.Now(), nil)
	commandOpenDocumentCommercialCancellation := NewCommandOpenDocumentCommercialCancellation(*documentId)
	documentCommercialCancellation := NewDocumentCommercialCancellation(*commandOpenDocumentCommercialCancellation, nil, nil)
	commands := documentCommercialCancellation.get()
	if len(commands) != 1 {
		t.Errorf("Expected 1 commands, got %d", len(commands))
	}

	fmt.Println("Completed testDocumentCommercialCancellation")
}

func TestDocumentPOSReturn(t *testing.T) {

	commandOpenDocumentPOSReturn := NewCommandOpenDocumentPOSReturn(time.Now())
	documentPOSReturn := NewDocumentPOSReturn(*commandOpenDocumentPOSReturn, nil, nil)
	commands := documentPOSReturn.get()
	if len(commands) != 1 {
		t.Errorf("Expected 1 commands, got %d", len(commands))
	}

	fmt.Println("Completed testDocumentPOSReturn")
}

func TestDocumentPOSCancellation(t *testing.T) {

	commandOpenDocumentPOSCancellation := NewCommandOpenDocumentPOSCancellation(time.Now())
	documentPOSCancellation := NewDocumentPOSCancellation(*commandOpenDocumentPOSCancellation, nil, nil)
	commands := documentPOSCancellation.get()
	if len(commands) != 1 {
		t.Errorf("Expected 1 commands, got %d", len(commands))
	}

	fmt.Println("Completed testDocumentPOSCancellation")
}

func TestDocumentInvoice(t *testing.T) {

	invId := 3
	commandOpenDocumentInvoice := NewCommandOpenInvoice(&invId)
	commandCustomerDetails := NewCommandInvoiceDetails("test")
	commandPayment, err := NewCommandPayment(TerminatorTypePaymentCash, nil, nil)
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	documentInvoice, err := NewDocumentInvoice(
		*commandOpenDocumentInvoice,
		[]CommandInvoiceDetails{*commandCustomerDetails},
		[]CommandProduct{
			*NewCommandProduct(750, nil, nil, nil),
		},
		[]CommandPayment{
			*commandPayment,
		},
	)
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	commands := documentInvoice.get()
	if len(commands) != 4 {
		t.Errorf("Expected 4 commands, got %d", len(commands))
	}

	fmt.Println("Completed testDocumentInvoice")
}

func TestDocumentCommercialWithInvoice(t *testing.T) {

	commandPayment, err := NewCommandPayment(TerminatorTypePaymentCash, nil, nil)
	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	documentCommercial := NewDocumentCommercial(
		[]CommandProduct{
			*NewCommandProduct(750, nil, nil, nil),
		},
		[]CommandPayment{
			*commandPayment,
		},
		nil,
		nil,
		nil,
		nil,
	)
	commDoc := 3
	commandOpenDocumentInvoice := NewCommandOpenInvoiceCommercialDocument(&commDoc)
	commandCustomerDetails := NewCommandInvoiceDetails("test")

	NewDocumentCommercialWithInvoice(
		*commandOpenDocumentInvoice,
		[]CommandInvoiceDetails{*commandCustomerDetails},
		*documentCommercial,
	)

	if err != nil {
		t.Errorf("Expected error = nil, got %s", err)
	}
	commands := documentCommercial.get()
	if len(commands) != 2 {
		t.Errorf("Expected 3 commands, got %d", len(commands))
	}

	fmt.Println("Completed testDocumentCommercialWithInvoice")
}
