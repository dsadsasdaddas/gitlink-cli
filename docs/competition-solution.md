# GitLink CLI Agent Workflow Enhancement Suite

## 1. Background

GitLink CLI serves both human maintainers and AI Agents. The competition focuses on intelligent open-source contribution workflows, where structured analysis, stable output, and safe automation matter more than raw command count.

## 2. Problem

Open-source maintenance often suffers from:

- Issue backlog and delayed triage
- High PR review cost
- Repetitive release note preparation
- Lack of structured repository health evaluation
- AI Agents needing stable, machine-readable output

## 3. Solution

This project extends GitLink CLI with the **GitLink CLI Agent Workflow Enhancement Suite**.

Implemented now:

- `workflow +triage`
- `workflow +health`
- read-only GitLink fetch layer for workflow triage and health
- expanded fetch boundary tests for empty responses, label and author normalization, error-in-body handling, alternative activity timestamps, release shapes, and CI unavailability
- local-first analysis with no LLM dependency
- stable Agent-facing JSON / table / markdown output

Planned next:

- `workflow +pr-summary`
- `workflow +release-notes`
- `workflow +stale`

## 4. Technical Route

- Go + Cobra + existing shortcut architecture
- rule-based analysis
- stable DTOs
- `json` / `table` / `markdown` renderers
- `en` / `zh-CN` message mapping
- no LLM dependency
- local-first, dry-run-safe workflow design

## 5. Implemented Features

### workflow +triage

- issue type detection
- priority scoring
- confidence scoring
- missing information detection
- risk flags
- recommended action
- suggested comment
- reasoning and matched rules

### workflow +health

- health score
- risk level
- metrics
- scoring notes
- recommendations
- unknown metric tolerance

## 6. Innovation Points

- Agent-native structured output
- rule-based intelligence without external LLM dependency
- explainable workflow decisions
- safety-first local analysis
- bilingual command output
- extensible workflow command design
- competition-friendly incremental PR path

## 7. Testing and Verification

- Unit tests cover triage, health scoring, messages, rendering, and command helpers.
- Fetch-layer tests cover issue normalization and repository health probing with `httptest`.
- Boundary tests cover empty responses, label and author normalization, error-in-body handling, alternative activity timestamps, release response shapes, and CI unavailability.
- Local command examples were executed successfully.
- Full repository testing passed in the current environment.
- Automated tests use `httptest` and do not depend on real remote API availability.

## 8. Demonstration Plan

### Official repository

Use `Gitlink/gitlink-cli` as the reference repository:

1. `workflow +triage` with English table output
2. `workflow +triage` with security JSON output
3. `workflow +triage` with Chinese markdown output
4. `workflow +health` with table output
5. `workflow +health` with risky JSON output
6. Explain how agents consume stable JSON

### Self-built test repository

Use a small demo repository to show:

- bug triage
- security triage
- docs triage
- healthy repo score
- risky repo score

## 9. Roadmap

- Phase 1: local workflow prototype, completed
- Phase 2: API fetch and normalization, completed
- Phase 3: `pr-summary`, `release-notes`, `stale`

## 10. PR Plan

- PR 1: workflow rule engine and local commands
- PR 2: documentation and tests
- PR 3: API fetch layer
- PR 4: `pr-summary` / `release-notes`
