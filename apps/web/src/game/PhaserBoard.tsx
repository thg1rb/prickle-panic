import { useEffect, useRef } from "react";
import type { LocalCard } from "./localSpeedBattle";
import { describeCardTh, th } from "../i18n/th";

type PhaserBoardProps = {
  cards: readonly LocalCard[];
  reference: number;
  disabled: boolean;
  onSelect: (cardId: string) => void;
};

type BoardInput = Pick<PhaserBoardProps, "cards" | "disabled" | "onSelect">;

export function PhaserBoard({
  cards,
  reference,
  disabled,
  onSelect,
}: PhaserBoardProps) {
  const hostRef = useRef<HTMLDivElement>(null);
  const redrawRef = useRef<() => void>(() => undefined);
  const latestRef = useRef<BoardInput>({ cards, disabled, onSelect });
  latestRef.current = { cards, disabled, onSelect };

  useEffect(() => {
    let disposed = false;
    let game: import("phaser").Game | undefined;
    redrawRef.current = () => undefined;

    void import("phaser").then(({ default: Phaser }) => {
      if (disposed || !hostRef.current) return;

      class CocoBoardScene extends Phaser.Scene {
        private graphics?: Phaser.GameObjects.Graphics;
        private zones: Phaser.GameObjects.Zone[] = [];

        constructor() {
          super({ key: "CocoBoardScene" });
        }

        create() {
          this.graphics = this.add.graphics();
          this.scale.on(Phaser.Scale.Events.RESIZE, this.draw, this);
          redrawRef.current = this.draw.bind(this);
          this.draw();
        }

        shutdown() {
          this.scale.off(Phaser.Scale.Events.RESIZE, this.draw, this);
          for (const zone of this.zones) zone.destroy();
          this.zones = [];
          redrawRef.current = () => undefined;
        }

        private draw() {
          if (!this.graphics) return;
          const input = latestRef.current;
          const graphics = this.graphics;
          graphics.clear();
          for (const zone of this.zones) zone.destroy();
          this.zones = [];
          this.input.enabled = !input.disabled;

          const bounds = 640;
          const margin = 24;
          const gap = 10;
          const cardSize = (bounds - margin * 2 - gap * 3) / 4;

          input.cards.forEach((card, index) => {
            const column = index % 4;
            const row = Math.floor(index / 4);
            const x = margin + column * (cardSize + gap);
            const y = margin + row * (cardSize + gap);
            drawCard(graphics, x, y, cardSize, card.variant);

            const zone = this.add.zone(
              x + cardSize / 2,
              y + cardSize / 2,
              cardSize,
              cardSize,
            );
            zone.setInteractive();
            zone.on("pointerup", () => {
              const current = latestRef.current;
              if (!current.disabled) current.onSelect(card.id);
            });
            this.zones.push(zone);
          });
        }
      }

      const scene = new CocoBoardScene();
      game = new Phaser.Game({
        type: Phaser.AUTO,
        parent: hostRef.current,
        width: 640,
        height: 640,
        transparent: true,
        scale: {
          mode: Phaser.Scale.FIT,
          autoCenter: Phaser.Scale.CENTER_BOTH,
        },
        scene,
        banner: false,
        audio: { noAudio: true },
      });
      game.events.once(Phaser.Core.Events.READY, () => {
        game?.canvas?.setAttribute("aria-hidden", "true");
        if (game?.canvas) game.canvas.tabIndex = -1;
      });
    });

    return () => {
      disposed = true;
      redrawRef.current = () => undefined;
      game?.destroy(true);
    };
  }, []);

  useEffect(() => {
    redrawRef.current();
  }, [cards, disabled]);

  return (
    <div className="board-shell" data-reference={reference}>
      <div className="phaser-board-host" ref={hostRef} aria-hidden="true" />
      <fieldset
        className="board-accessible-controls"
        aria-label={th.game.board}
      >
        {cards.map((card, index) => (
          <button
            className="board-card-control"
            type="button"
            key={card.id}
            data-variant={card.variant}
            aria-label={`${th.game.cardPrefix} ${index + 1} ${describeCardTh(card.variant)}`}
            disabled={disabled}
            onClick={() => onSelect(card.id)}
          />
        ))}
      </fieldset>
    </div>
  );
}

function drawCard(
  graphics: import("phaser").GameObjects.Graphics,
  x: number,
  y: number,
  size: number,
  variant: number,
) {
  const bodyCream = (variant & 0b10000) !== 0;
  const longQuills = (variant & 0b01000) !== 0;
  const glasses = (variant & 0b00100) !== 0;
  const scarf = (variant & 0b00010) !== 0;
  const strawberry = (variant & 0b00001) !== 0;
  const centerX = x + size / 2;
  const centerY = y + size / 2;

  graphics
    .fillStyle(0x654736, 0.13)
    .fillRoundedRect(x + 2, y + 5, size - 2, size - 2, 18);
  graphics.fillStyle(0xfffdf7, 1).lineStyle(2, 0xe7d6bd, 1);
  graphics.fillRoundedRect(x, y, size - 3, size - 4, 18);
  graphics.strokeRoundedRect(x, y, size - 3, size - 4, 18);

  const bodyColor = bodyCream ? 0xf1d9b8 : 0x79543d;
  const bodyShade = bodyCream ? 0xd5b28b : 0x553b2d;
  const faceY = centerY + 3;
  const faceRadius = size * 0.25;
  const quillTop = faceY - faceRadius - (longQuills ? 21 : 9);
  graphics.fillStyle(bodyShade, 1);
  for (let quill = 0; quill < 7; quill += 1) {
    const angle = Math.PI * (1.08 + quill * 0.14);
    const rootX = centerX + Math.cos(angle) * faceRadius * 0.86;
    const rootY = faceY + Math.sin(angle) * faceRadius * 0.86;
    const tipX = centerX + (quill - 3) * (longQuills ? 7 : 5);
    graphics.fillTriangle(
      rootX - 7,
      rootY + 5,
      rootX + 7,
      rootY + 4,
      tipX,
      quillTop + Math.abs(quill - 3) * 2,
    );
  }

  graphics
    .fillStyle(bodyColor, 1)
    .fillEllipse(centerX, faceY, faceRadius * 2.25, faceRadius * 1.8);
  graphics
    .lineStyle(2, bodyShade, 1)
    .strokeEllipse(centerX, faceY, faceRadius * 2.25, faceRadius * 1.8);

  const eyeY = faceY - 2;
  const eyeOffset = faceRadius * 0.38;
  if (glasses) {
    graphics.lineStyle(2.5, 0x42382f, 1);
    graphics.strokeCircle(centerX - eyeOffset, eyeY, 8);
    graphics.strokeCircle(centerX + eyeOffset, eyeY, 8);
    graphics.lineBetween(
      centerX - eyeOffset + 8,
      eyeY,
      centerX + eyeOffset - 8,
      eyeY,
    );
  } else {
    graphics.fillStyle(0x302820, 1);
    graphics.fillCircle(centerX - eyeOffset, eyeY, 2.8);
    graphics.fillCircle(centerX + eyeOffset, eyeY, 2.8);
  }
  graphics.lineStyle(2, 0x654736, 1).beginPath();
  graphics.arc(centerX, faceY + 9, 7, 0.2, Math.PI - 0.2, false).strokePath();

  if (scarf) {
    graphics.fillStyle(0xe57570, 1);
    graphics.fillRoundedRect(
      centerX - faceRadius * 0.95,
      faceY + faceRadius * 0.48,
      faceRadius * 1.9,
      8,
      4,
    );
    graphics.fillTriangle(
      centerX + 13,
      faceY + faceRadius * 0.57,
      centerX + 23,
      faceY + faceRadius * 0.57,
      centerX + 19,
      faceY + faceRadius * 1.05,
    );
  }
  if (strawberry) {
    const berryX = x + size * 0.78;
    const berryY = y + size * 0.74;
    graphics
      .fillStyle(0x6d965f, 1)
      .fillTriangle(
        berryX - 8,
        berryY - 8,
        berryX + 7,
        berryY - 8,
        berryX,
        berryY - 15,
      );
    graphics.fillStyle(0xe96962, 1).fillCircle(berryX, berryY, 10);
    graphics.fillStyle(0xfff4d9, 1).fillCircle(berryX - 3, berryY - 1, 1.4);
    graphics.fillCircle(berryX + 3, berryY + 2, 1.4);
  }
}
