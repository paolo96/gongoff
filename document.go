package gongoff

import (
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

type DocumentId string

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
	for _, command := range commandsProduct {
		commands = append(commands, &command)
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
	for _, command := range commandsPayment {
		commands = append(commands, &command)
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
	commands = append(commands, &CommandGeneric{
		data:       []Data{},
		terminator: Terminator{nil, TerminatorTypeOpenManagementDocument},
	})
	for _, row := range rows {
		if len(row) > 46 {
			row = row[:46]
		}
		commands = append(commands, &CommandGeneric{
			data: []Data{
				{variable: row, separator: SeparatorTypeDescription},
			},
			terminator: Terminator{
				variable:       nil,
				terminatorType: TerminatorTypeAdditionalDescription,
			},
		})
	}
	commands = append(commands, &CommandGeneric{
		data:       []Data{},
		terminator: Terminator{nil, TerminatorTypeCloseManagementDocument},
	})
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
