# cmus_query_waybar

A simple utility to query the status of the cmus music player from the command line.

Output is returned in a JSON format, specifically designed to be parsed by waybar.

Waybar configuration example:

```json
...
"custom/cmus": {
    "format": "{}",
    "tooltip": true,
    "interval": 1,
    "exec": "cmus_query_waybar -progress-bar -progress-bar-width 30",
    "return-type": "json",
    "on-click": "cmus-remote -Q | rg -q 'status playing' && cmus-remote --pause || cmus-remote --play",
    "on-click-right": "cmus-remote --stop",
    "on-scroll-up": "cmus-remote --prev",
    "on-scroll-down": "cmus-remote --next"
}
...
```

Assumed on-click actions:

| action            | cmus-remote command |
|-------------------|---------------------|
| left click        | play/pause          |
| right click       | stop                |
| scroll up         | previous track      |
| scroll down       | next track          |

## notes

- I used a small golang framework I use for a lot of quick utilities, it likely appears excessive for this use case.
