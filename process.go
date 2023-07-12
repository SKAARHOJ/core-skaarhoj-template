package main

import (
	"context"
	"fmt"
	"net"
	"time"

	ib "github.com/SKAARHOJ/ibeam-corelib-go"
	pb "github.com/SKAARHOJ/ibeam-corelib-go/ibeam-core"
	b "github.com/SKAARHOJ/ibeam-corelib-go/paramhelpers"
	"github.com/jpillora/backoff"
	log "github.com/s00500/env_logger"
)

/*
process.go holds the main implementation of the core, communicating with the manager over the two channels received on initializing the corelib

Right in the start we register all devices that are defined by the config and start go routines for them.
Messages from the manager are then routed to the individual devices go routines

Each of the sub go routines per device handle their respective state, network connections and reconnection handling

*/

func processDevices(r *ib.IBeamParameterRegistry, config CoreConfig, fromManager <-chan *pb.Parameter, toManager chan<- *pb.Parameter) {
	stateChannels := make(map[uint32]chan *pb.Parameter)

	for _, deviceConfig := range config.Devices {
		if !deviceConfig.Active {
			continue
		}

		_, err := r.RegisterDevice(deviceConfig.DeviceID, deviceConfig.ModelID)
		if log.Should(err) {
			continue
		}
		stateChan := make(chan *pb.Parameter, 10)
		stateChannels[deviceConfig.DeviceID] = stateChan

		go handleConnection(deviceConfig, stateChan, toManager, r)
	}

	for parameter := range fromManager { //determin correct device and send it
		select {
		case stateChannels[parameter.Id.Device] <- parameter:
		default:
			log.Errorf("Device %d: Channel is full", parameter.Id.Device)
		}
	}
}

func handleConnection(config DeviceConfig, fromManager <-chan *pb.Parameter, toManager chan<- *pb.Parameter, r *ib.IBeamParameterRegistry) {
	reconnTimer := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: 2,
		Jitter: false,
	}
	did := config.DeviceID

	toManager <- b.Param(r.PID("connection"), did, b.Bool(false))

	port := config.Port
	if port == 0 {
		port = 8080 // Provide a default port for this device in case it has not been set
	}

	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.IP, port))
		if err != nil {
			log.Errorf("Could not connect to '%s': %v", config.IP, err)
			time.Sleep(reconnTimer.Duration())
			continue
		}

		log.Infof("Client '%s' connected as deviceID %d", config.IP, did)

		ctx, cancel := context.WithCancel(context.Background())
		reconnTimer.Reset() // Consider connection successful

		toManager <- b.Param(r.PID("connection"), did, b.Bool(true))

		// Process incoming
		go func() {
			for {
				var err error
				// TODO: Process incoming data here
				log.Info("Processing Incoming")

				toManager <- b.Param(r.PID("iris"), did, b.Int(1000))

				time.Sleep(time.Second * 5)
				if err != nil {
					break
				}
			}
			cancel()
			conn.Close()
		}()

		// Process outgoing
		go func() {
			for {
				select {
				case parameter := <-fromManager:
					packets := createPacket(parameter, r)
					for _, packet := range packets {
						if len(packet) == 0 {
							continue
						}
						_, err := conn.Write(packet)
						if err != nil {
							log.Error(err)
							cancel()
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		// Process polling or any other kind of cycling requests
		go func() {
			t := time.NewTicker(time.Millisecond * 500)
			for {
				select {
				case <-t.C:
					_, err := conn.Write([]byte{1, 2, 4})
					if err != nil {
						log.Error(err)
						cancel()
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		<-ctx.Done()
		conn.Close()
		toManager <- b.Param(r.PID("connection"), did, b.Bool(false))
		retryIn := reconnTimer.Duration()
		log.Warnf("Client '%s' disconnected, trying to reconnect in %s", config.IP, time.Duration(retryIn).String())
		time.Sleep(retryIn)
	}
}

func createPacket(parameter *pb.Parameter, r *ib.IBeamParameterRegistry) [][]byte {
	packets := make([][]byte, 0)
	for _, val := range parameter.Value {
		switch r.ParameterNameByID(parameter.Id.Parameter) {
		case "some_parameter":
			log.Info("Will send some_parameter to device")
			packets = append(packets, []byte{1, 2, byte(val.GetInteger())})
		case "":
			log.Error("Invalid ParameterName")
		}
	}
	return packets
}
