package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	filedialog "github.com/sqweek/dialog"
	"github.com/sydneyowl/g90updatefw-gui/lib/g90updatefw"
	goserial "go.bug.st/serial"
	"os"
	"strconv"
	"time"
)

var (
	firmware = binding.NewString()
	device   = binding.NewString()
)

type fwgui struct {
	logEntry *widget.Entry
	window   fyne.Window
}

func (f *fwgui) log(text string) {
	curTime := time.Now().Format("15:04:05")
	f.logEntry.Append(fmt.Sprintf("[%s] %s\n", curTime, text))
	f.logEntry.CursorRow = len(f.logEntry.Text) - 1
}

func (f *fwgui) loadUI(myApp fyne.App) {
	myWindow := myApp.NewWindow(fmt.Sprintf("g90updatefw-gui %s", Version))
	myWindow.Resize(fyne.NewSize(600, 400))
	myWindow.SetFixedSize(true)

	f.logEntry = widget.NewMultiLineEntry()
	f.logEntry.Text = `Welcome to g90updatefw-gui!
This is an open-source project based on g90updatefw!
Repo address: github.com/sydneyowl/g90updatefw-gui
==========
This program is designed to write a firmware file to a Xiegu radio.
It can be used to update either the main unit or the display unit.
Choose a firmware file and the serial port connected to the Xiegu radio, then touch "start"!
==========
`

	firmwareEntry := widget.NewEntryWithData(firmware)
	firmwareEntry.SetPlaceHolder("Select a firmware file...(*.xgf, *.aes)")
	selectFirmwareButton := widget.NewButton("Select file", func() {
		filePath, err := filedialog.File().Filter("*.xgf, *.aes", "xgf", "aes").Load()
		if err != nil {
			f.log(err.Error())
			return
		}
		firmware.Set(filePath)
	})

	firmwareContainer := container.NewBorder(nil, nil, nil, selectFirmwareButton, firmwareEntry)

	// Fetch all ports
	ports, err := goserial.GetPortsList()
	if err != nil {
		f.log("Error: " + err.Error())
	}
	portSelect := widget.NewSelect(ports, func(s string) {
		device.Set(s)
	})
	portSelect.PlaceHolder = "Please choose a port..."
	refreshDeviceButton := widget.NewButton("Refresh", func() {
		portSelect.ClearSelected()
		ports, err := goserial.GetPortsList()
		if err != nil {
			f.log("Error: " + err.Error())
		}
		portSelect.Options = ports
	})

	portContainer := container.NewBorder(nil, nil, nil, refreshDeviceButton, portSelect)

	progressBar := widget.NewProgressBar()

	var startButton *widget.Button

	startButton = widget.NewButton("Start", func() {
		go func() {
			f.log("Starting..")
			startButton.Disable()
			defer startButton.Enable()

			devicePath, _ := device.Get()
			if devicePath == "" {
				dialog.NewError(errors.New("Please select a port!"), myWindow).Show()
				return
			}
			serial, err := g90updatefw.SerialOpen(devicePath, 115200)
			if err != nil {
				f.log("Error: " + err.Error())
				dialog.NewError(err, myWindow).Show()
				return
			}
			defer serial.Close()
			firmwarePath, _ := firmware.Get()
			if firmwarePath == "" {
				dialog.NewError(errors.New("Please select a firmware!"), myWindow).Show()
				return
			}
			data, err := os.ReadFile(firmwarePath)
			if err != nil {
				f.log("Error: " + err.Error())
				return
			}
			msgChan := make(chan g90updatefw.Message)
			f.log(`Please do as follows:
> 1. Disconnect power cable from the radio.
> 2. Reconnect power cable to the radio.
> 3. Press the volume button and while holding it in,
> 4. Press the power button until the radio begins erasing the existing firmware.
`)
			go g90updatefw.UpdateRadio(serial, data, msgChan)

			var blocks = len(data) / 1024

			for {
				msg := <-msgChan
				switch msg.MsgType {
				case g90updatefw.MSG_STD:
					f.log(msg.Content)
				case g90updatefw.MSG_ERR:
					f.log("Error: " + msg.Content)
					dialog.NewError(err, myWindow).Show()
					progressBar.SetValue(0)
					return
				case g90updatefw.MSG_PGS:
					curBlock, _ := strconv.Atoi(msg.Content)
					val := float64(curBlock) / float64(blocks)
					if val > 1 {
						val = 1
					}
					progressBar.SetValue(val)
				case g90updatefw.MSG_FIN:
					f.log("Done.")
					progressBar.SetValue(0)
					return
				}
			}
		}()
	})

	topContainer := container.NewVBox(
		container.NewGridWithRows(2, firmwareContainer, portContainer),
		startButton,
	)

	content := container.NewBorder(
		container.NewVBox(topContainer, progressBar),
		nil,
		nil,
		nil,
		f.logEntry,
	)

	myWindow.SetContent(content)
	myWindow.Show()
}

func newFwgui() *fwgui {
	return &fwgui{}
}
