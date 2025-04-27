package memory

import (
	"os"
	"strconv"
	"strings"

	momoError "momo/pkg/error"
)

type MemoMetric struct {
	MemTotal     uint64
	MemFree      uint64
	MemAvailable uint64
}

func New() (*MemoMetric, error) {
	dataRaw, err := getData()
	if err != nil {
		return nil, err
	}

	memTotal, memFree, memAvailable, err := extractData(dataRaw)

	return &MemoMetric{
		MemTotal:     memTotal,
		MemFree:      memFree,
		MemAvailable: memAvailable,
	}, nil
}

func extractData(dataRaw string) (uint64, uint64, uint64, error) {
	var memTotal, memFree, memAvailable uint64
	for _, line := range strings.Split(dataRaw, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := fields[0]

		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		switch key {
		case "MemTotal:":
			memTotal = value
		case "MemFree:":
			memFree = value
		case "MemAvailable:":
			memAvailable = value
		}
	}

	return memTotal, memFree, memAvailable, nil
}

func getData() (string, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", momoError.DebuggingErrorf("error to open /proc/meminfo the problem was %v", err)
	}
	return string(data), nil
}
