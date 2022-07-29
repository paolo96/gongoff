package gongoff

import (
	"errors"
	"strconv"
	"strings"
	"time"
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
func NewCommandProduct(unitPrice int, product *string, quantity *int, department *int) *CommandProduct {

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
		defaultDepartment := "1"
		commandProduct.terminator = Terminator{variable: &defaultDepartment, terminatorType: TerminatorTypeSold}
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

	paddedTrailer := trailer
	if len(trailer) < 39 {
		paddedTrailer += strings.Repeat(" ", 39-len(trailer))
	} else if len(trailer) > 46 {
		paddedTrailer = trailer[:46]
	}
	commandTrailer.data = []Data{
		{variable: paddedTrailer, separator: SeparatorTypeDescription},
	}
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
	commandCustomerIdentifier.data = []Data{
		{variable: customerIdentifier, separator: SeparatorTypeDescription},
	}
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
	commandDiscountPercentage.data = []Data{
		{variable: strconv.FormatFloat(discountPercentage, 'f', 2, 64), separator: SeparatorTypeDecimal},
	}
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
	commandDiscountAmount.data = []Data{
		{variable: strconv.Itoa(discountAmount), separator: SeparatorTypeValue},
	}
	commandDiscountAmount.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeDiscountValueTransaction}
	return commandDiscountAmount
}

type CommandBarcode struct {
	CommandGeneric
	barcode string
}

// NewCommandBarcode prints a barcode (EAN13, EAN8).
// Ex. ("1234567890123") -> "1234567890123"@39F -> Barcode is "1234567890123"
// Barcode is a 13- or 8-character string.
func NewCommandBarcode(barcode string) (*CommandBarcode, error) {
	commandBarcode := &CommandBarcode{
		barcode: barcode,
	}
	commandBarcode.data = []Data{
		{variable: barcode, separator: SeparatorTypeDescription},
	}
	if len(barcode) == 13 {
		commandBarcode.terminator = Terminator{variable: nil, terminatorType: TerminatorTypePrintBarcodeEAN13}
	} else if len(barcode) == 8 {
		commandBarcode.terminator = Terminator{variable: nil, terminatorType: TerminatorTypePrintBarcodeEAN8}
	} else {
		return nil, errors.New("barcode must be 13 or 8 characters long")
	}
	return commandBarcode, nil
}

type CommandOpenDocumentCommercialReturn struct {
	CommandGeneric
	documentId DocumentId
}

// NewCommandOpenDocumentCommercialReturn opens a commercial return document.
// Ex. (DocumentId[1, 2, date[01/01/2000], "1234567890123"]) -> "0001-0002-01-01-00-1234567890123"104M -> Open document 0001-0002-01-01-00 done with printer 1234567890123 for return.
func NewCommandOpenDocumentCommercialReturn(id DocumentId) *CommandOpenDocumentCommercialReturn {

	commandOpenDocumentCommercialReturn := &CommandOpenDocumentCommercialReturn{
		documentId: id,
	}
	commandOpenDocumentCommercialReturn.data = []Data{
		{variable: string(id), separator: SeparatorTypeDescription},
	}
	commandOpenDocumentCommercialReturn.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeOpenReturnDocumentCommercial}
	return commandOpenDocumentCommercialReturn
}

type CommandOpenDocumentCommercialCancellation struct {
	CommandGeneric
	documentId DocumentId
}

// NewCommandOpenDocumentCommercialCancellation opens a commercial cancellation document.
// Ex. (DocumentId[1, 2, date[01/01/2000], "1234567890123"]) -> "0001-0002-01-01-00-1234567890123"104M -> Open document 0001-0002-01-01-00 done with printer 1234567890123 for cancellation.
func NewCommandOpenDocumentCommercialCancellation(id DocumentId) *CommandOpenDocumentCommercialCancellation {

	commandOpenDocumentCommercialCancellation := &CommandOpenDocumentCommercialCancellation{
		documentId: id,
	}
	commandOpenDocumentCommercialCancellation.data = []Data{
		{variable: string(id), separator: SeparatorTypeDescription},
	}
	commandOpenDocumentCommercialCancellation.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeOpenCancellationDocumentCommercial}
	return commandOpenDocumentCommercialCancellation
}

type CommandOpenDocumentPOSReturn struct {
	CommandGeneric
	date time.Time
}

// NewCommandOpenDocumentPOSReturn opens a POS return document.
// Ex. (date[01/01/2000]) -> "01-01-00/POS"106M -> Open document done with printer for return.
func NewCommandOpenDocumentPOSReturn(date time.Time) *CommandOpenDocumentPOSReturn {
	commandOpenDocumentPOSReturn := &CommandOpenDocumentPOSReturn{
		date: date,
	}
	commandOpenDocumentPOSReturn.data = []Data{
		{variable: date.Format("01-02-06") + "/POS", separator: SeparatorTypeDescription},
	}
	commandOpenDocumentPOSReturn.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeOpenReturnDocumentPOS}
	return commandOpenDocumentPOSReturn
}

type CommandOpenDocumentPOSCancellation struct {
	CommandGeneric
	date time.Time
}

// NewCommandOpenDocumentPOSCancellation opens a POS cancellation document.
// Ex. (date[01/01/2000]) -> "01-01-00/POS"107M -> Open document done with printer for cancellation.
func NewCommandOpenDocumentPOSCancellation(date time.Time) *CommandOpenDocumentPOSCancellation {
	commandOpenDocumentPOSCancellation := &CommandOpenDocumentPOSCancellation{
		date: date,
	}
	commandOpenDocumentPOSCancellation.data = []Data{
		{variable: date.Format("01-02-06") + "/POS", separator: SeparatorTypeDescription},
	}
	commandOpenDocumentPOSCancellation.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeOpenCancellationDocumentPOS}
	return commandOpenDocumentPOSCancellation
}

type CommandOpenInvoice struct {
	CommandGeneric
	invoiceNumber *int
}

// NewCommandOpenInvoice opens an invoice document.
// Ex. (1) -> "00001"101M -> Open invoice n. 00001.
// If invoiceNumber is nil the numeration is delegated to the printer
func NewCommandOpenInvoice(invoiceNumber *int) *CommandOpenInvoice {
	commandOpenInvoice := &CommandOpenInvoice{
		invoiceNumber: invoiceNumber,
	}
	invNumString := "00000"
	if invoiceNumber != nil && *invoiceNumber > 0 && *invoiceNumber < 100000 {
		prefixedNum := strconv.Itoa(*invoiceNumber)
		invNumString = strings.Repeat("0", 5-len(prefixedNum)) + prefixedNum
	}
	commandOpenInvoice.data = []Data{
		{variable: invNumString, separator: SeparatorTypeDescription},
	}
	commandOpenInvoice.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeDirectInvoice}
	return commandOpenInvoice
}

type CommandOpenInvoiceCommercialDocument struct {
	CommandGeneric
	invoiceNumber *int
}

// NewCommandOpenInvoiceCommercialDocument opens an invoice document that follows a DocumentCommercial.
// Ex. (1) -> "00001"111M -> Open invoice n. 00001.
// If invoiceNumber is nil the numeration is delegated to the printer
func NewCommandOpenInvoiceCommercialDocument(invoiceNumber *int) *CommandOpenInvoiceCommercialDocument {
	commandOpenInvoice := &CommandOpenInvoiceCommercialDocument{
		invoiceNumber: invoiceNumber,
	}
	invNumString := "00000"
	if invoiceNumber != nil && *invoiceNumber > 0 && *invoiceNumber < 100000 {
		prefixedNum := strconv.Itoa(*invoiceNumber)
		invNumString = strings.Repeat("0", 5-len(prefixedNum)) + prefixedNum
	}
	commandOpenInvoice.data = []Data{
		{variable: invNumString, separator: SeparatorTypeDescription},
	}
	commandOpenInvoice.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeInvoiceCommercialDocument}
	return commandOpenInvoice
}

type CommandInvoiceDetails struct {
	CommandGeneric
	details string
}

// NewCommandInvoiceDetails prints the customer invoice details.
// Ex. ("Mario Rossi") -> "Mario Rossi                             "@38F -> Print invoice details.
// At least 40 characters are required. Max 46 characters.
func NewCommandInvoiceDetails(details string) *CommandInvoiceDetails {
	if len(details) < 40 {
		details = details + strings.Repeat(" ", 40-len(details))
	} else if len(details) > 46 {
		details = details[:46]
	}
	commandInvoiceDetails := &CommandInvoiceDetails{
		details: details,
	}
	commandInvoiceDetails.data = []Data{
		{variable: details, separator: SeparatorTypeDescription},
	}
	commandInvoiceDetails.terminator = Terminator{variable: nil, terminatorType: TerminatorTypeInvoiceCustomerDetails}
	return commandInvoiceDetails
}

type CommandDisplayMessage struct {
	CommandGeneric
	message string
	line    int
}

// NewCommandDisplayMessage shows a message on the display.
// Ex. ("Mario Rossi", 2) -> "Mario Rossi"2% -> Show message "Mario Rossi" on line 2.
func NewCommandDisplayMessage(message string, line int) *CommandDisplayMessage {
	commandDisplayMessage := &CommandDisplayMessage{
		message: message,
		line:    line,
	}
	commandDisplayMessage.data = []Data{
		{variable: message, separator: SeparatorTypeDescription},
	}
	lineTerminator := TerminatorTypeViewDescriptionOnDisplayFirstLine
	if line > 1 {
		lineTerminator = TerminatorTypeViewDescriptionOnDisplaySecondLine
	}
	commandDisplayMessage.terminator = Terminator{variable: nil, terminatorType: lineTerminator}
	return commandDisplayMessage
}
