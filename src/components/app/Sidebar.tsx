"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Activity,
  LayoutDashboard,
  Calendar,
  Heart,
  Lightbulb,
  Lock,
  Network,
  Moon,
  Clock,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";

const NAV = [
  { href: "/dashboard", label: "Overview", icon: LayoutDashboard },
  { href: "/dashboard#calendar", label: "Calendar", icon: Calendar },
  { href: "/log", label: "Mood Log", icon: Heart },
  { href: "/dashboard#insights", label: "Insights", icon: Lightbulb },
  { href: "/dashboard/sleep", label: "Sleep", icon: Moon },
  { href: "/dashboard/circadian", label: "Circadian", icon: Clock },
  {
    href: "/dashboard/constellation",
    label: "Constellation",
    icon: Network,
    badge: "New",
  },
];

export function Sidebar({ moodLoggedToday }: { moodLoggedToday: boolean }) {
  const pathname = usePathname();

  return (
    <aside className="fixed left-0 top-0 bottom-0 w-[220px] bg-pulse-surface border-r border-pulse-border flex flex-col z-40">
      <div className="p-4 border-b border-pulse-border flex items-center gap-2">
        <Activity className="w-5 h-5 text-pulse-primary" />
        <span className="font-light text-lg text-pulse-text-primary">
          Pulse
        </span>
        <span className="text-[10px] uppercase tracking-wider text-pulse-text-muted bg-pulse-bg px-1.5 py-0.5 rounded">
          Beta
        </span>
      </div>

      <nav className="flex-1 p-3 space-y-0.5">
        {NAV.map(
          ({
            href,
            label,
            icon: Icon,
            badge,
          }: {
            href: string;
            label: string;
            icon: React.ElementType;
            badge?: string;
          }) => {
            const active =
              pathname === href ||
              (href !== "/dashboard" && pathname?.startsWith(href));
            return (
              <Link
                key={href}
                href={href}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors ${
                  active
                    ? "bg-pulse-primary/10 text-pulse-primary border-l-2 border-pulse-primary -ml-[2px] pl-[14px]"
                    : "text-pulse-text-secondary hover:text-pulse-text-primary hover:bg-pulse-surface-raised"
                }`}
              >
                <Icon className="w-4 h-4 shrink-0" />
                <span className="flex-1">{label}</span>
                {badge && (
                  <span className="text-[9px] font-mono uppercase tracking-wider px-1.5 py-0.5 rounded-full bg-pulse-primary/15 text-pulse-primary border border-pulse-primary/20">
                    {badge}
                  </span>
                )}
              </Link>
            );
          },
        )}
      </nav>

      <div className="p-3 border-t border-pulse-border space-y-2">
        {!moodLoggedToday && (
          <Link href="/log" className="block">
            <Button
              size="sm"
              className="w-full bg-pulse-accent text-pulse-bg hover:bg-pulse-accent/90 relative"
            >
              {moodLoggedToday ? null : (
                <span className="absolute -top-0.5 -right-0.5 w-2 h-2 rounded-full bg-pulse-danger animate-pulse" />
              )}
              Log mood today
            </Button>
          </Link>
        )}
        <div className="flex items-center gap-2 px-2 py-1.5">
          <Avatar className="h-8 w-8 bg-pulse-primary/20">
            <AvatarFallback className="text-pulse-primary text-xs">
              A
            </AvatarFallback>
          </Avatar>
          <div className="min-w-0 flex-1">
            <p className="text-xs text-pulse-text-secondary truncate">
              Anonymous user
            </p>
            <button
              type="button"
              className="text-xs text-pulse-text-muted hover:text-pulse-text-secondary flex items-center gap-1"
            >
              <Lock className="w-3 h-3" /> Settings
            </button>
          </div>
        </div>
      </div>

      <p className="px-4 py-2 text-[10px] text-pulse-text-muted border-t border-pulse-border">
        If you&apos;re in crisis: 988 Suicide & Crisis Lifeline
      </p>
    </aside>
  );
}
