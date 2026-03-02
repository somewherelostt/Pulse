"use client";

import { Lock } from "lucide-react";

type EmptyStateProps = {
  title: string;
  subtitle?: string;
  progress?: { current: number; total: number };
  locked?: boolean;
  children?: React.ReactNode;
};

export function EmptyState({ title, subtitle, progress, locked, children }: EmptyStateProps) {
  return (
    <div className="bg-pulse-surface border border-pulse-border rounded-xl p-8 text-center">
      {locked && (
        <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-pulse-primary/10 text-pulse-primary mb-4">
          <Lock className="w-6 h-6" />
        </div>
      )}
      <h3 className="text-lg font-medium text-pulse-text-primary">{title}</h3>
      {subtitle && <p className="text-pulse-text-secondary text-sm mt-2">{subtitle}</p>}
      {progress != null && (
        <div className="mt-4">
          <div className="flex justify-between text-xs text-pulse-text-muted mb-1">
            <span>{progress.current} of {progress.total} days needed for pattern analysis</span>
          </div>
          <div className="h-2 bg-pulse-bg rounded-full overflow-hidden">
            <div
              className="h-full bg-pulse-primary rounded-full transition-all"
              style={{ width: `${Math.min(100, (progress.current / progress.total) * 100)}%` }}
            />
          </div>
        </div>
      )}
      {children && <div className="mt-6">{children}</div>}
    </div>
  );
}
