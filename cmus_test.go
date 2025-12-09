package main

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cmusTestCase struct {
	name              string
	input             string
	expectMeta        *cmusMetadata
	expectFormatted   *cmusFormatted
	expectWaybar      *cmusWaybar
	expectProgressBar string
}

//go:embed .tests/cmus_playing
var samplePlaying string

//go:embed .tests/cmus_paused
var samplePaused string

//go:embed .tests/cmus_stopped
var sampleStopped string

//go:embed .tests/cmus_notrunning
var sampleNotRunning string

var cmusTests = []cmusTestCase{
	{
		name:  "playing",
		input: samplePlaying,
		expectMeta: &cmusMetadata{
			status:   "playing",
			artist:   "Wintersun",
			title:    "Sons Of Winter And Stars 2.0",
			date:     "2024",
			album:    "TIME I 2.0",
			duration: 829,
			position: 641,
		},
		expectFormatted: &cmusFormatted{
			song:    "Wintersun - Sons Of Winter And Stars 2.0",
			time:    "10m41s/13m49s",
			album:   "[2024] TIME I 2.0",
			percent: 77.3220747889023,
		},
		expectWaybar: &cmusWaybar{
			Status:  "playing",
			Text:    "[▶] Wintersun - Sons Of Winter And Stars 2.0",
			Tooltip: "10m41s/13m49s (77.32%)\n[2024] TIME I 2.0",
		},
		expectProgressBar: " ███████████████░░░░░",
	},
	{
		name:  "paused",
		input: samplePaused,
		expectMeta: &cmusMetadata{
			status:   "paused",
			artist:   "Parius",
			title:    "Arecibo",
			date:     "2022-10-07",
			album:    "The Signal Heard Throughout Space",
			duration: 790,
			position: 0,
		},
		expectFormatted: &cmusFormatted{
			song:    "Parius - Arecibo",
			time:    "00m00s/13m10s",
			album:   "[2022-10-07] The Signal Heard Throughout Space",
			percent: 0.00,
		},
		expectWaybar: &cmusWaybar{
			Status:  "paused",
			Text:    "[⏸] Parius - Arecibo",
			Tooltip: "00m00s/13m10s (0.00%)\n[2022-10-07] The Signal Heard Throughout Space",
		},
		expectProgressBar: " ░░░░░░░░░░░░░░░░░░░░",
	},
	{
		name:  "stopped",
		input: sampleStopped,
		expectMeta: &cmusMetadata{
			status:   "stopped",
			artist:   "",
			title:    "",
			date:     "",
			album:    "",
			duration: 0,
			position: 0,
		},
		expectFormatted: &cmusFormatted{
			song:    " - ",
			time:    "00m00s/00m00s",
			album:   "[] ",
			percent: 0,
		},
		expectWaybar: &cmusWaybar{
			Status:  "stopped",
			Text:    "[⏹]",
			Tooltip: "",
		},
		expectProgressBar: " ░░░░░░░░░░░░░░░░░░░░",
	},
	{
		name:  "notrunning",
		input: sampleNotRunning,
		expectMeta: &cmusMetadata{
			status:   "unknown",
			artist:   "",
			title:    "",
			date:     "",
			album:    "",
			duration: 0,
			position: 0,
		},
		expectFormatted: &cmusFormatted{
			song:    " - ",
			time:    "00m00s/00m00s",
			album:   "[] ",
			percent: 0,
		},
		expectWaybar: &cmusWaybar{
			Status:  "unknown",
			Text:    "",
			Tooltip: "",
		},
		expectProgressBar: " ░░░░░░░░░░░░░░░░░░░░",
	},
}

func TestCmus(t *testing.T) {
	for _, tc := range cmusTests {
		// metadata
		var meta cmusMetadata
		t.Run(tc.name+"/_populateMetadata", func(t *testing.T) {
			meta = _populateMetadata(tc.input)
			assert.Equal(t, *tc.expectMeta, meta, "unexpected metadata")
		})
		// formatted metadata
		var formatted cmusFormatted
		t.Run(tc.name+"/_formatMetadata", func(t *testing.T) {
			formatted = _formatMetadata(meta)
			assert.Equal(t, *tc.expectFormatted, formatted, "unexpected formatting")
		})
		// final waybar output
		var waybar cmusWaybar
		t.Run(tc.name+"/_populateWaybar", func(t *testing.T) {
			waybar = _populateWaybar(meta.status, formatted)
			assert.Equal(t, *tc.expectWaybar, waybar, "unexpected waybar output")
		})
		// progress bar check
		progress_bar = true
		t.Run(tc.name+"/_progressBar", func(t *testing.T) {
			bar := _progressBar(formatted.percent, 20)
			assert.Equal(t, tc.expectProgressBar, bar, "unexpected progress bar")
		})
		progress_bar = false
	}
}
