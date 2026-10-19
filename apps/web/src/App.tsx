import { useCallback, useEffect, useState } from "react";
import { th } from "./i18n/th";
import { describeCard } from "./game/cocoCardVariant";
import {
  createLocalSpeedBattle,
  pickLocalCard,
  tickLocalSpeedBattle,
  type LocalSpeedBattleState,
} from "./game/localSpeedBattle";
import { LocalCardArt } from "./game/LocalCardArt";
import { PhaserBoard } from "./game/PhaserBoard";

type Screen = "home" | "game" | "results";

function App() {
  const [screen, setScreen] = useState<Screen>("home");
  const [game, setGame] = useState<LocalSpeedBattleState | null>(null);
  const [nowMs, setNowMs] = useState(() => Date.now());

  const startGame = useCallback(() => {
    const startTime = Date.now();
    setGame(createLocalSpeedBattle(createSeed(), startTime));
    setNowMs(startTime);
    setScreen("game");
  }, []);

  useEffect(() => {
    if (screen !== "game") return undefined;

    const timer = window.setInterval(() => {
      const now = Date.now();
      setNowMs(now);
      setGame((previous) =>
        previous ? tickLocalSpeedBattle(previous, now) : previous,
      );
    }, 100);
    return () => window.clearInterval(timer);
  }, [screen]);

  const selectCard = useCallback((cardId: string) => {
    const now = Date.now();
    setNowMs(now);
    setGame((previous) => {
      if (!previous) return previous;
      return pickLocalCard(previous, cardId, now).state;
    });
  }, []);

  if (screen === "game" && game?.status === "playing") {
    const referenceDescription = describeCard(game.reference);
    const cooldownActive = nowMs < game.cooldownUntilMs;
    return (
      <main className="app-shell game-shell" lang="th">
        <header className="game-topbar">
          <button
            className="text-button"
            type="button"
            onClick={() => setScreen("home")}
          >
            <span aria-hidden="true">←</span> {th.game.back}
          </button>
          <p className="game-mode-label">{th.home.mode}</p>
          <div
            className="time-pill"
            aria-label={`${th.game.timer} ${game.remainingSeconds} วินาที`}
          >
            <span className="time-berry" aria-hidden="true">
              ●
            </span>
            <span aria-hidden="true">{formatTime(game.remainingSeconds)}</span>
          </div>
        </header>

        <div className="game-layout">
          <section className="game-main" aria-label={th.game.board}>
            <div className="game-heading-row">
              <div>
                <p className="section-kicker">{th.picnic}</p>
                <h1>{th.game.instruction}</h1>
              </div>
              <div
                className="score-plate"
                aria-label={`${th.game.score} ${game.score}`}
              >
                <span>{th.game.score}</span>
                <strong>{game.score}</strong>
              </div>
            </div>

            <div className="playfield">
              <section
                className="reference-card"
                aria-labelledby="reference-title"
              >
                <div className="reference-copy">
                  <h2 id="reference-title">{th.game.reference}</h2>
                  <p>{referenceDescription}</p>
                </div>
                <LocalCardArt
                  className="reference-art"
                  variant={game.reference}
                />
              </section>

              <div className="board-and-feedback">
                <div className="board-heading">
                  <p>{th.game.board}</p>
                  <p
                    className="combo-indicator"
                    aria-label={`${th.game.combo} ${game.combo}`}
                  >
                    <span aria-hidden="true">✦</span> {game.combo}
                  </p>
                </div>
                <output className="visually-hidden" aria-live="polite">
                  {game.feedback === "correct"
                    ? th.game.correct
                    : game.feedback === "incorrect"
                      ? th.game.incorrect
                      : cooldownActive
                        ? th.game.cooldown
                        : ""}
                </output>
                <PhaserBoard
                  cards={game.cards}
                  reference={game.reference}
                  disabled={cooldownActive || game.status !== "playing"}
                  onSelect={selectCard}
                />
                <p className="game-hint">{th.game.instruction}</p>
              </div>
            </div>
          </section>

          <aside className="game-side-note" aria-label={th.game.combo}>
            <div className="combo-sticker" aria-hidden="true">
              ✦
            </div>
            <p>{th.game.combo}</p>
            <strong>{game.combo}</strong>
            <span>{th.times}</span>
          </aside>
        </div>
      </main>
    );
  }

  if ((screen === "results" || game?.status === "finished") && game) {
    return (
      <main className="app-shell result-shell" lang="th">
        <section className="result-card" aria-labelledby="result-title">
          <div className="result-hedgehog" aria-hidden="true">
            🦔
          </div>
          <p className="section-kicker">{th.picnic}</p>
          <h1 id="result-title">{th.result.title}</h1>
          <p className="result-score-label">{th.result.score}</p>
          <p className="result-score">{game.score}</p>
          <div className="result-actions">
            <button
              className="primary-button"
              type="button"
              onClick={startGame}
            >
              {th.result.replay}
            </button>
            <button
              className="secondary-button"
              type="button"
              onClick={() => setScreen("home")}
            >
              {th.result.home}
            </button>
          </div>
        </section>
      </main>
    );
  }

  return (
    <main className="app-shell home-shell" lang="th">
      <header className="home-header">
        <div className="brand-lockup">
          <span className="brand-mark" aria-hidden="true">
            ✿
          </span>
          <span>{th.picnic}</span>
        </div>
        <span className="home-duration">{th.home.duration}</span>
      </header>

      <section className="home-hero" aria-labelledby="home-title">
        <div className="hero-copy">
          <p className="section-kicker">{th.home.eyebrow}</p>
          <h1 id="home-title">{th.home.title}</h1>
          <p className="hero-description">{th.home.description}</p>
          <button className="primary-button" type="button" onClick={startGame}>
            {th.home.start}
            <span aria-hidden="true">↗</span>
          </button>
          <p className="mode-note">{th.home.mode}</p>
        </div>

        <div className="hero-portrait" aria-hidden="true">
          <span className="hero-sun" />
          <span className="hero-leaf hero-leaf-one">✦</span>
          <span className="hero-leaf hero-leaf-two">✿</span>
          <div className="hero-coco">
            <span className="hero-quills">⌁⌁⌁</span>
            <span className="hero-face">
              <span className="hero-eye hero-eye-left" />
              <span className="hero-eye hero-eye-right" />
              <span className="hero-nose" />
              <span className="hero-scarf" />
              <span className="hero-berry">●</span>
            </span>
          </div>
          <span className="portrait-caption">{th.coco}</span>
        </div>
      </section>

      <section className="home-instruction" aria-label={th.home.description}>
        <div className="instruction-art">
          <LocalCardArt variant={0b10110} />
        </div>
        <div>
          <p className="section-kicker">{th.tutorial.label}</p>
          <h2>{th.tutorial.title}</h2>
          <p>{th.home.description}</p>
        </div>
        <span className="instruction-spark" aria-hidden="true">
          ✦
        </span>
      </section>

      <footer className="home-footer">{th.brand} · เล่นในเครื่องนี้</footer>
    </main>
  );
}

function createSeed(): number {
  const values = new Uint32Array(1);
  if (typeof crypto !== "undefined" && "getRandomValues" in crypto) {
    crypto.getRandomValues(values);
    return values[0] ?? Date.now();
  }
  return Date.now() >>> 0;
}

function formatTime(seconds: number): string {
  const minutes = Math.floor(seconds / 60);
  const remainder = seconds % 60;
  return `${minutes}:${String(remainder).padStart(2, "0")}`;
}

export default App;
