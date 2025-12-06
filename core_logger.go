package main

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var custom_grey = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8")). // dark grey
	PaddingLeft(4)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	ReportTimestamp: true,
	ReportCaller:    false,
})

func logger_options() {
	// set level
	switch log_level {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	}

	// set level styles
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Bold(true).
		Foreground(lipgloss.Color("14")) // light teal
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Bold(true).
		Padding(0, 1, 0, 0).
		Foreground(lipgloss.Color("45")) // light blue
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN").
		Bold(true).
		Padding(0, 1, 0, 0).
		Foreground(lipgloss.Color("226")) // yellow
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Bold(true).
		Foreground(lipgloss.Color("196")) // red
	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATAL").
		Bold(true).
		Foreground(lipgloss.Color("99")) // purple

	// set custom key/value styles
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	styles.Values["err"] = lipgloss.NewStyle().Bold(true)
	styles.Keys["dry_run"] = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	styles.Values["dry_run"] = lipgloss.NewStyle().Foreground(lipgloss.Color("141")).Bold(true)

	logger.SetStyles(styles)

	if log_level == "debug" {
		logger.SetReportCaller(true)
	}

	logger.Debug("logger loaded", "log_level", log_level, "dry_run", dry_run)
}
