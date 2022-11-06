package parser

import (
	"github.com/jotaen/klog/klog/parser/txt"
)

type errData struct {
	id   string
	line int
	pos  int
	len  int
}

func (e HumanError) toErrData(line int, pos int, len int) errData {
	return errData{e.code, line, pos, len}
}

func toErrData(e txt.Error) errData {
	return errData{e.Code(), e.Context().LineNumber, e.Position(), e.Length()}
}