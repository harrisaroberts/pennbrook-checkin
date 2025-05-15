"use client";

import { useEffect, useState } from "react";

type Member = {
  id: number;
  name: string;
  age: number;
  member_type: string;
};

export default function HomePage() {
  const [checkins, setCheckins] = useState<Member[]>([]);

  useEffect(() => {
    const fetchCheckins = async () => {
      try {
        const res = await fetch("http://localhost:8080/checkins/today");
        const data = await res.json();
        setCheckins(data);
      } catch (err) {
        console.error("Failed to fetch today's check-ins:", err);
      }
    };

    fetchCheckins();
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">🏊 Welcome to Penn Brook</h1>
      <p className="mb-6 text-gray-700">Today's check-ins:</p>

      {!checkins || checkins.length === 0 ? (
        <p className="text-gray-500">No members checked in yet today.</p>
      ) : (
        <ul className="space-y-2">
          {checkins.map((m) => (
            <li key={m.id} className="border rounded p-3 shadow-sm">
              <div className="font-semibold">{m.name}</div>
              <div className="text-sm text-gray-600">
                {m.member_type}, Age {m.age}
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

