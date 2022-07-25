package gongoff

import (
	"errors"
	"strconv"
	"strings"
)

// Command anatomy:
// data(variable, separator), data(variable, separator), ..., terminator(*variable, terminatorType)

type Command interface {
	get() (string, error)
}

type CommandGeneric struct {
	data       []data
	terminator terminator
}

func (c *CommandGeneric) get() (string, error) {
	var command string
	for _, d := range c.data {
		s, err := d.get()
		if err != nil {
			return "", err
		}
		command += s
	}
	s, err := c.terminator.get()
	if err != nil {
		return "", err
	}
	command += s
	return command, nil
}

type CommandProduct struct {
	CommandGeneric
	product    *string
	unitPrice  int
	quantity   *int
	department *int
}

// NewCommandProduct prints a product with the given parameters.
// Ex. ("BREAD", 750, 2, 3) -> "BREAD"2*750H3R -> Sold 2 loaves of bread for 7,50€ each in department 3.
func NewCommandProduct(product *string, unitPrice int, quantity *int, department *int) *CommandProduct {

	commandProduct := &CommandProduct{
		product:    product,
		unitPrice:  unitPrice,
		quantity:   quantity,
		department: department,
	}
	commandProduct.data = []data{}

	if product != nil {
		commandProduct.data = append(commandProduct.data, data{variable: *product, separator: separatorTypeDescription})
	}

	if quantity != nil {
		commandProduct.data = append(commandProduct.data, data{variable: strconv.Itoa(*quantity), separator: separatorTypeMultiply})
	}

	if unitPrice != 0 {
		commandProduct.data = append(commandProduct.data, data{variable: strconv.Itoa(unitPrice), separator: separatorTypeValue})
	}

	if department != nil {
		departmentString := strconv.Itoa(*department)
		commandProduct.terminator = terminator{variable: &departmentString, terminatorType: terminatorTypeSold}
	} else {
		commandProduct.terminator = terminator{variable: nil, terminatorType: terminatorTypeSold}
	}

	return commandProduct
}

type CommandTrailer struct {
	CommandGeneric
	trailer string
}

// NewCommandMessage prints a line with the given message, forcing the width between 39 and 46 characters.
// Ex. ("TRAILER") -> "TRAILER"@40F -> "TRAILER                                "
func NewCommandMessage(trailer string) *CommandTrailer {
	commandTrailer := &CommandTrailer{
		trailer: trailer,
	}
	commandTrailer.data = []data{}

	paddedTrailer := trailer
	if len(trailer) < 39 {
		paddedTrailer += strings.Repeat(" ", 39-len(trailer))
	} else if len(trailer) > 46 {
		paddedTrailer = trailer[:46]
	}
	commandTrailer.data = append(commandTrailer.data, data{variable: paddedTrailer, separator: separatorTypeDescription})
	commandTrailer.terminator = terminator{variable: nil, terminatorType: terminatorTypePrintCourtesyMessage}
	return commandTrailer
}

type CommandPayment struct {
	CommandGeneric
	paymentMethod terminatorType
}

// NewCommandPayment prints a payment with the given parameters.
// Ex. (terminatorTypePaymentCards, 750) -> 750H3T -> Paid 750€ with cards.
// If no amount is given, the receipt is considered to be paid entirely with the given payment method.
// If the amount is given, change is applied accordingly.
func NewCommandPayment(paymentMethod terminatorType, amount *int, paymentMethodDescription *string) (*CommandPayment, error) {
	if !strings.HasSuffix(string(paymentMethod), "T") {
		return nil, errors.New("payment method terminator must end with 'T'")
	}
	commandPayment := &CommandPayment{
		paymentMethod: paymentMethod,
	}
	commandPayment.data = []data{}

	if amount != nil {
		commandPayment.data = append(commandPayment.data, data{variable: strconv.Itoa(*amount), separator: separatorTypeValue})
	}
	if paymentMethodDescription != nil {
		commandPayment.data = append(commandPayment.data, data{variable: *paymentMethodDescription, separator: separatorTypeDescription})
	}

	commandPayment.terminator = terminator{variable: nil, terminatorType: paymentMethod}
	return commandPayment, nil
}
