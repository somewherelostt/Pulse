"use client";

import React from "react";
import { motion } from "framer-motion";
import { Radar, Dna, Activity, Share2 } from "lucide-react";

const HowItWorks = () => {
  const features = [
    {
      icon: Radar,
      title: "Passive collection",
      body: "Connect your calendar, sleep tracker, and browser. Pulse runs in the background. You do nothing.",
      visual: (
        <div className="relative w-full h-32 flex items-center justify-center">
          <motion.div
            className="w-12 h-12 rounded-full border border-pulse-primary flex items-center justify-center bg-pulse-primary/10 relative z-10"
            animate={{ scale: [1, 1.1, 1] }}
            transition={{ duration: 2, repeat: Infinity }}
          >
            <Radar className="w-5 h-5 text-pulse-primary" />
          </motion.div>
          {Array.from({ length: 4 }).map((_, i) => (
             <motion.div
               key={i}
               className="absolute w-2 h-2 rounded-full bg-pulse-text-muted/30"
               initial={{ x: i % 2 === 0 ? -60 : 60, y: i < 2 ? -40 : 40 }}
               animate={{ x: 0, y: 0, opacity: [0, 1, 0] }}
               transition={{ duration: 2, repeat: Infinity, delay: i * 0.5, ease: "easeIn" }}
             />
          ))}
        </div>
      )
    },
    {
      icon: Dna,
      title: "Behavioral features",
      body: "We extract what matters: meeting density, sleep consistency, browsing entropy, circadian alignment. Not raw data — patterns.",
      visual: (
        <div className="w-full h-32 flex flex-col justify-center gap-2 px-8">
          {Array.from({ length: 3 }).map((_, i) => (
             <div key={i} className="h-2 w-full bg-pulse-bg border border-pulse-border rounded-full overflow-hidden">
               <motion.div
                 className="h-full bg-pulse-primary/40"
                 initial={{ width: "0%" }}
                 whileInView={{ width: `${30 + i * 20}%` }}
                 transition={{ duration: 1, delay: i * 0.2 }}
               />
             </div>
          ))}
        </div>
      )
    },
    {
      icon: Activity,
      title: "Drift detection",
      body: "Every day, Pulse compares your behavioral fingerprint to your personal baseline. Divergence detected silently.",
      visual: (
        <div className="w-full h-32 flex items-center justify-center">
          <svg width="120" height="60" viewBox="0 0 120 60">
            <motion.path
              d="M0 30 Q 30 30, 60 10 T 120 30"
              fill="none"
              stroke="rgba(110,123,242,0.3)"
              strokeWidth="2"
            />
            <motion.path
              d="M0 30 Q 30 30, 60 50 T 120 30"
              fill="none"
              stroke="#F7B731"
              strokeWidth="2"
              initial={{ pathLength: 0 }}
              whileInView={{ pathLength: 1 }}
              transition={{ duration: 1.5 }}
            />
          </svg>
        </div>
      )
    },
    {
      icon: Share2,
      title: "Peer matching",
      body: "When your pattern matches someone else's recovery story, we introduce you. No names. No profiles. Just pattern recognition and human connection.",
      visual: (
        <div className="w-full h-32 flex items-center justify-center gap-8 relative">
           <div className="w-8 h-8 rounded-full border border-pulse-border bg-pulse-surface" />
           <motion.div
             className="h-px bg-pulse-accent w-16"
             initial={{ scaleX: 0 }}
             whileInView={{ scaleX: 1 }}
             transition={{ duration: 1 }}
           />
           <div className="w-8 h-8 rounded-full border border-pulse-accent bg-pulse-surface-raised relative">
             <motion.div
               className="absolute inset-0 rounded-full bg-pulse-accent/20"
               animate={{ scale: [1, 1.4, 1], opacity: [0.5, 0, 0.5] }}
               transition={{ duration: 2, repeat: Infinity }}
             />
           </div>
        </div>
      )
    }
  ];

  return (
    <section id="how-it-works" className="py-32 max-w-[1200px] mx-auto px-6">
      <div className="text-center mb-20">
        <h2 className="text-4xl md:text-5xl font-light text-pulse-text-primary mb-4">How Pulse works</h2>
        <p className="text-lg text-pulse-text-secondary">Connect once. Understand everything.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {features.map((feature, idx) => (
          <motion.div
            key={idx}
            className="p-8 rounded-2xl border border-pulse-border bg-pulse-surface hover:bg-pulse-surface-raised transition-all group"
            whileHover={{ y: -4, borderColor: "rgba(110,123,242,0.3)" }}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: idx * 0.1 }}
          >
            <div className="mb-6 p-4 rounded-xl bg-pulse-bg border border-pulse-border w-fit group-hover:text-pulse-primary transition-colors">
              <feature.icon className="w-6 h-6" />
            </div>
            
            <h3 className="text-xl font-light text-pulse-text-primary mb-3">{feature.title}</h3>
            <p className="text-sm text-pulse-text-secondary leading-relaxed mb-8">
              {feature.body}
            </p>

            <div className="rounded-xl border border-pulse-border bg-pulse-bg/50 backdrop-blur-sm overflow-hidden">
              {feature.visual}
            </div>
          </motion.div>
        ))}
      </div>
    </section>
  );
};

export default HowItWorks;
