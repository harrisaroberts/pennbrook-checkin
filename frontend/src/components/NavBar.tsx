"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { isLoggedIn, logout } from "@/lib/auth";
import ThemeToggle from "@/components/ThemeToggle"; // update path if needed

export default function NavBar() {
  const router = useRouter();
  const pathname = usePathname();

  const [hydrated, setHydrated] = useState(false);
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setHydrated(true);
    setLoggedIn(isLoggedIn());
  }, []);

  if (!hydrated || !loggedIn || pathname === "/login") return null;

  return (
    <nav className="bg-blue-600 text-white py-2 px-4 shadow-sm relative">
      <div className="max-w-5xl mx-auto relative h-10 flex items-center justify-center">
        
        {/* LEFT: Theme toggle */}
        <div className="absolute left-0">
          <ThemeToggle />
        </div>

        {/* CENTER: Title */}
        <div className="text-lg font-semibold text-center mx-auto">
          Penn Brook Check-In System
        </div>

        {/* RIGHT: Navigation links */}
        <div className="absolute right-0 flex gap-4 text-sm items-center">
          <a href="/" className="hover:underline">Check-In</a>
          <button
            onClick={() => {
              logout();
              router.push("/login");
            }}
            className="hover:underline"
          >
            Logout
          </button>
        </div>
      </div>
    </nav>
  );
}

