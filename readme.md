# Hello

`wxigate` is a commandline tool that can send weather data from Ambient Weather stations into the APRS-IS network.

## Software Setup

1. Download the binary from the releases.
1. Unarchive
1. Run it, for example:
```
./wxigate -callsign YOUR_CALL -latitude 12.345 -longitude "-12.345"
```

## Weather Station setup

1. Open the 'awnet' application on on your phone.
1. Go to "Device List"
1. Select the station you want data from.
1. Click 'next' until on the 'Customized' view.
1. Enter the IP or hostname of the machine running `wxigate`.
1. In `Path` enter `/wxigate/awp/v1?`
1. Click 'Save'.
1. You should start seeing packets in the stdout of `wxigate`.

## Status

This tool was built for fun. Use at your own risk.

## Building

1. Install `go`
1. `go build .`