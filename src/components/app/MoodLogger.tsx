"use client";

import { useState, useCallback } from "react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Slider } from "@/components/ui/slider";
import type { MoodLog } from "@/lib/types";

const MOOD_TAGS = [
  "work-stress",
  "poor-sleep",
  "social",
  "exercise",
  "good-day",
  "burnout",
  "calm",
] as const;

const PLACEHOLDERS = [
  "Too many meetings today...",
  "Slept terribly last night...",
  "Actually felt clear-headed...",
];

type MoodLoggerProps = {
  onSubmit: (data: {
    score: number;
    energy?: number;
    anxiety?: number;
    note?: string;
    tags?: string[];
  }) => void | Promise<void>;
  existingEntry: MoodLog | null;
};

export function MoodLogger({ onSubmit, existingEntry }: MoodLoggerProps) {
  const [score, setScore] = useState(existingEntry?.score ?? 7);
  const [energy, setEnergy] = useState(existingEntry?.energy ?? undefined);
  const [anxiety, setAnxiety] = useState(existingEntry?.anxiety ?? undefined);
  const [note, setNote] = useState(existingEntry?.note ?? "");
  const [tags, setTags] = useState<string[]>(existingEntry?.tags ?? []);
  const [loading, setLoading] = useState(false);

  const placeholder = PLACEHOLDERS[Math.floor(Date.now() / 86400000) % PLACEHOLDERS.length];

  const toggleTag = useCallback((tag: string) => {
    setTags((prev) =>
      prev.includes(tag) ? prev.filter((t) => t !== tag) : prev.length >= 3 ? prev : [...prev, tag]
    );
  }, []);

  const handleSubmit = async () => {
    setLoading(true);
    try {
      await onSubmit({
        score,
        energy: energy ?? undefined,
        anxiety: anxiety ?? undefined,
        note: note.trim() || undefined,
        tags: tags.length ? tags : undefined,
      });
    } finally {
      setLoading(false);
    }
  };

  const scoreColor =
    score <= 3 ? "text-pulse-danger" : score <= 6 ? "text-pulse-accent-warm" : "text-pulse-accent";

  return (
    <div className="space-y-8">
      <div>
        <div className={`text-5xl font-mono font-light mb-2 ${scoreColor}`}>{score}</div>
        <Label className="text-pulse-text-muted text-sm">
          Rough (1) … Great (10)
        </Label>
        <Slider
          value={[score]}
          onValueChange={([v]) => setScore(v)}
          min={1}
          max={10}
          step={1}
          className="mt-2"
        />
      </div>

      <div>
        <Label className="text-pulse-text-secondary">Energy level (optional)</Label>
        <Slider
          value={[energy ?? 5]}
          onValueChange={([v]) => setEnergy(v)}
          min={1}
          max={10}
          step={1}
          className="mt-2"
        />
      </div>

      <div>
        <Label className="text-pulse-text-secondary">Anxiety level (optional) — 1 = calm, 10 = very anxious</Label>
        <Slider
          value={[anxiety ?? 5]}
          onValueChange={([v]) => setAnxiety(v)}
          min={1}
          max={10}
          step={1}
          className="mt-2"
        />
      </div>

      <div>
        <Label className="text-pulse-text-secondary">One sentence — what&apos;s driving this? (optional)</Label>
        <Textarea
          value={note}
          onChange={(e) => setNote(e.target.value.slice(0, 500))}
          placeholder={placeholder}
          className="mt-2 bg-pulse-surface border-pulse-border resize-none"
          rows={2}
        />
        <p className="text-xs text-pulse-text-muted mt-1">{note.length}/500</p>
      </div>

      <div>
        <Label className="text-pulse-text-secondary block mb-2">Tags (optional, max 3)</Label>
        <div className="flex flex-wrap gap-2">
          {MOOD_TAGS.map((tag) => (
            <button
              key={tag}
              type="button"
              onClick={() => toggleTag(tag)}
              className={`px-3 py-1.5 rounded-full text-sm border transition-colors ${
                tags.includes(tag)
                  ? "bg-pulse-primary/20 border-pulse-primary text-pulse-text-primary"
                  : "bg-pulse-surface border-pulse-border text-pulse-text-secondary hover:border-pulse-text-muted"
              }`}
            >
              {tag.replace(/-/g, " ")}
            </button>
          ))}
        </div>
      </div>

      <Button
        onClick={handleSubmit}
        disabled={loading}
        className="w-full"
      >
        {loading ? "Logging…" : "Log it →"}
      </Button>
    </div>
  );
}
