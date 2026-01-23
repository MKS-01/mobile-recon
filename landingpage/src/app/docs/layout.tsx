import Link from 'next/link';
import { Terminal, ChevronLeft, Github } from 'lucide-react';
import DocsSidebar from '@/components/DocsSidebar';
import DocsThemeToggle from '@/components/DocsThemeToggle';

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      {/* Top nav */}
      <header className="fixed top-0 left-0 right-0 z-50 h-14 bg-[var(--bg-primary)]/80 backdrop-blur-md border-b border-[var(--border)]">
        <div className="h-full px-6 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link
              href="/"
              className="flex items-center gap-2 text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition-colors"
            >
              <ChevronLeft className="w-4 h-4" />
              <span className="text-sm">Back</span>
            </Link>
            <div className="h-5 w-px bg-[var(--border)]" />
            <Link href="/docs" className="flex items-center gap-2">
              <div className="w-7 h-7 bg-[var(--accent)] rounded-lg flex items-center justify-center">
                <Terminal className="w-3.5 h-3.5 text-white" />
              </div>
              <span className="font-semibold text-[var(--text-primary)]">Docs</span>
            </Link>
          </div>
          <div className="flex items-center gap-2">
            <DocsThemeToggle />
            <a
              href="https://github.com/MKS-01/mobile-recon"
              target="_blank"
              rel="noopener noreferrer"
              className="w-8 h-8 flex items-center justify-center rounded-lg text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-secondary)] transition-colors"
            >
              <Github className="w-4 h-4" />
            </a>
          </div>
        </div>
      </header>

      <div className="flex pt-14">
        {/* Sidebar */}
        <aside className="fixed left-0 top-14 bottom-0 w-60 border-r border-[var(--border)] bg-[var(--bg-primary)] overflow-y-auto hidden md:block">
          <DocsSidebar />
        </aside>

        {/* Main content */}
        <main className="flex-1 md:ml-60 p-6 md:p-10">
          <div className="max-w-3xl mx-auto">
            <article className="docs-content">
              {children}
            </article>
          </div>
        </main>
      </div>
    </div>
  );
}
