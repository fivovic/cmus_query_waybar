package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// first layer of processing that holds the most basic data
type cmusMetadata struct {
	status   string
	artist   string
	title    string
	date     string
	album    string
	duration int
	position int
}

// second layer that holds formatted text in pieces
type cmusFormatted struct {
	song    string
	time    string
	album   string
	percent float64
}

// final layer that brings it all together for waybar
type cmusWaybar struct {
	Status  string `json:"status"`
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

func cmusQuery() string {
	var cmd string
	var output string
	// the actual command that scrapes cmus for information
	cmd_bin := "cmus-remote"
	cmd_args := "-Q"
	cmd = fmt.Sprintf("%s %s", cmd_bin, cmd_args)
	// script not complex enough for proper dryrun handling
	if dry_run {
		logger.Info("running cmus query", "dry_run", cmd)
		return ""
	} else {
		out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
		if err != nil {
			logger.Debug("command failed", "err", err, "output", string(out))
		} else {
			output = string(out)
		}
	}
	return output
}

func cmusParse(cmus_output string) cmusWaybar {
	// extract and populate the metadata
	meta := _populateMetadata(cmus_output)
	// format this metadata into useful pieces
	formatted := _formatMetadata(meta)
	// populate the final output for waybar
	waybar := _populateWaybar(meta.status, formatted)
	return waybar
}

func _populateMetadata(cmus_output string) cmusMetadata {
	var pop_status, pop_artist, pop_title, pop_date, pop_album, pop_duration, pop_position bool // track population
	meta := cmusMetadata{
		status: "unknown", // set default in case cmus is not running
	}
	// this is ordered by expected position in the output
	for rawLine := range strings.SplitSeq(cmus_output, "\n") {
		line := strings.TrimSpace(rawLine)
		if !pop_status {
			if status, ok := _extractStatus(line, "status "); ok {
				meta.status = status
				logger.Debug("found match", "status", status)
			}
		}
		if !pop_duration {
			if duration, ok := _extractNumber(line, "duration "); ok {
				meta.duration = duration
				logger.Debug("found match", "duration", duration)
			}
		}
		if !pop_position {
			if position, ok := _extractNumber(line, "position "); ok {
				meta.position = position
				logger.Debug("found match", "position", position)
			}
		}
		if !pop_artist {
			if artist, ok := _extractTag(line, "tag artist "); ok {
				meta.artist = artist
				logger.Debug("found match", "artist", artist)
			}
		}
		if !pop_album {
			if album, ok := _extractTag(line, "tag album "); ok {
				meta.album = album
				logger.Debug("found match", "album", album)
			}
		}
		if !pop_title {
			if title, ok := _extractTag(line, "tag title "); ok {
				meta.title = title
				logger.Debug("found match", "title", title)
			}
		}
		if !pop_date {
			if date, ok := _extractTag(line, "tag date "); ok {
				meta.date = date
				logger.Debug("found match", "date", date)
			}
		}
	}
	return meta
}

func _formatMetadata(meta cmusMetadata) cmusFormatted {
	formatted := cmusFormatted{}
	// extract and format the expected outputs
	formatted.song = fmt.Sprintf("%s - %s", meta.artist, meta.title)
	formatted.time = fmt.Sprintf("%s/%s", _formatSeconds(meta.position), _formatSeconds(meta.duration))
	formatted.album = fmt.Sprintf("[%s] %s", meta.date, meta.album)
	formatted.percent = _progressPercent(meta.position, meta.duration)
	return formatted
}

func _populateWaybar(status string, formatted cmusFormatted) cmusWaybar {
	waybar := cmusWaybar{}
	waybar.Status = status
	// handle each status, populating the final output for waybar
	switch status {
	case "playing":
		symbol := 0x25B6 // ▶
		bar := _progressBar(formatted.percent, progress_bar_width)
		waybar.Text = fmt.Sprintf("[%c] %s%s", symbol, formatted.song, bar)
		waybar.Tooltip = fmt.Sprintf("%s (%.2f%%)\n%s", formatted.time, formatted.percent, formatted.album)
	case "paused":
		symbol := 0x23F8 // ⏸
		bar := _progressBar(formatted.percent, progress_bar_width)
		waybar.Text = fmt.Sprintf("[%c] %s%s", symbol, formatted.song, bar)
		waybar.Tooltip = fmt.Sprintf("%s (%.2f%%)\n%s", formatted.time, formatted.percent, formatted.album)
	case "stopped":
		symbol := 0x23F9 // ⏹
		waybar.Text = fmt.Sprintf("[%c]", symbol)
	}
	return waybar
}

func _extractTag(line string, prefix string) (string, bool) {
	// catch: no match
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	// catch: not enough fields
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return "", false
	}
	// extract: all fields after the prefix
	value := strings.Join(fields[2:], " ")
	return value, true
}

func _extractStatus(line string, prefix string) (string, bool) {
	// catch: no match
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	// catch: not enough fields
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return "", false
	}
	// extract: specific status field
	value := fields[1]
	return value, true
}

func _extractNumber(line string, prefix string) (int, bool) {
	// catch: no match
	if !strings.HasPrefix(line, prefix) {
		return 0, false
	}
	// catch: not enough fields
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, false
	}
	// extract: convert to integer
	value, err := strconv.Atoi(fields[1])
	if err != nil {
		logger.Error("conversion failed", "err", err, "field", value)
		return 0, false
	}
	return value, true
}

func _formatSeconds(seconds int) (value string) {
	// example: 633 --> 10m33s
	m := (seconds % 3600) / 60
	s := seconds % 60
	value = fmt.Sprintf("%02dm%02ds", m, s)
	return value
}

func _progressBar(percent float64, width int) (bar string) {
	// example: 40% --> ████████░░░░░░░░░░░░
	if progress_bar {
		filledLength := int(percent / 100 * float64(width))
		// keep the leading space because without the bar flag you don't want a trailing one
		bar = " " + strings.Repeat("█", filledLength) + strings.Repeat("░", width-filledLength)
	} else {
		bar = ""
	}
	return bar
}

func _progressPercent(position int, duration int) float64 {
	// catch: avoid division by zero
	if position <= 0 || duration <= 0 {
		return 0
	}
	// calc: percentage
	value := float64(position) / float64(duration) * 100
	return value
}
