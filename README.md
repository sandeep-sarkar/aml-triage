# aml-triage

An AI engineering practice project: an alert-triage assistant for a
transaction-monitoring style workflow (inspired by AML systems), built to
learn and demonstrate applied AI engineering end to end — not just prompting,
but evals, RAG, and the platform around it.

## What it does

Given a flagged alert (a customer + a set of transactions), produce a
disposition: `escalate` or `close`, with reasoning and a cited policy
clause.

## Approach

1. **Rule-based baseline** — a scoring function in the spirit of
   traditional AML risk scoring, used as the floor the LLM has to beat.
2. **Eval harness** — precision, recall, and an explanation-quality rubric,
   run against a labelled synthetic dataset.
3. **LLM v1** — structured-output disposition, no retrieval.
4. **LLM + RAG** — retrieval over synthetic policy documents (pgvector),
   measured against the same eval harness.
5. **Platform** — LiteLLM gateway, vLLM on Kubernetes, Langfuse tracing,
   Grafana/DCGM dashboards, evals wired into CI.

All data is synthetic — no real financial data is used or required.

## Status

Early scaffolding. See `internal/` for current progress.

## Running locally

```
cd deploy
docker compose up -d
```

Requires `OPENAI_API_KEY` (or `ANTHROPIC_API_KEY`) set in your environment
or a `.env` file in `deploy/`.

## Layout

```
cmd/          entrypoints: triage app, data generator, eval harness
internal/     alert types, scoring, LLM client, RAG, eval metrics
data/         synthetic policy docs and generated alerts
deploy/       docker-compose, k8s manifests (later)
```
