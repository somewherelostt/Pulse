"use client";

import React from "react";
import { motion } from "framer-motion";
import { Moon, Calendar, Globe, ArrowDown, TrendingDown } from "lucide-react";
import { HoverVisual } from "./ui/HoverVisual";

const ProblemSection = () => {
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { staggerChildren: 0.2 },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, x: 20 },
    visible: {
      opacity: 1,
      x: 0,
      transition: { duration: 0.6, ease: "easeOut" },
    },
  };

  return (
    <section id="why-it-matters" className="py-32 max-w-[1200px] mx-auto px-6">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-20 items-center">
        <div>
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              className="mb-8"
            >
              <HoverVisual text="12" className="text-[120px] font-mono font-light text-pulse-primary leading-none">
                <div className="p-4 flex flex-col gap-3">
                  <div className="flex items-center gap-2">
                    <TrendingDown className="w-4 h-4 text-pulse-danger" />
                    <span className="text-[10px] font-mono text-pulse-text-secondary uppercase">Treatment Gap</span>
                  </div>
                  <div className="h-24 w-full flex items-end gap-1">
                    {[30, 40, 35, 60, 50, 80, 70, 95].map((h, i) => (
                      <div key={i} className="flex-1 bg-pulse-danger/20 rounded-t-[2px] relative group/bar" style={{ height: `${h}%` }}>
                         {i === 7 && <div className="absolute -top-4 left-1/2 -translate-x-1/2 text-[8px] text-pulse-danger font-mono font-bold">CRITICAL</div>}
                      </div>
                    ))}
                  </div>
                  <p className="text-[10px] text-pulse-text-muted mt-2">Behavioral signals often diverge years before clinical diagnosis.</p>
                </div>
              </HoverVisual>
              <div className="text-2xl font-mono text-pulse-text-secondary uppercase tracking-widest mt-2">Years</div>
            </motion.div>

          
          <h2 className="text-4xl md:text-5xl font-light leading-tight text-pulse-text-primary mb-6">
            Average delay between first symptoms and first treatment
          </h2>
          <p className="text-xs font-mono text-pulse-text-muted uppercase tracking-widest">
            Source: WHO, 2023
          </p>
        </div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          className="space-y-4 relative"
        >
          {/* Silo Cards */}
          {[
            { icon: Moon, text: "Your sleep data lives here", color: "text-emerald-400" },
            { icon: Calendar, text: "Your stress lives here", color: "text-blue-400" },
            { icon: Globe, text: "Your spirals live here", color: "text-purple-400" },
          ].map((silo, idx) => (
            <motion.div
              key={idx}
              variants={itemVariants}
              className="p-6 rounded-2xl border border-pulse-border bg-pulse-surface flex items-center gap-6 group hover:border-pulse-text-secondary transition-colors"
            >
              <div className={`p-3 rounded-xl bg-pulse-bg border border-pulse-border ${silo.color} group-hover:scale-110 transition-transform`}>
                <silo.icon className="w-6 h-6" />
              </div>
              <span className="text-lg text-pulse-text-secondary font-light">{silo.text}</span>
            </motion.div>
          ))}

          {/* Connection text */}
          <motion.div 
            className="flex flex-col items-center gap-2 py-4"
            variants={itemVariants}
          >
            <ArrowDown className="w-4 h-4 text-pulse-text-muted animate-bounce" />
            <span className="text-xs font-mono text-pulse-text-muted uppercase tracking-widest">No connection</span>
          </motion.div>

          {/* Solution Card */}
          <motion.div
            variants={itemVariants}
            className="p-8 rounded-2xl border border-pulse-primary/30 bg-pulse-surface-raised shadow-[0_0_30px_rgba(110,123,242,0.1)] relative overflow-hidden"
          >
            <div className="absolute top-0 left-0 w-1 h-full bg-pulse-primary" />
            <h3 className="text-2xl font-light text-pulse-text-primary mb-2">
              Pulse connects them
            </h3>
            <p className="text-pulse-text-secondary">
              And shows you the pattern before you crash.
            </p>
            
            <motion.div
              className="absolute -right-4 -bottom-4 w-24 h-24 bg-pulse-primary/10 rounded-full blur-2xl"
              animate={{ scale: [1, 1.2, 1], opacity: [0.3, 0.6, 0.3] }}
              transition={{ duration: 4, repeat: Infinity }}
            />
          </motion.div>
        </motion.div>
      </div>
    </section>
  );
};

export default ProblemSection;
