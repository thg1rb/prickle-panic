import {
  act,
  cleanup,
  fireEvent,
  render,
  screen,
} from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import App from "./App";
import type { LocalCard } from "./game/localSpeedBattle";
import { isValidPick } from "./game/cocoCardVariant";

vi.mock("./game/PhaserBoard", () => ({
  PhaserBoard: ({
    cards,
    reference,
    disabled,
    onSelect,
  }: {
    cards: readonly LocalCard[];
    reference: number;
    disabled: boolean;
    onSelect: (cardId: string) => void;
  }) => (
    <fieldset aria-label="กระดานการ์ดโคโค่">
      {cards.map((card, index) => (
        <button
          data-invalid={!isValidPick(reference, card.variant)}
          key={card.id}
          type="button"
          disabled={disabled}
          onClick={() => onSelect(card.id)}
        >
          การ์ดทดสอบ {index + 1}
        </button>
      ))}
    </fieldset>
  ),
}));

afterEach(() => {
  cleanup();
  vi.useRealTimers();
});

describe("Thai local gameplay flow", () => {
  it("opens a 16-card training board from the Thai home screen", () => {
    render(<App />);
    expect(
      screen.getByRole("heading", { name: "พร้อมลุยหรือยัง?" }),
    ).toBeTruthy();
    fireEvent.click(screen.getByRole("button", { name: "เริ่มเล่นคนเดียว" }));

    expect(
      screen.getByRole("heading", { name: /เลือกการ์ดที่ต่าง/ }),
    ).toBeTruthy();
    expect(screen.getAllByRole("button", { name: /การ์ดทดสอบ/ })).toHaveLength(
      16,
    );
    expect(screen.getByLabelText(/เวลาที่เหลือ 60 วินาที/)).toBeTruthy();
  });

  it("shows incorrect-pick feedback and blocks the board during cooldown", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date(1000));
    render(<App />);
    fireEvent.click(screen.getByRole("button", { name: "เริ่มเล่นคนเดียว" }));
    const button = document.querySelector<HTMLButtonElement>(
      'button[data-invalid="true"]',
    )!;
    fireEvent.click(button);

    expect(document.querySelector("output")?.textContent).toContain("ยังไม่ใช");
    expect(
      screen
        .getAllByRole("button", { name: /การ์ดทดสอบ/ })
        .every((card) => (card as HTMLButtonElement).disabled),
    ).toBe(true);
  });

  it("shows results when the round ends and starts a fresh round on replay", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date(0));
    render(<App />);
    fireEvent.click(screen.getByRole("button", { name: "เริ่มเล่นคนเดียว" }));
    act(() => vi.advanceTimersByTime(60000));

    expect(screen.getByRole("heading", { name: "หมดเวลาแล้ว!" })).toBeTruthy();
    fireEvent.click(screen.getByRole("button", { name: "เล่นอีกครั้ง" }));
    expect(
      screen.getByRole("heading", { name: /เลือกการ์ดที่ต่าง/ }),
    ).toBeTruthy();
    expect(screen.getAllByRole("button", { name: /การ์ดทดสอบ/ })).toHaveLength(
      16,
    );
  });
});
