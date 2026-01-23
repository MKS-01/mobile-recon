'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { BookOpen, Smartphone, Search, Apple, Network, Cpu, Zap, GitBranch, type LucideIcon } from 'lucide-react';

interface SidebarItem {
  href: string;
  label: string;
  icon?: LucideIcon;
}

interface SidebarSection {
  title: string;
  items: SidebarItem[];
}

const sidebarItems: SidebarSection[] = [
  {
    title: 'Getting Started',
    items: [
      { href: '/docs', label: 'Introduction', icon: BookOpen },
      { href: '/docs/installation', label: 'Installation' },
      { href: '/docs/quick-start', label: 'Quick Start' },
    ],
  },
  {
    title: 'Tools',
    items: [
      { href: '/docs/adb-toolkit', label: 'ADB Toolkit', icon: Smartphone },
      { href: '/docs/apk-analyzer', label: 'APK Analyzer', icon: Search },
      { href: '/docs/ios-toolkit', label: 'iOS Toolkit', icon: Apple },
      { href: '/docs/nmap-toolkit', label: 'Nmap Toolkit', icon: Network },
    ],
  },
  {
    title: 'Advanced',
    items: [
      { href: '/docs/frida-integration', label: 'Frida Integration', icon: Cpu },
      { href: '/docs/automation', label: 'Automation', icon: Zap },
      { href: '/docs/contributing', label: 'Contributing', icon: GitBranch },
    ],
  },
];

export default function DocsSidebar() {
  const pathname = usePathname();

  return (
    <nav className="p-4 space-y-6">
      {sidebarItems.map((section) => (
        <div key={section.title}>
          <h3 className="text-[var(--text-tertiary)] text-xs font-medium uppercase tracking-wider mb-2 px-3">
            {section.title}
          </h3>
          <ul className="space-y-0.5">
            {section.items.map((item) => {
              const isActive = pathname === item.href;
              const Icon = item.icon;

              return (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    className={`flex items-center gap-2.5 px-3 py-2 rounded-lg transition-colors text-sm ${
                      isActive
                        ? 'bg-[var(--accent-muted)] text-[var(--accent)] font-medium'
                        : 'text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--bg-secondary)]'
                    }`}
                  >
                    {Icon && <Icon className="w-4 h-4" />}
                    {item.label}
                  </Link>
                </li>
              );
            })}
          </ul>
        </div>
      ))}
    </nav>
  );
}
