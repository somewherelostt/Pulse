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
      className={`relative inline-block cursor-help group ${className}`}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <span className="relative z-10 transition-all duration-300 group-hover:font-serif group-hover:italic group-hover:text-pulse-primary">
        {text}
      </span>
      <AnimatePresence>
        {isHovered && (
          <motion.div
            initial={{ opacity: 0, scale: 0.5, y: 20, rotate: -10 }}
            animate={{ opacity: 1, scale: 1, y: -100, rotate: 0 }}
            exit={{ opacity: 0, scale: 0.5, y: 20, rotate: 10 }}
            transition={{ type: "spring", stiffness: 400, damping: 25 }}
            className="absolute left-1/2 -translate-x-1/2 z-50 pointer-events-none"
          >
            <div className="bg-pulse-surface-raised/80 border border-pulse-border/50 p-1.5 rounded-2xl shadow-[0_20px_50px_rgba(0,0,0,0.5),0_0_20px_rgba(110,123,242,0.15)] backdrop-blur-xl min-w-[240px] overflow-hidden">
              <div className="relative rounded-xl overflow-hidden bg-pulse-bg/50">
                {children}
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </span>
  );
};
