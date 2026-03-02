-- ============================================================
-- Pulse Demo Seed — abumaaz2004@gmail.com
-- Self-contained: creates the auth user if not already present.
-- Run this in the Supabase SQL Editor.
-- ============================================================

DO $$
DECLARE
  v_supabase_uid  UUID;
  v_user_id       UUID;
  v_today         DATE := CURRENT_DATE;
  v_day           DATE;
  i               INT;

  -- Mood pattern: starts okay (6-7), dips during days 15-22, recovers
  mood_scores     INT[]     := ARRAY[7,6,7,6,8,7,6,5,7,8,7,6,5,6,4,4,3,4,5,5,6,5,7,6,7,8,7,6,7,8];
  energy_scores   INT[]     := ARRAY[6,7,6,7,7,6,5,6,6,7,6,5,5,5,3,3,4,4,5,5,5,6,6,7,7,7,6,7,7,8];
  anxiety_scores  INT[]     := ARRAY[4,3,4,3,3,4,5,5,4,3,4,5,6,6,8,8,7,7,6,6,5,5,4,4,3,3,4,3,3,2];

  -- Calendar features: high meeting load weeks 2-3, lighter weeks 1 and 4
  meeting_density FLOAT[]   := ARRAY[0.42,0.48,0.51,0.44,0.38,0.55,0.60,0.62,0.58,0.61,0.65,0.68,0.72,0.70,0.75,0.78,0.73,0.71,0.68,0.65,0.60,0.58,0.55,0.52,0.48,0.44,0.40,0.42,0.45,0.43];
  meeting_counts  INT[]     := ARRAY[4,5,5,4,3,5,6,6,6,6,7,7,8,7,8,9,8,8,7,7,6,6,5,5,4,4,4,4,5,4];
  focus_mins      FLOAT[]   := ARRAY[78,72,68,75,90,65,60,55,62,58,52,50,45,48,42,38,40,42,50,52,58,60,65,70,75,80,85,78,72,76];
  frag_scores     FLOAT[]   := ARRAY[0.42,0.48,0.52,0.45,0.38,0.55,0.60,0.65,0.60,0.62,0.68,0.70,0.75,0.72,0.80,0.82,0.78,0.76,0.70,0.68,0.62,0.58,0.55,0.50,0.46,0.42,0.38,0.40,0.44,0.42];
  after_hours     FLOAT[]   := ARRAY[22,28,35,20,15,40,45,52,48,50,60,65,72,70,85,90,82,78,70,65,55,50,42,38,30,25,20,22,28,24];
  solo_pct        FLOAT[]   := ARRAY[0.62,0.55,0.52,0.58,0.65,0.50,0.44,0.42,0.45,0.42,0.38,0.35,0.32,0.34,0.28,0.26,0.30,0.32,0.38,0.40,0.45,0.48,0.52,0.55,0.60,0.63,0.67,0.62,0.58,0.62];
  b2b_counts      INT[]     := ARRAY[1,2,2,1,0,2,3,3,3,3,4,4,5,4,5,6,5,5,4,4,3,3,2,2,1,1,1,1,2,1];
  recovery_mins   FLOAT[]   := ARRAY[35,28,25,32,42,22,18,15,20,17,12,10,8,10,5,4,6,8,12,15,18,20,25,28,32,38,42,36,30,34];
  attendee_avg    FLOAT[]   := ARRAY[3.2,3.8,4.0,3.5,2.8,4.2,4.8,5.0,4.6,4.8,5.2,5.5,5.8,5.5,6.0,6.2,5.8,5.6,5.2,4.8,4.4,4.2,3.8,3.5,3.2,3.0,2.8,3.0,3.4,3.2];

BEGIN
  -- 1. Find or create the auth user
  SELECT id INTO v_supabase_uid
  FROM auth.users
  WHERE email = 'abumaaz2004@gmail.com'
  LIMIT 1;

  IF v_supabase_uid IS NULL THEN
    v_supabase_uid := gen_random_uuid();
    INSERT INTO auth.users (
      id, email, aud, role,
      email_confirmed_at,
      created_at, updated_at,
      raw_app_meta_data, raw_user_meta_data,
      is_sso_user, encrypted_password
    ) VALUES (
      v_supabase_uid,
      'abumaaz2004@gmail.com',
      'authenticated',
      'authenticated',
      NOW(),
      NOW(), NOW(),
      '{"provider":"email","providers":["email"]}'::jsonb,
      '{}'::jsonb,
      false,
      ''
    );
    RAISE NOTICE 'Created auth user: %', v_supabase_uid;
  ELSE
    RAISE NOTICE 'Found existing auth user: %', v_supabase_uid;
  END IF;

  -- 2. Upsert into public.users
  INSERT INTO public.users (supabase_uid, timezone, work_start_hour, work_end_hour, onboarding_done, consent_calendar, updated_at)
  VALUES (v_supabase_uid, 'Asia/Kolkata', 9, 19, TRUE, TRUE, NOW())
  ON CONFLICT (supabase_uid) DO UPDATE SET
    timezone = 'Asia/Kolkata',
    work_start_hour = 9,
    work_end_hour = 19,
    onboarding_done = TRUE,
    consent_calendar = TRUE,
    updated_at = NOW();

  SELECT id INTO v_user_id
  FROM public.users
  WHERE supabase_uid = v_supabase_uid;

  RAISE NOTICE 'public.users id: %', v_user_id;

  -- 3. Insert 30 days of mood logs (day 1 = 29 days ago, day 30 = yesterday)
  FOR i IN 1..30 LOOP
    v_day := v_today - (31 - i);

    INSERT INTO public.mood_logs (user_id, date, score, energy, anxiety, note, tags, logged_at)
    VALUES (
      v_user_id,
      v_day,
      mood_scores[i],
      energy_scores[i],
      anxiety_scores[i],
      CASE
        WHEN mood_scores[i] <= 4 THEN 'Feeling stretched thin. Too many meetings, no space to breathe.'
        WHEN mood_scores[i] >= 7 THEN 'Good day. Got into a solid flow in the afternoon.'
        ELSE NULL
      END,
      CASE
        WHEN mood_scores[i] <= 4 THEN ARRAY['overloaded', 'fragmented']
        WHEN mood_scores[i] >= 7 THEN ARRAY['focused', 'energized']
        ELSE ARRAY[]::TEXT[]
      END,
      v_day + INTERVAL '20 hours'
    )
    ON CONFLICT (user_id, date) DO UPDATE SET
      score = EXCLUDED.score,
      energy = EXCLUDED.energy,
      anxiety = EXCLUDED.anxiety,
      note = EXCLUDED.note,
      tags = EXCLUDED.tags,
      logged_at = EXCLUDED.logged_at;
  END LOOP;

  -- 4. Insert 30 days of daily features
  FOR i IN 1..30 LOOP
    v_day := v_today - (31 - i);

    INSERT INTO public.daily_features (
      user_id, date,
      meeting_density_pct, meeting_count,
      avg_focus_block_mins, max_focus_block_mins,
      fragmentation_score, after_hours_mins,
      weekend_meeting_mins, back_to_back_count,
      avg_recovery_mins, attendee_avg,
      solo_time_pct, source_event_count,
      computed_at
    )
    VALUES (
      v_user_id, v_day,
      meeting_density[i], meeting_counts[i],
      focus_mins[i], focus_mins[i] * 1.8,
      frag_scores[i], after_hours[i],
      CASE WHEN EXTRACT(DOW FROM v_day) IN (0,6) THEN after_hours[i] * 0.4 ELSE 0 END,
      b2b_counts[i],
      recovery_mins[i], attendee_avg[i],
      solo_pct[i], meeting_counts[i] + 2,
      NOW()
    )
    ON CONFLICT (user_id, date) DO UPDATE SET
      meeting_density_pct = EXCLUDED.meeting_density_pct,
      meeting_count = EXCLUDED.meeting_count,
      avg_focus_block_mins = EXCLUDED.avg_focus_block_mins,
      max_focus_block_mins = EXCLUDED.max_focus_block_mins,
      fragmentation_score = EXCLUDED.fragmentation_score,
      after_hours_mins = EXCLUDED.after_hours_mins,
      back_to_back_count = EXCLUDED.back_to_back_count,
      avg_recovery_mins = EXCLUDED.avg_recovery_mins,
      attendee_avg = EXCLUDED.attendee_avg,
      solo_time_pct = EXCLUDED.solo_time_pct,
      source_event_count = EXCLUDED.source_event_count,
      computed_at = NOW();
  END LOOP;

  -- 5. Insert a sync log so dashboard shows "calendar connected"
  INSERT INTO public.sync_log (user_id, provider, status, events_fetched, synced_at)
  VALUES (v_user_id, 'google', 'success', 187, NOW() - INTERVAL '2 hours')
  ON CONFLICT (user_id, provider) DO UPDATE SET
    status = 'success',
    events_fetched = 187,
    synced_at = NOW() - INTERVAL '2 hours';

  -- 6. Insert a pre-generated LLM insight
  INSERT INTO public.llm_insights (user_id, insight_type, week_start, content, model_used, created_at)
  VALUES (
    v_user_id,
    'pattern_analysis',
    DATE_TRUNC('week', CURRENT_DATE),
    jsonb_build_object(
      'summary', 'Your behavioral data shows a clear burnout spiral pattern over the past 3 weeks. Meeting density above 65% consistently preceded mood drops by 2–3 days — this lag is the fingerprint Pulse uses to predict your state before you feel it.',
      'recommendation', 'Protect two 90-minute blocks daily as non-negotiable focus time. The data shows your mood recovers fastest when back-to-back meetings drop below 3 per day.',
      'patterns', jsonb_build_array(
        jsonb_build_object('feature', 'meeting_density_pct', 'lag_days', 2, 'direction', 'negative', 'confidence', 0.78, 'severity', 'high', 'plain_english', 'High meeting load predicts mood dip 2 days later'),
        jsonb_build_object('feature', 'after_hours_mins', 'lag_days', 1, 'direction', 'negative', 'confidence', 0.71, 'severity', 'moderate', 'plain_english', 'Working late correlates with lower energy next day'),
        jsonb_build_object('feature', 'avg_focus_block_mins', 'lag_days', 0, 'direction', 'positive', 'confidence', 0.65, 'severity', 'moderate', 'plain_english', 'Longer focus blocks associate with better mood same day')
      ),
      'data_quality_note', 'Based on 30 days of calendar + mood data. Correlations computed at 7-day rolling window.',
      'disclaimer', 'This is behavioral pattern analysis, not medical advice. For mental health support, please speak with a professional.',
      'generated_at', NOW()::TEXT,
      'model_used', 'llama-3.3-70b-versatile'
    ),
    'llama-3.3-70b-versatile',
    NOW()
  )
  ON CONFLICT (user_id, insight_type, week_start) DO UPDATE SET
    content = EXCLUDED.content,
    model_used = EXCLUDED.model_used,
    created_at = NOW();

  RAISE NOTICE 'Seed complete for abumaaz2004@gmail.com (user_id: %)', v_user_id;
END $$;
