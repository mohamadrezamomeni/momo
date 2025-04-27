package cpu

import (
	"os"
	"strconv"
	"strings"

	momoError "momo/pkg/error"
)

type CpuMetric struct {
	idle  uint64
	total uint64
}

func New() (*CpuMetric, error) {
	data, err := readFileProcStat()
	if err != nil {
		return nil, momoError.DebuggingErrorf("something went wrong to open /proc/stat the problem was %v", err)
	}

	line, err := getCpuInfoLine(data)

	items := strings.Fields(line)
	if len(items) < 5 {
		return nil, momoError.DebuggingErrorf("we got unexpected error in proct/stat content")
	}
	idle, err := getIdle(items)
	total, err := getTotal(items)
	if err != nil {
		return nil, err
	}

	return &CpuMetric{idle: idle, total: total}, nil
}

func readFileProcStat() ([]byte, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, momoError.DebuggingErrorf("something went wrong to open /proc/stat the problem was %v", err)
	}
	return data, nil
}

func getCpuInfoLine(data []byte) (string, error) {
	lines := strings.Split(string(data), "\n")

	totalCpuline := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			totalCpuline = line
		}
	}
	if len(totalCpuline) == 0 {
		return "", momoError.DebuggingError("the total cput info wasn't found")
	}
	return totalCpuline, nil
}

func getIdle(items []string) (uint64, error) {
	idleStr := items[4]
	v, err := strconv.ParseUint(idleStr, 10, 64)
	if err != nil {
		return 0, momoError.DebuggingError("the problem was getting unexpected error while converting string to uint")
	}
	return v, nil
}

func getTotal(items []string) (uint64, error) {
	var total uint64 = 0
	for _, item := range items[1:] {
		v, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			return 0, momoError.DebuggingError("the problem was getting unexpected error while converting string to uint")
		}

		total += v
	}
	return total, nil
}
