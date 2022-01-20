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
	aprs.SwName = "wxigate-v"
	aprs.SwVers = Version

	wx := aprs.Wx{
		Lat:  latitude,
		Lon:  longitude,
		Type: comment,
	}

	wx.Timestamp = time.Now() // BUG(low) use "dateutc" in query

	query := req.URL.Query()

	for k, v := range query {
		// fmt.Printf("k/v: %v/%v\n", k, v)

		switch k {
		case "baromrelin":
			baromrelin, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.Altimeter = baromrelin
		case "tempf":
			temp, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.Temp = int(temp)
		case "humidity":
			humidity, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.Humidity = humidity
		case "hourlyrainin":
			hourlyrainin, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.RainLastHour = hourlyrainin
		case "24hourrainin":
			hourlyrainin24, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.RainLast24Hours = hourlyrainin24
		case "dailyrainin":
			dailyrainin, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.RainToday = dailyrainin
		case "solarradiation":
			solarradiation, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.SolarRad = int(solarradiation)
		case "winddir":
			winddir, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.WindDir = winddir
		case "windgustmph":
			windgustmph, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
			}
			wx.WindGust = int(windgustmph)
		case "windspeedmph":
			windspeedmph, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v", err) // BUG(medium) change
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

// BUG(medium-high) update
func errorHandler(w http.ResponseWriter, req *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404")
	}

	fmt.Println(req)
}

// BUG(medium-high) update
func catchall(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "welcome home")
}
