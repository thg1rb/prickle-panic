# Project Documentation

## Authority

[`REQUIREMENT.md`](REQUIREMENT.md) is the initial authoritative source for product behavior, gameplay, architecture direction, security, SDD, and workflow requirements. Derived documents explain or specialize those requirements; they do not replace them. Preserve requirement IDs and obtain Project Owner approval before changing confirmed decisions.

Requirements use the statuses CONFIRMED, BASELINE, TUNING, PROPOSED, DEFERRED, and OUT_OF_SCOPE as defined in `REQUIREMENT.md` §1.3. Do not silently promote a proposal or tuning value to a confirmed decision.

## Document ownership

| Area | Owner | Purpose | Authoritative references | Update when |
|---|---|---|---|---|
| Product and gameplay | Project Owner | `REQUIREMENT.md`; approves material product/rule changes | §1–17, §22, §28, §35–40 | An approved product or gameplay decision changes |
| Architecture | Lead/Architecture Agent | `architecture/overview.md` and ADRs | §18–20, §30–31 | An architecture decision is approved or superseded |
| Feature behavior and tasks | Work Item author | Approved specifications and task plans in `development/work-items/` | Relevant requirement IDs named in each Work Item | Scope, acceptance, dependencies, or approved design changes |
| Engineering and testing | Maintainers / Work Item authors | Standards and test strategy as introduced by their Work Items | §20–21, §27, §29 | Tooling, quality gates, or validation policy changes |
| Agent workflow and Skills | SDD Coordinator | `agents/`, templates, and Skills registry | §23–26, §33–34 | Workflow or approved Skill/revision changes |
| Review evidence | Independent Reviewer | Review decisions and findings associated with the reviewed commit | §25–26, §41 | After each independent review; preserve the reviewed revision and decision |

## Structure

- `agents/` — entry-point workflow, templates, and approved Skills registry.
- `architecture/` — system boundaries and decision records.
- `development/` — roadmap, Work Item specifications, and review evidence.
- `testing/` — validation policy and evidence conventions.
- `design/` and `standards/` — create when a Work Item needs authoritative detail in those areas.

Do not create empty directories or duplicate requirement content. Ownership, purpose, authoritative references, and update conditions may be inherited from the applicable family row above. Each feature specification and task plan must list its exact requirement IDs; other documents may rely on their family references unless they address a narrower requirement that needs an explicit citation.

## Review and maintenance

Work Item changes must update affected specifications, architecture decisions, test references, and this index when document ownership or structure changes. Reviewers check consistency with `REQUIREMENT.md`. If a conflict is found, record it, propose a resolution, and obtain the required approval before changing behavior.
