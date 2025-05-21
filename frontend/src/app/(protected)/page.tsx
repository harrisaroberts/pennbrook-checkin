"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

interface Member {
  id: number;
  name: string;
  age: number;
  member_type: string;
  membership_id: number;
}

export default function MemberSearchPage() {
  const [allMembers, setAllMembers] = useState<Member[]>([]);
  const [query, setQuery] = useState("");
  const [filtered, setFiltered] = useState<Member[]>([]);

  // Fetch all members once
  useEffect(() => {
    const fetchAll = async () => {
      try {
        const res = await fetch("http://localhost:8080/members?name=");
        const data = await res.json();
        setAllMembers(data);
      } catch (err) {
        console.error("Failed to fetch members:", err);
      }
    };
    fetchAll();
  }, []);

  // Filter members as the user types
  useEffect(() => {
    const q = query.toLowerCase();
    const matches = allMembers.filter((m) => m.name.toLowerCase().includes(q));
    setFiltered(matches);
  }, [query, allMembers]);

  return (
    <div className="p-8 max-w-xl mx-auto">
      <h1 className="text-center text-2xl font-bold mb-4">Search Members</h1>
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Start typing a name..."
        className="w-full px-4 py-2 border border-blue-500 rounded"
      />

      {query && (
        <ul className="mt-4 border rounded divide-y">
          {filtered.length > 0 ? (
            filtered.map((member) => (
              <li
                key={member.id}
                className="border p-4 rounded shadow hover:bg-gray-50 transition"
              >
                <Link href={`/memberships/${member.membership_id}`}>
                  <div>
                    <div className="font-semibold">{member.name}</div>
                    <div className="text-sm text-gray-600">
                      {member.member_type}, Age {member.age}
                    </div>
                  </div>
                </Link>
              </li>
            ))
          ) : (
            <li className="p-3 text-gray-500">No matches found</li>
          )}
        </ul>
      )}
    </div>
  );
}
