"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import { RefreshCw } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { api } from "@/lib/api";
import type { DashboardResponse } from "@/lib/types";
import { Sidebar } from "@/components/app/Sidebar";
import { MetricCard } from "@/components/app/MetricCard";
import { TimelineChart } from "@/components/app/TimelineChart";
import { InsightCard } from "@/components/app/InsightCard";
import { EmptyState } from "@/components/app/EmptyState";
import { Button } from "@/components/ui/button";

function greeting() {
  const h = new Date().getHours();
  if (h < 12) return "Good morning.";
  if (h < 17) return "Good afternoon.";
  return "Good evening.";
}

export default function DashboardPage() {
  const [token, setToken] = useState<string | null>(null);
  const [data, setData] = useState<DashboardResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [insightGenerating, setInsightGenerating] = useState(false);
  const [syncing, setSyncing] = useState(false);

  const supabase = createClient();

  const load = useCallback(async (t: string) => {
    try {
      const d = await api.getDashboard(t);
      setData(d);

      // Auto-fix timezone if it's still UTC but browser shows different timezone
      const browserTz = Intl.DateTimeFormat().resolvedOptions().timeZone;
      if (d.user.timezone === "UTC" && browserTz !== "UTC") {
        api.upsertMe(t, { timezone: browserTz }).catch(console.error);
      }
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      if (!session) {
        window.location.href = "/auth";
        return;
      }
      const t = session.access_token ?? null;
      setToken(t);
      if (t) load(t);
    });
  }, [supabase.auth, load]);

  const handleGenerateInsight = async () => {
    if (!token) return;
    setInsightGenerating(true);
    try {
      const insight = await api.generateInsight(token);
      setData((prev) => (prev ? { ...prev, latest_insight: insight } : null));
    } catch (e) {
      console.error(e);
    } finally {
      setInsightGenerating(false);
    }
  };

  const handleSync = async () => {
    if (!token) return;
    setSyncing(true);
    try {
      await api.syncCalendar(token);
      await load(token);
    } catch (e) {
      console.error(e);
    } finally {
      setSyncing(false);
    }
  };

  const moodLoggedToday = !!data?.timeline?.find(
    (d) => d.date === new Date().toISOString().slice(0, 10),
  )?.mood_score;

  return (
    <div className="min-h-screen bg-pulse-bg">
      <Sidebar moodLoggedToday={moodLoggedToday} />
      <main className="pl-[220px] min-h-screen">
        <div className="p-6 max-w-6xl">
          <div className="flex items-center justify-between mb-8">
            <div>
              <h1 className="text-2xl font-light text-pulse-text-primary">
                {greeting()}
              </h1>
              <p className="text-pulse-text-muted text-sm">
                {new Date().toLocaleDateString("en-US", {
                  weekday: "long",
                  month: "long",
                  day: "numeric",
                  year: "numeric",
                })}
              </p>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-xs text-pulse-text-muted">
                Last synced:{" "}
                {data?.sync_status?.last_synced_at
                  ? new Date(data.sync_status.last_synced_at).toLocaleString()
                  : "Never"}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={handleSync}
                disabled={syncing || !token}
                className="border-pulse-border"
              >
                <RefreshCw
                  className={`w-4 h-4 ${syncing ? "animate-spin" : ""}`}
                />
              </Button>
            </div>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              {[1, 2, 3, 4].map((i) => (
                <div
                  key={i}
                  className="h-28 bg-pulse-surface border border-pulse-border rounded-xl animate-pulse"
                />
              ))}
            </div>
          ) : data && data.metrics.data_days >= 7 ? (
            <>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
                <MetricCard
                  label="Meeting Load"
                  value={`${Math.round((data.metrics.avg_meeting_density_pct ?? 0) * 100)}%`}
                  delta={
                    data.metrics.meeting_density_trend === "up"
                      ? "↑ from last week"
                      : data.metrics.meeting_density_trend === "down"
                        ? "↓ from last week"
                        : "Stable"
                  }
                  deltaDirection={
                    data.metrics.meeting_density_trend === "up"
                      ? "up-bad"
                      : "down-good"
                  }
                  severity={
                    (data.metrics.avg_meeting_density_pct ?? 0) > 0.7
                      ? "critical"
                      : (data.metrics.avg_meeting_density_pct ?? 0) > 0.5
                        ? "warning"
                        : "healthy"
                  }
                  tooltip="% of your work window in meetings"
                />
                <MetricCard
                  label="Avg Focus Block"
                  value={`${Math.round(data.timeline?.[0]?.avg_focus_block_mins ?? 0)} min`}
                  severity="healthy"
                />
                <MetricCard
                  label="7-Day Mood Avg"
                  value={(data.metrics.avg_mood_score ?? 0).toFixed(1)}
                  delta={data.metrics.mood_trend}
                  severity={
                    (data.metrics.avg_mood_score ?? 0) < 4
                      ? "critical"
                      : (data.metrics.avg_mood_score ?? 0) < 6
                        ? "warning"
                        : "healthy"
                  }
                />
                <MetricCard
                  label="After-Hours / Day"
                  value={`${((data.timeline?.[0]?.after_hours_mins ?? 0) / 60).toFixed(1)} hrs`}
                  severity="healthy"
                />
              </div>

              <div className="mb-8">
                <h2 className="text-lg font-medium text-pulse-text-primary mb-2">
                  30-Day Behavioral Fingerprint
                </h2>
                <p className="text-pulse-text-muted text-sm mb-4">
                  {new Date().toLocaleDateString("en-US", {
                    month: "long",
                    year: "numeric",
                  })}{" "}
                  — last 30 days
                </p>
                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                  <TimelineChart
                    data={data.timeline ?? []}
                    correlationLag={data.top_correlations?.[0]?.lag_days}
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                <InsightCard
                  insight={data.latest_insight}
                  onGenerate={handleGenerateInsight}
                  isGenerating={insightGenerating}
                />
                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                  <h3 className="font-medium text-pulse-text-primary mb-2">
                    Today&apos;s status
                  </h3>
                  {!moodLoggedToday ? (
                    <div>
                      <p className="text-pulse-text-secondary text-sm mb-2">
                        Log today&apos;s mood to improve your analysis.
                      </p>
                      <Link href="/log">
                        <Button
                          size="sm"
                          className="bg-pulse-accent text-pulse-bg"
                        >
                          Log now →
                        </Button>
                      </Link>
                    </div>
                  ) : (
                    <p className="text-pulse-text-secondary text-sm">
                      Mood logged. Keep it up.
                    </p>
                  )}
                </div>
              </div>

              {data.top_correlations && data.top_correlations.length > 0 && (
                <div>
                  <h2 className="text-lg font-medium text-pulse-text-primary mb-2">
                    What predicts your mood
                  </h2>
                  <p className="text-pulse-text-muted text-sm mb-4">
                    Patterns found across your last 30 days
                  </p>
                  <ul className="space-y-3">
                    {data.top_correlations.map((c, i) => (
                      <li key={i} className="flex items-center gap-4 text-sm">
                        <span className="text-pulse-text-primary w-40">
                          {c.feature_name}
                        </span>
                        <span className="text-pulse-text-muted">
                          Affects mood {c.lag_days} days later
                        </span>
                        <div className="flex-1 h-2 bg-pulse-bg rounded-full max-w-[120px] overflow-hidden">
                          <div
                            className="h-full rounded-full bg-pulse-primary"
                            style={{
                              width: `${Math.min(100, Math.abs(c.correlation) * 100)}%`,
                            }}
                          />
                        </div>
                        <span className="font-mono text-pulse-text-secondary">
                          r={c.correlation.toFixed(2)}
                        </span>
                        <span className="text-pulse-text-muted">
                          {c.direction === "negative" ? "↓ mood" : "↑ mood"}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </>
          ) : (
            <>
              <EmptyState
                title="Building your behavioral baseline"
                progress={{ current: data?.metrics.data_days ?? 0, total: 7 }}
                subtitle="Log mood and connect calendar to unlock pattern analysis"
              />
              <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
                <EmptyState
                  title="Calendar syncing"
                  subtitle={`${data?.sync_status?.events_fetched ?? 0} events fetched`}
                />
                <EmptyState
                  title="Mood logs"
                  subtitle={`${data?.metrics.data_days ?? 0} of 7 logged`}
                />
              </div>
              {data?.timeline && data.timeline.length > 0 && (
                <div className="mt-8">
                  <h2 className="text-lg font-medium text-pulse-text-primary mb-4">
                    30-Day Behavioral Fingerprint
                  </h2>
                  <TimelineChart data={data.timeline} />
                </div>
              )}
              <div className="mt-8">
                <EmptyState
                  title="Pattern analysis unlocks after 7 mood logs"
                  locked
                />
              </div>
            </>
          )}
        </div>
      </main>
    </div>
  );
}
