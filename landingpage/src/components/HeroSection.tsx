'use client';

import { motion } from 'motion/react';
import { ArrowRight, Sparkles } from 'lucide-react';
import MobileDevice from './MobileDevice';

export default function HeroSection() {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden pt-20">
      {/* Subtle background pattern */}
      <div className="absolute inset-0 dot-pattern opacity-50" />

      {/* Gradient orbs */}
      <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-[var(--accent)]/5 rounded-full blur-3xl" />
      <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-[var(--accent)]/5 rounded-full blur-3xl" />

      <div className="relative z-10 max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="grid lg:grid-cols-2 gap-16 items-center">
          {/* Left content */}
          <div className="text-center lg:text-left">
            {/* Badge */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5 }}
              className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-[var(--accent-muted)] border border-[var(--accent)]/20 mb-6"
            >
              <Sparkles className="w-3.5 h-3.5 text-[var(--accent)]" />
              <span className="text-xs font-medium text-[var(--accent)]">Weekend Project</span>
            </motion.div>

            {/* Main title */}
            <motion.h1
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.1 }}
              className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight mb-6"
            >
              <span className="text-[var(--text-primary)]">Mobile security</span>
              <br />
              <span className="text-[var(--accent)]">made simple</span>
            </motion.h1>

            {/* Subtitle */}
            <motion.p
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.2 }}
              className="text-lg text-[var(--text-secondary)] mb-8 max-w-lg mx-auto lg:mx-0 leading-relaxed"
            >
              A fun toolkit for mobile security testing on Android & iOS.
              Built for security researchers and curious developers who like to tinker.
            </motion.p>

            {/* CTA Buttons */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.3 }}
              className="flex flex-col sm:flex-row gap-3 justify-center lg:justify-start"
            >
              <a
                href="#install"
                className="group inline-flex items-center justify-center gap-2 px-6 py-3 bg-[var(--accent)] text-white font-medium rounded-xl hover:bg-[var(--accent-dark)] transition-all shadow-lg shadow-[var(--accent)]/20"
              >
                Get Started
                <ArrowRight className="w-4 h-4 group-hover:translate-x-0.5 transition-transform" />
              </a>

              <a
                href="/docs"
                className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-[var(--bg-secondary)] text-[var(--text-primary)] font-medium rounded-xl hover:bg-[var(--bg-tertiary)] transition-all border border-[var(--border)]"
              >
                Documentation
              </a>
            </motion.div>

            {/* Quick stats */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.5, delay: 0.4 }}
              className="flex items-center gap-8 mt-12 justify-center lg:justify-start"
            >
              {[
                { value: '4', label: 'Tools' },
                { value: '50+', label: 'Commands' },
                { value: 'MIT', label: 'License' },
              ].map((stat, i) => (
                <div key={i} className="text-center lg:text-left">
                  <div className="text-2xl font-bold text-[var(--text-primary)]">{stat.value}</div>
                  <div className="text-sm text-[var(--text-tertiary)]">{stat.label}</div>
                </div>
              ))}
            </motion.div>
          </div>

          {/* Right content - Mobile devices */}
          <motion.div
            initial={{ opacity: 0, x: 50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.7, delay: 0.3 }}
            className="hidden lg:flex items-center justify-center relative"
          >
            <div className="relative flex items-end gap-6">
              {/* Android device - slightly back */}
              <div className="transform -translate-y-8">
                <MobileDevice type="android" />
              </div>

              {/* iOS device - slightly forward */}
              <div className="transform translate-y-8 -ml-8">
                <MobileDevice type="ios" />
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
