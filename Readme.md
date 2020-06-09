* go get github.com/DATA-DOG/godog/colors
* go test
* go run drumbeats.go track.json

Notes:
1. If the plaform is darwin (macos) - the program will attempt to play real music 
using afplay command. b.mp3, h.mp3 and s.mp3 should be in the executable dir
2. Please see track.json for sample input
3. Valid values: 0-15 for step slot for a given instrument in the play sequence
4. Tempo: Valid values 1-1000
5. duration-secs: Specifies how long to run the pattern and music in seconds
