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
