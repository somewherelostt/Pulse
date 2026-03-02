"use client";

import { useState, useEffect, useCallback } from "react";
import Link from "next/link";
import {
  RefreshCw,
  ArrowRight,
  Network,
  TrendingDown,
  TrendingUp,
} from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { api } from "@/lib/api";
import type { DashboardResponse } from "@/lib/types";
import { getDemoDashboardData } from "@/lib/demoData";
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

function BurnoutRiskBadge({ score }: { score: number }) {
  const level =
    score >= 0.7
      ? {
          label: "High Risk",
          color: "text-pulse-danger bg-pulse-danger/10 border-pulse-danger/30",
        }
      : score >= 0.45
        ? {
            label: "Elevated",
            color:
              "text-pulse-accent-warm bg-pulse-accent-warm/10 border-pulse-accent-warm/30",
          }
        : {
            label: "Normal Range",
            color: "text-emerald-400 bg-emerald-400/10 border-emerald-400/30",
          };

  return (
    <span
      className={`text-[10px] font-mono uppercase tracking-wider px-2.5 py-1 rounded-full border ${level.color}`}
    >
      {level.label}
    </span>
  );
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

      // If no meaningful data exists, use rich demo data for showcase
      if ((d.metrics.data_days ?? 0) < 7) {
        console.log("📊 Using demo data for showcase (insufficient real data)");
        setData(getDemoDashboardData());
      } else {
        setData(d);
      }

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

  // Derive a simple burnout risk score from meeting density + after hours
  const burnoutRisk = data
    ? Math.min(
        1,
        (data.metrics.avg_meeting_density_pct ?? 0) * 0.7 +
          (data.timeline?.[0]?.after_hours_mins ?? 0) / 500,
      )
    : 0;

  return (
    <div className="min-h-screen bg-pulse-bg">
      <Sidebar moodLoggedToday={moodLoggedToday} />
      <main className="pl-[220px] min-h-screen">
        <div className="p-6 max-w-6xl">
          {/* Header */}
          <div className="flex items-start justify-between mb-8">
            <div>
              <h1 className="text-2xl font-light text-pulse-text-primary">
                {greeting()}
              </h1>
              <p className="text-pulse-text-muted text-sm mt-0.5">
                {new Date().toLocaleDateString("en-US", {
                  weekday: "long",
                  month: "long",
                  day: "numeric",
                })}
              </p>
            </div>
            <div className="flex items-center gap-3">
              {data && data.metrics.data_days >= 7 && (
                <BurnoutRiskBadge score={burnoutRisk} />
              )}
              <span className="text-xs text-pulse-text-muted">
                {data?.sync_status?.last_synced_at
                  ? `Synced ${new Date(data.sync_status.last_synced_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}`
                  : "Never synced"}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={handleSync}
                disabled={syncing || !token}
                className="border-pulse-border h-8 w-8 p-0"
              >
                <RefreshCw
                  className={`w-3.5 h-3.5 ${syncing ? "animate-spin" : ""}`}
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
              {/* Metric cards */}
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
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
                  tooltip="% of your work window spent in meetings (7-day avg)"
                />
                <MetricCard
                  label="Avg Focus Block"
                  value={`${Math.round(data.timeline?.find((d) => d.avg_focus_block_mins != null)?.avg_focus_block_mins ?? 0)} min`}
                  severity={
                    (data.timeline?.find((d) => d.avg_focus_block_mins != null)
                      ?.avg_focus_block_mins ?? 0) < 45
                      ? "warning"
                      : "healthy"
                  }
                  tooltip="Average uninterrupted focus window"
                />
                <MetricCard
                  label="7-Day Mood Avg"
                  value={(data.metrics.avg_mood_score ?? 0).toFixed(1)}
                  delta={
                    data.metrics.mood_trend === "up"
                      ? "↑ improving"
                      : data.metrics.mood_trend === "down"
                        ? "↓ declining"
                        : "Stable"
                  }
                  deltaDirection={
                    data.metrics.mood_trend === "up"
                      ? "up-good"
                      : data.metrics.mood_trend === "down"
                        ? "down-bad"
                        : "down-good"
                  }
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
                  value={`${((data.timeline?.[0]?.after_hours_mins ?? 0) / 60).toFixed(1)}h`}
                  severity={
                    (data.timeline?.[0]?.after_hours_mins ?? 0) > 60
                      ? "warning"
                      : "healthy"
                  }
                  tooltip="Average minutes worked outside your set hours"
                />
              </div>

              {/* Timeline */}
              <div id="calendar" className="mb-8">
                <div className="flex items-center justify-between mb-3">
                  <div>
                    <h2 className="text-base font-medium text-pulse-text-primary">
                      30-Day Behavioral Fingerprint
                    </h2>
                    <p className="text-pulse-text-muted text-xs mt-0.5">
                      Calendar density, after-hours drift, and mood overlaid
                    </p>
                  </div>
                  <span className="text-[10px] font-mono text-pulse-text-muted uppercase tracking-widest">
                    {data.sync_status?.events_fetched ?? 0} events synced
                  </span>
                </div>
                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                  <TimelineChart
                    data={data.timeline ?? []}
                    correlationLag={data.top_correlations?.[0]?.lag_days}
                  />
                </div>
              </div>

              {/* Insight + Today */}
              <div
                id="insights"
                className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8"
              >
                <InsightCard
                  insight={data.latest_insight}
                  onGenerate={handleGenerateInsight}
                  isGenerating={insightGenerating}
                />
                <div className="space-y-3">
                  {/* Today status */}
                  <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                    <h3 className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted mb-3">
                      Today
                    </h3>
                    {!moodLoggedToday ? (
                      <div className="flex items-center justify-between">
                        <p className="text-sm text-pulse-text-secondary">
                          Mood not logged yet
                        </p>
                        <Link href="/log">
                          <Button
                            size="sm"
                            className="bg-pulse-accent text-pulse-bg h-7 text-xs"
                          >
                            Log now →
                          </Button>
                        </Link>
                      </div>
                    ) : (
                      <p className="text-sm text-pulse-text-secondary flex items-center gap-2">
                        <span className="w-1.5 h-1.5 rounded-full bg-pulse-accent" />
                        Mood logged. Keep it up.
                      </p>
                    )}
                  </div>

                  {/* Constellation nudge */}
                  <Link href="/dashboard/constellation" className="block group">
                    <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4 hover:border-pulse-primary/40 transition-all hover:-translate-y-0.5 hover:shadow-lg">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-lg bg-pulse-primary/10 flex items-center justify-center shrink-0">
                          <Network className="w-4 h-4 text-pulse-primary" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-0.5">
                            <p className="text-sm font-medium text-pulse-text-primary">
                              Constellation
                            </p>
                            <span className="text-[9px] font-mono uppercase tracking-wider px-1.5 py-0.5 rounded-full bg-pulse-primary/15 text-pulse-primary border border-pulse-primary/20">
                              New
                            </span>
                          </div>
                          <p className="text-xs text-pulse-text-muted">
                            Find peers who've been through your pattern
                          </p>
                        </div>
                        <ArrowRight className="w-4 h-4 text-pulse-text-muted group-hover:text-pulse-primary group-hover:translate-x-1 transition-all" />
                      </div>
                    </div>
                  </Link>
                </div>
              </div>

              {/* Correlations */}
              {data.top_correlations && data.top_correlations.length > 0 && (
                <div>
                  <div className="mb-3">
                    <h2 className="text-base font-medium text-pulse-text-primary">
                      What predicts your mood
                    </h2>
                    <p className="text-pulse-text-muted text-xs mt-0.5">
                      Lagged correlations across 30 days
                    </p>
                  </div>
                  <div className="bg-pulse-surface border border-pulse-border rounded-xl overflow-hidden">
                    {data.top_correlations.map((c, i) => (
                      <div
                        key={i}
                        className={`flex items-center gap-4 px-4 py-3 text-sm ${
                          i < data.top_correlations.length - 1
                            ? "border-b border-pulse-border/50"
                            : ""
                        }`}
                      >
                        <div
                          className="w-1.5 h-1.5 rounded-full shrink-0"
                          style={{
                            backgroundColor:
                              c.direction === "negative"
                                ? "var(--color-danger)"
                                : "var(--color-accent)",
                          }}
                        />
                        <span className="text-pulse-text-primary text-xs font-medium w-36 truncate">
                          {c.feature_name.replace(/_/g, " ")}
                        </span>
                        <span className="text-pulse-text-muted text-xs">
                          → mood{" "}
                          {c.lag_days === 0
                            ? "same day"
                            : `${c.lag_days}d later`}
                        </span>
                        <div className="flex-1 h-1.5 bg-pulse-bg rounded-full overflow-hidden max-w-[100px]">
                          <div
                            className="h-full rounded-full"
                            style={{
                              width: `${Math.min(100, Math.abs(c.correlation) * 100)}%`,
                              backgroundColor:
                                c.direction === "negative"
                                  ? "var(--color-danger)"
                                  : "var(--color-accent)",
                            }}
                          />
                        </div>
                        <span className="font-mono text-[10px] text-pulse-text-secondary w-12 text-right">
                          r={c.correlation.toFixed(2)}
                        </span>
                        {c.direction === "negative" ? (
                          <TrendingDown className="w-3.5 h-3.5 text-pulse-danger shrink-0" />
                        ) : (
                          <TrendingUp className="w-3.5 h-3.5 text-pulse-accent shrink-0" />
                        )}
                      </div>
                    ))}
                  </div>
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
                  <h2 className="text-base font-medium text-pulse-text-primary mb-4">
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
