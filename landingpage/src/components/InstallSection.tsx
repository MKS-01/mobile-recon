'use client';

import { motion } from 'motion/react';
import { useState } from 'react';
import { Copy, Check, Terminal, Download, Package } from 'lucide-react';
import TerminalWindow from './TerminalWindow';

const installMethods = [
  {
    id: 'quick',
    label: 'Quick',
    icon: Download,
    commands: [
      'git clone https://github.com/MKS-01/mobile-recon.git',
      'cd mobile-recon && ./scripts/install.sh',
    ],
  },
  {
    id: 'manual',
    label: 'Manual',
    icon: Terminal,
    commands: [
      'git clone https://github.com/MKS-01/mobile-recon.git',
      'cd mobile-recon/go-tools/mobile-recon-cli',
      'go build -o mobile-recon && mv mobile-recon ~/go/bin/',
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
      className="p-1.5 text-[var(--text-tertiary)] hover:text-[var(--accent)] transition-colors rounded"
    >
      {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
    </button>
  );
}

export default function InstallSection() {
  const [activeMethod, setActiveMethod] = useState('quick');
  const currentMethod = installMethods.find((m) => m.id === activeMethod)!;

  return (
    <section id="install" className="relative py-24">
      <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-12"
        >
          <h2 className="text-3xl sm:text-4xl font-bold text-[var(--text-primary)] mb-4">
            Get started
          </h2>
          <p className="text-[var(--text-secondary)]">
            Pick your favorite way to install. Takes less than a minute.
          </p>
        </motion.div>

        {/* Method tabs */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="flex justify-center gap-2 mb-8"
        >
          {installMethods.map((method) => (
            <button
              key={method.id}
              onClick={() => setActiveMethod(method.id)}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all ${
                activeMethod === method.id
                  ? 'bg-[var(--accent)] text-white'
                  : 'bg-[var(--bg-secondary)] text-[var(--text-secondary)] hover:text-[var(--text-primary)] border border-[var(--border)]'
              }`}
            >
              <method.icon className="w-4 h-4" />
              {method.label}
            </button>
          ))}
        </motion.div>

        {/* Terminal */}
        <motion.div
          key={activeMethod}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
        >
          <TerminalWindow title="install">
            <div className="space-y-2">
              {currentMethod.commands.map((cmd, i) => (
                <div key={i} className="group flex items-center justify-between gap-4">
                  <div className="flex items-center gap-2 font-mono text-sm overflow-x-auto">
                    <span className="text-[var(--accent)] shrink-0">$</span>
                    <code className="text-[var(--text-primary)]">{cmd}</code>
                  </div>
                  <div className="opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
                    <CopyButton text={cmd} />
                  </div>
                </div>
              ))}
            </div>
          </TerminalWindow>
        </motion.div>

        {/* Requirements */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="mt-12 flex flex-wrap justify-center gap-4"
        >
          {[
            { name: 'Go 1.21+', desc: 'Build' },
            { name: 'ADB', desc: 'Android' },
            { name: 'Nmap', desc: 'Network' },
          ].map((req, i) => (
            <div
              key={i}
              className="px-4 py-2 bg-[var(--bg-secondary)] rounded-lg border border-[var(--border)] text-center"
            >
              <span className="text-sm font-medium text-[var(--text-primary)]">{req.name}</span>
              <span className="text-[var(--text-tertiary)] text-xs ml-2">({req.desc})</span>
            </div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}
