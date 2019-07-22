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
	"github.com/rai-project/evaluation/writer"
)

type Writer struct {
	outputs         map[string]io.Writer
	outputFileNames map[string]string
	tbl             *tablewriter.Table
	csv             *csv.Writer
	jsonRows        []interface{}
	opts            writer.Options
}

type Rower interface {
	Header(...writer.Option) []string
	Row(...writer.Option) []string
}

type Rowers interface {
	Rower
	Rows(...writer.Option) [][]string
}

func getOutput(outputFileName string) io.Writer {
	var output io.Writer = os.Stdout
	if outputFileName != "" {
		output = &bytes.Buffer{}
	}
	return output
}

func NewWriter(rower Rower, opts ...writer.Option) *Writer {
	baseOpts := []writer.Option{writer.Format(outputFormat)}
	wr := &Writer{
		outputs:         make(map[string]io.Writer),
		outputFileNames: make(map[string]string),
		opts:            writer.NewOptions(append(baseOpts, opts...)...),
	}
	if wr.hasFormat("table") {
		output := getOutput(outputFileName)
		wr.outputs["table"] = output
		tOutputFileName := outputFileName + ".tbl"
		wr.outputFileNames["table"] = tOutputFileName
		wr.tbl = tablewriter.NewWriter(output)
	}
	if wr.hasFormat("csv") {
		output := getOutput(outputFileName)
		wr.outputs["csv"] = output
		wr.outputFileNames["csv"] = outputFileName + ".csv"
		wr.csv = csv.NewWriter(output)
	}
	if wr.hasFormat("json") {
		output := getOutput(outputFileName)
		wr.outputs["json"] = output
		wr.outputFileNames["json"] = outputFileName + ".json"
		wr.jsonRows = []interface{}{}
	}
	if rower != nil && (!noHeader || appendOutput) {
		wr.Header(rower)
	}
	return wr
}

func (w *Writer) Header(rower Rower) error {
	if w.hasFormat("table") {
		w.tbl.SetHeader(rower.Header(writer.FromOptions(w.opts)))
	}
	if w.hasFormat("csv") {
		w.csv.Write(rower.Header(writer.FromOptions(w.opts)))
	}
	return nil
}

func (w *Writer) Row(rower Rower) error {
	if w.hasFormat("table") {
		w.tbl.Append(rower.Row(writer.FromOptions(w.opts)))
	}

	if w.hasFormat("csv") {
		w.csv.Write(rower.Row(writer.FromOptions(w.opts)))
	}

	if w.hasFormat("json") {
		w.jsonRows = append(w.jsonRows, rower)
	}
	return nil
}

func (w *Writer) Rows(rower Rowers) error {
	if w.hasFormat("table") {
		for _, r := range rower.Rows(writer.FromOptions(w.opts)) {
			w.tbl.Append(r)
		}
	}
	if w.hasFormat("csv") {
		for _, r := range rower.Rows(writer.FromOptions(w.opts)) {
			w.csv.Write(r)
		}
	}

	if w.hasFormat("json") {
		for _, r := range rower.Rows(writer.FromOptions(w.opts)) {
			w.jsonRows = append(w.jsonRows, r)
		}
	}
	return nil
}

func (w *Writer) Flush() {
	if w.hasFormat("table") {
		w.tbl.Render()
	}
	if w.hasFormat("csv") {
		w.csv.Flush()
	}
	if w.hasFormat("json") {
		data := []interface{}{}
		outputFileName := w.outputFileNames["json"]
		if com.IsFile(outputFileName) && appendOutput {
			buf, err := ioutil.ReadFile(outputFileName)
			if err != nil {
				log.WithError(err).
					WithField("file", outputFileName).
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

		w.outputs["json"].Write(b)
	}
}

func (w *Writer) Close() {
	w.Flush()
	if outputFileName != "" {
		for format, output := range w.outputs {
			com.WriteFile(w.outputFileNames[format], output.(*bytes.Buffer).Bytes())
			//pp.Println("Finish writing = ", outputFileName)
		}
	}
}

func (w *Writer) hasFormat(name string) bool {
	name = strings.ToLower(name)
	for _, f := range w.opts.Formats {
		if f == name {
			return true
		}
	}
	return false
}
