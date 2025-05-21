"use client";

import { useEffect, useState } from "react";
import { MoonIcon, SunIcon } from "@heroicons/react/24/solid";

export default function ThemeToggle() {
  const [isDark, setIsDark] = useState(false);

  useEffect(() => {
    const savedTheme = localStorage.getItem("theme");
    if (savedTheme === "dark") {
      document.documentElement.classList.add("dark");
      setIsDark(true);
    }
  }, []);

  const toggleTheme = () => {
    const root = document.documentElement;
    if (isDark) {
      root.classList.remove("dark");
      localStorage.setItem("theme", "light");
      setIsDark(false);
    } else {
      root.classList.add("dark");
      localStorage.setItem("theme", "dark");
      setIsDark(true);
    }
  };

  return (
    <label className="relative inline-block w-14 h-8 cursor-pointer select-none">
      <input
        type="checkbox"
        checked={isDark}
        onChange={toggleTheme}
        className="sr-only peer"
      />
      {/* OUTER BACKGROUND (switch track) */}
      <div className="w-full h-full bg-white dark:bg-black border border-gray-300 rounded-full transition-colors duration-300"></div>

      {/* MOVING THUMB */}
      <div
        className={`
          absolute top-0.5 left-0.5 w-7 h-7 bg-gray-100 dark:bg-neutral-800 rounded-full shadow-md flex items-center justify-center
          transition-transform duration-300
          peer-checked:translate-x-6
        `}
      >
        {isDark ? (
          <MoonIcon className="h-4 w-4 text-white" />
        ) : (
          <SunIcon className="h-4 w-4 text-yellow-400" />
        )}
      </div>
    </label>
  );
}

