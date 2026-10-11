package main

import (
	"fmt"
)

func main() {
	payment := Payment{ID: "pay_1042", AmountCents: 125000}
	fmt.Println("payment requires confirmation:", NeedsConfirmation(payment))
}
