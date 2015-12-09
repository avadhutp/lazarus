# Lazarus [![Build Status](https://img.shields.io/travis/avadhutp/lazarus/master.svg?style=flat)](https://travis-ci.org/avadhutp/lazarus) [![CodeCov](https://img.shields.io/codecov/c/github/avadhutp/lazarus.svg?style=flat)](https://codecov.io/github/avadhutp/lazarus) [![GoDoc](https://godoc.org/github.com/avadhutp/lazarus?status.png)](https://godoc.org/github.com/avadhutp/lazarus)

Lazarus plays the most recent *HOT* songs from `r/ListenToThis` subreddit.

![screenshot](http://i.imgur.com/7g6Pscd.png)

# Requirements
1. youtube-dl
2. a terminal-based mp3 players like
  * afplayer: on OS X
  * cvlc: on Linux (:warning:note: While using cvlc specify the command as `cvlc --play-and-exit` to avoid stalling Lazarus.)
  * mplayer

# Settings
Lazarus works off of an `ini` file. This supports the following configs:

1. `tmp_location`: The location where Lazarus can download temporary mp3s. Ideally, this location should be absolute as Lazarus does not deal well with relative paths currently. Also, the trailing slash is necessary. Example, `tmp_location = /tmp/lazarus/`.
2. `player_cmd`: The command-line music player to use for playback. This needs to be installed on your system and accessible by the current user. Try to use one which supports a wide range of formats suchas `m4a`, `mp3`, `opus`, etc. Also, specify the required arguments, if any, along with the command. For example, `player_cmd = cvlc --play-and-exit`.

# Keyboard shortcuts
Shortcut | Purpose
---------|--------
q | Quit Lazarus
s | Skip song
