"use client";

import React from "react";
import { motion } from "framer-motion";
import { Check, ArrowRight, Shield, Database, Sparkles } from "lucide-react";
import { HoverVisual } from "./ui/HoverVisual";

const PrivacyCTA = () => {
  const trustCards = [
    "Anonymous sessions",
    "Row-level security",
    "Peer sessions encrypted",
    "One-click delete",
    "Open source",
    "WCAG 2.1 AA",
  ];

  return (
    <div className="bg-pulse-bg">
      {/* Privacy Section */}
      <section className="py-32 max-w-[1200px] mx-auto px-6 grid grid-cols-1 lg:grid-cols-2 gap-20 items-center">
        <div>
          <h2 className="text-4xl md:text-5xl font-light text-pulse-text-primary mb-6 leading-tight">
            Everything runs on your data. <br />
            <span className="text-pulse-text-secondary">Nothing runs on your identity.</span>
          </h2>
          <p className="text-lg text-pulse-text-secondary mb-8 leading-relaxed">
            We process behavioral features, not raw content. Your journal entries never leave your device. Your peer sessions are never recorded. You own every byte.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {trustCards.map((card, idx) => (
            <motion.div
              key={idx}
              className="p-4 rounded-xl border border-pulse-border bg-pulse-surface flex items-center gap-3"
              initial={{ opacity: 0, x: 20 }}
              whileInView={{ opacity: 1, x: 0 }}
              transition={{ delay: idx * 0.1 }}
            >
              <div className="w-5 h-5 rounded-full bg-pulse-accent/10 flex items-center justify-center">
                <Check className="w-3 h-3 text-pulse-accent" />
              </div>
              <span className="text-sm text-pulse-text-secondary">{card}</span>
            </motion.div>
          ))}
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-32 relative overflow-hidden">
        <div className="absolute inset-0 bg-pulse-primary/5 -z-10" />
        
        <div className="max-w-[800px] mx-auto px-6 text-center">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            className="mb-12"
          >
            <h2 className="text-5xl md:text-7xl font-light text-pulse-text-primary mb-6 leading-tight">
              Your <HoverVisual text="data">
                <div className="p-4 flex flex-col gap-3">
                  <div className="flex items-center gap-2">
                    <Database className="w-4 h-4 text-pulse-primary" />
                    <span className="text-[10px] font-mono text-pulse-text-secondary uppercase">Behavioral Inventory</span>
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    {["Calendar", "Sleep", "Health", "Focus"].map((item, i) => (
                      <div key={i} className="px-2 py-1 rounded bg-pulse-bg/50 border border-pulse-border/50 text-[9px] font-mono text-pulse-text-muted">
                        ● {item}
                      </div>
                    ))}
                  </div>
                </div>
              </HoverVisual> is already telling a <HoverVisual text="story">
                <div className="p-4 max-w-[200px]">
                  <Sparkles className="w-4 h-4 text-pulse-accent mb-2" />
                  <p className="text-[10px] text-pulse-text-secondary leading-relaxed italic">
                    "A pattern of late-night browsing followed by calendar avoidance suggests an emerging burnout spiral."
                  </p>
                </div>
              </HoverVisual>. <br />
              <span className="pulse-gradient font-semibold">
                Let <HoverVisual text="Pulse">
                  <div className="p-4">
                    <div className="flex items-center gap-2 mb-3">
                      <div className="w-2 h-2 rounded-full bg-pulse-accent animate-pulse" />
                      <span className="text-[10px] font-mono text-pulse-accent uppercase tracking-widest">ACTIVE MONITORING</span>
                    </div>
                    <div className="space-y-1">
                      <div className="h-0.5 w-full bg-pulse-accent/20 rounded-full overflow-hidden">
                        <motion.div animate={{ x: ["-100%", "100%"] }} transition={{ repeat: Infinity, duration: 1.5, ease: "linear" }} className="h-full w-1/3 bg-pulse-accent" />
                      </div>
                      <div className="h-0.5 w-full bg-pulse-accent/20 rounded-full overflow-hidden">
                        <motion.div animate={{ x: ["-100%", "100%"] }} transition={{ repeat: Infinity, duration: 2, ease: "linear", delay: 0.5 }} className="h-full w-1/4 bg-pulse-accent" />
                      </div>
                    </div>
                  </div>
                </HoverVisual> read it.
              </span>
            </h2>
            <p className="text-lg md:text-xl text-pulse-text-secondary leading-relaxed">
              Join early access. Anonymous. Free to start. <br />
              No commitment. No credit card.
            </p>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            whileInView={{ opacity: 1, scale: 1 }}
            className="flex flex-col md:flex-row items-center gap-4 max-w-[500px] mx-auto mb-8"
          >
            <div className="w-full relative">
              <input
                type="email"
                placeholder="abumaaz2004@email.com (optional)"
                className="w-full h-14 bg-pulse-surface border border-pulse-border rounded-full px-8 text-pulse-text-primary focus:outline-none focus:border-pulse-primary transition-colors pr-32"
              />
              <button className="absolute right-2 top-2 bottom-2 px-6 rounded-full bg-pulse-primary text-white font-medium flex items-center gap-2 hover:shadow-[0_0_15px_rgba(110,123,242,0.4)] transition-all">
                Join
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          </motion.div>

          <p className="text-xs text-pulse-text-muted mb-12">
            We hate spam more than burnout. We'll only email you when it's ready.
          </p>

          <div className="flex flex-wrap justify-center gap-4">
            {["Free tier available", "Open source", "Privacy first"].map((pill, idx) => (
               <div key={idx} className="px-4 py-1.5 rounded-full border border-pulse-border bg-pulse-surface text-[10px] font-mono text-pulse-text-muted uppercase tracking-widest">
                 {pill}
               </div>
            ))}
          </div>
        </div>
      </section>
    </div>
  );
};

export default PrivacyCTA;
