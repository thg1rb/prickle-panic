import { useCallback, useEffect, useRef, useState } from "react";
import {
  createGuestSession,
  displayNameStorageKey,
  GuestSessionError,
  recoverGuestSession,
  type GuestSession,
  type GuestSessionErrorCode,
} from "./client";

export type GuestSessionState =
  | { status: "INITIALIZING" }
  | {
      status: "UNAUTHENTICATED";
      prefill: string;
      error?: { code: GuestSessionErrorCode; retryAfterSeconds?: number };
    }
  | { status: "CREATING_SESSION"; username: string }
  | { status: "AUTHENTICATED"; session: GuestSession }
  | { status: "RECOVERY_ERROR"; code: GuestSessionErrorCode };

export interface GuestSessionController {
  state: GuestSessionState;
  create(username: string): Promise<void>;
  retryRecovery(): void;
}

export function useGuestSession(): GuestSessionController {
  const [state, setState] = useState<GuestSessionState>({
    status: "INITIALIZING",
  });
  const [recoveryAttempt, setRecoveryAttempt] = useState(0);
  const operationId = useRef(0);
  const createInFlight = useRef(false);
  const mounted = useRef(false);
  const controllers = useRef(new Set<AbortController>());

  useEffect(() => {
    const activeControllers = controllers.current;
    mounted.current = true;
    return () => {
      mounted.current = false;
      operationId.current += 1;
      for (const controller of activeControllers) controller.abort();
      activeControllers.clear();
    };
  }, []);

  useEffect(() => {
    const activeControllers = controllers.current;
    const id = (operationId.current += 1 + recoveryAttempt);
    const controller = new AbortController();
    activeControllers.add(controller);

    void recoverGuestSession(controller.signal)
      .then((session) => {
        if (!mounted.current || id !== operationId.current) return;
        rememberDisplayName(session.username);
        setState({ status: "AUTHENTICATED", session });
      })
      .catch((error: unknown) => {
        if (!mounted.current || id !== operationId.current || isAbort(error))
          return;
        if (isUnauthenticated(error)) {
          setState({
            status: "UNAUTHENTICATED",
            prefill: readDisplayName(),
          });
          return;
        }
        setState({ status: "RECOVERY_ERROR", code: errorCode(error) });
      })
      .finally(() => activeControllers.delete(controller));

    return () => {
      controller.abort();
      activeControllers.delete(controller);
    };
  }, [recoveryAttempt]);

  useEffect(() => {
    if (
      state.status !== "UNAUTHENTICATED" ||
      !state.error?.retryAfterSeconds ||
      state.error.retryAfterSeconds <= 0
    ) {
      return undefined;
    }
    const interval = window.setInterval(() => {
      setState((current) => {
        if (
          current.status !== "UNAUTHENTICATED" ||
          current.error?.retryAfterSeconds === undefined
        ) {
          return current;
        }
        return {
          ...current,
          error: {
            ...current.error,
            retryAfterSeconds: Math.max(0, current.error.retryAfterSeconds - 1),
          },
        };
      });
    }, 1000);
    return () => window.clearInterval(interval);
  }, [state]);

  const create = useCallback(async (username: string) => {
    if (createInFlight.current) return;
    createInFlight.current = true;
    const id = ++operationId.current;
    const controller = new AbortController();
    controllers.current.add(controller);
    setState({ status: "CREATING_SESSION", username });

    try {
      const session = await createGuestSession(username, controller.signal);
      if (!mounted.current || id !== operationId.current) return;
      rememberDisplayName(session.username);
      setState({ status: "AUTHENTICATED", session });
    } catch (error) {
      if (!mounted.current || id !== operationId.current || isAbort(error))
        return;
      const code = errorCode(error);
      setState({
        status: "UNAUTHENTICATED",
        prefill: username,
        error: {
          code,
          ...(error instanceof GuestSessionError &&
          error.code === "RATE_LIMITED" &&
          error.retryAfterSeconds !== undefined
            ? { retryAfterSeconds: error.retryAfterSeconds }
            : {}),
        },
      });
    } finally {
      controllers.current.delete(controller);
      createInFlight.current = false;
    }
  }, []);

  const retryRecovery = useCallback(() => {
    if (state.status !== "RECOVERY_ERROR") return;
    setState({ status: "INITIALIZING" });
    setRecoveryAttempt((attempt) => attempt + 1);
  }, [state.status]);

  return { state, create, retryRecovery };
}

function readDisplayName(): string {
  try {
    const value = window.localStorage.getItem(displayNameStorageKey);
    if (value === null || new TextEncoder().encode(value).length > 256)
      return "";
    const normalized = value.normalize("NFC").trim();
    const count = countGraphemes(normalized);
    if (count !== undefined && (count < 2 || count > 16)) return "";
    return value;
  } catch {
    return "";
  }
}

function rememberDisplayName(username: string): void {
  try {
    window.localStorage.setItem(displayNameStorageKey, username);
  } catch {
    // Storage is only a convenience; session authorization remains server-side.
  }
}

export function validateUsernameLength(
  value: string,
): { valid: true } | { valid: false; reason: "tooLong" | "graphemeCount" } {
  if (new TextEncoder().encode(value).length > 256) {
    return { valid: false, reason: "tooLong" };
  }
  const candidate = value.normalize("NFC").trim();
  const count = countGraphemes(candidate);
  if (count === undefined) return { valid: true };
  if (count < 2 || count > 16) {
    return { valid: false, reason: "graphemeCount" };
  }
  return { valid: true };
}

function countGraphemes(value: string): number | undefined {
  if (typeof Intl.Segmenter !== "function") return undefined;
  return [
    ...new Intl.Segmenter("th", { granularity: "grapheme" }).segment(value),
  ].length;
}

function isUnauthenticated(error: unknown): boolean {
  return (
    error instanceof GuestSessionError &&
    error.status === 401 &&
    error.code === "SESSION_REQUIRED"
  );
}

function errorCode(error: unknown): GuestSessionErrorCode {
  if (error instanceof GuestSessionError) return error.code;
  return "SESSION_UNAVAILABLE";
}

function isAbort(error: unknown): boolean {
  return error instanceof DOMException && error.name === "AbortError";
}
