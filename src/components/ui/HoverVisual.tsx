"use client";

import React, { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

interface HoverVisualProps {
  text: string;
  children: React.ReactNode;
  className?: string;
}

export const HoverVisual = ({ text, children, className }: HoverVisualProps) => {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <span
      className={`relative inline-block ${className ?? ""}`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Text — smooth glow on hover, no font change */}
      <span
        className="relative inline-block cursor-default transition-colors duration-300"
        style={{
          color: isHovered ? "#6E7BF2" : "inherit",
          textShadow: isHovered ? "0 0 40px rgba(110,123,242,0.4)" : "none",
        }}
      >
        {text}
        {/* Slide-in underline */}
        <motion.span
          aria-hidden
          className="absolute bottom-0 left-0 h-[2px] rounded-full bg-gradient-to-r from-pulse-primary to-pulse-accent"
          initial={{ scaleX: 0, opacity: 0 }}
          animate={{ scaleX: isHovered ? 1 : 0, opacity: isHovered ? 1 : 0 }}
          transition={{ duration: 0.25, ease: "easeOut" }}
          style={{ transformOrigin: "left center", display: "block", width: "100%" }}
        />
      </span>

      {/* Popup card — appears above the word */}
      <AnimatePresence>
        {isHovered && (
          <motion.div
            initial={{ opacity: 0, y: 8, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 6, scale: 0.96 }}
            transition={{ duration: 0.18, ease: [0.16, 1, 0.3, 1] }}
            className="absolute left-1/2 -translate-x-1/2 bottom-full mb-5 z-50 pointer-events-none"
          >
            {/* Glow halo */}
            <div className="absolute inset-0 rounded-2xl bg-pulse-primary/8 blur-2xl scale-110 pointer-events-none" />

            <div className="relative bg-[#111118]/95 border border-[#2A2A3A]/80 rounded-2xl shadow-[0_24px_60px_rgba(0,0,0,0.6),0_0_0_1px_rgba(110,123,242,0.1)] backdrop-blur-2xl min-w-[240px] overflow-hidden">
              {/* Top accent line */}
              <div className="h-px bg-gradient-to-r from-transparent via-pulse-primary/50 to-transparent" />
              {children}
            </div>

            {/* Arrow pointer */}
            <div className="absolute -bottom-[5px] left-1/2 -translate-x-1/2 w-2.5 h-2.5 bg-[#111118] border-r border-b border-[#2A2A3A]/80 rotate-45 shadow-[2px_2px_4px_rgba(0,0,0,0.3)]" />
          </motion.div>
        )}
      </AnimatePresence>
    </span>
  );
};
