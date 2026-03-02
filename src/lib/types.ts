export type User = {
  id: string;
  supabase_uid: string;
  timezone: string;
  onboarding_done: boolean;
  consent_calendar: boolean;
  work_start_hour: number;
  work_end_hour: number;
};

export type MoodLog = {
  id: string;
  date: string;
  score: number;
  energy: number | null;
  anxiety: number | null;
  note: string | null;
  tags: string[];
  logged_at: string;
};

export type DailyFeatures = {
  date: string;
  meeting_density_pct: number | null;
  meeting_count: number | null;
  avg_focus_block_mins: number | null;
  max_focus_block_mins: number | null;
  fragmentation_score: number | null;
  after_hours_mins: number | null;
  back_to_back_count: number | null;
  avg_recovery_mins: number | null;
  solo_time_pct: number | null;
  attendee_avg: number | null;
};

export type TimelinePoint = DailyFeatures & {
  mood_score: number | null;
  mood_energy: number | null;
};

export type CorrelationResult = {
  feature_name: string;
  lag_days: number;
  correlation: number;
  direction: "positive" | "negative";
  strength: "weak" | "moderate" | "strong";
  plain_english: string;
};

export type DetectedPattern = {
  feature: string;
  lag_days: number;
  direction: "positive" | "negative";
  confidence: number;
  plain_english: string;
  severity: "low" | "moderate" | "high";
};

export type InsightData = {
  patterns: DetectedPattern[];
  summary: string;
  recommendation: string;
  data_quality_note: string;
  disclaimer: string;
  generated_at: string;
  model_used: string;
};

export type DashboardResponse = {
  user: Pick<User, "onboarding_done" | "timezone"> & {
    calendar_connected: boolean;
  };
  metrics: {
    avg_meeting_density_pct: number;
    meeting_density_trend: "up" | "down" | "stable";
    avg_fragmentation_score: number;
    avg_mood_score: number;
    mood_trend: "up" | "down" | "stable";
    data_days: number;
  };
  timeline: TimelinePoint[];
  top_correlations: CorrelationResult[];
  latest_insight: InsightData | null;
  sync_status: {
    last_synced_at: string | null;
    events_fetched: number;
    status: "success" | "failed" | "never";
  };
};
