package main

import (
	"flag"
)

var (
	progress_bar       bool
	progress_bar_width int
	log_level          string
	dry_run            bool
)

func parse_flags() {

	flag.BoolVar(&progress_bar, "progress-bar", false, "Calculate a progress bar within the output.\n   ")
	flag.IntVar(&progress_bar_width, "progress-bar-width", 20, "Set the width of the progress bar (in characters).\n   ")
	flag.StringVar(&log_level, "log", "info", "Set the logging level.\n    (options: \"debug\", \"info\", \"warn\", \"error\", \"fatal\")\n   ")
	flag.BoolVar(&dry_run, "dry-run", false, "Run in dry run mode. Logs only what-if output.\n   ")
	flag.Parse()

	if !validate_log_level(log_level) {
		logger.Fatal("invalid log level specified, see \"--help\" for options", "log_level", log_level)
	}
}

func validate_log_level(level string) bool {
	switch level {
	case "debug", "info", "warn", "error", "fatal":
		return true
	default:
		return false
	}
}
