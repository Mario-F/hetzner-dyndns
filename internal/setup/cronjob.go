package setup

import (
	"fmt"

	"github.com/Mario-F/hetzner-dyndns/internal/hetzner"
	"github.com/Mario-F/hetzner-dyndns/internal/network"
	"github.com/manifoldco/promptui"
)

var ipVersion network.IPVersion

func SetVersion(version network.IPVersion) {
	ipVersion = version
}

func Cronjob() error {

	prompt := promptui.Prompt{
		Label:    "Please enter your Hetzner token for access dns",
		Validate: checkHetznerToken,
	}

	hetznerToken, err := prompt.Run()

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
