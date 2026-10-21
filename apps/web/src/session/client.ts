export const guestSessionPath = "/api/v1/guest-session";
export const displayNameStorageKey = "prickle-panic:guest-display-name";

export interface GuestSession {
  guestId: string;
  username: string;
  expiresAt: string;
}

export type GuestSessionErrorCode =
  | "USERNAME_INVALID"
  | "REQUEST_INVALID"
  | "REQUEST_TOO_LARGE"
  | "ORIGIN_NOT_ALLOWED"
  | "RATE_LIMITED"
  | "SESSION_REQUIRED"
  | "SESSION_UNAVAILABLE"
  | "UNEXPECTED_RESPONSE";

export class GuestSessionError extends Error {
  readonly code: GuestSessionErrorCode;
  readonly status: number;
  readonly retryAfterSeconds?: number;

  constructor(
    code: GuestSessionErrorCode,
    status: number,
    retryAfterSeconds?: number,
  ) {
    super(code);
    this.name = "GuestSessionError";
    this.code = code;
    this.status = status;
    this.retryAfterSeconds = retryAfterSeconds;
  }
}

export async function recoverGuestSession(
  signal: AbortSignal,
): Promise<GuestSession> {
  return requestGuestSession("GET", undefined, signal);
}

export async function createGuestSession(
  username: string,
  signal: AbortSignal,
): Promise<GuestSession> {
  return requestGuestSession("POST", username, signal);
}

async function requestGuestSession(
  method: "GET" | "POST",
  username: string | undefined,
  signal: AbortSignal,
): Promise<GuestSession> {
  const response = await fetch(guestSessionPath, {
    method,
    credentials: "include",
    signal,
    ...(method === "POST"
      ? {
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ username }),
        }
      : {}),
  });

  if (method === "GET" && response.status === 200) {
    return readSession(response);
  }
  if (
    method === "POST" &&
    (response.status === 200 || response.status === 201)
  ) {
    return readSession(response);
  }

  throw await readApiError(response);
}

async function readSession(response: Response): Promise<GuestSession> {
  let value: unknown;
  try {
    value = await response.json();
  } catch {
    throw new GuestSessionError("UNEXPECTED_RESPONSE", response.status);
  }

  if (
    typeof value !== "object" ||
    value === null ||
    !("guestId" in value) ||
    typeof value.guestId !== "string" ||
    value.guestId.length === 0 ||
    !("username" in value) ||
    typeof value.username !== "string" ||
    value.username.length === 0 ||
    !("expiresAt" in value) ||
    typeof value.expiresAt !== "string" ||
    !Number.isFinite(Date.parse(value.expiresAt))
  ) {
    throw new GuestSessionError("UNEXPECTED_RESPONSE", response.status);
  }

  return {
    guestId: value.guestId,
    username: value.username,
    expiresAt: value.expiresAt,
  };
}

async function readApiError(response: Response): Promise<GuestSessionError> {
  let code: GuestSessionErrorCode = "UNEXPECTED_RESPONSE";
  try {
    const value: unknown = await response.json();
    if (
      typeof value === "object" &&
      value !== null &&
      "error" in value &&
      typeof value.error === "object" &&
      value.error !== null &&
      "code" in value.error &&
      isGuestSessionErrorCode(value.error.code)
    ) {
      code = value.error.code;
    }
  } catch {
    // Keep an unrecognized body behind a stable, player-safe error code.
  }

  const retryAfter = response.headers.get("Retry-After");
  const seconds =
    retryAfter && /^\d+$/.test(retryAfter) ? Number(retryAfter) : undefined;
  return new GuestSessionError(code, response.status, seconds);
}

function isGuestSessionErrorCode(
  value: unknown,
): value is GuestSessionErrorCode {
  return (
    value === "USERNAME_INVALID" ||
    value === "REQUEST_INVALID" ||
    value === "REQUEST_TOO_LARGE" ||
    value === "ORIGIN_NOT_ALLOWED" ||
    value === "RATE_LIMITED" ||
    value === "SESSION_REQUIRED" ||
    value === "SESSION_UNAVAILABLE"
  );
}
