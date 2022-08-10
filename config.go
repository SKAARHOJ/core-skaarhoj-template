package main

import (
	skconfig "github.com/SKAARHOJ/ibeam-lib-config"
)

// CoreConfig main config structure of the core
type CoreConfig struct {
	Devices []DeviceConfig `ibDispatch:"devices" ibDescription:"Configure your device settings here"`
}

// DeviceConfig configuration of an individual device
type DeviceConfig struct {
	skconfig.BaseDeviceConfig
	IP       string `ibValidate:"ip" ibOrder:"1"`                    // Markup for skaarOS webui
	Port     uint16 `ibValidate:"port" ibOrder:"2" ibDefault:"8000"` // Markup for skaarOS webui
	Username string `ibOrder:"3" ibDefault:"skaarhoj"`
	Password string `ibLabel:"Device Password" ibValidate:"password" ibDescription:"The password that is set on the device" ibOrder:"4"` // Markup for skaarOS webui
}

func defaultConfig() CoreConfig {
	return CoreConfig{
		Devices: []DeviceConfig{
			{
				BaseDeviceConfig: skconfig.BaseDeviceConfig{
					DeviceID:    1,
					Description: "Main PowerPDU 4C in Rack",
					Active:      true,
				},
				IP:       "192.168.10.28",
				Username: "write",
				Password: "24A42C39352B",
			},
		},
	}
}
