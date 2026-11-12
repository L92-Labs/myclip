package native

import (
	"bytes"
	"fmt"
)

// CF_HTML ("HTML Format") wraps an HTML fragment in a text header with
// byte offsets into the payload. These two are pure and tested everywhere.

// cfhtmlParse returns the payload's HTML fragment, or the payload itself
// when it has no readable header.
func cfhtmlParse(b []byte) []byte {
	off := func(name string) (int, bool) {
		i := bytes.Index(b, []byte(name+":"))
		if i < 0 {
			return 0, false
		}
		i += len(name) + 1
		n := 0
		for i < len(b) && b[i] >= '0' && b[i] <= '9' {
			n = n*10 + int(b[i]-'0')
			i++
		}
		if n <= 0 || n >= len(b) {
			return 0, false
		}
		return n, true
	}
	start, okStart := off("StartFragment")
	end, okEnd := off("EndFragment")
	if !okStart || !okEnd || start >= end {
		return b
	}
	return b[start:end]
}

// cfhtmlBuild wraps fragment in a CF_HTML payload. The header fields are
// fixed width, so the offsets are known before the payload is assembled:
// one line block of four 10-digit fields, then the HTML around the
// fragment, with an ignored byte marking where the fragment starts.
func cfhtmlBuild(fragment []byte) []byte {
	head := fmt.Sprintf("Version:0.9\r\nStartHTML:%010d\r\nEndHTML:%010d\r\nStartFragment:%010d\r\nEndFragment:%010d\r\n",
		0, 0, 0, 0)
	const (
		mark = "<html><body><!--StartFragment-->"
		tail = "<!--EndFragment--></body></html>"
	)
	startHTML := len(head)
	startFrag := startHTML + len(mark)
	endFrag := startFrag + len(fragment)
	endHTML := endFrag + len(tail)

	payload := make([]byte, 0, endHTML)
	payload = append(payload, fmt.Sprintf("Version:0.9\r\nStartHTML:%010d\r\nEndHTML:%010d\r\nStartFragment:%010d\r\nEndFragment:%010d\r\n",
		startHTML, endHTML, startFrag, endFrag)...)
	payload = append(payload, mark...)
	payload = append(payload, fragment...)
	payload = append(payload, tail...)
	return payload
}
