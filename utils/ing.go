package utils

//import "strings"

/*GetQRCode generates the content for a ING Mobiel Bankieren QR code
The content will be transformed into an QR image client side.
*/
func GetQRCode(ownerName string, ownerIBAN string, amount string) (QRCodeText string) {
	//if strings.Contains(strings.ToUpper(debtorIBAN), "INGB") {
	QRCodeText = "BCD\n002\n1\nSCT\n\n" + ownerName + "\n" + ownerIBAN + "\nEUR" + amount
	return QRCodeText
}
