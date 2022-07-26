package gongoff

import (
	"fmt"
	"testing"
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
