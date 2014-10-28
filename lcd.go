// Package serial_lcd controls an Adafruit 16x2 serial-backpack LCD display, typically connected over USB.
//
// This package is specifically designed to work with the Adafruit serial
// backpack LCD kit (http://www.adafruit.com/products/784).
//
// Typical usage is (with all error handling omitted):
//
//   lcd, _ := serial_lcd.Open("COM2", 9600) // or "/dev/tty.usbmodem1451"
//   defer lcd.Close()
//   lcd.SetSize(16,2)
//   lcd.SetBrightness(255)
//   lcd.SetContrast(200)
//   lcd.SetCursor(UNDERLINE_CURSOR_OFF, BLOCK_CURSOR_OFF)
//   lcd.SetBG(0,0,255) // R,G,B
//   lcd.Clear()
//   lcd.Home()
//   fmt.Fprint(lcd, "Hi there!")
//
package serial_lcd

import (
	"io"

	"github.com/tarm/goserial"
)

type UnderlineCursorState uint8
type BlockCursorState uint8
type AutoscrollState uint8

type LCD struct{ io.ReadWriteCloser }

func Open(port string, baud int) (LCD, error) {
	s, err := serial.OpenPort(&serial.Config{Name: port, Baud: baud})
	return LCD{s}, err
}

// dropN ignores the number of bytes written and just returns the error.
func dropN(n int, e error) error { return e }

// Raw writes a series of raw bytes to the LCD.
func (l LCD) Raw(bytes ...byte) error { return dropN(l.Write(bytes)) }

// SetBG sets the background color.  The RGB values should each be 0-255.
func (l LCD) SetBG(r, g, b uint8) error { return l.Raw(COMMAND, SET_RGB_BACKLIGHT_COLOR, r, g, b) }

// Off turns the LCD backlight off.
func (l LCD) Off() error { return l.Raw(COMMAND, BACKLIGHT_OFF) }

// On turns the LCD backlight on.
func (l LCD) On() error { return l.Raw(COMMAND, BACKLIGHT_ON, 0) }

// SetBrightness sets the LCD backlight brightness.  0-255 where 255 is the brightest.
func (l LCD) SetBrightness(b uint8) error { return l.Raw(COMMAND, BRIGHTNESS, b) }

// SetContrast sets the LCD backlight contrast. 0-255, usually 200 is a nice value.
func (l LCD) SetContrast(c uint8) error { return l.Raw(COMMAND, CONTRAST, c) }

// Autoscrolls determines how the LCD handles more text than fits on the
// display.  When on, if more text is received than fits it will immediately be
// scrolled so that the newest text is always at the bottom.  When off, as more
// text is received the display wraps around to the beginning.
func (l LCD) SetAutoscroll(a AutoscrollState) error { return l.Raw(COMMAND, byte(a)) }
func (l LCD) SetSize(cols, rows uint8) error        { return l.Raw(COMMAND, SET_LCD_SIZE, cols, rows) }
func (l LCD) Clear() error                          { return l.Raw(COMMAND, CLEAR) }

func (l LCD) SetCursor(u UnderlineCursorState, b BlockCursorState) error {
	return l.Raw(COMMAND, byte(u), COMMAND, byte(b))
}

// Move the cursor home (to 1,1).
func (l LCD) Home() error { return l.Raw(COMMAND, GO_HOME) }

// Set the cursor position.  Row/col number starts at 1,1.
func (l LCD) MoveTo(col, row uint8) error { return l.Raw(COMMAND, SET_CURSOR_POSITION, col, row) }
func (l LCD) MoveForward() error          { return l.Raw(COMMAND, CURSOR_FORWARD) }
func (l LCD) MoveBack() error             { return l.Raw(COMMAND, CURSOR_BACK) }

func (l LCD) CreateCustomChar(spot uint8, c Char) error {
	return l.Raw(append([]byte{COMMAND, CREATE_CUSTOM_CHARACTER, spot}, c[:]...)...)
}

// Characters are 5x8 pixels.  The first 5 bits of each byte defines the pixels
// for that row.
type Char [8]byte

// MakeChar converts an array of 5 character long string lines into the Char
// byte array defining a single char.  Any symbol other than a " " (space) or
// "." will be an ON pixel.
//
// For example:
//   heart := MakeChar([8]string{
//   	".....",
//   	".*.*.",
//   	"*.*.*",
//   	"*...*",
//   	"*...*",
//   	".*.*.",
//   	"..*..",
//   	".....",
//   })
//   lcd.CreateCustomChar(0, heart)
//
func MakeChar(lines [8]string) Char {
	var charDef Char
	for i, line := range lines {
		var pixels byte
		for _, c := range line {
			pixels = pixels << 1
			if c != '.' && c != ' ' {
				pixels |= 1
			}
		}
		charDef[i] = pixels
	}
	return charDef
}

const (
	// All commands start with the COMMAND byte.
	COMMAND = 0xFE

	// ---------------------------------------------------------------
	// Basic commands:
	BACKLIGHT_ON  = 0x42 // Turns the backlight on.  expect extra arg that is ignored.
	BACKLIGHT_OFF = 0x46 // Turns the backlight off.
	BRIGHTNESS    = 0x99 // Set brightness: expects arg for brightness 0-255
	CONTRAST      = 0x91 // Set contrast: expects arg for contrast 0-255
	CLEAR         = 0x58 // Clear the display.

	// This will make it so when text is received and there's no more space on
	// the display, the text will automatically 'scroll' so the second line
	// becomes the first line, etc. and new text is always at the bottom of the
	// display.
	AUTOSCROLL_ON = AutoscrollState(0x51)
	// This will make it so when text is received and there's no more space on
	// the display, the text will wrap around to start at the top of the
	// display.
	AUTOSCROLL_OFF = AutoscrollState(0x52)

	// after sending this command, write up to 32 characters (for 16x2) or up to
	// 80 characters (for 20x4) that will appear as the splash screen during
	// startup. If you don't want a splash screen, write a bunch of spaces.
	SET_STARTUP_SPLASH = 0x40

	// ---------------------------------------------------------------
	// Moving and changing the cursor:

	// set the position of text entry cursor. Column and row numbering starts
	// with 1 so the first position in the very top left is (1, 1)
	SET_CURSOR_POSITION = 0x47
	// place the cursor at location (1, 1)
	GO_HOME = 0x48
	// move cursor back one space, if at location (1,1) it will 'wrap' to the
	// last position.
	CURSOR_BACK = 0x4C
	// move cursor back one space, if at the last location location it will
	// 'wrap' to the (1,1) position.
	CURSOR_FORWARD = 0x4D
	// turn on the underline cursor
	UNDERLINE_CURSOR_ON = UnderlineCursorState(0x4A)
	// turn off the underline cursor
	UNDERLINE_CURSOR_OFF = UnderlineCursorState(0x4B)
	// turn on the blinking block cursor
	BLOCK_CURSOR_ON = BlockCursorState(0x53)
	// turn off the blinking block cursor
	BLOCK_CURSOR_OFF = BlockCursorState(0x54)

	// ---------------------------------------------------------------
	// RGB Backlight and LCD size

	// Sets the backlight to the red, green and blue component colors. The
	// values of can range from 0 to 255 (one byte). This is saved to EEPROM.
	// Each color R, G, and B is represented by a byte following the command.
	// Color values range from 0 to 255. To set the backlight to Red, the
	// command is 0xFE 0xD0 0x255 0x0 0x0. Blue is 0xFE 0xD0 0x0 0x0 0x255.
	// White is 0xFE 0xD0 0x255 0x255 0x255.
	SET_RGB_BACKLIGHT_COLOR = 0xD0
	// You can configure the backpack to what size display is attached. This is
	// saved to EEPROM so you only have to do it once.Custom Characters
	SET_LCD_SIZE = 0xD1
	// this will create a custom character in spot # can be between 0 and 7 (8
	// spots). 8 bytes are sent which indicate how the character should appear
	CREATE_CUSTOM_CHARACTER = 0x4E
	// this will save the custom character to EEPROM bank for later use. There
	// are 4 banks and 8 locations per bank.
	SAVE_CUSTOM_CHARACTER_TO_EEPROM_BANK = 0xC1
	// this will load all 8 characters saved to an EEPROM bank into the LCD's
	// memoryGeneral Purpose Output
	LOAD_CUSTOM_CHARACTERS_FROM_EEPROM_BANK = 0xC0
)
