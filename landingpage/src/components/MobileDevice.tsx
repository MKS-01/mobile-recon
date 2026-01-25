'use client';

import { motion } from 'motion/react';

interface MobileDeviceProps {
  type: 'android' | 'ios';
  className?: string;
}

export default function MobileDevice({ type, className = '' }: MobileDeviceProps) {
  const isAndroid = type === 'android';

  // Samsung S25 Ultra style for Android
  if (isAndroid) {
    return (
      <motion.div
        className={`relative ${className}`}
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.6 }}
      >
        {/* S25 Ultra Frame - Sharp corners, thin bezels */}
        <div className="relative w-[220px] h-[460px] rounded-[28px] bg-gradient-to-b from-[#2a2a2f] to-[#1f1f24] p-[3px] shadow-[0_25px_50px_-12px_rgba(0,0,0,0.5)]">
          {/* Titanium frame effect */}
          <div className="absolute inset-0 rounded-[28px] bg-gradient-to-r from-[#4a4a50] via-[#3a3a40] to-[#4a4a50] opacity-50" />

          {/* Inner bezel */}
          <div className="relative w-full h-full rounded-[25px] bg-[#0d0d0f] overflow-hidden">
            {/* Screen */}
            <div className="absolute inset-[2px] rounded-[23px] bg-gradient-to-br from-[#1a1a1f] to-[#0f0f12] overflow-hidden">
              {/* Status bar */}
              <div className="flex items-center justify-between px-5 pt-3 pb-2">
                <span className="text-[10px] text-[var(--text-secondary)] font-medium">9:41</span>
                {/* Camera punch hole - centered */}
                <div className="w-3 h-3 rounded-full bg-[#0a0a0c] border border-[#2a2a30]" />
                <div className="flex items-center gap-1">
                  <div className="flex gap-[2px]">
                    {[1, 2, 3, 4].map((i) => (
                      <div key={i} className="w-[3px] h-2 rounded-sm bg-[var(--text-tertiary)]" />
                    ))}
                  </div>
                  <span className="text-[10px] text-[var(--text-secondary)]">5G</span>
                </div>
              </div>

              {/* Screen content */}
              <div className="px-4 pt-6">
                {/* App header simulation */}
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-[var(--accent)] flex items-center justify-center">
                    <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                  </div>
                  <div>
                    <div className="text-sm font-semibold text-[var(--text-primary)]">Security Scan</div>
                    <div className="text-xs text-[var(--text-tertiary)]">Analyzing...</div>
                  </div>
                </div>

                {/* Scan progress */}
                <div className="space-y-3">
                  {['Permissions', 'Network', 'Storage', 'Crypto'].map((item, i) => (
                    <motion.div
                      key={item}
                      initial={{ opacity: 0, x: -10 }}
                      whileInView={{ opacity: 1, x: 0 }}
                      transition={{ delay: i * 0.1 }}
                      className="flex items-center justify-between p-3 rounded-xl bg-[var(--bg-secondary)]"
                    >
                      <span className="text-xs text-[var(--text-secondary)]">{item}</span>
                      <div className="w-4 h-4 rounded-full bg-[var(--accent)]/20 flex items-center justify-center">
                        <svg className="w-2.5 h-2.5 text-[var(--accent)]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                        </svg>
                      </div>
                    </motion.div>
                  ))}
                </div>
              </div>

              {/* Navigation bar */}
              <div className="absolute bottom-2 left-1/2 -translate-x-1/2">
                <div className="w-28 h-1 rounded-full bg-[var(--text-tertiary)]/50" />
              </div>
            </div>
          </div>

          {/* Side buttons */}
          <div className="absolute right-[-2px] top-24 w-[3px] h-16 rounded-r-sm bg-[#3a3a40]" />
          <div className="absolute left-[-2px] top-28 w-[3px] h-8 rounded-l-sm bg-[#3a3a40]" />
          <div className="absolute left-[-2px] top-40 w-[3px] h-8 rounded-l-sm bg-[#3a3a40]" />
        </div>

        {/* Label */}
        <div className="text-center mt-4">
          <span className="text-sm font-medium text-[var(--text-secondary)]">Android</span>
        </div>
      </motion.div>
    );
  }

  // iPhone 17 style - Dynamic Island, thinner bezels
  return (
    <motion.div
      className={`relative ${className}`}
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ duration: 0.6 }}
    >
      {/* iPhone 17 Frame - Titanium finish, rounded corners */}
      <div className="relative w-[210px] h-[450px] rounded-[44px] bg-gradient-to-b from-[#e8e8ed] to-[#d1d1d6] p-[3px] shadow-[0_25px_50px_-12px_rgba(0,0,0,0.4)]">
        {/* Titanium frame effect */}
        <div className="absolute inset-0 rounded-[44px] bg-gradient-to-r from-[#f5f5f7] via-[#e8e8ed] to-[#f5f5f7] opacity-60" />

        {/* Inner bezel */}
        <div className="relative w-full h-full rounded-[41px] bg-[#0d0d0f] overflow-hidden">
          {/* Screen */}
          <div className="absolute inset-[2px] rounded-[39px] bg-gradient-to-br from-[#1a1a1f] to-[#0f0f12] overflow-hidden">
            {/* Dynamic Island */}
            <div className="flex justify-center pt-3">
              <motion.div
                className="w-[90px] h-[28px] rounded-full bg-[#0a0a0c] flex items-center justify-center gap-2"
                animate={{ width: [90, 95, 90] }}
                transition={{ duration: 2, repeat: Infinity }}
              >
                <div className="w-2 h-2 rounded-full bg-[#1a1a1f]" />
                <div className="w-2.5 h-2.5 rounded-full bg-[#1a1a1f] border border-[#2a2a30]" />
              </motion.div>
            </div>

            {/* Screen content */}
            <div className="px-4 pt-8">
              {/* App header simulation */}
              <div className="flex items-center gap-3 mb-6">
                <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-[var(--accent)] to-[var(--accent-dark)] flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                  </svg>
                </div>
                <div>
                  <div className="text-sm font-semibold text-[var(--text-primary)]">App Shield</div>
                  <div className="text-xs text-[var(--text-tertiary)]">Protected</div>
                </div>
              </div>

              {/* Security stats */}
              <div className="grid grid-cols-2 gap-2">
                {[
                  { label: 'Scanned', value: '24' },
                  { label: 'Issues', value: '0' },
                  { label: 'Blocked', value: '12' },
                  { label: 'Safe', value: '100%' },
                ].map((stat, i) => (
                  <motion.div
                    key={stat.label}
                    initial={{ opacity: 0, scale: 0.9 }}
                    whileInView={{ opacity: 1, scale: 1 }}
                    transition={{ delay: i * 0.1 }}
                    className="p-3 rounded-2xl bg-[var(--bg-secondary)]"
                  >
                    <div className="text-lg font-bold text-[var(--text-primary)]">{stat.value}</div>
                    <div className="text-[10px] text-[var(--text-tertiary)]">{stat.label}</div>
                  </motion.div>
                ))}
              </div>
            </div>

            {/* Home indicator */}
            <div className="absolute bottom-2 left-1/2 -translate-x-1/2">
              <div className="w-32 h-1 rounded-full bg-[var(--text-tertiary)]/50" />
            </div>
          </div>
        </div>

        {/* Side button */}
        <div className="absolute right-[-2px] top-28 w-[3px] h-12 rounded-r-sm bg-[#d1d1d6]" />
        {/* Volume buttons */}
        <div className="absolute left-[-2px] top-24 w-[3px] h-6 rounded-l-sm bg-[#d1d1d6]" />
        <div className="absolute left-[-2px] top-32 w-[3px] h-10 rounded-l-sm bg-[#d1d1d6]" />
      </div>

      {/* Label */}
      <div className="text-center mt-4">
        <span className="text-sm font-medium text-[var(--text-secondary)]">iOS</span>
      </div>
    </motion.div>
  );
}
