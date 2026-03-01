package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "NTP err: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(ntpTime)
}
