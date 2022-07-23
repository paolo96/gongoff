package gongoff

import (
	"fmt"
	"strconv"
)

type dataPiece interface {
	get() (string, error)
}

type data struct {
	variable  string
	separator separatorType
}

func (d *data) get() (string, error) {
	switch d.separator {
	case separatorTypeValue, separatorTypeDecimal, separatorTypeMultiply:
		_, err := strconv.Atoi(d.variable)
		if err != nil {
			return "", err
		}
		return d.variable + string(d.separator), nil
	case separatorTypeDescription:
		return string(d.separator) + d.variable + string(d.separator), nil
	case separatorTypeDescriptionDoubleHeight:
		return string(d.separator) + d.variable + string(separatorTypeDescription), nil
	default:
		return "", fmt.Errorf("separatorType %s is not supported", d.separator)
	}
}

type separatorType string

const (
	separatorTypeValue                   separatorType = "H"
	separatorTypeDecimal                 separatorType = "."
	separatorTypeMultiply                separatorType = "*"
	separatorTypeDescription             separatorType = "\""
	separatorTypeDescriptionDoubleHeight separatorType = "~\""
)

type terminator struct {
	variable       *string
	terminatorType terminatorType
}

func (t *terminator) get() (string, error) {
	if t.variable != nil {
		return *t.variable + string(t.terminatorType), nil
	} else {
		return string(t.terminatorType), nil
	}
}

//Refer to epson docs for more information on each terminatorType
type terminatorType string

const (
	terminatorTypeSold                                   terminatorType = "R"
	terminatorTypeDiscountDepartment                     terminatorType = "r"
	terminatorTypeSoldPLU                                terminatorType = "P"
	terminatorTypeCancellation                           terminatorType = "0M"
	terminatorTypeDiscountPercentTransaction             terminatorType = "1M"
	terminatorTypeDiscountPercentSubtotal                terminatorType = "2M"
	terminatorTypeDiscountValueTransaction               terminatorType = "3M"
	terminatorTypeDiscountValueSubtotal                  terminatorType = "4M"
	terminatorTypeIncreasePercentTransaction             terminatorType = "5M"
	terminatorTypeIncreasePercentSubtotal                terminatorType = "6M"
	terminatorTypeIncreaseValueTransaction               terminatorType = "7M"
	terminatorTypeIncreaseValueSubtotal                  terminatorType = "8M"
	terminatorTypeReturn                                 terminatorType = "9M" //Deprecated
	terminatorTypeCashIncome                             terminatorType = "10M"
	terminatorTypeCashOutflow                            terminatorType = "11M"
	terminatorTypePaymentWithCredit                      terminatorType = "12M" //Same as 4T
	terminatorTypeCashCreditRecovery                     terminatorType = "13M"
	terminatorTypeAdvancePayment                         terminatorType = "16M"
	terminatorTypeGift                                   terminatorType = "17M"
	terminatorTypeOneTimeCoupon                          terminatorType = "18M"
	terminatorTypeDirectInvoice                          terminatorType = "101M"
	terminatorTypeOpenCreditNote                         terminatorType = "102M" //Deprecated
	terminatorTypeOpenReturnDocumentCommercial           terminatorType = "104M"
	terminatorTypeOpenCancellationDocumentCommercial     terminatorType = "105M"
	terminatorTypeOpenReturnDocumentPOS                  terminatorType = "106M"
	terminatorTypeOpenCancellationDocumentPOS            terminatorType = "107M"
	terminatorTypeCommercialInvoice                      terminatorType = "111M"
	terminatorTypeCancelDocumentOrInvoice                terminatorType = "k"
	terminatorTypeSubtotal                               terminatorType = "="
	terminatorTypePaymentCash                            terminatorType = "1T"
	terminatorTypePaymentCheck                           terminatorType = "2T"
	terminatorTypePaymentCards                           terminatorType = "3T"
	terminatorTypePaymentCredit                          terminatorType = "4T"
	terminatorTypePaymentTicket                          terminatorType = "5T"
	terminatorTypePaymentCash2                           terminatorType = "6T" //Same as 1T
	terminatorTypePaymentUncollectedAssets               terminatorType = "7T"
	terminatorTypePaymentTicket2                         terminatorType = "20T" //Same as 5T
	terminatorTypePaymentTicket3                         terminatorType = "21T" //Same as 5T
	terminatorTypePaymentTicket4                         terminatorType = "22T" //Same as 5T
	terminatorTypePaymentUncollectedServices             terminatorType = "50T"
	terminatorTypePaymentUncollectedInvoice              terminatorType = "51T"
	terminatorTypePaymentUncollectedSSN                  terminatorType = "52T"
	terminatorTypePaymentDiscountGeneric                 terminatorType = "53T"
	terminatorTypePaymentOneTimeCoupon                   terminatorType = "54T"
	terminatorTypeAdditionalDescription                  terminatorType = "@"
	terminatorTypeLotteryCode                            terminatorType = "@37F" //Unused
	terminatorTypeInvoiceCustomerDetails                 terminatorType = "@38F"
	terminatorTypePrintVATorIDorLotteryCode              terminatorType = "@39F"
	terminatorTypePrintCourtesyMessage                   terminatorType = "@40F"
	terminatorTypePrintTrailerAfterLogo                  terminatorType = "@41F"
	terminatorTypePrintBarcodeEAN13                      terminatorType = "1Z"
	terminatorTypePrintBarcodeEAN8                       terminatorType = "2Z"
	terminatorTypePrintBarcodeCODE39                     terminatorType = "3Z"
	terminatorTypePrintNotCalculated                     terminatorType = "#"
	terminatorTypeOpenCashRegister                       terminatorType = "a"
	terminatorTypeClear                                  terminatorType = "K"
	terminatorTypeSelectOperator                         terminatorType = "O"
	terminatorTypeLockKeyboard                           terminatorType = "y"
	terminatorTypeUnlockKeyboard                         terminatorType = "Y"
	terminatorTypeOpenManagementDocument                 terminatorType = "j"
	terminatorTypePrintTextLine                          terminatorType = "@"
	terminatorTypeCloseManagementDocument                terminatorType = "J"
	terminatorTypeViewDescriptionOnDisplayFirstLine      terminatorType = "1%"
	terminatorTypeViewDescriptionOnDisplaySecondLine     terminatorType = "2%"
	terminatorTypeFinancialReportNoZeroing               terminatorType = "1f"
	terminatorTypeDepartmentReportNoZeroing              terminatorType = "2f"
	terminatorTypePLUReportNoZeroing                     terminatorType = "3f"
	terminatorTypeOperatorsReportNoZeroing               terminatorType = "4f"
	terminatorTypeFinancialReportZeroing                 terminatorType = "1F"
	terminatorTypeDepartmentReportZeroing                terminatorType = "2F"
	terminatorTypePLUReportZeroing                       terminatorType = "3F"
	terminatorTypeOperatorsReportZeroing                 terminatorType = "4F"
	terminatorTypeFinancialReportAndFiscalClosureZeroing terminatorType = "8F"
	terminatorTypeResetInvoiceNumber                     terminatorType = "9F"
	terminatorTypePrintFiscalMemoryAll                   terminatorType = "1w"
	terminatorTypePrintFiscalMemoryByDate                terminatorType = "2w"
	terminatorTypePrintFiscalMemoryByClosureNumber       terminatorType = "3w"
	terminatorTypePrintDetailsMemoryAll                  terminatorType = "4w"
	terminatorTypePrintDetailsMemoryByDate               terminatorType = "5w"
	terminatorTypePrintDetailsMemoryByClosureNumber      terminatorType = "6w"
	terminatorTypeSetDateTime                            terminatorType = "D"
	terminatorTypeDisableXonXoff                         terminatorType = "E"
	terminatorTypeDisableXonXoff2                        terminatorType = "1492E"
)

// Used for custom payments
func getTerminatorTypePaymentLight(paymentMethodCode string) (terminatorType, error) {
	if len(paymentMethodCode) != 3 {
		return "", fmt.Errorf("paymentMethodCode must be 3 characters long")
	}
	return terminatorType(paymentMethodCode + "T"), nil
}
