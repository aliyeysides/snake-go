package main

import (
	"fmt"
	"image/color"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const screenWidth int = 800
const screenHeight int = 500

const cellHeight int32 = 25
const cellWidth int32 = 25

const gridCols int = 20
const gridRows int = 20

type Grid [gridRows][gridCols]string

type Coordinate [2]int
type Snake []Coordinate

var grid Grid
var snake Snake = make(Snake, 3)
var score int = 0
var speed int = 10

var food *Coordinate

type Direction int
type GameStatus int

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

const (
	LIVE GameStatus = iota
	OVER
	PAUSED
)

func (d Direction) String() string {
	return []string{"UP", "RIGHT", "DOWN", "LEFT"}[d]
}

func (g GameStatus) String() string {
	return []string{"LIVE", "OVER", "PAUSED"}[g]
}

var direction Direction = RIGHT
var nextDirection Direction = RIGHT
var gameStatus GameStatus = LIVE

func InitSnake() Snake {
	snake = make(Snake, 3)
	snake[0] = Coordinate{1, 1}
	snake[1] = Coordinate{1, 2}
	snake[2] = Coordinate{1, 3}

	return snake
}

func SetFood() {
	if food == nil {
		food = &Coordinate{rand.Intn(gridRows), rand.Intn(gridCols)}
	}
}

func SetSpeed() {
	if score > 0 && score <= 50 {
		speed = 9
	}

	if score > 50 && score <= 100 {
		speed = 8
	}

	if score > 100 && score <= 150 {
		speed = 7
	}

	if score > 150 && score <= 200 {
		speed = 6
	}

	if score > 200 {
		speed = 5
	}
}

func UpdateScore() {
	score += 10
	SetSpeed()
}

func DetectKeyPress() {
	if rl.IsKeyPressed(rl.KeyUp) && direction != DOWN {
		nextDirection = UP
	}

	if rl.IsKeyPressed(rl.KeyRight) && direction != LEFT {
		nextDirection = RIGHT
	}

	if rl.IsKeyPressed(rl.KeyDown) && direction != UP {
		nextDirection = DOWN
	}

	if rl.IsKeyPressed(rl.KeyLeft) && direction != RIGHT {
		nextDirection = LEFT
	}
}

func IsFoodSquare(square Coordinate) bool {
	squareRow := square[0]
	squareCol := square[1]

	if food != nil && squareRow == food[0] && squareCol == food[1] {
		return true
	}

	return false
}

func IsSnakeSquare(snake Snake, square Coordinate) bool {
	for _, segment := range snake {
		if segment[0] == square[0] && segment[1] == square[1] {
			return true
		}
	}
	return false
}

func IsBorderSquare(square Coordinate) bool {
	squareRow := square[0]
	squareCol := square[1]

	if squareRow < 0 || squareRow >= gridRows || squareCol < 0 || squareCol >= gridCols {
		return true
	}

	return false
}

func MoveSnake(s *Snake) {
	snake := *s
	head := snake[len(snake)-1]

	direction = nextDirection

	// Calculate new head position
	var newHead Coordinate
	switch direction {
	case UP:
		newHead = Coordinate{head[0] - 1, head[1]}
	case RIGHT:
		newHead = Coordinate{head[0], head[1] + 1}
	case DOWN:
		newHead = Coordinate{head[0] + 1, head[1]}
	case LEFT:
		newHead = Coordinate{head[0], head[1] - 1}
	}

	// Check for collisions before updating the snake
	if IsBorderSquare(newHead) || IsSnakeSquare(snake, newHead) {
		gameStatus = OVER
		return
	}

	// Add new head to the snake
	snake = append(snake, newHead)

	// Check if food is eaten
	if IsFoodSquare(newHead) {
		food = nil
		UpdateScore()
		SetFood()
		// Do not remove tail to grow the snake
	} else {
		// Remove tail to move the snake forward
		snake = snake[1:]
	}

	// Update the snake
	*s = snake
}

func ClearGrid() Grid {
	for i := 0; i < gridRows; i++ {
		for j := 0; j < gridCols; j++ {
			grid[i][j] = "EMPTY"
		}
	}

	return grid
}

func DrawSquare(offsetX int32, offsetY int32, color color.RGBA) {
	rl.DrawRectangle(offsetX, offsetY, cellWidth, cellHeight, color)
	rl.DrawLine(offsetX, offsetY, offsetX+cellWidth, offsetY, rl.Black)
	rl.DrawLine(offsetX, offsetY, offsetX, offsetY+cellHeight, rl.Black)
	rl.DrawLine(offsetX+cellWidth, offsetY+cellHeight, offsetX-cellWidth, offsetY+cellHeight, rl.Black)
	rl.DrawLine(offsetX+cellWidth, offsetY+cellHeight, offsetX+cellWidth, offsetY-cellHeight, rl.Black)
}

func DrawEmptySquare(offsetX int32, offsetY int32) {
	DrawSquare(offsetX, offsetY, rl.White)
}

func DrawSnakeSquare(offsetX int32, offsetY int32) {
	DrawSquare(offsetX, offsetY, rl.Black)
}

func DrawFoodSquare(offsetX int32, offsetY int32) {
	DrawSquare(offsetX, offsetY, rl.Red)
}

func DrawBoard(s *Snake) {
	snake := *s

	// Clear Grid
	grid := ClearGrid()

	// Add Snake to Grid
	for _, segment := range snake {
		row := segment[0]
		col := segment[1]

		// Create Grid Boundary
		if row >= 0 && row < gridRows && col >= 0 && col < gridCols {
			grid[row][col] = "SNAKE"
		}
	}

	// Add Food to Grid
	if food != nil {
		grid[food[0]][food[1]] = "FOOD"
	}

	// Draw Grid
	for i := range grid {
		for j := range grid[i] {
			offsetX := int32(j) * cellWidth
			offsetY := int32(i) * cellHeight

			switch grid[i][j] {
			case "EMPTY":
				DrawEmptySquare(offsetX, offsetY)
			case "SNAKE":
				DrawSnakeSquare(offsetX, offsetY)
			case "FOOD":
				DrawFoodSquare(offsetX, offsetY)
			}
		}
	}
}

func DrawScore() {
	rl.DrawText(fmt.Sprint("Score: ", score), int32(gridCols)*cellWidth+cellWidth, 10, 24, rl.Green)
}

func GameOver() {
	rl.DrawText(fmt.Sprint("Game Over!"), int32(screenWidth)/2-100, int32(screenHeight)/2-100, 36, rl.Red)
	rl.DrawText(fmt.Sprint("press Enter to restart"), int32(screenWidth)/2-150, int32(screenHeight)/2-50, 24, rl.Red)
}

func RestartGame() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		gameStatus = LIVE
		direction = RIGHT
		nextDirection = RIGHT
		snake = InitSnake()
		food = nil
		SetFood()
		score = 0
		speed = 10
	}
}

func GameLoop(s *Snake, frameCounter *int) {
	DetectKeyPress()
	SetFood()

	DrawBoard(s)

	if *frameCounter%speed == 0 && gameStatus == LIVE {
		MoveSnake(s)
	}

	if gameStatus == OVER {
		GameOver()
		RestartGame()
	}

	DrawScore()
}

func main() {
	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Go Snake")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	InitSnake()
	frameCounter := 0

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		GameLoop(&snake, &frameCounter)
		frameCounter++

		rl.EndDrawing()
	}
}
