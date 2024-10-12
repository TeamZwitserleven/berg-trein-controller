package main

import "machine"

// Motor controller state
type motorController struct {
	min1Pin machine.Pin
	min2Pin machine.Pin
	min1PWM pwm
	min2PWM pwm
}

const (
	motorFreq   = 20_000 // 20KHz
	motorPeriod = uint64(1e9) / motorFreq
)

// Construct and initialize a motor controller wih default values
func newMotorController() *motorController {
	return &motorController{
		min1Pin: MIN1,
		min2Pin: MIN2,
		min1PWM: MIN1PWM,
		min2PWM: MIN2PWM,
	}
}

// Set the motor into driving forward
func (m *motorController) DriveForward(fullSpeed bool) {
	m.min1PWM.Configure(machine.PWMConfig{Period: motorPeriod})
	ch, _ := m.min1PWM.Channel(m.min1Pin)
	m.min1PWM.Set(ch, speedValue(m.min1PWM, fullSpeed))
	m.min1PWM.Enable(true)

	m.min2Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	m.min2Pin.Low()
}

// Set the motor into driving backwards
func (m *motorController) DriveBackward(fullSpeed bool) {
	m.min2PWM.Configure(machine.PWMConfig{Period: motorPeriod})
	ch, _ := m.min2PWM.Channel(m.min2Pin)
	m.min2PWM.Set(ch, speedValue(m.min2PWM, fullSpeed))
	m.min2PWM.Enable(true)

	m.min1Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	m.min1Pin.Low()
}

// Set the motor into full stop
func (m *motorController) FullStop() {
	m.min1Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	m.min2Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	m.min1Pin.Low()
	m.min2Pin.Low()
}

// Set the pins according to the given loc direction
func (m *motorController) SetLocDirection(reverse bool) {
	if reverse {
		m.min1Pin = MIN2
		m.min2Pin = MIN1
		m.min1PWM = MIN2PWM
		m.min2PWM = MIN1PWM
	} else {
		m.min1Pin = MIN1
		m.min2Pin = MIN2
		m.min1PWM = MIN1PWM
		m.min2PWM = MIN2PWM
	}
}

func speedValue(pwm pwm, fullSpeed bool) uint32 {
	if fullSpeed {
		return (pwm.Top() / 3) * 2
	} else {
		return pwm.Top() / 3
	}
}
