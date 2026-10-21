import { StrictMode } from "react";
import { act, cleanup, renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { displayNameStorageKey, guestSessionPath } from "./client";
import { useGuestSession, validateUsernameLength } from "./useGuestSession";

const guest = {
  guestId: "gst_same_guest",
  username: "โคโค่ Berry",
  expiresAt: "2026-10-22T10:00:00Z",
};

function jsonResponse(value: unknown, status = 200, headers?: HeadersInit) {
  return Response.json(value, { status, headers });
}

afterEach(() => {
  cleanup();
  localStorage.clear();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe("Guest Session bootstrap and recovery", () => {
  it("creates only after explicit submission, stores only the canonical name, then recovers the same identity", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        jsonResponse({ error: { code: "SESSION_REQUIRED" } }, 401),
      )
      .mockResolvedValueOnce(
        jsonResponse({ error: { code: "SESSION_REQUIRED" } }, 401),
      )
      .mockResolvedValueOnce(
        jsonResponse({ ...guest, username: "โคโค่ Berry" }, 201),
      );
    vi.stubGlobal("fetch", fetchMock);

    const first = renderHook(() => useGuestSession(), {
      wrapper: StrictMode,
    });
    await waitFor(() =>
      expect(first.result.current.state.status).toBe("UNAUTHENTICATED"),
    );
    expect(
      fetchMock.mock.calls.every(([url]) => url === guestSessionPath),
    ).toBe(true);
    expect(
      fetchMock.mock.calls.filter(
        ([, init]) => (init as RequestInit).method === "POST",
      ),
    ).toHaveLength(0);

    await act(async () => {
      await first.result.current.create("โคโค่");
    });
    await waitFor(() =>
      expect(first.result.current.state.status).toBe("AUTHENTICATED"),
    );
    expect(first.result.current.state).toMatchObject({
      status: "AUTHENTICATED",
      session: { guestId: "gst_same_guest", username: "โคโค่ Berry" },
    });
    expect(localStorage.getItem(displayNameStorageKey)).toBe("โคโค่ Berry");
    expect(localStorage.length).toBe(1);
    expect(localStorage.key(0)).toBe(displayNameStorageKey);

    first.unmount();
    const recoveryFetch = vi.fn().mockResolvedValueOnce(jsonResponse(guest));
    vi.stubGlobal("fetch", recoveryFetch);
    const second = renderHook(() => useGuestSession());
    await waitFor(() =>
      expect(second.result.current.state.status).toBe("AUTHENTICATED"),
    );
    expect(second.result.current.state).toMatchObject({
      status: "AUTHENTICATED",
      session: { guestId: "gst_same_guest", username: "โคโค่ Berry" },
    });
    expect(recoveryFetch).toHaveBeenCalledWith(
      guestSessionPath,
      expect.objectContaining({ method: "GET", credentials: "include" }),
    );
  });

  it("keeps recovery failures distinct from an unauthenticated session and retries", async () => {
    const fetchMock = vi
      .fn()
      .mockRejectedValueOnce(new TypeError("offline"))
      .mockResolvedValueOnce(jsonResponse(guest));
    vi.stubGlobal("fetch", fetchMock);
    const { result } = renderHook(() => useGuestSession());

    await waitFor(() =>
      expect(result.current.state).toEqual({
        status: "RECOVERY_ERROR",
        code: "SESSION_UNAVAILABLE",
      }),
    );
    act(() => result.current.retryRecovery());
    await waitFor(() =>
      expect(result.current.state.status).toBe("AUTHENTICATED"),
    );
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("does not treat a saved display name as authorization", async () => {
    localStorage.setItem(displayNameStorageKey, "Coco");
    vi.stubGlobal(
      "fetch",
      vi
        .fn()
        .mockResolvedValueOnce(
          jsonResponse({ error: { code: "SESSION_REQUIRED" } }, 401),
        ),
    );
    const { result } = renderHook(() => useGuestSession());

    await waitFor(() =>
      expect(result.current.state).toEqual({
        status: "UNAUTHENTICATED",
        prefill: "Coco",
      }),
    );
  });

  it("tolerates blocked localStorage reads", async () => {
    vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
      throw new DOMException("blocked", "SecurityError");
    });
    vi.stubGlobal(
      "fetch",
      vi
        .fn()
        .mockResolvedValueOnce(
          jsonResponse({ error: { code: "SESSION_REQUIRED" } }, 401),
        ),
    );
    const { result } = renderHook(() => useGuestSession());

    await waitFor(() =>
      expect(result.current.state).toEqual({
        status: "UNAUTHENTICATED",
        prefill: "",
      }),
    );
  });

  it("honors Retry-After and prevents a second create call while one is pending", async () => {
    let finishCreate: ((response: Response) => void) | undefined;
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        jsonResponse({ error: { code: "SESSION_REQUIRED" } }, 401),
      )
      .mockImplementationOnce(
        () =>
          new Promise<Response>((resolve) => {
            finishCreate = resolve;
          }),
      );
    vi.stubGlobal("fetch", fetchMock);
    const { result } = renderHook(() => useGuestSession());
    await waitFor(() =>
      expect(result.current.state.status).toBe("UNAUTHENTICATED"),
    );

    act(() => {
      void result.current.create("Coco");
      void result.current.create("Coco");
    });
    await waitFor(() =>
      expect(result.current.state.status).toBe("CREATING_SESSION"),
    );
    expect(
      fetchMock.mock.calls.filter(
        ([, init]) => (init as RequestInit).method === "POST",
      ),
    ).toHaveLength(1);

    await act(async () => {
      finishCreate?.(
        jsonResponse({ error: { code: "RATE_LIMITED" } }, 429, {
          "Retry-After": "2",
        }),
      );
    });
    await waitFor(() => {
      expect(result.current.state).toMatchObject({
        status: "UNAUTHENTICATED",
        error: { code: "RATE_LIMITED", retryAfterSeconds: 2 },
      });
    });
  });

  it("does not crash when display-name storage is unavailable", async () => {
    const setItem = vi
      .spyOn(Storage.prototype, "setItem")
      .mockImplementation(() => {
        throw new DOMException("blocked", "SecurityError");
      });
    vi.stubGlobal("fetch", vi.fn().mockResolvedValueOnce(jsonResponse(guest)));
    const { result } = renderHook(() => useGuestSession());

    await waitFor(() =>
      expect(result.current.state.status).toBe("AUTHENTICATED"),
    );
    expect(setItem).toHaveBeenCalledWith(displayNameStorageKey, guest.username);
  });

  it("uses byte and grapheme lengths without approximating graphemes with string length", () => {
    expect(validateUsernameLength("กิข")).toEqual({ valid: true });
    expect(validateUsernameLength("Coco")).toEqual({ valid: true });
    expect(validateUsernameLength("a")).toEqual({
      valid: false,
      reason: "graphemeCount",
    });
    expect(validateUsernameLength("a".repeat(257))).toEqual({
      valid: false,
      reason: "tooLong",
    });
    expect(validateUsernameLength("😀😀")).toEqual({ valid: true });
  });
});
