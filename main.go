package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

var (
	colorBoot    = color.RGBA{R: 165, G: 165, B: 165}
	colorByState = []color.RGBA{
		color.RGBA{R: 32, G: 32, B: 32},
		color.RGBA{R: 165, G: 0, B: 165},
		color.RGBA{R: 165, G: 0, B: 0},
		color.RGBA{R: 165, G: 165, B: 0},
		color.RGBA{R: 0, G: 165, B: 0},
		color.RGBA{R: 0, G: 255, B: 0},
		color.RGBA{R: 0, G: 0, B: 255},
	}
)

var (
	S1 = machine.GPIO29
	S2 = machine.GPIO28
	S3 = machine.GPIO27
	S4 = machine.GPIO26

	MIN1    = machine.GPIO5
	MIN2    = machine.GPIO7
	MIN1PWM = machine.PWM2
	MIN2PWM = machine.PWM3

	REV = machine.GPIO2

	led ws2812.Device
)

const (
	i2cAddress = uint8(0x28)
)

type pwm interface {
	// Configure enables and configures this PWM.
	Configure(config machine.PWMConfig) error
	// Channel returns a PWM channel for the given pin. If pin does
	// not belong to PWM peripheral ErrInvalidOutputPin error is returned.
	// It also configures pin as PWM output.
	Channel(pin machine.Pin) (channel uint8, err error)
	// SetPeriod updates the period of this PWM peripheral in nanoseconds.
	// To set a particular frequency, use the following formula:
	//
	//	period = 1e9 / frequency
	//
	// Where frequency is in hertz. If you use a period of 0, a period
	// that works well for LEDs will be picked.
	//
	// SetPeriod will try not to modify TOP if possible to reach the target period.
	// If the period is unattainable with current TOP SetPeriod will modify TOP
	// by the bare minimum to reach the target period. It will also enable phase
	// correct to reach periods above 130ms.
	SetPeriod(period uint64) error
	// Top returns the current counter top, for use in duty cycle calculation.
	//
	// The value returned here is hardware dependent. In general, it's best to treat
	// it as an opaque value that can be divided by some number and passed to Set
	// (see Set documentation for more information).
	Top() uint32
	// Set updates the channel value. This is used to control the channel duty
	// cycle, in other words the fraction of time the channel output is high (or low
	// when inverted). For example, to set it to a 25% duty cycle, use:
	//
	//	pwm.Set(channel, pwm.Top() / 4)
	//
	// pwm.Set(channel, 0) will set the output to low and pwm.Set(channel,
	// pwm.Top()) will set the output to high, assuming the output isn't inverted.
	Set(channel uint8, value uint32)
	// Enable enables or disables PWM peripheral channels.
	Enable(enable bool)
}

func main() {
	// Configure sensor inputs
	S1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	S2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	S3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	S4.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Configure reverse config input
	REV.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Configure motor control pins
	MIN1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	MIN1.Low()
	MIN2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	MIN2.Low()

	// Configure neopixel
	machine.NEOPIXEL.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led = ws2812.New(machine.NEOPIXEL)
	led.WriteColors([]color.RGBA{colorBoot})

	// Prepare motor controller
	mc := newMotorController()

	// Prepare state machine
	m := stateMachine{mc: mc}

	// Start listening to I2C
	go func() {
		for {
			if err := listenForIncomingI2CRequests(machine.I2C0, i2cAddress); err != nil {
				println("listenForIncomingI2CRequests failed: ", err)
				time.Sleep(time.Second)
			}
		}
	}()

	// Run state machine
	m.Run()
}
