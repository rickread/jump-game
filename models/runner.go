package models

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	_ "image/png"
)

type Runner struct {
	gameObject
	runningSprites 			[]*ebiten.Image
	jumpingSprite 			*ebiten.Image
	fallingSprite 			*ebiten.Image
	currentRunningSprite  	int
	isJumping            	bool
	isGrounded            	bool
	jumpStartTime           int64
	weight			        float64
	jumpSound  				*audio.Player
}

func CreateRunner(audioContext *audio.Context) *Runner {

	var (
		err error
		runningImage *ebiten.Image
		jumpingImage *ebiten.Image
	)

	runningImage, _, err = ebitenutil.NewImageFromFile("assets/run-right.png")
	if err != nil {
		log.Fatal(err)
	}

	runningSprites := make([]*ebiten.Image, 4)

	for i := 0; i < r_spriteAmount; i++ {
    	runningSprites[i] = runningImage.SubImage(image.Rect(r_width * i, 0, (r_width) * (i+1), r_height)).(*ebiten.Image)
	}

	jumpingImage, _, err = ebitenutil.NewImageFromFile("assets/jump.png")
	if err != nil {
		log.Fatal(err)
	}

	jumpingSprite := jumpingImage.SubImage(image.Rect(0, 0, r_width, r_height)).(*ebiten.Image)
	fallingSprite := jumpingImage.SubImage(image.Rect(r_width, 0, r_width * 2, r_height)).(*ebiten.Image)

	jumpAudioBytes, err := ioutil.ReadFile("assets/jump.mp3")
    if err != nil {
        log.Fatal(err)
    }

	jumpSound, err := mp3.Decode(audioContext, bytes.NewReader(jumpAudioBytes))
	if err != nil {
		log.Fatal(err)
	}
	jumpPlayer, err := audioContext.NewPlayer(jumpSound)
	if err != nil {
		log.Fatal(err)
	}

	return &Runner{
		gameObject{g_screenWidth / 2, f_YPos - r_height, r_velocity},
		runningSprites, 
		jumpingSprite, 
		fallingSprite, 
		0, 
		false, 
		true, 
		0, 
		r_weight,
		jumpPlayer,
	}
}

func (runner *Runner) Update() {

	if runner.isGrounded && ebiten.IsKeyPressed(ebiten.KeySpace) {
		runner.currentRunningSprite = 0;
		runner.isJumping = true
		runner.isGrounded = false
		runner.jumpStartTime = time.Now().UnixMilli()
		runner.jumpSound.Rewind()
		runner.jumpSound.Play()
				
	} else if runner.yPos >= r_groundedYPos {
			runner.isJumping = false
			runner.isGrounded = true
			runner.velocity = r_velocity		
	}
    
    var jumpForce float64
	delta := g_delta

	if runner.isJumping {
		if inpututil.IsKeyJustReleased(ebiten.KeySpace) || runner.yPos <= (g_screenHeight / 4.5) {
			runner.isJumping = false;
		} else if ebiten.IsKeyPressed(ebiten.KeySpace) {

			runner.weight++

			gravity := runner.weight * 0.3

			jumpForce = (r_jumpForce * 2) - (gravity)

			//log.Printf("F: %v", gravity)
			//log.Printf("G: %v", jumpForce)
		}
	}

	if (!runner.isJumping && !runner.isGrounded) {
		if (runner.weight != r_weight) {
				runner.jumpStartTime = time.Now().UnixMilli()
				runner.weight = r_weight
		} else {
			delta = -delta			
			jumpForce = (r_jumpForce) - (float64(time.Now().UnixMilli() - runner.jumpStartTime) * 0.4)
			runner.velocity = r_velocity
		}
	}

	// https://pavcreations.com/jumping-controls-in-2d-pixel-perfect-platformers/
	if !runner.isGrounded {
		runner.velocity = runner.velocity + jumpForce * g_gravity * delta
		runner.yPos = runner.yPos - (runner.velocity * delta)

		if runner.yPos >= r_groundedYPos {
			runner.yPos = r_groundedYPos
		}
	}
	
	
}

func (runner *Runner) Draw(game *Game, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(runner.xPos, runner.yPos)

	if runner.currentRunningSprite >= r_spriteAmount {
		runner.currentRunningSprite = 0;
	}

	var sprite *ebiten.Image

	if (runner.isGrounded) { 
		sprite = runner.runningSprites[runner.currentRunningSprite]		
	} else if (runner.isJumping) {
		sprite = runner.jumpingSprite
	} else {
		sprite = runner.fallingSprite
	}

	if game.count % 3 == 0 && !g_isGameOver {
		runner.currentRunningSprite++;
	}

	screen.DrawImage(sprite, op)
}