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

	out, err := os.Create("fcc.cgz64")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	r := csv.NewReader(f)
	for {
		var csvbuf, gzbuf, b64buf bytes.Buffer
		w := csv.NewWriter(&csvbuf)
		gw := gzip.NewWriter(&gzbuf)
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		err = w.Write(record)
		if err != nil {
			log.Fatal(err)
		}
		w.Flush()
		i, err := gw.Write(csvbuf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		log.Print(i)
		gw.Flush()
		encoder := base64.NewEncoder(base64.StdEncoding, &b64buf)
		i, err = encoder.Write(gzbuf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		log.Print(i)
		encoder.Close()
		_, err = out.Write(b64buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		out.WriteString("\n")
	}
}
