package rlog

import "testing"

func TestDebug(t *testing.T) {
	Debug("test")
	Debug("%s", "this is a test")
}
