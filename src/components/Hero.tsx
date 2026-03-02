"use client";

import React from "react";
import { motion } from "framer-motion";
import { ArrowRight, AlertTriangle, Moon, Calendar, Activity } from "lucide-react";
import { HoverVisual } from "./ui/HoverVisual";

const Hero = () => {
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.15,
        delayChildren: 0.3,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.8, ease: [0.16, 1, 0.3, 1] as any },
    },
  };

  return (
    <section className="relative min-h-screen flex flex-col items-center justify-center pt-32 pb-20 overflow-hidden">
      {/* Background radial glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] bg-pulse-primary/5 rounded-full blur-[120px] -z-10" />

      <motion.div
        className="max-w-[800px] mx-auto px-6 text-center"
        variants={containerVariants}
        initial="hidden"
        animate="visible"
      >
        <motion.div variants={itemVariants} className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-pulse-border bg-pulse-surface-raised mb-8">
          <motion.div
            className="w-1.5 h-1.5 rounded-full bg-pulse-primary"
            animate={{ opacity: [1, 0.4, 1] }}
            transition={{ duration: 1.5, repeat: Infinity, ease: "easeInOut" }}
          />
          <span className="text-[11px] font-mono uppercase tracking-widest text-pulse-primary">
            Live Behavioral Monitoring
          </span>
        </motion.div>

        <motion.h1 variants={itemVariants} className="text-5xl md:text-7xl lg:text-8xl font-light tracking-tight text-pulse-text-primary leading-[1.1] mb-8">
          Your mind leaves <HoverVisual text="traces" className="mx-2">
            <div className="p-4 flex flex-col gap-2">
              <div className="flex items-center gap-2 mb-2">
                <div className="w-1.5 h-1.5 rounded-full bg-pulse-accent-warm animate-pulse" />
                <span className="text-[10px] font-mono text-pulse-text-secondary uppercase">Spiral pattern detected</span>
              </div>
              <div className="flex gap-1 h-8 items-end">
                {[40, 70, 45, 90, 65, 30, 85].map((h, i) => (
                  <motion.div
                    key={i}
                    className="flex-1 bg-pulse-accent-warm/20 rounded-t-[2px]"
                    initial={{ height: 0 }}
                    animate={{ height: `${h}%` }}
                    transition={{ duration: 0.5, delay: i * 0.05 }}
                  />
                ))}
              </div>
              <p className="text-[10px] text-pulse-text-muted mt-2">Correlated to high meeting volume (78%)</p>
            </div>
          </HoverVisual>
          everywhere.<br />
          <HoverVisual text="Pulse" className="pulse-gradient font-semibold mr-2">
            <div className="p-4">
              <div className="flex items-center justify-between mb-3">
                <span className="text-[10px] font-mono text-pulse-text-secondary tracking-widest">SIGNAL QUALITY</span>
                <span className="text-[10px] font-mono text-pulse-primary">92%</span>
              </div>
              <div className="relative h-1 w-full bg-pulse-border rounded-full overflow-hidden">
                <motion.div 
                  className="absolute left-0 top-0 bottom-0 bg-pulse-primary" 
                  initial={{ width: 0 }} 
                  animate={{ width: "92%" }}
                  transition={{ duration: 1.2, delay: 0.2 }}
                />
              </div>
              <div className="mt-4 flex gap-2">
                <div className="px-2 py-0.5 rounded-sm bg-pulse-primary/10 border border-pulse-primary/20 text-[8px] font-mono text-pulse-primary uppercase">E2E ENCRYPTED</div>
                <div className="px-2 py-0.5 rounded-sm bg-pulse-accent/10 border border-pulse-accent/20 text-[8px] font-mono text-pulse-accent uppercase">LOCAL OPS</div>
              </div>
            </div>
          </HoverVisual> 
          reads them.
        </motion.h1>

        <motion.p variants={itemVariants} className="max-w-[600px] mx-auto text-lg md:text-xl text-pulse-text-secondary leading-relaxed mb-10">
          Before you feel it, your calendar, sleep, and browser already know.
        </motion.p>

        <motion.div variants={itemVariants} className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16">
          <button className="group h-12 px-8 rounded-full bg-pulse-primary text-white font-medium flex items-center gap-2 hover:scale-[1.02] transition-all shadow-[0_0_20px_rgba(110,123,242,0.3)] hover:shadow-[0_0_30px_rgba(110,123,242,0.5)]">
            Get Early Access
            <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
          </button>
          <button className="h-12 px-6 text-pulse-text-secondary hover:text-pulse-text-primary transition-colors flex items-center gap-2 group">
            See how it works 
            <span className="group-hover:translate-y-1 transition-transform">↓</span>
          </button>
        </motion.div>

        <motion.p variants={itemVariants} className="text-xs text-pulse-text-muted mb-20">
          Anonymous by default. No email required to start.
        </motion.p>
      </motion.div>

      {/* Hero Visual Dashboard Preview */}
      <motion.div
        className="relative w-full max-w-[1000px] mx-auto px-6"
        initial={{ opacity: 0, y: 40 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 1.2, delay: 0.8, ease: [0.16, 1, 0.3, 1] }}
      >
        <div className="relative rounded-[20px] border border-pulse-border bg-pulse-surface p-8 shadow-2xl overflow-hidden group">
          <div className="absolute top-0 left-0 right-0 h-px bg-gradient-to-right from-transparent via-pulse-primary/20 to-transparent" />
          
          <div className="flex items-center justify-between mb-12">
            <h3 className="text-sm font-mono text-pulse-text-secondary uppercase tracking-widest">Unified Timeline</h3>
            <div className="flex gap-4">
              <div className="flex items-center gap-2 text-[10px] text-pulse-text-muted font-mono">
                <div className="w-2 h-2 rounded-sm bg-blue-500/50" /> CALENDAR
              </div>
              <div className="flex items-center gap-2 text-[10px] text-pulse-text-muted font-mono">
                <div className="w-2 h-2 rounded-sm bg-emerald-500/50" /> SLEEP
              </div>
              <div className="flex items-center gap-2 text-[10px] text-pulse-text-muted font-mono">
                <div className="w-2 h-2 rounded-sm bg-pulse-primary" /> MOOD
              </div>
            </div>
          </div>

          <div className="space-y-10 relative">
            {/* Red vertical line */}
            <motion.div 
              className="absolute left-[80%] top-0 bottom-0 w-px border-l border-dashed border-pulse-danger z-10"
              initial={{ height: 0 }}
              animate={{ height: "100%" }}
              transition={{ delay: 2.5, duration: 1 }}
            >
              <div className="absolute -top-6 left-1/2 -translate-x-1/2 whitespace-nowrap text-[10px] text-pulse-danger font-mono bg-pulse-surface px-1">
                YOU FELT THIS HERE
              </div>
            </motion.div>

            {/* Amber shaded region */}
            <motion.div 
              className="absolute left-[60%] w-[15%] top-0 bottom-0 bg-pulse-accent-warm/5 border-x border-pulse-accent-warm/20 z-0"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 2, duration: 1 }}
            >
              <div className="absolute -top-6 left-1/2 -translate-x-1/2 whitespace-nowrap text-[10px] text-pulse-accent-warm font-mono bg-pulse-surface px-1">
                PULSE SAW THIS HERE
              </div>
            </motion.div>

            {/* Row 1: Calendar */}
            <div className="relative h-12 flex items-end gap-1">
              {Array.from({ length: 30 }).map((_, i) => (
                <motion.div
                  key={i}
                  className="flex-1 bg-blue-500/20 rounded-t-sm"
                  initial={{ height: 0 }}
                  animate={{ height: `${20 + Math.random() * 60}%` }}
                  transition={{ duration: 1.5, delay: 1 + i * 0.02 }}
                />
              ))}
            </div>

            {/* Row 2: Sleep */}
            <div className="relative h-12">
              <svg className="w-full h-full" preserveAspectRatio="none">
                <motion.path
                  d="M0 48 Q 50 10, 100 40 T 200 20 T 300 45 T 400 15 T 500 35 T 600 25 T 700 40 T 800 10 T 900 30 T 1000 48"
                  fill="none"
                  stroke="rgba(16, 185, 129, 0.3)"
                  strokeWidth="2"
                  initial={{ pathLength: 0 }}
                  animate={{ pathLength: 1 }}
                  transition={{ duration: 1.5, delay: 1.2 }}
                />
                <motion.path
                   d="M0 48 Q 50 10, 100 40 T 200 20 T 300 45 T 400 15 T 500 35 T 600 25 T 700 40 T 800 10 T 900 30 T 1000 48 V 48 H 0 Z"
                   fill="rgba(16, 185, 129, 0.05)"
                   initial={{ opacity: 0 }}
                   animate={{ opacity: 1 }}
                   transition={{ duration: 1, delay: 2.2 }}
                />
              </svg>
            </div>

            {/* Row 3: Mood */}
            <div className="relative h-12">
               <svg className="w-full h-full" preserveAspectRatio="none">
                <motion.path
                  d="M0 24 Q 50 20, 100 26 T 200 22 T 300 28 T 400 20 T 500 32 T 600 18 T 700 40 T 800 45 T 900 42 T 1000 44"
                  fill="none"
                  stroke="#6E7BF2"
                  strokeWidth="2"
                  initial={{ pathLength: 0 }}
                  animate={{ pathLength: 1 }}
                  transition={{ duration: 1.5, delay: 1.4 }}
                />
              </svg>
            </div>
          </div>
        </div>

        {/* Stat Pills */}
        <div className="flex flex-wrap justify-center gap-4 mt-8">
          {[
            { label: "73% sleep consistency", icon: Moon, color: "text-emerald-500", trend: "↓" },
            { label: "68% meeting density", icon: Calendar, color: "text-blue-500", trend: "↑" },
            { label: "drift detected", icon: AlertTriangle, color: "text-pulse-accent-warm", trend: "⚠" }
          ].map((pill, idx) => (
            <motion.div
              key={idx}
              className="px-4 py-2 rounded-full border border-pulse-border bg-pulse-surface-raised flex items-center gap-2 text-sm text-pulse-text-secondary"
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: 2 + idx * 0.1 }}
            >
              <pill.icon className={`w-3.5 h-3.5 ${pill.color}`} />
              <span className="font-mono text-xs uppercase tracking-wider">
                <span className={pill.color}>{pill.trend}</span> {pill.label}
              </span>
            </motion.div>
          ))}
        </div>
      </motion.div>
    </section>
  );
};

export default Hero;
