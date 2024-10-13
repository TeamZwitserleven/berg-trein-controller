package main

import (
	"fmt"
	"machine"
)

var (
	// Current firward version
	version = []byte{0, 1, 0} // Major.Minor.Patch
)

// Single i2c message sent to the incoming i2c port
type incomingI2CEvent struct {
	Event    machine.I2CTargetEvent
	HasValue bool
	Value    uint8
}

// Listen for incoming I2C requests.
// This listeners impersonates a pcf8574 devices.
// Data bits:
// - P0: Reverse mode is active (ro)
// - P1: S1 is active (ro)
// - P2: S2 is active (ro)
// - P3: S3 is active (ro)
// - P4: S4 is active (ro)
// - P5: Loc is driving (ro)
func listenForIncomingI2CRequests(i2c *machine.I2C, i2cAddress uint8, m *stateMachine) error {
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
		var replyBuf [1]uint8
		for {
			select {
			case evt := <-events:
				// Handle event
				switch evt.Event {
				case machine.I2CReceive:
					// Ignore
				case machine.I2CRequest:
					// Send current state
					replyBuf[0] = m.GetI2CByte()
					i2c.Reply(replyBuf[:])
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
			Event:    evt,
			HasValue: count >= 1,
			Value:    buf[0],
		}
	}
}
