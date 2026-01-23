package chainutils

import (
	"fmt"
	"math"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/mathz"
	"strings"
)

func ShortenAddress(address string, count ...int) string {

	c := 4
	if len(count) > 0 {
		c = count[0]
	}

	if len(address) <= c*2 {
		return address
	}

	return address[:c] + "..." + address[len(address)-4:]
}

func FormatValue(raw uint64, decimals int) float64 {
	return mathz.ToFixed8(float64(raw) / math.Pow(10.0, float64(decimals)))
}

type FormatDecimalOptions struct {
	ZeroCount           int
	MaxSignificantCount int
}

func FormatDecimal(raw float64, options ...FormatDecimalOptions) string {

	if raw > 1000000 {
		//return `${(value / 1000000).toFixed(2)}M `
		return fmt.Sprintf("%.2fM", raw/1000000)
	}

	if raw > 1000 {
		return fmt.Sprintf("%.2fK", raw/1000)
	}

	zeroCount := 4
	maxSignificantCount := 5
	if len(options) > 0 {
		if options[0].MaxSignificantCount > 0 {
			maxSignificantCount = options[0].MaxSignificantCount
		}

		if options[0].ZeroCount > 0 {
			zeroCount = options[0].ZeroCount
		}
	}

	if raw*math.Pow(10, float64(zeroCount)) >= 1 {
		return formatWithSignificantDecimals(raw, maxSignificantCount)
	}

	return formatDecimal(formatWithSignificantDecimals(raw, maxSignificantCount))

}

func formatWithSignificantDecimals(value float64, maxSignificantCount ...int) string {

	maxSignificant := 5
	if len(maxSignificantCount) > 0 {
		maxSignificant = maxSignificantCount[0]
	}

	// 将浮点数转为字符串
	strValue := fmt.Sprintf("%.20f", value)

	// 去掉尾部的零
	strValue = strings.TrimRight(strValue, "0")
	// 将 "0." 替换为 "0" 以处理非零的情况
	if strings.HasPrefix(strValue, "0.") {
		strValue = "0" + strValue[1:]
	}

	// 找到小数点的位置
	dotIndex := strings.Index(strValue, ".")
	if dotIndex == -1 {
		return strValue // 如果没有小数点，直接返回
	}

	// 获取小数和整数部分
	integerPart := strValue[:dotIndex]
	decimalPart := strValue[dotIndex+1:]

	// 有效位计数
	significantCount := 0
	var resultDec string

	// 处理小数部分，保留五位有效位
	for _, ch := range decimalPart {
		if ch != '0' {
			significantCount++
		}
		if significantCount <= maxSignificant {
			resultDec += string(ch)
		}
		if significantCount >= maxSignificant {
			break
		}
	}

	// 如果有效位数字不到五位，保留原有的小数位
	if significantCount < 5 {
		return fmt.Sprintf("%s.%s", integerPart, resultDec)
	}

	// 如果整数部分是零，格式化时要确保是以 "0" 开头
	if integerPart == "0" {
		return fmt.Sprintf("0.%s", resultDec)
	}

	return fmt.Sprintf("%s.%s", integerPart, resultDec)
}

// eg.  0.000006600 转换成 0.{5}66
func formatDecimal(input string) string {
	// 去除前导的零
	if strings.HasPrefix(input, "0.") {

		input = input[2:]

		// 找到小数点后第一个非零数的位置 即0前导0的
		firstNonZeroIndex := strings.IndexFunc(input, func(r rune) bool {
			return r != '0'
		})

		if firstNonZeroIndex == -1 {
			return "0"
		}

		nonZeroPart := input[firstNonZeroIndex:]

		return fmt.Sprintf("0.{%d}%s", firstNonZeroIndex, nonZeroPart)
	}

	return input
}

// todo
func FormatPercentage(rate float64) string {
	return conv.String(conv.Int(conv.Float64(rate)*100)) + "%"
}
