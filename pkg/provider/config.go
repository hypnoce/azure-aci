package provider

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/virtual-kubelet/node-cli/provider"
)

type providerConfig struct {
	ResourceGroup   string
	Region          string
	OperatingSystem string
	CPU             string
	Memory          string
	Pods            string
	SubnetName      string
	SubnetCIDR      string
}

func (p *ACIProvider) loadConfig(r io.Reader) error {
	var config providerConfig
	if _, err := toml.DecodeReader(r, &config); err != nil {
		return err
	}
	p.region = config.Region
	p.resourceGroup = config.ResourceGroup

	// Default to 20 mcpu
	p.cpu = "20"
	if config.CPU != "" {
		p.cpu = config.CPU
	}
	// Default to 100Gi
	p.memory = "100Gi"
	if config.Memory != "" {
		p.memory = config.Memory
	}
	// Default to 20 pods
	p.pods = "20"
	if config.Pods != "" {
		p.pods = config.Pods
	}

	// Default to Linux if the operating system was not defined in the config.
	if config.OperatingSystem == "" {
		config.OperatingSystem = provider.OperatingSystemLinux
	} else {
		// Validate operating system from config.
		ok := provider.ValidOperatingSystems[config.OperatingSystem]
		if !ok {
			return fmt.Errorf("%q is not a valid operating system, try one of the following instead: %s", config.OperatingSystem, strings.Join(provider.ValidOperatingSystems.Names(), " | "))
		}
	}

	// default subnet name
	if config.SubnetName != "" {
		p.providernetwork.SubnetName = config.SubnetName
	}
	if config.SubnetCIDR != "" {
		if config.SubnetName == "" {
			return fmt.Errorf("subnet CIDR is set but no subnet name provided, must provide a subnet name in order to set a subnet CIDR")
		}
		if _, _, err := net.ParseCIDR(config.SubnetCIDR); err != nil {
			return fmt.Errorf("error parsing provided subnet CIDR: %v", err)
		}
	}

	p.operatingSystem = config.OperatingSystem
	return nil
}
