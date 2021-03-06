// A tiny web server that allows interactively manipulating the LCD.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/augustoroman/serial_lcd"
	"github.com/go-martini/martini"
	"github.com/pwaller/go-hexcolor"
)

func main() {
	port := flag.String("port", "/dev/tty.usbmodem1451", "COM port that LCD is on.")
	baud := flag.Int("baud", 9600, "Baud rate to communicate at.")
	addr := flag.String("addr", ":12000", "Web address to bind to.")
	flag.Parse()
	lcd, err := serial_lcd.Open(*port, *baud)
	if err != nil {
		log.Fatal(err)
	}
	lcd.On()
	lcd.SetSize(16, 2)
	lcd.Clear()
	lcd.Home()
	fmt.Fprintf(lcd, "Hi there!")

	s := &server{lcd}

	m := martini.Classic()
	m.Handlers(martini.Recovery())
	m.Get("/", func() string { return home })
	m.Post("/set", s.Set)
	http.ListenAndServe(*addr, m)
}

type server struct{ lcd serial_lcd.LCD }

func (s *server) Set(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if b, ok := getByte("brightness", r.Form); ok {
		s.lcd.SetBrightness(b)
	}
	if c, ok := getByte("contrast", r.Form); ok {
		s.lcd.SetContrast(c)
	}
	if r, g, b, ok := getRGB(r.Form); ok {
		s.lcd.SetBG(r, g, b)
	}
	if a, ok := getText("autoscroll", r.Form); ok {
		s.lcd.SetAutoscroll(a == "true")
	}

	if on, ok := getText("on", r.Form); ok {
		s.lcd.SetOn(on == "true")
	}

	if vals, ok := r.Form["txt"]; ok && len(vals) == 1 {
		s.lcd.Clear()
		s.lcd.Home()
		fmt.Fprint(s.lcd, vals[0])
	}
}

func getByte(key string, vals url.Values) (byte, bool) {
	if txt, ok := getText(key, vals); ok {
		num, err := strconv.ParseUint(txt, 10, 8)
		if err != nil {
			return 0, false
		}
		return uint8(num), true
	}
	return 0, false
}

func getRGB(vals url.Values) (r, g, b byte, ok bool) {
	if txt, ok := getText("background", vals); ok {
		r, g, b, _ := hexcolor.HexToRGBA(hexcolor.Hex(txt))
		return r, g, b, true
	}
	return 0, 0, 0, false
}

func getText(key string, vals url.Values) (string, bool) {
	if val, ok := vals[key]; ok && len(val) == 1 {
		return val[0], true
	}
	return "", false
}

var home = `<html>
<body>
LCD Control:
<hr>
Text:<br><textarea rows=2 cols=16 oninput="set({txt:this.value})" onchange="set({txt:this.value})">Hi there!</textarea><br>
Brightness: <input min=0 max=255 step=1 type=range
  oninput="set({brightness:this.value})" onchange="set({brightness:this.value})"><br>
Contrast: <input min=0 max=255 step=1 type=range
  oninput="set({contrast:this.value})" onchange="set({contrast:this.value})"><br>
Background: <input type=color oninput="set({background:this.value})" onchange="set({background:this.value})"><br>
Autoscroll: <input type=checkbox onchange="set({autoscroll:this.checked})"><br>
On: <input type=checkbox onchange="set({on:this.checked})"><br>
<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
<script>
function set(vals) { $.post("/set", vals); }
</script>
</body>
</html>
`
