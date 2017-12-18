package cmd

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"

	"github.com/Unknwon/com"
	"github.com/olekukonko/tablewriter"
)

type Writer struct {
	format         string
	output         io.Writer
	outputFileName string
	tbl            *tablewriter.Table
	csv            *csv.Writer
}

func NewWriter(header []string) *Writer {
	var output io.Writer = os.Stdout
	if outputFileName != "" {
		output = &bytes.Buffer{}
	}
	wr := &Writer{
		format:         outputFormat,
		output:         output,
		outputFileName: outputFileName,
	}
	switch wr.format {
	case "table":
		wr.tbl = tablewriter.NewWriter(output)
	case "csv":
		wr.csv = csv.NewWriter(output)
	}
	if header != nil && !noHeader {
		wr.Header(header)
	}
	return wr
}

func (w *Writer) Header(header []string) {
	switch w.format {
	case "table":
		w.tbl.SetHeader(header)
	case "csv":
		w.csv.Write(header)
	}
}

func (w *Writer) Row(row []string) {
	switch w.format {
	case "table":
		w.tbl.Append(row)
	case "csv":
		w.csv.Write(row)
	}
}

func (w *Writer) Flush() {
	switch w.format {
	case "table":
		w.tbl.Render()
	case "csv":
		w.csv.Flush()
	}
}

func (w *Writer) Close() {
	w.Flush()
	if w.outputFileName != "" {
		com.WriteFile(w.outputFileName, w.output.(*bytes.Buffer).Bytes())
	}
}
