package listener

import (
	"strconv"
	"strings"
)

var memoryUnitMultiplier = map[string]uint64{
	"b":  1,
	"kb": 1024,
	"mb": 1024 * 1024,
	"gb": 1024 * 1024 * 1024,
	"tb": 1024 * 1024 * 1024 * 1024,
}

func parseMemorySize(sizeStr string) (uint, error) {
	sizeStr = strings.TrimSpace(sizeStr)

	for unit, multiplier := range memoryUnitMultiplier {
		if strings.HasSuffix(sizeStr, unit) {
			numberStr := strings.TrimSuffix(sizeStr, unit)
			number, err := strconv.ParseFloat(numberStr, 64)
			if err != nil {
				return 0, err
			}

			return uint(number * float64(multiplier)), nil
		}
	}

	number, err := strconv.ParseUint(sizeStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(number), nil
}
