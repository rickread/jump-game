package models

// Game
const (
	g_titleText = "Go Runner Go!!!"
	g_goText = "Go!"
	g_gameOverText = "Game Over"
	g_pressStart = "Press [SPACE] to start"
	g_gameOverDelay = 2250
    g_screenWidth, g_screenHeight = 800.0, 310.0
	g_delta = 0.016
	g_gravity = 300.0
	g_velocity = 0.1		
)

var g_isGameStart = false
var g_isGameOver = false
var g_textAlpha uint16 = 0

var g_universalVelocity = 1.0

var score = 0

// Floor
const (
	f_width, f_height = 300.0, 34.0
	f_spriteAmount = 5
	f_velocity = 4.0
	f_YPos = g_screenHeight - f_height
)

// Runner
const (
	r_width, r_height = 55.0, 74.0
	r_spriteAmount = 4	
	r_jumpForce = 6.0
	r_velocity = 5
	r_weight = 0
	r_groundedYPos = g_screenHeight - f_height - r_height
)

// Hurdles
type hurdleSize struct {
	back float64
	front float64
}

const (
	h_spriteAmount = 4
	h_startPosX = -150
	h_width = 12
)

var h_minDelay int64 = 5000
var h_maxDelay int64 = 8000

var h_varMinSpeedDelay int64 = 1250

func h_getSizes() (map[int]hurdleSize) {
	return map[int]hurdleSize{
		0: {62, 74},
		1: {84, 104},
		2: {95, 119},
		3: {105, 134},
	}
}