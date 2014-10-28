// A simple tool to test communication with the lcd.
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/augustoroman/serial_lcd"
)

func main() {
	port := flag.String("port", "/dev/tty.usbmodem1451", "COM port that LCD is on.")
	baud := flag.Int("baud", 9600, "Baud rate to communicate at.")
	flag.Parse()
	lcd, err := serial_lcd.Open(*port, *baud)
	if err != nil {
		log.Fatal(err)
	}
	defer lcd.Close()

	lcd.Clear()
	lcd.On()

	lcd.SetCursor(serial_lcd.UNDERLINE_CURSOR_OFF, serial_lcd.BLOCK_CURSOR_OFF)
	lcd.MoveTo(8, 2)
	fmt.Fprint(lcd, "xyz")
	delay(1000)

	go loop(lcd)

	for i := 0; i < 100; i++ {
		x := uint8(i % 16)
		y := uint8((i / 16) % 2)
		lcd.MoveTo(x, y)
		lcd.Raw('*' + uint8(i)%5)
		delay(10)
	}

	delay(250)

	setup(lcd)

	delay(3000)

	lcd.SetBG(100, 0, 100)
}

func setup(lcd serial_lcd.LCD) {
	lcd.SetSize(16, 2)
	delay(10)

	lcd.SetContrast(200)
	delay(10)

	lcd.SetBrightness(255)
	delay(10)

	// turn off cursors
	lcd.SetCursor(serial_lcd.UNDERLINE_CURSOR_OFF, serial_lcd.BLOCK_CURSOR_OFF)
	delay(10)

	// create a custom character
	heart := serial_lcd.MakeChar([8]string{
		".....",
		".*.*.",
		"*.*.*",
		"*...*",
		"*...*",
		".*.*.",
		"..*..",
		".....",
	})
	lcd.CreateCustomChar(0, heart)
	delay(10) // we suggest putting delays after each command

	// clear screen
	lcd.Clear()
	delay(10) // we suggest putting delays after each command

	// go 'home'
	lcd.Home()
	delay(10) // we suggest putting delays after each command

	fmt.Fprint(lcd, " We \x00 Arduino!")
	fmt.Fprint(lcd, "     - Adafruit")
	delay(10) // we suggest putting delays after each command
}

func loop(lcd serial_lcd.LCD) {
	var red, green, blue byte
	// adjust colors
	for red = 0; red < 255; red++ {
		lcd.SetBG(red, 0, 255-red)
		delay(1) // give it some time to adjust the backlight!
	}

	for green = 0; green < 255; green++ {
		lcd.SetBG(255-green, green, 0)
		delay(1) // give it some time to adjust the backlight!
	}

	for blue = 0; blue < 255; blue++ {
		lcd.SetBG(0, 255-blue, blue)
		delay(1) // give it some time to adjust the backlight!
	}
}

func delay(ms int) { time.Sleep(time.Duration(ms) * time.Millisecond) }
