const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "";

export async function apiFetch<T>(
  path: string,
  options?: RequestInit & { token?: string }
): Promise<T> {
  const { token, ...rest } = options ?? {};
  const res = await fetch(`${API_URL}${path}`, {
    ...rest,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(rest.headers as Record<string, string>),
    },
  });

  if (!res.ok) {
    const error = await res.json().catch(() => ({ message: "Request failed" }));
    const err = error as { message?: string; error?: string };
    throw new Error(err.message ?? err.error ?? `HTTP ${res.status}`);
  }

  return res.json() as Promise<T>;
}

import type { User, MoodLog, DashboardResponse, InsightData } from "./types";

export const api = {
  getMe: (token: string) =>
    apiFetch<User>("/api/v1/users/me", { token }),

  upsertMe: (token: string, data: Partial<Pick<User, "work_start_hour" | "work_end_hour" | "timezone">>) =>
    apiFetch<User>("/api/v1/users/me", {
      method: "POST",
      body: JSON.stringify(data),
      token,
    }),

  getCalendarConnectUrl: (token: string) =>
    apiFetch<{ url: string }>("/api/v1/calendar/connect", { token }),

  getCalendarStatus: (token: string) =>
    apiFetch<{
      connected: boolean;
      last_synced_at: string | null;
      status: string;
    }>("/api/v1/calendar/status", { token }),

  syncCalendar: (token: string) =>
    apiFetch<{ status: string }>("/api/v1/calendar/sync", {
      method: "POST",
      token,
    }),

  logMood: (
    token: string,
    data: Partial<Pick<MoodLog, "score" | "energy" | "anxiety" | "note" | "tags">>
  ) =>
    apiFetch<MoodLog>("/api/v1/mood", {
      method: "POST",
      body: JSON.stringify(data),
      token,
    }),

  getTodayMood: (token: string) =>
    apiFetch<MoodLog | null>("/api/v1/mood/today", { token }),

  getDashboard: (token: string) =>
    apiFetch<DashboardResponse>("/api/v1/dashboard", { token }),

  getLatestInsight: (token: string) =>
    apiFetch<InsightData | null>("/api/v1/insights/latest", { token }),

  generateInsight: (token: string) =>
    apiFetch<InsightData>("/api/v1/insights/generate", {
      method: "POST",
      token,
    }),
};
