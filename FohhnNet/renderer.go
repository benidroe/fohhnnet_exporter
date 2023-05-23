package FohhnNet

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func RenderFohhnNetScan(ids []int8) {

	fmt.Println("==== FohhnNet scan result ===")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, id := range ids {
		fmt.Fprintf(w, "ID  \t| STATE\n")
		fmt.Fprintf(w, "=======\t|=======\n")
		fmt.Fprintf(w, "%d \t| ONLINE\n", id)

	}
	w.Flush()
	fmt.Println("===========================================================")
}

func RenderFohhnDevice(fds *fohhnDeviceStateSet) {

	stby := "STANDBY OFF"
	if fds.Standby {
		stby = "STANDBY ON"
	}

	fmt.Printf("Device ID:      \t %d\n", fds.Id)
	fmt.Printf("Device Alias:    \t %s\n", fds.Alias)
	fmt.Printf("Device Model:       \t %s\n", GetModelNameByNumber(fds.Device))
	fmt.Printf("Device Version:  \t %s\n", fds.Version)
	fmt.Printf("Device Standby:  \t %s\n", stby)
	fmt.Printf("Temperature:    \t %.1f °C\n", fds.Temperature)
	fmt.Printf("Operation Time:    \t %d:%dh\n", fds.OperatingTimeHours, fds.OperatingTimeMinutes)
	fmt.Println("===========================================================")
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	tableHeading := "Channel"
	tablePreset := "Speaker Preset"
	tableName := "Channel Name"
	tableProtect := "Protect"
	for i := int8(0); i < GetNumOfModelChannels(fds.Device); i++ {
		tableHeading += fmt.Sprintf(" \t| CH %d", i+1)
		tablePreset += fmt.Sprintf(" \t| %s", fds.SpeakerPreset[i])
		tableName += fmt.Sprintf(" \t| %s", fds.OutputChannelName[i])
		tableProtect += fmt.Sprintf(" \t| %s", boolToOk(fds.Protect[i]))
	}

	fmt.Fprintf(w, tableHeading+"\n")
	fmt.Fprintf(w, tablePreset+"\n")
	fmt.Fprintf(w, tableName+"\n")
	fmt.Fprintf(w, tableProtect+"\n")
	w.Flush()
	fmt.Println("===========================================================")

}
