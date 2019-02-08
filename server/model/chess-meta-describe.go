package model

import "github.com/team142/angrychess/util"

type MoveDescription struct {
	XDiff          int
	YDiff          int
	Down           bool //Y is decreasing
	Diagonal       bool
	BeingPlaced    bool
	BeingRemoved   bool
	MovingBoards   bool
	PawnOnSpawn    bool
	LastTwoRows    bool
	OtherBoard     bool
	LandingOnPiece *Piece
	PiecesBetween  []*Piece
}

func CalcMoveDescription(game *Game, player *Player, piece *Piece, move *MessageMove) *MoveDescription {
	result := &MoveDescription{}

	//Change in vertical tiles
	result.YDiff = util.Abs(piece.Y - move.ToY)

	//Change in horizontal tiles
	result.XDiff = util.Abs(piece.X - move.ToX)

	//Is the piece moving down the board
	result.Down = piece.Y > move.ToY

	//Is it going from off the board to on the board
	result.BeingPlaced = piece.Cache && !move.Cache

	//Is it going from the board to the cache
	result.BeingRemoved = !piece.Cache && move.Cache

	//Is it moving to another board
	result.MovingBoards = piece.Board != move.Board

	//Is the pawn on it's starting row
	result.PawnOnSpawn = player.Team == 1 && piece.Y == 7 || player.Team == 2 && piece.Y == 2

	//Is the placement on the last two of rows for that players direction
	result.LastTwoRows = (1 == player.Team && 2 >= piece.Y) || (piece.Y >= 7 && 2 == player.Team)

	//Is it not the player's board
	result.OtherBoard = player.Board != move.Board

	//Pieces between the from tile and to tile
	CalcPiecesBetween(game, player, piece, move, result)

Outer:
	for _, pl := range game.Players {
		for _, pi := range pl.Pieces {
			if pi.IsEqual(move) {
				result.LandingOnPiece = pi
				break Outer
			}
		}
	}
	return result

}

func CalcPiecesBetween(game *Game, player *Player, piece *Piece, move *MessageMove, result *MoveDescription) {
	//Don't worry about one tile
	if result.XDiff+result.YDiff <= 1 {
		return
	}
	//Knights don't worry about pieces between
	if piece.Identity == identityKnight {
		return
	}

	//Horizontal moves and is also greater than 1 diff
	if result.XDiff > 1 && result.YDiff == 0 {
		rx1, _, rx2, _ := util.OrderPoints(move.ToX, move.ToY, piece.X, piece.Y)
		for x := rx1 + 1; x < rx2; x++ {
			found, p := game.GetPieceAtPoint(piece.Board, piece.X, piece.Y)
			if found {
				result.PiecesBetween = append(result.PiecesBetween, p)
			}
		}
	} else if result.XDiff == 0 && result.YDiff > 1 {
		//Vertical moves
		_, ry1, _, ry2 := util.OrderPoints(move.ToX, move.ToY, piece.X, piece.Y)
		for y := ry1 + 1; y < ry2; y++ {
			found, p := game.GetPieceAtPoint(piece.Board, piece.X, piece.Y)
			if found {
				result.PiecesBetween = append(result.PiecesBetween, p)
			}
		}
	}

	//TODO: do diagonal moves
	//

}
