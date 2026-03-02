-- Pulse Layer 1 schema for Supabase Postgres
-- Run in Supabase SQL editor

-- Extensions
create extension if not exists "uuid-ossp";

-- Users (mirrors Supabase auth.users)
create table if not exists public.users (
  id                uuid primary key default uuid_generate_v4(),
  supabase_uid      uuid unique not null,
  created_at        timestamptz not null default now(),
  updated_at        timestamptz not null default now(),
  timezone          text not null default 'UTC',
  onboarding_done   boolean not null default false,
  consent_calendar  boolean not null default false,
  work_start_hour   int not null default 9,
  work_end_hour     int not null default 18
);

-- OAuth tokens per provider
create table if not exists public.oauth_tokens (
  id              uuid primary key default uuid_generate_v4(),
  user_id         uuid not null references public.users(id) on delete cascade,
  provider        text not null,
  access_token    text not null,
  refresh_token   text,
  token_expiry    timestamptz,
  scope           text,
  created_at      timestamptz not null default now(),
  updated_at      timestamptz not null default now(),
  unique(user_id, provider)
);

-- Raw calendar events (privacy: title stored as hash only)
create table if not exists public.raw_calendar_events (
  id                uuid primary key default uuid_generate_v4(),
  user_id           uuid not null references public.users(id) on delete cascade,
  google_event_id   text not null,
  title_hash        text,
  start_time        timestamptz not null,
  end_time          timestamptz not null,
  attendee_count    int not null default 0,
  is_all_day        boolean not null default false,
  is_recurring      boolean not null default false,
  is_after_hours    boolean not null default false,
  is_weekend        boolean not null default false,
  duration_mins     float not null default 0,
  fetched_at        timestamptz not null default now(),
  unique(user_id, google_event_id)
);

-- Extracted behavioral features per day
create table if not exists public.daily_features (
  id                      uuid primary key default uuid_generate_v4(),
  user_id                 uuid not null references public.users(id) on delete cascade,
  date                    date not null,
  meeting_density_pct     float,
  meeting_count           int,
  avg_focus_block_mins    float,
  max_focus_block_mins    float,
  fragmentation_score     float,
  after_hours_mins        float,
  weekend_meeting_mins    float,
  back_to_back_count      int,
  avg_recovery_mins       float,
  attendee_avg            float,
  solo_time_pct           float,
  source_event_count      int,
  computed_at             timestamptz not null default now(),
  unique(user_id, date)
);

-- Daily mood logs
create table if not exists public.mood_logs (
  id          uuid primary key default uuid_generate_v4(),
  user_id     uuid not null references public.users(id) on delete cascade,
  date        date not null,
  score       int not null check (score between 1 and 10),
  energy      int check (energy between 1 and 10),
  anxiety     int check (anxiety between 1 and 10),
  note        text check (char_length(note) <= 500),
  tags        text[],
  logged_at   timestamptz not null default now(),
  unique(user_id, date)
);

-- LLM generated insights (cached)
create table if not exists public.llm_insights (
  id              uuid primary key default uuid_generate_v4(),
  user_id         uuid not null references public.users(id) on delete cascade,
  insight_type    text not null,
  week_start      date not null,
  content         jsonb not null,
  model_used      text not null,
  prompt_version  text not null default 'v1',
  created_at      timestamptz not null default now(),
  unique(user_id, insight_type, week_start)
);

-- Sync status per connector
create table if not exists public.sync_log (
  id              uuid primary key default uuid_generate_v4(),
  user_id         uuid not null references public.users(id) on delete cascade,
  provider        text not null,
  synced_at       timestamptz not null default now(),
  events_fetched  int not null default 0,
  status          text not null default 'success',
  error_message   text,
  unique(user_id, provider)
);

-- RLS: enable on all tables
alter table public.users enable row level security;
alter table public.oauth_tokens enable row level security;
alter table public.raw_calendar_events enable row level security;
alter table public.daily_features enable row level security;
alter table public.mood_logs enable row level security;
alter table public.llm_insights enable row level security;
alter table public.sync_log enable row level security;

-- RLS policies (drop if exists to allow re-running migration)
drop policy if exists "users_own" on public.users;
create policy "users_own" on public.users
  for all using (supabase_uid = auth.uid());

drop policy if exists "oauth_own" on public.oauth_tokens;
create policy "oauth_own" on public.oauth_tokens
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

drop policy if exists "calendar_own" on public.raw_calendar_events;
create policy "calendar_own" on public.raw_calendar_events
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

drop policy if exists "features_own" on public.daily_features;
create policy "features_own" on public.daily_features
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

drop policy if exists "mood_own" on public.mood_logs;
create policy "mood_own" on public.mood_logs
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

drop policy if exists "insights_own" on public.llm_insights;
create policy "insights_own" on public.llm_insights
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

drop policy if exists "sync_own" on public.sync_log;
create policy "sync_own" on public.sync_log
  for all using (user_id in (
    select id from public.users where supabase_uid = auth.uid()
  ));

-- Indexes
create index if not exists idx_daily_features_user_date
  on public.daily_features(user_id, date desc);

create index if not exists idx_mood_logs_user_date
  on public.mood_logs(user_id, date desc);

create index if not exists idx_raw_calendar_user_start
  on public.raw_calendar_events(user_id, start_time desc);

create index if not exists idx_llm_insights_user_week
  on public.llm_insights(user_id, week_start desc);
