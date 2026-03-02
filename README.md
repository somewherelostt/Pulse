# Pulse

**Mental health early-warning system powered by behavioral data analysis.**

---

## The Problem

Mental health deterioration doesn't happen overnight—it shows up in patterns weeks before people realize they need help:

- **12-year average delay** between onset of symptoms and seeking treatment
- **Behavioral signals go unnoticed**: calendar overload, sleep disruption, browsing patterns, screen time changes
- **Isolation amplifies decline**: people struggle alone without knowing others have survived similar patterns
- **Traditional mood trackers miss context**: logging emotions without understanding the behavioral factors behind them

The gap isn't awareness—it's **detection and connection**. By the time someone realizes they're struggling, weeks of behavioral decline have already occurred.

---

## The Solution

Pulse aggregates passive behavioral data to detect early warning signals and connects users to peer support:

### 1. **Multi-Source Behavioral Detection**

- **Calendar analysis**: event density, fragmentation, boundary violations, recovery time
- **Sleep tracking**: consistency, debt, circadian rhythm disruption
- **Digital behavior**: screen time patterns, browsing activity (local-first, privacy-preserving)
- **Mood correlation**: cross-reference self-reported mood with behavioral patterns

### 2. **AI-Powered Pattern Analysis**

- LLM-generated insights from 30-day behavioral fingerprints
- Drift detection using vector embeddings (14-day baseline comparison)
- Predictive warnings when current trajectory matches historical decline patterns
- Weekly narrative reports highlighting correlations and intervention points

### 3. **Constellation: Anonymous Peer Matching**

- Match users with peers who **recovered from similar behavioral patterns**
- Anonymous, encrypted WebRTC sessions (no recording, no data retention)
- Behavioral fingerprint matching—not demographics or symptoms
- "Someone who survived this pattern" > generic support groups

### 4. **Privacy-First Architecture**

- **Anonymous by default**: no email required
- **Row-level security (RLS)**: user-scoped data access only
- **Local-first extension**: browser data never leaves device raw
- **One-click data export/delete**: full user control
- **Crisis escalation**: always surfaces professional resources, never blocks access

---

## Use Cases

### Early Intervention

**Scenario**: User's calendar shows increasing meeting density, sleep tracker indicates 3+ nights of poor rest, mood logs show gradual decline.

**Pulse Action**: Generates drift alert with pattern visualization, suggests checking recovery time between high-load days, surfaces professional resources.

---

### Pattern Recognition

**Scenario**: User can't pinpoint why they feel off. Mood is "fine" but energy is depleted.

**Pulse Action**: Weekly LLM narrative reveals correlation between late-night screen time spikes, next-day calendar fragmentation, and mood dips 48 hours later.

---

### Peer Support Discovery

**Scenario**: User experiencing burnout symptoms, unsure where to start.

**Pulse Action**: Constellation matches them with a peer who recovered from a similar behavioral fingerprint (calendar overload + sleep debt + social withdrawal). 30-minute anonymous session provides lived experience and validation.

---

### Episode Analysis

**Scenario**: User wants to understand past decline to prevent future recurrence.

**Pulse Action**: Episode replay visualizes behavioral data from low-mood periods, allows saving as "warning signature" for future drift detection.

---

## Tech Stack

### Frontend

- **Framework**: Next.js 14 (App Router)
- **Styling**: CSS Modules, vanilla CSS (no frameworks)
- **UI**: Chart.js (time series), D3.js (correlations, network graphs)
- **Auth**: Supabase Auth (anonymous + magic link)
- **Real-time**: Supabase Realtime, WebRTC (peer sessions)

### Backend

- **Language**: Go 1.21+
- **Framework**: Chi (routing), standard library
- **Deployment**: Render (single backend: API + cron + WebRTC signaling)
- **OAuth**: Google Calendar API (official Go client)
- **Scheduling**: Chi cron (weekly reports, drift checks)

### Database & Storage

- **Primary DB**: Supabase (Postgres 15 + pgvector)
- **Auth**: Supabase Auth (anonymous, magic link)
- **Storage**: Supabase Storage (embeddings, user exports)
- **Security**: Row-Level Security (RLS) on all tables

### AI/ML

- **LLM Primary**: Groq (llama-3.1-70b)
- **LLM Fallback**: Cerebras (same prompts, failover)
- **Embeddings**: OpenAI `text-embedding-3-small` or Voyage AI (TBD)
- **Vector Search**: pgvector (cosine similarity for drift + peer matching)

### Browser Extension

- **Platforms**: Chrome, Firefox (Manifest V3)
- **Local-First**: Events batched locally every 60s, aggregated before send
- **Privacy**: No raw URLs or content sent to backend—only domain categories and time-on-site
- **Optional**: Layer 1 (calendar + mood) works without it

---

## Architecture Highlights

### Data Flow

```
User devices (browser + extension) 
  → Go API (feature extraction, OAuth, drift scoring)
  → Supabase (Postgres + pgvector + Auth + RLS)
  → LLM (Groq → Cerebras fallback)
  → Dashboard (visualizations + insights + peer matching)
```

### 5 LLM Inference Calls

1. **Pattern Analysis** (weekly): 30-day features + mood → behavioral patterns, lag correlations, plain English summary
2. **Circadian Narrative** (weekly): Sleep consistency + debt → one micro-intervention
3. **Predictive Warning** (drift alert): Current trajectory + historical episodes → confidence score, watch signals
4. **Weekly Report** (Friday): Correlations + mood trends → 3-paragraph insight
5. **Match Context** (Constellation): Seeker + supporter anonymized patterns → warm session framing

### Privacy Guarantees

- **No raw data storage**: Extension sends aggregated domain categories, not URLs
- **User-scoped access**: RLS ensures users can only query their own rows
- **No session recording**: Peer sessions via WebRTC, no backend logging
- **Transparent AI**: Every LLM card includes "AI-generated, not a diagnosis" disclaimer
- **Crisis-safe**: Drift alerts always include professional support links

---

## Differentiation

| **vs. Mood Trackers** | Adds passive behavioral context (calendar, sleep, browsing) instead of mood-only logs |
| **vs. Wearables-Only Apps** | Multi-source correlation and LLM narrative generation |
| **vs. Therapy Platforms** | Early signal detection + peer connection, not a replacement for professionals |
| **vs. Generic Support Groups** | Peer matching on **behavioral fingerprint** (recovered-from-similar-pattern), not demographics |

**Core Bet**: Behavioral data already signals decline weeks before people seek help. Connecting that signal to peer support shortens the path from "I'm not okay" to "I found someone who gets it."

---

## Ethical Commitments

1. **Anonymous by default** — No email required; magic link or anonymous session
2. **No medical claims** — "Insight and support," never diagnosis or treatment
3. **Crisis escalation** — Always offer professional resources; never block user choice
4. **Opt-in peer matching** — Constellation is optional; opt-out anytime
5. **Transparent AI** — All LLM outputs labeled as AI-generated, not clinical advice
6. **User control** — One-click data export (JSON) and delete

---

## Build Layers (Incremental Delivery)

1. **Layer 1**: Calendar + Mood + Correlation + LLM insight *(MVP)*
2. **Layer 2**: Sleep tracking + Circadian analysis
3. **Layer 3**: Browser extension + Screen time features
4. **Layer 4**: Embeddings + Drift detection (pgvector)
5. **Layer 5**: Peer matching + WebRTC (Constellation)

Dependencies are sequential: Layer 4 requires features from 1-3; Layer 5 requires embeddings and drift from Layer 4.

---

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.22+
- Supabase account (free tier works)
- Groq API key (free at [console.groq.com](https://console.groq.com))
- Google Cloud project with Calendar API enabled (for calendar sync)

---

### Step 1 — Database

Run the migration SQL files in your **Supabase SQL Editor** (project → SQL Editor), in order:

1. `backend/internal/db/migrations/001_layer1.sql` — core tables (users, mood, calendar, features, insights)
2. `backend/internal/db/migrations/002_sleep_tables.sql` — sleep and circadian tables

---

### Step 2 — Backend environment

```bash
cd backend
cp .env.example .env
```

Open `backend/.env` and fill in your values:

```env
# Required — Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_KEY=your-service-role-key
SUPABASE_JWT_SECRET=your-jwt-secret
DATABASE_URL=postgresql://postgres:your-password@db.your-project.supabase.co:5432/postgres

# Required — Google OAuth (Calendar)
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/calendar/callback

# Required — LLM (Groq primary, Cerebras fallback)
GROQ_API_KEY=gsk_your-groq-key
CEREBRAS_API_KEY=your-cerebras-key

# Optional — Sleep providers (Layer 2)
OURA_PERSONAL_TOKEN=           # from cloud.ouraring.com/personal-access-tokens
FITBIT_CLIENT_ID=              # from dev.fitbit.com
FITBIT_CLIENT_SECRET=
```

> The server starts without Google/LLM/sleep keys — only those features will be unavailable until configured.

---

### Step 3 — Frontend environment

Create `src/.env.local` (or `.env.local` in the project root):

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

---

### Step 4 — Install frontend dependencies

```bash
# From the project root
npm install
# or
bun install
```

---

### Step 5 — Run (two terminals)

**Terminal 1 — Go backend** (port 8080):

```bash
cd backend
go run ./cmd/server
```

**Terminal 2 — Next.js frontend** (port 3000):

```bash
npm run dev
# or
bun dev
```

---

### Where things run

| URL | What |
|-----|------|
| `http://localhost:8080/connect` | Onboarding — create session, connect Google Calendar, sync |
| `http://localhost:8080/log` | Daily mood check-in |
| `http://localhost:8080/dashboard` | Main dashboard — timeline, correlations, AI insights |
| `http://localhost:8080/circadian` | Sleep health — timeline, narrative, manual sleep log |
| `http://localhost:3000` | Next.js frontend |
| `http://localhost:8080/health` | API health check |

---

### Other useful commands

```bash
# Build backend binary
cd backend && go build -o bin/server ./cmd/server

# Run Go tests
cd backend && go test ./...

# Tidy Go modules
cd backend && go mod tidy

# Build frontend for production
npm run build && npm start
```

---

## Project Structure

```
pulse/
├── src/
│   ├── app/              # Next.js App Router pages
│   ├── components/       # UI components
│   └── lib/              # API client, Supabase helpers, types
├── backend/
│   ├── cmd/server/       # Entry point (main.go)
│   ├── internal/
│   │   ├── api/          # HTTP handlers (calendar, mood, dashboard, insights, sleep, circadian)
│   │   ├── collectors/
│   │   │   ├── google/   # Google Calendar OAuth + event fetching
│   │   │   └── sleep/    # Oura, Fitbit, Google Fit, manual entry
│   │   ├── correlation/  # Pearson, lagged correlation, significance, matrix
│   │   ├── db/           # Database layer + migrations/
│   │   ├── features/
│   │   │   ├── calendar/ # Meeting density, fragmentation, focus blocks
│   │   │   └── circadian/# Rhythm consistency, sleep debt, social jetlag
│   │   ├── llm/          # Groq/Cerebras clients, pattern analysis, circadian narrative
│   │   ├── middleware/   # Auth (JWT), CORS, logger, recover
│   │   ├── pipeline/     # Sync orchestration, cron scheduler
│   │   └── config/       # Environment config
│   ├── web/              # Standalone HTML UI served by Go
│   │   ├── connect.html  # Onboarding (session → OAuth → sync)
│   │   ├── log.html      # Daily mood check-in
│   │   ├── dashboard.html# Main dashboard (Chart.js timeline + insights)
│   │   └── circadian.html# Sleep health + narrative
│   ├── .env.example      # All environment variables documented
│   ├── Makefile          # dev, build, test, tidy targets
│   └── go.mod
├── public/               # Static assets
└── package.json
```

---

## License

MIT

---

## Disclaimer

**Pulse is not a medical device or diagnostic tool.** It provides behavioral insights and peer support, not professional mental health treatment. Always consult licensed healthcare providers for medical advice. In crisis, contact emergency services or a crisis hotline immediately.
