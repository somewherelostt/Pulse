"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Activity, Calendar, Check } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { api } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Progress } from "@/components/ui/progress";
import { MoodLogger } from "@/components/app/MoodLogger";

const HOURS = Array.from({ length: 17 }, (_, i) => i + 6); // 6–22

export default function OnboardingPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [token, setToken] = useState<string | null>(null);
  const [workStart, setWorkStart] = useState(9);
  const [workEnd, setWorkEnd] = useState(18);
  const [syncing, setSyncing] = useState(false);
  const [pollStatus, setPollStatus] = useState(false);

  const supabase = createClient();

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const s = params.get("step");
    const sync = params.get("syncing") === "true";
    if (s) setStep(parseInt(s, 10) || 1);
    if (sync) {
      setSyncing(true);
      setPollStatus(true);
    }
  }, []);

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
    if (!token || !pollStatus) return;
    const t = setInterval(async () => {
      try {
        const st = await api.getCalendarStatus(token);
        if (st.status === "success") {
          setPollStatus(false);
          setSyncing(false);
          setStep(3);
        }
      } catch {
        // keep polling
      }
    }, 3000);
    return () => clearInterval(t);
  }, [token, pollStatus]);

  const handleStep1Next = async () => {
    if (!token) return;
    try {
      await api.upsertMe(token, {
        work_start_hour: workStart,
        work_end_hour: workEnd,
        timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      });
      setStep(2);
    } catch (e) {
      console.error(e);
    }
  };

  const handleConnectCalendar = async () => {
    if (!token) return;
    try {
      const { url } = await api.getCalendarConnectUrl(token);
      if (url) window.location.href = url;
    } catch (e) {
      console.error(e);
    }
  };

  const handleMoodSubmit = () => {
    router.replace("/dashboard");
  };

  const progressPct = (step / 3) * 100;

  return (
    <main className="min-h-screen bg-pulse-bg flex flex-col">
      <Link href="/" className="absolute top-6 left-6 flex items-center gap-2 text-pulse-text-secondary hover:text-pulse-text-primary">
        <Activity className="w-5 h-5 text-pulse-primary" />
        <span className="font-light text-lg">Pulse</span>
      </Link>

      <div className="flex-1 flex flex-col items-center justify-center p-6 max-w-lg mx-auto w-full">
        <Progress value={progressPct} className="w-full mb-10 h-1.5" />

        {step === 1 && (
          <div className="w-full space-y-6">
            <h2 className="text-2xl font-light text-pulse-text-primary">Set your work hours</h2>
            <p className="text-pulse-text-secondary">
              When does your work day typically start and end? This helps us understand your calendar load accurately.
            </p>
            <div className="flex gap-4">
              <div className="flex-1">
                <label className="text-xs text-pulse-text-muted uppercase tracking-wider block mb-2">Start</label>
                <Select value={String(workStart)} onValueChange={(v) => setWorkStart(parseInt(v, 10))}>
                  <SelectTrigger className="bg-pulse-surface border-pulse-border">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {HOURS.filter((h) => h < workEnd).map((h) => (
                      <SelectItem key={h} value={String(h)}>
                        {h === 12 ? "12 pm" : h < 12 ? `${h} am` : `${h - 12} pm`}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="flex-1">
                <label className="text-xs text-pulse-text-muted uppercase tracking-wider block mb-2">End</label>
                <Select value={String(workEnd)} onValueChange={(v) => setWorkEnd(parseInt(v, 10))}>
                  <SelectTrigger className="bg-pulse-surface border-pulse-border">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {HOURS.filter((h) => h > workStart).map((h) => (
                      <SelectItem key={h} value={String(h)}>
                        {h === 12 ? "12 pm" : h < 12 ? `${h} am` : `${h - 12} pm`}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
            <Button onClick={handleStep1Next} className="w-full">Next</Button>
          </div>
        )}

        {step === 2 && (
          <Card className="w-full bg-pulse-surface border-pulse-border">
            <CardHeader>
              <div className="flex items-center gap-2">
                <Calendar className="w-6 h-6 text-pulse-primary animate-pulse" />
                <CardTitle>Connect your calendar</CardTitle>
              </div>
              <CardDescription>
                We&apos;ll pull 30 days of events to find patterns. We only see event times and attendee counts — never titles, descriptions, or content.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-xs text-pulse-text-muted">
                Event titles are hashed immediately on our server. We store a fingerprint, not your words.
              </p>
              {syncing ? (
                <div className="flex flex-col items-center gap-3 py-4">
                  <div className="w-8 h-8 border-2 border-pulse-primary border-t-transparent rounded-full animate-spin" />
                  <p className="text-pulse-text-secondary">Syncing your last 30 days…</p>
                  <p className="text-sm text-pulse-text-muted">Fetching events… Extracting patterns…</p>
                </div>
              ) : (
                <Button onClick={handleConnectCalendar} className="w-full" disabled={!token}>
                  Connect Google Calendar →
                </Button>
              )}
            </CardContent>
          </Card>
        )}

        {step === 3 && (
          <div className="w-full space-y-6">
            <h2 className="text-2xl font-light text-pulse-text-primary">Log your first mood</h2>
            <p className="text-pulse-text-secondary">
              Before we show insights, tell us how you&apos;re feeling today.
            </p>
            <MoodLogger onSubmit={handleMoodSubmit} existingEntry={null} />
          </div>
        )}
      </div>
    </main>
  );
}
