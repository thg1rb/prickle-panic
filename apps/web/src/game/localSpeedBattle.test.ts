import { describe, expect, it } from "vitest";
import {
  DEFAULT_LOCAL_CONFIG,
  createLocalSpeedBattle,
  pickLocalCard,
  tickLocalSpeedBattle,
} from "./localSpeedBattle";
import { isValidPick } from "./cocoCardVariant";

describe("local Speed Battle adapter", () => {
  it("creates a reproducible, playable board with configurable defaults", () => {
    const first = createLocalSpeedBattle(1337, 5000);
    const second = createLocalSpeedBattle(1337, 5000);
    expect(first).toEqual(second);
    expect(first.cards).toHaveLength(DEFAULT_LOCAL_CONFIG.boardSize);
    expect(
      first.cards.some((card) => isValidPick(first.reference, card.variant)),
    ).toBe(true);
    expect(first.remainingSeconds).toBe(DEFAULT_LOCAL_CONFIG.durationSeconds);
    expect(first.deadlineMs).toBe(65000);
  });

  it("scores a valid selection, advances the reference, and keeps a solvable board", () => {
    const initial = createLocalSpeedBattle(41, 1000);
    const selected = initial.cards.find((card) =>
      isValidPick(initial.reference, card.variant),
    );
    expect(selected).toBeDefined();

    const result = pickLocalCard(initial, selected!.id, 1200);
    expect(result.outcome).toBe("correct");
    expect(result.state.score).toBe(1);
    expect(result.state.combo).toBe(1);
    expect(result.state.reference).toBe(selected!.variant);
    expect(result.state.cards).toHaveLength(initial.cards.length);
    expect(
      result.state.cards.some((card) =>
        isValidPick(result.state.reference, card.variant),
      ),
    ).toBe(true);
    expect(
      result.state.cards.find((card) => card.id === selected!.id),
    ).toBeUndefined();
    expect(tickLocalSpeedBattle(result.state, 1300).feedback).toBe("correct");
  });

  it("applies an incorrect-pick cooldown, clears combo, and does not extend it", () => {
    const initial = createLocalSpeedBattle(8, 0);
    const primed = { ...initial, combo: 4 };
    const invalid = primed.cards.find(
      (card) => !isValidPick(primed.reference, card.variant),
    );
    expect(invalid).toBeDefined();

    const wrong = pickLocalCard(primed, invalid!.id, 100);
    expect(wrong.outcome).toBe("incorrect");
    expect(wrong.state.score).toBe(0);
    expect(wrong.state.combo).toBe(0);
    expect(wrong.state.cooldownUntilMs).toBe(900);
    expect(wrong.state.reference).toBe(primed.reference);
    expect(wrong.state.cards).toEqual(primed.cards);

    const blocked = pickLocalCard(wrong.state, invalid!.id, 500);
    expect(blocked.outcome).toBe("cooldown");
    expect(blocked.state.cooldownUntilMs).toBe(900);
    expect(tickLocalSpeedBattle(blocked.state, 900).cooldownUntilMs).toBe(0);
  });

  it("rejects unknown cards and finishes at the exact deadline", () => {
    const initial = createLocalSpeedBattle(2026, 1000);
    expect(pickLocalCard(initial, "missing", 1100).outcome).toBe("inactive");
    const finished = tickLocalSpeedBattle(initial, initial.deadlineMs);
    expect(finished.status).toBe("finished");
    expect(finished.remainingSeconds).toBe(0);
    expect(
      pickLocalCard(finished, initial.cards[0]!.id, finished.deadlineMs)
        .outcome,
    ).toBe("inactive");
  });

  it("rejects invalid tuning values", () => {
    expect(() =>
      createLocalSpeedBattle(1, 0, { ...DEFAULT_LOCAL_CONFIG, boardSize: 0 }),
    ).toThrow(RangeError);
  });
});
