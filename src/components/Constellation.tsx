"use client";

import React, { useMemo } from "react";
import { motion } from "framer-motion";
import { Shield, User, Heart } from "lucide-react";

const Constellation = () => {
  const nodes = useMemo(() => [
    { id: 1, x: 20, y: 30, initial: "A", label: "Recovered" },
    { id: 2, x: 40, y: 70, initial: "M", label: "Active" },
    { id: 3, x: 70, y: 20, initial: "S", label: "Similar pattern" },
    { id: 4, x: 85, y: 55, initial: "R", label: "Recovered" },
    { id: 5, x: 55, y: 40, initial: "YOU", isUser: true },
    { id: 6, x: 30, y: 15, initial: "K", label: "Recovered" },
    { id: 7, x: 15, y: 60, initial: "J", label: "Active" },
    { id: 8, x: 80, y: 80, initial: "T", label: "Recovered" },
  ], []);

  const connections = [
    [5, 1], [5, 3], [5, 4], [1, 6], [7, 2], [4, 8]
  ];

  return (
    <section className="py-32 bg-pulse-bg relative pulse-dot-grid">
      <div className="max-w-[1200px] mx-auto px-6 text-center mb-20">
        <h2 className="text-4xl md:text-5xl font-light text-pulse-text-primary mb-4 leading-tight">
          The hardest part isn't the pattern. <br />
          <span className="text-pulse-text-secondary">It's feeling like nobody else has been through it.</span>
        </h2>
      </div>

      <div className="max-w-[1000px] mx-auto h-[500px] relative mb-24 overflow-hidden border border-pulse-border rounded-3xl bg-pulse-surface/30 backdrop-blur-sm">
         <svg className="absolute inset-0 w-full h-full">
            {connections.map(([fromId, toId], idx) => {
              const from = nodes.find(n => n.id === fromId)!;
              const to = nodes.find(n => n.id === toId)!;
              return (
                <motion.line
                  key={idx}
                  x1={`${from.x}%`}
                  y1={`${from.y}%`}
                  x2={`${to.x}%`}
                  y2={`${to.y}%`}
                  stroke="rgba(110,123,242,0.2)"
                  strokeWidth="1"
                  initial={{ pathLength: 0 }}
                  whileInView={{ pathLength: 1 }}
                  transition={{ duration: 2, delay: idx * 0.2 }}
                />
              );
            })}
         </svg>

         {nodes.map((node) => (
           <motion.div
             key={node.id}
             className={`absolute -translate-x-1/2 -translate-y-1/2 flex flex-col items-center group`}
             style={{ left: `${node.x}%`, top: `${node.y}%` }}
             initial={{ opacity: 0, scale: 0 }}
             whileInView={{ opacity: 1, scale: 1 }}
             transition={{ type: "spring", stiffness: 100, delay: node.id * 0.1 }}
           >
             <motion.div
               className={`w-10 h-10 rounded-full flex items-center justify-center text-xs font-mono border transition-all ${
                 node.isUser 
                   ? "bg-pulse-primary text-white border-pulse-primary shadow-[0_0_20px_rgba(110,123,242,0.5)]" 
                   : "bg-pulse-surface-raised text-pulse-text-secondary border-pulse-border group-hover:border-pulse-primary"
               }`}
               animate={node.isUser ? { scale: [1, 1.1, 1] } : { y: [0, -4, 0] }}
               transition={{ duration: node.isUser ? 2 : 3 + Math.random() * 2, repeat: Infinity, ease: "easeInOut" }}
             >
               {node.initial}
             </motion.div>
             {node.label && (
                <div className="mt-2 opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap bg-pulse-bg border border-pulse-border rounded-full px-2 py-0.5 text-[10px] text-pulse-text-muted uppercase tracking-widest pointer-events-none">
                  {node.label}
                </div>
             )}
           </motion.div>
         ))}

         {/* Annotations */}
         <motion.div
            className="absolute left-10 bottom-10 p-4 rounded-xl border border-pulse-accent/20 bg-pulse-surface-raised flex items-center gap-3"
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            transition={{ delay: 1.5 }}
         >
            <div className="w-2 h-2 rounded-full bg-pulse-accent animate-pulse" />
            <div className="text-[10px] font-mono text-pulse-accent uppercase tracking-widest">
               Similar pattern · 3 weeks ago · Recovered
            </div>
         </motion.div>
      </div>

      <div className="max-w-[1200px] mx-auto px-6 grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
        {[
          { icon: Shield, title: "Anonymous by default", body: "No name, no profile, no history shared." },
          { icon: User, title: "Pattern-matched", body: "Behavioral fingerprint, not diagnosis." },
          { icon: Heart, title: "Human, not AI", body: "Real people, encrypted sessions." },
        ].map((principle, idx) => (
          <motion.div
            key={idx}
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ delay: idx * 0.1 }}
          >
            <principle.icon className="w-6 h-6 text-pulse-primary mx-auto mb-4" />
            <h3 className="text-lg font-light text-pulse-text-primary mb-2">{principle.title}</h3>
            <p className="text-sm text-pulse-text-secondary">{principle.body}</p>
          </motion.div>
        ))}
      </div>

      <div className="mt-20 text-center">
        <button className="px-8 py-3 rounded-full border border-pulse-accent text-pulse-accent hover:bg-pulse-accent hover:text-pulse-bg transition-all text-sm font-medium">
          Join the peer pool
        </button>
      </div>
    </section>
  );
};

export default Constellation;
