"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

interface Member {
  id: number;
  name: string;
  age: number;
  member_type: string;
  swim_test_passed: boolean;
  parent_note_on_file: boolean;
}

export default function MemberFamilyPage() {
  const { id } = useParams();
  const [family, setFamily] = useState<Member[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    const fetchFamily = async () => {
      try {
        const res = await fetch(`http://localhost:8080/members/family?id=${id}`);
        if (!res.ok) throw new Error("Failed to fetch family");
        const data = await res.json();
        setFamily(data);
      } catch (err) {
        setError("Failed to load family data");
      }
    };
    fetchFamily();
  }, [id]);

  if (error) {
    return <div className="p-6 text-red-500">{error}</div>;
  }

  if (!family) {
    return <div className="p-6 text-gray-500">Loading family info...</div>;
  }

  return (
    <div className="p-6 max-w-2xl mx-auto">
      <h1 className="text-2xl font-bold mb-4">👨‍👩‍👧 Family Members</h1>
      <ul className="space-y-4">
        {family.map((member) => (
          <li key={member.id} className="border p-4 rounded shadow">
            <div className="text-lg font-semibold">{member.name}</div>
            <div className="text-sm text-gray-700">
              <p>Age: {member.age}</p>
              <p>Type: {member.member_type}</p>
              <p>Swim Test: {member.swim_test_passed ? "✅ Passed" : "❌ Not Passed"}</p>
              <p>Parent Note: {member.parent_note_on_file ? "📝 On File" : "❌ None"}</p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
} 

