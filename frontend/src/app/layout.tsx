import type { Metadata } from "next";
import "./globals.css";
import "./landing.css";

export const metadata: Metadata = {
  title: "Entré — 散らかった就活、ぜんぶ1枚に。",
  description:
    "マイナビ、リクナビ、ワンキャリ、企業HP…応募ごとに1箇所へ集める、就活の台帳。",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja" className="h-full antialiased">
      <body className="min-h-full flex flex-col">{children}</body>
    </html>
  );
}
