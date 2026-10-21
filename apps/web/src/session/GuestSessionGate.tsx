import { useState, type FormEvent } from "react";
import { th } from "../i18n/th";
import type { GuestSessionState } from "./useGuestSession";
import { validateUsernameLength } from "./useGuestSession";

interface GuestSessionGateProps {
  state: GuestSessionState;
  onCreate: (username: string) => Promise<void>;
  onRetryRecovery: () => void;
}

export function GuestSessionGate({
  state,
  onCreate,
  onRetryRecovery,
}: GuestSessionGateProps) {
  if (state.status === "INITIALIZING") {
    return (
      <main className="app-shell guest-shell" lang="th" aria-busy="true">
        <section className="guest-card" aria-live="polite">
          <p className="section-kicker">{th.picnic}</p>
          <p>{th.guest.loading}</p>
        </section>
      </main>
    );
  }

  if (state.status === "RECOVERY_ERROR") {
    return (
      <main className="app-shell guest-shell" lang="th">
        <section className="guest-card" aria-labelledby="guest-title">
          <p className="section-kicker">{th.picnic}</p>
          <h1 id="guest-title">{th.guest.recoveryTitle}</h1>
          <p role="alert">{guestErrorMessage(state.code)}</p>
          <button
            className="primary-button"
            type="button"
            onClick={onRetryRecovery}
          >
            {th.guest.retry}
          </button>
        </section>
      </main>
    );
  }

  if (state.status === "AUTHENTICATED") return null;

  return (
    <UsernameEntry
      key={
        state.status === "UNAUTHENTICATED"
          ? `entry-${state.prefill}`
          : `creating-${state.username}`
      }
      state={state}
      onCreate={onCreate}
    />
  );
}

type EntryState = Extract<
  GuestSessionState,
  { status: "UNAUTHENTICATED" | "CREATING_SESSION" }
>;

function UsernameEntry({
  state,
  onCreate,
}: {
  state: EntryState;
  onCreate: (username: string) => Promise<void>;
}) {
  const [username, setUsername] = useState(
    state.status === "UNAUTHENTICATED" ? state.prefill : state.username,
  );
  const [touched, setTouched] = useState(false);
  const isCreating = state.status === "CREATING_SESSION";
  const validation = validateUsernameLength(username);
  const waitSeconds =
    state.status === "UNAUTHENTICATED"
      ? (state.error?.retryAfterSeconds ?? 0)
      : 0;
  const errorCode =
    state.status === "UNAUTHENTICATED" ? state.error?.code : undefined;

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (isCreating || !validation.valid || waitSeconds > 0) return;
    void onCreate(username);
  }

  return (
    <main className="app-shell guest-shell" lang="th">
      <section className="guest-card" aria-labelledby="guest-title">
        <p className="section-kicker">{th.picnic}</p>
        <div className="guest-mascot" aria-hidden="true">
          🦔
        </div>
        <h1 id="guest-title">{th.guest.title}</h1>
        <p className="guest-description">{th.guest.description}</p>
        {errorCode && (
          <p className="guest-error" role="alert">
            {guestErrorMessage(errorCode, waitSeconds)}
          </p>
        )}
        <form onSubmit={submit} noValidate>
          <label htmlFor="guest-username">{th.guest.usernameLabel}</label>
          <input
            autoComplete="off"
            id="guest-username"
            name="username"
            placeholder={th.guest.placeholder}
            value={username}
            onChange={(event) => {
              setUsername(event.currentTarget.value);
              setTouched(true);
            }}
            onBlur={() => setTouched(true)}
            aria-describedby="guest-username-hint"
            aria-invalid={touched && !validation.valid}
            disabled={isCreating}
          />
          <p id="guest-username-hint" className="guest-hint">
            {th.guest.usernameHint}
          </p>
          {touched && !validation.valid && (
            <p className="guest-error" role="alert">
              {validation.reason === "tooLong"
                ? th.guest.tooLong
                : th.guest.graphemeCount}
            </p>
          )}
          <button
            className="primary-button guest-submit"
            type="submit"
            disabled={isCreating || !validation.valid || waitSeconds > 0}
          >
            {isCreating
              ? th.guest.creating
              : waitSeconds > 0
                ? `${th.guest.wait} ${waitSeconds}`
                : th.guest.submit}
          </button>
          {isCreating && (
            <output className="guest-loading" aria-live="polite">
              {th.guest.creating}
            </output>
          )}
        </form>
      </section>
    </main>
  );
}

function guestErrorMessage(code: string, waitSeconds = 0): string {
  if (code === "USERNAME_INVALID") return th.guest.errors.usernameInvalid;
  if (code === "REQUEST_INVALID") return th.guest.errors.requestInvalid;
  if (code === "REQUEST_TOO_LARGE") return th.guest.errors.requestTooLarge;
  if (code === "ORIGIN_NOT_ALLOWED") return th.guest.errors.originNotAllowed;
  if (code === "RATE_LIMITED") {
    return waitSeconds > 0
      ? `${th.guest.errors.rateLimited} ${waitSeconds} ${th.guest.seconds}`
      : th.guest.errors.rateLimited;
  }
  if (code === "SESSION_REQUIRED") return th.guest.errors.sessionRequired;
  if (code === "SESSION_UNAVAILABLE") return th.guest.errors.sessionUnavailable;
  return th.guest.errors.unexpected;
}
