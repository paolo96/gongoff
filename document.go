package gongoff

import (
	"errors"
	"fmt"
	"time"
)

type Document interface {
	get() []Command
}

type DocumentGeneric struct {
	commands []Command
}

func (d *DocumentGeneric) get() []Command {
	return d.commands
}

// DocumentId is the unique document identifier, optionally with the serial number of the printer which generated the document.
type DocumentId string

// NewDocumentId creates a new document identifier.
func NewDocumentId(dailyClosureNumber int, documentNumber int, documentDate time.Time, printerSerialNumber *string) *DocumentId {

	dailyClosureNumberString := "9999"
	if dailyClosureNumber > 0 && dailyClosureNumber < 10000 {
		dailyClosureNumberString = fmt.Sprintf("%04d", dailyClosureNumber)
	}
	documentNumberString := "9999"
	if documentNumber > 0 && documentNumber < 10000 {
		documentNumberString = fmt.Sprintf("%04d", documentNumber)
	}
	documentDateString := "30-01-20"
	if documentDate.Year() > 0 {
		documentDateString = documentDate.Format("02-01-06")
	}
	printerSerialNumberString := ""
	if printerSerialNumber != nil {
		printerSerialNumberString = "-" + *printerSerialNumber
	}
	id := DocumentId(dailyClosureNumberString + "-" + documentNumberString + "-" + documentDateString + printerSerialNumberString)
	return &id

}

// DocumentCommercial is commonly known as a fiscal receipt.
type DocumentCommercial struct {
	DocumentGeneric
}

func NewDocumentCommercial(
	commandsProduct []CommandProduct,
	commandsPayment []CommandPayment,
	commandDiscountAmount *CommandDiscountAmount,
	commandDiscountPercentage *CommandDiscountPercentage,
	commandCI *CommandCustomerIdentifier,
	commandTrailer *CommandTrailer) *DocumentCommercial {

	var commands []Command
	for i := range commandsProduct {
		commands = append(commands, &commandsProduct[i])
	}
	if commandDiscountAmount != nil {
		commands = append(commands, commandDiscountAmount)
	}
	if commandDiscountPercentage != nil {
		commands = append(commands, commandDiscountPercentage)
	}
	if commandCI != nil {
		commands = append(commands, commandCI)
	}
	if commandTrailer != nil {
		commands = append(commands, commandTrailer)
	}
	for i := range commandsPayment {
		commands = append(commands, &commandsPayment[i])
	}
	return &DocumentCommercial{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
	}
}

// DocumentManagement is a generic document useful for testing purposes and generic text print.
type DocumentManagement struct {
	DocumentGeneric
	rows []string
}

func NewDocumentManagement(rows []string) *DocumentManagement {

	var commands []Command
	commands = append(commands, NewCommandGeneric([]Data{}, Terminator{nil, TerminatorTypeOpenManagementDocument}))
	for _, row := range rows {
		if len(row) > 46 {
			row = row[:46]
		}
		commands = append(commands, NewCommandGeneric(
			[]Data{
				{variable: row, separator: SeparatorTypeDescription},
			},
			Terminator{
				variable:       nil,
				terminatorType: TerminatorTypeAdditionalDescription,
			}),
		)
	}
	commands = append(commands, NewCommandGeneric([]Data{}, Terminator{nil, TerminatorTypeCloseManagementDocument}))
	return &DocumentManagement{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
	}

}

// DocumentCommercialReturn is used when a customer returns a product.
// If no product is given, the whole content of the document is considered returned.
type DocumentCommercialReturn struct {
	DocumentGeneric
}

func NewDocumentCommercialReturn(
	commandOpen CommandOpenDocumentCommercialReturn,
	commandProduct *CommandProduct,
	commandPayment *CommandPayment) *DocumentCommercialReturn {

	commands := []Command{
		&commandOpen,
	}
	if commandProduct != nil {
		commands = append(commands, commandProduct)
	}
	if commandPayment != nil {
		commands = append(commands, commandPayment)
	}
	return &DocumentCommercialReturn{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
	}
}

// DocumentCommercialCancellation is used when a fiscal receipt is cancelled.
type DocumentCommercialCancellation struct {
	DocumentGeneric
}

func NewDocumentCommercialCancellation(
	commandOpen CommandOpenDocumentCommercialCancellation,
	commandProduct *CommandProduct,
	commandPayment *CommandPayment) *DocumentCommercialCancellation {

	commands := []Command{
		&commandOpen,
	}
	if commandProduct != nil {
		commands = append(commands, commandProduct)
	}
	if commandPayment != nil {
		commands = append(commands, commandPayment)
	}
	return &DocumentCommercialCancellation{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
	}
}

// DocumentPOSReturn is used when a customer returns a product.
// If no product is given, the whole content of the document is considered returned.
// The document was generated by a POS terminal.
type DocumentPOSReturn struct {
	DocumentGeneric
}

func NewDocumentPOSReturn(
	commandOpen CommandOpenDocumentPOSReturn,
	commandProduct *CommandProduct,
	commandPayment *CommandPayment) *DocumentPOSReturn {

	commands := []Command{
		&commandOpen,
	}
	if commandProduct != nil {
		commands = append(commands, commandProduct)
	}
	if commandPayment != nil {
		commands = append(commands, commandPayment)
	}
	return &DocumentPOSReturn{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
	}
}

// DocumentPOSCancellation is used when a fiscal receipt is cancelled.
// The document was generated by a POS terminal.
type DocumentPOSCancellation struct {
	DocumentGeneric
}

func NewDocumentPOSCancellation(
	commandOpen CommandOpenDocumentPOSCancellation,
	commandProduct *CommandProduct,
	commandPayment *CommandPayment) *DocumentPOSCancellation {

	commands := []Command{
		&commandOpen,
	}
	if commandProduct != nil {
		commands = append(commands, commandProduct)
	}
	if commandPayment != nil {
		commands = append(commands, commandPayment)
	}
	return &DocumentPOSCancellation{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
	}
}

// DocumentInvoice is a direct invoice document.
type DocumentInvoice struct {
	DocumentGeneric
}

func NewDocumentInvoice(
	commandOpen CommandOpenInvoice,
	customerDetails []CommandInvoiceDetails,
	products []CommandProduct,
	payments []CommandPayment) (*DocumentInvoice, error) {

	if len(customerDetails) == 0 || len(customerDetails) > 5 {
		return nil, errors.New("invalid number of customer details commands, must be between 1 and 5")
	}
	if len(products) == 0 {
		return nil, errors.New("invalid number of products commands, must be at least 1")
	}
	if len(payments) == 0 {
		return nil, errors.New("invalid number of payments commands, must be at least 1")
	}

	var commands []Command
	for i := range customerDetails {
		commands = append(commands, &customerDetails[i])
	}
	commands = append(commands, &commandOpen)
	for i := range products {
		commands = append(commands, &products[i])
	}
	for i := range payments {
		commands = append(commands, &payments[i])
	}

	return &DocumentInvoice{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
	}, nil
}

type DocumentCommercialWithInvoice struct {
	DocumentGeneric
	commercialDocument DocumentCommercial
}

func NewDocumentCommercialWithInvoice(
	commandOpen CommandOpenInvoiceCommercialDocument,
	customerDetails []CommandInvoiceDetails,
	commercialDocument DocumentCommercial) *DocumentCommercialWithInvoice {

	var commands []Command
	for i := range customerDetails {
		commands = append(commands, &customerDetails[i])
	}
	commands = append(commands, &commandOpen)

	return &DocumentCommercialWithInvoice{
		DocumentGeneric: DocumentGeneric{
			commands: commands,
		},
		commercialDocument: commercialDocument,
	}
}
