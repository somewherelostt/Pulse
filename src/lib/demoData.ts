import type { DashboardResponse } from "./types";

/**
 * Rich demo data for showcasing Pulse features
 * Shows a realistic burnout pattern over 30 days
 */
export function getDemoDashboardData(): DashboardResponse {
  const now = new Date();
  const generateTimeline = () => {
    const timeline = [];

    // Mood pattern: starts okay (6-7), dips during days 15-22, recovers
    const moodScores = [
      7, 6, 7, 6, 8, 7, 6, 5, 7, 8, 7, 6, 5, 6, 4, 4, 3, 4, 5, 5, 6, 5, 7, 6, 7,
      8, 7, 6, 7, 8,
    ];
    const energyScores = [
      6, 7, 6, 7, 7, 6, 5, 6, 6, 7, 6, 5, 5, 5, 3, 3, 4, 4, 5, 5, 5, 6, 6, 7, 7,
      7, 6, 7, 7, 8,
    ];

    // Calendar features: high meeting load weeks 2-3, lighter weeks 1 and 4
    const meetingDensity = [
      0.42, 0.48, 0.51, 0.44, 0.38, 0.55, 0.6, 0.62, 0.58, 0.61, 0.65, 0.68,
      0.72, 0.7, 0.75, 0.78, 0.73, 0.71, 0.68, 0.65, 0.6, 0.58, 0.55, 0.52,
      0.48, 0.44, 0.4, 0.42, 0.45, 0.43,
    ];
    const meetingCounts = [
      4, 5, 5, 4, 3, 5, 6, 6, 6, 6, 7, 7, 8, 7, 8, 9, 8, 8, 7, 7, 6, 6, 5, 5, 4,
      4, 4, 4, 5, 4,
    ];
    const focusMins = [
      78, 72, 68, 75, 90, 65, 60, 55, 62, 58, 52, 50, 45, 48, 42, 38, 40, 42,
      50, 52, 58, 60, 65, 70, 75, 80, 85, 78, 72, 76,
    ];
    const fragScores = [
      0.42, 0.48, 0.52, 0.45, 0.38, 0.55, 0.6, 0.65, 0.6, 0.62, 0.68, 0.7, 0.75,
      0.72, 0.8, 0.82, 0.78, 0.76, 0.7, 0.68, 0.62, 0.58, 0.55, 0.5, 0.46, 0.42,
      0.38, 0.4, 0.44, 0.42,
    ];
    const afterHours = [
      22, 28, 35, 20, 15, 40, 45, 52, 48, 50, 60, 65, 72, 70, 85, 90, 82, 78,
      70, 65, 55, 50, 42, 38, 30, 25, 20, 22, 28, 24,
    ];
    const soloPct = [
      0.62, 0.55, 0.52, 0.58, 0.65, 0.5, 0.44, 0.42, 0.45, 0.42, 0.38, 0.35,
      0.32, 0.34, 0.28, 0.26, 0.3, 0.32, 0.38, 0.4, 0.45, 0.48, 0.52, 0.55, 0.6,
      0.63, 0.67, 0.62, 0.58, 0.62,
    ];
    const b2bCounts = [
      1, 2, 2, 1, 0, 2, 3, 3, 3, 3, 4, 4, 5, 4, 5, 6, 5, 5, 4, 4, 3, 3, 2, 2, 1,
      1, 1, 1, 2, 1,
    ];
    const recoveryMins = [
      35, 28, 25, 32, 42, 22, 18, 15, 20, 17, 12, 10, 8, 10, 5, 4, 6, 8, 12, 15,
      18, 20, 25, 28, 32, 38, 42, 36, 30, 34,
    ];
    const attendeeAvg = [
      3.2, 3.8, 4.0, 3.5, 2.8, 4.2, 4.8, 5.0, 4.6, 4.8, 5.2, 5.5, 5.8, 5.5, 6.0,
      6.2, 5.8, 5.6, 5.2, 4.8, 4.4, 4.2, 3.8, 3.5, 3.2, 3.0, 2.8, 3.0, 3.4, 3.2,
    ];

    for (let i = 0; i < 30; i++) {
      const date = new Date(now);
      date.setDate(date.getDate() - (29 - i));

      timeline.push({
        date: date.toISOString().split("T")[0],
        mood_score: moodScores[i],
        mood_energy: energyScores[i],
        meeting_density_pct: meetingDensity[i],
        meeting_count: meetingCounts[i],
        avg_focus_block_mins: focusMins[i],
        max_focus_block_mins: focusMins[i] * 1.8,
        fragmentation_score: fragScores[i],
        after_hours_mins: afterHours[i],
        back_to_back_count: b2bCounts[i],
        avg_recovery_mins: recoveryMins[i],
        solo_time_pct: soloPct[i],
        attendee_avg: attendeeAvg[i],
      });
    }

    return timeline;
  };

  return {
    user: {
      onboarding_done: true,
      calendar_connected: true,
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
    },
    metrics: {
      avg_meeting_density_pct: 0.62,
      meeting_density_trend: "down",
      avg_fragmentation_score: 0.58,
      avg_mood_score: 6.2,
      mood_trend: "up",
      data_days: 30,
    },
    timeline: generateTimeline(),
    top_correlations: [
      {
        feature_name: "meeting_density_pct",
        lag_days: 2,
        correlation: -0.78,
        direction: "negative",
        strength: "strong",
        plain_english: "High meeting load predicts mood dip 2 days later",
      },
      {
        feature_name: "after_hours_mins",
        lag_days: 1,
        correlation: -0.71,
        direction: "negative",
        strength: "strong",
        plain_english: "Working late correlates with lower energy next day",
      },
      {
        feature_name: "avg_focus_block_mins",
        lag_days: 0,
        correlation: 0.65,
        direction: "positive",
        strength: "moderate",
        plain_english:
          "Longer focus blocks associate with better mood same day",
      },
      {
        feature_name: "fragmentation_score",
        lag_days: 2,
        correlation: -0.68,
        direction: "negative",
        strength: "moderate",
        plain_english:
          "Calendar fragmentation leads to mood decline within 2 days",
      },
      {
        feature_name: "back_to_back_count",
        lag_days: 1,
        correlation: -0.62,
        direction: "negative",
        strength: "moderate",
        plain_english: "Multiple consecutive meetings reduce next-day energy",
      },
    ],
    latest_insight: {
      summary:
        "Your behavioral data shows a clear burnout spiral pattern over the past 3 weeks. Meeting density above 65% consistently preceded mood drops by 2–3 days — this lag is the fingerprint Pulse uses to predict your state before you feel it.",
      recommendation:
        "Protect two 90-minute blocks daily as non-negotiable focus time. The data shows your mood recovers fastest when back-to-back meetings drop below 3 per day.",
      patterns: [
        {
          feature: "meeting_density_pct",
          lag_days: 2,
          direction: "negative",
          confidence: 0.78,
          severity: "high",
          plain_english: "High meeting load predicts mood dip 2 days later",
        },
        {
          feature: "after_hours_mins",
          lag_days: 1,
          direction: "negative",
          confidence: 0.71,
          severity: "moderate",
          plain_english: "Working late correlates with lower energy next day",
        },
        {
          feature: "avg_focus_block_mins",
          lag_days: 0,
          direction: "positive",
          confidence: 0.65,
          severity: "moderate",
          plain_english:
            "Longer focus blocks associate with better mood same day",
        },
      ],
      data_quality_note:
        "Based on 30 days of calendar + mood data. Correlations computed at 7-day rolling window.",
      disclaimer:
        "This is behavioral pattern analysis, not medical advice. For mental health support, please speak with a professional.",
      generated_at: new Date().toISOString(),
      model_used: "llama-3.3-70b-versatile",
    },
    sync_status: {
      last_synced_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(),
      events_fetched: 187,
      status: "success",
    },
  };
}
