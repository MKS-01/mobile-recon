'use client';

import { motion } from 'motion/react';
import { Copy, Check } from 'lucide-react';
import { useState } from 'react';
import TerminalWindow from './TerminalWindow';

const usageExamples = [
  {
    title: 'Quick Start',
    commands: [
      { cmd: 'mobile-recon list', out: 'adb, apk, ios, nmap' },
      { cmd: 'mobile-recon adb device list', out: 'emulator-5554  device' },
    ],
  },
  {
    title: 'Android',
    commands: [
      { cmd: 'mobile-recon adb device list', out: 'emulator-5554  device' },
      { cmd: 'mobile-recon adb app pull com.target.app', out: 'Saved: com.target.app.apk' },
    ],
  },
  {
    title: 'Analysis',
    commands: [
      { cmd: 'mobile-recon apk strings app.apk', out: 'Found 3 potential secrets...' },
      { cmd: 'mobile-recon apk permissions app.apk', out: '12 dangerous permissions' },
    ],
  },
  {
    title: 'Network',
    commands: [
      { cmd: 'mobile-recon nmap scan quick 192.168.1.0/24', out: 'Discovered 8 hosts' },
      { cmd: 'mobile-recon nmap mobile adb 192.168.1.0/24', out: 'Found 2 ADB devices' },
    ],
  },
];

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <button
      onClick={handleCopy}
      className="p-1 text-[var(--text-tertiary)] hover:text-[var(--accent)] transition-colors"
    >
      {copied ? <Check className="w-3.5 h-3.5" /> : <Copy className="w-3.5 h-3.5" />}
    </button>
  );
}

export default function UsageSection() {
  return (
    <section id="usage" className="relative py-24 bg-[var(--bg-secondary)]">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-3xl sm:text-4xl font-bold text-[var(--text-primary)] mb-4">
            How it works
          </h2>
          <p className="text-[var(--text-secondary)] max-w-2xl mx-auto">
            Simple commands, powerful results. Here are some things you can do.
          </p>
        </motion.div>

        {/* Terminal examples */}
        <div className="grid md:grid-cols-2 gap-6">
          {usageExamples.map((example, i) => (
            <TerminalWindow key={example.title} title={example.title.toLowerCase()}>
              <div className="space-y-3">
                {example.commands.map((item, j) => (
                  <div key={j} className="group">
                    <div className="flex items-center justify-between gap-2">
                      <div className="flex items-center gap-2 font-mono text-sm">
                        <span className="text-[var(--accent)]">$</span>
                        <span className="text-[var(--text-primary)]">{item.cmd}</span>
                      </div>
                      <div className="opacity-0 group-hover:opacity-100 transition-opacity">
                        <CopyButton text={item.cmd} />
                      </div>
                    </div>
                    <div className="text-[var(--text-tertiary)] text-sm pl-4 mt-1">
                      {item.out}
                    </div>
                  </div>
                ))}
              </div>
            </TerminalWindow>
          ))}
        </div>
      </div>
    </section>
  );
}
