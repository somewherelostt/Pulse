"use client";

import {
  ComposedChart,
  Area,
  Bar,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  ReferenceArea,
} from "recharts";
import type { TimelinePoint } from "@/lib/types";

type TimelineChartProps = {
  data: TimelinePoint[];
  correlationLag?: number | null;
};

export function TimelineChart({ data, correlationLag }: TimelineChartProps) {
  const chartData = data.map((d) => ({
    date: d.date,
    meetingLoad: d.meeting_density_pct != null ? Math.round(d.meeting_density_pct * 100) : null,
    afterHours: d.after_hours_mins ?? 0,
    mood: d.mood_score,
  }));

  return (
    <div className="w-full h-[280px]">
      <ResponsiveContainer width="100%" height="100%">
        <ComposedChart data={chartData} margin={{ top: 8, right: 8, left: 0, bottom: 0 }}>
          <defs>
            <linearGradient id="meetingFill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="rgba(110,123,242,0.25)" />
              <stop offset="100%" stopColor="rgba(110,123,242,0.02)" />
            </linearGradient>
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="var(--border-subtle)" vertical={false} />
          <XAxis
            dataKey="date"
            tick={{ fontSize: 10, fill: "var(--text-muted)" }}
            tickFormatter={(v) => (v ? new Date(v).toLocaleDateString("en-US", { month: "short", day: "numeric" }) : "")}
          />
          <YAxis
            yAxisId="left"
            domain={[0, 100]}
            tick={{ fontSize: 10, fill: "var(--text-muted)" }}
            tickFormatter={(v) => `${v}%`}
          />
          <YAxis yAxisId="right" orientation="right" hide />
          <Tooltip
            contentStyle={{ backgroundColor: "var(--bg-raised)", border: "1px solid var(--border)" }}
            labelFormatter={(v) => new Date(v).toLocaleDateString()}
            formatter={(value: number, name: string) => [
              name === "meetingLoad" ? `${value}%` : name === "afterHours" ? `${value} min` : value,
              name === "meetingLoad" ? "Meeting Load" : name === "afterHours" ? "After Hours" : "Mood",
            ]}
          />
          <Area
            yAxisId="left"
            type="monotone"
            dataKey="meetingLoad"
            name="Meeting Load"
            fill="url(#meetingFill)"
            stroke="var(--primary)"
            strokeWidth={1.5}
          />
          <Bar yAxisId="right" dataKey="afterHours" fill="rgba(247,183,49,0.4)" name="After Hours" barSize={4} radius={2} />
          <Line
            yAxisId="right"
            type="monotone"
            dataKey="mood"
            name="Mood Score"
            stroke="var(--accent)"
            strokeWidth={2}
            dot={{ fill: "var(--accent)", r: 3 }}
            connectNulls={false}
          />
          <Legend
            wrapperStyle={{ fontSize: 11 }}
            formatter={(value) => (
              <span className="text-pulse-text-secondary text-xs">{value}</span>
            )}
          />
        </ComposedChart>
      </ResponsiveContainer>
    </div>
  );
}
