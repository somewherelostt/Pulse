"use client";

import React from "react";
import { Activity, Twitter, Github } from "lucide-react";

const Footer = () => {
  return (
    <footer className="py-20 border-t border-pulse-border bg-pulse-bg">
      <div className="max-w-[1200px] mx-auto px-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12 mb-20">
          <div className="col-span-1 md:col-span-1">
            <div className="flex items-center gap-2 mb-6">
              <Activity className="w-5 h-5 text-pulse-primary" />
              <span className="font-light text-xl">Pulse</span>
            </div>
            <p className="text-sm text-pulse-text-muted leading-relaxed">
              Behavioral intelligence for a resilient digital life.
            </p>
          </div>

          <div className="col-span-1 md:col-span-2 flex flex-wrap gap-10 md:gap-20 justify-start md:justify-center">
            <div className="space-y-4">
              <h4 className="text-xs font-mono text-pulse-text-primary uppercase tracking-widest">Platform</h4>
              <ul className="space-y-2 text-sm text-pulse-text-secondary">
                <li><a href="#" className="hover:text-pulse-primary transition-colors">How it works</a></li>
                <li><a href="#" className="hover:text-pulse-primary transition-colors">Why it matters</a></li>
                <li><a href="#" className="hover:text-pulse-primary transition-colors">Early access</a></li>
              </ul>
            </div>
            <div className="space-y-4">
              <h4 className="text-xs font-mono text-pulse-text-primary uppercase tracking-widest">Company</h4>
              <ul className="space-y-2 text-sm text-pulse-text-secondary">
                <li><a href="#" className="hover:text-pulse-primary transition-colors">Privacy</a></li>
                <li><a href="#" className="hover:text-pulse-primary transition-colors">GitHub</a></li>
                <li><a href="#" className="hover:text-pulse-primary transition-colors">Ethics</a></li>
              </ul>
            </div>
          </div>

          <div className="col-span-1 flex flex-col items-start md:items-end gap-6">
             <div className="flex gap-4">
                <a href="#" className="p-2 rounded-lg border border-pulse-border bg-pulse-surface hover:border-pulse-primary transition-colors">
                  <Twitter className="w-4 h-4 text-pulse-text-secondary" />
                </a>
                <a href="#" className="p-2 rounded-lg border border-pulse-border bg-pulse-surface hover:border-pulse-primary transition-colors">
                  <Github className="w-4 h-4 text-pulse-text-secondary" />
                </a>
             </div>
             <div className="text-xs text-pulse-text-muted font-mono">
               build by maaz
             </div>
          </div>
        </div>

        <div className="pt-10 border-t border-pulse-border text-center space-y-6">
          <p className="text-[10px] font-mono text-pulse-text-muted uppercase tracking-widest max-w-[600px] mx-auto leading-relaxed">
            Pulse is not a medical device. Not a replacement for professional mental health care. Always seek professional support if you're in crisis.
          </p>
          <div className="p-4 rounded-xl border border-pulse-danger/30 bg-pulse-danger/5 inline-block">
             <span className="text-xs text-pulse-danger font-medium tracking-wide uppercase">
               If you're in crisis: 988 Suicide & Crisis Lifeline
             </span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
