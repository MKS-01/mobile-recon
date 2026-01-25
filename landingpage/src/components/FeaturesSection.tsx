'use client';

import { motion } from 'motion/react';
import {
  Smartphone,
  Search,
  Network,
  Apple,
  Shield,
  Bug,
  Key,
  Wifi,
  Lock,
  Eye,
} from 'lucide-react';

const features = [
  {
    icon: Smartphone,
    title: 'ADB Toolkit',
    description: 'Talk to Android devices like a pro. List apps, pull APKs, monitor logs, and more.',
    color: 'bg-emerald-500',
  },
  {
    icon: Search,
    title: 'APK Analyzer',
    description: 'Peek inside APK files. Find hardcoded secrets, check permissions, spot misconfigs.',
    color: 'bg-blue-500',
  },
  {
    icon: Apple,
    title: 'iOS Toolkit',
    description: 'Play with iOS simulators. No jailbreak needed. Frida integration included.',
    color: 'bg-gray-500',
  },
  {
    icon: Network,
    title: 'Nmap Toolkit',
    description: 'Scan networks for mobile devices. Find ADB ports, detect proxies, enumerate services.',
    color: 'bg-purple-500',
  },
];

const capabilities = [
  { icon: Shield, label: 'Security Testing' },
  { icon: Bug, label: 'Vuln Detection' },
  { icon: Key, label: 'Secret Discovery' },
  { icon: Eye, label: 'Runtime Analysis' },
  { icon: Wifi, label: 'Network Recon' },
  { icon: Lock, label: 'SSL Testing' },
];

export default function FeaturesSection() {
  return (
    <section id="features" className="relative py-24">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-16"
        >
          <h2 className="text-3xl sm:text-4xl font-bold text-[var(--text-primary)] mb-4">
            What&apos;s in the box?
          </h2>
          <p className="text-[var(--text-secondary)] max-w-2xl mx-auto">
            Four tools that make mobile security testing less painful. Each one does one thing well.
          </p>
        </motion.div>

        {/* Feature cards */}
        <div className="grid md:grid-cols-2 gap-6 mb-20">
          {features.map((feature, i) => (
            <motion.div
              key={feature.title}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.1 }}
              className="card p-6 group"
            >
              <div className={`w-12 h-12 ${feature.color} rounded-xl flex items-center justify-center mb-4 group-hover:scale-110 transition-transform`}>
                <feature.icon className="w-6 h-6 text-white" />
              </div>
              <h3 className="text-xl font-semibold text-[var(--text-primary)] mb-2">
                {feature.title}
              </h3>
              <p className="text-[var(--text-secondary)] leading-relaxed">
                {feature.description}
              </p>
            </motion.div>
          ))}
        </div>

        {/* Capabilities */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-center mb-8"
        >
          <h3 className="text-xl font-semibold text-[var(--text-primary)] mb-8">
            Things you can do
          </h3>
        </motion.div>

        <div className="flex flex-wrap justify-center gap-3">
          {capabilities.map((cap, i) => (
            <motion.div
              key={cap.label}
              initial={{ opacity: 0, scale: 0.9 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              className="flex items-center gap-2 px-4 py-2.5 bg-[var(--bg-secondary)] rounded-full border border-[var(--border)] hover:border-[var(--accent)]/30 transition-colors"
            >
              <cap.icon className="w-4 h-4 text-[var(--accent)]" />
              <span className="text-sm text-[var(--text-secondary)]">{cap.label}</span>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}
