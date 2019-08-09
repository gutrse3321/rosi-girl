/**
 * @Author: Tomonori
 * @Date: 2019/6/3 10:57
 * @File: logger
 * @Desc:
 */
package util

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/pio"
	"os"
	"time"
)

var pioCtx *pio.Printer

func init() {
	golog.SetTimeFormat("")
	golog.SetLevel("debug")
	pioCtx = pio.NewTextPrinter("color", os.Stdout)
}

type ILogger interface {
	Normal(opt ...interface{})
	Info(opt ...interface{})
	Underline(opt ...interface{})
	Complate(str string)
}

type Logger struct {
}

func (l *Logger) Normal(opt ...interface{}) {
	l.printDatetime()
	golog.Println(opt...)
}

func (l *Logger) Info(opt ...interface{}) {
	l.printDatetime()
	golog.Info(opt...)
}

func (l *Logger) Underline(opt ...interface{}) {
	l.printDatetime()
	golog.Warn(opt...)
}

func (l *Logger) Complate(str string) {
	l.printDatetime()
	str = fmt.Sprintf("[COMPLATE] %s", str)
	pioCtx.Println(pio.Green(str))
}

func (l *Logger) printDatetime() {
	time := fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))
	pioCtx.Print(pio.Yellow(time))
}
