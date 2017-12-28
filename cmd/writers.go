package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Unknwon/com"
	"github.com/olekukonko/tablewriter"
)

type Writer struct {
	format         string
	output         io.Writer
	outputFileName string
	tbl            *tablewriter.Table
	csv            *csv.Writer
	json           []string
}

type Rower interface {
	Header() []string
	Row() []string
}

func NewWriter(rower Rower) *Writer {
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
	case "json":
		wr.json = []string{}
	}
	if rower != nil && (!noHeader || appendOutput) {
		wr.Header(rower)
	}
	return wr
}

func (w *Writer) Header(rower Rower) error {
	switch w.format {
	case "table":
		w.tbl.SetHeader(rower.Header())
	case "csv":
		w.csv.Write(rower.Header())
	}
	return nil
}

func (w *Writer) Row(rower Rower) error {
	switch w.format {
	case "table":
		w.tbl.Append(rower.Row())
	case "csv":
		w.csv.Write(rower.Row())
	case "json":
		buf, err := json.Marshal(rower)
		if err != nil {
			log.WithError(err).Error("failed to marshal json data...")
			return err
		}
		w.json = append(w.json, string(buf))
	}
	return nil
}

func (w *Writer) Flush() {
	switch w.format {
	case "table":
		w.tbl.Render()
	case "csv":
		w.csv.Flush()
	case "json":
		prevData := ""
		if com.IsFile(w.outputFileName) && appendOutput {
			buf, err := ioutil.ReadFile(w.outputFileName)
			if err == nil {
				prevData = string(buf)
			}
		}
		prevData = strings.TrimSpace(prevData)
		js := "["
		if prevData != "" && prevData != "[]" {
			js += strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(prevData, "["), "],"))
		}
		toAdd := strings.TrimSpace(strings.Join(w.json, ","))
		if toAdd != "" {
			js += ",\n" + toAdd + "\n"
		}
		js += "]"
		w.output.Write([]byte(js))
	}
}

func (w *Writer) Close() {
	w.Flush()
	if w.outputFileName != "" {
		com.WriteFile(w.outputFileName, w.output.(*bytes.Buffer).Bytes())
		//pp.Println("Finish writing = ", outputFileName)
	}
}
