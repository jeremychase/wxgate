package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ebarkie/aprs"
)

func v1(w http.ResponseWriter, req *http.Request) {
	wx := aprs.Wx{
		Lat:  41.38265020114848,  // BUG(high) flag in
		Lon:  -72.43430504603869, // BUG(high) flag in
		Type: "DvsVP2+",          // BUG(medium) flag
	}
	wx.Timestamp = time.Now() // BUG(low) use "dateutc" in query

	wx.Altimeter = 42.0 // BUG(high) flag in

	query := req.URL.Query()

	for k, v := range query {
		// fmt.Printf("k/v: %v/%v\n", k, v)

		switch k {
		case "tempf":
			temp, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.Temp = int(temp)
		case "humidity":
			humidity, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.Humidity = humidity
		case "hourlyrainin":
			hourlyrainin, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.RainLastHour = hourlyrainin
		case "24hourrainin":
			hourlyrainin24, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.RainLast24Hours = hourlyrainin24
		case "dailyrainin":
			dailyrainin, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.RainToday = dailyrainin
		case "solarradiation":
			solarradiation, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.SolarRad = int(solarradiation)
		case "winddir":
			winddir, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.WindDir = winddir
		case "windgustmph":
			windgustmph, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.WindGust = int(windgustmph)
		case "windspeedmph":
			windspeedmph, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(high) fix
			}
			wx.WindSpeed = int(windspeedmph)
		}
	}

	fmt.Printf("sending, temp(%d): %v\n", wx.Temp, wx.String())

	f := aprs.Frame{
		Dst:  aprs.Addr{Call: "APRS"},
		Src:  aprs.Addr{Call: fmt.Sprintf("%s-%s", callsign, ssid)},
		Path: aprs.Path{aprs.Addr{Call: "TCPIP", Repeated: true}},
		Text: wx.String(),
	}
	err := f.SendIS("tcp://cwop.aprs.net:14580", -1) //BUG(medium) flag
	if err != nil {
		log.Printf("Upload error: %s", err) // BUG(medium) handle
	}
}

// BUG(high) gross
func errorHandler(w http.ResponseWriter, req *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404")
	}

	fmt.Println(req)
}

// BUG(high) gross
func catchall(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "welcome home")
}
