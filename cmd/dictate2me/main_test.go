package main

import (
    "bytes"
    "os"
    "testing"
)

func TestPrintVersionAndUsage(t *testing.T) {
    // Capture stdout
    old := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    printVersion()
    printUsage()

    w.Close()
    var buf bytes.Buffer
    _, _ = buf.ReadFrom(r)
    os.Stdout = old

    out := buf.String()
    if out == "" {
        t.Fatalf("expected output from version and usage, got empty")
    }
}
