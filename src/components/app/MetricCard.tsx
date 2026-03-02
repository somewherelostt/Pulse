"use client";

import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

type Severity = "healthy" | "warning" | "critical";
type DeltaDirection = "up-bad" | "up-good" | "down-bad" | "down-good";

type MetricCardProps = {
  label: string;
  value: string | number;
  delta?: string;
  deltaDirection?: DeltaDirection;
  severity?: Severity;
  tooltip?: string;
};

const severityBar: Record<Severity, string> = {
  healthy: "bg-pulse-accent",
  warning: "bg-pulse-accent-warm",
  critical: "bg-pulse-danger",
};

const deltaColor: Record<DeltaDirection, string> = {
  "up-bad": "text-pulse-danger",
  "up-good": "text-pulse-accent",
  "down-bad": "text-pulse-danger",
  "down-good": "text-pulse-accent",
};

export function MetricCard({
  label,
  value,
  delta,
  deltaDirection = "up-bad",
  severity = "healthy",
  tooltip,
}: MetricCardProps) {
  const content = (
    <div
      className="bg-pulse-surface border border-pulse-border rounded-xl p-4 transition-all hover:-translate-y-0.5 hover:shadow-lg overflow-hidden"
      style={{ boxShadow: "var(--shadow-sm)" }}
    >
      <div className={`h-[3px] w-full -mx-4 -mt-4 mb-3 ${severityBar[severity]}`} />
      <div className="font-mono text-2xl md:text-3xl text-pulse-text-primary">{value}</div>
      <p className="text-[11px] uppercase tracking-widest text-pulse-text-muted mt-1">{label}</p>
      {delta && (
        <p className={`text-xs mt-1 ${deltaColor[deltaDirection]}`}>{delta}</p>
      )}
    </div>
  );

  if (tooltip) {
    return (
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger asChild>
            <div className="cursor-help">{content}</div>
          </TooltipTrigger>
          <TooltipContent>{tooltip}</TooltipContent>
        </Tooltip>
      </TooltipProvider>
    );
  }
  return content;
}
