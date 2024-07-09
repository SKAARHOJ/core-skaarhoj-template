package main

import (
	skconfig "github.com/SKAARHOJ/ibeam-lib-config"
)

/*
In this file we create a configstructure that holds devices and core settings.
This will be dumped as a toml file when the core is running
When it is registered correctly in main reactor will be able to adjust config remotely and restart the core

You are free to use other means of configuration if needed, but using ibeam-lib-config will ensure your core is the most compatible it can be
*/

// CoreConfig main config structure of the core
type CoreConfig struct {
	Devices []DeviceConfig `ibDispatch:"devices" ibDescription:"Configure your device settings here"`
}

// DeviceConfig configuration of an individual device
type DeviceConfig struct {
	skconfig.BaseDeviceConfig
	IP       string `ibDispatch:"ip" ibValidate:"ip" ibOrder:"1"`
	Port     uint16 `ibDispatch:"port" ibValidate:"port" ibOrder:"2" ibDefault:"8000"`
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
				Password: "somepassword",
			},
		},
	}
}
