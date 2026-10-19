# Project Documentation

## Authority

[`REQUIREMENT.md`](REQUIREMENT.md) is the initial authoritative source for product behavior, gameplay, architecture direction, security, SDD, and workflow requirements. Derived documents explain or specialize those requirements; they do not replace them. Preserve requirement IDs and obtain Project Owner approval before changing confirmed decisions.

Requirements use the statuses CONFIRMED, BASELINE, TUNING, PROPOSED, DEFERRED, and OUT_OF_SCOPE as defined in `REQUIREMENT.md` §1.3. Do not silently promote a proposal or tuning value to a confirmed decision.

## Document ownership

| Area | Owner | Purpose |
|---|---|---|
| Product and gameplay | Project Owner | `REQUIREMENT.md`; approves material product/rule changes |
| Architecture | Lead/Architecture Agent | `architecture/overview.md` and ADRs |
| Feature behavior and tasks | Work Item author | Approved specifications and task plans in `development/work-items/` |
| Engineering and testing | Maintainers / Work Item authors | Standards and test strategy as introduced by their Work Items |
| Agent workflow and Skills | SDD Coordinator | `agents/`, templates, and Skills registry |
| Review evidence | Independent Reviewer | Review decisions and findings associated with the reviewed commit |

## Structure

- `agents/` — entry-point workflow, templates, and approved Skills registry.
- `architecture/` — system boundaries and decision records.
- `development/` — roadmap, Work Item specifications, and review evidence.
- `testing/` — validation policy and evidence conventions.
- `design/` and `standards/` — create when a Work Item needs authoritative detail in those areas.

Do not create empty directories or duplicate requirement content. Each derived document must identify its owner, purpose, relevant requirement references, and update conditions.

## Review and maintenance

Work Item changes must update affected specifications, architecture decisions, test references, and this index when document ownership or structure changes. Reviewers check consistency with `REQUIREMENT.md`. If a conflict is found, record it, propose a resolution, and obtain the required approval before changing behavior.
