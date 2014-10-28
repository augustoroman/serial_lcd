Serial LCD
==========

[![GoDoc](https://godoc.org/github.com/augustoroman/serial_lcd?status.svg)](http://godoc.org/github.com/augustoroman/serial_lcd)

A library for communicating with an serial LCD controller.  This is specifically geared towards the Adafruit 16x2 serial backpack LCD kit (http://www.adafruit.com/products/784).

Simple usage:

```go
lcd, err := serial_lcd.Open("COM2", 9600)
lcd.SetSize(16,2)
lcd.SetBG(255,0,0) // R,G,B
lcd.Clear()
fmt.Fprint(lcd, "Hi there!")
```
