# Lazarus [![Build Status](https://drone.io/github.com/avadhutp/lazarus/status.png)](https://drone.io/github.com/avadhutp/lazarus/latest)

Lazarus plays the most recent *HOT* songs from `r/ListenToThis` subreddit.

# Settings
Lazarus works off of an `ini` file. This supports the following configs:
1. `tmp_location`: The location where Lazarus can download temporary mp3s. Ideally, this location should be absolute as Lazarus does not deal well with relative paths currently. Also, the trailing slash is necessary. Example, `tmp_location = /tmp/lazarus/`.

# Requirements
1. youtube-dl
2. cvlc
