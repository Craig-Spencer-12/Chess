package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//Constants and variables
const (
	scale        = 1
	screenWidth  = 800
	screenHeight = 800
)

var (
	err error

	//Board values
	lightColor  color.RGBA = color.RGBA{228, 198, 162, 255}
	lightSquare *ebiten.Image
	darkColor   color.RGBA = color.RGBA{189, 129, 55, 255}
	darkSquare  *ebiten.Image
	highColor   color.RGBA = color.RGBA{217, 94, 106, 255}
	highSquare  *ebiten.Image

	//Pieces
	pawnWImg   *ebiten.Image
	rookWImg   *ebiten.Image
	knightWImg *ebiten.Image
	bishopWImg *ebiten.Image
	queenWImg  *ebiten.Image
	kingWImg   *ebiten.Image
	pawnBImg   *ebiten.Image
	rookBImg   *ebiten.Image
	knightBImg *ebiten.Image
	bishopBImg *ebiten.Image
	queenBImg  *ebiten.Image
	kingBImg   *ebiten.Image
)

//----------Helper functions----------
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// TODO: implement niche elemtents of board state
func loadBoardFromFen(fen string, board *Board) {
	mode := 0
	position := 0
	var unit byte
	num := 0

	//Set all castles to false
	board.wkc = false
	board.wqc = false
	board.bkc = false
	board.bqc = false

	for i := 0; i < len(fen); i++ {
		unit = fen[i]

		if unit == ' ' {
			mode++
		} else {
			switch mode {
			case 0:
				num, err = strconv.Atoi(string(unit))
				if num > 0 {
					position += num - 1
					num = 0
				} else {
					switch unit {
					case 'P':
						board.layout[position%8][position/8] = pawnW
					case 'R':
						board.layout[position%8][position/8] = rookW
					case 'N':
						board.layout[position%8][position/8] = knightW
					case 'B':
						board.layout[position%8][position/8] = bishopW
					case 'Q':
						board.layout[position%8][position/8] = queenW
					case 'K':
						board.layout[position%8][position/8] = kingW
					case 'p':
						board.layout[position%8][position/8] = pawnB
					case 'r':
						board.layout[position%8][position/8] = rookB
					case 'n':
						board.layout[position%8][position/8] = knightB
					case 'b':
						board.layout[position%8][position/8] = bishopB
					case 'q':
						board.layout[position%8][position/8] = queenB
					case 'k':
						board.layout[position%8][position/8] = kingB
					default:
						continue
					}
				}
			case 1:
				if unit == 'w' {
					board.whitesTurn = true
				} else {
					board.whitesTurn = false
				}
			case 2:
				switch unit {
				case 'K':
					board.wkc = true
				case 'Q':
					board.wqc = true
				case 'k':
					board.bkc = true
				case 'q':
					board.bqc = true
				}
			case 3:
			case 4:
			case 5:

			}
		}
		position++
	}
}

func createFen(b Board) string {
	fen := ""
	emptyCounter := 0

	//Add board postition to the fen string
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			switch b.layout[i][j] {
			case none:
				emptyCounter++
			case pawnW:
				fen += "P"
				emptyCounter = 0
			case rookW:
				fen += "R"
				emptyCounter = 0
			case knightW:
				fen += "N"
				emptyCounter = 0
			case bishopW:
				fen += "B"
				emptyCounter = 0
			case queenW:
				fen += "Q"
				emptyCounter = 0
			case kingW:
				fen += "K"
				emptyCounter = 0
			case pawnB:
				fen += "p"
				emptyCounter = 0
			case rookB:
				fen += "r"
				emptyCounter = 0
			case knightB:
				fen += "n"
				emptyCounter = 0
			case bishopB:
				fen += "b"
				emptyCounter = 0
			case queenB:
				fen += "q"
				emptyCounter = 0
			case kingB:
				fen += "k"
				emptyCounter = 0
			default:
				//TODO: change to actual error message
				fen += "^error^"
			}
		}
		fen += "/"
	}

	return fen
}

func loadImg(location string) image.Image {
	reader, err := os.Open(location)
	check(err)
	defer reader.Close()

	img, err := png.Decode(reader)
	check(err)

	return img
}

//----------Data Storage----------
type piece int

const (
	none        piece = 0
	pawnW       piece = 1
	rookW       piece = 2
	knightW     piece = 3
	bishopW     piece = 4
	queenW      piece = 5
	kingW       piece = 6
	whiteIfLess piece = 8
	pawnB       piece = 9
	rookB       piece = 10
	knightB     piece = 11
	bishopB     piece = 12
	queenB      piece = 13
	kingB       piece = 14
)

type Board struct {
	layout     [8][8]piece
	enPassant  int
	whitesTurn bool
	halftimer  int
	fulltimer  int

	//castles
	wkc bool
	wqc bool
	bkc bool
	bqc bool
}

type Game struct {
	fen         string
	gameBoard   Board
	highlighted int
	movePreped  piece
	lastHigh    int
}

// //----------Mouse movements----------
// type StrokeSource interface {
// 	Position() (int, int)
// 	IsJustReleased() bool
// }

// type MouseStrokeSource struct{}

// func (m *MouseStrokeSource) Position() (int, int) {
// 	return ebiten.CursorPosition()
// }

// func (m *MouseStrokeSource) IsJustReleased() bool {
// 	return inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft)
// }

//----------Core Ebiten Functionality----------
func InitializeBoard(g *Game) {
	g.fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	loadBoardFromFen(g.fen, &g.gameBoard)
	g.highlighted = -1
	g.movePreped = none

	//Init board squares
	darkSquare = ebiten.NewImage(screenWidth/8, screenHeight/8)
	darkSquare.Fill(darkColor)
	lightSquare = ebiten.NewImage(screenWidth/8, screenHeight/8)
	lightSquare.Fill(lightColor)
	highSquare = ebiten.NewImage(screenWidth/8, screenHeight/8)
	highSquare.Fill(highColor)

	//Load piece images
	pawnWImg = ebiten.NewImageFromImage(loadImg("assets/pieces/pawnW.png"))
	rookWImg = ebiten.NewImageFromImage(loadImg("assets/pieces/rookW.png"))
	knightWImg = ebiten.NewImageFromImage(loadImg("assets/pieces/knightW.png"))
	bishopWImg = ebiten.NewImageFromImage(loadImg("assets/pieces/bishopW.png"))
	queenWImg = ebiten.NewImageFromImage(loadImg("assets/pieces/queenW.png"))
	kingWImg = ebiten.NewImageFromImage(loadImg("assets/pieces/kingW.png"))
	pawnBImg = ebiten.NewImageFromImage(loadImg("assets/pieces/pawnB.png"))
	rookBImg = ebiten.NewImageFromImage(loadImg("assets/pieces/rookB.png"))
	knightBImg = ebiten.NewImageFromImage(loadImg("assets/pieces/knightB.png"))
	bishopBImg = ebiten.NewImageFromImage(loadImg("assets/pieces/bishopB.png"))
	queenBImg = ebiten.NewImageFromImage(loadImg("assets/pieces/queenB.png"))
	kingBImg = ebiten.NewImageFromImage(loadImg("assets/pieces/kingB.png"))
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		//Find which square was just clicked and flag it to be highlighted
		g.lastHigh = g.highlighted
		clickX, clickY := ebiten.CursorPosition()
		col := clickX / (screenWidth / 8)
		row := clickY / (screenHeight / 8)
		g.highlighted = row*8 + col

		//Perform action if players own piece is clicked
		pieceClicked := g.gameBoard.layout[g.highlighted%8][g.highlighted/8]
		if pieceClicked != none && pieceClicked < whiteIfLess && g.gameBoard.whitesTurn {
			g.movePreped = pieceClicked
		} else if pieceClicked > whiteIfLess && !g.gameBoard.whitesTurn {
			g.movePreped = pieceClicked
		} else {
			if g.movePreped != none {
				//TODO: create a function to check if the move is valid
				g.gameBoard.layout[g.highlighted%8][g.highlighted/8] = g.movePreped
				g.gameBoard.layout[g.lastHigh%8][g.lastHigh/8] = none
				g.gameBoard.whitesTurn = !g.gameBoard.whitesTurn
				fmt.Println(createFen(g.gameBoard))
			}

			g.movePreped = none
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//Draw the game board
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(i*screenWidth/8), float64(j*screenHeight/8))
			if i+8*j == g.highlighted {
				screen.DrawImage(highSquare, op)
			} else if (i%2 == 0 && j%2 == 0) || (i%2 == 1 && j%2 == 1) {
				screen.DrawImage(lightSquare, op)
			} else {
				screen.DrawImage(darkSquare, op)
			}
		}
	}

	//Draw the pieces
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {

			var pieceImg *ebiten.Image
			var hasPiece bool = true

			//Pick the correct piece
			switch g.gameBoard.layout[i][j] {
			case none:
				hasPiece = false
			case pawnW:
				pieceImg = pawnWImg
			case rookW:
				pieceImg = rookWImg
			case knightW:
				pieceImg = knightWImg
			case bishopW:
				pieceImg = bishopWImg
			case queenW:
				pieceImg = queenWImg
			case kingW:
				pieceImg = kingWImg
			case pawnB:
				pieceImg = pawnBImg
			case rookB:
				pieceImg = rookBImg
			case knightB:
				pieceImg = knightBImg
			case bishopB:
				pieceImg = bishopBImg
			case queenB:
				pieceImg = queenBImg
			case kingB:
				pieceImg = kingBImg
			default:
				pieceImg = kingBImg //TODO: change to an error
			}

			if hasPiece {
				//Display the piece
				var scaleX float64 = float64((screenWidth / 8.0) / float64(pieceImg.Bounds().Dx()))
				var scaleY float64 = float64((screenHeight / 8.0) / float64(pieceImg.Bounds().Dy()))

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(i*screenWidth/8)*(1/scaleX), float64(j*screenHeight/8)*(1/scaleY))
				op.GeoM.Scale(scaleX, scaleY)
				screen.DrawImage(pieceImg, op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

//----------Main Method----------
func main() {

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Chess")

	game := &Game{}
	InitializeBoard(game)

	err = ebiten.RunGame(game)
	check(err)
}
