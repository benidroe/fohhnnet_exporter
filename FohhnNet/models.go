package FohhnNet

/*
Convert Modelnumber to Model
0x0C40 = MA-4.100 ANA
0x0C50 = MA-4.100 DAN
0x0C60 = MA-4.600 ANA
0x0C70 = MA-4.600 DAN
0x0C80 = MA-2.1200 ANA
0x0C90 = MA-2.1200 DAN
0x0330 = D-2.1500
0x0331 = D-2.1500
0x0332 = D-2.1500
0x0333 = D-2.1500
0x0300 = D-2.750
0x0301 = D-2.750
0x0302 = D-2.750
0x0303 = D-2.750
0x0304 = D-2.750
0x0310 = D-4.750
0x0311 = D-4.750
0x0312 = D-4.750
0x0321 = D-4.1200
0x0322 = D-4.1200
0x0D20 = DLI-130
0x0D30 = DLI-230
0x0D40 = DLI-330
0x0D50 = DLI-430
0x0D21 = DLI-130
0x0D31 = DLI-230
0x0C00 = DI-2.2000
0x0C10 = DI-4.1000
0x0C20 = DI-4.2000
0x0C21 = DI-4.2000
0x0C30 = DI-2.4000
0x0C31 = DI-2.4000
0x0C01 = DI-2.2000
0x0C11 = DI-4.1000

*/

type model struct {
	name     string
	channels int8
}

var models = map[string]model{
	"\x0c\x40": model{name: "MA-4.100 ANA", channels: 4},
	"\x0c\x50": model{name: "MA-4.100 DAN", channels: 4},
	"\x0c\x60": model{name: "MA-4.600 ANA", channels: 4},
	"\x0c\x70": model{name: "MA-4.600 DAN", channels: 4},
	"\x0c\x80": model{name: "MA-2.1200 ANA", channels: 2},
	"\x0c\x90": model{name: "MA-2.1200 DAN", channels: 2},
	"\x03\x30": model{name: "D-2.1500", channels: 2},
	"\x03\x00": model{name: "D-2.750", channels: 2},
	"\x03\x10": model{name: "D-4.750", channels: 4},
	"\x03\x20": model{name: "D-4.1200", channels: 4},
	"\x0D\x20": model{name: "DLI-130", channels: 1},
	"\x0D\x30": model{name: "DLI-230", channels: 1},
	"\x0D\x40": model{name: "DLI-330", channels: 1},
	"\x0D\x50": model{name: "DLI-430", channels: 1},
	"\x0C\x00": model{name: "DI-2.2000", channels: 2},
	"\x0C\x10": model{name: "DI-4.1000", channels: 4},
	"\x0C\x20": model{name: "DI-4.2000", channels: 4},
	"\x0C\x30": model{name: "DI-2.4000", channels: 2},
}

func GetModelNameByNumber(hex string) string {

	if len(hex) > 1 {
		// Mask hardware version (last four bits 0xFFF0)
		i := hex[0]
		j := hex[1] & '\xF0'
		masked := string(i) + string(j)
		val, ok := models[masked]
		if ok {
			return val.name
		}
	}
	return "-"
}

func GetNumOfModelChannels(hex string) int8 {

	if len(hex) > 1 {
		// Mask hardware version (last four bits 0xFFF0)
		i := hex[0]
		j := hex[1] & '\xF0'
		masked := string(i) + string(j)
		val, ok := models[masked]
		if ok {
			return val.channels
		}
	}

	return 1
}
