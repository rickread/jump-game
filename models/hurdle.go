package models

import (
	"bytes"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/png"
)

type hurdle struct {
	gameObject
	image 		*ebiten.Image	
	size		float64
}

type HurdlePair struct {
	front 		hurdle
	back		hurdle
	isScored    bool
	hitSound  	*audio.Player
}

type Hurdles struct {
	pairs		[]HurdlePair
}

func CreateHurdles(audioContext *audio.Context) Hurdles {

	var (
		err error
		imageFront *ebiten.Image
		imageBack *ebiten.Image
	)

	hurdlePairs := make([]HurdlePair, h_spriteAmount)
	hurdleSizes := h_getSizes()


	for i := 0; i < h_spriteAmount; i++ {
		imageFront, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("assets/hurdle-front-%v.png", i))
		if err != nil {
			log.Fatal(err)
		}
	
		imageBack, _, err = ebitenutil.NewImageFromFile(fmt.Sprintf("assets/hurdle-back-%v.png", i))
		if err != nil {
			log.Fatal(err)
		}

		size := hurdleSizes[i]

		front := hurdle{gameObject{h_startPosX * 2, f_YPos - size.front, f_velocity + 1}, imageFront, size.front}
		back := hurdle{gameObject{h_startPosX, f_YPos - size.back, f_velocity}, imageBack, size.back}

		hurdleAudioBytes, err := ioutil.ReadFile("assets/hit.mp3")
		if err != nil {
			log.Fatal(err)
		}
	
		hurdleSound, err := mp3.Decode(audioContext, bytes.NewReader(hurdleAudioBytes))
		if err != nil {
			log.Fatal(err)
		}
		hurdlePlayer, err := audioContext.NewPlayer(hurdleSound)
		if err != nil {
			log.Fatal(err)
		}

		hurdlePairs[i] = HurdlePair{front, back, false, hurdlePlayer}
	}

	return Hurdles{hurdlePairs}
}

func (pair *HurdlePair) Update(game *Game) {
	if (game.count % 2 == 0) {
		pair.front.xPos += pair.front.velocity * g_universalVelocity
		pair.back.xPos += pair.back.velocity * g_universalVelocity
	}
}

func (pair *HurdlePair) Draw(game *Game, screen *ebiten.Image) {
	backOp := &ebiten.DrawImageOptions{}
	backOp.GeoM.Translate(pair.back.xPos, pair.back.yPos)
	screen.DrawImage(pair.back.image, backOp)

	lineColor := color.RGBA{255, 37, 37, uint8(color.Alpha16{ 255 }.A) }
	
	hurdleBackBehind := pair.back.xPos + 2
	hurdleBackFront := pair.back.xPos + h_width - 2
	
	hurdleFrontBehind := pair.front.xPos + 2
	hurdleFrontFront := pair.front.xPos + h_width - 2

	hurdleRatio := pair.front.size / pair.back.size

	for i := 1.0; i < pair.back.size / 5; i++ {
		ebitenutil.DrawLine(screen, hurdleBackBehind, pair.back.yPos + (5 * i), hurdleFrontBehind, (pair.front.yPos + (5 * i) * hurdleRatio), lineColor)
		ebitenutil.DrawLine(screen, hurdleBackFront, pair.back.yPos + (5 * i), hurdleFrontFront, (pair.front.yPos + (5 * i) * hurdleRatio), lineColor)
	}
	frontOp := &ebiten.DrawImageOptions{}
	frontOp.GeoM.Translate(pair.front.xPos, pair.front.yPos)
	screen.DrawImage(pair.front.image, frontOp)
	
}

func (pair *HurdlePair) IsOffScreen() (bool) {
	return pair.back.xPos > g_screenWidth
}

func (pair *HurdlePair) HasHitRunner(runner *Runner) (bool) {

	runnerFront := runner.xPos + 14
	runnerBehind := runner.xPos + (r_width / 2)
	runnerShoes := runner.yPos + (r_height - 5)

	hurdleFront := pair.front.xPos + h_width - 2
	hurdleBehind := pair.front.xPos + 2
	hurdleBar := pair.front.yPos + 15

	if (hurdleFront < runnerFront || hurdleBehind > runnerBehind) {
		return false

	}

	if (hurdleFront >= runnerFront && runnerShoes > hurdleBar) {		
		return true
	}

	if (!pair.isScored) {
		pair.isScored = true
		score += int(pair.front.size) * 10
	}

	return false
}