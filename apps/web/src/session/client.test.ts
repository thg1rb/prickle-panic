import { afterEach, describe, expect, it, vi } from "vitest";
import {
  createGuestSession,
  guestSessionPath,
  recoverGuestSession,
} from "./client";

afterEach(() => vi.unstubAllGlobals());

describe("Guest Session HTTP client", () => {
  it("recovers through same-origin credentials and validates the response", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      Response.json({
        guestId: "gst_example",
        username: "โคโค่",
        expiresAt: "2026-10-22T10:00:00Z",
      }),
    );
    vi.stubGlobal("fetch", fetchMock);

    await expect(
      recoverGuestSession(new AbortController().signal),
    ).resolves.toEqual({
      guestId: "gst_example",
      username: "โคโค่",
      expiresAt: "2026-10-22T10:00:00Z",
    });
    expect(fetchMock).toHaveBeenCalledWith(
      guestSessionPath,
      expect.objectContaining({
        method: "GET",
        credentials: "include",
        signal: expect.any(AbortSignal),
      }),
    );
  });

  it("creates or reuses a session with the exact request shape and no JS-visible credential", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      Response.json(
        {
          guestId: "gst_example",
          username: "Coco",
          expiresAt: "2026-10-22T10:00:00Z",
        },
        { status: 201 },
      ),
    );
    vi.stubGlobal("fetch", fetchMock);

    const session = await createGuestSession(
      "Coco",
      new AbortController().signal,
    );

    expect(session.guestId).toBe("gst_example");
    const [, request] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(fetchMock).toHaveBeenCalledWith(
      guestSessionPath,
      expect.any(Object),
    );
    expect(request).toMatchObject({
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: "Coco" }),
    });
    expect(JSON.stringify(session)).not.toMatch(/credential|token|cookie/i);
  });

  it("accepts POST reuse status and maps stable server errors without leaking response details", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        Response.json(
          {
            guestId: "gst_example",
            username: "Coco",
            expiresAt: "2026-10-22T10:00:00Z",
          },
          { status: 200 },
        ),
      )
      .mockResolvedValueOnce(
        Response.json(
          {
            error: {
              code: "RATE_LIMITED",
              detail: "private server diagnostic",
            },
          },
          { status: 429, headers: { "Retry-After": "7" } },
        ),
      );
    vi.stubGlobal("fetch", fetchMock);

    await expect(
      createGuestSession("Coco", new AbortController().signal),
    ).resolves.toMatchObject({
      guestId: "gst_example",
    });
    await expect(
      createGuestSession("Coco", new AbortController().signal),
    ).rejects.toMatchObject({
      code: "RATE_LIMITED",
      status: 429,
      retryAfterSeconds: 7,
    });
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("rejects malformed success payloads and does not expose raw error bodies", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(
        Response.json({ guestId: "gst", username: "Coco", expiresAt: "bad" }),
      )
      .mockResolvedValueOnce(
        Response.json(
          { error: { code: "NOT_A_PUBLIC_CODE", detail: "secret" } },
          { status: 500 },
        ),
      );
    vi.stubGlobal("fetch", fetchMock);

    await expect(
      recoverGuestSession(new AbortController().signal),
    ).rejects.toMatchObject({
      code: "UNEXPECTED_RESPONSE",
    });
    await expect(
      recoverGuestSession(new AbortController().signal),
    ).rejects.toMatchObject({
      code: "UNEXPECTED_RESPONSE",
      status: 500,
    });
  });
});
