package journal

import (
	"fmt"
	"os"
)

func Logfln(format string, v ...interface{}) {
	mux(fmt.Sprintf(format, v...), true)
}

func Logf(format string, v ...interface{}) {
	mux(fmt.Sprintf(format, v...), false)
}

func mux(log string, ret bool) {
	stdout(log, ret)
}

func stdout(log string, ret bool) {
	fmt.Fprintf(os.Stdout, "%v", log)
	if ret {
		fmt.Printf("\n")
	}
}
