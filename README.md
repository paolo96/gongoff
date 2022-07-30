# GoOn-GoOff

GoOn-GoOff (gongoff) is a go implementation of the [Epson Xon-Xoff](https://download.epson-biz.com/modules/pos/index.php?page=single_doc&cid=6735&pcat=51&pid=5811) protocol.

It's used to communicate with fiscal printers that support the Xon-Xoff protocol.

Supported connections: serial and network

Supported platforms: Windows, Linux, macOS, Android (only network), iOS (only network)

## How to use

There are 2 ways to use GoOn-GoOff to control a printer:

### Documents

Documents are a set of commands commonly sent to a printer together.

Supported documents are:
* DocumentCommercial (fiscal receipt)
* DocumentManagement
* DocumentCommercialReturn
* DocumentCommercialCancellation
* DocumentPOSReturn
* DocumentPOSCancellation
* DocumentInvoice
* DocumentCommercialWithInvoice

### Commands

Commands are a less abstract way to use the library and allow to use all the features of Xon-Xoff.

There are a handful of commands with predefined implementations. 
All the other commands can be created using a GenericCommand with the appropriate parameters and terminator. 

## Usage examples

#### Printing a test document through serial port (RS-232)
```go
// Create a SerialPrinter object and open the serial port COM3. 
printer := gongoff.NewSerialPrinter("COM3")
err := printer.Open()
if err != nil {
    panic(err)
}
defer printer.Close()

// Create a management document.
doc := NewDocumentManagement([]string{"test", "test2", "test3"})

// Print it.
err = printer.PrintDocument(doc)
if err != nil {
    panic(err)
}
```

#### Printing a commercial document (fiscal receipt) through network
```go
// Create a NetworkPrinter object and open the connection.
printer := gongoff.NewNetworkPrinter("192.168.1.100", 9100)
err := printer.Open()
if err != nil {
    panic(err)
}
defer printer.Close()

// Create a receipt document.
// In this example, the document is a 7,50 euro receipt for a product called "Bread" paid with cash.
testProduct := "BREAD"
commandProduct := NewCommandProduct(750, &testProduct, nil, nil)
commandPayment, err := NewCommandPayment(TerminatorTypePaymentCash, nil, nil)
if err != nil {
	panic(err)
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
    nil,
)

// Print it.
err = printer.PrintDocument(commercialDoc)
if err != nil {
    panic(err)
}
```

#### Showing test message on the display
```go
// Suppose the printer object is already created and opened.
// Since there is no document that fits our needs, we use an array of commands.
commandFirstLine := NewCommandDisplayMessage("First line", 1)
commandSecondLine := NewCommandDisplayMessage("Second line", 2)

// Print the commands.
err := printer.PrintCommands([]Command{commandFirstLine, commandSecondLine})
if err != nil {
    panic(err)
}
```

#### Print financial report with zeroing
```go
// Suppose the printer object is already created and opened.
// Since there is neither a document nor a command that fits our needs, we generate a generic command with the right terminator.
var commands []gongoff.Command
commands = append(commands, 
	gongoff.NewCommandGeneric(
            []gongoff.Data{},
            gongoff.Terminator{nil, gongoff.TerminatorTypeFinancialReportZeroing}, 
	),
)

// Print the commands.
err := printer.PrintCommands(commands)
if err != nil {
    panic(err)
}
```