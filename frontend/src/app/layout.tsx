import "./globals.css";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="" suppressHydrationWarning>
      <body className="transition-colors duration-300">
        {children}
      </body>
    </html>
  );
}

