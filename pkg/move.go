package libra

const (
	MoveQuiet = iota
	MoveCapture
	MoveOnPassant
	MoveCastle
	MovePromotion
)

type Move struct {
	From byte
	To   byte
	Kind byte
	Code byte
}

func NewMove(from byte, to byte, kind byte) *Move {
	return &Move{
		From: from,
		To:   to,
		Kind: kind,
		Code: 0,
	}
}

func NewPromotionMove(from byte, to byte, kind byte, code byte) *Move {
	return &Move{
		From: from,
		To:   to,
		Kind: kind,
		Code: code,
	}
}
