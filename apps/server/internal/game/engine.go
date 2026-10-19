package game

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

const defaultBoardSize = 16

var (
	ErrInvalidConfig = errors.New("invalid game configuration")
	ErrUnknownCard   = errors.New("card instance is not on the board")
	ErrScoreOverflow = errors.New("score or combo counter overflow")
)

// Config contains baseline gameplay values that may be tuned by callers.
type Config struct {
	BoardSize        int
	PointsPerCorrect uint64
}

// DefaultConfig returns the documented Speed Battle baseline values.
func DefaultConfig() Config {
	return Config{BoardSize: defaultBoardSize, PointsPerCorrect: 1}
}

// CardInstance has an identity distinct from its visual variation.
type CardInstance struct {
	ID      string
	Variant Variant
}

// Snapshot is a copy-safe view of the authoritative local engine state.
type Snapshot struct {
	Reference        Variant
	Cards            []CardInstance
	Score            uint64
	Combo            uint64
	BoardVersion     uint64
	ReferenceVersion uint64
}

// Outcome describes a successfully processed gameplay selection.
type Outcome string

const (
	OutcomeCorrect   Outcome = "correct"
	OutcomeIncorrect Outcome = "incorrect"
)

// PickResult records a gameplay outcome. Technical errors are returned as Go
// errors and do not create a gameplay outcome.
type PickResult struct {
	Outcome Outcome
	State   Snapshot
}

// Engine owns one player's private board and scoring state. It is not safe for
// concurrent use; a Room event loop or another caller must serialize actions.
type Engine struct {
	config           Config
	rng              *rand.Rand
	nextCardID       uint64
	reference        Variant
	cards            []CardInstance
	score            uint64
	combo            uint64
	boardVersion     uint64
	referenceVersion uint64
}

// NewEngine creates a reproducible initial board from a caller-provided seed.
// The caller must keep the seed and random source private from clients.
func NewEngine(config Config, seed int64) (*Engine, error) {
	if config.BoardSize < 1 || config.PointsPerCorrect == 0 {
		return nil, fmt.Errorf("%w: board size and correct-pick points must be positive", ErrInvalidConfig)
	}
	engine := &Engine{
		config: config,
		rng:    newSeededRandom(seed),
	}
	engine.reference = Variant(engine.rng.Intn(VariantCount))
	engine.cards = make([]CardInstance, config.BoardSize)
	for index := range engine.cards {
		engine.cards[index] = engine.newCard(Variant(engine.rng.Intn(VariantCount)))
	}
	if !hasMatch(engine.reference, engine.cards) {
		engine.cards[len(engine.cards)-1] = engine.newCard(engine.randomMatchingVariant(engine.reference))
	}
	return engine, nil
}

func newSeededRandom(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// Snapshot returns a deep copy of the engine's board state.
func (e *Engine) Snapshot() Snapshot {
	cards := append([]CardInstance(nil), e.cards...)
	return Snapshot{
		Reference:        e.reference,
		Cards:            cards,
		Score:            e.score,
		Combo:            e.combo,
		BoardVersion:     e.boardVersion,
		ReferenceVersion: e.referenceVersion,
	}
}

// Pick processes a card selection. Unknown IDs are technical errors and leave
// all state, including the private random source, unchanged. Caller layers
// enforce cooldown and stale-version rejection.
func (e *Engine) Pick(cardID string) (PickResult, error) {
	index := -1
	for candidateIndex, card := range e.cards {
		if card.ID == cardID {
			index = candidateIndex
			break
		}
	}
	if index < 0 {
		return PickResult{}, ErrUnknownCard
	}

	selected := e.cards[index]
	if !Matches(e.reference, selected.Variant) {
		e.combo = 0
		return PickResult{Outcome: OutcomeIncorrect, State: e.Snapshot()}, nil
	}
	if e.score > math.MaxUint64-e.config.PointsPerCorrect || e.combo == math.MaxUint64 || e.boardVersion == math.MaxUint64 || e.referenceVersion == math.MaxUint64 {
		return PickResult{}, ErrScoreOverflow
	}

	newReference := selected.Variant
	replacement := e.newCard(Variant(e.rng.Intn(VariantCount)))
	candidateBoard := append([]CardInstance(nil), e.cards...)
	candidateBoard[index] = replacement
	if !hasMatch(newReference, candidateBoard) {
		candidateBoard[index] = e.newCard(e.randomMatchingVariant(newReference))
	}

	e.reference = newReference
	e.cards = candidateBoard
	e.score += e.config.PointsPerCorrect
	e.combo++
	e.boardVersion++
	e.referenceVersion++
	return PickResult{Outcome: OutcomeCorrect, State: e.Snapshot()}, nil
}

func (e *Engine) newCard(variant Variant) CardInstance {
	e.nextCardID++
	return CardInstance{ID: fmt.Sprintf("card-%016x", e.nextCardID), Variant: variant}
}

func (e *Engine) randomMatchingVariant(reference Variant) Variant {
	valid := [AttributeCount + 1]Variant{reference}
	for bit := 0; bit < AttributeCount; bit++ {
		valid[bit+1] = reference ^ Variant(1<<(AttributeCount-1-bit))
	}
	return valid[e.rng.Intn(len(valid))]
}

func hasMatch(reference Variant, cards []CardInstance) bool {
	for _, card := range cards {
		if Matches(reference, card.Variant) {
			return true
		}
	}
	return false
}
