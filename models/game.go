package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"

	"image/color"
	_ "image/png"
)

// assets
var (
	background *ebiten.Image
	floor *Floor
	playerOne *Runner
	ship *Ship
	hurdles Hurdles
	titleFont *truetype.Font
	titleFontFace font.Face
	subTitleFont *truetype.Font
	subTitleFontFace font.Face	
	scoreFontFace font.Face
)

// state
var timeSinceLastHurdle int64
var activeHurdles []HurdlePair

var hurdleDelay int64 = 0

var gameOverTime int64

var musicPlayer *audio.Player

type gameObject struct {	
	xPos, yPos float64
	velocity   float64  
}

type Game struct {
	count int
}

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	background, _, err = ebitenutil.NewImageFromFile("assets/background.png")
	if err != nil {
		log.Fatal(err)
	}

	audioContext := audio.NewContext(48000)

	floor = CreateFloor()
	playerOne = CreateRunner(audioContext)
	ship = CreateShip()
	hurdles = CreateHurdles(audioContext)

	timeSinceLastHurdle = time.Now().UnixMilli() + 5000
	activeHurdles = make([]HurdlePair, 0)
	
    titleFontBytes, err := ioutil.ReadFile("assets/FasterOne-Regular.ttf")
    if err != nil {
        log.Println(err)
        return
    }
    titleFont, err = truetype.Parse(titleFontBytes)
    if err != nil {
        log.Println(err)
        return
    }

	const dpi = 72
	titleFontFace = truetype.NewFace(titleFont, &truetype.Options{
		Size:    72,
		DPI:     dpi,
		Hinting: font.HintingFull,		
	})	

	subTitleFontBytes, err := ioutil.ReadFile("assets/PressStart2P-Regular.ttf")
    if err != nil {
        log.Println(err)
        return
    }

    subTitleFont, err = truetype.Parse(subTitleFontBytes)
    if err != nil {
        log.Println(err)
        return
    }

	subTitleFontFace = truetype.NewFace(subTitleFont, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,		
	})

	scoreFontFace = truetype.NewFace(subTitleFont, &truetype.Options{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,		
	})

	musicBytes, err := ioutil.ReadFile("assets/music.mp3")
	if err != nil {
		log.Fatal(err)
	}
	
	music, err := mp3.Decode(audioContext, bytes.NewReader(musicBytes))
	if err != nil {
		log.Fatal(err)
	}

	musicPlayer, err = audioContext.NewPlayer(music)
	if err != nil {
		log.Fatal(err)
	}
}

func (game *Game) Update() error {
	game.count++

	if (game.count == 100) {
		game.count = 0
	}

	if (!g_isGameOver) {

		playerOne.Update()

		tmp := activeHurdles[:0]
		for _, pair := range activeHurdles {
			pair.Update(game)
			if !pair.IsOffScreen() {
				tmp = append(tmp, pair)
			}

			if (pair.HasHitRunner(playerOne)) {
				pair.hitSound.Rewind()
				pair.hitSound.Play()
				endGame()
			}
		}

		if (g_isGameStart) {
			activeHurdles = tmp
            		
			hurdleDelay = (h_minDelay + h_maxDelay)  / 2

			diffSinceLastHurdle := time.Now().UnixMilli() - timeSinceLastHurdle

			if (diffSinceLastHurdle > hurdleDelay) {
				timeSinceLastHurdle = time.Now().UnixMilli()
				pair := generateRandomHurdles()
				activeHurdles = append(activeHurdles, pair)		

				if (g_universalVelocity < 6) {
					g_universalVelocity += g_velocity
				}

				if (h_minDelay > 750) {
					reducer := rand.Int63n(h_varMinSpeedDelay)

					if (h_minDelay - reducer > 750) {
						h_minDelay -= reducer
					}
				}

				if (h_maxDelay > 750) {					
					reducer := rand.Int63n(h_varMinSpeedDelay)

					if (h_maxDelay - reducer > 750) {
						h_maxDelay -= reducer
					}
				}

				h_varMinSpeedDelay += 3

				log.Printf("%v and %v and %v", h_minDelay, h_maxDelay, hurdleDelay)
				
			}
		} else if (ebiten.IsKeyPressed(ebiten.KeySpace)) {
			startGame()
		}
	
		floor.Update(game)

	} else if (g_isGameOver && ebiten.IsKeyPressed(ebiten.KeySpace)) {
		if (time.Now().UnixMilli() - gameOverTime > g_gameOverDelay) {
			resetGame()
		}
	}

	ship.Update(game)
	
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)
	ship.Draw(game, screen)
	floor.Draw(game, screen)
	playerOne.Draw(game, screen)
	
	if (g_isGameStart) {
		for _, pair := range activeHurdles {
			pair.Draw(game, screen)
		}
		
		if (!g_isGameOver) {
			text.Draw(screen, g_goText, titleFontFace, 60, 130, color.RGBA{ 255, 222, 108, uint8(color.Alpha16{ g_textAlpha }.A) })
			increaseTextTransparency()
		}

	} else {
		text.Draw(screen, g_titleText, titleFontFace, 65, 130, color.RGBA{ 255, 222, 108, uint8(color.Alpha16{ g_textAlpha }.A) })
		text.Draw(screen, g_pressStart, subTitleFontFace, 130, 180, color.RGBA{ 255, 222, 255, uint8(color.Alpha16{ g_textAlpha }.A) })
		decreaseTextTransparency()
	}

	if (g_isGameOver) {
		text.Draw(screen, g_gameOverText, titleFontFace, 165, 130, color.RGBA{ 255, 222, 108, uint8(color.Alpha16{ g_textAlpha }.A) })
		decreaseTextTransparency()
	}

	if g_isGameStart {
		text.Draw(screen, fmt.Sprintf("Score: %v", score), scoreFontFace, 15, 30, color.RGBA{ 255, 222, 255, uint8(color.Alpha16{ 255 }.A) })
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.GetDefaultWindowSize()
}

func (g *Game) GetDefaultWindowSize() (screenWidth, screenHeight int) {
	return g_screenWidth, g_screenHeight
}

func generateRandomHurdles() (hurdle HurdlePair) {
	return hurdles.pairs[rand.Intn(h_spriteAmount)]
}

func resetGame() {
	timeSinceLastHurdle = time.Now().UnixMilli() + 4250
	gameOverTime = 0
	activeHurdles = nil	
	activeHurdles = make([]HurdlePair, 0)
	playerOne.isGrounded = true
	playerOne.isJumping = false
	playerOne.yPos = f_YPos - r_height
	playerOne.currentRunningSprite = 0
	playerOne.weight = r_weight
	g_universalVelocity = 1.0
	h_minDelay = 5000
	h_maxDelay = 8000
	score = 0
	
	startGame()
}

func startGame() {
	g_isGameOver = false
	g_isGameStart = true

	musicPlayer.Rewind()
	musicPlayer.Play()
}

func endGame() {
	g_isGameOver = true	
	gameOverTime = time.Now().UnixMilli()

	musicPlayer.Pause()
}

func decreaseTextTransparency() {	
	g_textAlpha += 2
	if (g_textAlpha > 255) {
		g_textAlpha = 255
	}
}

func increaseTextTransparency() {
	if (g_textAlpha > 0) {	
		g_textAlpha--
	}
}