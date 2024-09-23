package main

import (
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/CorentinB/warc"
	"github.com/spf13/cobra"
)

func ls(cmd *cobra.Command, files []string) {
	for _, filepath := range files {

		f, err := os.Open(filepath)
		if err != nil {
			slog.Error("failed to open file", "file", filepath, "err", err.Error())
			return
		}

		reader, err := warc.NewReader(f)
		if err != nil {
			slog.Error("warc.NewReader failed for file", "file", filepath, "err", err.Error())
			return
		}

		for {
			record, eol, err := reader.ReadRecord()
			if eol {
				break
			}
			if err != nil {
				slog.Error("failed to read all record content", "file", filepath, "err", err.Error())
				return
			}
			lsRecord(record)
		}
	}
}

func lsRecord(record *warc.Record) {
	defer record.Content.Close()

	// Only process Content-Type: application/http; msgtype=response (no reason to process requests or other records)
	if !strings.Contains(record.Header.Get("Content-Type"), "msgtype=response") {
		slog.Debug("skipping record with", "Content-Type", record.Header.Get("Content-Type"), "recordID", record.Header.Get("WARC-Record-ID"))
		return
	}

	if record.Header.Get("WARC-Type") == "revisit" {
		slog.Debug("skipping revisit record", "recordID", record.Header.Get("WARC-Record-ID"))
		return
	}

	fmt.Println(record.Header.Get("WARC-Target-URI"))
}
