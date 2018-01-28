package prof

import (
	"time"
	"os"
	"fmt"
	"strings"
	"bytes"
	"os/exec"
	_ "net/http/pprof"
)

const (
	fileSuffix = ".prof"
)

var pprofArgs = []string{
	"tool",
	"pprof",
	"-raw"
}

func raw(opts ...string) ([]byte, error) {

	allArgs := make([]string, len(pprofArgs)+len(opts))
	allArgs = append(allArgs, pprofArgs, opts)

	log.Println("run pprof command: go %s", strings.Join(allArgs, " "))
	
	var buf bytes.Buffer
	cmd := exec.Command("go", allArgs)
	cmd.Stderr = &buf

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run pprof: %s\nSTDERR: %s", err.Error(), buf.String())
	}
	return out, nil
}

func save(filepath string, raw []byte) error {
	f, err := os.Create(filepath+getFileName())
	if err != nil {
		return fmt.Errorf("failed to create profile file: ", err.Error())
	}

	_, err = f.Write(raw)
	if err != nil {
		return fmt.Errorf("failed to write pprof to file: ", err.Error())
	}
	log.Printf("write pprof raw to file: %s success\n", f.Name())
	return nil
}

func getFileName() string {
	return time.Now().Format("20060102_150405")+fileSuffix
}