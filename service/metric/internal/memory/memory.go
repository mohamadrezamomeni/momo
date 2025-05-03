package memory

import (
	"os"
	"strconv"
	"strings"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
)

type MemoMetric struct {
	memInfoPath string
}

func New() *MemoMetric {
	return &MemoMetric{
		memInfoPath: "/proc/meminfo",
	}
}

func (m *MemoMetric) GetData() (uint64, uint64, uint64, error) {
	dataRaw, err := m.getData()
	if err != nil {
		return 0, 0, 0, err
	}

	return m.extractData(dataRaw)
}

func (m *MemoMetric) extractData(dataRaw string) (uint64, uint64, uint64, error) {
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

func (m *MemoMetric) getData() (string, error) {
	scope := "memoryMetric.getData"
	data, err := os.ReadFile(m.memInfoPath)
	if err != nil {
		return "", momoError.Wrap(err).Scope(scope).DebuggingErrorf("the items must be 5")
	}
	return string(data), nil
}
