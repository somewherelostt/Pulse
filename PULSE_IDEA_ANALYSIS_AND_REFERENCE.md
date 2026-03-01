# Pulse — Idea Analysis & Future Reference

**Document purpose:** Analysis of the Pulse concept (no code). Use this for product decisions, pitch prep, and implementation reference.

---

## 1. Executive Summary

**Pulse** is a mental-health early-warning system that:
- **Detects** behavioral decline from passive data (calendar, sleep, screen, browsing).
- **Surfaces** patterns via a unified dashboard and LLM-generated insights.
- **Connects** users to anonymized peers who recovered from similar patterns (Constellation).

The core bet: *behavioral data already signals decline weeks before people seek help; connecting that signal to peer support shortens the path from “I’m not okay” to “I found someone who gets it.”*

---

## 2. Idea Analysis

### 2.1 Problem–Solution Fit

| Aspect | Assessment |
|--------|------------|
| **Problem clarity** | Strong. “Deterioration in the gaps” and “12-year delay to help” are concrete and research-backed. |
| **Data availability** | Real. Calendar, sleep trackers, and browser/screen data are widely available; the gap is aggregation and interpretation. |
| **Solution fit** | Detection + pattern display + peer match addresses both “notice earlier” and “bridge to support.” |
| **Edge case** | Users with no wearable: manual fallback and Calendar-only path (Layer 1) keep the product viable. |

**Verdict:** Problem and solution are well aligned. Layer 1 (calendar + mood + correlation + LLM) is the right minimum to validate “does unified behavioral + mood insight add value?”

---

### 2.2 Differentiation

- **vs. mood trackers:** Pulse adds passive behavioral context (calendar load, sleep, browsing) instead of mood-only logs.
- **vs. wearables-only apps:** Multi-source (calendar, screen, browser) and cross-source correlation + LLM narrative.
- **vs. therapy platforms:** Not a replacement; focuses on early signal and peer connection, with explicit crisis escalation to professional resources.
- **Constellation:** Peer matching on *behavioral fingerprint* (recovered-from-similar-pattern) is a clear differentiator; “someone who survived this pattern” is a strong narrative.

**Risk:** “Behavioral fingerprint” and drift scoring could feel surveillance-like if not framed as user-controlled, opt-in, and deletable. Your ethical commitments (anonymous, RLS, one-click export/delete) directly address this.

---

### 2.3 Technical Coherence

- **Single backend (Go on Render):** API, cron, feature extraction, drift, and WebRTC signaling in one place — simpler ops and one deployment.
- **Supabase as core:** Auth (anonymous + magic link), Postgres, pgvector, RLS, and storage fit the “privacy-first, user-scoped data” model.
- **LLM strategy:** Groq primary + Cerebras fallback with same prompts keeps latency and reliability acceptable; five well-scoped LLM calls limit cost and complexity.
- **Frontend (vanilla HTML/CSS/JS):** Aligns with “ship fast, readable, no framework lock-in”; Chart.js + D3 cover the planned visualizations.
- **Extension:** Local-first (batched every 60s), no raw content to backend — good for privacy and compliance.

**Gaps to decide later:**
- Embedding model choice (OpenAI vs Voyage) and where it runs (Supabase vs Go).
- Exact drift threshold and baseline (14-day) tuning; likely needs product/data iteration.

---

### 2.4 Ethical Design

Your commitments cover the main failure modes:

- **Anonymous by default + RLS:** Limits exposure and ensures user-scoped access.
- **No peer session content stored:** Reduces legal and trust risk.
- **Crisis path:** Drift alerts surface professional resources; safety layer scans for crisis signals and never blocks, only offers alternatives.
- **Transparency:** “AI-generated, not a diagnosis” and pattern-level (not raw data) matching context.
- **Constellation opt-in + opt-out anytime:** Consent and control are explicit.
- **Export + delete:** Supports GDPR-style expectations.

**Suggestion:** Document a short “safety playbook” (when to show crisis resources, when to avoid pushing peer match, how to handle “I’m fine” dismissals) so future you and any team stay consistent.

---

## 3. Architecture Summary (Reference)

### 3.1 Stack at a Glance

| Layer | Choice | Role |
|-------|--------|------|
| Frontend | Vercel, HTML/CSS/JS, Chart.js, D3, Supabase JS | Dashboard, charts, auth, WebRTC client |
| Backend | Go (chi, cron, oauth2, supabase-go, pgx) on Render | API, crons, features, drift, signaling |
| DB | Supabase (Postgres, pgvector, Auth, RLS, Storage) | All persistent data, vectors, auth |
| LLM | Groq → Cerebras fallback | All 5 inference calls |
| Embeddings | OpenAI or Voyage (TBD) | One embedding per user per day |
| Extension | Chrome/Firefox MV3, vanilla JS | Browser events → batched POST to Go |

### 3.2 Data Flow (Conceptual)

```
User devices (browser + extension) 
  → Go server (API + cron) 
  → Supabase (raw + features + embeddings + insights)
  → LLM (Groq/Cerebras) for narratives and match context
  → User (dashboard, alerts, peer session)
```

### 3.3 The 5 LLM Calls (Quick Reference)

1. **Pattern analysis** — 30-day features + mood → patterns, lags, confidence, plain English (weekly + on-demand).
2. **Circadian narrative** — Sleep consistency, debt, constraints → one micro-intervention (weekly).
3. **Predictive warning** — Current trajectory + episode signatures → trajectory, confidence, watch/positive signals (drift alert).
4. **Weekly narrative** — Correlations + mood + patterns → 3-paragraph insight (Friday).
5. **Match context** — Seeker + supporter anonymized patterns → warm session framing (Constellation).

### 3.4 Build Order (As in Brief)

1. **Layer 1** — Calendar + Mood + Correlation + LLM insight *(ship first)*  
2. **Layer 2** — Sleep + Circadian  
3. **Layer 3** — Browser extension + Screen features  
4. **Layer 4** — Embeddings + Drift (pgvector)  
5. **Layer 5** — Peer matching + WebRTC (Constellation)  

Dependencies are consistent: Layer 4 needs features from 1–3; Layer 5 needs embeddings and drift.

---

## 4. Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Drift threshold too sensitive | Tune with real data; allow user-adjustable sensitivity or “learning period” before first alert. |
| Peer pool too small at launch | Start with “pattern insight only” and add Constellation when N is meaningful; or soft-launch in one community. |
| LLM cost at scale | Cap context size, cache weekly narratives, use smaller model for match context. |
| Extension adoption | Layer 1 works without it; position extension as “deeper insight” not required. |
| Regulatory (health/wellness) | Keep language “insight and support,” not diagnosis or treatment; crisis path to professionals. |
| WebRTC reliability (NAT/firewalls) | TURN server or fallback to “schedule a time” if P2P fails; document in architecture. |

---

## 5. Key Product Decisions (Locked for Reference)

- **Anonymous-first:** No email required; magic link or anonymous session.
- **Calendar mandatory for Layer 1;** sleep and extension optional.
- **14-day baseline** for drift; **30-day** for pattern analysis.
- **Weekly report:** Friday; **drift:** triggered by threshold, not scheduled.
- **Peer sessions:** 30 min, no names, no recording, rating + “talk again?” only.
- **Crisis:** Always offer professional resources; never block user choice.

---

## 6. Glossary (Future Reference)

| Term | Meaning |
|------|--------|
| **Drift score** | Cosine similarity (or derived metric) of today’s embedding vs 14-day baseline; high drift → alert. |
| **Behavioral fingerprint** | Normalized feature vector (and/or embedding) representing pattern; used for peer match only. |
| **Constellation** | Peer-matching product: match on similar past pattern, prioritize “recovered” peers, WebRTC session. |
| **Episode replay** | View behavioral data from past low-mood period; optional “save as warning signature.” |
| **RLS** | Row Level Security (Supabase); user can only access own rows. |

---

## 7. What’s Not in This Doc

- API specs, table DDL, or environment variables (to be created in implementation).
- UI/UX wireframes or copy (to be created with Layer 1).
- Legal or compliance checklist (recommended before public launch).
- Go package layout or file structure (to be created when writing Layer 1).

---

## 8. Conclusion

Pulse is a coherent idea with a clear problem (decline in the gaps, long delay to support), a plausible solution (passive detection + pattern display + peer connection), and a stack and build order that support incremental delivery. The main product risk is calibration (drift threshold, narrative tone); the main technical risk is scaling LLM and WebRTC. Ethical commitments and safety design are explicitly addressed.

**Recommended next step:** Implement Layer 1 (Calendar + Mood + Correlation + one LLM insight) and validate that users find the unified pattern view and narrative useful before adding sleep, extension, drift, and Constellation.


# Pulse Design System — Cursor Reference

## Colors

--bg-base:        #0A0A0F
--bg-surface:     #111118
--bg-raised:      #1A1A24
--border:         #2A2A3A
--border-subtle:  #1E1E2A
--primary:        #6E7BF2
--primary-glow:   rgba(110, 123, 242, 0.15)
--primary-dim:    rgba(110, 123, 242, 0.08)
--accent:         #4ECDC4
--accent-glow:    rgba(78, 205, 196, 0.15)
--warning:        #F7B731
--warning-glow:   rgba(247, 183, 49, 0.15)
--danger:         #FF6B6B
--text-primary:   #F0F0F5
--text-secondary: #8888A0
--text-muted:     #4A4A5A

## Typography

font-display:  'Geist', 'Inter Display', sans-serif
font-body:     'Inter', sans-serif
font-mono:     'Geist Mono', 'JetBrains Mono', monospace

## Font Weights

Thin headlines:   300
Body:             400
Emphasis:         500
Buttons/labels:   600
Data values:      mono 500

## Font Sizes (clamp)

--text-xs:   12px
--text-sm:   14px
--text-base: 16px
--text-lg:   18px
--text-xl:   24px
--text-2xl:  36px
--text-3xl:  48px
--text-4xl:  clamp(48px, 6vw, 72px)
--text-5xl:  clamp(64px, 8vw, 96px)

## Spacing

Base unit: 4px
--space-1:  4px   --space-2:  8px   --space-3: 12px
--space-4: 16px   --space-5: 20px   --space-6: 24px
--space-8: 32px   --space-10: 40px  --space-12: 48px
--space-16: 64px  --space-20: 80px  --space-24: 96px

## Border Radius

--radius-sm:   8px   (tags, badges, inputs)
--radius-md:  12px   (inner cards)
--radius-lg:  16px   (cards)
--radius-xl:  20px   (large cards, modals)
--radius-full: 9999px (pills)

## Shadows

--shadow-sm:  0 2px 8px rgba(0,0,0,0.4)
--shadow-md:  0 8px 24px rgba(0,0,0,0.5)
--shadow-lg:  0 20px 40px rgba(0,0,0,0.6)
--shadow-primary: 0 8px 32px rgba(110,123,242,0.25)
--shadow-accent:  0 8px 32px rgba(78,205,196,0.2)

## Animation Tokens

--spring-soft:   type spring, stiffness 100, damping 20
--spring-snappy: type spring, stiffness 400, damping 25
--duration-fast:   0.15s
--duration-base:   0.3s
--duration-slow:   0.6s
--duration-data:   1.5s  (chart draw animations)

## Component Patterns

### Card

background: var(--bg-surface)
border: 1px solid var(--border)
border-radius: var(--radius-lg)
padding: 24px
hover: y -4px, shadow-primary (if interactive)

### Metric Card

Same as card + top accent bar (3px, primary or warning color)
Value: font-mono, text-2xl, text-primary
Label: text-sm, text-muted, uppercase, tracking 0.08em
Delta: text-sm, with ↑ (accent) or ↓ (danger) prefix

### Insight Card (LLM generated)

Left border: 3px solid primary
background: var(--bg-raised)
Small label top: "AI INSIGHT · NOT A DIAGNOSIS"
  font-mono, 10px, text-muted, required on every LLM card

### Button Primary

background: var(--primary)
color: white
height: 48px
padding: 0 24px
border-radius: var(--radius-sm)
font-weight: 600
hover: scale 1.02, shadow-primary

### Button Ghost  

background: transparent
border: 1px solid var(--border)
color: var(--text-secondary)
hover: border-color primary, color text-primary

### Status Pill

background: var(--bg-raised)
border: 1px solid var(--border)
border-radius: var(--radius-full)
padding: 4px 12px
font-mono, 11px, uppercase, tracking 0.1em
With colored dot prefix (pulse animation for live states)

### Drift Alert

background: rgba(247,183,49,0.08)
border: 1px solid rgba(247,183,49,0.3)
border-radius: var(--radius-md)
Icon: ⚠ in --warning
Always includes: "Seek professional support" link at bottom

## Layout

max-width: 1200px
section-padding-y: 120px desktop, 80px tablet, 60px mobile
grid: 12 col, 24px gutter desktop / 4 col, 16px gutter mobile

## Ethical UI Rules

1. Every LLM-generated card MUST include disclaimer
2. Drift alerts MUST include professional resource link
3. Constellation UI MUST show "Anonymous · Encrypted · Not recorded"
4. No dark patterns — no auto-opt-in, no hidden consent
5. Crisis resources always visible in footer
6. Data deletion always one click away

## Copy Voice

- Direct, never preachy
- Data-confident, not clinical  
- Warm but not soft
- Never: "journey", "wellness journey", "start your journey"
- Never: generic wellness adjectives (calming, soothing, peaceful)
- Do: specific, data-referenced, honest about limitations
- Always: remind it's not a replacement for professional care

---



*Document generated from project brief for analysis and future reference. No code.*


