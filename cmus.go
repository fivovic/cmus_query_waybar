package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type cmusWaybar struct {
	Status  string `json:"status"`
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

type cmusMetadata struct {
	artist   string
	title    string
	date     string
	album    string
	duration int
	position int
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
	// set defaults in case cmus is not running
	status := "unknown"
	text := ""
	tooltip := ""
	// check the status first
	for line := range strings.SplitSeq(cmus_output, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "status ") {
			fields := strings.Fields(line)
			status = fields[1]
			break
		}
	}
	// handle each status, extracting the metadata
	switch status {
	case "playing":
		symbol := 0x25B6 // ▶
		song, time, album, percent := _cmusExtract(cmus_output)
		if progress_bar {
			bar := _progressBar(percent, progress_bar_width)
			text = fmt.Sprintf("[%c] %s %s", symbol, song, bar)
		} else {
			text = fmt.Sprintf("[%c] %s", symbol, song)
		}
		tooltip = fmt.Sprintf("%s (%.2f%%)\n%s", time, percent, album)
	case "paused":
		symbol := 0x23F8 // ⏸
		song, time, album, percent := _cmusExtract(cmus_output)
		if progress_bar {
			bar := _progressBar(percent, progress_bar_width)
			text = fmt.Sprintf("[%c] %s %s", symbol, song, bar)
		} else {
			text = fmt.Sprintf("[%c] %s", symbol, song)
		}
		tooltip = fmt.Sprintf("%s (%.2f%%)\n%s", time, percent, album)
	case "stopped":
		symbol := 0x23F9 // ⏹
		text = fmt.Sprintf("[%c]", symbol)
	}
	// the final output will be an easily parsed json for waybar
	return cmusWaybar{
		Status:  status,
		Text:    text,
		Tooltip: tooltip,
	}
}

func _cmusExtract(cmus_output string) (song string, time string, album string, percent float64) {
	meta := _populateMetadata(cmus_output)
	// extract and format the expected outputs
	song = fmt.Sprintf("%s - %s", meta.artist, meta.title)
	time = fmt.Sprintf("%s/%s", _formatSeconds(meta.position), _formatSeconds(meta.duration))
	album = fmt.Sprintf("[%s] %s", meta.date, meta.album)
	percent = _progressPercent(meta.position, meta.duration)
	return song, time, album, percent
}

func _populateMetadata(cmus_output string) cmusMetadata {
	var pop_artist, pop_title, pop_date, pop_album, pop_duration, pop_position bool // track population
	meta := cmusMetadata{}

	for rawLine := range strings.SplitSeq(cmus_output, "\n") {
		line := strings.TrimSpace(rawLine)
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
	return strings.Join(fields[2:], " "), true
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
	filledLength := int(percent / 100 * float64(width))
	bar = strings.Repeat("█", filledLength) + strings.Repeat("░", width-filledLength)
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
