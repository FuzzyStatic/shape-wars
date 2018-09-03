/*
 * @Author: Allen Flickinger (allen.flickinger@gmail.com)
 * @Date: 2018-03-31 20:22:37
 * @Last Modified by: FuzzyStatic
 * @Last Modified time: 2018-03-31 20:23:45
 */

package main

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// Audio constants
const (
	WAVSAMPLERATE = 44100
	WAVEXPLOSION  = "./audio/explosion.wav"
	WAVMISSILE    = "./audio/missile.wav"
)

func initAudo() {
	var format beep.Format

	format.SampleRate = WAVSAMPLERATE

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/50))

}

func playWavAudio(file string) {
	var (
		f *os.File
		s beep.StreamSeekCloser
	)

	f, _ = os.Open(file)
	s, _, _ = wav.Decode(f)

	speaker.Play(s)
}
