"use client";

export const dynamic = "force-dynamic";

import { useState } from "react";
import Link from "next/link";
import { Activity, Ghost, Mail } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

export default function AuthPage() {
  const [loading, setLoading] = useState<string | null>(null);
  const [emailSent, setEmailSent] = useState(false);
  const [email, setEmail] = useState("");

  const supabase = createClient();

  const handleAnonymous = async () => {
    setLoading("anonymous");
    try {
      const { error } = await supabase.auth.signInAnonymously();
      if (error) throw error;
      window.location.href = "/onboarding";
    } catch (e) {
      console.error(e);
      setLoading(null);
    }
  };

  const handleEmail = async () => {
    if (!email.trim()) return;
    setLoading("email");
    try {
      const { error } = await supabase.auth.signInWithOtp({
        email: email.trim(),
        options: { emailRedirectTo: `${window.location.origin}/onboarding` },
      });
      if (error) throw error;
      setEmailSent(true);
    } catch (e) {
      console.error(e);
    }
    setLoading(null);
  };

  return (
    <main className="min-h-screen bg-pulse-bg flex flex-col items-center justify-center p-6">
      <Link href="/" className="absolute top-6 left-6 flex items-center gap-2 text-pulse-text-secondary hover:text-pulse-text-primary">
        <Activity className="w-5 h-5 text-pulse-primary" />
        <span className="font-light text-lg">Pulse</span>
      </Link>

      <div className="w-full max-w-md space-y-10">
        <div className="text-center space-y-2">
          <h1 className="text-4xl md:text-5xl font-light tracking-tight text-pulse-text-primary" style={{ fontFamily: "var(--font-display)" }}>
            Start anonymously.
          </h1>
          <p className="text-pulse-text-secondary text-lg">
            No email. No account. Just connect your calendar.
          </p>
        </div>

        <div className="grid gap-4">
          <Card className="bg-pulse-surface-raised border-pulse-primary/30 border shadow-lg" style={{ boxShadow: "var(--shadow-primary)" }}>
            <CardHeader>
              <div className="flex items-center gap-2">
                <Ghost className="w-5 h-5 text-pulse-primary" />
                <CardTitle>Continue anonymously</CardTitle>
              </div>
              <CardDescription>
                No email required. Your data stays private. You can add an email later.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button
                onClick={handleAnonymous}
                disabled={loading !== null}
                className="w-full bg-primary hover:bg-primary/90"
              >
                {loading === "anonymous" ? "Starting…" : "Start — no account needed"}
              </Button>
            </CardContent>
          </Card>

          <Card className="bg-pulse-surface border-pulse-border">
            <CardHeader>
              <div className="flex items-center gap-2">
                <Mail className="w-5 h-5 text-pulse-text-secondary" />
                <CardTitle>Sign in with email</CardTitle>
              </div>
              <CardDescription>
                Get your data on any device.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              {emailSent ? (
                <p className="text-sm text-pulse-accent">Check your email for the sign-in link.</p>
              ) : (
                <>
                  <Input
                    type="email"
                    placeholder="you@example.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="bg-pulse-bg border-pulse-border"
                  />
                  <Button
                    variant="outline"
                    onClick={handleEmail}
                    disabled={loading !== null}
                    className="w-full border-pulse-border text-pulse-text-primary hover:bg-pulse-surface-raised"
                  >
                    {loading === "email" ? "Sending…" : "Continue with email"}
                  </Button>
                </>
              )}
            </CardContent>
          </Card>
        </div>

        <p className="text-center text-xs text-pulse-text-muted">
          Both options create identical anonymous data environments. We never sell data.
        </p>
      </div>
    </main>
  );
}
