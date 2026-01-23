'use client';

import { motion } from 'motion/react';

interface TerminalWindowProps {
  title?: string;
  children: React.ReactNode;
  className?: string;
}

export default function TerminalWindow({ title = 'terminal', children, className = '' }: TerminalWindowProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ duration: 0.5 }}
      className={`terminal ${className}`}
    >
      {/* Title bar */}
      <div className="terminal-header">
        <div className="flex gap-2">
          <div className="terminal-dot bg-[#ff5f57]" />
          <div className="terminal-dot bg-[#febc2e]" />
          <div className="terminal-dot bg-[#28c840]" />
        </div>
        <span className="ml-3 text-xs text-[var(--text-tertiary)]">{title}</span>
      </div>
      {/* Content */}
      <div className="terminal-body">
        {children}
      </div>
    </motion.div>
  );
}
