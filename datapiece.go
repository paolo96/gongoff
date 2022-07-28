package gongoff

import (
	"fmt"
	"strconv"
	"strings"
)

type dataPiece interface {
	get() (string, error)
}

type Data struct {
	variable  string
	separator SeparatorType
}

func (d *Data) get() (string, error) {
	switch d.separator {
	case SeparatorTypeValue, SeparatorTypeMultiply:
		_, err := strconv.Atoi(d.variable)
		if err != nil {
			return "", err
		}
		return d.variable + string(d.separator), nil
	case SeparatorTypeDecimal:
		if !strings.Contains(d.variable, ".") {
			return "", fmt.Errorf("decimal variable must contain '.'")
		}
		return d.variable, nil
	case SeparatorTypeDescription:
		return string(d.separator) + d.variable + string(d.separator), nil
	case SeparatorTypeDescriptionDoubleHeight:
		return string(d.separator) + d.variable + string(SeparatorTypeDescription), nil
	default:
		return "", fmt.Errorf("separatorType %s is not supported", d.separator)
	}
}

type SeparatorType string

const (
	SeparatorTypeValue                   SeparatorType = "H"
	SeparatorTypeDecimal                 SeparatorType = "."
	SeparatorTypeMultiply                SeparatorType = "*"
	SeparatorTypeDescription             SeparatorType = "\""
	SeparatorTypeDescriptionDoubleHeight SeparatorType = "~\""
)

type Terminator struct {
	variable       *string
	terminatorType TerminatorType
}

func (t *Terminator) get() (string, error) {
	if t.variable != nil {
		return *t.variable + string(t.terminatorType), nil
	} else {
		return string(t.terminatorType), nil
	}
}

// TerminatorType is documented in the epson docs, for more information refer to it
type TerminatorType string

const (
	TerminatorTypeSold                                   TerminatorType = "R"
	TerminatorTypeDiscountDepartment                     TerminatorType = "r"
	TerminatorTypeSoldPLU                                TerminatorType = "P"
	TerminatorTypeCancellation                           TerminatorType = "0M"
	TerminatorTypeDiscountPercentTransaction             TerminatorType = "1M"
	TerminatorTypeDiscountPercentSubtotal                TerminatorType = "2M"
	TerminatorTypeDiscountValueTransaction               TerminatorType = "3M"
	TerminatorTypeDiscountValueSubtotal                  TerminatorType = "4M"
	TerminatorTypeIncreasePercentTransaction             TerminatorType = "5M"
	TerminatorTypeIncreasePercentSubtotal                TerminatorType = "6M"
	TerminatorTypeIncreaseValueTransaction               TerminatorType = "7M"
	TerminatorTypeIncreaseValueSubtotal                  TerminatorType = "8M"
	TerminatorTypeReturn                                 TerminatorType = "9M" //Deprecated
	TerminatorTypeCashIncome                             TerminatorType = "10M"
	TerminatorTypeCashOutflow                            TerminatorType = "11M"
	TerminatorTypePaymentWithCredit                      TerminatorType = "12M" //Same as 4T
	TerminatorTypeCashCreditRecovery                     TerminatorType = "13M"
	TerminatorTypeAdvancePayment                         TerminatorType = "16M"
	TerminatorTypeGift                                   TerminatorType = "17M"
	TerminatorTypeOneTimeCoupon                          TerminatorType = "18M"
	TerminatorTypeDirectInvoice                          TerminatorType = "101M"
	TerminatorTypeOpenCreditNote                         TerminatorType = "102M" //Deprecated
	TerminatorTypeOpenReturnDocumentCommercial           TerminatorType = "104M"
	TerminatorTypeOpenCancellationDocumentCommercial     TerminatorType = "105M"
	TerminatorTypeOpenReturnDocumentPOS                  TerminatorType = "106M"
	TerminatorTypeOpenCancellationDocumentPOS            TerminatorType = "107M"
	TerminatorTypeInvoiceCommercialDocument              TerminatorType = "111M"
	TerminatorTypeCancelDocumentOrInvoice                TerminatorType = "k"
	TerminatorTypeSubtotal                               TerminatorType = "="
	TerminatorTypePaymentCash                            TerminatorType = "1T"
	TerminatorTypePaymentCheck                           TerminatorType = "2T"
	TerminatorTypePaymentCards                           TerminatorType = "3T"
	TerminatorTypePaymentCredit                          TerminatorType = "4T"
	TerminatorTypePaymentTicket                          TerminatorType = "5T"
	TerminatorTypePaymentCash2                           TerminatorType = "6T" //Same as 1T
	TerminatorTypePaymentUncollectedAssets               TerminatorType = "7T"
	TerminatorTypePaymentTicket2                         TerminatorType = "20T" //Same as 5T
	TerminatorTypePaymentTicket3                         TerminatorType = "21T" //Same as 5T
	TerminatorTypePaymentTicket4                         TerminatorType = "22T" //Same as 5T
	TerminatorTypePaymentUncollectedServices             TerminatorType = "50T"
	TerminatorTypePaymentUncollectedInvoice              TerminatorType = "51T"
	TerminatorTypePaymentUncollectedSSN                  TerminatorType = "52T"
	TerminatorTypePaymentDiscountGeneric                 TerminatorType = "53T"
	TerminatorTypePaymentOneTimeCoupon                   TerminatorType = "54T"
	TerminatorTypeAdditionalDescription                  TerminatorType = "@"
	TerminatorTypeLotteryCode                            TerminatorType = "@37F" //Unused
	TerminatorTypeInvoiceCustomerDetails                 TerminatorType = "@38F"
	TerminatorTypePrintCustomerIdentifier                TerminatorType = "@39F"
	TerminatorTypePrintCourtesyMessage                   TerminatorType = "@40F"
	TerminatorTypePrintTrailerAfterLogo                  TerminatorType = "@41F"
	TerminatorTypePrintBarcodeEAN13                      TerminatorType = "1Z"
	TerminatorTypePrintBarcodeEAN8                       TerminatorType = "2Z"
	TerminatorTypePrintBarcodeCODE39                     TerminatorType = "3Z"
	TerminatorTypePrintNotCalculated                     TerminatorType = "#"
	TerminatorTypeOpenCashRegister                       TerminatorType = "a"
	TerminatorTypeClear                                  TerminatorType = "K"
	TerminatorTypeSelectOperator                         TerminatorType = "O"
	TerminatorTypeLockKeyboard                           TerminatorType = "y"
	TerminatorTypeUnlockKeyboard                         TerminatorType = "Y"
	TerminatorTypeOpenManagementDocument                 TerminatorType = "j"
	TerminatorTypePrintTextLine                          TerminatorType = "@"
	TerminatorTypeCloseManagementDocument                TerminatorType = "J"
	TerminatorTypeViewDescriptionOnDisplayFirstLine      TerminatorType = "1%"
	TerminatorTypeViewDescriptionOnDisplaySecondLine     TerminatorType = "2%"
	TerminatorTypeFinancialReportNoZeroing               TerminatorType = "1f"
	TerminatorTypeDepartmentReportNoZeroing              TerminatorType = "2f"
	TerminatorTypePLUReportNoZeroing                     TerminatorType = "3f"
	TerminatorTypeOperatorsReportNoZeroing               TerminatorType = "4f"
	TerminatorTypeFinancialReportZeroing                 TerminatorType = "1F"
	TerminatorTypeDepartmentReportZeroing                TerminatorType = "2F"
	TerminatorTypePLUReportZeroing                       TerminatorType = "3F"
	TerminatorTypeOperatorsReportZeroing                 TerminatorType = "4F"
	TerminatorTypeFinancialReportAndFiscalClosureZeroing TerminatorType = "8F"
	TerminatorTypeResetInvoiceNumber                     TerminatorType = "9F"
	TerminatorTypePrintFiscalMemoryAll                   TerminatorType = "1w"
	TerminatorTypePrintFiscalMemoryByDate                TerminatorType = "2w"
	TerminatorTypePrintFiscalMemoryByClosureNumber       TerminatorType = "3w"
	TerminatorTypePrintDetailsMemoryAll                  TerminatorType = "4w"
	TerminatorTypePrintDetailsMemoryByDate               TerminatorType = "5w"
	TerminatorTypePrintDetailsMemoryByClosureNumber      TerminatorType = "6w"
	TerminatorTypeSetDateTime                            TerminatorType = "D"
	TerminatorTypeDisableXonXoff                         TerminatorType = "E"
	TerminatorTypeDisableXonXoff2                        TerminatorType = "1492E"
)

// GetTerminatorTypePaymentLight is used for custom payments
func GetTerminatorTypePaymentLight(paymentMethodCode string) (TerminatorType, error) {
	if len(paymentMethodCode) != 3 {
		return "", fmt.Errorf("paymentMethodCode must be 3 characters long")
	}
	return TerminatorType(paymentMethodCode + "T"), nil
}
