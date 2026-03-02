"use client";

import React, { useState, useEffect } from "react";
import { motion, useScroll, useTransform } from "framer-motion";
import { Activity } from "lucide-react";

const Navigation = () => {
  const { scrollY } = useScroll();
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const unsub = scrollY.on("change", (latest) => {
      setScrolled(latest > 50);
    });
    return () => unsub();
  }, [scrollY]);

  return (
    <motion.nav
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        scrolled
          ? "bg-pulse-bg/80 backdrop-blur-md h-16 border-b border-pulse-border"
          : "bg-transparent h-20"
      }`}
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
    >
      <div className="max-w-[1200px] mx-auto px-6 h-full flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Activity className="w-4 h-4 text-pulse-primary" />
          <span className="font-light tracking-tight text-xl text-pulse-text-primary">
            Pulse
          </span>
        </div>

        <div className="flex items-center gap-8">
          <div className="hidden md:flex items-center gap-8">
            <a
              href="#how-it-works"
              className="text-sm text-pulse-text-secondary hover:text-pulse-text-primary transition-colors"
            >
              How it works
            </a>
            <a
              href="#why-it-matters"
              className="text-sm text-pulse-text-secondary hover:text-pulse-text-primary transition-colors"
            >
              Why it matters
            </a>
          </div>
          <a
            href="/auth"
            className="text-sm px-4 py-2 rounded-full border border-pulse-border text-pulse-text-secondary hover:text-pulse-text-primary hover:border-pulse-text-primary transition-all"
          >
            Get early access
          </a>
        </div>
      </div>
    </motion.nav>
  );
};

export default Navigation;
