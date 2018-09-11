package go_linux_apis

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"time"
)

type Uptime struct {
	UpTime, IdleTime time.Duration
}

const maxDuration = time.Duration(^uint64(0) >> 1)
const minDuration = time.Duration(^maxDuration)

var procUptimeRgx = regexp.MustCompile(`\A(\d+(?:\.\d+)?) (\d+(?:\.\d+)?)\b`)

func GetUptime() (Uptime, error) {
	content, errRF := ioutil.ReadFile("/proc/uptime")
	if errRF != nil {
		return Uptime{}, errRF
	}

	match := procUptimeRgx.FindSubmatch(content)
	if match == nil {
		// The Linux guys have broken the userspace...
		return Uptime{minDuration, minDuration}, nil
	}

	return Uptime{seconds2duration(string(match[1])), seconds2duration(string(match[2]))}, nil
}

func seconds2duration(secs string) time.Duration {
	if seconds, errPF := strconv.ParseFloat(secs, 64); errPF == nil {
		return time.Duration(seconds*1000000000.0) * time.Nanosecond
	}

	// Someone hasn't restarted their Linux for a very long time...
	return maxDuration
}
