package cpu

import (
	"os"
	"strconv"
	"strings"

	momoError "momo/pkg/error"
)

type CpuMetric struct {
	statFilePath string
}

func New() *CpuMetric {
	return &CpuMetric{
		statFilePath: "/proc/stat",
	}
}

func (c *CpuMetric) getData() (uint64, uint64, error) {
	data, err := c.readFileProcStat()
	if err != nil {
		return 0, 0, momoError.DebuggingErrorf("something went wrong to open /proc/stat the problem was %v", err)
	}

	line, err := c.getCpuInfoLine(data)

	items := strings.Fields(line)
	if len(items) < 5 {
		return 0, 0, momoError.DebuggingErrorf("we got unexpected error in proct/stat content")
	}
	idle, err := c.getIdle(items)
	total, err := c.getTotal(items)
	if err != nil {
		return 0, 0, err
	}

	return total, idle, nil
}

func (c *CpuMetric) readFileProcStat() ([]byte, error) {
	data, err := os.ReadFile(c.statFilePath)
	if err != nil {
		return nil, momoError.DebuggingErrorf("something went wrong to open /proc/stat the problem was %v", err)
	}
	return data, nil
}

func (c *CpuMetric) getCpuInfoLine(data []byte) (string, error) {
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

func (c *CpuMetric) getIdle(items []string) (uint64, error) {
	idleStr := items[4]
	v, err := strconv.ParseUint(idleStr, 10, 64)
	if err != nil {
		return 0, momoError.DebuggingError("the problem was getting unexpected error while converting string to uint")
	}
	return v, nil
}

func (c *CpuMetric) getTotal(items []string) (uint64, error) {
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
