package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/DATA-DOG/godog/colors"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "runtime"
    "strconv"
    "time"
)

const steps = 16
const Maxbpm int = 1000
type Instrument string         // a instrument
type Instruments []Instrument // a set of instruments like snare+bass+hitop
type BarSeq []Instruments     // an array representing a set of instruments to be play over a time sequence

const Bass Instrument = "bass"
const Hitop Instrument = "hitop"
const Snare Instrument = "snare"

const BeatBass string = "b.mp3"
const BeatHitop string = "h.mp3"
const BeatSnare string = "s.mp3"
const delimiter string = "~"

const MPlayer string = "/usr/bin/afplay"
var playmusic bool = false
/**
Transforms a map to a BarSeq
input map:
    {
        "hitop" : [1],
        "bass" : [0,4,8,12],
        "snare" : [1],
    }
    transforms to sequence (array)
    =>
    |Bass|Hitop+Snare|_|_|Bass|_|_|_|Bass|_|_|_|Bass|_|_|_|
       0      1      2 3  4           8         12      15 (array indexes)
*/
func mapToBarSeq(m map[Instrument][]uint8) (BarSeq, error) {
	arr := make(BarSeq, steps)
	for k, v := range m {
		for _, slotNum := range v {
			if slotNum < 0 || slotNum >= steps {
				return nil, errors.New(fmt.Sprint("Invalid: ", k, v))
			}
			arr[slotNum] = append(arr[slotNum], k)
		}
	}
	return arr, nil
}

type Track struct {
    DurationSecs float64 `json:"duration-secs"`
	Name string `json:"name"`
	//map of instrument and the various slots it will be played in the steps
	Pmap   map[Instrument][]uint8 `json:"instruments"`
	Tempo  int `json:"tempo"`
	Volume int `json:"volume"`
}

type Bar struct {
	// Example BarSeq array: |Bass|Hitop+Snare|_|_|Bass|_|_|_|Bass|_|_|_|Bass|_|_|_|
	BarSequence BarSeq
	Tick        <-chan time.Time
}

func playMusic(vol int, beatfilepath string) {
    if !playmusic {
        return
    }
    v := strconv.Itoa(vol)
    if err := exec.Command(MPlayer, "-v", v, beatfilepath).Run(); err != nil {
        //log.Print(err)
        return
    }
}
// Example BarSeq array: |Bass|Hitop+Snare|_|_|Bass|_|_|_|Bass|_|_|_|Bass|_|_|_|
func (bar *Bar) viz(durationSecs float64, volume int) {
	barslot := 0
	go func() {
	    for range bar.Tick {
            instruments := bar.BarSequence[barslot]
            fmt.Print("(")
            if len(instruments) == 0 {
                fmt.Print(delimiter)
            }
            for _, k := range instruments {
                switch k {
                case Bass:
                    fmt.Print(colors.Red(delimiter+string(Bass)+delimiter))
                    go playMusic(volume, BeatBass)
                case Snare:
                    fmt.Print(colors.Green(delimiter+string(Snare)+delimiter))
                    go playMusic(volume, BeatSnare)
                case Hitop:
                    fmt.Print(colors.Yellow(delimiter+string(Hitop)+delimiter))
                    go playMusic(volume, BeatHitop)
                default:
                    fmt.Print("...")
                }
            }
            fmt.Print(")")
            barslot++
            if barslot%8 == 0 {
                fmt.Print("\n")
            }
            if barslot > steps-1 {
                barslot = 0
            }
        }
    }()
	if durationSecs >= 1 {
        time.Sleep(time.Duration(durationSecs) * time.Second)
    }
}

func (track *Track) validate() error {
	if track.Name == "" {
		return errors.New("track name cannot be empty")
	}
	if track.Pmap == nil {
		return errors.New("no beats configured")
	}
	if track.Tempo < 1 || track.Tempo > Maxbpm {
		return errors.New("bpm must be between 1 and " + strconv.Itoa(Maxbpm))
	}
    if track.DurationSecs < 1 {
        return errors.New("track duration cannnot be less than 1 sec")
    }
    if track.Volume < 1 {
        track.Volume = 5
    }
	return nil
}

func (track *Track) play() error {
	if err := track.validate(); err != nil {
		return err
	}
	bpm := track.Tempo
	fmt.Println("Track: ", track.Name, "\n", "BPM: ", bpm)
	seq, err := mapToBarSeq(track.Pmap)
	if err != nil {
		return err
	}
	tick := time.Tick(steps * time.Second / time.Duration(bpm))
	bar := &Bar{
		BarSequence: seq,
		Tick:        tick,
	}
	bar.viz(track.DurationSecs, track.Volume)
	return nil
}
func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}
func main() {
    if runtime.GOOS == "darwin" && fileExists(BeatHitop) && fileExists(BeatSnare) && fileExists(BeatBass) && fileExists(MPlayer){
        playmusic = true
    }
    //get the json input file - see example track.json
    if len(os.Args) < 2 {
        fmt.Println("Missing parameter, provide file name!")
        return
    }
    // if we os.Open returns an error then handle it
    var jsonFile *os.File
    var err error
    if jsonFile, err = os.Open(os.Args[1]); err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := jsonFile.Close(); err != nil {
            log.Fatal(err)
        }
    }()
    byteValue, _ := ioutil.ReadAll(jsonFile)
    trk := &Track{}
    if err := json.Unmarshal(byteValue, trk); err != nil {
        log.Fatal(err)
    }
    if err := trk.play(); err != nil {
        log.Fatal(err)
    }
}
