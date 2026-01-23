import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import { ThemeProvider } from "@/components/ThemeProvider";
import "./globals.css";

const inter = Inter({
  variable: "--font-sans",
  subsets: ["latin"],
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  variable: "--font-mono",
  subsets: ["latin"],
  display: "swap",
});

export const metadata: Metadata = {
  title: "mobile-recon | Mobile Security Toolkit",
  description: "A fun toolkit for mobile security testing on Android & iOS. Built for security researchers and curious developers.",
  keywords: ["mobile security", "penetration testing", "android security", "ios security", "apk analysis"],
  authors: [{ name: "MKS-01" }],
  openGraph: {
    title: "mobile-recon | Mobile Security Toolkit",
    description: "A fun toolkit for mobile security testing on Android & iOS. Built for security researchers and curious developers.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.variable} ${jetbrainsMono.variable} antialiased`}>
        <ThemeProvider>
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
