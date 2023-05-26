package FohhnNet

import "strings"

func boolToOk(ok bool) string {
	if ok {
		return "OK"
	}
	return "ERR"
}

func decodeString(inputstr string) string {

	resp := strings.Replace(inputstr, "\xFF\x00", "\xF0", -1)

	return strings.Replace(resp, "\xFF\x01", "\xFF", -1)
}
