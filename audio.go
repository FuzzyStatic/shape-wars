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
	WAVEXPLOSION = "./audio/explosion.wav"
	WAVMISSILE   = "./audio/missile.wav"
)

func playWavAudio(file string) {
	go func() {
		var (
			done   chan struct{}
			f      *os.File
			s      beep.StreamSeekCloser
			format beep.Format
		)

		done = make(chan struct{})
		f, _ = os.Open(file)
		s, format, _ = wav.Decode(f)
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/50))
		speaker.Play(beep.Seq(s, beep.Callback(func() {
			close(done)
		})))
		<-done
	}()

	return
}
