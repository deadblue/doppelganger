package daemon

import "github.com/spf13/pflag"

type Flags struct {
	ConfigFile string
}

func (f *Flags) ParseCommandLine() {
	pflag.StringVarP(&f.ConfigFile, "config-file", "C", "", "Config file for daemon.")
	pflag.Parse()
}
