package go_linux_apis

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strconv"
)

type NetDevReceive struct {
	Bytes, Packets, Errs, Drop, Fifo, Frame, Compressed, Multicast uint64
}

type NetDevTransmit struct {
	Bytes, Packets, Errs, Drop, Fifo, Colls, Carrier, Compressed uint64
}

type NetDev struct {
	Receive  NetDevReceive
	Transmit NetDevTransmit
}

var lf = [1]byte{'\n'}
var colon = [1]byte{':'}

var counter = regexp.MustCompile(`\b\d+\b`)

func GetNetDev() (map[string]NetDev, error) {
	content, errRF := ioutil.ReadFile("/proc/net/dev")
	if errRF != nil {
		return nil, errRF
	}

	lines := bytes.Split(content, lf[:])
	content = nil

	if len(lines) < 3 {
		// The Linux guys have broken the userspace...
		return map[string]NetDev{}, nil
	}

	lines = lines[2:]

	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	res := make(map[string]NetDev, len(lines))

	for _, line := range lines {
		if ifAndStats := bytes.SplitN(line, colon[:], 2); len(ifAndStats) > 1 {
			c := [16]uint64{}

			for i, n := range counter.FindAll(ifAndStats[1], 16) {
				m, errUint := strconv.ParseUint(string(n), 10, 64)
				if errUint != nil {
					m = ^uint64(0)
				}

				c[i] = m
			}

			res[string(bytes.TrimSpace(ifAndStats[0]))] = NetDev{
				NetDevReceive{c[0], c[1], c[2], c[3], c[4], c[5], c[6], c[7]},
				NetDevTransmit{c[8], c[9], c[10], c[11], c[12], c[13], c[14], c[15]},
			}
		}
	}

	return res, nil
}
