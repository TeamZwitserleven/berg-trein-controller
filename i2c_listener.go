package main

import (
	"fmt"
	"machine"
)

var (
	// Current firward version
	version = []byte{0, 1, 0} // Major.Minor.Patch
)

const (
	// Register addresses
	RegVersionMajor   = 0x00 // No input, returns 1 version
	RegVersionMinor   = 0x01 // No input, returns 1 version
	RegVersionPatch   = 0x02 // No input, returns 1 version
	RegCarSensorCount = 0x03 // No input, returns 1 byte giving the number of detected car sensor bits (0..8)
	RegI2COutputCount = 0x04 // No input, returns 1 byte giving the number of detected I2C binary output pins (0, 8, 16, ..., 256)
	RegCarSensorState = 0x10 // No input, returns 1 byte with 8-bit car detection sensor state
	RegOutput         = 0x20 // 1 byte input, targeting 8 on-pcb output pins
	RegOutputI2C0     = 0x21 // 1 byte input, targeting 8 output pins on PCF8574 output device 0
	RegOutputI2C1     = 0x22 // 1 byte input, targeting 8 output pins on PCF8574 output device 1
	RegOutputI2C2     = 0x23 // 1 byte input, targeting 8 output pins on PCF8574 output device 2
	RegOutputI2C3     = 0x24 // 1 byte input, targeting 8 output pins on PCF8574 output device 3
	RegOutputI2C4     = 0x25 // 1 byte input, targeting 8 output pins on PCF8574 output device 4
	RegOutputI2C5     = 0x26 // 1 byte input, targeting 8 output pins on PCF8574 output device 5
	RegOutputI2C6     = 0x27 // 1 byte input, targeting 8 output pins on PCF8574 output device 6
	RegOutputI2C7     = 0x28 // 1 byte input, targeting 8 output pins on PCF8574 output device 7
	RegConfigurePWM0  = 0x30 // 1 byte input, pwm-value (0-256) of pin 0
	RegConfigurePWM1  = 0x31 // 1 byte input, pwm-value (0-256) of pin 1
	RegConfigurePWM2  = 0x32 // 1 byte input, pwm-value (0-256) of pin 2
	RegConfigurePWM3  = 0x33 // 1 byte input, pwm-value (0-256) of pin 3
	RegConfigurePWM4  = 0x34 // 1 byte input, pwm-value (0-256) of pin 4
	RegConfigurePWM5  = 0x35 // 1 byte input, pwm-value (0-256) of pin 5
	RegConfigurePWM6  = 0x36 // 1 byte input, pwm-value (0-256) of pin 6
	RegConfigurePWM7  = 0x37 // 1 byte input, pwm-value (0-256) of pin 7

	pwmPeriod = uint64(1e9) / 60
)

// Single i2c message sent to the incoming i2c port
type incomingI2CEvent struct {
	Event       machine.I2CTargetEvent
	HasRegister bool
	Register    uint8
	HasValue    bool
	Value       uint8
}

// Listen for incoming I2C requests.
func listenForIncomingI2CRequests(i2c *machine.I2C, i2cAddress uint8) error {
	// Configure i2c bus as target
	if err := i2c.Configure(machine.I2CConfig{
		Mode: machine.I2CModeTarget,
	}); err != nil {
		return fmt.Errorf("Failed to configure i2c bus: %w", err)
	}

	// Start listening on the i2c bus
	if err := i2c.Listen(uint16(i2cAddress)); err != nil {
		return fmt.Errorf("Failed to listen on i2c bus: %w", err)
	}
	println("Listening on i2c address: ", i2cAddress)

	// Process events & status changes
	events := make(chan incomingI2CEvent)
	go func() {
		for {
			select {
			case evt := <-events:
				// Handle event
				switch evt.Event {
				case machine.I2CReceive:
					// TODO
					switch evt.Register {
					// TODO
					default:
						println("I2C:Receive: Invalid register ", evt.Register, evt.HasValue, evt.Value)
					}
				case machine.I2CRequest:
					// TODO
				case machine.I2CFinish:
					// No response needed
				}
			}
		}
	}()
	var buf [8]uint8
	for {
		// Wait for event
		evt, count, err := i2c.WaitForEvent(buf[:])
		if err != nil {
			return fmt.Errorf("Failed to wait for event: %w", err)
		}

		// Handle event
		events <- incomingI2CEvent{
			Event:       evt,
			HasRegister: count >= 1,
			Register:    buf[0],
			HasValue:    count >= 2,
			Value:       buf[1],
		}
	}
}
