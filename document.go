package gongoff

type Document interface {
	get() []Command
}

type DocumentGeneric struct {
	commands []Command
}

func (d *DocumentGeneric) get() []Command {
	return d.commands
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

type DocumentCommercialReturn struct {
	DocumentGeneric
}
