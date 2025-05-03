package cpu

import (
	"os"
	"strconv"
	"strings"

	momoError "github.com/mohamadrezamomeni/momo/pkg/error"
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
	scope := "cpuMetric.getData"
	data, err := c.readFileProcStat()
	if err != nil {
		return 0, 0, momoError.Wrap(err).Scope(scope).DebuggingError()
	}

	line, err := c.getCpuInfoLine(data)

	items := strings.Fields(line)
	if len(items) < 5 {
		return 0, 0, momoError.Scope(scope).DebuggingErrorf("the items must be 5")
	}
	idle, err := c.getIdle(items)
	total, err := c.getTotal(items)
	if err != nil {
		return 0, 0, err
	}

	return total, idle, nil
}

func (c *CpuMetric) readFileProcStat() ([]byte, error) {
	scope := "cpuMetric.readFileProcStat"

	data, err := os.ReadFile(c.statFilePath)
	if err != nil {
		return nil, momoError.Wrap(err).Scope(scope).DebuggingError()
	}
	return data, nil
}

func (c *CpuMetric) getCpuInfoLine(data []byte) (string, error) {
	scope := "cpuMetric.getCpuInfoLine"

	lines := strings.Split(string(data), "\n")

	totalCpuline := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			totalCpuline = line
		}
	}
	if len(totalCpuline) == 0 {
		return "", momoError.Scope(scope).DebuggingErrorf("the total cput info wasn't found")
	}
	return totalCpuline, nil
}

func (c *CpuMetric) getIdle(items []string) (uint64, error) {
	scope := "cpuMetric.getIdle"

	idleStr := items[4]
	v, err := strconv.ParseUint(idleStr, 10, 64)
	if err != nil {
		return 0, momoError.Wrap(err).Scope(scope).DebuggingErrorf("the problem was getting unexpected error while converting string to uint")
	}
	return v, nil
}

func (c *CpuMetric) getTotal(items []string) (uint64, error) {
	scope := "cpuMetric.getTotal"

	var total uint64 = 0
	for _, item := range items[1:] {
		v, err := strconv.ParseUint(item, 10, 64)
		if err != nil {
			return 0, momoError.Wrap(err).Scope(scope).DebuggingErrorf("the problem was getting unexpected error while converting string to uint")
		}

		total += v
	}
	return total, nil
}
