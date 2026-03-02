-- Pulse Layer 2 schema — Sleep data and circadian features
-- Run in Supabase SQL editor after 001_layer1.sql

-- Raw sleep sessions from any provider
create table if not exists public.sleep_sessions (
  id              uuid primary key default uuid_generate_v4(),
  user_id         uuid not null references public.users(id) on delete cascade,
  provider        text not null,  -- 'oura' | 'fitbit' | 'googlefit' | 'manual'
  session_date    date not null,  -- the calendar date the night started
  bedtime         timestamptz,
  wake_time       timestamptz,
  total_sleep_mins int,
  rem_mins        int,
  deep_mins       int,
  light_mins      int,
  awake_mins      int,
  sleep_score     int,            -- provider native score 0-100
  hrv_rmssd       float,          -- ms
  resting_hr      float,          -- bpm
  source_id       text,           -- provider's own session ID
  raw_json        jsonb,
  fetched_at      timestamptz not null default now(),
  unique(user_id, provider, session_date)
);

-- Extracted circadian features per day (aggregated from sleep_sessions)
create table if not exists public.circadian_features (
  id                      uuid primary key default uuid_generate_v4(),
  user_id                 uuid not null references public.users(id) on delete cascade,
  date                    date not null,
  sleep_duration_mins     float,
  sleep_efficiency_pct    float,   -- total_sleep / (wake - bedtime)
  sleep_debt_mins         float,   -- rolling 7-day vs target (480 min)
  mid_sleep_hour          float,   -- hour of day at midpoint of sleep (e.g. 3.25)
  rhythm_consistency_pct  float,   -- 0-1, stdev-based vs 7-day mean midpoint
  social_jetlag_mins      float,   -- |workday mid_sleep - weekend mid_sleep| * 60
  rem_pct                 float,
  deep_pct                float,
  hrv_rmssd               float,
  resting_hr              float,
  sleep_score             float,   -- normalized 0-100
  light_hygiene_score     float,   -- 0-100 based on timing
  computed_at             timestamptz not null default now(),
  unique(user_id, date)
);

-- LLM circadian narratives (weekly)
create table if not exists public.circadian_insights (
  id              uuid primary key default uuid_generate_v4(),
  user_id         uuid not null references public.users(id) on delete cascade,
  week_start      date not null,
  narrative       text not null,
  interventions   jsonb,           -- [{title, description, priority}]
  model_used      text not null,
  created_at      timestamptz not null default now(),
  unique(user_id, week_start)
);

-- RLS
alter table public.sleep_sessions enable row level security;
alter table public.circadian_features enable row level security;
alter table public.circadian_insights enable row level security;

-- RLS policies (drop if exists to allow re-running migration)
drop policy if exists "sleep_own" on public.sleep_sessions;
create policy "sleep_own" on public.sleep_sessions
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

drop policy if exists "circadian_features_own" on public.circadian_features;
create policy "circadian_features_own" on public.circadian_features
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

drop policy if exists "circadian_insights_own" on public.circadian_insights;
create policy "circadian_insights_own" on public.circadian_insights
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

-- Indexes
create index if not exists idx_sleep_sessions_user_date
  on public.sleep_sessions(user_id, session_date desc);

create index if not exists idx_circadian_features_user_date
  on public.circadian_features(user_id, date desc);

create index if not exists idx_circadian_insights_user_week
  on public.circadian_insights(user_id, week_start desc);

-- Extend oauth_tokens for Fitbit (provider='fitbit' row will be created on auth)
-- No schema change needed since oauth_tokens already has a generic provider column.

-- Sync log entries for sleep providers use the existing sync_log table.
-- Provider values: 'oura', 'fitbit', 'googlefit'
