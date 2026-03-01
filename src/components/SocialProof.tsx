"use client";

import React from "react";
import { motion } from "framer-motion";

const SocialProof = () => {
  const items = [
    "Privacy first",
    "Local processing",
    "No data sold",
    "Anonymous by default",
    "Open source",
    "WCAG compliant",
    "End-to-end encrypted sessions",
  ];

  const marqueeItems = [...items, ...items, ...items];

  return (
    <section className="h-20 border-y border-pulse-border bg-pulse-bg overflow-hidden flex items-center">
      <div className="flex w-full overflow-hidden">
        <motion.div
          className="flex whitespace-nowrap gap-12 items-center"
          animate={{ x: ["0%", "-33.33%"] }}
          transition={{
            duration: 30,
            repeat: Infinity,
            ease: "linear",
          }}
        >
          {marqueeItems.map((item, idx) => (
            <div
              key={idx}
              className="flex items-center gap-4 text-[11px] font-mono text-pulse-text-muted uppercase tracking-[0.2em]"
            >
              <span className="text-pulse-primary">●</span>
              {item}
            </div>
          ))}
        </motion.div>
      </div>
    </section>
  );
};

export default SocialProof;
