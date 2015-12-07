# Lazarus [![Build Status](https://img.shields.io/travis/avadhutp/lazarus/master.svg?style=flat)](https://travis-ci.org/avadhutp/lazarus)

Lazarus plays the most recent *HOT* songs from `r/ListenToThis` subreddit.

![screenshot](http://i.imgur.com/7g6Pscd.png)

# Requirements
1. youtube-dl
2. cvlc

# Settings
Lazarus works off of an `ini` file. This supports the following configs:
1. `tmp_location`: The location where Lazarus can download temporary mp3s. Ideally, this location should be absolute as Lazarus does not deal well with relative paths currently. Also, the trailing slash is necessary. Example, `tmp_location = /tmp/lazarus/`.

# Keyboard shortcuts
Shortcut | Purpose
---------|--------
q | Quit Lazarus
s | Skip song
