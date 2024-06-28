package utils

import (
	"log"
	"strconv"
)

func ParseInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func ParseFloat64(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func ParseIntToString(value int) string {
	return strconv.Itoa(value)
}

func ParseIntPointer(s string) *int {
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func ParseFloat64Pointer(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}

func ParseStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ParseIntOrZero(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ParseFloatOrZero(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
