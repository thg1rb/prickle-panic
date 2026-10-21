import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import { GuestSessionGate } from "./GuestSessionGate";
import type { GuestSessionErrorCode } from "./client";

afterEach(cleanup);

const errorMessages: [GuestSessionErrorCode, string][] = [
  ["USERNAME_INVALID", "ชื่อนี้ยังใช้ไม่ได้"],
  ["REQUEST_INVALID", "ข้อมูลชื่อไม่ถูกต้อง"],
  ["REQUEST_TOO_LARGE", "ข้อมูลชื่อยาวเกินไป"],
  ["ORIGIN_NOT_ALLOWED", "ไม่สามารถส่งคำขอจากหน้านี้ได้"],
  ["RATE_LIMITED", "มีการเริ่มเซสชันหลายครั้งเกินไป"],
  ["SESSION_REQUIRED", "เซสชันหมดอายุแล้ว"],
  ["SESSION_UNAVAILABLE", "บริการยังไม่พร้อม"],
];

describe("Thai Guest Session entry", () => {
  it("shows a labeled username field and submits on Enter", async () => {
    const user = userEvent.setup();
    const onCreate = vi.fn().mockResolvedValue(undefined);
    render(
      <GuestSessionGate
        state={{ status: "UNAUTHENTICATED", prefill: "" }}
        onCreate={onCreate}
        onRetryRecovery={vi.fn()}
      />,
    );

    expect(
      screen.getByRole("heading", { name: "ยินดีต้อนรับสู่ปิกนิกของโคโค่" }),
    ).toBeTruthy();
    const input = screen.getByRole("textbox", { name: "ชื่อผู้เล่น" });
    await user.type(input, "โคโค่");
    await user.keyboard("{Enter}");

    expect(onCreate).toHaveBeenCalledTimes(1);
    expect(onCreate).toHaveBeenCalledWith("โคโค่");
  });

  it("announces client-side length feedback and disables invalid submission", async () => {
    const onCreate = vi.fn();
    render(
      <GuestSessionGate
        state={{ status: "UNAUTHENTICATED", prefill: "" }}
        onCreate={onCreate}
        onRetryRecovery={vi.fn()}
      />,
    );
    fireEvent.change(screen.getByLabelText("ชื่อผู้เล่น"), {
      target: { value: "ก" },
    });

    expect(screen.getByRole("alert").textContent).toContain("2–16");
    expect(screen.getByRole("button", { name: "เริ่มปิกนิก" })).toHaveProperty(
      "disabled",
      true,
    );
    expect(onCreate).not.toHaveBeenCalled();
  });

  it("shows busy state during creation and Thai retry guidance for service failure", () => {
    const onCreate = vi.fn();
    const { rerender } = render(
      <GuestSessionGate
        state={{ status: "CREATING_SESSION", username: "Coco" }}
        onCreate={onCreate}
        onRetryRecovery={vi.fn()}
      />,
    );
    expect(
      screen.getByRole("button", { name: "กำลังเตรียมเซสชัน..." }),
    ).toHaveProperty("disabled", true);
    expect(
      screen.getByText("กำลังเตรียมเซสชัน...", { selector: "output" }),
    ).toBeTruthy();

    rerender(
      <GuestSessionGate
        state={{ status: "RECOVERY_ERROR", code: "SESSION_UNAVAILABLE" }}
        onCreate={onCreate}
        onRetryRecovery={vi.fn()}
      />,
    );
    expect(screen.getByRole("alert").textContent).toContain(
      "บริการยังไม่พร้อม",
    );
    expect(
      screen.getByRole("button", { name: "ลองตรวจสอบอีกครั้ง" }),
    ).toBeTruthy();
  });

  it.each(errorMessages)(
    "maps %s to player-safe Thai feedback",
    (code, message) => {
      render(
        <GuestSessionGate
          state={{
            status: "UNAUTHENTICATED",
            prefill: "Coco",
            error: { code },
          }}
          onCreate={vi.fn()}
          onRetryRecovery={vi.fn()}
        />,
      );

      expect(screen.getByRole("alert").textContent).toContain(message);
    },
  );
});
