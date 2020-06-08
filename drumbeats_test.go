package main

import "testing"

func TestValidTrackToPlay(t *testing.T) {
    m := map[Instrument][]uint8{}
    m[Bass] = []uint8{0,4,8,12}
    m[Hitop] = []uint8{1}
    m[Snare] = []uint8{2}
    trk := &Track{
        Name:   "Four on the floor 1",
        Pmap:   m,
        Tempo:  0,
        Volume: 10,
    }
    err := trk.validate()
    if err == nil {
        t.Errorf("should not allow tempo less than 1")
    }

    trk.Name = ""
    trk.Pmap = m
    trk.Tempo = 1
    err = trk.validate()
    if err == nil {
        t.Errorf("must not allow track name to be empty")
    }

    trk.Name = "t1"
    trk.Pmap = nil
    trk.Tempo = 1
    err = trk.validate()
    if err == nil {
        t.Errorf("must not allow track map to be nil")
    }
}
func TestMapToBarSeq(t *testing.T) {
    m := map[Instrument][]uint8{}
    m[Bass] = []uint8{0}
    m[Hitop] = []uint8{1}
    m[Snare] = []uint8{2}
    barseq, err := mapToBarSeq(m)
    if err != nil {
        t.Errorf(err.Error())
    }
    if len(barseq) != steps {
        t.Errorf("bar sequence must have %d slots", steps)
    }
    instruments := barseq[0]
    if instruments[0] != Bass {
        t.Errorf("barseq invalid")
    }
    instruments = barseq[1]
    if instruments[0] != Hitop {
        t.Errorf("barseq invalid")
    }
    instruments = barseq[2]
    if instruments[0] != Snare {
        t.Errorf("barseq invalid")
    }
    x  := uint8(steps+1)
    m[Snare] = []uint8{x}
    _, err = mapToBarSeq(m)

    if err == nil {
        t.Errorf("BarSlot of %d is greater than %d", x,steps)
    }
}
