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

	MIN1 = machine.GPIO5
	MIN2 = machine.GPIO7

	led ws2812.Device
)

const (
	i2cAddress = uint8(0x28)
)

func main() {
	// Configure sensor inputs
	S1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	S2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	S3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	S4.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Configure motor control pins
	MIN1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	MIN1.Low()
	MIN2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	MIN2.Low()

	// Configure neopixel
	machine.NEOPIXEL.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led = ws2812.New(machine.NEOPIXEL)
	led.WriteColors([]color.RGBA{colorBoot})

	go func() {
		for {
			if err := listenForIncomingI2CRequests(machine.I2C0, i2cAddress); err != nil {
				println("listenForIncomingI2CRequests failed: ", err)
				time.Sleep(time.Second)
			}
		}
	}()

	var m stateMachine
	m.Run()
}
