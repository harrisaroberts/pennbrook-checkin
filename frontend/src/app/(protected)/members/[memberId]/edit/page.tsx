"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { getUserRole } from "@/lib/auth";

type Member = {
  id: number;
  name: string;
  age: number;
  member_type: string;
  swim_test_passed: boolean;
  parent_note_on_file: boolean;
  membership_id: number;
};

export default function EditMemberPage() {
  const { memberId } = useParams();
  const router = useRouter();
  const [role, setRole] = useState<string | null>(null);
  const [member, setMember] = useState<Member | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    setRole(getUserRole());
    fetch(`http://localhost:8080/members/${memberId}`)
      .then((res) => res.json())
      .then((data) => setMember(data))
      .catch(() => setError("Failed to load member data."));
  }, [memberId]);

  const handleSubmit = async () => {
    try {
      const res = await fetch(`http://localhost:8080/members/${memberId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(member),
      });

      if (res.ok) {
        router.push(`/memberships/${member.membership_id}`);
      } else {
        setError("Failed to update member.");
      }
    } catch (err) {
      setError("Error submitting form.", err);
    }
  };

  const updateField = (field: keyof Member, value: unkown) => {
    setMember((prev) => (prev ? { ...prev, [field]: value } : prev));
  };

  if (!member) return <p>Loading...</p>;

  return (
    <div className="max-w-xl mx-auto mt-10 p-6 border rounded shadow">
      <h1 className="text-xl font-bold mb-4">Edit Member</h1>

      {error && <p className="text-red-600 mb-4">{error}</p>}

      {role === "admin" && (
        <>
          <label>Name</label>
          <input
            className="w-full mb-4 p-2 border rounded"
            value={member.name}
            onChange={(e) => updateField("name", e.target.value)}
          />

          <label>Age</label>
          <input
            type="number"
            className="w-full mb-4 p-2 border rounded"
            value={member.age}
            onChange={(e) => updateField("age", parseInt(e.target.value))}
          />

          <label>Type</label>
          <select
            className="w-full mb-4 p-2 border rounded"
            value={member.member_type}
            onChange={(e) => updateField("member_type", e.target.value)}
          >
            <option value="adult">Adult</option>
            <option value="child">Child</option>
            <option value="caregiver">Caregiver</option>
          </select>
        </>
      )}

      {(role === "guard" || role === "admin") && (
        <>
          <label>Swim Test Passed</label>
          <input
            type="checkbox"
            className="ml-2"
            checked={member.swim_test_passed}
            onChange={(e) => updateField("swim_test_passed", e.target.checked)}
          />

          <br />

          <label>Parent Note on File</label>
          <input
            type="checkbox"
            className="ml-2"
            checked={member.parent_note_on_file}
            onChange={(e) =>
              updateField("parent_note_on_file", e.target.checked)
            }
          />
        </>
      )}

      <button
        onClick={handleSubmit}
        className="mt-6 w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
      >
        Save Changes
      </button>
    </div>
  );
}
