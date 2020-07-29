package cli

import "github.com/spf13/pflag"

type Flags struct {
	// Server address
	ServerAddr string
	// Command task
	TaskCmd  string
	TaskArgs []string
	// HTTP task
	TaskUrl     string
	TaskMethod  string
	TaskHeaders []string
}

func (f *Flags) ParseCommandLine() {
	pflag.StringVarP(&f.ServerAddr, "server", "s", "/tmp/doppelganger.sock", "Server address")

	pflag.StringVarP(&f.TaskCmd, "cmd-name", "n", "", "Command name")
	pflag.StringArrayVarP(&f.TaskArgs, "cmd-arg", "a", []string{}, "Command argument")

	pflag.StringVarP(&f.TaskUrl, "http-url", "u", "", "HTTP URL")
	pflag.StringVarP(&f.TaskMethod, "http-method", "m", "", "HTTP Method")
	pflag.StringArrayVarP(&f.TaskHeaders, "http-header", "h", []string{}, "HTTP Header")

	pflag.Parse()
}
