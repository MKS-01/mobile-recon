'use client';

import Link from 'next/link';
import { Github, Twitter, Heart, Terminal } from 'lucide-react';

export default function Footer() {
  return (
    <footer className="border-t border-[var(--border)] bg-[var(--bg-secondary)]">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="flex flex-col md:flex-row items-center justify-between gap-8">
          {/* Logo & tagline */}
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-[var(--accent)] flex items-center justify-center">
              <Terminal className="w-4 h-4 text-white" />
            </div>
            <div>
              <span className="font-semibold text-[var(--text-primary)]">mobile-recon</span>
              <p className="text-xs text-[var(--text-tertiary)]">A weekend project for fun</p>
            </div>
          </div>

          {/* Links */}
          <div className="flex items-center gap-6">
            <Link href="/docs" className="text-sm text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors">
              Docs
            </Link>
            <a
              href="https://github.com/MKS-01/mobile-recon"
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors flex items-center gap-1.5"
            >
              <Github className="w-4 h-4" />
              GitHub
            </a>
          </div>
        </div>

        {/* Disclaimer */}
        <div className="mt-8 pt-8 border-t border-[var(--border)]">
          <p className="text-xs text-[var(--text-tertiary)] text-center max-w-2xl mx-auto">
            This toolkit is for authorized security testing and educational purposes only.
            Always get proper permission before testing. MIT License.
          </p>
        </div>

        {/* Bottom */}
        <div className="mt-6 flex items-center justify-center gap-1 text-xs text-[var(--text-tertiary)]">
          <span>Made with</span>
          <Heart className="w-3 h-3 text-red-500 fill-red-500" />
          <span>for tinkerers</span>
        </div>
      </div>
    </footer>
  );
}
