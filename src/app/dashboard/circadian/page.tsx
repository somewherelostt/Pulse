"use client";

import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { Clock, Zap, Moon, Sun, ArrowRight } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { api } from "@/lib/api";
import { Sidebar } from "@/components/app/Sidebar";
import { Button } from "@/components/ui/button";
import { getDemoCircadianDashboard } from "@/lib/demoData";

export default function CircadianPage() {
  const [token, setToken] = useState<string | null>(null);
  const [dashboard, setDashboard] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [extracting, setExtracting] = useState(false);

  const supabase = createClient();

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      if (!session) {
        window.location.href = "/auth";
        return;
      }
      setToken(session.access_token ?? null);
    });
  }, [supabase.auth]);

  useEffect(() => {
    if (!token) return;
    loadDashboard();
  }, [token]);

  const loadDashboard = async () => {
    if (!token) return;
    setLoading(true);
    try {
      const data = await api.circadianDashboard(token);
      // Use demo data if no insights exist
      if (!data || !data.summary || data.summary.sleep_sessions_count === 0) {
        setDashboard(getDemoCircadianDashboard());
      } else {
        setDashboard(data);
      }
    } catch (error) {
      console.error("Failed to load circadian dashboard:", error);
      // Fallback to demo data on error for showcase
      setDashboard(getDemoCircadianDashboard());
    } finally {
      setLoading(false);
    }
  };

  const handleExtractFeatures = async () => {
    if (!token) return;
    setExtracting(true);
    try {
      await api.circadianExtractFeatures(token);
      await loadDashboard();
    } catch (error) {
      console.error("Failed to extract features:", error);
      alert(
        "Failed to extract features. Make sure you have sleep data logged.",
      );
    } finally {
      setExtracting(false);
    }
  };

  const handleGenerateNarrative = async () => {
    if (!token) return;
    setGenerating(true);
    try {
      const result = await api.circadianGenerateNarrative(token);
      setDashboard({
        ...dashboard,
        latest_insight: {
          narrative: result.narrative,
          interventions: result.interventions,
          model_used: result.model_used,
        },
      });
    } catch (error) {
      console.error("Failed to generate narrative:", error);
      alert("Failed to generate analysis. Need at least 3 days of sleep data.");
    } finally {
      setGenerating(false);
    }
  };

  const summary = dashboard?.summary || {};

  return (
    <div className="min-h-screen bg-pulse-bg">
      <Sidebar moodLoggedToday={false} />
      <main className="pl-[220px] min-h-screen">
        <div className="p-6 max-w-6xl">
          {/* Header */}
          <div className="flex items-start justify-between mb-8">
            <div>
              <div className="flex items-center gap-2 mb-3">
                <Clock className="w-5 h-5 text-pulse-primary" />
                <span className="text-[10px] font-mono uppercase tracking-widest text-pulse-primary">
                  Circadian Analysis
                </span>
              </div>
              <h1 className="text-2xl font-light text-pulse-text-primary mb-1">
                Circadian Rhythm
              </h1>
              <p className="text-pulse-text-muted text-sm max-w-xl">
                Deep analysis of your sleep patterns, consistency, and
                biological rhythm.
              </p>
            </div>
            <div className="flex gap-2">
              <Button
                onClick={handleExtractFeatures}
                disabled={extracting}
                variant="outline"
                className="flex items-center gap-2"
              >
                {extracting ? (
                  <>
                    <div className="w-3.5 h-3.5 border-2 border-pulse-primary/30 border-t-pulse-primary rounded-full animate-spin" />
                    Extracting...
                  </>
                ) : (
                  <>
                    <Zap className="w-4 h-4" />
                    Extract Features
                  </>
                )}
              </Button>
              <Button
                onClick={handleGenerateNarrative}
                disabled={generating}
                className="bg-pulse-primary text-white flex items-center gap-2"
              >
                {generating ? (
                  <>
                    <div className="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    Analyzing...
                  </>
                ) : (
                  <>
                    <Zap className="w-4 h-4" />
                    Generate Analysis
                  </>
                )}
              </Button>
            </div>
          </div>

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {[1, 2, 3].map((i) => (
                <div
                  key={i}
                  className="h-32 bg-pulse-surface border border-pulse-border rounded-xl animate-pulse"
                />
              ))}
            </div>
          ) : (
            <>
              {/* Summary Cards */}
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Moon className="w-4 h-4 text-pulse-primary" />
                    <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted">
                      Avg Sleep
                    </p>
                  </div>
                  <p className="text-2xl font-light text-pulse-text-primary">
                    {summary.avg_sleep_duration_mins
                      ? `${Math.floor(summary.avg_sleep_duration_mins / 60)}h ${summary.avg_sleep_duration_mins % 60}m`
                      : "—"}
                  </p>
                  <p className="text-xs text-pulse-text-muted mt-1">
                    Last 7 days
                  </p>
                </div>

                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Clock className="w-4 h-4 text-pulse-primary" />
                    <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted">
                      Consistency
                    </p>
                  </div>
                  <p className="text-2xl font-light text-pulse-text-primary">
                    {summary.avg_rhythm_consistency_pct
                      ? `${Math.round(summary.avg_rhythm_consistency_pct)}%`
                      : "—"}
                  </p>
                  <p className="text-xs text-pulse-text-muted mt-1">
                    Rhythm score
                  </p>
                </div>

                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Sun className="w-4 h-4 text-pulse-primary" />
                    <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted">
                      Sleep Debt
                    </p>
                  </div>
                  <p className="text-2xl font-light text-pulse-text-primary">
                    {summary.avg_sleep_debt_mins
                      ? `${Math.round(summary.avg_sleep_debt_mins)}m`
                      : "—"}
                  </p>
                  <p className="text-xs text-pulse-text-muted mt-1">
                    Average deficit
                  </p>
                </div>

                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Zap className="w-4 h-4 text-pulse-primary" />
                    <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted">
                      Sleep Score
                    </p>
                  </div>
                  <p className="text-2xl font-light text-pulse-text-primary">
                    {summary.avg_sleep_score
                      ? Math.round(summary.avg_sleep_score)
                      : "—"}
                  </p>
                  <p className="text-xs text-pulse-text-muted mt-1">
                    Quality rating
                  </p>
                </div>
              </div>

              {/* Latest Insight */}
              {dashboard?.latest_insight ? (
                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-6 mb-8">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h2 className="text-sm font-medium text-pulse-text-primary">
                        AI-Generated Analysis
                      </h2>
                      <p className="text-xs text-pulse-text-muted mt-0.5">
                        From{" "}
                        {new Date(
                          dashboard.latest_insight.week_start,
                        ).toLocaleDateString()}{" "}
                        · {dashboard.latest_insight.model_used}
                      </p>
                    </div>
                  </div>
                  <p className="text-sm text-pulse-text-secondary leading-relaxed mb-4">
                    {dashboard.latest_insight.narrative}
                  </p>

                  {dashboard.latest_insight.interventions &&
                    dashboard.latest_insight.interventions.length > 0 && (
                      <div className="space-y-3 pt-4 border-t border-pulse-border">
                        <p className="text-xs font-medium text-pulse-text-primary uppercase tracking-wider">
                          Recommended Interventions
                        </p>
                        {dashboard.latest_insight.interventions.map(
                          (int: any, i: number) => (
                            <div
                              key={i}
                              className="flex items-start gap-3 p-3 rounded-lg bg-pulse-bg/50"
                            >
                              <div
                                className={`w-1.5 h-1.5 rounded-full mt-1.5 shrink-0 ${
                                  int.priority === "high"
                                    ? "bg-pulse-danger"
                                    : int.priority === "medium"
                                      ? "bg-pulse-accent-warm"
                                      : "bg-pulse-primary"
                                }`}
                              />
                              <div className="flex-1 min-w-0">
                                <p className="text-sm font-medium text-pulse-text-primary">
                                  {int.title}
                                </p>
                                <p className="text-xs text-pulse-text-muted mt-0.5">
                                  {int.description}
                                </p>
                              </div>
                              <span
                                className={`text-[9px] font-mono uppercase tracking-wider px-2 py-0.5 rounded-full shrink-0 ${
                                  int.priority === "high"
                                    ? "bg-pulse-danger/10 text-pulse-danger border border-pulse-danger/20"
                                    : int.priority === "medium"
                                      ? "bg-pulse-accent-warm/10 text-pulse-accent-warm border border-pulse-accent-warm/20"
                                      : "bg-pulse-primary/10 text-pulse-primary border border-pulse-primary/20"
                                }`}
                              >
                                {int.priority}
                              </span>
                            </div>
                          ),
                        )}
                      </div>
                    )}
                </div>
              ) : (
                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-12 mb-8 text-center">
                  <Clock className="w-12 h-12 text-pulse-text-muted mx-auto mb-3 opacity-30" />
                  <p className="text-sm text-pulse-text-muted">
                    No circadian analysis yet
                  </p>
                  <p className="text-xs text-pulse-text-muted mt-1 mb-4">
                    Log sleep data and click "Generate Analysis" to get
                    AI-powered insights
                  </p>
                  <Button
                    onClick={handleGenerateNarrative}
                    disabled={generating}
                    className="bg-pulse-primary text-white mx-auto flex items-center gap-2"
                  >
                    {generating ? (
                      <>
                        <div className="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                        Analyzing...
                      </>
                    ) : (
                      <>
                        <Zap className="w-4 h-4" />
                        Generate Analysis
                        <ArrowRight className="w-3 h-3" />
                      </>
                    )}
                  </Button>
                </div>
              )}

              {/* Timeline */}
              {dashboard?.timeline && dashboard.timeline.length > 0 && (
                <div className="bg-pulse-surface border border-pulse-border rounded-xl p-6">
                  <h2 className="text-sm font-medium text-pulse-text-primary mb-4">
                    30-Day Timeline
                  </h2>
                  <div className="space-y-2">
                    {dashboard.timeline
                      .slice(0, 10)
                      .map((day: any, i: number) => (
                        <motion.div
                          key={i}
                          initial={{ opacity: 0, x: -10 }}
                          animate={{ opacity: 1, x: 0 }}
                          transition={{ delay: i * 0.05 }}
                          className="flex items-center gap-4 p-3 rounded-lg hover:bg-pulse-bg/50 transition-colors"
                        >
                          <p className="text-sm text-pulse-text-secondary w-24">
                            {new Date(day.date).toLocaleDateString("en-US", {
                              month: "short",
                              day: "numeric",
                            })}
                          </p>
                          <div className="flex-1 flex items-center gap-2">
                            <div className="flex-1 h-1.5 bg-pulse-border rounded-full overflow-hidden">
                              <div
                                className="h-full bg-pulse-primary rounded-full"
                                style={{
                                  width: `${(day.sleep_duration_mins / 600) * 100}%`,
                                }}
                              />
                            </div>
                            <span className="text-xs text-pulse-text-muted w-16 text-right">
                              {Math.floor(day.sleep_duration_mins / 60)}h{" "}
                              {day.sleep_duration_mins % 60}m
                            </span>
                          </div>
                          {day.rhythm_consistency_pct && (
                            <span className="text-xs text-pulse-text-muted w-12 text-right">
                              {Math.round(day.rhythm_consistency_pct)}%
                            </span>
                          )}
                        </motion.div>
                      ))}
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </main>
    </div>
  );
}
