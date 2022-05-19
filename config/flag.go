package config

import "flag"

func Flag(Args *CommandLineArgs) {
	flag.StringVar(&Args.Targets, "t", "", "IP and CNAME to check, separated by ','")
	flag.StringVar(&Args.Filepath, "f", "", "CDN data")
	flag.Parse()
}
