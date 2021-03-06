package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/csv"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("fcc.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = CSVToCZ64(f, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func CSVToCZ64(in io.Reader, out io.Writer) error {
	r := csv.NewReader(in)
	for {
		var csvbuf, gzbuf, b64buf bytes.Buffer
		w := csv.NewWriter(&csvbuf)
		gw, err := gzip.NewWriterLevel(&gzbuf, gzip.BestCompression)
		if err != nil {
			return err
		}
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = w.Write(record)
		if err != nil {
			return err
		}
		w.Flush()
		_, err = gw.Write(csvbuf.Bytes())
		if err != nil {
			return err
		}
		gw.Flush()
		encoder := base64.NewEncoder(base64.StdEncoding, &b64buf)
		_, err = encoder.Write(gzbuf.Bytes())
		if err != nil {
			return err
		}
		encoder.Close()
		_, err = out.Write(b64buf.Bytes())
		if err != nil {
			return err
		}
		out.Write([]byte("\n"))
	}
	return nil
}
