// Copyright 2020-2022 Dale Farnsworth. All rights reserved.

// Dale Farnsworth
// 1007 W Mendoza Ave
// Mesa, AZ  85210
// USA
//
// dale@farnsworth.org

// This program is free software: you can redistribute it and/or modify
// it under the terms of version 3 of the GNU General Public License
// as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// MODIFIED!
package g90updatefw

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/dalefarnsworth/go-xmodem/xmodem"
)

const (
	buflen = 64 * 1024

	MSG_STD = 0
	MSG_ERR = 1
	MSG_PGS = 2
	MSG_FIN = 3
)

type Message struct {
	MsgType uint8
	Content string
}

func makeStdMessage(msg string) Message {
	return Message{
		MsgType: MSG_STD,
		Content: msg,
	}
}

func makeErrMessage(msg string) Message {
	return Message{
		MsgType: MSG_ERR,
		Content: msg,
	}
}

func makeProgressMessage(msg string) Message {
	return Message{
		MsgType: MSG_PGS,
		Content: msg,
	}
}

func makeFinishMessage(msg string) Message {
	return Message{
		MsgType: MSG_FIN,
		Content: msg,
	}
}

func readString(serial *Serial, msgChan chan<- Message) (string, error) {
	buf := make([]byte, buflen)

	i := 0
	lastReadZeroBytes := false
	for {
		n, err := serial.Read(buf[i:])
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			if i == 0 {
				continue
			}
			if lastReadZeroBytes {
				// only return
				break
			}
			lastReadZeroBytes = true
			continue
		}
		lastReadZeroBytes = false
		msgChan <- makeStdMessage(string(buf[i : i+n]))
		i += n
	}

	if i >= len(buf) {
		return "", fmt.Errorf("read buffer overrun")
	}

	return string(buf[0:i]), nil
}

func expect(serial *Serial, expects []string, msgChan chan<- Message) (expectIndex int, err error) {
	previousStr := ""

	for {
		str, err := readString(serial, msgChan)
		if err != nil {
			return 0, err
		}
		for i, expect := range expects {
			if strings.Contains(previousStr+str, expect) {
				return i, nil
			}
		}
		previousStr = str
	}
}

func expectSend(serial *Serial, expects []string, sends []string, msgChan chan<- Message) (whichExpect string, err error) {
	if len(sends) != len(expects) {
		return "", fmt.Errorf("length of sends array does not equal length of expects array")
	}

	msgChan <- makeStdMessage(fmt.Sprintf("> Waiting for '%s'...\n\n", strings.Join(expects, "' or '")))

	expectIndex, err := expect(serial, expects, msgChan)
	if err != nil {
		return "", err
	}
	send := sends[expectIndex]

	if len(send) != 0 {
		_, err := serial.Write([]byte(send))
		if err != nil {
			return "", err
		}
	}

	return expects[expectIndex], nil
}
func UpdateRadio(serial *Serial, data []byte, msgChan chan<- Message) {
	attentionTimeout := 10 * time.Millisecond
	menuTimeout := 50 * time.Millisecond
	eraseTimeout := 50 * time.Millisecond
	uploadTimeout := 10 * time.Second
	cleanupTimeout := 500 * time.Millisecond

	banner := "Hit a key to abort"
	menu := "1.Update FW"
	waitFW := "Wait FW file"

	attentionGrabber := " "
	menuSelector := "1"

	serial.Flush()

	expects := []string{banner, menu}
	sends := []string{attentionGrabber, menuSelector}

	serial.SetReadTimeout(attentionTimeout)
	found, err := expectSend(serial, expects, sends, msgChan)
	if err != nil {
		msgChan <- makeErrMessage(err.Error())
		return
	}

	if found != menu {
		serial.SetReadTimeout(menuTimeout)
		_, err := expectSend(serial, []string{menu}, []string{menuSelector}, msgChan)
		if err != nil {
			msgChan <- makeErrMessage(err.Error())
			return
		}
	}

	serial.SetReadTimeout(eraseTimeout)
	_, err = expectSend(serial, []string{waitFW}, []string{""}, msgChan)
	if err != nil {
		msgChan <- makeErrMessage(err.Error())
		return
	}
	msgChan <- makeStdMessage(fmt.Sprintf("\n\n> Uploading %d bytes.\n", len(data)))

	serial.SetReadTimeout(uploadTimeout)
	//counter := 0
	//previousBlock := -1
	callback := func(block int) {
		//if counter%40 == 0 {
		//	if counter != 0 {
		//		msgChan <- makeStdMessage("\n")
		//	}
		//	msgChan <- makeStdMessage("> ")
		//}
		//marker := "."
		//if block != previousBlock+1 {
		//	marker = "R"
		//}
		//msgChan <- makeStdMessage(marker)
		//counter++
		//previousBlock = block
		msgChan <- makeProgressMessage(strconv.Itoa(block))
	}
	err = xmodem.ModemSend1K(serial, data, callback)
	if err != nil {
		msgChan <- makeErrMessage(err.Error())
		return
	}

	serial.SetReadTimeout(cleanupTimeout)
	readString(serial, msgChan)

	msgChan <- makeFinishMessage("\n> Upload complete.")
}
