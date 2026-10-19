# Development Roadmap

The Work Item IDs below preserve the baseline in `docs/REQUIREMENT.md` §32. Each item requires an approved specification and an independent review before completion. Speed Battle is the MVP critical path; Classic Rush and Party Chaos remain post-MVP.

| Phase | Work Item | Depends on | Outcome |
|---|---|---|---|
| Foundation | W001 — Repository Bootstrap | Initial repository seed | Reproducible web and Go workspaces |
| Foundation | W002 — SDD Documentation Foundation | W001 | Entry-point instructions, traceability, templates, architecture overview |
| Foundation | W003 — Skill Discovery & Approval | W001, W002 | Reviewed project-local Skills and registry |
| Foundation | W004 — Development Standards & CI | W001–W003 | Standards, workflows, protected CI checks |
| Core engine | W005 — Protocol Contracts | W002, W004 | Versioned messages, schemas, errors, compatibility rules |
| Core engine | W006 — Card Matching Engine | W004, W005 | Pure, exhaustively tested card rules and generation |
| Core engine | W007 — Local Gameplay Prototype | W001, W006 | Local Thai Speed Battle interaction and responsive prototype |
| Multiplayer foundation | W008 — Guest Session | W004, W005 | Secure guest identity, username validation and recovery |
| Multiplayer foundation | W009 — Room Management | W005, W008 | Private Rooms, invites, Lobby, Ready and host permissions |
| Multiplayer foundation | W010 — WebSocket Foundation | W005, W008, W009 | Authorized connection and bounded state synchronization |
| Online gameplay | W011 — Multiplayer Speed Battle | W006–W010 | Server-authoritative multi-player Speed Battle |
| Online gameplay | W012 — Results & Rematch | W011 | Results, winners, return to Lobby and reset Match |
| Online gameplay | W013 — Reconnection & Anti-Griefing | W010–W012 | Grace period, forfeit, host migration and rate limits |
| Product experience | W014 — Coco Assets & Gameplay Polish | W007 and integrated MVP | Consistent assets, audio, animation and responsive polish |
| MVP validation | W015 — Closed Alpha Readiness | W011–W014 | Independent end-to-end, security, compatibility and performance acceptance |
| Post-MVP | New Work Items for Classic Rush and Party Chaos | W015 | Assign IDs after MVP review; preserve FR-045–063 and AC-011–016 |

Dependencies are planning guidance; changes require documented rationale and traceability. Do not start all modes together or mark an implementation-only item DONE.
