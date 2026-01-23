import type { MDXComponents } from 'mdx/types'

export function useMDXComponents(components: MDXComponents): MDXComponents {
  return {
    h1: ({ children }) => (
      <h1 className="text-4xl font-bold text-matrix-green mb-6 font-mono">{children}</h1>
    ),
    h2: ({ children }) => (
      <h2 className="text-2xl font-bold text-matrix-green mb-4 mt-8 font-mono border-b border-matrix-green/30 pb-2">{children}</h2>
    ),
    h3: ({ children }) => (
      <h3 className="text-xl font-semibold text-matrix-green/90 mb-3 mt-6 font-mono">{children}</h3>
    ),
    p: ({ children }) => (
      <p className="text-matrix-green/70 mb-4 leading-relaxed">{children}</p>
    ),
    code: ({ children }) => (
      <code className="bg-matrix-green/10 text-matrix-green px-1.5 py-0.5 rounded font-mono text-sm border border-matrix-green/20">{children}</code>
    ),
    pre: ({ children }) => (
      <pre className="bg-black/50 border border-matrix-green/30 rounded-lg p-4 overflow-x-auto mb-4 font-mono text-sm">{children}</pre>
    ),
    ul: ({ children }) => (
      <ul className="list-none space-y-2 mb-4 ml-4">{children}</ul>
    ),
    li: ({ children }) => (
      <li className="text-matrix-green/70 flex items-start gap-2">
        <span className="text-matrix-green mt-1">▸</span>
        <span>{children}</span>
      </li>
    ),
    a: ({ href, children }) => (
      <a href={href} className="text-matrix-green hover:text-matrix-bright underline underline-offset-2 transition-colors">{children}</a>
    ),
    blockquote: ({ children }) => (
      <blockquote className="border-l-2 border-matrix-green/50 pl-4 italic text-matrix-green/60 my-4">{children}</blockquote>
    ),
    ...components,
  }
}
