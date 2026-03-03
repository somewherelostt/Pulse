<div align="center">

# Pulse

**The Mental Health Platform That Sees Burnout Before You Do**

*Behavioral intelligence that detects, understands, and connects — powered by the patterns in your calendar, sleep, and mood*

[Report Bug](https://github.com/yourusername/pulse/issues) · [Request Feature](https://github.com/yourusername/pulse/issues)

</div>

---

## The "Why" (Our Story)

We've all been there: you're burning out, but you don't realize it until it's too late. Your calendar looks busy, not bad. Your sleep seems "mostly okay." Your mood feels "fine."

But somewhere in the data — the meeting density, the fragmented calendar, the slowly shifting sleep schedule — the truth was hiding. Your behavior knew you were struggling weeks before you did.

We built Pulse because we wanted a system that doesn't wait for you to crash. We wanted something that reads the invisible patterns across your calendar, sleep, and mood — and tells you what's happening before it becomes a crisis.

This isn't another mood tracker that asks "how do you feel?" and stops there. This is behavioral intelligence that finds the two-day lag between late-night screen time and mood crashes. That matches you with peers who recovered from your exact pattern. That generates plain-English insights from data you already create.

**Mental health tools are built around symptoms. Pulse is built around signals.**

<div align="center">

![The Gap](public/1.png)

</div>

---

## The Problem — Four Real Stories

### 📅 Arjun, 27 — Software Engineer, Bengaluru

> *"78 meetings in 19 days. No single one was unreasonable. His calendar looked busy, not bad."*

Arjun submitted his resignation letter at 2am, six weeks later. His recovery windows between high-load days had quietly shrunk to zero. No tool flagged the trajectory — not his manager, not any app.

**What he needed:** Something that reads meeting density + recovery gaps + mood drift *together*, before the crash.

---

### 🛌 Mira, 31 — PhD Candidate, Amsterdam

> *"Her circadian rhythm shifted 3.2 hours over five months. She thought it was dissertation stress."*

Her wearable gave weekly averages — numbers with no narrative. The hidden pattern: every 45-min rightward sleep shift = 22% drop in cognitive clarity next day. She owned the data. She just couldn't read it.

**What she needed:** Circadian analysis that generates a *plain-English narrative* with one actionable intervention — not another chart.

---

### 🤝 James, 34 — Co-founder, London

> *"He knew all the warning signs. He still couldn't stop. What he needed wasn't information."*

He searched for someone who had been in the same place — same calendar collapse, same sleep data. Not a therapist following a protocol. Not a generic support group. He called it "finding a needle in a haystack while the haystack is on fire."

**What he needed:** Anonymous peer matching on *behavioral fingerprint* — connected to someone who recovered from a statistically similar pattern.

---

### 🧩 Priya, 29 — UX Designer, Mumbai

> *"Mood: fine. Sleep: mostly okay. Work: manageable. But something was wrong."*

She kept a mood journal — it captured states, not causes. Hidden thread: late-night screen time spike → fragmented next-day calendar → mood dip 48h later. A two-day lag that no journal can surface on its own.

**What she needed:** Cross-source behavioral correlation that finds lagged relationships human introspection misses.

---

## How Pulse Solves This

### 📅 For Arjun — Calendar Intelligence + Drift Detection

Pulse doesn't just count meetings. It tracks:

- 🗓️ **Meeting density score** — total event hours, normalized daily
- 🧩 **Calendar fragmentation** — average uninterrupted focus blocks
- 🌙 **Boundary violations** — after-hours events, frequency and trend
- ⚡ **Drift alert** — when your 14-day window diverges from your baseline, Pulse flags it with confidence score, watch signals, and always-visible professional resources

**The system notices before Arjun does.**

---

### 🛌 For Mira — Sleep & Circadian Analysis

Wearables give you numbers. Pulse gives you a narrative:

- 📡 **Device integrations**: Oura Ring, Fitbit, Google Fit, manual entry
- 🔄 **Circadian consistency** — nightly variation in sleep onset and wake time
- 💤 **Sleep debt** — rolling deficit vs. personal baseline
- 🌞 **Social jetlag** — weekday vs. weekend sleep timing divergence
- 🧠 **LLM circadian narrative** — 7-day reading → one specific micro-intervention, in plain English

**Not more data. Language for what the data has been saying.**

---

### � For Priya — Cross-Source Correlation + LLM Insight

The human brain can't detect two-day lags across three data sources. Pulse can:

- 📊 **Pearson + lagged correlation** across sleep × calendar × mood × screen time
- 🔍 Detects 48-hour lag relationships human introspection misses
- 🧠 **Weekly LLM pattern analysis** — 30 days of features → strongest correlation, lag relationship, 3-paragraph narrative, one time-specific intervention point

### 🤝 For James — Constellation: Peer-to-Peer Support

Most support systems match you by age, location, or diagnosis. Constellation matches you by behavioral pattern — the data signature of what you're actually experiencing.

**How it works:**

1. **Safety First** — Before entering the pool, we check for crisis indicators (severe drift scores, sustained low mood, crisis keywords). If detected, we show professional resources immediately.
2. **Behavioral Fingerprinting** — Your last 14 days of data (calendar patterns, sleep consistency, mood trajectory, circadian rhythm, screen time, correlation strength) gets encoded into a 6-dimensional vector stored in pgvector.
3. **Smart Matching** — When you opt in, we run cosine similarity against the peer pool. Scoring is 70% behavioral similarity ("your patterns look like mine did") + 30% temporal overlap ("I was where you are 3-6 months ago").
4. **AI-Generated Context** — Before the session, an LLM reads both fingerprints and generates 2-3 sentences of warm context: "You're both experiencing calendar fragmentation with late-night recovery attempts. One of you has been through this pattern and came out the other side."
5. **Anonymous WebRTC Session** — 45-minute peer-to-peer video/audio call. No server recording. No content logging. Just two people, one pattern, zero judgment.
6. **Ephemeral by Design** — Session ends, connection closes, no chat history, no replay. Only metadata logged: session start/end, peer IDs (hashed), rating (1-5).

**Why this matters:** Therapists follow protocols. Hotlines follow scripts. But someone who survived your *exact* behavioral collapse? They know the 2am thoughts, the meeting dread, the "I'm fine" lies. That's not a hotline. That's a lifeline.

<div align="center">

![How Pulse Works](public/2.png)

</div>

---

## What Pulse Does

**Layer 1: Behavioral Intelligence** — Calendar sync (Google), mood check-ins, meeting density analysis, boundary violation tracking, correlation detection across all data sources.

**Layer 2: Sleep & Circadian** — Integrations with Oura, Fitbit, Google Fit, or manual logging. Tracks sleep debt, circadian consistency, social jetlag. Generates plain-English interventions.

**Layer 3: Screen Time** — Browser extension (Chrome/Firefox) that tracks domain categories only, never URLs. Local-first, batched sync every 60s.

**Layer 4: Drift Detection** — 6-dimensional behavioral fingerprint per day. Compares your 14-day window to baseline via pgvector. Flags drift before crisis.

**Layer 5: Constellation (P2P Matching)** — Behavioral peer matching with safety checks, WebRTC sessions, zero content retention. See detailed section above.

<div align="center">

![Dashboard](public/3.png)

*One dashboard. Four data sources. Finally connected.*

</div>

---

<div align="center">

![Constellation](public/4.png)

*The hardest part isn't the pattern. It's feeling like nobody else has been through it.*

</div>

---

## Tech Stack

### Frontend

- **Next.js 14** (App Router) — React framework
- **Tailwind CSS** + CSS Modules — styling
- **Chart.js** — timeline visualizations
- **D3.js** — correlation matrices
- **Three.js** — 3D landing page animations
- **Supabase Auth** — anonymous + magic link auth

### Backend

- **Go 1.22** — API server
- **Chi router** — HTTP routing
- **Supabase Postgres 15** — main database
- **pgvector** — vector similarity for peer matching
- **JWT (HS256)** — authentication tokens
- **Row-Level Security (RLS)** — database-level access control
- **WebSocket (gorilla/websocket)** — Constellation signaling

### AI & Integrations

- **Groq (llama-3.1-70b)** — primary LLM (pattern analysis, insights)
- **Cerebras** — LLM failover
- **Google Calendar API** — calendar sync
- **Oura Cloud API** — sleep data
- **Fitbit OAuth** — sleep data
- **Google Fit** — sleep data

### Infrastructure

- **Render** — deployment (API + cron + WebSocket)
- **Supabase** — managed Postgres + Auth

### Database Schema

```
001_layer1.sql        users, oauth_tokens, calendar_events, daily_features,
                      mood_logs, llm_insights, sync_log
002_sleep_tables.sql  sleep_sessions, circadian_features, circadian_insights
003_constellation.sql peer_pool (vector(6)), constellation_sessions,
                      constellation_session_log
```

**LLM Usage (5 specific calls):**

1. **Pattern Analysis** (weekly) — 30-day behavioral narrative
2. **Circadian Narrative** (weekly) — sleep intervention suggestions
3. **Predictive Warning** (on drift) — confidence score + watch signals
4. **Weekly Report** (Friday) — 3-paragraph correlation insight
5. **Match Context** (Constellation) — warm P2P session framing

All AI outputs labeled "AI-generated, not clinical advice."

---

## Privacy & Ethics

- 🔐 **Anonymous by default** — no email required
- 🚫 **No raw data storage** — browser extension sends domain categories only, never URLs or page content
- 📵 **No session recording** — peer sessions use WebRTC, zero content logging, ephemeral by design
- 🔒 **Row-Level Security** — database enforces that users can only access their own data
- 🏷️ **Transparent AI** — every LLM output labeled "AI-generated, not clinical advice"
- 🆘 **Crisis-safe** — professional resources visible in every drift alert, never paywalled
- 🗑️ **Full data export + deletion** — one-click JSON export and account deletion

**We provide insight and peer support, not diagnosis or treatment.**

<div align="center">

![Privacy](public/5.png)

*Everything runs on your data. Nothing runs on your identity.*

</div>

---

## Why Pulse Is Different

| Feature | Pulse | Mood Trackers | Wearables | Therapy Apps |
|---|:---:|:---:|:---:|:---:|
| Passive behavioral data | ✅ | ❌ | partial | ❌ |
| Cross-source correlation | ✅ | ❌ | ❌ | ❌ |
| LLM narrative (not just charts) | ✅ | ❌ | ❌ | ❌ |
| Drift detection before crisis | ✅ | ❌ | ❌ | ❌ |
| Peer match on behavioral pattern | ✅ | ❌ | ❌ | ❌ |
| Anonymous — no account required | ✅ | ❌ | ❌ | ❌ |

---

## Getting Started

### Prerequisites

- **Node.js 18+** and **Go 1.22+**
- **Supabase account** (free tier works)
- **Groq API key** — free at [console.groq.com](https://console.groq.com)
- **Google Cloud project** with Calendar API enabled

### Step 1 — Database Setup 🗄️

Run these migrations in your **Supabase SQL Editor**, in order:

```sql
backend/internal/db/migrations/001_layer1.sql
backend/internal/db/migrations/002_sleep_tables.sql
backend/internal/db/migrations/003_constellation.sql
```

💡 **Tip**: For 30 days of demo data, run `seed-demo.sql` after your first sign-in.

### Step 2 — Backend Configuration ⚙️

1. Copy the example environment file:

   ```bash
   cd backend && cp .env.example .env
   ```

2. Fill in your environment variables:

   ```env
   SUPABASE_URL=https://your-project.supabase.co
   SUPABASE_ANON_KEY=
   SUPABASE_SERVICE_KEY=
   SUPABASE_JWT_SECRET=
   DATABASE_URL=postgresql://postgres:password@db.your-project.supabase.co:5432/postgres

   GOOGLE_CLIENT_ID=
   GOOGLE_CLIENT_SECRET=
   GOOGLE_REDIRECT_URI=http://localhost:8080/api/v1/calendar/callback

   GROQ_API_KEY=
   CEREBRAS_API_KEY=

   # Optional — sleep providers
   OURA_PERSONAL_TOKEN=
   FITBIT_CLIENT_ID=
   FITBIT_CLIENT_SECRET=
   ```

### Step 3 — Frontend Configuration 🖥️

Create `.env.local` in the root directory:

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=
```

### Step 4 — Running the App 🚀

You can run both the backend and frontend together:

```bash
# Install dependencies
npm install

# Terminal 1 — Go backend (port 8080)
cd backend && go run ./cmd/server

# Terminal 2 — Next.js frontend (port 3000)
npm run dev
```

Visit `http://localhost:3000` 🎉

**Key Routes:**

- `/` — Landing page
- `/auth` — Sign in
- `/dashboard` — Main dashboard (timeline, metrics, insights)
- `/dashboard/sleep` — Sleep & circadian analysis
- `/dashboard/constellation` — Peer matching
- `/log` — Mood check-in

---

## Project Structure

```
pulse/
├── src/                           # Next.js Frontend
│   ├── app/
│   │   ├── dashboard/             Main dashboard, sleep, circadian, constellation
│   │   ├── log/                   Mood check-in
│   │   └── auth/                  Authentication
│   ├── components/                UI components (Hero, charts, forms)
│   └── lib/                       API client, types, utilities
│
├── backend/                       # Go API
│   ├── cmd/server/                Entry point
│   ├── internal/
│   │   ├── api/                   REST endpoints
│   │   ├── constellation/         Peer matching engine
│   │   ├── correlation/           Statistical analysis
│   │   ├── features/              Calendar & circadian extraction
│   │   ├── llm/                   AI insight generation
│   │   └── db/                    Database + migrations
│   └── web/                       Standalone HTML views
│
└── seed-demo.sql                  30-day demo data
```

---

## Disclaimer

**Pulse is not a medical device.** It provides behavioral insight and anonymous peer connection — not diagnosis or treatment.

Always consult licensed healthcare providers for clinical concerns. If you are in crisis, contact emergency services or a crisis line immediately — Pulse will always show you how.

---

<div align="center">

*Built for the Arjuns and Miras and Jameses and Priyas*  
*— the ones whose data knew before they did*

**Made with ❤️ for mental health**

</div>
