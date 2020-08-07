package cmd

import (
	"fmt"

	"github.com/Mario-F/hetzner-dyndns/internal/externalip"
)

func cmdMyIP() {
	res, err := externalip.GetExternalIP()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Your external IP is: %v\n", res)
}
