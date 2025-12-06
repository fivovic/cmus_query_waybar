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

func cmus_query() string {
	var cmd string
	var output string

	cmd_bin := "cmus-remote"
	cmd_args := "-Q"
	cmd = fmt.Sprintf("%s %s", cmd_bin, cmd_args)

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

func cmus_parse(cmus_output string) cmusWaybar {
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

	// handle each status
	switch status {
	case "playing":
		symbol := 0x25B6 // ▶
		song, time, album, percent := _cmus_extract(cmus_output)
		if progress_bar {
			bar := _progressBar(percent, progress_bar_width)
			text = fmt.Sprintf("[%c] %s %s", symbol, song, bar)
		} else {
			text = fmt.Sprintf("[%c] %s", symbol, song)
		}
		tooltip = fmt.Sprintf("%s (%.2f%%)\n%s", time, percent, album)
	case "paused":
		symbol := 0x23F8 // ⏸
		song, time, album, percent := _cmus_extract(cmus_output)
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

	return cmusWaybar{
		Status:  status,
		Text:    text,
		Tooltip: tooltip,
	}
}

func _cmus_extract(cmus_output string) (song string, time string, album string, percent float64) {
	var artist, title, date, position, duration string
	var match_artist, match_title, match_date, match_album, match_position, match_duration bool
	var duration_raw, position_raw int
	var err error

	for line := range strings.SplitSeq(cmus_output, "\n") {
		line = strings.TrimSpace(line)
		if !match_artist && strings.HasPrefix(line, "tag artist ") {
			fields := strings.Fields(line)
			artist = strings.Join(fields[2:], " ")
			match_artist = true
			logger.Debug("found match", "artist", artist)
		}
		if !match_title && strings.HasPrefix(line, "tag title ") {
			fields := strings.Fields(line)
			title = strings.Join(fields[2:], " ")
			match_title = true
			logger.Debug("found match", "title", title)
		}
		if !match_date && strings.HasPrefix(line, "tag date ") {
			fields := strings.Fields(line)
			date = strings.Join(fields[2:], " ")
			match_date = true
			logger.Debug("found match", "date", date)
		}
		if !match_album && strings.HasPrefix(line, "tag album ") {
			fields := strings.Fields(line)
			album = strings.Join(fields[2:], " ")
			match_album = true
			logger.Debug("found match", "album", album)
		}
		if !match_duration && strings.HasPrefix(line, "duration ") {
			fields := strings.Fields(line)
			duration_raw, err = strconv.Atoi(fields[1])
			if err != nil {
				logger.Error("conversion failed", "err", err, "duration", duration_raw)
			}
			duration = _formatSeconds(duration_raw)
			match_duration = true
			logger.Debug("found match", "duration", duration_raw)
		}
		if !match_position && strings.HasPrefix(line, "position ") {
			fields := strings.Fields(line)
			position_raw, err = strconv.Atoi(fields[1])
			if err != nil {
				logger.Error("conversion failed", "err", err, "position", position_raw)
			}
			position = _formatSeconds(position_raw)
			match_position = true
			logger.Debug("found match", "position", position_raw)
		}
	}

	if position_raw > 1 {
		percent = float64(position_raw) / float64(duration_raw) * 100
	} else {
		percent = 0
	}
	song = fmt.Sprintf("%s - %s", artist, title)
	time = fmt.Sprintf("%s/%s", position, duration)
	album = fmt.Sprintf("[%s] %s", date, album)

	return song, time, album, percent
}

func _formatSeconds(seconds int) (formatted string) {
	m := (seconds % 3600) / 60
	s := seconds % 60
	formatted = fmt.Sprintf("%02dm%02ds", m, s) // example: 10m33s
	return formatted
}

func _progressBar(percent float64, width int) (bar string) {
	filledLength := int(percent / 100 * float64(width))
	bar = strings.Repeat("█", filledLength) + strings.Repeat("░", width-filledLength)
	return bar
}
