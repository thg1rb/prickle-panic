package game

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestConfigValidation(t *testing.T) {
	for _, config := range []Config{
		{BoardSize: 0, PointsPerCorrect: 1},
		{BoardSize: 16, PointsPerCorrect: 0},
	} {
		if _, err := NewEngine(config, 1); !errors.Is(err, ErrInvalidConfig) {
			t.Errorf("NewEngine(%+v) error = %v, want ErrInvalidConfig", config, err)
		}
	}
}

func TestSeededBoardsAreReproducibleSolvableAndHaveUniqueInstances(t *testing.T) {
	config := DefaultConfig()
	for seed := int64(0); seed < 200; seed++ {
		first, err := NewEngine(config, seed)
		if err != nil {
			t.Fatal(err)
		}
		second, err := NewEngine(config, seed)
		if err != nil {
			t.Fatal(err)
		}
		left, right := first.Snapshot(), second.Snapshot()
		if !reflect.DeepEqual(left, right) {
			t.Fatalf("same seed %d generated different states", seed)
		}
		if len(left.Cards) != config.BoardSize {
			t.Fatalf("seed %d board has %d cards, want %d", seed, len(left.Cards), config.BoardSize)
		}
		if !hasMatch(left.Reference, left.Cards) {
			t.Fatalf("seed %d board has no valid choice", seed)
		}
		ids := make(map[string]bool, len(left.Cards))
		for _, card := range left.Cards {
			if card.ID == "" || ids[card.ID] {
				t.Fatalf("seed %d has empty or duplicate card ID %q", seed, card.ID)
			}
			ids[card.ID] = true
		}
	}
}

func TestInitialBoardFallbackMakesUnsuitableSampleSolvable(t *testing.T) {
	config := DefaultConfig()
	var fallbackSeed int64 = -1
	for seed := int64(0); seed < 1000; seed++ {
		rng := newSeededRandom(seed)
		reference := Variant(rng.Intn(VariantCount))
		cards := make([]CardInstance, config.BoardSize)
		for index := range cards {
			cards[index] = CardInstance{Variant: Variant(rng.Intn(VariantCount))}
		}
		if !hasMatch(reference, cards) {
			fallbackSeed = seed
			break
		}
	}
	if fallbackSeed < 0 {
		t.Fatal("test could not find a seed requiring the solvability fallback")
	}
	engine, err := NewEngine(config, fallbackSeed)
	if err != nil {
		t.Fatal(err)
	}
	if !hasMatch(engine.reference, engine.cards) {
		t.Fatalf("fallback did not make seed %d solvable", fallbackSeed)
	}
}

func TestCorrectPickUpdatesScoreReferenceComboAndReplacement(t *testing.T) {
	variants := make([]Variant, 16)
	variants[0] = 1
	for index := 1; index < len(variants); index++ {
		variants[index] = 6
	}
	engine := newTestEngine(0, variants)
	result, err := engine.Pick("card-0000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != OutcomeCorrect {
		t.Fatalf("outcome = %q, want correct", result.Outcome)
	}
	if result.State.Score != 1 || result.State.Combo != 1 {
		t.Fatalf("score/combo = %d/%d, want 1/1", result.State.Score, result.State.Combo)
	}
	if result.State.Reference != 1 || result.State.ReferenceVersion != 1 || result.State.BoardVersion != 1 {
		t.Fatalf("reference and versions not advanced: %+v", result.State)
	}
	if len(result.State.Cards) != 16 || !hasMatch(result.State.Reference, result.State.Cards) {
		t.Fatalf("replacement broke board cardinality or solvability: %+v", result.State)
	}
	ids := make(map[string]bool, len(result.State.Cards))
	for _, card := range result.State.Cards {
		if card.ID == "card-0000000000000001" {
			t.Fatal("selected card instance remained on the board")
		}
		if ids[card.ID] {
			t.Fatalf("replacement duplicated active card ID %q", card.ID)
		}
		ids[card.ID] = true
	}
}

func TestIncorrectPickResetsComboWithoutChangingOtherGameplayState(t *testing.T) {
	engine := newTestEngine(0, []Variant{1, 6})
	engine.score = 9
	engine.combo = 4
	engine.boardVersion = 7
	engine.referenceVersion = 3
	before := engine.Snapshot()
	result, err := engine.Pick("card-0000000000000002")
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != OutcomeIncorrect {
		t.Fatalf("outcome = %q, want incorrect", result.Outcome)
	}
	if result.State.Score != before.Score || result.State.Reference != before.Reference || result.State.BoardVersion != before.BoardVersion || result.State.ReferenceVersion != before.ReferenceVersion {
		t.Fatalf("incorrect pick changed score/reference/versions: before=%+v after=%+v", before, result.State)
	}
	if result.State.Combo != 0 || !reflect.DeepEqual(result.State.Cards, before.Cards) {
		t.Fatalf("incorrect pick did not only reset combo: before=%+v after=%+v", before, result.State)
	}
}

func TestCorrectReplacementFallbackRestoresSolvability(t *testing.T) {
	variants := make([]Variant, 16)
	variants[0] = 1
	for index := 1; index < len(variants); index++ {
		variants[index] = 6
	}
	engine := newTestEngine(0, variants)
	for seed := int64(0); seed < 100; seed++ {
		if !Matches(1, Variant(newSeededRandom(seed).Intn(VariantCount))) {
			engine.rng = newSeededRandom(seed)
			break
		}
	}
	if !hasMatch(engine.reference, engine.cards) {
		t.Fatal("test board must initially be solvable")
	}
	result, err := engine.Pick("card-0000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	if !hasMatch(result.State.Reference, result.State.Cards) {
		t.Fatal("correct replacement did not restore a valid selection")
	}
	if got := result.State.Cards[0].ID; got != "card-0000000000000012" {
		t.Fatalf("fallback card ID = %q, want card-...0018", got)
	}
}

func TestTechnicalErrorIsAtomicAndDoesNotConsumeRandomState(t *testing.T) {
	left := newTestEngine(0, []Variant{1, 6})
	right := newTestEngine(0, []Variant{1, 6})
	before := left.Snapshot()
	if _, err := left.Pick("missing-card"); !errors.Is(err, ErrUnknownCard) {
		t.Fatalf("Pick error = %v, want ErrUnknownCard", err)
	}
	if got := left.Snapshot(); !reflect.DeepEqual(got, before) {
		t.Fatalf("technical error mutated state: before=%+v after=%+v", before, got)
	}
	leftResult, leftErr := left.Pick("card-0000000000000001")
	rightResult, rightErr := right.Pick("card-0000000000000001")
	if leftErr != nil || rightErr != nil || !reflect.DeepEqual(leftResult, rightResult) {
		t.Fatalf("failed action affected future deterministic output: left=%+v/%v right=%+v/%v", leftResult, leftErr, rightResult, rightErr)
	}
}

func TestSnapshotDoesNotExposeMutableBoardStorage(t *testing.T) {
	engine := newTestEngine(0, []Variant{1, 6})
	snapshot := engine.Snapshot()
	snapshot.Cards[0].ID = "mutated"
	snapshot.Cards[1] = CardInstance{}
	if engine.cards[0].ID == "mutated" || engine.cards[1].ID == "" {
		t.Fatal("snapshot mutation changed engine-owned state")
	}
}

func TestDuplicateVariationsRemainDistinctInstances(t *testing.T) {
	engine := newTestEngine(0, []Variant{1, 1})
	if engine.cards[0].Variant != engine.cards[1].Variant {
		t.Fatal("test setup requires duplicate variations")
	}
	if engine.cards[0].ID == engine.cards[1].ID {
		t.Fatal("duplicate variations must have distinct card instance IDs")
	}
}

func newTestEngine(reference Variant, variants []Variant) *Engine {
	engine := &Engine{
		config:    Config{BoardSize: len(variants), PointsPerCorrect: 1},
		rng:       newSeededRandom(10),
		reference: reference,
		cards:     make([]CardInstance, len(variants)),
	}
	for index, variant := range variants {
		engine.cards[index] = CardInstance{ID: fmt.Sprintf("card-%016x", index+1), Variant: variant}
	}
	engine.nextCardID = uint64(len(variants))
	return engine
}
