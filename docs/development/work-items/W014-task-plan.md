# Task Plan — W014: Coco Assets & Gameplay Polish

Specification: `docs/development/work-items/W014-coco-assets-gameplay-polish.md` (READY_FOR_REVIEW; asset production requires approved specification and licensing review)
Branch: `feature/W014-coco-assets-gameplay-polish-spec`

| Task | Change | Requirement / acceptance refs | Deliverable | Validation |
|---|---|---|---|---|
| T01 | Approve modular Coco design, asset catalog/manifest, and source/licensing policy | §§4.1–4.10, 8.1–8.10, 22; W014-AC-01–10 | Approved technical/art direction and manifest schema | Owner approval where required; independent design/spec review |
| T02 | Produce/review modular Coco and 32 card variant assets with editable sources | §§4.5–4.9; W014-AC-01–03/08 | Approved character, card assets, source files and provenance records | Exhaustive variant mapping; visual/readability review; licensing audit |
| T03 | Integrate state-driven animation and Thai gameplay feedback | §§4.6, 22; W014-AC-04–05 | Shared visual feedback and accessible UI polish | Component/browser checks; reduced motion, keyboard, mobile, Phaser cleanup |
| T04 | Add basic sound set, controls, and robust asset-loading fallbacks | §§3.1, 22.7; W014-AC-06–07 | Local audio assets/settings/fallbacks | Gesture/autoplay, mute, missing asset, no-audio checks |
| T05 | Measure responsive behavior/performance, validate asset licenses, record evidence, and obtain exact-head independent review | W014-AC-03/08–10 | Validation report, manifest and PR evidence | `pnpm check`, Playwright, visual/browser and security/license review |

## Dependencies

- W005 Protocol Contracts and W006 Card Matching Engine (DONE); W007 Local Gameplay Prototype (DONE).
- Coordinate final integration boundaries with approved W011 multiplayer gameplay before replacing shared visual interfaces.
- W012 Results and W013 reconnect/error states should consume the shared asset/UI interfaces where their approved specs require them.

## Out of scope

- Gameplay rule changes, new card attributes, changing the 4×4 / 32-variant contract, or changing score/timer authority.
- Classic Rush, Party Chaos, English localization, cosmetics economy, or remote asset service.
- Shipping unreviewed or license-uncertain source assets.

## Completion evidence

- [ ] Approved design and manifest schema recorded
- [ ] All shipped asset sources, licenses, attribution and provenance recorded
- [ ] Actual visual, responsive, accessibility, audio and performance evidence recorded
- [ ] Independent review PASS for exact implementation commit
- [ ] Blocking findings resolved and re-reviewed
- [ ] Documentation synchronized
- [ ] PR merged to `develop` with Merge Commit; feature branch retained
