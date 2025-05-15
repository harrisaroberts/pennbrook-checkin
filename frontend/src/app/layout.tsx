import Link from "next/link";
import "./globals.css"; // or wherever your styles are

export const metadata = {
  title: "Penn Brook Check-In",
  description: "Pool check-in and member management system",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <nav className="bg-blue-600 text-white p-4 flex gap-6">
          <Link href="/" className="hover:underline">🏠 Home</Link>
          <Link href="/members" className="hover:underline">🔍 Check-In</Link>
        </nav>
        <main className="p-6">{children}</main>
      </body>
    </html>
  );
}

