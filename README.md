# dab-cli   
A fast, minimal Go-based terminal client for the DAB Music API.

The API can be found here: https://sixnine-dotdev.github.io/dab-api-docs/

## Features
- Search for songs
- Stream and play tracks in the terminal
- Download songs via ffmpeg
- Login and token-based access

## Dependencies
- FFMPEG
- MPV

## TODO

### Playback & Queue
- [ ] Add advanced playback controls: shuffle, repeat, loop
- [ ] Keyboard shortcuts for playback: next, previous, pause/resume, volume up/down, mute
- [ ] Seek to specific timestamps in tracks
- [ ] Persistent playlist and queue management
  - [ ] Save/load playlists locally
  - [ ] Multi-queue support (current queue + saved playlists)
  - [ ] Drag-and-drop rearrangement in TUI

### Interactive TUI
- [ ] Multiple panels: library, now-playing, queue, playlist
- [ ] Highlight currently playing track
- [ ] Search bar with live results
- [ ] Progress bar with elapsed/remaining time
- [ ] Dark/light themes and customizable keybindings

### Metadata & Visual Enhancements
- [ ] Display song metadata: artist, album, release year, bitrate, duration
- [ ] Show album art in terminal (ASCII/Unicode or terminal image libs)
- [ ] Add lyrics display in CLI

### Download & Offline Support
- [ ] Batch download playlists or multiple tracks
- [ ] Resume interrupted downloads
- [ ] Option to select download quality/bitrate
- [ ] Cache songs and metadata for offline browsing

### Cloud & User Experience
- [ ] User authentication (token/OAuth, multi-account support)
- [ ] Cloud playlist sync and queue state across sessions
- [ ] Optional scrobbling or listening stats tracking
- [ ] System/terminal notifications for track changes or queue updates

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
