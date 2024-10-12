package main

import (
	"image/color"
	"time"
)

type state uint8

const (
	// In initial state, wait until a sensor is activated
	// The first sensor activated determined where the train is.
	stateInitial state = iota
	// Driving at full speed towards S2
	stateDrivingToS2
	// Driving at slow speed towards S1
	stateDrivingToS1
	// Waiting in S1
	stateStoppedInS1
	// Driving in full speed towards S3
	stateDrivingToS3
	// Driving in slow speed towards S4
	stateDrivingToS4
	// Waiting in S4
	stateStoppedInS4
)

// State of state machine
type stateMachine struct {
	// Motor controller
	mc *motorController
	// Current state
	state state
	// When the current state was reached
	lastChange time.Time
}

// Run the state machine
func (m *stateMachine) Run() {
	// Initialize
	m.changeState(stateInitial)
	waitTimeInS1 := time.Second * 5
	waitTimeInS4 := time.Second * 10
	lastRevActive := false

	for {
		s1Active := !S1.Get()
		s2Active := !S2.Get()
		s3Active := !S3.Get()
		s4Active := !S4.Get()
		revActive := !REV.Get()

		if revActive != lastRevActive {
			// REV switch changes
			m.mc.SetLocDirection(revActive)
			lastRevActive = revActive
			m.mc.FullStop()
			m.changeState(stateInitial)
		}

		switch m.state {
		// In initial state, wait until a sensor is activated
		// The first sensor activated determined where the train is.
		case stateInitial:
			if s1Active || s2Active {
				m.driveToS3()
			} else if s3Active || s4Active {
				m.driveToS2()
			}

			// Driving at full speed towards S2
		case stateDrivingToS2:
			if s2Active {
				// Reached S2, slow down
				m.driveToS1()
			} else if s1Active {
				// Already reached S1, stop
				m.waitInS1()
			}

			// Driving at slow speed towards S1
		case stateDrivingToS1:
			if s1Active {
				// Reached S1
				m.waitInS1()
			}

			// Waiting in S1
		case stateStoppedInS1:
			if time.Since(m.lastChange) >= waitTimeInS1 {
				m.driveToS3()
			}

			// Driving in full speed towards S3
		case stateDrivingToS3:
			if s3Active {
				// Reached S3, slow down
				m.driveToS4()
			} else if s4Active {
				// Already reached S4, stop
				m.waitInS4()
			}

			// Driving in slow speed towards S4
		case stateDrivingToS4:
			if s4Active {
				// Reached S1
				m.waitInS4()
			}

			// Waiting in S4
		case stateStoppedInS4:
			if time.Since(m.lastChange) >= waitTimeInS4 {
				m.driveToS2()
			}
		}

		// Wait a bit
		time.Sleep(time.Millisecond)
	}
}

// Start driving at full speed to S3
func (m *stateMachine) driveToS3() {
	m.mc.DriveForward(true)
	m.changeState(stateDrivingToS3)
}

// Start driving at full speed to S2
func (m *stateMachine) driveToS2() {
	m.mc.DriveBackward(true)
	m.changeState(stateDrivingToS2)
}

// Start driving at slow speed to S4
func (m *stateMachine) driveToS4() {
	m.mc.DriveForward(false)
	m.changeState(stateDrivingToS4)
}

// Start driving at slow speed to S1
func (m *stateMachine) driveToS1() {
	m.mc.DriveBackward(false)
	m.changeState(stateDrivingToS1)
}

// Stop driving and wait in S1
func (m *stateMachine) waitInS1() {
	m.mc.FullStop()
	m.changeState(stateStoppedInS1)
}

// Stop driving and wait in S4
func (m *stateMachine) waitInS4() {
	m.mc.FullStop()
	m.changeState(stateStoppedInS4)
}

// Set the new state of the state machine
func (m *stateMachine) changeState(newState state) {
	m.state = newState
	m.lastChange = time.Now()
	led.WriteColors([]color.RGBA{colorByState[m.state]})
}
