"use client";

import { useState } from "react";

export default function MemberSearchPage() {
    const [query, setQuery] = useState("");
    const [members, setMembers] = useState<any[]>([]);

    const handleSearch = async () => {
        console.log("Search triggered:", query); // should log when button is clicked

        const res = await fetch(`http://localhost:8080/members?name=${query}`);
        const text = await res.text();
        console.log("Raw response:", text);
    };

    return (
        <div className="p-8 max-w-xl mx-auto">
            <h1 className="text-2xl font-bold mb-4">🔍 Search Members</h1>
            <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Enter member name..."
                className="w-full px-4 py-2 border rounded mb-4"
            />
            <button
                onClick={handleSearch}
                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
            >
                Search
            </button>

            <ul className="mt-6 space-y-2">
                {members.map((member) => (
                    <li key={member.id} className="border p-4 rounded shadow">
                        <strong>{member.name}</strong> — {member.member_type}, Age {member.age}
                    </li>
                ))}
            </ul>
        </div>
    );
}

