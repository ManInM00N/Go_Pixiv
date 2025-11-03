package test

import (
	"os"
	"testing"
)

func TestLogMk(t *testing.T) {

	os.MkdirAll("trace/1145-14-19.log", 0777)

}
