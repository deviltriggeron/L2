package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	currTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintln(os.Stderr, "NTP error:", err)
		os.Exit(1)
	}

	fmt.Println("Точное время:", currTime.Format(time.RFC1123))
}
