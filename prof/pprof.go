package prof

import (
	"bytes"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	pprofArgs = []string{
		"tool",
		"pprof",
		"-raw",
	}

	fileSuffix string
)

// Start pprof to collect web server profile
func Start(opts *Options) {
	log.Println("pprof start: ", *opts)

	// setting options here
	fileSuffix = opts.profSuffix

	ticker := time.NewTicker(time.Duration(opts.tickerSecond) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				b, err := raw([]string{opts.serverURL})
				if err != nil {
					log.Println(err)
					continue
				}
				filename, err := save(opts.filepath, b)
				if err != nil {
					log.Println(err)
					continue
				}
				err = render(filename, opts.targetSVG)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()

}

func raw(opts []string) ([]byte, error) {

	allArgs := make([]string, 0, len(pprofArgs)+len(opts))
	allArgs = append(allArgs, pprofArgs...)
	allArgs = append(allArgs, opts...)

	log.Printf("run pprof command: go %s\n", strings.Join(allArgs, " "))

	var buf bytes.Buffer
	cmd := exec.Command("go", allArgs...)
	cmd.Stderr = &buf

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run pprof: %s\nSTDERR: %s", err.Error(), buf.String())
	}
	return out, nil
}

func save(filepath string, raw []byte) (string, error) {
	f, err := os.Create(filepath + getFileName())
	if err != nil {
		return "", fmt.Errorf("failed to create profile file: %s", err.Error())
	}

	_, err = f.Write(raw)
	if err != nil {
		return "", fmt.Errorf("failed to write pprof to file: %s", err.Error())
	}
	log.Printf("write pprof raw to file: %s success\n", f.Name())
	return f.Name(), nil
}

// render flamegraph to svg file
func render(src, dst string) error {
	var buf bytes.Buffer
	cmd := exec.Command("./build_flamegraph.sh", src, dst)
	cmd.Stderr = &buf

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to render flamegraph file: %s\nSTDERR: %s\n", err.Error(), buf.String())
	}
	log.Println("render flamegraph file success: ", string(out))
	return nil
}

func getFileName() string {
	return time.Now().Format("20060102_150405") + fileSuffix
}
