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

/**
 * Demo sleep data for showcase (30 days)
 */
export function getDemoSleepData() {
  const now = new Date();
  const sleepData = [];

  // Pattern: good sleep initially, deteriorates mid-month, recovers
  const sleepDurations = [
    450, 440, 455, 435, 470, 420, 410, 400, 415, 425, 405, 390, 380, 370, 360,
    355, 365, 375, 390, 400, 410, 425, 435, 450, 460, 465, 455, 445, 450, 460,
  ]; // minutes
  const sleepScores = [
    85, 82, 87, 80, 90, 78, 75, 72, 76, 78, 74, 68, 65, 62, 58, 56, 60, 65, 70,
    72, 75, 78, 82, 85, 88, 90, 87, 84, 86, 89,
  ];
  const bedtimes = [
    23, 23.5, 23, 0, 23, 0.5, 1, 1.5, 0.5, 0, 1, 1.5, 2, 2.5, 2, 1.5, 1, 0.5, 0,
    23.5, 23, 23, 22.5, 22.5, 23, 23, 23.5, 23, 22.5, 23,
  ]; // hours (decimal)

  for (let i = 0; i < 30; i++) {
    const date = new Date(now);
    date.setDate(date.getDate() - (29 - i));

    const bedtime = bedtimes[i];
    const bedtimeHour = Math.floor(bedtime);
    const bedtimeMin = Math.round((bedtime % 1) * 60);

    const sleepMins = sleepDurations[i];
    const wakeTime = (bedtime + sleepMins / 60) % 24;
    const wakeHour = Math.floor(wakeTime);
    const wakeMin = Math.round((wakeTime % 1) * 60);

    const fmtTime = (h: number, m: number) => {
      const hour = h % 12 || 12;
      const ampm = h < 12 ? "AM" : "PM";
      return `${hour}:${m.toString().padStart(2, "0")} ${ampm}`;
    };
    sleepData.push({
      session_id: `demo-sleep-${i}`,
      date: date.toISOString().split("T")[0],
      bedtime_hour: bedtimeHour,
      bedtime_min: bedtimeMin,
      wake_hour: wakeHour,
      wake_min: wakeMin,
      bedtime: fmtTime(bedtimeHour, bedtimeMin),
      wake_time: fmtTime(wakeHour, wakeMin),
      total_mins: sleepMins,
      sleep_score: sleepScores[i],
      source: "manual",
      created_at: date.toISOString(),
    });
  }

  return sleepData;
}

/**
 * Demo circadian dashboard data for showcase
 */
export function getDemoCircadianDashboard() {
  return {
    summary: {
      avg_sleep_duration_mins: 415,
      avg_rhythm_consistency_pct: 72,
      avg_sleep_debt_mins: 35,
      avg_recovery_score: 78,
      avg_sleep_score: 77,
      sleep_sessions_count: 30,
      features_extracted_count: 28,
    },
    latest_insight: {
      narrative:
        "Your circadian analysis reveals a bi-phasic pattern with notable mid-month disruption. Sleep onset showed progressive delay from 11 PM to 2 AM during days 10-16, correlating with the meeting density spike observed in your behavioral data. Recovery phase (days 20-30) demonstrates homeostatic regulation — your body naturally stabilized bedtime back to 11 PM when workload normalized. The consistency score of 72% indicates moderate rhythm entrainment; ideal range is 80-90%. Sleep debt accumulated to 8.5 hours during peak stress, cleared within 6 days post-recovery.",
      interventions: [
        {
          title: "Screen cutoff on high-meeting days",
          description:
            "Implement strict 10:30 PM screen cutoff on high-meeting days (>6 meetings). Data shows 90-minute delay in sleep onset when meetings exceed this threshold.",
          priority: "high",
        },
        {
          title: "Circadian anchor wake time",
          description:
            "Schedule 'circadian anchor' — wake at same time (±15 min) even on weekends. This single intervention can boost consistency score by 15-20%.",
          priority: "medium",
        },
        {
          title: "Evening meeting-free zone",
          description:
            "Block 7-8 PM as meeting-free zone. Evening meetings correlate with +45 min bedtime delay and -12 point sleep score reduction.",
          priority: "medium",
        },
      ],
      model_used: "llama-3.3-70b-versatile",
    },
    timeline: generateDemoCircadianTimeline(),
  };
}

function generateDemoCircadianTimeline() {
  const now = new Date();
  const timeline = [];

  const sleepDurations = [
    450, 440, 455, 435, 470, 420, 410, 400, 415, 425, 405, 390, 380, 370, 360,
    355, 365, 375, 390, 400, 410, 425, 435, 450, 460, 465, 455, 445, 450, 460,
  ];
  const consistencyScores = [
    85, 82, 80, 78, 75, 72, 70, 68, 65, 63, 60, 58, 55, 53, 50, 52, 55, 60, 65,
    68, 72, 75, 78, 80, 82, 85, 87, 85, 83, 85,
  ];
  const sleepDebt = [
    10, 15, 5, 20, -10, 25, 35, 40, 30, 25, 40, 50, 60, 70, 80, 85, 75, 65, 50,
    45, 35, 25, 15, 5, -5, -10, 0, 10, 5, 0,
  ];
  const recoveryScores = [
    85, 82, 87, 80, 90, 78, 75, 72, 76, 78, 74, 68, 65, 62, 58, 56, 60, 65, 70,
    72, 75, 78, 82, 85, 88, 90, 87, 84, 86, 89,
  ];

  for (let i = 0; i < 30; i++) {
    const date = new Date(now);
    date.setDate(date.getDate() - (29 - i));

    timeline.push({
      date: date.toISOString().split("T")[0],
      sleep_duration_mins: sleepDurations[i],
      sleep_score: recoveryScores[i],
      rhythm_consistency_pct: consistencyScores[i],
      sleep_debt_mins: sleepDebt[i],
      recovery_score: recoveryScores[i],
    });
  }

  return timeline;
}
