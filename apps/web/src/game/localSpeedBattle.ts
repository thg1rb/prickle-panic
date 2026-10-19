import { BOARD_SIZE, VARIANT_COUNT, isValidPick } from "./cocoCardVariant";

export type LocalCard = {
  id: string;
  variant: number;
};

export type LocalGameStatus = "playing" | "finished";
export type PickOutcome = "correct" | "incorrect" | "cooldown" | "inactive";

export type LocalSpeedBattleConfig = {
  boardSize: number;
  durationSeconds: number;
  pointsPerCorrect: number;
  cooldownMs: number;
};

export type LocalSpeedBattleState = {
  status: LocalGameStatus;
  reference: number;
  cards: LocalCard[];
  score: number;
  combo: number;
  remainingSeconds: number;
  deadlineMs: number;
  cooldownUntilMs: number;
  feedback: "" | "correct" | "incorrect";
  rngState: number;
  nextCardId: number;
  config: LocalSpeedBattleConfig;
};

export type PickTransition = {
  state: LocalSpeedBattleState;
  outcome: PickOutcome;
};

export const DEFAULT_LOCAL_CONFIG: Readonly<LocalSpeedBattleConfig> = {
  boardSize: BOARD_SIZE,
  durationSeconds: 60,
  pointsPerCorrect: 1,
  cooldownMs: 800,
};

export function createLocalSpeedBattle(
  seed: number,
  nowMs: number,
  config: LocalSpeedBattleConfig = { ...DEFAULT_LOCAL_CONFIG },
): LocalSpeedBattleState {
  if (
    !Number.isInteger(config.boardSize) ||
    config.boardSize < 1 ||
    !Number.isInteger(config.durationSeconds) ||
    config.durationSeconds < 1 ||
    !Number.isInteger(config.pointsPerCorrect) ||
    config.pointsPerCorrect < 1 ||
    !Number.isInteger(config.cooldownMs) ||
    config.cooldownMs < 0
  ) {
    throw new RangeError("Local Speed Battle configuration is invalid.");
  }

  let rngState = seed >>> 0;
  if (rngState === 0) rngState = 0x6d2b79f5;

  const draw = (): number => {
    rngState ^= rngState << 13;
    rngState ^= rngState >>> 17;
    rngState ^= rngState << 5;
    return (rngState >>> 0) / 0x1_0000_0000;
  };
  let nextCardId = 0;
  const createCard = (variant: number): LocalCard => {
    nextCardId += 1;
    return { id: `card-${nextCardId}`, variant };
  };
  const reference = Math.floor(draw() * VARIANT_COUNT);
  const cards = Array.from({ length: config.boardSize }, () =>
    createCard(Math.floor(draw() * VARIANT_COUNT)),
  );

  if (!cards.some((card) => isValidPick(reference, card.variant))) {
    cards[cards.length - 1] = createCard(
      randomMatchingVariant(reference, draw),
    );
  }

  return {
    status: "playing",
    reference,
    cards,
    score: 0,
    combo: 0,
    remainingSeconds: config.durationSeconds,
    deadlineMs: nowMs + config.durationSeconds * 1000,
    cooldownUntilMs: 0,
    feedback: "",
    rngState,
    nextCardId,
    config: { ...config },
  };
}

export function tickLocalSpeedBattle(
  state: LocalSpeedBattleState,
  nowMs: number,
): LocalSpeedBattleState {
  if (state.status !== "playing") return state;

  const remainingSeconds = Math.max(
    0,
    Math.ceil((state.deadlineMs - nowMs) / 1000),
  );
  const cooldownActive = nowMs < state.cooldownUntilMs;
  const cooldownUntilMs = cooldownActive ? state.cooldownUntilMs : 0;
  const feedback =
    cooldownActive || state.feedback === "correct" ? state.feedback : "";
  const status = remainingSeconds === 0 ? "finished" : "playing";

  if (
    remainingSeconds === state.remainingSeconds &&
    cooldownUntilMs === state.cooldownUntilMs &&
    feedback === state.feedback &&
    status === state.status
  ) {
    return state;
  }

  return { ...state, remainingSeconds, cooldownUntilMs, feedback, status };
}

export function pickLocalCard(
  currentState: LocalSpeedBattleState,
  cardId: string,
  nowMs: number,
): PickTransition {
  const state = tickLocalSpeedBattle(currentState, nowMs);
  if (state.status !== "playing") return { state, outcome: "inactive" };
  if (nowMs < state.cooldownUntilMs) return { state, outcome: "cooldown" };

  const index = state.cards.findIndex((card) => card.id === cardId);
  if (index < 0) return { state, outcome: "inactive" };

  const selected = state.cards[index];
  if (!isValidPick(state.reference, selected.variant)) {
    return {
      outcome: "incorrect",
      state: {
        ...state,
        combo: 0,
        cooldownUntilMs: nowMs + state.config.cooldownMs,
        feedback: "incorrect",
      },
    };
  }

  const nextRandom = createRandom(state.rngState);
  const newReference = selected.variant;
  const cards = state.cards.slice();
  cards[index] = makeCard(
    state.nextCardId + 1,
    Math.floor(nextRandom.draw() * VARIANT_COUNT),
  );
  let nextCardId = state.nextCardId + 1;

  if (!cards.some((card) => isValidPick(newReference, card.variant))) {
    nextCardId += 1;
    cards[index] = makeCard(
      nextCardId,
      randomMatchingVariant(newReference, nextRandom.draw),
    );
  }

  return {
    outcome: "correct",
    state: {
      ...state,
      reference: newReference,
      cards,
      score: state.score + state.config.pointsPerCorrect,
      combo: state.combo + 1,
      remainingSeconds: state.remainingSeconds,
      cooldownUntilMs: 0,
      feedback: "correct",
      rngState: nextRandom.state(),
      nextCardId,
    },
  };
}

function makeCard(id: number, variant: number): LocalCard {
  return { id: `card-${id}`, variant };
}

function randomMatchingVariant(reference: number, draw: () => number): number {
  const valid = [reference];
  for (let bit = 0; bit < 5; bit += 1) valid.push(reference ^ (1 << bit));
  return valid[Math.floor(draw() * valid.length)];
}

function createRandom(seed: number) {
  let rngState = seed >>> 0;
  if (rngState === 0) rngState = 0x6d2b79f5;

  return {
    draw: () => {
      rngState ^= rngState << 13;
      rngState ^= rngState >>> 17;
      rngState ^= rngState << 5;
      return (rngState >>> 0) / 0x1_0000_0000;
    },
    state: () => rngState,
  };
}
