"use client";

import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { Moon, TrendingUp, TrendingDown, Plus } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { api } from "@/lib/api";
import { Sidebar } from "@/components/app/Sidebar";
import { Button } from "@/components/ui/button";
import { getDemoSleepData } from "@/lib/demoData";

export default function SleepPage() {
  const [token, setToken] = useState<string | null>(null);
  const [sleepData, setSleepData] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [showAddModal, setShowAddModal] = useState(false);
  const [formData, setFormData] = useState({
    date: new Date().toISOString().split("T")[0],
    bedtime_hour: 23,
    bedtime_min: 0,
    wake_hour: 7,
    wake_min: 0,
    sleep_score: 75,
  });

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
    loadSleepData();
  }, [token]);

  const loadSleepData = async () => {
    if (!token) return;
    setLoading(true);
    try {
      const data = await api.sleepGetRange(token);
      // Use demo data if no real data exists
      if (!data || data.length === 0) {
        setSleepData(getDemoSleepData());
      } else {
        setSleepData(data);
      }
    } catch (error) {
      console.error("Failed to load sleep data:", error);
      // Fallback to demo data on error for showcase
      setSleepData(getDemoSleepData());
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;

    try {
      await api.sleepLogManual(token, formData);
      setShowAddModal(false);
      loadSleepData();
      setFormData({
        date: new Date().toISOString().split("T")[0],
        bedtime_hour: 23,
        bedtime_min: 0,
        wake_hour: 7,
        wake_min: 0,
        sleep_score: 75,
      });
    } catch (error) {
      console.error("Failed to log sleep:", error);
      alert("Failed to log sleep. Please try again.");
    }
  };

  const avgSleepDuration =
    sleepData.length > 0
      ? Math.round(
          sleepData.reduce((sum, d) => sum + (d.total_mins || 0), 0) /
            sleepData.length,
        )
      : 0;

  const avgSleepScore =
    sleepData.length > 0
      ? Math.round(
          sleepData.reduce((sum, d) => sum + (d.sleep_score || 0), 0) /
            sleepData.length,
        )
      : 0;

  const formatTime = (hour: number, min: number) => {
    const h = hour % 12 || 12;
    const ampm = hour < 12 ? "AM" : "PM";
    return `${h}:${min.toString().padStart(2, "0")} ${ampm}`;
  };

  return (
    <div className="min-h-screen bg-pulse-bg">
      <Sidebar moodLoggedToday={false} />
      <main className="pl-[220px] min-h-screen">
        <div className="p-6 max-w-6xl">
          {/* Header */}
          <div className="flex items-start justify-between mb-8">
            <div>
              <div className="flex items-center gap-2 mb-3">
                <Moon className="w-5 h-5 text-pulse-primary" />
                <span className="text-[10px] font-mono uppercase tracking-widest text-pulse-primary">
                  Sleep Tracking
                </span>
              </div>
              <h1 className="text-2xl font-light text-pulse-text-primary mb-1">
                Sleep Log
              </h1>
              <p className="text-pulse-text-muted text-sm max-w-xl">
                Track your sleep patterns to understand circadian rhythm and
                recovery quality.
              </p>
            </div>
            <Button
              onClick={() => setShowAddModal(true)}
              className="bg-pulse-primary text-white flex items-center gap-2"
            >
              <Plus className="w-4 h-4" />
              Log Sleep
            </Button>
          </div>

          {/* Summary Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
            <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
              <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted mb-2">
                Avg Sleep Duration
              </p>
              <p className="text-2xl font-light text-pulse-text-primary">
                {Math.floor(avgSleepDuration / 60)}h {avgSleepDuration % 60}m
              </p>
              <p className="text-xs text-pulse-text-muted mt-1">Last 30 days</p>
            </div>
            <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
              <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted mb-2">
                Avg Sleep Score
              </p>
              <p className="text-2xl font-light text-pulse-text-primary">
                {avgSleepScore}
                <span className="text-sm text-pulse-text-muted">/100</span>
              </p>
              <p className="text-xs text-pulse-text-muted mt-1">
                Quality rating
              </p>
            </div>
            <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
              <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted mb-2">
                Total Logs
              </p>
              <p className="text-2xl font-light text-pulse-text-primary">
                {sleepData.length}
              </p>
              <p className="text-xs text-pulse-text-muted mt-1">
                Sessions tracked
              </p>
            </div>
          </div>

          {/* Sleep Log Table */}
          <div className="bg-pulse-surface border border-pulse-border rounded-xl overflow-hidden">
            <div className="p-4 border-b border-pulse-border">
              <h2 className="text-sm font-medium text-pulse-text-primary">
                Sleep Sessions
              </h2>
              <p className="text-xs text-pulse-text-muted mt-0.5">
                Your recent sleep history
              </p>
            </div>
            {loading ? (
              <div className="p-12 flex items-center justify-center">
                <div className="w-8 h-8 border-2 border-pulse-primary border-t-transparent rounded-full animate-spin" />
              </div>
            ) : sleepData.length === 0 ? (
              <div className="p-12 text-center">
                <Moon className="w-12 h-12 text-pulse-text-muted mx-auto mb-3 opacity-30" />
                <p className="text-sm text-pulse-text-muted">
                  No sleep data yet
                </p>
                <p className="text-xs text-pulse-text-muted mt-1">
                  Click "Log Sleep" to add your first entry
                </p>
              </div>
            ) : (
              <div className="divide-y divide-pulse-border/50">
                {sleepData.map((session, i) => (
                  <motion.div
                    key={i}
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: i * 0.05 }}
                    className="p-4 flex items-center gap-4 hover:bg-pulse-bg/50 transition-colors"
                  >
                    <div className="flex-1">
                      <p className="text-sm font-medium text-pulse-text-primary">
                        {new Date(session.date).toLocaleDateString("en-US", {
                          weekday: "short",
                          month: "short",
                          day: "numeric",
                        })}
                      </p>
                      {session.bedtime && session.wake_time && (
                        <p className="text-xs text-pulse-text-muted mt-0.5">
                          {session.bedtime} → {session.wake_time}
                        </p>
                      )}
                    </div>
                    <div className="text-right">
                      <p className="text-sm font-medium text-pulse-text-primary">
                        {Math.floor(session.total_mins / 60)}h{" "}
                        {session.total_mins % 60}m
                      </p>
                      <p className="text-xs text-pulse-text-muted mt-0.5">
                        {session.provider || "Manual"}
                      </p>
                    </div>
                    {session.sleep_score && (
                      <div className="w-16 text-center">
                        <div
                          className={`inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs ${
                            session.sleep_score >= 80
                              ? "bg-emerald-400/10 text-emerald-400"
                              : session.sleep_score >= 60
                                ? "bg-pulse-accent/10 text-pulse-accent"
                                : "bg-pulse-danger/10 text-pulse-danger"
                          }`}
                        >
                          {session.sleep_score}
                        </div>
                      </div>
                    )}
                  </motion.div>
                ))}
              </div>
            )}
          </div>
        </div>
      </main>

      {/* Add Sleep Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-6">
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-pulse-surface border border-pulse-border rounded-xl p-6 max-w-md w-full"
          >
            <h2 className="text-lg font-medium text-pulse-text-primary mb-4">
              Log Sleep Session
            </h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="text-xs text-pulse-text-muted block mb-1">
                  Date
                </label>
                <input
                  type="date"
                  value={formData.date}
                  onChange={(e) =>
                    setFormData({ ...formData, date: e.target.value })
                  }
                  className="w-full px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-sm text-pulse-text-primary focus:outline-none focus:border-pulse-primary"
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-xs text-pulse-text-muted block mb-1">
                    Bedtime (Hour)
                  </label>
                  <input
                    type="number"
                    min="0"
                    max="23"
                    value={formData.bedtime_hour}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        bedtime_hour: parseInt(e.target.value),
                      })
                    }
                    className="w-full px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-sm text-pulse-text-primary focus:outline-none focus:border-pulse-primary"
                  />
                </div>
                <div>
                  <label className="text-xs text-pulse-text-muted block mb-1">
                    Bedtime (Min)
                  </label>
                  <input
                    type="number"
                    min="0"
                    max="59"
                    step="15"
                    value={formData.bedtime_min}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        bedtime_min: parseInt(e.target.value),
                      })
                    }
                    className="w-full px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-sm text-pulse-text-primary focus:outline-none focus:border-pulse-primary"
                  />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-xs text-pulse-text-muted block mb-1">
                    Wake Time (Hour)
                  </label>
                  <input
                    type="number"
                    min="0"
                    max="23"
                    value={formData.wake_hour}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        wake_hour: parseInt(e.target.value),
                      })
                    }
                    className="w-full px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-sm text-pulse-text-primary focus:outline-none focus:border-pulse-primary"
                  />
                </div>
                <div>
                  <label className="text-xs text-pulse-text-muted block mb-1">
                    Wake Time (Min)
                  </label>
                  <input
                    type="number"
                    min="0"
                    max="59"
                    step="15"
                    value={formData.wake_min}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        wake_min: parseInt(e.target.value),
                      })
                    }
                    className="w-full px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-sm text-pulse-text-primary focus:outline-none focus:border-pulse-primary"
                  />
                </div>
              </div>
              <div>
                <label className="text-xs text-pulse-text-muted block mb-1">
                  Sleep Quality Score (0-100)
                </label>
                <input
                  type="number"
                  min="0"
                  max="100"
                  value={formData.sleep_score}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      sleep_score: parseInt(e.target.value),
                    })
                  }
                  className="w-full px-3 py-2 rounded-lg bg-pulse-bg border border-pulse-border text-sm text-pulse-text-primary focus:outline-none focus:border-pulse-primary"
                />
              </div>
              <div className="flex gap-3 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setShowAddModal(false)}
                  className="flex-1"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  className="flex-1 bg-pulse-primary text-white"
                >
                  Save Sleep Log
                </Button>
              </div>
            </form>
          </motion.div>
        </div>
      )}
    </div>
  );
}
