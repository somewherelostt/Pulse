"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Activity } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { api } from "@/lib/api";
import { MoodLogger } from "@/components/app/MoodLogger";
import type { MoodLog } from "@/lib/types";

export default function LogPage() {
  const router = useRouter();
  const [token, setToken] = useState<string | null>(null);
  const [todayEntry, setTodayEntry] = useState<MoodLog | null>(null);
  const [loading, setLoading] = useState(true);

  const supabase = createClient();

  useEffect(() => {
    supabase.auth.getSession().then(({ data: { session } }) => {
      if (!session) {
        router.replace("/auth");
        return;
      }
      setToken(session.access_token ?? null);
    });
  }, [router, supabase.auth]);

  useEffect(() => {
    if (!token) return;
    api
      .getTodayMood(token)
      .then(setTodayEntry)
      .catch(() => setTodayEntry(null))
      .finally(() => setLoading(false));
  }, [token]);

  const handleSubmit = async (data: Parameters<Parameters<typeof MoodLogger>[0]["onSubmit"]>[0]) => {
    if (!token) return;
    await api.logMood(token, data);
    if (typeof navigator !== "undefined" && navigator.vibrate) {
      navigator.vibrate(50);
    }
    router.push("/dashboard");
  };

  const today = new Date().toISOString().slice(0, 10);

  return (
    <main className="min-h-screen bg-pulse-bg flex flex-col items-center justify-center p-6">
      <Link href="/dashboard" className="absolute top-6 left-6 flex items-center gap-2 text-pulse-text-secondary hover:text-pulse-text-primary">
        <Activity className="w-5 h-5 text-pulse-primary" />
        <span className="font-light text-lg">Pulse</span>
      </Link>

      <div className="w-full max-w-[480px]">
        <h1 className="text-2xl font-light text-pulse-text-primary mb-1">How are you today?</h1>
        <p className="text-pulse-text-muted text-sm mb-8">
          {new Date().toLocaleDateString("en-US", { weekday: "long", month: "long", day: "numeric", year: "numeric" })}
        </p>

        {loading ? (
          <div className="h-64 flex items-center justify-center">
            <div className="w-8 h-8 border-2 border-pulse-primary border-t-transparent rounded-full animate-spin" />
          </div>
        ) : todayEntry ? (
          <div className="space-y-6">
            <p className="text-pulse-text-secondary">You already logged today.</p>
            <MoodLogger onSubmit={handleSubmit} existingEntry={todayEntry} />
          </div>
        ) : (
          <MoodLogger onSubmit={handleSubmit} existingEntry={null} />
        )}
      </div>
    </main>
  );
}
