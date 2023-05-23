package FohhnNet

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"
)

// Alle Commands ohne Startbyte und ohne Device ID und ohne Databyte Count \xF0 \xID \x01

/*
1. Byte Startbyte <SB>
2. Byte Device ID <ID>
3. Byte Databyte Count <COUNT>
4. Byte Command Byte <CMD>
5. Byte Address MSB <ADR_MSB>
6. Byte Address LSB <ADR_LSB>
7. Byte Databyte 1 <DATA> // min. one databyte
N. Bytes

*/

type fohhnNetSession struct {
	fohhnDialer    net.Dialer
	IsConnected    bool
	connection     net.Conn
	reader         *bufio.Reader
	responseChanel chan string
	responsebuf    []byte
	hasFailed      bool
}

type fohhnDeviceStateSet struct {
	Id                   int8
	Version              string
	Device               string
	Alias                string
	OperatingTimeHours   uint32
	OperatingTimeMinutes uint8
	Protect              []bool
	Standby              bool
	Temperature          float32
	SpeakerPreset        []string
	OutputChannelName    []string
}

func NewFohhnNetTcpSession(host string, port int) (*fohhnNetSession, error) {

	f := fohhnNetSession{}
	f.fohhnDialer = net.Dialer{Timeout: 1 * time.Second}
	conn, err := f.fohhnDialer.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error while connecting")
		return nil, errors.New("Konnte keine Verbindung zu FohhnNet Node aufbauen")
	} else {
		f.IsConnected = true
		f.hasFailed = false
		f.connection = conn
		f.reader = bufio.NewReader(f.connection)
		//fmt.Println("Connected")

		f.responseChanel = make(chan string, 1)
		go responseGoRoutine(&f)

	}
	return &f, nil
}

func ScrapeFohhnDevice(fohhnNetSession *fohhnNetSession, id int8) (*fohhnDeviceStateSet, error) {

	hasAnswered := false
	_ = hasAnswered

	fds := fohhnDeviceStateSet{}

	fds.Id = id

	if fohhnNetSession.IsConnected {

		protect1, protect2, protect3, protect4, temperature, err := GetControls(fohhnNetSession, id)

		if err == nil {
			hasAnswered = true
			fds.Protect = append(fds.Protect, protect1)
			fds.Protect = append(fds.Protect, protect2)
			fds.Protect = append(fds.Protect, protect3)
			fds.Protect = append(fds.Protect, protect4)
			fds.Temperature = temperature
		} else {
			hasAnswered = false
			return nil, errors.New("Gerät hat nicht geantwortet")
		}

		deviceInfo, version, err := GetDeviceInfo(fohhnNetSession, id)
		if err == nil {
			fds.Version = version
			fds.Device = deviceInfo
		} else {
			return nil, err
		}

		operatingHours, operatingMinutes, err := GetOperatingTime(fohhnNetSession, id)
		if err == nil {
			fds.OperatingTimeHours = operatingHours
			fds.OperatingTimeMinutes = operatingMinutes
		} else {
			return nil, err
		}

		deviceAlias, err := GetDeviceAlias(fohhnNetSession, id)
		if err == nil {
			fds.Alias = deviceAlias
		} else {
			return nil, err
		}

		deviceStandby, err := GetStandby(fohhnNetSession, id)
		if err == nil {
			fds.Standby = deviceStandby
		} else {
			return nil, err
		}

		numOfChannels := GetNumOfModelChannels(fds.Device)

		for i := int8(1); i <= numOfChannels; i++ {
			speakerPresetName, err := GetCurrentSpeakerPresetName(fohhnNetSession, id, i)
			if err == nil {
				fds.SpeakerPreset = append(fds.SpeakerPreset, speakerPresetName)
			} else {
				return nil, err
			}
		}

		for i := int8(1); i <= numOfChannels; i++ {
			outputChannelName, err := GetOutputChannelName(fohhnNetSession, id, i)
			if err == nil {
				fds.OutputChannelName = append(fds.OutputChannelName, outputChannelName)
			} else {
				return nil, err
			}
		}

	}

	return &fds, nil

}

func Close(fohhnNetSession *fohhnNetSession) {

	fohhnNetSession.connection.Close()
}

func ScanFohhnNet(fohhnNetSession *fohhnNetSession, idMin int8, idMax int8) []int8 {

	foundDevs := []int8{}
	for i := idMin; i <= idMax; i++ {

		_, _, err := GetDeviceInfo(fohhnNetSession, i)

		if err == nil {
			//fmt.Printf("Version: %s, Device: %d\n", version, deviceInfo)
			//fmt.Println("Device is responding. ID=", i)
			foundDevs = append(foundDevs, i)
		} else {
			//fmt.Println("Device is not responding. ID=", i)
		}

		//time.Sleep(1 * time.Second)
	}

	return foundDevs
}

func GetControls(fohhnNetSession *fohhnNetSession, deviceId int8) (bool, bool, bool, bool, float32, error) {

	// because of byte loss, we need sometimes a second try.
	for i := 0; i < 2; i++ {

		sendMessage(fohhnNetSession, deviceId, "\x07", "\x00", "\x00", "\x00")
		msg, err := receiveStateDelimited(fohhnNetSession)

		if err != nil {
			fmt.Println("Error", err)
		} else {

			// retry, if message length does not match 6 bytes.
			if len(msg) != 6 {
				time.Sleep(350 * time.Millisecond)
				continue
			}

			if len(msg) > 2 {

				// Temperature als unsigned word einlesen
				numBytes := []byte{msg[1], msg[2]}
				u := binary.BigEndian.Uint16(numBytes)

				// Und die eingelesene Bitfolge als signed word behandeln und durch 10 teilen
				temperature := float32(int16(u))

				// Bit 0,1,2,3 steht jeweils für den Protect des Kanals. Diesen mit einem Bitshift und AND ermitteln
				protect1 := msg[0]&(1<<0) == 0
				protect2 := msg[0]&(1<<1) == 0
				protect3 := msg[0]&(1<<2) == 0
				protect4 := msg[0]&(1<<3) == 0

				return protect1, protect2, protect3, protect4, temperature, nil
			}
		}
	}

	return false, false, false, false, 0, errors.New("Keine gültige Antwort empfangen")

}

func GetDeviceInfo(fohhnNetSession *fohhnNetSession, deviceId int8) (string, string, error) {

	sendMessage(fohhnNetSession, deviceId, "\x20", "\x00", "\x00", "\x01")
	msg, err := receiveStateDelimited(fohhnNetSession)

	if err != nil {
		//fmt.Println("Error", err)
	} else {

		if len(msg) > 4 {

			device := msg[0:2]

			version := fmt.Sprintf("%d.%d.%d", msg[2], msg[3], msg[4])

			return device, version, nil
		}

	}

	return "", "", errors.New("Keine gültige Antwort empfangen")

}

func GetOperatingTime(fohhnNetSession *fohhnNetSession, deviceId int8) (uint32, uint8, error) {

	sendMessage(fohhnNetSession, deviceId, "\x0B", "\x01", "\x00", "\x00")
	msg, err := receiveStateDelimited(fohhnNetSession)

	if err != nil {
		fmt.Println("Error", err)
	} else {

		if len(msg) > 3 {

			// Betriebsstunden als unsigned einlesen
			numBytes := []byte{0, msg[0], msg[1], msg[2]}
			operatingTime := binary.BigEndian.Uint32(numBytes)
			operatingMinutes := uint8(msg[3])

			return operatingTime, operatingMinutes, nil
		}
	}

	return 0, 0, errors.New("Keine gültige Antwort empfangen")

}

func GetDeviceAlias(fohhnNetSession *fohhnNetSession, deviceId int8) (string, error) {

	sendMessage(fohhnNetSession, deviceId, "\x90", "\x01", "\x00", "\x00")
	msg, err := receiveStateDelimited(fohhnNetSession)

	if err != nil {
		fmt.Println("Error", err)
	} else {

		if len(msg) > 17 {
			deviceAlias := msg[2:17]
			return deviceAlias, nil
		}

	}

	return "nil", errors.New("Keine gültige Antwort empfangen")

}

func GetStandby(fohhnNetSession *fohhnNetSession, deviceId int8) (bool, error) {

	sendMessage(fohhnNetSession, deviceId, "\x0A", "\x00", "\x00", "\x0C")
	msg, err := receiveStateDelimited(fohhnNetSession)

	if err != nil {
		fmt.Println(false, err)
	} else {

		if len(msg) > 0 {

			standby := msg[0]
			if standby == 1 {
				return true, nil
			} else {
				return false, nil
			}
		}

	}

	return false, errors.New("Keine gültige Antwort empfangen")

}

func GetCurrentSpeakerPresetName(fohhnNetSession *fohhnNetSession, deviceId int8, channel int8) (string, error) {

	sendMessage(fohhnNetSession, deviceId, "\x22", string(channel), "\x00", "\x02")
	msg, err := receiveStateDelimited(fohhnNetSession)

	if err != nil {
		fmt.Println("Error", err)
	} else {

		if len(msg) >= 37 {
			speakerPreset := msg[22:37]
			return speakerPreset, nil
		}
	}

	return "nil", errors.New("GetCurrentSpeakerPresetName - Keine gültige Antwort empfangen")

}

func GetOutputChannelName(fohhnNetSession *fohhnNetSession, deviceId int8, channel int8) (string, error) {

	sendMessage(fohhnNetSession, deviceId, "\x94", string(channel), "\x01", "\x00")
	msg, err := receiveStateDelimited(fohhnNetSession)

	if err != nil {
		fmt.Println("Error", err)
	} else {

		if len(msg) >= 18 {
			channelName := msg[2:18]
			return channelName, nil
		}

	}

	return "nil", errors.New("GetOutputChannelName - Keine gültige Antwort empfangen")

}

func sendMessage(fohhnNetSession *fohhnNetSession, deviceId int8, commandByte string, adrMsb string, adrLsb string, data string) {
	fohhnNetSession.responsebuf = nil
	databyteCount := len(data)
	message := "\xF0" + string(rune(deviceId)) + string(rune(databyteCount)) + commandByte + adrMsb + adrLsb + data // send request to device
	fmt.Fprintf(fohhnNetSession.connection, message)                                                                // send request to device
}

/**
Diese Go Routine empfängt Strings von der FohhnNetSession und schreibt diese auf den responseChanel.
Gestartet wird sie im Constructor gleich nach Aufbau der Verbindung.
*/
func responseGoRoutine(fohhnNetSession *fohhnNetSession) {
	for {

		maxlen := 42
		fohhnNetSession.responsebuf = nil
		for {

			message, err := fohhnNetSession.reader.ReadByte() // read response

			if err != nil {
				// return for example when the connection is closed. Otherwise goroutine will run forever.
				return
			}
			fohhnNetSession.responsebuf = append(fohhnNetSession.responsebuf, message)
			maxlen -= 1
			if message == '\xF0' {
				fohhnNetSession.responseChanel <- string(fohhnNetSession.responsebuf)
				break
			}
			if maxlen < 1 {
				break
			}
		}

	}
}

/**
Holt letzte Antwort vom responseChanel. Nach erreichen des Timeouts wird ein error zurückgeliefert.
*/

func receiveStateDelimited(fohhnNetSession *fohhnNetSession) (string, error) {

	// Select Timeout after some milliseconds...
	select {
	case msg := <-fohhnNetSession.responseChanel:

		pattern, _ := regexp.Compile(`.*f0`)             // Das Pattern muss mit dem SB \xF0 enden
		if pattern.MatchString(fmt.Sprintf("%x", msg)) { // Weil MatchString mit 4-Byte Unicode arbeitet, muss Hex als Hex-String Formatiert werden und kann erst dann gematched werden.
			return msg, nil
		}
		fohhnNetSession.hasFailed = true
		return "", errors.New("Error: Response does not match SB")

	case <-time.After(150 * time.Millisecond):
		// fmt.Println("Error - No delimited reponse")
		fohhnNetSession.hasFailed = true
		return "", errors.New("Error: No delimited response from device - timeout")
	}

}
