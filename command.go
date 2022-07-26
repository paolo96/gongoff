package gongoff

import (
	"errors"
	"strconv"
	"strings"
)

// Command anatomy:
// Data(variable, separator), Data(variable, separator), ..., Terminator(*variable, terminatorType)

type Command interface {
	get() (string, error)
}

type CommandGeneric struct {
	data       []Data
	terminator Terminator
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

	if product != nil && len(*product) > 38 {
		productDesc := (*product)[:38]
		product = &productDesc
	}

	commandProduct := &CommandProduct{
		product:    product,
		unitPrice:  unitPrice,
		quantity:   quantity,
		department: department,
	}
	commandProduct.data = []Data{}

	if product != nil {
		commandProduct.data = append(commandProduct.data, Data{variable: *product, separator: SeparatorTypeDescription})
	}

	if quantity != nil {
		commandProduct.data = append(commandProduct.data, Data{variable: strconv.Itoa(*quantity), separator: SeparatorTypeMultiply})
	}

	if unitPrice != 0 {
		commandProduct.data = append(commandProduct.data, Data{variable: strconv.Itoa(unitPrice), separator: SeparatorTypeValue})
	}

	if department != nil {
		departmentString := strconv.Itoa(*department)
		commandProduct.terminator = Terminator{variable: &departmentString, terminatorType: TerminatorTypeSold}
	} else {
		commandProduct.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeSold}
	}

	return commandProduct
}

type CommandTrailer struct {
	CommandGeneric
	trailer string
}

// NewCommandTrailer prints a line with the given message.
// Ex. ("TRAILER") -> "TRAILER"@40F -> "TRAILER                                "
// Width forced between 39 and 46 characters.
func NewCommandTrailer(trailer string) *CommandTrailer {
	commandTrailer := &CommandTrailer{
		trailer: trailer,
	}
	commandTrailer.data = []Data{}

	paddedTrailer := trailer
	if len(trailer) < 39 {
		paddedTrailer += strings.Repeat(" ", 39-len(trailer))
	} else if len(trailer) > 46 {
		paddedTrailer = trailer[:46]
	}
	commandTrailer.data = append(commandTrailer.data, Data{variable: paddedTrailer, separator: SeparatorTypeDescription})
	commandTrailer.terminator = Terminator{variable: nil, terminatorType: TerminatorTypePrintCourtesyMessage}
	return commandTrailer
}

type CommandPayment struct {
	CommandGeneric
	paymentMethod TerminatorType
}

// NewCommandPayment prints a payment with the given parameters.
// Ex. (terminatorTypePaymentCards, 750) -> 750H3T -> Paid 750€ with cards.
// If no amount is given, the receipt is considered to be paid entirely with the given payment method.
// If the amount is given, change is applied accordingly.
func NewCommandPayment(paymentMethod TerminatorType, amount *int, paymentMethodDescription *string) (*CommandPayment, error) {
	if !strings.HasSuffix(string(paymentMethod), "T") {
		return nil, errors.New("payment method Terminator must end with 'T'")
	}
	commandPayment := &CommandPayment{
		paymentMethod: paymentMethod,
	}
	commandPayment.data = []Data{}

	if amount != nil {
		commandPayment.data = append(commandPayment.data, Data{variable: strconv.Itoa(*amount), separator: SeparatorTypeValue})
	}
	if paymentMethodDescription != nil {
		commandPayment.data = append(commandPayment.data, Data{variable: *paymentMethodDescription, separator: SeparatorTypeDescription})
	}

	commandPayment.terminator = Terminator{variable: nil, terminatorType: paymentMethod}
	return commandPayment, nil
}

type CommandCustomerIdentifier struct {
	CommandGeneric
	customerIdentifier string
}

// NewCommandCustomerIdentifier prints a customer identifier (CF, VAT, LotteryTicket)
// Ex. ("CF", "RSSMRA00A01F205F") -> "RSSMRA00A01F205F"@39F -> CF is "RSSMRA00A01F205F"
// CF (codice fiscale) is a 16-character string.
// VAT (vat number) is a 11-character string (without country prefix).
// LotteryTicket (lottery ticket) is a 8-character string.
func NewCommandCustomerIdentifier(customerIdentifier string) (*CommandCustomerIdentifier, error) {
	if len(customerIdentifier) != 16 && len(customerIdentifier) != 11 && len(customerIdentifier) != 8 {
		return nil, errors.New("customer identifier must be 16, 11 or 8 characters long")
	}

	commandCustomerIdentifier := &CommandCustomerIdentifier{
		customerIdentifier: customerIdentifier,
	}
	commandCustomerIdentifier.data = []Data{}
	commandCustomerIdentifier.data = append(commandCustomerIdentifier.data, Data{variable: customerIdentifier, separator: SeparatorTypeDescription})
	commandCustomerIdentifier.terminator = Terminator{variable: nil, terminatorType: TerminatorTypePrintCustomerIdentifier}
	return commandCustomerIdentifier, nil
}

type CommandDiscountPercentage struct {
	CommandGeneric
	discountPercentage float64
}

// NewCommandDiscountPercentage adds a discount percentage to the receipt.
// Ex. (10) -> 10.001M -> 10% discount on whole transaction.
func NewCommandDiscountPercentage(discountPercentage float64) *CommandDiscountPercentage {
	commandDiscountPercentage := &CommandDiscountPercentage{
		discountPercentage: discountPercentage,
	}
	commandDiscountPercentage.data = []Data{}
	commandDiscountPercentage.data = append(commandDiscountPercentage.data, Data{variable: strconv.FormatFloat(discountPercentage, 'f', 2, 64), separator: SeparatorTypeDecimal})
	commandDiscountPercentage.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeDiscountPercentTransaction}
	return commandDiscountPercentage
}

type CommandDiscountAmount struct {
	CommandGeneric
	discountAmount int
}

// NewCommandDiscountAmount adds a fixed value discount to the receipt.
// Ex. (1000) -> 10003M -> 10.00€ discount on whole transaction.
func NewCommandDiscountAmount(discountAmount int) *CommandDiscountAmount {
	commandDiscountAmount := &CommandDiscountAmount{
		discountAmount: discountAmount,
	}
	commandDiscountAmount.data = []Data{}
	commandDiscountAmount.data = append(commandDiscountAmount.data, Data{variable: strconv.Itoa(discountAmount), separator: SeparatorTypeValue})
	commandDiscountAmount.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeDiscountValueTransaction}
	return commandDiscountAmount
}
