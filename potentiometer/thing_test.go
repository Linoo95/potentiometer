package main

import (
	"testing"
)

func TestAdcToPercent_Max(t *testing.T) {
	if adcToPercent(4095) != 100 {
		t.Fail()
	}
}
func TestAdcToPercent_Min(t *testing.T) {
	if adcToPercent(0) != 0 {
		t.Fail()
	}
}

func TestAdcToPercent_Half(t *testing.T) {
	if adcToPercent(2048) != 50 {
		t.Fail()
	}
}

func TestCleanLine_NoNumber(t *testing.T) {
	if cleanLine("woopada") != "" {
		t.Fail()
	}
}

func TestCleanLine_NestedNumber(t *testing.T) {
	if cleanLine("kak2048ada\r\n") != "2048" {
		t.Fail()
	}
}
