"use client";

import { useState, useEffect, useMemo } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Shield, Network, Heart, Zap, ArrowRight, Users, AlertCircle } from "lucide-react";
import { Sidebar } from "@/components/app/Sidebar";

// Demo fingerprint data — represents a behavioral profile
const DEMO_FINGERPRINT = {
  meeting_load: 0.72,
  mood_volatility: 0.58,
  after_hours: 0.44,
  recovery_speed: 0.61,
  focus_depth: 0.67,
  circadian_alignment: 0.53,
};

const FINGERPRINT_LABELS: Record<string, string> = {
  meeting_load: "Meeting Load",
  mood_volatility: "Mood Volatility",
  after_hours: "After-Hours Work",
  recovery_speed: "Recovery Speed",
  focus_depth: "Focus Depth",
  circadian_alignment: "Circadian Alignment",
};

// Demo peer matches
const DEMO_MATCHES = [
  {
    id: "peer-1",
    initials: "M",
    similarity: 0.91,
    patterns: ["High meeting density", "Evening energy dip"],
    status: "Recovered",
    statusColor: "text-emerald-400",
    context: "Similar burnout spiral 3 weeks ago — emerged after reducing back-to-back meetings by 40%.",
    weeksAgo: 3,
  },
  {
    id: "peer-2",
    initials: "S",
    similarity: 0.87,
    patterns: ["After-hours drift", "Focus fragmentation"],
    status: "Active",
    statusColor: "text-pulse-accent-warm",
    context: "Currently navigating the same pattern. Connected 2 others through it already.",
    weeksAgo: null,
  },
  {
    id: "peer-3",
    initials: "R",
    similarity: 0.83,
    patterns: ["Calendar overload", "Mood dip lag"],
    status: "Recovered",
    statusColor: "text-emerald-400",
    context: "Hit the same wall during a product launch. Found anchor points that helped.",
    weeksAgo: 6,
  },
];

const PRINCIPLES = [
  { icon: Shield, title: "Anonymous by default", body: "No name, no photo. Your identity stays yours." },
  { icon: Network, title: "Pattern-matched", body: "Behavioral fingerprint only — not a diagnosis." },
  { icon: Heart, title: "Human, not AI", body: "Real people who've been where you are." },
];

function FingerprintRadar({ data }: { data: Record<string, number> }) {
  const keys = Object.keys(data);
  const n = keys.length;
  const cx = 100;
  const cy = 100;
  const r = 70;

  const points = keys.map((_, i) => {
    const angle = (i / n) * 2 * Math.PI - Math.PI / 2;
    return {
      x: cx + r * Math.cos(angle),
      y: cy + r * Math.sin(angle),
      labelX: cx + (r + 18) * Math.cos(angle),
      labelY: cy + (r + 18) * Math.sin(angle),
    };
  });

  const dataPoints = keys.map((k, i) => {
    const angle = (i / n) * 2 * Math.PI - Math.PI / 2;
    const val = data[k];
    return {
      x: cx + r * val * Math.cos(angle),
      y: cy + r * val * Math.sin(angle),
    };
  });

  const polyline = dataPoints.map((p) => `${p.x},${p.y}`).join(" ");
  const outerPolygon = points.map((p) => `${p.x},${p.y}`).join(" ");

  return (
    <svg viewBox="0 0 200 200" className="w-full max-w-[260px]">
      {/* Grid rings */}
      {[0.25, 0.5, 0.75, 1].map((t) => (
        <polygon
          key={t}
          points={points.map((p) => {
            const angle = Math.atan2(p.y - cy, p.x - cx);
            const rr = r * t;
            return `${cx + rr * Math.cos(angle)},${cy + rr * Math.sin(angle)}`;
          }).join(" ")}
          fill="none"
          stroke="rgba(110,123,242,0.1)"
          strokeWidth="1"
        />
      ))}
      {/* Axis lines */}
      {points.map((p, i) => (
        <line key={i} x1={cx} y1={cy} x2={p.x} y2={p.y} stroke="rgba(110,123,242,0.15)" strokeWidth="1" />
      ))}
      {/* Data polygon */}
      <motion.polygon
        points={polyline}
        fill="rgba(110,123,242,0.15)"
        stroke="rgba(110,123,242,0.6)"
        strokeWidth="1.5"
        initial={{ opacity: 0, scale: 0 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.8, ease: "easeOut" }}
        style={{ transformOrigin: `${cx}px ${cy}px` }}
      />
      {/* Data points */}
      {dataPoints.map((p, i) => (
        <motion.circle
          key={i}
          cx={p.x}
          cy={p.y}
          r="3"
          fill="var(--color-primary, #6E7BF2)"
          initial={{ opacity: 0, r: 0 }}
          animate={{ opacity: 1, r: 3 }}
          transition={{ delay: 0.6 + i * 0.1, duration: 0.3 }}
        />
      ))}
    </svg>
  );
}

export default function ConstellationPage() {
  const [joined, setJoined] = useState(false);
  const [loading, setLoading] = useState(false);
  const [selectedMatch, setSelectedMatch] = useState<string | null>(null);
  const [moodLoggedToday] = useState(false);

  // Check if user is in peer pool (demo: use localStorage)
  useEffect(() => {
    const inPool = localStorage.getItem("constellation_joined") === "1";
    setJoined(inPool);
  }, []);

  const handleJoin = async () => {
    setLoading(true);
    // Simulate network join
    await new Promise((r) => setTimeout(r, 1200));
    localStorage.setItem("constellation_joined", "1");
    setJoined(true);
    setLoading(false);
  };

  const handleLeave = async () => {
    localStorage.removeItem("constellation_joined");
    setJoined(false);
    setSelectedMatch(null);
  };

  return (
    <div className="min-h-screen bg-pulse-bg">
      <Sidebar moodLoggedToday={moodLoggedToday} />
      <main className="pl-[220px] min-h-screen">
        <div className="p-6 max-w-5xl">
          {/* Header */}
          <div className="mb-8">
            <div className="flex items-center gap-2 mb-3">
              <div className="w-1.5 h-1.5 rounded-full bg-pulse-primary animate-pulse" />
              <span className="text-[10px] font-mono uppercase tracking-widest text-pulse-primary">
                Constellation · Layer 5
              </span>
            </div>
            <h1 className="text-2xl font-light text-pulse-text-primary mb-1">
              Peer Matching
            </h1>
            <p className="text-pulse-text-muted text-sm max-w-xl">
              Anonymously match with people who've been where you are now — based on behavioral patterns, not words.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Left col — fingerprint + pool status */}
            <div className="lg:col-span-1 space-y-4">
              {/* Behavioral fingerprint card */}
              <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                <div className="flex items-center justify-between mb-4">
                  <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted">
                    Your Fingerprint
                  </p>
                  <span className="text-[10px] font-mono text-pulse-primary bg-pulse-primary/10 px-2 py-0.5 rounded-full">
                    6-dim
                  </span>
                </div>
                <div className="flex justify-center mb-4">
                  <FingerprintRadar data={DEMO_FINGERPRINT} />
                </div>
                <div className="space-y-2">
                  {Object.entries(DEMO_FINGERPRINT).map(([k, v]) => (
                    <div key={k} className="flex items-center gap-2">
                      <span className="text-[10px] text-pulse-text-muted w-28 truncate">{FINGERPRINT_LABELS[k]}</span>
                      <div className="flex-1 h-1 bg-pulse-border rounded-full overflow-hidden">
                        <motion.div
                          className="h-full bg-pulse-primary rounded-full"
                          initial={{ width: 0 }}
                          animate={{ width: `${v * 100}%` }}
                          transition={{ duration: 0.8, delay: 0.2 }}
                        />
                      </div>
                      <span className="text-[10px] font-mono text-pulse-text-muted w-8 text-right">
                        {Math.round(v * 100)}%
                      </span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Pool status card */}
              <div className={`border rounded-xl p-4 transition-all ${
                joined
                  ? "bg-pulse-primary/5 border-pulse-primary/30"
                  : "bg-pulse-surface border-pulse-border"
              }`}>
                <div className="flex items-center gap-2 mb-3">
                  <div className={`w-2 h-2 rounded-full ${joined ? "bg-pulse-primary animate-pulse" : "bg-pulse-border"}`} />
                  <span className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted">
                    {joined ? "In peer pool" : "Not in pool"}
                  </span>
                </div>
                {joined ? (
                  <div>
                    <p className="text-xs text-pulse-text-secondary mb-3">
                      Your anonymized behavioral fingerprint is active. You'll be matched when a compatible peer is found.
                    </p>
                    <div className="flex items-center gap-2 text-[10px] font-mono text-pulse-text-muted mb-4">
                      <Users className="w-3 h-3" />
                      <span>~140 active in pool · 12 matched today</span>
                    </div>
                    <button
                      onClick={handleLeave}
                      className="text-xs text-pulse-text-muted hover:text-pulse-danger transition-colors"
                    >
                      Leave pool →
                    </button>
                  </div>
                ) : (
                  <div>
                    <p className="text-xs text-pulse-text-secondary mb-4">
                      Share your behavioral fingerprint anonymously to find people who've navigated the same patterns.
                    </p>
                    <button
                      onClick={handleJoin}
                      disabled={loading}
                      className="w-full py-2 px-4 rounded-lg bg-pulse-primary text-white text-sm font-medium hover:bg-pulse-primary/90 transition-all disabled:opacity-50 flex items-center justify-center gap-2"
                    >
                      {loading ? (
                        <>
                          <div className="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                          Joining...
                        </>
                      ) : (
                        <>
                          Join peer pool
                          <ArrowRight className="w-3.5 h-3.5" />
                        </>
                      )}
                    </button>
                  </div>
                )}
              </div>

              {/* Safety note */}
              <div className="flex items-start gap-2 p-3 rounded-lg bg-pulse-surface-raised border border-pulse-border/50">
                <AlertCircle className="w-3.5 h-3.5 text-pulse-text-muted shrink-0 mt-0.5" />
                <p className="text-[10px] text-pulse-text-muted leading-relaxed">
                  Sessions are 45 min max, end-to-end encrypted, and never recorded. If you're in crisis: call 988.
                </p>
              </div>
            </div>

            {/* Right col — matches + session */}
            <div className="lg:col-span-2 space-y-4">
              {!joined ? (
                /* Pre-join — show principles */
                <div className="space-y-4">
                  <div className="bg-pulse-surface border border-pulse-border rounded-xl p-6">
                    <h2 className="text-lg font-light text-pulse-text-primary mb-2">
                      The hardest part isn't the pattern.
                    </h2>
                    <p className="text-pulse-text-secondary text-sm mb-6">
                      It's feeling like nobody else has been through it.
                    </p>

                    {/* Mini constellation visual */}
                    <div className="relative h-48 rounded-xl border border-pulse-border bg-pulse-bg/50 overflow-hidden mb-6">
                      <ConstellationVisual />
                    </div>

                    <div className="grid grid-cols-3 gap-4">
                      {PRINCIPLES.map(({ icon: Icon, title, body }, i) => (
                        <motion.div
                          key={i}
                          className="text-center"
                          initial={{ opacity: 0, y: 10 }}
                          animate={{ opacity: 1, y: 0 }}
                          transition={{ delay: i * 0.1 }}
                        >
                          <Icon className="w-5 h-5 text-pulse-primary mx-auto mb-2" />
                          <p className="text-xs font-medium text-pulse-text-primary mb-1">{title}</p>
                          <p className="text-[10px] text-pulse-text-muted">{body}</p>
                        </motion.div>
                      ))}
                    </div>
                  </div>
                </div>
              ) : (
                /* Post-join — show matches */
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <h2 className="text-sm font-medium text-pulse-text-primary">Pattern matches</h2>
                    <span className="text-[10px] font-mono text-pulse-text-muted uppercase tracking-widest">
                      {DEMO_MATCHES.length} found
                    </span>
                  </div>

                  {DEMO_MATCHES.map((match, i) => (
                    <motion.div
                      key={match.id}
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: i * 0.15 }}
                    >
                      <button
                        className={`w-full text-left bg-pulse-surface border rounded-xl p-4 transition-all hover:-translate-y-0.5 hover:shadow-lg ${
                          selectedMatch === match.id
                            ? "border-pulse-primary/50 bg-pulse-primary/5"
                            : "border-pulse-border hover:border-pulse-border/80"
                        }`}
                        onClick={() => setSelectedMatch(selectedMatch === match.id ? null : match.id)}
                      >
                        <div className="flex items-start gap-4">
                          <div className="w-10 h-10 rounded-full bg-pulse-surface-raised border border-pulse-border flex items-center justify-center text-sm font-mono text-pulse-text-secondary shrink-0">
                            {match.initials}
                          </div>
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 mb-1">
                              <span className="text-[10px] font-mono text-pulse-text-muted">
                                {Math.round(match.similarity * 100)}% pattern match
                              </span>
                              <div className="flex-1 h-px bg-pulse-border" />
                              <span className={`text-[10px] font-mono ${match.statusColor}`}>
                                {match.status}
                                {match.weeksAgo && ` · ${match.weeksAgo}w ago`}
                              </span>
                            </div>
                            <div className="flex flex-wrap gap-1.5 mb-2">
                              {match.patterns.map((p) => (
                                <span key={p} className="text-[10px] px-2 py-0.5 rounded-full bg-pulse-primary/10 text-pulse-text-secondary border border-pulse-primary/10">
                                  {p}
                                </span>
                              ))}
                            </div>

                            <AnimatePresence>
                              {selectedMatch === match.id && (
                                <motion.div
                                  initial={{ opacity: 0, height: 0 }}
                                  animate={{ opacity: 1, height: "auto" }}
                                  exit={{ opacity: 0, height: 0 }}
                                  transition={{ duration: 0.2 }}
                                >
                                  <p className="text-xs text-pulse-text-secondary mb-3 pt-2 border-t border-pulse-border/50">
                                    {match.context}
                                  </p>
                                  <button className="flex items-center gap-2 text-xs font-medium text-pulse-primary hover:text-pulse-primary/80 transition-colors">
                                    <Zap className="w-3.5 h-3.5" />
                                    Request anonymous session
                                    <ArrowRight className="w-3 h-3" />
                                  </button>
                                </motion.div>
                              )}
                            </AnimatePresence>
                          </div>
                        </div>
                      </button>
                    </motion.div>
                  ))}

                  {/* How it works */}
                  <div className="bg-pulse-surface border border-pulse-border rounded-xl p-4">
                    <p className="text-[10px] font-mono uppercase tracking-widest text-pulse-text-muted mb-3">
                      How sessions work
                    </p>
                    <div className="space-y-2">
                      {[
                        { step: "01", text: "Request is sent anonymously. Peer accepts or declines." },
                        { step: "02", text: "End-to-end encrypted video room opens. 45-min limit." },
                        { step: "03", text: "No names, no recording. Only the conversation exists." },
                        { step: "04", text: "Rate the session (1–5) to improve matching for others." },
                      ].map(({ step, text }) => (
                        <div key={step} className="flex gap-3 items-start">
                          <span className="text-[10px] font-mono text-pulse-primary/60 w-5 shrink-0">{step}</span>
                          <p className="text-[11px] text-pulse-text-secondary">{text}</p>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

// Mini constellation visualization
function ConstellationVisual() {
  const nodes = useMemo(() => [
    { id: 1, x: 18, y: 28, initial: "A" },
    { id: 2, x: 38, y: 72, initial: "M" },
    { id: 3, x: 68, y: 18, initial: "S" },
    { id: 4, x: 82, y: 58, initial: "R" },
    { id: 5, x: 52, y: 42, initial: "YOU", isUser: true },
    { id: 6, x: 28, y: 14, initial: "K" },
    { id: 7, x: 14, y: 62, initial: "J" },
  ], []);

  const connections = [[5, 1], [5, 3], [5, 4], [1, 6], [7, 2]];

  return (
    <>
      <svg className="absolute inset-0 w-full h-full">
        {connections.map(([fromId, toId], idx) => {
          const from = nodes.find((n) => n.id === fromId)!;
          const to = nodes.find((n) => n.id === toId)!;
          return (
            <motion.line
              key={idx}
              x1={`${from.x}%`} y1={`${from.y}%`}
              x2={`${to.x}%`} y2={`${to.y}%`}
              stroke="rgba(110,123,242,0.25)"
              strokeWidth="1"
              initial={{ pathLength: 0 }}
              animate={{ pathLength: 1 }}
              transition={{ duration: 1.5, delay: idx * 0.2 }}
            />
          );
        })}
      </svg>
      {nodes.map((node) => (
        <motion.div
          key={node.id}
          className="absolute -translate-x-1/2 -translate-y-1/2"
          style={{ left: `${node.x}%`, top: `${node.y}%` }}
          initial={{ opacity: 0, scale: 0 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ type: "spring", stiffness: 120, delay: node.id * 0.08 }}
        >
          <motion.div
            className={`w-8 h-8 rounded-full flex items-center justify-center text-[9px] font-mono border ${
              node.isUser
                ? "bg-pulse-primary text-white border-pulse-primary shadow-[0_0_16px_rgba(110,123,242,0.5)]"
                : "bg-pulse-surface-raised text-pulse-text-muted border-pulse-border"
            }`}
            animate={node.isUser ? { scale: [1, 1.08, 1] } : { y: [0, -3, 0] }}
            transition={{ duration: node.isUser ? 2 : 3 + node.id * 0.3, repeat: Infinity, ease: "easeInOut" }}
          >
            {node.initial}
          </motion.div>
        </motion.div>
      ))}
    </>
  );
}
