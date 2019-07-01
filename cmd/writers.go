package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/Unknwon/com"
	"github.com/olekukonko/tablewriter"
)

//easyjson:json
type Writer struct {
	format         string
	output         io.Writer
	outputFileName string
	tbl            *tablewriter.Table
	csv            *csv.Writer
	jsonRows       []interface{}
}

type Rower interface {
	Header() []string
	Row() []string
}

type Rowers interface {
	Rower
	Rows() [][]string
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
		wr.jsonRows = []interface{}{}
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
		w.jsonRows = append(w.jsonRows, rower)
	}
	return nil
}

func (w *Writer) Rows(rower Rowers) error {
	switch w.format {
	case "table":
		for _, r := range rower.Rows() {
			w.tbl.Append(r)
		}
	case "csv":
		for _, r := range rower.Rows() {
			w.csv.Write(r)
		}
	case "json":
		for _, r := range rower.Rows() {
			w.jsonRows = append(w.jsonRows, r)
		}
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
		data := []interface{}{}
		if com.IsFile(w.outputFileName) && appendOutput {
			buf, err := ioutil.ReadFile(w.outputFileName)
			if err != nil {
				log.WithError(err).
					WithField("file", w.outputFileName).
					Error("failed to read output file")
				return
			}
			if err := json.Unmarshal(buf, &data); err != nil {
				log.WithError(err).Error("failed to unmarshal data")
				return
			}
		}

		data = append(data, w.jsonRows...)

		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.WithError(err).Error("failed to marshal indent data")
			return
		}

		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

		w.output.Write(b)
	}
}

func (w *Writer) Close() {
	w.Flush()
	if w.outputFileName != "" {
		com.WriteFile(w.outputFileName, w.output.(*bytes.Buffer).Bytes())
		//pp.Println("Finish writing = ", outputFileName)
	}
}
