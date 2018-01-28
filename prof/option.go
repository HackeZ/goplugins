package prof

import (
	"fmt"
	"net/url"
	"os"
)

const (
	defaultFileSuffix   = ".prof"
	defaultTickerSecond = 30
	defaultFilepath     = "./"
	defaultServerURL    = "http://localhost:8080"
)

// Option for setting options functional
type Option func(o *Options)

// Options for get profile files from web server and save them
type Options struct {
	serverURL    string // the url with ip and port for web server
	profSuffix   string // the url of pprof profile
	tickerSecond int32  // the second of cllection profile
	filepath     string // the filepath to save profile
}

func (o *Options) String() string {
	return fmt.Sprintf("%#v", *o)
}

// NewOption return a avilable option with default value
func NewOption(opt ...Option) *Options {
	opts := &Options{
	// init base type which mush be malloc, like map/slice/channel
	}

	for _, o := range opt {
		o(opts)
	}

	if opts.serverURL == "" {
		opts.serverURL = defaultServerURL
	}

	if opts.profSuffix == "" {
		opts.profSuffix = defaultFileSuffix
	}

	if opts.tickerSecond == 0 {
		opts.tickerSecond = defaultTickerSecond
	}

	if opts.filepath == "" {
		opts.filepath = defaultFilepath
	}

	return opts
}

// SetServerURL setting your web server url with ip and port
func SetServerURL(ru string) Option {
	u, e := url.Parse(ru)
	if e != nil {
		panic(fmt.Sprintf("invalid url: %s\n", e.Error()))
	}

	return func(o *Options) {
		o.serverURL = u.Path
	}
}

func SetProfSuffix(suf string) Option {
	if suf[0] != '.' {
		panic("invalid pprof file suffix which must start with '.'")
	}

	return func(o *Options) {
		o.profSuffix = suf
	}
}

func SetTickerSencond(sec int32) Option {
	if sec <= 0 {
		panic("invalid ticker second must be a natural number")
	}

	return func(o *Options) {
		o.tickerSecond = sec
	}
}

func SetFilepath(path string) Option {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		panic(fmt.Sprintf("file path: %s not exist\n", path))
	}

	return func(o *Options) {
		o.filepath = path
	}
}
