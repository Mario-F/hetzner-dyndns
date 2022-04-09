package setup

import (
	"fmt"

	"github.com/Mario-F/hetzner-dyndns/internal/hetzner"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
	"github.com/manifoldco/promptui"
)

var ipVersion network.IPVersion

func Cronjob() error {

	prompt := promptui.Select{
		Label: "Select ip version to use",
		Items: []string{"ipv6", "ipv4"},
	}

	_, option, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("Select failed %v\n", err)
	}

	switch option {
	case "ipv6":
		ipVersion = network.IPv6
	case "ipv4":
		ipVersion = network.IPv4
	}

	promptToken := promptui.Prompt{
		Label:    "Please enter your Hetzner token for access dns",
		Validate: checkHetznerToken,
	}

	hetznerToken, err := promptToken.Run()

	if err != nil {
		return fmt.Errorf("Prompt failed %v\n", err)
	}
	fmt.Print(hetznerToken)

	return nil
}

func checkHetznerToken(token string) error {
	if len(token) < 1 {
		return fmt.Errorf("Token is required")
	}

	return hetzner.CheckToken(token)
}
