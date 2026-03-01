"use client";

import React from "react";
import { motion } from "framer-motion";
import { Activity, Moon, Calendar, Globe, LayoutDashboard, Share2, Compass, LogOut, TrendingDown, ArrowUpRight, ArrowDownRight, MessageSquare, Sparkles } from "lucide-react";

const DashboardPreview = () => {
  const metrics = [
    { label: "Sleep Consistency", value: "67%", trend: "↓ from 84%", status: "warning", icon: Moon },
    { label: "Calendar Load", value: "71%", trend: "↑ critical", status: "danger", icon: Calendar },
    { label: "Digital Entropy", value: "6.2", trend: "↑ compulsive signal", status: "warning", icon: Globe },
    { label: "Mood Trend", value: "Declining", trend: "↓ 12% week over week", status: "danger", icon: Activity },
  ];

  return (
    <section className="py-32 bg-pulse-bg relative overflow-hidden">
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[80%] h-px bg-gradient-to-right from-transparent via-pulse-primary/20 to-transparent" />
      
      <div className="max-w-[1200px] mx-auto px-6 mb-20 text-center">
        <h2 className="text-4xl md:text-5xl font-light mb-4">
          One dashboard. <br />
          Four data sources. <br />
          <span className="pulse-gradient font-semibold">Finally connected.</span>
        </h2>
      </div>

      <div className="max-w-[1200px] mx-auto px-6">
        <div className="rounded-[24px] border border-pulse-border bg-pulse-surface overflow-hidden shadow-2xl flex min-h-[600px] group relative">
          {/* Subtle inner glow */}
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_0%,rgba(110,123,242,0.05),transparent_50%)]" />

          {/* Sidebar */}
          <div className="hidden md:flex w-[200px] border-r border-pulse-border flex-col p-6 z-10">
            <div className="flex items-center gap-2 mb-12">
              <Activity className="w-5 h-5 text-pulse-primary" />
              <span className="font-light text-lg">Pulse</span>
            </div>

            <div className="space-y-6 flex-1">
              {[
                { icon: LayoutDashboard, label: "Overview", active: true },
                { icon: Moon, label: "Sleep" },
                { icon: Calendar, label: "Calendar" },
                { icon: Globe, label: "Browser" },
                { icon: Compass, label: "Constellation" },
              ].map((item, idx) => (
                <div key={idx} className={`flex items-center gap-3 text-sm ${item.active ? "text-pulse-primary font-medium" : "text-pulse-text-muted hover:text-pulse-text-secondary"} cursor-pointer transition-colors`}>
                  <item.icon className="w-4 h-4" />
                  {item.label}
                </div>
              ))}
            </div>

            <div className="flex items-center gap-3 text-sm text-pulse-text-muted mt-auto pt-6 border-t border-pulse-border">
              <div className="w-6 h-6 rounded-full bg-pulse-primary/20 flex items-center justify-center text-[10px] text-pulse-primary border border-pulse-primary/30">
                A
              </div>
              <span>Anonymous user</span>
            </div>
          </div>

          {/* Main Content */}
          <div className="flex-1 p-8 z-10 flex flex-col gap-8">
            {/* Top Row: Metrics */}
            <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
              {metrics.map((metric, idx) => (
                <div key={idx} className="p-4 rounded-xl border border-pulse-border bg-pulse-bg relative overflow-hidden group/card hover:border-pulse-primary/30 transition-colors">
                  <div className="flex items-center justify-between mb-4">
                    <metric.icon className="w-4 h-4 text-pulse-text-muted" />
                    <div className={`w-1.5 h-1.5 rounded-full ${metric.status === "warning" ? "bg-pulse-accent-warm" : "bg-pulse-danger"}`} />
                  </div>
                  <div className="text-2xl font-mono text-pulse-text-primary mb-1">{metric.value}</div>
                  <div className="text-[10px] font-mono text-pulse-text-muted uppercase tracking-wider mb-2">{metric.label}</div>
                  <div className={`text-[10px] font-mono flex items-center gap-1 ${metric.status === "warning" ? "text-pulse-accent-warm" : "text-pulse-danger"}`}>
                    {metric.trend.includes("↑") ? <ArrowUpRight className="w-3 h-3" /> : <ArrowDownRight className="w-3 h-3" />}
                    {metric.trend}
                  </div>
                </div>
              ))}
            </div>

            {/* Middle: Unified Timeline */}
            <div className="flex-1 p-6 rounded-2xl border border-pulse-border bg-pulse-bg flex flex-col">
              <div className="flex items-center justify-between mb-8">
                <h4 className="text-sm font-mono text-pulse-text-secondary uppercase tracking-widest">30-Day behavioral fingerprint</h4>
                <div className="text-[10px] font-mono text-pulse-text-muted uppercase">February 2026</div>
              </div>

              <div className="flex-1 relative min-h-[160px]">
                {/* Y-axis labels */}
                <div className="absolute left-0 top-0 bottom-0 w-8 flex flex-col justify-between text-[8px] font-mono text-pulse-text-muted py-2 border-r border-pulse-border/50">
                   <span>100</span>
                   <span>75</span>
                   <span>50</span>
                   <span>25</span>
                   <span>0</span>
                </div>

                <div className="ml-12 h-full flex items-end gap-[2px]">
                   {Array.from({ length: 60 }).map((_, i) => (
                      <motion.div
                        key={i}
                        className="flex-1 bg-pulse-primary/10 rounded-t-[1px]"
                        initial={{ height: 0 }}
                        whileInView={{ height: `${20 + Math.sin(i * 0.2) * 40 + 30}%` }}
                        transition={{ duration: 1, delay: i * 0.01 }}
                      />
                   ))}

                   {/* Drift detected annotation */}
                   <motion.div
                      className="absolute left-[80%] top-0 bottom-0 w-px bg-pulse-accent-warm/50 border-l border-dashed border-pulse-accent-warm z-20"
                      initial={{ opacity: 0 }}
                      whileInView={{ opacity: 1 }}
                      transition={{ delay: 1.5 }}
                   >
                     <div className="absolute top-4 left-2 whitespace-nowrap bg-pulse-bg border border-pulse-accent-warm/30 rounded px-2 py-1 text-[8px] font-mono text-pulse-accent-warm">
                       DRIFT DETECTED
                     </div>
                   </motion.div>
                </div>
              </div>
            </div>

            {/* Bottom Row: AI Insights & Peer Match */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="p-6 rounded-2xl border border-pulse-border bg-pulse-bg relative overflow-hidden flex gap-4">
                <div className="absolute left-0 top-0 bottom-0 w-1 bg-gradient-to-bottom from-pulse-primary to-pulse-accent" />
                <div className="p-2 rounded-lg bg-pulse-primary/10 border border-pulse-primary/20 h-fit">
                   <Sparkles className="w-4 h-4 text-pulse-primary" />
                </div>
                <div>
                   <h5 className="text-sm font-light text-pulse-text-primary mb-2">Pattern detected</h5>
                   <p className="text-xs text-pulse-text-secondary leading-relaxed">
                     Your mood typically declines 2-3 days after weeks where meetings exceed 65% of work hours. Last week: 71%.
                   </p>
                </div>
              </div>

              <div className="p-6 rounded-2xl border border-pulse-primary/30 bg-pulse-primary/5 relative overflow-hidden flex gap-4">
                <div className="p-2 rounded-lg bg-pulse-accent/10 border border-pulse-accent/20 h-fit">
                   <div className="relative">
                      <MessageSquare className="w-4 h-4 text-pulse-accent" />
                      <motion.div
                        className="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-pulse-accent"
                        animate={{ scale: [1, 1.5, 1], opacity: [1, 0.4, 1] }}
                        transition={{ duration: 2, repeat: Infinity }}
                      />
                   </div>
                </div>
                <div className="flex-1">
                   <h5 className="text-sm font-light text-pulse-text-primary mb-1">Peer match available</h5>
                   <p className="text-xs text-pulse-text-secondary mb-4">
                     Someone navigated this pattern 3 weeks ago. They're available to talk.
                   </p>
                   <button className="text-[10px] font-mono uppercase tracking-widest text-pulse-accent hover:text-white transition-colors flex items-center gap-1 group">
                      Talk to them 
                      <ArrowRight className="w-3 h-3 group-hover:translate-x-1 transition-transform" />
                   </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

const ArrowRight = ({ className }: { className?: string }) => (
  <svg className={className} width="16" height="16" viewBox="0 0 16 16" fill="none">
    <path d="M1 8H15M15 8L8 1M15 8L8 15" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
  </svg>
);

export default DashboardPreview;
