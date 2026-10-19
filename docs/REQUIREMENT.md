# 🦔 Prickle Panic!
# Project Requirements & Development Specification

| Field | Value |
|---|---|
| Project | Prickle Panic! |
| Document | Project Requirements & Development Specification |
| Version | 1.0.0 |
| Status | BASELINE |
| Last Updated | 2026-10-19 |
| Development Methodology | Agent-first Spec-Driven Development (SDD) |
| Platform | Web Browser |
| Game Genre | Online Multiplayer Party / Puzzle / Casual |
| Primary Language | Thai |
| Development Strategy | MVP-first |
| Documentation Source | `docs/REQUIREMENT.md` |

---

# 1. Document Purpose & Governance

## 1.1 Purpose

This document defines the product requirements, gameplay rules,
technical direction, development standards, and agent operating
requirements for Prickle Panic!.

It is the initial authoritative specification for the project.

The project will be developed primarily through AI Coding Agents
using Spec-Driven Development (SDD).

All agents must understand and follow this document before
designing, implementing, reviewing, or modifying the system.

## 1.2 Single Source of Truth

During Project Bootstrap:

`docs/REQUIREMENT.md`

is the authoritative source of project requirements.

Agents may subsequently create specialized documentation under
`/docs`, including:

- Product specifications
- Feature specifications
- Technical architecture
- Architecture Decision Records (ADR)
- Development plans
- Testing documentation
- Agent instructions
- Skill registry
- Independent review reports

However, derived documentation must remain consistent with this
document.

No agent may silently introduce requirements that contradict
approved decisions.

When documentation is reorganized, each requirement must retain
its identifier and traceability.

This document must not be discarded merely because the agent has
generated more detailed specifications.

## 1.3 Requirement Status

Requirements and decisions use the following statuses.

### CONFIRMED

Explicitly accepted by the Project Owner.

Agents must not change these decisions without authorization.

### BASELINE

Current agreed design direction.

May be refined during detailed technical design without changing
the fundamental product requirements.

Material behavior changes require review.

### TUNING

Configurable gameplay values that require playtesting.

Examples:

- Match duration
- Board size
- Cooldown duration
- Event duration
- Score modifiers

Agents must use the documented initial values.

Changes must be supported by testing or an approved design
adjustment.

### PROPOSED

Suggested implementation details that have not yet been approved.

### DEFERRED

Requirements intentionally excluded from the current milestone.

### OUT_OF_SCOPE

Features not planned for the current product scope.

## 1.4 Requirement Priority

| Priority | Meaning |
|---|---|
| P0 | Required for the MVP |
| P1 | Required for planned post-MVP game modes or enhancements |
| P2 | Optional future improvement |
| P3 | Long-term consideration |

## 1.5 Change Control

Agents must not modify confirmed requirements without
Project Owner approval.

When a specification conflict is discovered:

1. Identify the conflicting requirements.
2. Explain the impact.
3. Propose possible resolutions.
4. Record the decision.
5. Update the affected specifications.
6. Update the implementation plan.
7. Continue development.

Do not modify requirements merely to justify an implementation
that does not meet the original specification.

---

# 2. Product Overview

## 2.1 Game Concept

Prickle Panic! is a fast-paced online multiplayer party game
inspired by visual card-matching mechanics.

Players compete by observing hedgehog cards and identifying
which cards satisfy the matching rule.

The game combines:

- Pattern recognition
- Observation
- Quick decision-making
- Accuracy
- Competitive multiplayer
- Cute character design

The game should be easy to understand but rewarding to master.

## 2.2 Product Vision

Create a cute, accessible, and competitive multiplayer party
game that groups of friends can enjoy directly through a
web browser without installing an application.

The intended experience is:

> Easy to learn, fun to master, and enjoyable with friends.

## 2.3 Product Principles

### Easy to Access

Players should enter the game through a browser.

No application installation is required.

### Easy to Learn

The matching rule should be understandable through a short
interactive tutorial.

### Quick to Play

Matches should be short enough to encourage rematches.

### Social Competition

Players should be able to create rooms and invite friends.

### Cute but Competitive

The game should combine a friendly visual identity with
competitive gameplay.

### Fair Gameplay

Game outcomes must be determined by valid gameplay actions
and authoritative server rules.

### Low Friction

Joining a game should require as few steps as possible.

### Mobile Friendly

The game must remain usable on smartphones.

## 2.4 Target Audience

Primary audiences:

- Casual gamers
- Groups of friends
- Party game enthusiasts
- Players who enjoy visual puzzles
- Players who prefer short multiplayer sessions

The product does not initially target professional esports
or highly competitive ranked matchmaking.

## 2.5 Supported Platforms

The game is Web-first.

Required:

- Desktop Web Browsers
- Mobile Web Browsers
- Tablet Web Browsers

Initial browser targets:

- Chrome
- Safari
- Edge
- Firefox

The application should use responsive layouts and support
both mouse and touch input.

Native desktop and mobile applications are not required.

## 2.6 Player Count

Minimum players per Match: 2

Maximum players per Room: 6

The MVP focuses on private multiplayer rooms.

## 2.7 Game Name

Working title:

**Prickle Panic!**

The name is the current design baseline.

Final commercial use remains subject to naming, domain,
and trademark verification.

## 2.8 Mascot

Name: **Coco (โคโค่)**

Character: Hedgehog

Coco is the primary mascot and visual identity of the game.

Coco should appear throughout the product, including:

- Main Menu
- Tutorial
- Lobby
- Gameplay feedback
- Results screen
- Loading states
- Error messages

---

# 3. Product Scope

## 3.1 MVP Scope

The MVP must include:

- Web-based game client
- Thai user interface
- Guest username
- Guest session
- Create Room
- Join Room by Code
- Join Room by Invitation URL
- Lobby
- Player list
- Ready system
- Host-controlled Match start
- Speed Battle
- Core card matching
- Real-time score synchronization
- Match timer
- Match results
- Rematch
- Basic room moderation
- Server-authoritative validation
- Basic anti-cheat protection
- Reconnection
- Responsive desktop and mobile interface
- Coco character assets
- Basic sound effects
- Basic background music
- Automated tests
- Independent agent review

## 3.2 Planned Post-MVP Scope

- Classic Rush
- Party Chaos
- Expanded animations
- Additional audio assets
- More advanced room management
- Expanded gameplay statistics
- English localization

Both Classic Rush and Party Chaos must be documented and
considered in the architecture from the beginning.

However, they must not be implemented before the core
multiplayer MVP is sufficiently stable.

## 3.3 Out of Scope for MVP

The MVP does not require:

- User registration
- Email/password login
- OAuth
- Persistent user accounts
- Global ranked matchmaking
- Global leaderboard
- Public matchmaking
- Text chat
- Voice chat
- Spectator mode
- Paid cosmetics
- In-app purchases
- Persistent inventory
- Native applications
- Microservices
- Distributed game servers
- Persistent match history
- AI-controlled opponents

These features may be reconsidered in future milestones.

---

# 4. Visual Identity & Art Direction

## 4.1 Visual Theme

The visual concept is:

**Coco's Berry Picnic**

The game should feel like a cute and playful woodland world
filled with warm colors and soft illustrations.

## 4.2 Art Style

Selected direction:

- Cute 2D Cartoon
- Cozy Woodland
- Warm Pastel
- Soft Shading
- Rounded Shapes
- Minimal Visual Clutter
- Playful Animations

The game should look friendly and charming without
compromising gameplay readability.

## 4.3 Color Palette

Initial color palette:

| Role | Color | Hex |
|---|---|---|
| Background | Warm Cream | `#FFF7E9` |
| Primary | Strawberry Coral | `#F27F7A` |
| Secondary | Sage Green | `#88A889` |
| Accent | Honey Yellow | `#F4C66B` |
| Character | Cocoa Brown | `#654736` |
| Text | Dark Cocoa | `#302820` |

These values are design tokens and may be refined
during implementation.

Text contrast and accessibility take precedence
over purely decorative color choices.

## 4.4 Typography

Primary font:

**Mitr**

Used for:

- Buttons
- Menus
- Player names
- Scores
- Room information
- Settings
- Body text
- Gameplay UI

Decorative font:

**Mali**

Used for:

- Large headings
- Coco's messages
- Victory and defeat messages
- Decorative labels

The fonts must support Thai correctly.

Font files should be self-hosted where practical,
with licensing information preserved.

## 4.5 Coco Character Identity

Coco should have:

- A rounded body
- Short limbs
- Soft rounded quills
- Large expressive eyes
- Small nose
- Cute facial expressions
- A recognizable silhouette

Character personality:

- Playful
- Curious
- Energetic
- Friendly
- Slightly mischievous

## 4.6 Coco Animations

Planned character reactions:

| Event | Animation |
|---|---|
| Idle | Blink / Small movement |
| Correct Pick | Happy reaction |
| Incorrect Pick | Surprised reaction |
| Combo | Excited reaction |
| Countdown | Anticipation |
| Victory | Celebration |
| Defeat | Lighthearted disappointment |
| Loading | Playful waiting animation |

Not all animations are mandatory in the initial MVP.

## 4.7 Card Visual Requirements

Cards must:

- Use a consistent visual layout
- Have rounded corners
- Have clear distinguishing attributes
- Avoid unnecessary decorative elements
- Remain readable on mobile
- Support mouse and touch interaction
- Maintain a stable aspect ratio
- Avoid visually ambiguous attributes

Gameplay cards should not contain decorative accessories
that may be confused with gameplay attributes.

## 4.8 AI-assisted Asset Production

AI Agents may assist with:

- Character concept art
- Character variations
- Card illustrations
- UI assets
- Icons
- Backgrounds
- Animations
- Background music
- Sound effects
- Asset optimization

Assets must be reviewed before use.

Generated artwork must be visually consistent.

The project should use a modular character system rather
than independently generating 32 unrelated character images.

## 4.9 Asset Source Requirements

Keep editable source assets when practical.

Examples:

- SVG
- Layered image sources
- Animation source data
- Audio source projects
- Asset generation scripts

Generated production assets may include:

- PNG
- WebP
- SVG
- Audio files
- Sprite sheets

Agent-generated assets must be tested in the actual
gameplay layout before final acceptance.

## 4.10 Asset Licensing

Do not copy original artwork, branding, or protected assets
from Pick a Dog or another existing game without permission.

Third-party assets must have acceptable licensing.

For every externally sourced asset, record:

- Asset name
- Source
- Creator or provider
- License
- Relevant attribution
- Usage restrictions

For AI-generated assets, record:

- Generation method
- Source tool or service
- Relevant usage terms
- Original prompts or generation references when useful

Do not assume that all AI-generated music or images
automatically have unrestricted commercial rights.

---

# 5. Localization

## LOC-001 — Primary Language

Thai is the primary language for the MVP.

All player-facing text should be written in Thai.

## LOC-002 — Future English Support

English localization is deferred until the core game
is playable.

## LOC-003 — Translation-ready Structure

User-facing text must not be scattered throughout
business logic.

Use centralized translation keys or equivalent structure.

Example:

```typescript
{
  room: {
    create: "สร้างห้อง",
    join: "เข้าร่วมห้อง"
  },
  game: {
    correct: "ถูกต้อง!",
    incorrect: "ลองใหม่อีกครั้ง!"
  }
}
```

## LOC-004 — Server Messages

Server error codes must remain language-independent.

The client translates codes into appropriate
Thai messages.

## LOC-005 — Thai Unicode Support

Usernames and text validation must correctly handle
Thai Unicode characters and combining marks.

---

# 6. Guest Player System

## 6.1 No Login Requirement

Players must be able to play without creating accounts.

The system must not require:

- Email
- Password
- Phone number
- Social login

## 6.2 Username

Players must enter a username before accessing gameplay.

Initial username policy:

| Rule | Value |
|---|---|
| Minimum Length | 2 Unicode extended grapheme clusters |
| Maximum Length | 16 Unicode extended grapheme clusters |
| Normalization | NFC after trimming surrounding whitespace |
| Allowed characters | Thai letters with valid associated Thai marks, English letters, ASCII digits, Thai digits, and single internal ASCII spaces |
| Duplicate in Same Room | Not Allowed |
| Duplicate Across Rooms | Allowed |

The server must reject empty or whitespace-only values, controls, newlines,
tabs, invisible format and zero-width controls, bidirectional controls,
standalone or malformed combining marks, emoji, unsupported scripts,
unsupported punctuation, markup, and overlong raw input. The raw input byte
limit is 256 bytes before normalization. Thai marks associated with valid
Thai letters remain permitted. The server is authoritative. The frontend
must render usernames as text, never markup.

For same-Room duplicate detection, compare a normalized key that treats
English letters case-insensitively while preserving the NFC display casing.
Usernames are display names and are never authentication credentials.

## 6.3 Username Persistence

The client may remember the last username locally.

Local storage may contain the display username.

Authentication secrets must not be stored in localStorage.

## 6.4 Guest Identity

The server must generate a unique guest session.

Username must not be used as an authentication credential.

Guest identity and display username must be separate.

## 6.5 Guest Session Security

Guest sessions must use unpredictable session credentials.

The server must issue an opaque credential containing at least 256 bits of
cryptographically secure randomness. It must not encode identity or sensitive
information, be derived from a username or timestamp, be logged, exposed in a
URL, or stored in localStorage. The server should store a cryptographic hash
of the credential and must support server-side invalidation. Reconnecting a
WebSocket must not rotate the credential.

Session state is in memory for MVP; recovery across a server restart is not
guaranteed. The configurable inactivity timeout is 24 hours since the last
qualifying authenticated application activity, and the configurable absolute
lifetime is 7 days from creation. The earlier deadline expires the session.
Successful authenticated application activity and explicit session recovery
qualify; WebSocket ping/pong, passive snapshots, background heartbeats, and
failed requests do not. Expiration revokes authorization consistently for
HTTP and WebSocket actions. It is separate from the 15-second Match network
reconnection grace period; if the session expires during an active Match,
authorization is revoked and the player immediately forfeits without that
network grace period.

Production cookie name is `__Host-prickle_guest`, with:

- HttpOnly
- Secure
- SameSite=Lax
- Path=/
- No Domain attribute
- Max-Age consistent with the 7-day absolute lifetime

Local HTTP development must use an explicit development-only cookie
configuration and must not weaken the production policy. Production uses
HTTPS/WSS. Cookie-authenticated WebSockets require strict Origin validation
and message-level authorization; SameSite is not the sole CSRF or WebSocket
security boundary.

Create no more than 10 new Guest Sessions per 10 minutes per source IP as
the configurable initial baseline. Reuse an existing valid session on page
refresh. Apply the limit server-side and return a rate-limit response.
Room creation is initially limited to 5 attempts per 10 minutes per Guest
Session and joining to 20 attempts per minute per Guest Session. Additional
network-source controls may mitigate repeated session creation; IP addresses
are not permanent player identities. Limits require integration and abuse
testing and must avoid unfairly blocking shared networks.

The server must validate the session before authorizing
room and gameplay actions.

## 6.6 Multiple Browser Tabs

A guest session must not control the same RoomPlayer
through multiple active connections simultaneously.

If a valid session reconnects from another tab,
the server must establish which connection is active.

The older connection must lose gameplay authority.

## 6.7 Username Changes

Players may change their username in the Lobby.

Username changes are not permitted during an active Match.

Changing a username must not change the player's
identity, score, or session ownership.

---

# 7. Room Management

## 7.1 Room Creation

Any guest player may create a private room.

The creator becomes the initial Host.

The server generates:

- Unique internal Room ID
- Shareable Room Code
- Invitation URL

## 7.2 Room Code

Initial proposal:

Six-character code.

Example:

`K7P2MX`

Codes must be generated using a cryptographically secure
random generator.

Avoid confusing characters when practical.

Room codes must be unique among active rooms.

## 7.3 Invitation URL

Example:

`https://example.com/room/K7P2MX`

The actual production domain is not yet selected.

Opening an invitation URL should:

1. Read the Room Code.
2. Check the guest session.
3. Request a username if necessary.
4. Attempt to join the room.
5. Show the Lobby if successful.

The invitation destination must be preserved while
the user enters their username.

## 7.4 Room Capacity

Minimum Match size: 2 players.

Maximum Room capacity: 6 players.

Disconnected players with an active reserved slot
must count toward the capacity limit.

## 7.5 Room States

The room supports:

- LOBBY
- COUNTDOWN
- IN_PROGRESS
- RESULTS
- CLOSED

Valid primary flow:

LOBBY
→ COUNTDOWN
→ IN_PROGRESS
→ RESULTS
→ LOBBY

A room may transition to CLOSED when it is terminated
or expires.

## 7.6 Lobby

The Lobby must display:

- Room Code
- Share Link
- Player list
- Host indicator
- Ready status
- Player count
- Selected game mode
- Match settings
- Start control
- Leave control

## 7.7 Ready System

Every participant, including the Host, must be Ready.

The Host may start a Match only when:

- At least 2 players are present.
- No more than 6 players are present.
- Every participant is connected.
- Every participant is Ready.
- Room settings are valid.
- No Match is currently running.

## 7.8 Settings Changes

If the Host changes gameplay settings:

All players' Ready states must be reset.

Players must confirm readiness under the new settings.

## 7.9 Host Permissions

The Host may:

- Select game mode in Lobby
- Change allowed Match settings in Lobby
- Start the Match when everyone is Ready
- Cancel Countdown
- Remove a player from Lobby
- Lock the Room
- Regenerate Room Code in Lobby

The Host must not:

- Change scores
- Modify card rules during Match
- Kick players during Match
- Change Match duration during Match
- Change Event Schedule during Match
- Cancel Match merely to avoid losing

## 7.10 Kick Player

The Host may remove a player during Lobby.

The removed guest session must not be allowed to
rejoin the same room using that session.

Because players do not have verified accounts,
permanent identity-based bans are not guaranteed.

## 7.11 Lock Room

A locked room rejects new participants.

Authorized reconnections must remain possible.

## 7.12 Regenerate Room Code

The Host may regenerate the Room Code in Lobby.

The old code and invitation URL become invalid.

Existing members remain in the room.

## 7.13 Host Migration

If the Host voluntarily leaves:

The server transfers Host ownership to the eligible
remaining player with the earliest join order.

For unexpected disconnection, apply the reconnect
grace period before transferring Host ownership.

Host migration must not modify active gameplay results.

## 7.14 Joining During Match

New players cannot join an active Match.

For the MVP, there is no Spectator Mode or Waiting Queue.

An authorized reconnecting player is not considered
a new participant.

## 7.15 Room Expiration

Initial configuration:

- Empty room expiration: approximately 5 minutes
- Inactive Lobby expiration: approximately 30 minutes

Exact values are TUNING parameters.

Active Matches must not be terminated solely by
the Lobby expiration policy.

Room cleanup must release relevant resources.

---

# 8. Core Card & Matching System

## 8.1 Card Attributes

Each Coco card contains five binary attributes.

| Attribute | Value 0 | Value 1 |
|---|---|---|
| Body Color | Brown | Cream |
| Quills | Short | Long |
| Glasses | No | Yes |
| Scarf | No | Yes |
| Strawberry | No | Yes |

The number of possible variations is:

2^5 = 32

## 8.2 Card Representation

A card variation may be represented using five bits.

Example:

`10110`

Each bit represents one attribute.

## 8.3 Card Instance

A card variation is not the same as a card instance.

Multiple card instances may share the same variation.

Each card instance must have an identifier unique within
the relevant game-state scope.

Example:

```text
Card Instance
  ID: card-012
  Variant: 10110
  Position: 5
```

## 8.4 Reference Card

The Reference Card is the card against which a player
compares available cards.

It is displayed separately from selectable board cards.

The Reference Card itself cannot be selected.

## 8.5 Matching Rule

A card is valid when it differs from the Reference Card
in zero or one attribute.

A card is invalid when it differs in two or more attributes.

The formal rule is:

`HammingDistance(reference, selected) <= 1`

For binary card representations, XOR and bit counting
may be used.

## 8.6 Correct Selection

When a selection is valid:

1. Server validates the action.
2. Server awards the applicable score.
3. Server updates the player's Reference Card.
4. Server removes the selected Card Instance.
5. Server applies mode-specific board updates.
6. Server updates Combo state.
7. Server publishes the accepted state.

These changes must be committed consistently.

## 8.7 Incorrect Selection

When a selection is an actual invalid match:

- No points are awarded.
- Reference Card does not change.
- Board does not change.
- Combo resets.
- Mistake Cooldown is applied.

Default cooldown:

0.8 seconds.

## 8.8 Cooldown Behavior

Actions received during cooldown must not:

- Award points
- Update the board
- Update the reference
- Extend the existing cooldown

Gameplay cooldown and network abuse rate limiting
must be handled separately.

## 8.9 Technical Rejections

The following are not gameplay mistakes:

- Duplicate request
- Unavailable card
- Stale state
- Connection transition
- Event transition
- Match not active
- Previously processed request

Technical rejections must not create gameplay
mistake penalties.

## 8.10 Combo

The system counts consecutive valid selections.

Example:

Correct → 1

Correct → 2

Correct → 3

Incorrect → 0

For Speed Battle and Classic Rush, Combo initially
affects only audiovisual feedback.

Party Chaos may apply event-specific bonuses.

---

# 9. Game Mode — Speed Battle

## 9.1 Overview

Speed Battle is the primary MVP game mode.

Players compete to earn the highest score within
a fixed time limit.

Each player has an independent card board.

## 9.2 Default Settings

| Setting | Value |
|---|---|
| Players | 2–6 |
| Duration | 60 seconds |
| Board | 4×4 |
| Cards | 16 |
| Correct Pick | +1 |
| Incorrect Pick | 0 |
| Mistake Cooldown | 0.8 seconds |
| Combo | Visual only |
| Card Refill | Continuous |
| Winner | Highest score |

These values are configurable TUNING parameters.

## 9.3 Starting Conditions

At Match start:

- Every player has zero points.
- Every player has Combo zero.
- Every player receives the same starting Reference Card.
- Every player receives the same starting board.
- Every player has an independent board state.

## 9.4 Board Ownership

A player's board belongs exclusively to that player.

Selecting a card must not remove cards from another
player's board.

## 9.5 Correct Pick

When a player selects a valid card:

1. Award one point.
2. Increase Combo.
3. Change Reference Card.
4. Remove selected Card Instance.
5. Generate a replacement card.
6. Restore the board to 16 cards.
7. Publish the updated player state.

## 9.6 Card Refill

A replacement card should occupy the selected
card's position.

The server must ensure that a valid selectable card
always remains available.

Invariant:

`exists card in board where ValidMatch(reference, card)`

This applies whenever the board is open for interaction.

## 9.7 Random Generation

The server controls card generation.

Clients must not generate authoritative replacement cards.

All players begin with the same starting configuration.

The server must use a consistent generation policy.

Random generation should support deterministic reproduction
for testing.

Do not reveal future card sequences or sensitive random
generation state to clients.

Identical starting conditions do not guarantee identical
difficulty after players make different choices.

This is an accepted design characteristic.

## 9.8 Incorrect Pick

An actual incorrect pick:

- Awards zero points.
- Resets Combo.
- Applies cooldown.
- Does not change board or reference.

## 9.9 AFK Behavior

An inactive player does not stop the Match timer.

The player's score remains unchanged.

The player may resume selecting cards while still eligible
and while the Match remains active.

## 9.10 End Conditions

The Match ends when:

1. The Match timer expires.
2. Only one eligible competitor remains after forfeits.
3. The server cancels the Match because a valid game
   state cannot be maintained.

A disconnected player inside the grace period remains
eligible until the applicable timeout or Match-end rule
resolves their status.

The Match must not award an early forfeit victory while
another player remains eligible to reconnect.

## 9.11 Winner

The highest-scoring eligible player wins.

Ties are allowed.

Players tied at the highest score are joint winners.

Forfeited players are not eligible to win.

## 9.12 Invalid Board Recovery

If an impossible board state is detected:

- Do not penalize the player.
- Attempt safe server-side correction.
- Synchronize the corrected board.
- Record the anomaly.

If the state cannot be recovered safely, cancel the Match
without declaring a winner.

---

# 10. Game Mode — Classic Rush

## 10.1 Overview

Classic Rush is a real-time shared-board competition.

Players compete to select cards from the same board.

Once a card is successfully selected, it disappears
for everyone.

## 10.2 Default Settings

| Setting | Value |
|---|---|
| Players | 2–6 |
| Waves | 3 |
| Wave Duration | 30 seconds |
| Board | Shared 4×4 |
| Cards per Wave | 16 |
| Correct Pick | +1 |
| Incorrect Pick | 0 |
| Mistake Cooldown | 0.8 seconds |
| Card Refill | None during Wave |
| Winner | Highest cumulative score |

## 10.3 Starting a Wave

At the start of each Wave:

- Generate a new shared board.
- Generate the initial Reference Card.
- Give every participant the same initial Reference Card.
- Reset wave-specific Combo.
- Reset wave-specific completion state.

The starting board should contain at least as many valid
selectable Card Instances as the number of participants.

This is an initial fairness measure, not a guarantee that
all players will retain a valid choice throughout the Wave.

## 10.4 Shared Board

Every player sees the same available Card Instances.

A successful selection removes the selected instance
from the shared board.

Cards are not refilled during the Wave.

## 10.5 Individual Reference Cards

Each player maintains an individual Reference Card.

After a successful pick, only the selecting player's
Reference Card changes.

Other players' Reference Cards remain unchanged.

## 10.6 Contested Selection

If two players select the same Card Instance:

The server processes authoritative board mutations
in a deterministic order.

Only one request may successfully claim the card.

The later unsuccessful request receives a technical
rejection without a mistake penalty.

## 10.7 Incorrect Pick

An invalid match:

- Awards no points.
- Resets Combo.
- Applies 0.8-second cooldown.
- Does not change the shared board.
- Does not change the player's reference.

## 10.8 No-Match Declaration

Players may declare:

"ไม่มีใบที่ตรงแล้ว!"

The declaration is valid only when no currently available
shared-board card matches that player's Reference Card.

## 10.9 Correct No-Match Declaration

If the declaration is correct:

- No score penalty.
- The player enters WAVE_FINISHED.
- The player can no longer select cards during this Wave.
- Other players continue.

A correct declaration does not automatically end
the entire Wave.

## 10.10 Incorrect No-Match Declaration

If at least one valid card exists:

- Deduct one point.
- Score must not go below zero.
- Reset Combo.
- Apply a two-second cooldown.
- Allow the player to resume afterward.

## 10.11 No-Match Abuse Prevention

Initial policy:

Maximum two incorrect No-Match declarations per player
per Wave.

After reaching the limit:

- Disable the declaration for the rest of that Wave.
- The player may continue selecting cards after cooldown.
- The server may still automatically finish the Wave.

This limit is a TUNING parameter.

## 10.12 No-Match Race Conditions

The server evaluates the declaration against the
authoritative shared board at processing time.

If another player removed the last valid card before
the declaration was processed, the declaration may
be accepted as correct.

## 10.13 Automatic Wave Completion

A Wave ends when:

- The timer expires.
- All shared cards have been selected.
- No active eligible participant can make another valid pick.

A player who is AFK must not prevent Wave completion
indefinitely.

The Wave timer remains authoritative.

## 10.14 Wave Transition

When a Wave ends:

1. Stop processing actions for the old Wave.
2. Preserve cumulative scores.
3. Clear the old shared board.
4. Generate a new board.
5. Generate a new starting Reference Card.
6. Reset wave-specific state.
7. Begin the next Wave countdown.

Old Wave actions must never be accepted in the new Wave.

## 10.15 Match Completion

After Wave 3:

The server calculates cumulative scores.

The highest score wins.

Ties are allowed.

## 10.16 Network Fairness

Classic Rush is a casual competitive mode.

Network latency may affect contested card selections.

The MVP does not promise complete latency equality.

The server must consistently apply the documented
action-ordering policy.

---

# 11. Game Mode — Party Chaos

## 11.1 Overview

Party Chaos extends the private-board gameplay
of Speed Battle with server-controlled special events.

The mode emphasizes variety and playful competition.

## 11.2 Default Settings

| Setting | Value |
|---|---|
| Players | 2–6 |
| Active Gameplay Duration | 75 seconds |
| Board | Private 4×4 |
| Correct Pick | +1 |
| Mistake Cooldown | 0.8 seconds |
| Card Refill | Continuous |
| Special Events | 3 |
| Winner | Highest score |

## 11.3 Event Principles

Special events must:

- Be controlled by the server.
- Follow a Match-level schedule.
- Apply consistent rules to all eligible players.
- Avoid targeting individual opponents.
- Have clearly communicated effects.
- Preserve authoritative score calculation.

Clients and Hosts cannot activate arbitrary events.

## 11.4 Basket Shuffle

Thai name:

"สลับตะกร้า!"

Effect:

Rearrange the positions of cards on the player's board.

Rules:

- Card identities remain unchanged.
- Card attributes remain unchanged.
- Reference Card remains unchanged.
- Score remains unchanged.
- Card validity remains unchanged.

The event must notify players before activation.

Initial warning period:

3 seconds.

Initial transition duration:

0.4 seconds.

During the synchronized transition:

- New picks are temporarily disabled.
- Active gameplay time pauses.
- All players receive equivalent transition treatment.
- Stale actions are rejected without gameplay penalty.

## 11.5 Berry Boost

Thai name:

"เบอร์รีคูณสอง!"

Effect:

Correct picks award two points instead of one.

Default duration:

8 seconds.

Incorrect picks still award no points.

Combo behavior remains subject to ordinary rules.

The server determines whether a pick qualifies
for the bonus using authoritative event timing.

## 11.6 Coco Combo

Thai name:

"โคโค่เชียร์อยู่!"

Effect:

Award a bonus when a player makes three consecutive
correct picks during the event.

Initial configuration:

- Required Combo: 3
- Bonus: +2 points
- Duration: 10 seconds
- Reward Limit: Once per player per event

The event-specific Combo starts at zero.

Combo accumulated before event activation does not count.

An incorrect pick resets the event-specific Combo.

## 11.7 Event Schedule

Initial 75-second active-time schedule:

| Active Time | Event |
|---|---|
| 0–15s | Normal Gameplay |
| 15s | Basket Shuffle |
| 32–40s | Berry Boost |
| 55–65s | Coco Combo |
| 65–75s | Final Rush |

Basket Shuffle transition time is excluded from
active gameplay time.

Server time remains authoritative.

The schedule may be refined during playtesting.

## 11.8 Event Fairness

Events must not randomly provide better rules
to individual players.

Every participant receives the same event rules.

Board states may still differ based on individual choices.

## 11.9 Disallowed Mechanics for Initial Version

Do not introduce:

- Score stealing
- Opponent freezing
- Host-controlled event targeting
- Arbitrary score deductions
- Forced penalties targeting individual players
- Hidden rule changes without warnings

These features require separate design review.

## 11.10 Incorrect Actions During Events

During an event transition:

Invalid timing or stale state must produce a technical
rejection, not a gameplay mistake.

Normal invalid card matching outside transitions
still receives the ordinary mistake penalty.

## 11.11 Winner

The highest-scoring eligible player wins.

Final score includes:

- Normal pick points
- Event-modified pick points
- Earned event bonuses

Ties are allowed.

Forfeited players cannot win.

---

# 12. Match Lifecycle

## 12.1 Match States

Each Match supports:

- PREPARING
- COUNTDOWN
- IN_PROGRESS
- COMPLETED
- CANCELLED

Match state must be controlled by the server.

## 12.2 Match Preparation

When the Host starts a Match:

1. Validate Room state.
2. Validate player count.
3. Validate Ready states.
4. Validate Match settings.
5. Create a new Match ID.
6. Generate initial game state.
7. Validate the initial board.
8. Schedule the authoritative start time.

## 12.3 Countdown

Default countdown:

3 seconds.

During Countdown:

- No gameplay selections are accepted.
- New participants cannot join the Match.
- Match settings are locked.
- The server may cancel Countdown if required conditions fail.

## 12.4 In Progress

During gameplay:

- The server processes authorized actions.
- The server controls timers.
- The server calculates scores.
- The server validates card rules.
- The server publishes authoritative state changes.

## 12.5 Completion

A Match transitions to COMPLETED when the documented
end condition is satisfied.

After completion:

- Gameplay input is rejected.
- Scores are finalized.
- Results are published.
- No further score mutations are allowed.

## 12.6 Cancellation

A Match may be cancelled for unrecoverable
server-side consistency errors.

A cancelled Match has no winner.

The system must not invent a winner from corrupted state.

## 12.7 Results

The Results screen must display:

- Final scores
- Player names
- Rankings
- Winners
- Tie status when applicable
- Forfeited players
- Rematch action

Rankings must be calculated by the server.

For equal scores, use shared placement.

Example:

1st, 1st, 3rd.

## 12.8 Rematch

A Rematch returns players to Lobby.

The next Match must use a new Match ID.

Reset:

- Scores
- Combo
- Boards
- Reference Cards
- Match Timer
- Event Schedule
- Match-specific state

Preserve:

- Room identity
- Room Code
- Connected Room members
- Current Host
- Selected settings unless changed

All players must Ready again.

---

# 13. Disconnection & Reconnection

## 13.1 Unexpected Disconnection

Unexpected disconnection includes:

- Network interruption
- Browser refresh
- Browser process termination
- Temporary connection loss
- Mobile connectivity interruption

The server must not treat every unexpected disconnection
as an immediate voluntary forfeit.

Session expiration is distinct from unexpected network disconnection. If a
Guest Session expires during an active Match, revoke its authorization and
forfeit that participant immediately; do not apply the 15-second network
reconnection grace period. W008 owns expiry and authorization revocation;
W013 applies the active-Match forfeit through the authoritative Room event
loop. A newly created session does not regain the expired session's identity
or Match authority.

## 13.2 Grace Period

Initial reconnection grace period:

15 seconds.

The Match continues for other players.

The disconnected player's state is temporarily retained.

## 13.3 Successful Reconnection

If the player reconnects with an authorized session
within the grace period:

- Restore their current server-side game state.
- Preserve their score.
- Preserve their reference and board.
- Preserve relevant cooldown state.
- Continue the Match.

Match time is not refunded.

## 13.4 Grace Period Expiration

If the player fails to reconnect before the grace
period expires while the Match remains active:

The player forfeits.

## 13.5 Match End During Grace Period

If the normal Match deadline occurs before the
disconnection grace period expires:

The player remains eligible for the result under
the current Match-end policy.

Use their server-recorded score.

Do not delay normal Match completion solely to wait
for the remainder of the grace period.

## 13.6 Voluntary Leave

If the player explicitly leaves during gameplay:

- Forfeit immediately.
- Do not grant a reconnection grace period.
- Do not allow re-entry into the same active Match.

## 13.7 Reconnection Authorization

A reconnecting client must prove ownership
of the original guest session.

Username alone is insufficient.

## 13.8 Connection Replacement

A newly authorized active connection must invalidate
gameplay authority of the older connection.

## 13.9 Server Restart Limitation

The MVP stores active Match state in memory.

A server restart may destroy active Rooms and Matches.

Full recovery after server process loss is not required
for the MVP.

The UI should report a clear recoverable error rather
than displaying stale gameplay indefinitely.

---

# 14. Anti-Griefing & Anti-Cheat

## 14.1 General Policy

No player should be able to unfairly manipulate
another player's Match result through ordinary controls.

The system must prevent:

- Unauthorized score changes
- Duplicate score updates
- Unauthorized gameplay actions
- Invalid card selection acceptance
- Host abuse during Matches
- Unlimited action spam
- Room overcapacity
- Session impersonation
- Resource exhaustion from unbounded queues

## 14.2 Server Authority

The server owns:

- Card generation
- Card validation
- Scores
- Combo
- Timers
- Match outcomes
- Event schedules
- Room membership
- Player permissions

Clients submit intentions, not authoritative results.

## 14.3 Request Idempotency

Every gameplay action must have an identifier.

Repeated requests representing the same logical action
must not modify game state more than once.

## 14.4 Rate Limiting

The system must use rate limiting for:

- Room creation
- Room joining
- Card selection
- No-Match declarations
- Username changes
- Repeated malformed messages

Initial suggested limits:

| Operation | Initial Limit |
|---|---|
| Card Pick | 12 requests/second/player |
| Guest Session Creation | At most 10 new sessions/10 minutes/source IP |
| Create Room | 5 attempts/10 minutes/Guest Session |
| Join Room | 20 attempts/minute/session |

The session creation, Room creation, and Room joining limits are
Project Owner-approved initial configurable baselines. They must be enforced
server-side and reviewed during integration and abuse testing. Card Pick and
other operation limits remain tuning values subject to testing.

Protection must not rely exclusively on guest session IDs,
because users can create new sessions.

Network-source-based protection may be used carefully,
without relying on IP address as a permanent identity.

Reuse an existing valid Guest Session on page refresh. Network-source
controls may mitigate repeated session creation, but must account for shared
networks and must not treat source IP as permanent player identity.

## 14.5 Gameplay Cooldown

Gameplay penalties and security rate limiting
must remain separate mechanisms.

Cooldown violations must not repeatedly extend
the cooldown timer.

## 14.6 Host Abuse

The Host must not:

- Kick players during active Matches.
- Edit scores.
- Force changes to Match rules.
- Remove losing or winning players.
- Cancel an active Match to avoid defeat.

## 14.7 Room Abuse

Provide:

- Room capacity limits
- Lobby-only Kick
- Room Lock
- Room Code regeneration
- Room expiration
- Join rate limiting

The system cannot guarantee permanent bans without
verified player identity.

## 14.8 Bot Limitation

Server-authoritative validation prevents direct
score forgery but does not fully prevent automated
programs from analyzing visible card data.

The MVP is casual and private-room focused.

Do not claim complete bot prevention.

## 14.9 Security Logging

Security-related events may be logged for diagnostics.

Do not log:

- Raw session secrets
- Authentication cookies
- Sensitive credentials
- Unnecessary personal information

Logs must have reasonable retention and cleanup policies.

---

# 15. Player Action & Edge Case Matrix

The implementation must handle the following scenarios.

| ID | Scenario | Expected Behavior |
|---|---|---|
| EC-001 | Username empty | Reject |
| EC-002 | Username invalid | Display validation error |
| EC-003 | Duplicate username in room | Reject or request change |
| EC-004 | Invalid Room Code | Show room not found |
| EC-005 | Expired Room Code | Show room unavailable |
| EC-006 | Room full | Reject new participant |
| EC-007 | Room locked | Reject new participant |
| EC-008 | Join during Match | Reject new participant |
| EC-009 | Authorized reconnect | Restore session |
| EC-010 | Host leaves Lobby | Transfer Host |
| EC-011 | Host disconnects | Apply grace period |
| EC-012 | Host changes settings | Reset Ready |
| EC-013 | Start while not all Ready | Reject |
| EC-014 | Start request repeated | Start only one Match |
| EC-015 | Player leaves Countdown | Revalidate and cancel Countdown |
| EC-016 | Correct card pick | Accept and score |
| EC-017 | Incorrect card pick | Apply mistake penalty |
| EC-018 | Pick during cooldown | Reject without extending cooldown |
| EC-019 | Duplicate pick request | No duplicate score |
| EC-020 | Pick unknown Card ID | Reject |
| EC-021 | Pick removed card | Technical rejection |
| EC-022 | Pick using stale state | Resynchronize or validate safely |
| EC-023 | Pick after Match end | Reject |
| EC-024 | Pick before Match start | Reject |
| EC-025 | Pick during event transition | Reject without mistake penalty |
| EC-026 | Same shared card picked concurrently | One accepted selection |
| EC-027 | Correct No-Match | Finish player's Wave |
| EC-028 | Incorrect No-Match | Apply defined penalty |
| EC-029 | Excessive No-Match declarations | Enforce limit |
| EC-030 | AFK player | Match continues |
| EC-031 | Unexpected disconnect | Start grace period |
| EC-032 | Reconnect within grace | Resume |
| EC-033 | Grace expires | Forfeit |
| EC-034 | Voluntary leave | Immediate forfeit |
| EC-035 | All players forfeit | No winner |
| EC-036 | Only one eligible player remains | End with winner |
| EC-037 | Timer expires | Finalize results |
| EC-038 | Scores tied | Joint winners |
| EC-039 | Invalid Board state | Recover or cancel safely |
| EC-040 | Server process loses state | Show recovery error |
| EC-041 | Same guest opens second tab | Single active connection |
| EC-042 | Client modifies score payload | No authoritative effect |
| EC-043 | Client modifies clock | No authoritative effect |
| EC-044 | Client sends malformed JSON | Reject safely |
| EC-045 | Client sends oversized payload | Reject or disconnect |
| EC-046 | Host attempts Kick during Match | Reject |
| EC-047 | Host attempts settings change during Match | Reject |
| EC-048 | Party Event starts | Apply server schedule |
| EC-049 | Event ends | Restore normal rule |
| EC-050 | Rematch | Create new Match and reset state |

The Agent may add test cases for uncovered behavior,
but may not silently change the expected outcomes above.

---

# 16. Functional Requirements

This section consolidates the functional requirements.

## 16.1 Guest Management

| ID | Requirement | Priority |
|---|---|---|
| FR-001 | Enter username before playing | P0 |
| FR-002 | Validate username | P0 |
| FR-003 | Create guest session | P0 |
| FR-004 | Remember display username locally | P0 |
| FR-005 | Change username in Lobby | P1 |
| FR-006 | Prevent duplicate usernames in the same Room | P0 |
| FR-007 | Recover authorized session after Refresh | P0 |
| FR-008 | Enforce one active connection per RoomPlayer | P0 |

## 16.2 Room Management

| ID | Requirement | Priority |
|---|---|---|
| FR-009 | Create private room | P0 |
| FR-010 | Generate unique active Room Code | P0 |
| FR-011 | Join by Room Code | P0 |
| FR-012 | Join by Invitation URL | P0 |
| FR-013 | Enforce Room capacity | P0 |
| FR-014 | Display players and Ready states | P0 |
| FR-015 | Toggle Ready state | P0 |
| FR-016 | Host selects game mode | P0 |
| FR-017 | Host configures permitted Match settings | P0 |
| FR-018 | Host starts Match | P0 |
| FR-019 | Host kicks player in Lobby | P0 |
| FR-020 | Host locks Room | P0 |
| FR-021 | Host regenerates Room Code | P1 |
| FR-022 | Automatic Host migration | P0 |
| FR-023 | Room expiration and cleanup | P0 |

## 16.3 Card & Game Engine

| ID | Requirement | Priority |
|---|---|---|
| FR-024 | Support five binary card attributes | P0 |
| FR-025 | Support 32 character variations | P0 |
| FR-026 | Assign Card Instance IDs | P0 |
| FR-027 | Generate Reference Cards | P0 |
| FR-028 | Generate game boards | P0 |
| FR-029 | Calculate attribute differences | P0 |
| FR-030 | Accept cards differing by 0–1 attributes | P0 |
| FR-031 | Reject cards differing by 2+ attributes | P0 |
| FR-032 | Update Reference after valid pick | P0 |
| FR-033 | Preserve Reference after invalid pick | P0 |
| FR-034 | Calculate authoritative scores | P0 |
| FR-035 | Manage Combo state | P0 |

## 16.4 Speed Battle

| ID | Requirement | Priority |
|---|---|---|
| FR-036 | Independent board per player | P0 |
| FR-037 | Identical starting conditions | P0 |
| FR-038 | Server validates selections | P0 |
| FR-039 | Refill board after valid pick | P0 |
| FR-040 | Guarantee an available valid selection | P0 |
| FR-041 | Apply mistake cooldown | P0 |
| FR-042 | Synchronize player scores | P0 |
| FR-043 | End Match by configured conditions | P0 |
| FR-044 | Calculate winners and rankings | P0 |

## 16.5 Classic Rush

| ID | Requirement | Priority |
|---|---|---|
| FR-045 | Shared board | P1 |
| FR-046 | Remove claimed card for all players | P1 |
| FR-047 | Prevent duplicate card claims | P1 |
| FR-048 | Support multiple Waves | P1 |
| FR-049 | No card refill during Wave | P1 |
| FR-050 | Support No-Match declaration | P1 |
| FR-051 | Validate No-Match | P1 |
| FR-052 | Penalize incorrect declaration | P1 |
| FR-053 | End Waves correctly | P1 |
| FR-054 | Accumulate Wave scores | P1 |

## 16.6 Party Chaos

| ID | Requirement | Priority |
|---|---|---|
| FR-055 | Use private player boards | P1 |
| FR-056 | Server-controlled Event Schedule | P1 |
| FR-057 | Consistent event rules for players | P1 |
| FR-058 | Basket Shuffle | P1 |
| FR-059 | Berry Boost | P1 |
| FR-060 | Coco Combo | P1 |
| FR-061 | Event warning | P1 |
| FR-062 | Safe event transitions | P1 |
| FR-063 | Authoritative event scoring | P1 |

## 16.7 Multiplayer & Lifecycle

| ID | Requirement | Priority |
|---|---|---|
| FR-064 | WebSocket state synchronization | P0 |
| FR-065 | Server-authoritative game state | P0 |
| FR-066 | Server-owned Match timing | P0 |
| FR-067 | Complete Match lifecycle | P0 |
| FR-068 | Reconnection support | P0 |
| FR-069 | Forfeit handling | P0 |
| FR-070 | Non-blocking AFK behavior | P0 |
| FR-071 | Duplicate action prevention | P0 |
| FR-072 | Race condition protection | P0 |
| FR-073 | Stale action handling | P0 |
| FR-074 | Authoritative Match Results | P0 |
| FR-075 | Rematch | P0 |
| FR-076 | Separate gameplay errors from technical rejections | P0 |

---

# 17. Critical Game Invariants

The system must preserve the following conditions.

| ID | Invariant |
|---|---|
| INV-001 | Maximum 6 Room participants |
| INV-002 | Match requires at least 2 Ready participants |
| INV-003 | All participants must be Ready before Match start |
| INV-004 | Guest Session cannot control multiple active Room connections |
| INV-005 | Invalid Match cannot change Reference Card |
| INV-006 | One accepted logical action cannot score twice |
| INV-007 | Scores cannot fall below zero |
| INV-008 | Private playable board always contains a valid choice |
| INV-009 | Shared Card Instance can be claimed once |
| INV-010 | Completed Match cannot accept gameplay actions |
| INV-011 | Host cannot modify Match results |
| INV-012 | Technical rejection does not cause mistake penalty |
| INV-013 | Event timing is server-authoritative |
| INV-014 | Duplicate requests do not duplicate mutations |
| INV-015 | Every Match has a bounded termination path |
| INV-016 | Forfeited player cannot resume the same Match |
| INV-017 | Reconnection requires authorized session ownership |
| INV-018 | Final Results match authoritative server state |

These invariants must be covered by automated tests.

---

# 18. Technical Architecture Requirements

## 18.1 Architecture Style

Selected baseline:

**Modular Monolith**

A single Go Game Server will manage multiple Rooms.

Modules must remain logically separated.

Microservices are not required for MVP.

## 18.2 Frontend Stack

Selected technologies:

- React
- TypeScript
- Phaser 4
- Tailwind CSS
- Vite
- pnpm

The Agent must verify compatible dependency versions
during Project Bootstrap.

Do not automatically upgrade major versions without review.

## 18.3 Frontend Responsibilities

React manages:

- Username
- Main Menu
- Lobby
- Forms
- Room information
- Scoreboard
- Match timer display
- Results
- Settings
- Responsive layout

Phaser manages:

- Card rendering
- Card interaction
- Gameplay animations
- Visual effects
- Game-related audio effects

## 18.4 React–Phaser Integration

Use a controlled interface or Game Bridge.

Phaser should receive render state and emit user intentions.

Phaser must not independently manage authoritative
scores or Match results.

## 18.5 Backend Stack

Selected:

- Go
- Standard `net/http`
- WebSocket
- In-memory Match State

Preferred WebSocket library:

`github.com/coder/websocket`

The Agent must verify API compatibility and security
configuration before implementation.

## 18.6 Backend Modules

Recommended modules:

- Guest Session
- Room Management
- Match Coordinator
- Action Processor
- Card Engine
- Game Mode Engine
- Score Engine
- WebSocket Gateway
- Security
- Configuration
- State Management

## 18.7 Authoritative Game Engine

Core gameplay rules must exist on the server.

The game engine should not depend directly on:

- React
- Phaser
- Browser APIs
- UI rendering
- HTTP transport details

Core rules must be independently testable.

## 18.8 Concurrency Model

Recommended baseline:

**Room-owned Event Loop**

Each Room processes authoritative state mutations
in a deterministic sequence.

Different Rooms may operate concurrently.

Shared board actions must not create data races.

## 18.9 State Storage

MVP storage:

In-memory.

No database is initially required.

A PostgreSQL database may be added for future features
requiring persistence.

## 18.10 Single-instance Deployment

MVP assumes one authoritative Game Server instance.

Do not deploy multiple independent game server instances
without addressing Room ownership and state routing.

## 18.11 Server State Limitations

Active Match state is not durable in MVP.

Server crash recovery is not guaranteed.

This limitation must be documented.

---

# 19. Networking & Protocol Requirements

## 19.1 Communication

Use:

- HTTPS for Web Application and HTTP APIs
- WSS for production WebSocket connections

## 19.2 API Responsibilities

HTTP may be used for:

- Guest Session initialization
- Health checks
- Application bootstrap

WebSocket is used for:

- Room actions
- Ready status
- Match start
- Card selection
- Game events
- Score updates
- Match results
- State synchronization

## 19.3 Protocol Format

Initial protocol:

Versioned JSON messages.

Example:

```json
{
  "v": 1,
  "type": "game.pick_card",
  "requestId": "req-1042",
  "payload": {
    "matchId": "match-001",
    "cardId": "card-012",
    "boardVersion": 18,
    "referenceVersion": 7
  }
}
```

## 19.4 Server Validation

The server must validate:

- Session
- Room membership
- Match membership
- Player eligibility
- Match state
- Request ID
- Action type
- Payload schema
- Card availability
- Reference state
- Applicable cooldown
- Matching rule
- Applicable event state

## 19.5 Error Codes

Suggested error codes:

- INVALID_MATCH
- CARD_UNAVAILABLE
- STALE_ACTION
- ACTION_DUPLICATE
- ACTION_COOLDOWN
- MATCH_NOT_ACTIVE
- PLAYER_FORFEITED
- RATE_LIMITED
- INVALID_SESSION
- ROOM_FULL
- ROOM_LOCKED
- ROOM_NOT_FOUND
- NOT_AUTHORIZED
- INVALID_PAYLOAD

Errors must have predictable handling.

## 19.6 State Synchronization

Use:

**Authoritative Snapshots + Incremental Events**

On connection or reconnection:

Send an authorized snapshot.

During gameplay:

Send relevant state updates.

If synchronization cannot be maintained safely:

Request or send a fresh snapshot.

## 19.7 Private State

Do not expose another player's private board
unless specifically required by an approved feature.

Public score information may be broadcast to the Room.

## 19.8 Action Versioning

Use board and reference versions where applicable.

Version checks must consider mode-specific ownership.

Classic Rush must not reject every action merely because
another player removed a different shared card.

The server should validate whether the selected instance
remains available and satisfies the player's reference.

## 19.9 Input Responsiveness

Immediate local visual feedback is allowed.

Authoritative scoring remains server-side.

The initial implementation may allow one unresolved
pick action per player to simplify synchronization.

More advanced bounded input queuing or prediction may be
introduced only after latency testing demonstrates a need.

Client prediction must never be treated as confirmed scoring.

## 19.10 Message Ordering

The server must process authoritative actions safely.

Use sequence numbers or equivalent synchronization mechanisms.

Client actions from previous Matches or Waves must not
modify a newer Match or Wave.

## 19.11 Time Authority

Server timing controls:

- Match start
- Match end
- Wave duration
- Cooldown
- Event activation
- Reconnection deadline

Client clock changes must have no authoritative effect.

## 19.12 Slow Client Protection

A slow or disconnected client must not block
the Room's authoritative event processing.

Use bounded queues and explicit slow-client handling.

---

# 20. Security Requirements

## SEC-001 — Session Authorization

Every protected action must be associated with
an authorized Guest Session.

## SEC-002 — Secure Session Credentials

Session credentials must be unpredictable.

Do not expose session secrets in logs or UI.

## SEC-003 — WebSocket Origin Validation

Validate allowed Origins for WebSocket connections.

Do not rely on cookies alone.

## SEC-004 — Input Validation

Validate untrusted input on the server.

Includes:

- Username
- Room Code
- Message type
- Message payload
- Card ID
- Match ID
- Request ID

## SEC-005 — Message Size Limits

Enforce maximum message sizes.

Oversized messages must not exhaust server resources.

## SEC-006 — Rate Limiting

Apply appropriate rate limits to sensitive operations.

## SEC-007 — Authorization

Host-only operations must verify Host authority.

Gameplay operations must verify Match participation.

## SEC-008 — Duplicate Action Protection

Repeated logical actions must not duplicate state mutations.

## SEC-009 — Race Condition Protection

Concurrent actions must not corrupt Room or Match state.

## SEC-010 — XSS Protection

Usernames and other user-supplied content must be
rendered safely.

Do not inject untrusted content as raw HTML.

## SEC-011 — Dependency Security

Dependencies must be reviewed and kept reasonably updated.

## SEC-012 — Skill Supply-chain Security

Third-party skills must be reviewed before installation
or execution.

## SEC-013 — Least Privilege

Agents and runtime components should receive only
the permissions required for their tasks.

## SEC-014 — Safe Logging

Sensitive credentials must not be logged.

## SEC-015 — Resource Limits

Bound:

- Connection counts
- Room counts where appropriate
- Message sizes
- Action queues
- Outgoing queues
- Session lifetimes
- Room lifetimes

---

# 21. Non-functional Requirements

| ID | Category | Requirement |
|---|---|---|
| NFR-001 | Performance | Gameplay interactions should feel responsive |
| NFR-002 | Performance | Room processing must not unnecessarily block other Rooms |
| NFR-003 | Capacity | Support 6 players per Room |
| NFR-004 | Reliability | Connection failures must not crash the server |
| NFR-005 | Reliability | Support authorized reconnection |
| NFR-006 | Reliability | Duplicate actions must not duplicate scores |
| NFR-007 | Security | Server validates gameplay actions |
| NFR-008 | Security | Client score and timer values are untrusted |
| NFR-009 | Security | Session credentials are unpredictable |
| NFR-010 | Security | Rate limiting is required |
| NFR-011 | Security | WebSocket Origin must be validated |
| NFR-012 | Security | User-generated text must be safely rendered |
| NFR-013 | Usability | Support mouse and touch |
| NFR-014 | Usability | Gameplay must not require hover |
| NFR-015 | Usability | All five card attributes must be distinguishable |
| NFR-016 | Usability | Mobile gameplay should not require board scrolling |
| NFR-017 | Accessibility | Readable text and sufficient contrast |
| NFR-018 | Accessibility | Reduced-motion support |
| NFR-019 | Localization | Correct Thai Unicode handling |
| NFR-020 | Maintainability | Game rules separated from presentation |
| NFR-021 | Maintainability | Game modes reuse shared card rules |
| NFR-022 | Maintainability | Tuning values are configurable |
| NFR-023 | Testability | Core engine supports automated testing |
| NFR-024 | Privacy | Minimize personal data collection |
| NFR-025 | Observability | Errors and system events are diagnosable |
| NFR-026 | Deployment | Reproducible build and deployment |
| NFR-027 | Compatibility | Support target browsers |
| NFR-028 | Scalability | Permit future scaling without unnecessary initial complexity |
| NFR-029 | Accessibility | Essential actions must be accessible without precise hover interactions |
| NFR-030 | Documentation | Implementation remains traceable to approved requirements |

## 21.1 Initial Performance Targets

The following are test targets, not guaranteed outcomes.

| Metric | Initial Target |
|---|---|
| Frame Rate | Aim for 60 FPS on capable devices |
| Players per Room | 6 |
| Board Size | 16 cards |
| Reconnection Grace | 15 seconds |
| Mobile Width Testing | Starting at 320 CSS px |
| Touch Target | Approximately 44×44 CSS px where practical |
| Network Testing | 0, 100, 200, 300ms latency |

Record benchmark results before making performance claims.

---

# 22. Responsive UI/UX Requirements

## 22.1 Mobile-first

Gameplay must be designed for smartphone screens.

Desktop layouts may use additional space.

## 22.2 Card Interaction

Primary gameplay action:

Single tap or mouse click.

Do not require drag-and-drop.

## 22.3 Gameplay Board Visibility

The intended board size is 4×4.

Players should be able to view the playable board
without scrolling.

If 16 cards are not sufficiently readable on a
small smartphone, the design must be revised through
playtesting.

Do not silently provide different gameplay card counts
to competing players based solely on device type.

## 22.4 Animation

Animations must not change card attributes.

Important attributes should remain visually stable.

## 22.5 Reduced Motion

Support reduced-motion preferences where practical.

Reduced motion must not create an unfair gameplay advantage.

## 22.6 Thai UI

Examples of primary UI text:

- สร้างห้อง
- เข้าร่วมห้อง
- พร้อมแล้ว!
- เริ่มเกม!
- ถูกต้อง!
- ลองใหม่อีกครั้ง!
- หมดเวลา!
- เล่นอีกครั้ง
- ออกจากห้อง
- ไม่มีใบที่ตรงแล้ว!

## 22.7 Audio

Audio should reinforce gameplay without distracting players.

Provide:

- Background music
- Correct pick SFX
- Incorrect pick SFX
- Countdown SFX
- Victory SFX
- Basic volume controls
- Mute controls

Browser autoplay restrictions must be respected.

Audio playback should begin only after an appropriate
user interaction when required.

---

# 23. Agent-first SDD Requirements

## 23.1 Development Methodology

The project MUST follow:

**Agent-first Spec-Driven Development**

AI Coding Agents are the primary implementation workforce.

However, implementation must remain controlled by:

- Requirements
- Approved specifications
- Technical design
- Task plans
- Automated tests
- Independent reviews

## 23.2 Initial Documentation

The Project Owner provides:

`docs/REQUIREMENT.md`

as the initial authoritative source.

The Agent is responsible for organizing additional
project documentation during bootstrap.

## 23.3 Documentation Location

All authoritative project documentation must be located
under:

`/docs`

An optional root-level `AGENTS.md` may act as a small
entry point for coding agents.

It should reference documents under `/docs` rather than
duplicate the entire requirement set.

## 23.4 Recommended Documentation Structure

After bootstrap, the Agent may create:

```text
docs/
├── REQUIREMENT.md
├── product/
├── specifications/
├── architecture/
├── design/
├── standards/
├── testing/
├── development/
│   └── reviews/
└── agents/
```

The Agent must determine the exact structure during
Documentation Foundation.

Do not create unnecessary files merely to fill directories.

## 23.5 SDD Workflow

Every implementation work item must follow:

Requirements
→ Specification
→ Technical Design
→ Task Breakdown
→ Implementation
→ Developer Verification
→ Independent Sub-Agent Review
→ Remediation if required
→ Final Validation
→ Completion

## 23.6 Feature Specifications

Before feature implementation, the Agent must define:

- Feature goals
- Scope
- Non-goals
- Relevant requirement IDs
- Expected behavior
- Edge cases
- Acceptance criteria
- Technical design
- Test plan
- Implementation tasks

## 23.7 Documentation Traceability

Implementation tasks must reference relevant requirement IDs.

Tests should reference applicable requirements or
acceptance criteria.

Material design decisions should be documented through ADRs.

## 23.8 No Silent Scope Expansion

Agents must not:

- Add unrelated features
- Change game rules without approval
- Replace selected core technologies without justification
- Introduce unnecessary infrastructure
- Rewrite requirements to match defective code

## 23.9 Documentation Synchronization

Approved implementation changes must be reflected
in applicable documentation.

Documentation must not contradict accepted implementation.

---

# 24. Agent Skills Management

## 24.1 Skills Policy

The project must use appropriate Agent Skills to improve:

- Code quality
- Language best practices
- Architecture
- Testing
- Security
- UX/UI
- Documentation
- DevOps

## 24.2 Mandatory Initial Skill

Use:

**find-skills**

Official source:

https://github.com/vercel-labs/skills

The Project Bootstrap Agent must inspect the source
and installation behavior before installing it.

## 24.3 Project-local Installation Only

ALL skills must be installed within the Repository.

Global installation is prohibited.

User-level installation is prohibited.

Preferred canonical location:

`.agents/skills/`

Agent-specific project-local paths may be used when required
by the selected coding agent.

## 24.4 Installation Restrictions

DO NOT:

- Use `-g`
- Use `--global`
- Install into user home directories
- Modify global agent settings
- Install unreviewed third-party skills
- Automatically update skills without review

## 24.5 Initial Skill Discovery Procedure

During Project Bootstrap:

1. Inspect the official `find-skills` source.
2. Verify project-local installation support.
3. Install `find-skills` at project scope only.
4. Verify the actual installation path.
5. Search for skills relevant to the technology stack.
6. Evaluate candidate skills.
7. Approve trusted candidates.
8. Install approved skills locally.
9. Record versions and source repositories.
10. Validate installation and compatibility.

Example project-local discovery:

```bash
npx skills find "typescript best practices"
npx skills find "golang"
npx skills find "react"
npx skills find "web design"
npx skills find "playwright testing"
npx skills find "web security"
npx skills find "code review"
```

Do not run installation commands until their effects
and scope have been verified.

## 24.6 Skill Categories

Initial discovery should cover:

| Category | Focus |
|---|---|
| TypeScript | Language best practices |
| Go | Language best practices |
| React | Component quality and performance |
| Phaser | Game development |
| UX/UI | Visual and interaction design |
| Accessibility | Responsive and accessible UI |
| Testing | Unit, Integration, E2E |
| Security | Web, WebSocket, Session |
| Code Quality | Review and refactoring |
| Architecture | Design patterns and maintainability |
| DevOps | Docker and GitHub Actions |
| Documentation | SDD and technical writing |

## 24.7 Skill Trust Evaluation

Every candidate must be evaluated for:

- Maintainer reputation
- Source authenticity
- Community adoption
- Maintenance activity
- Documentation quality
- Technology relevance
- Security risks
- License
- Agent compatibility
- Conflicts with existing standards

Popularity is a useful signal, not proof of safety.

Prefer official or established maintainers.

## 24.8 Skill Approval

A Skill must not be treated as approved merely because
it appears in search results.

Before approval, inspect:

- SKILL.md
- Bundled scripts
- Referenced commands
- External dependencies
- Network behavior
- Filesystem operations
- Required permissions

## 24.9 Skill Registry

The Agent must maintain a Skill Registry under `/docs`.

For each approved Skill, record:

- Name
- Source
- Maintainer
- Purpose
- Version or commit
- Installation path
- Security review status
- Approval status
- Review date

## 24.10 Skill Updates

Do not update skills automatically.

Updates require:

1. Change inspection
2. Security review
3. Compatibility validation
4. Registry update
5. Relevant tests

## 24.11 Missing Skills

If no trusted Skill exists for a required domain:

Use official documentation and established
project standards instead.

Do not install low-quality skills merely to satisfy
a category checklist.

---

# 25. Independent Sub-Agent Review

## 25.1 Mandatory Review

Every completed development iteration or Work Item
MUST undergo independent review before being marked DONE.

The implementing agent cannot independently certify
its own implementation.

## 25.2 Reviewer Independence

Preferred:

A separate Sub-Agent provided by the coding environment.

Alternative:

A separate isolated Agent session.

The reviewer must independently inspect evidence.

Do not claim that a Sub-Agent review occurred
unless an independent execution actually took place.

## 25.3 Reviewer Responsibilities

The Reviewer must validate:

- Requirement compliance
- Feature specification compliance
- Technical design compliance
- Acceptance criteria
- Code quality
- Security
- Error handling
- Concurrency correctness
- Regression risks
- Test evidence
- Documentation consistency

## 25.4 Review Inputs

The Reviewer must receive access to:

- Requirement references
- Feature Specification
- Technical Design
- Task Plan
- Implementation Diff
- Relevant Source Code
- Tests
- Test Results
- Developer Implementation Report

The Developer Agent's summary is supporting information,
not sufficient evidence by itself.

## 25.5 Review Permissions

Prefer read-only reviewer operation.

The Reviewer may:

- Inspect code
- Read documentation
- Analyze tests
- Run non-destructive checks where authorized
- Produce findings

The Reviewer must not silently:

- Modify implementation
- Rewrite requirements
- Install dependencies
- Install Skills
- Merge code
- Deploy changes

## 25.6 Review Severity

### BLOCKER

Critical issue preventing acceptance.

Examples:

- Core gameplay violation
- Score manipulation vulnerability
- Critical security defect
- Data corruption
- Incorrect Match result
- Unrecoverable concurrency defect

### MAJOR

Significant issue requiring correction.

Examples:

- Missing required behavior
- Failed critical acceptance criterion
- Incorrect state transition
- Important test gap
- Material architecture violation

### MINOR

Non-critical issue.

Examples:

- Minor UX inconsistency
- Small maintainability issue
- Minor documentation discrepancy

### SUGGESTION

Optional improvement outside mandatory acceptance.

## 25.7 Review Decisions

Supported decisions:

- PASS
- FAIL
- BLOCKED

### PASS

All required acceptance criteria are satisfied.

No unresolved BLOCKER or MAJOR findings remain.

### FAIL

Blocking defects or failed required acceptance criteria exist.

### BLOCKED

Review cannot be completed due to missing evidence,
missing access, unavailable tools, or unresolved ambiguity.

BLOCKED is not PASS.

## 25.8 Remediation

If review fails:

1. Reviewer documents findings.
2. Developer investigates root causes.
3. Developer fixes the implementation.
4. Developer updates tests.
5. Developer runs affected validation.
6. Reviewer independently verifies fixes.
7. Reviewer issues updated decision.

Repeat until the review gate is satisfied.

## 25.9 Post-review Changes

Material code changes after a PASS require
additional validation and review.

A previous PASS does not automatically approve
new implementation changes.

## 25.10 Review Reports

Each development Work Item must have an independent
review record.

Recommended location:

`docs/development/reviews/WXXX-review.md`

Each review must include:

- Work Item ID
- Reviewer execution reference
- Reviewed commit or diff
- Requirements reviewed
- Acceptance criteria
- Tests inspected or executed
- Findings
- Severity
- Resolution status
- Remaining risks
- Final decision

## 25.11 No Fake Validation

Agents must never:

- Invent review results
- Invent test results
- Claim unexecuted tests passed
- Claim unavailable reviewers participated
- Suppress failing test evidence

If independent review is unavailable, report BLOCKED.

---

# 26. Definition of Done

A Work Item is DONE only when all applicable conditions
are satisfied.

## 26.1 Specification

- [ ] Relevant requirements identified
- [ ] Specification approved
- [ ] Implementation follows specification
- [ ] Scope respected
- [ ] Architecture decisions respected

## 26.2 Implementation

- [ ] Code follows project standards
- [ ] Errors handled
- [ ] Untrusted inputs validated
- [ ] No unnecessary dependencies
- [ ] No unauthorized rule changes
- [ ] No critical regressions

## 26.3 Testing

- [ ] Formatting passed
- [ ] Lint passed
- [ ] Type checking passed
- [ ] Unit tests passed
- [ ] Applicable integration tests passed
- [ ] Applicable E2E tests passed
- [ ] Acceptance criteria verified
- [ ] Security checks completed where applicable

Non-applicable checks must be marked N/A with a reason.

Unexecuted required checks must not be represented as PASS.

## 26.4 Independent Review

- [ ] Independent reviewer executed
- [ ] Code diff reviewed
- [ ] Requirements checked
- [ ] Acceptance criteria checked
- [ ] Tests inspected
- [ ] BLOCKER findings resolved
- [ ] MAJOR findings resolved
- [ ] Re-review completed when needed
- [ ] Final review decision PASS

## 26.5 Documentation

- [ ] Required documentation updated
- [ ] Requirement references remain correct
- [ ] ADRs updated if needed
- [ ] Review report recorded

## 26.6 Skill Compliance

- [ ] Approved Skills used where relevant
- [ ] New Skills reviewed
- [ ] Installed Skills are project-local
- [ ] Skill Registry updated if applicable
- [ ] No global installation performed

## 26.7 Handoff

- [ ] Implementation report complete
- [ ] Test results recorded
- [ ] Review result recorded
- [ ] Known limitations recorded
- [ ] Follow-up tasks recorded
- [ ] Final status accurate

No Agent may bypass mandatory review to accelerate delivery.

---

# 27. Testing Strategy

## 27.1 Testing Levels

Required testing categories:

- Unit Testing
- Integration Testing
- End-to-End Testing
- Concurrency Testing
- Security Testing
- Performance Testing
- Browser Compatibility Testing
- Visual and UX Validation

## 27.2 Recommended Tools

Frontend:

- Vitest
- Playwright

Backend:

- Go testing
- Go race detector
- Go fuzz testing where appropriate

CI:

- GitHub Actions

## 27.3 Core Engine Testing

The Card Matching Engine must test all possible
card variation pairs.

There are:

32 × 32 = 1024 pairs.

All pairs should be checked automatically.

## 27.4 Deterministic Testing

The Game Engine should support controllable:

- Random generation
- Time
- Match state
- Action ordering

Tests must not unnecessarily wait for real Match durations.

## 27.5 Concurrency Testing

Critical scenarios include:

- Multiple players choosing the same shared card
- Multiple Ready requests
- Duplicate Match start requests
- Reconnection during active Match
- Match ending while an action is in progress
- Event transition during card selection

## 27.6 E2E Testing

Use multiple browser contexts to represent
different players.

Important flows:

- Create Room
- Join Room
- Ready
- Start Match
- Pick Cards
- Update Scoreboard
- Match End
- Results
- Rematch
- Reconnect

## 27.7 Network Testing

Simulate:

- Network latency
- Disconnects
- Reconnects
- Delayed messages
- Duplicate messages
- Slow clients

Network-related failures must not violate
server-side game invariants.

## 27.8 Visual Testing

Validate:

- Thai text rendering
- Mobile layout
- Card readability
- Character consistency
- Attribute visibility
- Touch usability
- Responsive behavior

## 27.9 Security Testing

Validate:

- Unauthorized actions
- Session misuse
- Invalid messages
- Oversized payloads
- Rate limiting
- Duplicate requests
- WebSocket Origin handling
- Host permission boundaries

---

# 28. Acceptance Criteria

The following are minimum acceptance scenarios.

## AC-001 — Guest Entry

Given a first-time visitor,

When the visitor enters a valid username,

Then the system creates an authorized Guest Session
without requiring Login.

## AC-002 — Create Room

Given an authorized Guest,

When the Guest creates a Room,

Then the system creates a unique Room and assigns
the Guest as Host.

## AC-003 — Join by Invitation

Given a valid Room invitation,

When another Guest opens the link,

Then the Guest can enter the correct Lobby
subject to Room rules.

## AC-004 — Room Capacity

Given a Room containing six reserved participants,

When a seventh new Guest attempts to join,

Then the server rejects the request.

## AC-005 — Ready Validation

Given at least one participant is not Ready,

When the Host attempts to start the Match,

Then the server rejects the action.

## AC-006 — Correct Match

Given Reference Card 10110,

When the player chooses 10100,

Then the server accepts the card.

## AC-007 — Incorrect Match

Given Reference Card 10110,

When the player chooses 10000,

Then the server rejects the card
and applies the mistake penalty.

## AC-008 — Duplicate Action

Given an accepted Pick request,

When the same logical request is repeated,

Then the score changes only once.

## AC-009 — Speed Battle Independence

Given two players with private boards,

When Player A picks a card,

Then Player B's board is not modified by A's selection.

## AC-010 — Board Solvability

Given an active private board,

After any accepted pick,

The resulting playable board contains
at least one valid card.

## AC-011 — Shared Card Collision

Given two players select the same shared Card Instance,

When the server processes both requests,

Then at most one request succeeds.

## AC-012 — Correct No-Match

Given no valid shared card remains for a player,

When the player declares No-Match,

Then the declaration is accepted.

## AC-013 — Incorrect No-Match

Given a valid shared card exists,

When the player declares No-Match,

Then the documented penalty is applied.

## AC-014 — Basket Shuffle

Given Basket Shuffle begins,

When card positions change,

Then all Card IDs, attributes, and References remain correct.

## AC-015 — Berry Boost

Given Berry Boost is active,

When the server accepts a correct pick,

Then the configured bonus scoring applies.

## AC-016 — Coco Combo

Given Coco Combo is active,

When a player satisfies the required sequence,

Then the bonus is awarded at most once for that event.

## AC-017 — Reconnection

Given a player disconnects unexpectedly,

When the same authorized session reconnects
within the grace period,

Then the player resumes their current server-side state.

## AC-018 — Forfeit

Given a player voluntarily leaves an active Match,

Then the player forfeits immediately.

## AC-019 — Host Restriction

Given a Match is active,

When the Host attempts to Kick a competitor,

Then the server rejects the action.

## AC-020 — Match Timer

Given the authoritative Match deadline is reached,

Then the server finalizes the Match
and rejects further scoring actions.

## AC-021 — Tie

Given two eligible players have the highest equal score,

Then both players are declared joint winners.

## AC-022 — Rematch

Given a completed Match,

When players return to Lobby and start another Match,

Then the new Match uses a new Match ID
and reset Match state.

## AC-023 — Multiple Tabs

Given one Guest Session opens the Room in two tabs,

Then only one connection retains gameplay authority.

## AC-024 — Invalid Client Score

Given a client submits an unauthorized score value,

Then the server does not apply it.

## AC-025 — Independent Review

Given a development Work Item is implementation-complete,

When independent review has not passed,

Then the Work Item cannot be marked DONE.

## AC-026 — Project-local Skills

Given an Agent installs a Skill,

Then the installation must remain within the Repository
and match the approved Skill Registry.

---

# 29. CI/CD & Quality Gates

## 29.1 Continuous Integration

Use GitHub Actions.

Pull Request validation should include:

- Formatting
- Linting
- Type checking
- Unit tests
- Applicable integration tests
- Build verification
- Security checks where appropriate

## 29.2 Go Race Detection

Use the Go race detector for concurrency-sensitive code.

Example:

```bash
go test -race ./...
```

## 29.3 Review Gate

Code should not be considered ready for completion
without independent validation.

A Pull Request review by a human or platform reviewer
may supplement but does not automatically replace
the required independent Agent review.

## 29.4 Release Authorization

Agents must not independently publish releases
or deploy production changes without authorization
from the applicable project workflow.

---

# 30. Deployment Requirements

## 30.1 Deployment Model

Initial deployment:

Single-region, single authoritative Game Server.

## 30.2 Frontend Deployment

The frontend should support static hosting.

React and Phaser assets should be built
into deployable web assets.

## 30.3 Backend Deployment

The Go Game Server must support:

- HTTP
- WebSocket
- Graceful shutdown
- Health checks
- Resource limits
- Environment-based configuration

## 30.4 Same-origin Production Routing

Preferred production routing:

```text
https://example.com/
https://example.com/room/:code
https://example.com/api/...
wss://example.com/ws
```

Use a reverse proxy or equivalent configuration.

## 30.5 Containerization

Docker may be used for reproducible backend deployment.

Do not introduce container orchestration without
a demonstrated requirement.

## 30.6 Initial Infrastructure

Database:

Not required.

Redis:

Not required.

Microservices:

Not required.

Kubernetes:

Not required.

## 30.7 Future Scaling

Future scaling may require:

- Room ownership routing
- Distributed coordination
- Persistent storage
- Multi-instance infrastructure

These are not MVP requirements.

---

# 31. Suggested Repository Architecture

The following structure is a technical baseline.

The Agent may refine it during Project Bootstrap,
provided approved architectural principles are preserved.

```text
prickle-panic/
│
├── docs/
│   ├── REQUIREMENT.md
│   └── ...
│
├── .agents/
│   └── skills/
│
├── apps/
│   ├── web/
│   │   ├── src/
│   │   ├── public/
│   │   └── tests/
│   │
│   └── server/
│       ├── cmd/
│       ├── internal/
│       └── go.mod
│
├── packages/
│   └── protocol/
│
├── assets/
│   ├── source/
│   └── generated/
│
├── tests/
│   ├── e2e/
│   ├── integration/
│   └── load/
│
├── infra/
│
├── .github/
│   └── workflows/
│
├── AGENTS.md
├── package.json
├── pnpm-workspace.yaml
└── README.md
```

The Agent must avoid creating empty architecture
layers without an implementation need.

---

# 32. Development Roadmap

## Phase 00 — Foundation & Agent Setup

### W001 — Repository Bootstrap

- Initialize repository structure.
- Configure frontend and backend workspaces.
- Establish development commands.
- Verify local toolchain.
- Establish initial documentation entry points.

### W002 — SDD Documentation Foundation

- Analyze this REQUIREMENT.md.
- Generate detailed specifications as needed.
- Establish SDD workflow.
- Establish document templates.
- Define requirement traceability.
- Define ADR workflow.

### W003 — Skill Discovery & Approval

- Verify find-skills.
- Install find-skills locally.
- Discover relevant Skills.
- Evaluate candidate sources.
- Review security.
- Install approved Skills locally.
- Create Skill Registry.
- Verify installation paths.

### W004 — Development Standards & CI

- Configure formatting.
- Configure linting.
- Configure testing.
- Configure CI.
- Define quality gates.

## Phase 01 — Core Game Engine

### W005 — Protocol Contracts

- Define messages.
- Define schemas.
- Define error codes.
- Establish contract validation.

### W006 — Card Matching Engine

- Implement five-attribute cards.
- Implement validation.
- Implement card generation.
- Test all 1024 variation pairs.

### W007 — Local Gameplay Prototype

- Implement Speed Battle locally.
- Validate card interaction.
- Validate scoring.
- Validate basic mobile usability.

## Phase 02 — Multiplayer Foundation

### W008 — Guest Session

- Guest identity
- Session handling
- Username validation
- Session recovery

### W009 — Room Management

- Create Room
- Join Room
- Invitation URL
- Lobby
- Ready
- Host permissions

### W010 — WebSocket Foundation

- Connection lifecycle
- Authorized actions
- State synchronization
- Room concurrency
- Duplicate action handling

## Phase 03 — Online Gameplay

### W011 — Multiplayer Speed Battle

- Authoritative game state
- Real-time card selections
- Real-time scoring
- Match timer
- Multiplayer synchronization

### W012 — Results & Rematch

- Match completion
- Winner calculation
- Results screen
- Rematch workflow

### W013 — Reconnection & Anti-Griefing

- Grace period
- Session reconnection
- Host migration
- Rate limiting
- Error recovery
- Room moderation

## Phase 04 — Product Experience

### W014 — Coco Assets & Gameplay Polish

- Coco character variants
- Card assets
- UI polish
- Animation
- Audio
- Responsive improvements

## Phase 05 — MVP Validation

### W015 — Closed Alpha Readiness

- End-to-end testing
- Multiplayer testing
- Browser compatibility
- Security review
- Performance testing
- Independent acceptance review

## Phase 06 — Additional Game Modes

After MVP validation:

- Classic Rush
- Party Chaos
- Additional balancing
- Additional audiovisual polish

Each new mode must follow the same SDD and review process.

---

# 33. Work Item Execution Contract

Every Agent development Work Item must follow this process.

## Step 1 — Read

Read:

- REQUIREMENT.md
- Relevant detailed specifications
- Architecture decisions
- Applicable standards
- Relevant Skills

## Step 2 — Analyze

Identify:

- Work Item ID
- Requirements
- Scope
- Dependencies
- Risks
- Acceptance criteria

## Step 3 — Design

Produce or update the technical design.

## Step 4 — Plan

Break work into verifiable tasks.

## Step 5 — Implement

Implement only approved scope.

## Step 6 — Self-Test

Run applicable tests.

## Step 7 — Independent Review

Launch a separate Reviewer Sub-Agent or equivalent
independent Agent execution.

## Step 8 — Remediate

Resolve BLOCKER and MAJOR findings.

## Step 9 — Re-Review

Validate corrected implementation independently.

## Step 10 — Complete

Update documentation.

Record evidence.

Mark DONE only after passing the Definition of Done.

---

# 34. Required Agent Implementation Report

Every development Work Item must conclude with a report
containing:

## Work Item

- ID
- Title
- Objective
- Scope

## Implementation

- Summary
- Files changed
- Requirements addressed
- Design decisions

## Validation

- Tests executed
- Results
- Tests not executed
- Acceptance criteria results

## Independent Review

- Reviewer execution reference
- Reviewed revision
- Review decision
- Findings
- Resolved issues
- Remaining risks

## Documentation

- Documents created
- Documents updated
- Requirement traceability changes

## Final Status

One of:

- IN_PROGRESS
- READY_FOR_REVIEW
- CHANGES_REQUESTED
- REVIEW_BLOCKED
- VERIFIED
- DONE

Do not report DONE unless required review has passed.

---

# 35. MVP Definition of Success

The MVP is considered functionally successful when
the following end-to-end scenario works correctly.

## Scenario

Four friends open Prickle Panic! using different browsers
or devices.

One player creates a Room.

The Host shares the invitation URL.

Other players join without creating accounts.

Every player chooses a username.

The Lobby displays all players correctly.

All players select Ready.

The Host starts Speed Battle.

Players receive their private card boards.

Players select cards simultaneously.

The server validates selections and calculates scores.

All players receive appropriate score updates.

A player refreshes their browser and reconnects
using the original authorized session.

The Match continues for everyone else.

The server ends the Match at the configured deadline.

The Results screen shows correct rankings.

Players return to Lobby.

A Rematch begins without creating a new Room.

## Required Completion Evidence

- Functional tests pass.
- Core game invariants pass.
- Relevant integration tests pass.
- Multiplayer E2E tests pass.
- Required browser checks pass.
- Security checks have no unresolved blockers.
- Independent Reviewer approves the Work Item.
- Documentation matches implementation.

---

# 36. Known Risks & Limitations

## RISK-001 — Browser Automation

Players may use automated programs to inspect card data
and calculate valid answers.

Server authority does not fully eliminate this risk.

Mitigation:

Focus on casual private multiplayer for MVP.

## RISK-002 — Network Latency

High latency may affect interaction responsiveness.

Classic Rush is especially sensitive to contested actions.

Mitigation:

Server-side ordering, responsive visual feedback,
and explicit network testing.

## RISK-003 — Mobile Card Readability

A 4×4 board may be difficult to read on small screens.

Mitigation:

Prototype and usability testing before final acceptance.

## RISK-004 — Server Restart

In-memory Match state is lost after server restart.

Mitigation:

Accept as a documented MVP limitation.

## RISK-005 — AI-generated Asset Consistency

AI-generated illustrations may have inconsistent
proportions or attributes.

Mitigation:

Modular asset system and visual validation.

## RISK-006 — Third-party Skill Security

External Skills may contain unsafe instructions or scripts.

Mitigation:

Source review, project-local installation,
approved registry, and controlled updates.

## RISK-007 — Requirement Drift

Agents may unintentionally implement behavior
that differs from the approved design.

Mitigation:

SDD, traceability, independent review,
and mandatory change control.

## RISK-008 — Guest Identity Abuse

Without accounts, players may create new sessions
to bypass session-level restrictions.

Mitigation:

Room controls, rate limiting, and realistic
security limitations.

## RISK-009 — Scope Expansion

Developing all three modes simultaneously could
delay the playable MVP.

Mitigation:

MVP-first development, beginning with Speed Battle.

---

# 37. Decisions Requiring Playtesting

The following parameters are not permanently frozen.

| Decision | Initial Value |
|---|---|
| Speed Battle Board | 4×4 |
| Speed Battle Duration | 60 seconds |
| Classic Rush Waves | 3 |
| Classic Rush Wave Duration | 30 seconds |
| Party Chaos Duration | 75 active seconds |
| Correct Pick | +1 |
| Mistake Cooldown | 0.8 seconds |
| No-Match Penalty | -1 point |
| No-Match Cooldown | 2 seconds |
| No-Match Failure Limit | 2 per Wave |
| Reconnection Grace | 15 seconds |
| Basket Shuffle Transition | 0.4 seconds |
| Berry Boost Duration | 8 seconds |
| Coco Combo Requirement | 3 correct picks |
| Coco Combo Bonus | +2 |

Agents must implement these as configurable values
where appropriate.

Gameplay tuning must not silently change
confirmed core game rules.

---

# 38. Instructions for the Initial Coding Agent

You are the Project Bootstrap and Planning Agent
for Prickle Panic!.

Your first responsibility is NOT to implement gameplay.

You must first establish the project's SDD foundation
using this document.

## Required Initial Actions

1. Read this entire document.
2. Summarize the confirmed product requirements.
3. Identify the current architecture baseline.
4. Identify implementation risks and specification gaps.
5. Establish the repository documentation structure.
6. Create appropriate Agent instructions.
7. Establish the SDD workflow.
8. Establish independent review procedures.
9. Discover and evaluate project-local Skills.
10. Establish coding and testing standards.
11. Produce a phased development plan.
12. Define the first implementable Work Item.

## Documentation Instructions

Use this document as the initial authoritative source.

You may organize more detailed documentation under `/docs`.

Preserve requirement IDs.

Do not silently alter confirmed requirements.

Do not create contradictory rules.

Do not fragment documentation merely to create files.

## Skill Instructions

Use `find-skills` during the initial setup.

Verify source and installation scope before installation.

Install Skills only within the project.

Never install Skills globally.

Do not trust unreviewed Skill instructions.

## Review Instructions

After every development iteration:

- Run applicable validation.
- Launch an independent Reviewer Agent.
- Review against requirements and the approved plan.
- Resolve blocking findings.
- Re-review when necessary.
- Record the outcome.

Never mark the iteration DONE without the required
independent review.

## Architecture Instructions

Preserve:

- Web-first development
- React + TypeScript + Phaser
- Go Game Server
- WebSocket communication
- Server-authoritative gameplay
- Modular Monolith
- In-memory MVP Match State
- Agent-first SDD

Changes to these core decisions require documented
justification and approval.

## Implementation Priority

Do not start by implementing all three game modes.

Start with the project foundation and then
the Core Card Engine and Speed Battle MVP.

## Reporting Instructions

When completing any assigned task, report:

- What was completed
- What changed
- Which requirements were addressed
- What was tested
- What remains incomplete
- What the independent reviewer found
- Whether the work passed review
- What should happen next

Be precise and evidence-based.

---

# 39. Final Project Principles

Prickle Panic! must remain:

1. Easy to learn.
2. Fun to master.
3. Accessible through browsers.
4. Friendly for casual players.
5. Competitive without unnecessary frustration.
6. Fair under documented game rules.
7. Secure against basic client manipulation.
8. Maintainable through clear architecture.
9. Testable through automated verification.
10. Consistent through Spec-Driven Development.
11. Reviewable through independent Agent validation.
12. Extensible without premature complexity.

The Project Owner retains authority over material
product decisions.

Coding Agents implement approved specifications.

Independent Reviewer Agents verify compliance.

No Agent may replace evidence with assumptions.

---

# 40. Document Completion Status

| Area | Status |
|---|---|
| Product Vision | BASELINE |
| Platform | CONFIRMED |
| Guest Multiplayer | CONFIRMED |
| Coco Mascot | CONFIRMED |
| Visual Direction | BASELINE |
| Card Matching Rules | BASELINE |
| Speed Battle | BASELINE |
| Classic Rush | BASELINE |
| Party Chaos | BASELINE |
| Room Management | BASELINE |
| Anti-Griefing | BASELINE |
| Security Requirements | BASELINE |
| Technical Architecture | BASELINE |
| Agent-first SDD | CONFIRMED |
| Project-local Skills | CONFIRMED |
| Independent Sub-Agent Review | CONFIRMED |
| Development Roadmap | BASELINE |
| Gameplay Balance | TUNING |
| Implementation | NOT STARTED |

This document is sufficiently defined to begin
Project Bootstrap and detailed specification work.

It does not authorize arbitrary changes to the
confirmed product or gameplay requirements.

---

# 41. Git & GitHub Development Workflow

## 41.1 Overview

Prickle Panic! uses Git and GitHub for source code management,
version control, collaboration, and development traceability.

The project follows a Feature Branch Workflow.

All source code changes must follow the approved branching
strategy and Pull Request process.

This workflow is mandatory for all development agents.

It applies to:

- Source code
- Documentation
- Tests
- Configuration
- Game assets
- Infrastructure
- CI/CD
- Security-related changes

Direct changes to protected branches are prohibited.

## 41.2 Branch Strategy

The repository uses three categories of branches:

| Branch | Purpose | Lifetime |
|---|---|---|
| main | Production / Release | Permanent |
| develop | Development / Integration | Permanent |
| feature/* | Feature and Work Item development | Retained after merge |

There may be multiple feature branches.

The main and develop branches must always exist.

### BRANCH-001 — main

The `main` branch represents the current production
release source code.

Responsibilities:

- Store released game versions.
- Provide production release source.
- Maintain release history.
- Serve as the source for production deployment.
- Support version tagging.

Rules:

- Direct pushes are prohibited.
- Force pushes are prohibited.
- Branch deletion is prohibited.
- Changes must arrive through Pull Requests from develop.
- Pull Requests require successful validation.
- Independent Sub-Agent Review is mandatory.
- Production release requires appropriate authorization.

The main branch must not be used for ordinary development.

### BRANCH-002 — develop

The `develop` branch is the integration branch.

Responsibilities:

- Combine completed feature work.
- Maintain the current development version.
- Provide the base for new feature branches.
- Support integration testing.
- Prepare changes for future releases.

Rules:

- Direct pushes are prohibited.
- Force pushes are prohibited.
- Branch deletion is prohibited.
- Changes must arrive through Pull Requests from feature/*.
- Pull Requests require successful validation.
- Independent Sub-Agent Review is mandatory.

The develop branch may contain functionality that is not yet
available in production.

### BRANCH-003 — feature/*

Feature branches are used for individual development
work items.

Responsibilities:

- Implement approved features.
- Fix defects.
- Add tests.
- Update documentation.
- Perform technical improvements.
- Prepare changes for integration.

Rules:

- Create new feature branches from develop.
- Use one primary feature branch per Work Item.
- Do not directly merge into main.
- Open a Pull Request targeting develop.
- Run required tests before merge.
- Complete Independent Sub-Agent Review before merge.
- Do not force push unless explicitly authorized
  through an exceptional recovery procedure.
- Do not delete the branch after merge.

Feature branches must remain available after successful merge.

---

## 41.3 Branch Naming Convention

All feature branches must follow:

`feature/<work-item-id>-<short-description>`

Examples:

```text
feature/W001-repository-bootstrap
feature/W002-sdd-foundation
feature/W003-skill-discovery
feature/W004-ci-setup
feature/W005-protocol-contracts
feature/W006-card-engine
feature/W007-speed-battle-prototype
feature/W008-guest-session
feature/W009-room-management
feature/W010-websocket-foundation
```

Naming rules:

- Use the feature/ prefix.
- Include the Work Item ID.
- Use lowercase descriptions.
- Use hyphens between words.
- Keep names concise and meaningful.
- Avoid spaces and unnecessary special characters.

Bug fixes, documentation updates, and infrastructure work
must also use feature branches.

Do not introduce additional permanent branch categories
without approval.

---

## 41.4 Branch Creation

Before starting a Work Item:

1. Check the repository status.
2. Verify the current branch.
3. Fetch the latest remote references.
4. Ensure the local develop branch is up to date.
5. Verify there are no unrelated uncommitted changes.
6. Create a new feature branch from develop.
7. Begin work on the new branch.

Example:

```bash
git fetch origin
git switch develop
git pull --ff-only origin develop

git switch -c feature/W006-card-engine
```

The Agent must verify that the new branch was created
from the intended develop revision.

Do not start feature development directly on main or develop.

---

## 41.5 Feature Development Lifecycle

Each Work Item follows:

```text
Create Feature Branch
        |
        v
Read Requirements
        |
        v
Prepare Specification
        |
        v
Technical Design
        |
        v
Task Planning
        |
        v
Implementation
        |
        v
Developer Self-Test
        |
        v
Commit Changes
        |
        v
Push Feature Branch
        |
        v
Open Pull Request
        |
        v
Independent Sub-Agent Review
        |
   +----+----+
   |         |
  PASS      FAIL
   |         |
   |         v
   |      Fix Issues
   |         |
   |         v
   |      Re-Review
   |
   v
Required CI Checks
        |
        v
Authorized PR Merge
        |
        v
        develop
        |
        v
Retain Feature Branch
```

The Work Item must not be considered fully complete
until the required review and merge conditions are satisfied.

---

## 41.6 Pull Request Policy

All changes entering develop or main must use Pull Requests.

Allowed Pull Request directions:

| Source | Target | Purpose |
|---|---|---|
| feature/* | develop | Feature integration |
| develop | main | Production release |

Prohibited Pull Request directions:

- feature/* → main
- main → develop as an ordinary feature integration path
- Unapproved alternative branches → main

Exceptional synchronization or recovery operations require
an explicitly documented and authorized procedure.

### PR-001 — Feature Pull Request

After completing implementation and developer validation,
the Developer Agent must open a Pull Request:

`feature/* → develop`

The PR must include:

- Work Item ID
- Objective
- Implementation summary
- Requirements addressed
- Specification references
- Files changed
- Tests executed
- Test results
- Known limitations
- Relevant screenshots when applicable
- Review status

### PR-002 — Release Pull Request

A production release must use:

`develop → main`

The Release PR must include:

- Proposed release version
- Release summary
- Included features
- Fixed defects
- Relevant Work Items
- Known limitations
- Testing evidence
- Regression test results
- Independent review results
- Deployment considerations

### PR-003 — Draft Pull Requests

Draft Pull Requests may be used for incomplete work.

A Draft PR must not be merged.

Before requesting final review:

- Implementation must be complete.
- Required tests must be executed.
- Documentation must be synchronized.
- The PR must be ready for review.

### PR-004 — No Direct Merge

Agents must not bypass the Pull Request process
by directly merging local branches and pushing
the resulting merge commit to develop or main.

Use GitHub Pull Requests as the authoritative merge path.

---

## 41.7 Independent Pull Request Review

Every Pull Request must undergo independent review
before merge.

This requirement applies to both:

- feature/* → develop
- develop → main

### REVIEW-001 — Separate Reviewer

The implementation must be reviewed by an independent
Sub-Agent or equivalent separate Agent execution.

The implementing Agent cannot approve its own changes.

### REVIEW-002 — Review Against Requirements

The Reviewer must validate the PR against:

- docs/REQUIREMENT.md
- Relevant feature specifications
- Approved technical design
- Development task plan
- Acceptance criteria
- Applicable architecture decisions
- Coding standards
- Security standards
- Testing requirements

### REVIEW-003 — Code Inspection

The Reviewer must inspect:

- Actual code changes
- Implementation diff
- Test changes
- Dependency changes
- Configuration changes
- Security implications
- Documentation changes

The Reviewer must not rely solely on the Developer Agent's
implementation summary.

### REVIEW-004 — Test Verification

The Reviewer must assess:

- Whether required tests were executed
- Whether the tests actually validate requirements
- Whether important edge cases are covered
- Whether test results are genuine
- Whether regressions were introduced

Where practical, the Reviewer should execute targeted
non-destructive validation independently.

### REVIEW-005 — Review Decisions

Allowed review decisions:

- PASS
- FAIL
- BLOCKED

PASS:

All mandatory conditions are satisfied.

FAIL:

Blocking issues or failed requirements remain.

BLOCKED:

The review cannot be completed because of missing
evidence, tools, access, or unresolved specification issues.

### REVIEW-006 — Required Remediation

If a Reviewer returns FAIL:

1. Record the findings.
2. Return the PR to the Developer Agent.
3. Fix the identified problems.
4. Add or update relevant tests.
5. Re-run validation.
6. Push the fixes to the same feature branch.
7. Request independent re-review.

The PR cannot be merged until blocking findings
are resolved and the review decision is PASS.

### REVIEW-007 — Changes After Approval

If implementation code changes after a PASS:

- The previous PASS does not automatically apply.
- Affected tests must be re-run.
- The updated diff must be reviewed.
- A new PASS must cover the revision being merged.

### REVIEW-008 — Review Evidence

Each review must identify:

- Pull Request number
- Work Item ID
- Reviewed commit SHA
- Reviewer execution reference
- Requirements checked
- Tests evaluated
- Findings
- Final decision

Review evidence must be associated with the PR.

The project may also store review reports under:

`docs/development/reviews/`

A review report committed to the repository must be included
in the final reviewed content.

A later documentation-only change must be checked for
consistency before merge.

### REVIEW-009 — GitHub Approval vs Agent Review

Independent Sub-Agent Review is a project quality gate.

It is not automatically equivalent to GitHub's native
Pull Request approval mechanism.

Native GitHub approval requirements may depend on
repository permissions and available integrations.

The project must not claim that GitHub branch protection
enforces Sub-Agent Review unless a trusted mechanism
actually validates the review evidence.

Where supported, an independent review validation check
should be integrated into CI.

No Agent may mark a review as PASS merely by editing
a status field or creating an unsupported claim.

---

## 41.8 Pull Request Merge Requirements

A Pull Request may be merged only when:

1. Its source and target branches are correct.
2. The PR is not in Draft status.
3. The implementation matches approved specifications.
4. Required CI checks pass.
5. Independent Sub-Agent Review returns PASS.
6. No unresolved BLOCKER findings remain.
7. No unresolved MAJOR findings remain.
8. Required documentation is updated.
9. The reviewed revision matches the implementation
   being merged.
10. Merge authorization has been obtained according
    to the applicable project workflow.

If any condition is unmet, do not merge.

### Merge Authorization

Feature PR:

Merge only after all required gates pass and the
authorized project workflow permits the merge.

Release PR:

Production promotion requires explicit Project Owner
authorization unless a separately approved release
policy delegates that authority.

A Reviewer PASS is not permission to deploy production.

---

## 41.9 Merge Strategy

The recommended merge method is:

**Create a Merge Commit**

This method should be used for:

- feature/* → develop
- develop → main

Reasons:

- Preserve feature branch commit history.
- Preserve explicit integration points.
- Improve traceability.
- Retain review and implementation history.
- Keep release integration visible in Git history.

Avoid squash or rebase merging as the default strategy.

Do not enable a branch rule requiring strictly linear
history if it conflicts with the approved Merge Commit
strategy.

Changes to the default merge strategy require review.

---

## 41.10 Feature Branch Retention

This requirement is CONFIRMED.

### BRANCH-RETENTION-001

Feature branches MUST NOT be automatically deleted
after Pull Request merge.

### BRANCH-RETENTION-002

Developer Agents and Reviewer Agents MUST NOT delete
merged feature branches.

### BRANCH-RETENTION-003

The Repository must disable:

`Automatically delete head branches`

in GitHub Pull Request settings.

### BRANCH-RETENTION-004

Merged feature branches must remain available on GitHub.

### BRANCH-RETENTION-005

A completed feature branch should not normally be reused
for unrelated development.

New Work Items should use new feature branches
created from the latest develop.

### BRANCH-RETENTION-006

Do not rewrite the history of merged feature branches.

### BRANCH-RETENTION-007

If a feature branch is accidentally deleted,
report the event and restore the branch when possible.

Branch deletion protection should be enabled where
supported by the repository configuration.

The purpose of branch retention is traceability
and preservation of development history.

---

## 41.11 GitHub Repository Protection

The Agent must propose or configure appropriate GitHub
Branch Protection Rules or Repository Rulesets.

Configuration changes affecting repository security
must be authorized before application.

### main Protection

Recommended protections:

- Require Pull Request before merging.
- Require required CI checks.
- Require review according to repository policy.
- Block force pushes.
- Block branch deletion.
- Restrict direct updates.
- Resolve required review conversations.
- Prevent unauthorized bypass.

### develop Protection

Recommended protections:

- Require Pull Request before merging.
- Require required CI checks.
- Require review according to repository policy.
- Block force pushes.
- Block branch deletion.
- Restrict direct updates.
- Resolve required review conversations.

### feature/* Protection

Recommended protections:

- Prevent accidental branch deletion.
- Prevent force pushes.
- Preserve merged branches.

Ensure feature development remains possible.

### Repository Settings

Recommended:

- Enable Merge Commits.
- Disable automatic head branch deletion.
- Configure required status checks.
- Configure appropriate review policies.
- Restrict protection bypasses.
- Protect main and develop from direct updates.

If required settings cannot be enforced because of
repository permissions or platform limitations,
document the limitation.

Do not falsely report protections as enabled.

---

## 41.12 CI/CD Integration

GitHub Actions must run relevant validation
for Pull Requests.

### Feature Pull Request Checks

For:

`feature/* → develop`

Recommended checks:

- Formatting
- Linting
- Type checking
- Unit tests
- Relevant integration tests
- Build verification
- Protocol validation
- Security checks where applicable

### Release Pull Request Checks

For:

`develop → main`

Recommended checks:

- All required feature checks
- Full regression tests
- Multiplayer integration tests
- E2E tests
- Browser compatibility checks
- Release build validation
- Relevant security validation
- Independent release review

### Required Status Checks

Required CI checks should be enforced through
GitHub repository protection where supported.

Checks must be associated with the revision being merged.

Do not bypass failing checks merely to accelerate development.

### Review Evidence Gate

Where supported, create a trusted mechanism to verify
that mandatory independent review has completed.

Review evidence must identify the reviewed commit SHA.

A self-declared PASS inside an untrusted PR file
must not be treated as sufficient proof by itself.

---

## 41.13 Development Branch Synchronization

Before opening a PR, the Developer Agent should check
whether the feature branch remains compatible with
the current develop branch.

If develop has advanced:

1. Fetch the latest remote branches.
2. Review the changes introduced into develop.
3. Update the feature branch when necessary.
4. Resolve conflicts carefully.
5. Re-run affected tests.
6. Update review evidence if necessary.

Do not force push or rewrite reviewed history merely
to make branch synchronization easier.

The implementation must be tested against the
integration state intended for merging.

---

## 41.14 Release Workflow

Production releases must be promoted from develop
to main through Pull Requests.

Release procedure:

1. Identify the release scope.
2. Verify included Work Items.
3. Confirm required feature reviews passed.
4. Verify develop integration health.
5. Prepare release notes.
6. Open PR from develop to main.
7. Run full release validation.
8. Launch Independent Release Reviewer.
9. Resolve blocking findings.
10. Obtain release authorization.
11. Merge using Merge Commit.
12. Tag the approved release commit on main.
13. Deploy the approved release.
14. Validate production deployment.
15. Record release results.

### Versioning

Recommended version format:

Semantic Versioning.

Examples:

- v0.1.0
- v0.2.0
- v0.3.0

Pre-release versions may use:

- v0.1.0-alpha.1
- v0.1.0-beta.1
- v0.1.0-rc.1

Version tags must identify release commits
reachable from main.

### Production Deployment

Production deployment must use approved
main-branch release source.

A successful merge into main does not automatically
prove that production deployment succeeded.

Post-deployment verification is required.

---

## 41.15 Hotfix Workflow

The project does not introduce a separate permanent
hotfix branch category.

Emergency fixes must follow the approved branching model.

Recommended procedure:

1. Identify the production defect.
2. Assess whether develop is suitable for release.
3. Create a feature branch for the fix.
4. Implement and validate the correction.
5. Open a PR targeting develop.
6. Complete independent review.
7. Merge into develop.
8. Prepare a develop-to-main Release PR.
9. Complete release validation.
10. Obtain release authorization.
11. Merge into main.
12. Publish the corrected release.

Example:

`feature/WXXX-fix-room-reconnect`

### Important Limitation

If develop contains unfinished changes that are not
safe for production, the normal release flow may be
unable to promote only the urgent fix.

Do not silently bypass the approved workflow.

The Project Owner must approve an exceptional recovery
or release strategy before proceeding.

Maintaining develop in a release-ready state reduces
this risk.

---

## 41.16 Agent Git Responsibilities

### Developer Agent

Must:

- Create feature branches from develop.
- Work only within the assigned branch.
- Commit implementation changes.
- Run required tests.
- Push feature branches.
- Open PRs against develop.
- Provide review evidence.
- Resolve review findings.
- Avoid unauthorized merges.

Must not:

- Push directly to main.
- Push directly to develop.
- Delete merged feature branches.
- Force push without authorization.
- Merge a PR with unresolved blocking findings.
- Claim review passed without independent evidence.

### Reviewer Sub-Agent

Must:

- Review the PR independently.
- Inspect the implementation diff.
- Verify requirements and tests.
- Report findings with evidence.
- Issue PASS, FAIL, or BLOCKED.

Must not:

- Approve its own implementation.
- Modify code without authorization.
- Merge PRs solely because its review passed.
- Delete branches.
- Override required GitHub checks.

### Project Owner

Retains authority over:

- Material requirement changes
- Core architecture changes
- Exceptional workflow changes
- Production release authorization
- Exceptions to branch and security policies

---

## 41.17 Pull Request Naming

### Feature PR

Recommended title:

`[WXXX] <Feature Description>`

Examples:

```text
[W001] Initialize Monorepo
[W006] Implement Card Matching Engine
[W011] Implement Multiplayer Speed Battle
```

### Release PR

Recommended title:

`[Release] vX.Y.Z — <Release Description>`

Example:

```text
[Release] v0.1.0 — Closed Alpha
```

---

## 41.18 Pull Request Description Template

Every feature PR should include:

```markdown
## Work Item

WXXX — <Title>

## Objective

<Description>

## Related Requirements

- FR-XXX
- BR-XXX
- NFR-XXX

## Implementation Summary

<Changes>

## Files Changed

<Relevant files>

## Tests

| Test | Result |
|---|---|
| Formatting | PASS / FAIL / NOT_RUN |
| Lint | PASS / FAIL / NOT_RUN |
| Unit Tests | PASS / FAIL / NOT_RUN |
| Integration Tests | PASS / FAIL / NOT_RUN |
| E2E Tests | PASS / FAIL / NOT_RUN |

## Independent Review

Reviewer:
<Agent execution reference>

Reviewed Revision:
<Commit SHA>

Decision:
PASS / FAIL / BLOCKED

Review Evidence:
<PR comment, check, or review report>

## Known Limitations

<Limitations>

## Documentation Changes

<Documentation references>

## Merge Readiness

READY / NOT_READY / BLOCKED
```

---

## 41.19 Git Workflow Acceptance Criteria

### GIT-AC-001 — Feature Branch Creation

Given an approved Work Item,

When the Developer Agent starts implementation,

Then a new feature branch must be created from develop.

### GIT-AC-002 — No Direct Push

Given protected main and develop branches,

When an Agent attempts to push directly,

Then the repository must reject the operation
where protection is supported.

### GIT-AC-003 — Feature Pull Request

Given a completed feature implementation,

When the Agent requests integration,

Then the PR must target develop.

### GIT-AC-004 — Independent Review

Given a feature PR,

When independent review has not passed,

Then the PR must not be merged.

### GIT-AC-005 — Required CI

Given a feature PR with failing required checks,

Then the PR must not be merged.

### GIT-AC-006 — Branch Retention

Given a feature PR has been merged,

Then the source feature branch must remain
available on GitHub.

### GIT-AC-007 — Release Promotion

Given an approved production release,

When changes are promoted,

Then the PR source must be develop
and the target must be main.

### GIT-AC-008 — Release Review

Given a develop-to-main PR,

Then independent release review and required
validation must pass before merge.

### GIT-AC-009 — No Unauthorized Release

Given the release has not been authorized,

Then the Agent must not deploy to production.

### GIT-AC-010 — Review Revision Integrity

Given a PR received PASS,

When material implementation changes are pushed afterward,

Then affected validation and independent review
must be repeated before merge.

---

## 41.20 Git Workflow Definition of Done

A development Work Item may be marked DONE only when:

- [ ] Feature branch created from develop.
- [ ] Approved scope implemented.
- [ ] Relevant specifications satisfied.
- [ ] Required tests executed and passed.
- [ ] Changes committed to feature branch.
- [ ] Feature branch pushed to GitHub.
- [ ] PR opened against develop.
- [ ] Independent Sub-Agent Review completed.
- [ ] Review decision PASS.
- [ ] No unresolved BLOCKER findings.
- [ ] No unresolved MAJOR findings.
- [ ] Required CI checks passed.
- [ ] Documentation synchronized.
- [ ] PR merge authorized.
- [ ] PR merged into develop.
- [ ] Feature branch retained after merge.
- [ ] Final implementation report recorded.

A Work Item that is implementation-complete but not yet
merged must remain in an appropriate pre-completion state.

Recommended status:

READY_FOR_REVIEW or VERIFIED.

Do not mark DONE merely because code has been written
or reviewed.

---

## 41.21 Final Git Workflow Rules

The following rules are mandatory:

1. main is the Production branch.
2. develop is the Integration branch.
3. feature/* is used for development.
4. New feature branches originate from develop.
5. Feature changes enter develop through Pull Requests.
6. Production changes enter main from develop through Pull Requests.
7. Independent Sub-Agent Review is mandatory before merge.
8. Required CI checks must pass before merge.
9. Direct pushes to main and develop are prohibited.
10. Force pushes to protected branches are prohibited.
11. Feature branches must not be deleted after merge.
12. GitHub automatic head branch deletion must be disabled.
13. Merge Commit is the default merge strategy.
14. Material changes after review require re-review.
15. Production release requires authorization.
16. All Git operations must preserve requirement traceability.
17. Agent development reports must include PR and review evidence.
18. Git Workflow violations must be reported, not silently bypassed.

This workflow is part of the mandatory Agent-first SDD
development process for Prickle Panic!.

# END OF DOCUMENT
