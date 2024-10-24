package memory

import (
	"strconv"
	"strings"
)

type Suffix string

const (
	SuffixKB Suffix = "kb"
	SuffixMB Suffix = "mb"
	SuffixGB Suffix = "gb"
	SuffixTB Suffix = "tb"
)

var memoryUnitMultiplier = map[Suffix]uint{
	SuffixKB: 1 << 10, // 1024
	SuffixMB: 1 << 20, // 1024 * 1024
	SuffixGB: 1 << 30, // 1024 * 1024 * 1024
	SuffixTB: 1 << 40, // 1024 * 1024 * 1024 * 1024
}

func ParseSize(sizeStr string) (uint, error) {
	sizeStr = strings.ToLower(strings.TrimSpace(sizeStr))

	for unit, multiplier := range memoryUnitMultiplier {
		if strings.HasSuffix(sizeStr, string(unit)) {
			numberStr := strings.TrimSuffix(sizeStr, string(unit))
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
