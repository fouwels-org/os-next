// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package journal

import (
	"fmt"
	"os"
	"time"
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
	if ret {
		fmt.Fprintf(os.Stdout, "\n%v ", time.Now().Format("15:04:05.000000"))
	}
	fmt.Fprintf(os.Stdout, "%v", log)
}
