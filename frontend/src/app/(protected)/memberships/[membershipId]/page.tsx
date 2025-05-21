"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { getUserRole } from "@/lib/auth";

type Member = {
  id: number;
  name: string;
  age: number;
  member_type: string;
  swim_test_passed: boolean;
  parent_note_on_file: boolean;
  is_checked_in?: boolean;
  membership_id?: number;
};

type Guest = {
  id: number;
  guest_name: string;
  visit_date: string;
};

export default function FamilyPage() {
  const { membershipId } = useParams();
  const [family, setFamily] = useState<Member[]>([]);
  const [guestNames, setGuestNames] = useState("");
  const [guestCount, setGuestCount] = useState<number | null>(null);
  const [activeMembershipId, setActiveMembershipId] = useState<number | null>(
    null,
  );
  const [guestFeedback, setGuestFeedback] = useState("");
  const [todayGuests, setTodayGuests] = useState<Guest[]>([]);
  const [userRole, setUserRole] = useState<string | null>(null);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  useEffect(() => {
    if (typeof membershipId === "string") {
      fetchFamily(membershipId);
    }
  }, [membershipId]);

  useEffect(() => {
    if (activeMembershipId !== null && !isNaN(activeMembershipId)) {
      fetchGuestCount(activeMembershipId);
      fetchTodayGuests(activeMembershipId);
    }
  }, [activeMembershipId]);

  const fetchFamily = async (membershipId: string) => {
    try {
      const res = await fetch(
        `http://localhost:8080/members/by-membership?id=${membershipId}`,
      );
      const text = await res.text(); // 👈 get raw text
      console.log("🔍 Raw response:", text);

      const data = JSON.parse(text); // 👈 now parse manually
      setFamily(data);
      setActiveMembershipId(parseInt(membershipId));
      console.log("📥 Family data received:", data);
    } catch (err) {
      console.error("Failed to load family data:", err);
    }
  };
  const fetchGuestCount = async (membershipId: number) => {
    try {
      const res = await fetch(
        `http://localhost:8080/guests/monthly-total?id=${membershipId}`,
      );
      const data = await res.json();
      console.log("👥 Guest count fetched:", data);
      setGuestCount(data.guest_count);
    } catch (err) {
      console.error("Failed to fetch guest count:", err);
    }
  };

  const fetchTodayGuests = async (membershipId: number) => {
    try {
      const res = await fetch(
        `http://localhost:8080/guests/today?membership_id=${membershipId}`,
      );
      const data = await res.json();
      setTodayGuests(data);
    } catch (err) {
      console.error("Failed to fetch today's guests:", err);
    }
  };

  const handleCheckin = async (memberId: number, isCheckedIn: boolean) => {
    console.log("🟢 Check-in clicked:", { memberId, isCheckedIn });
    const url = `http://localhost:8080/checkins${isCheckedIn ? `/${memberId}` : ""}`;
    const res = await fetch(url, {
      method: isCheckedIn ? "DELETE" : "POST",
      headers: { "Content-Type": "application/json" },
      body: isCheckedIn ? undefined : JSON.stringify({ member_id: memberId }),
    });

    if (res.ok) {
      fetchFamily(membershipId as string);
    } else {
      console.warn("❌ Check-in failed");
    }
  };

  const handleGuestSubmit = async () => {
    if (!guestNames.trim() || activeMembershipId === null) return;

    const nameList = guestNames
      .split(",")
      .map((name) => name.trim())
      .filter(Boolean);

    const guestCountBeingAdded = nameList.length;

    try {
      const res = await fetch("http://localhost:8080/guests", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          membership_id: activeMembershipId,
          guest_names: guestNames,
        }),
      });

      if (res.ok) {
        setGuestNames("");
        setGuestFeedback(
          `${guestCountBeingAdded} guest(s) signed in successfully.`,
        );
        fetchGuestCount(activeMembershipId);
        fetchTodayGuests(activeMembershipId);
        setTimeout(() => setGuestFeedback(""), 4000);
      } else {
        setGuestFeedback("Failed to sign in guests.");
      }
    } catch (err) {
      setGuestFeedback("Something went wrong while adding guests.", err);
    }
  };

  const handleDeleteGuest = async (guestId: number) => {
    try {
      const res = await fetch(`http://localhost:8080/guests/${guestId}`, {
        method: "DELETE",
      });

      if (res.status === 204 && activeMembershipId !== null) {
        setTodayGuests((prev) => prev.filter((g) => g.id !== guestId));
        await fetchGuestCount(activeMembershipId);
      } else {
        console.error("Failed to delete guest");
      }
    } catch (err) {
      console.error("Error deleting guest:", err);
    }
  };

  // (Optional) Placeholder functions for toggles
  const handleSwimTestToggle = (id: number, current: boolean) => {
    console.log("Swim test toggle for", id, "=>", !current);
  };

  const handleParentNoteToggle = (id: number, current: boolean) => {
    console.log("Parent note toggle for", id, "=>", !current);
  };

  return (
    <div className="px-6 sm:px-12 md:px-20 lg:px-32 xl:px-48 py-6">
      <div className="hidden md:grid grid-cols-6 font-semibold text-gray-600 px-2 py-2 border-b">
        <span>Name</span>
        <span>Type</span>
        <span>Age</span>
        <span>Swim Test</span>
        <span>Parent Note</span>
        <span className="text-right">Check-In</span>
      </div>
      <ul className="space-y-2">
        {family.map((member) => (
          <li
            key={member.id}
            className={`grid grid-cols-6 items-center text-sm px-2 py-3 border rounded shadow-sm ${
              member.is_checked_in ? "bg-green-100" : ""
            }`}
          >
            {/* Name & Edit */}
            <div className="font-medium flex items-center gap-2">
              {member.name}
              {(userRole === "guard" || userRole === "admin") && (
                <a
                  href={`/members/${member.id}/edit`}
                  className="text-xs text-blue-600 underline hover:text-blue-800"
                >
                  Edit
                </a>
              )}
            </div>

            <div>{member.member_type}</div>
            <div>{member.age}</div>

            {/* Swim Test */}
            {member.age < 13 ? (
              <div className="flex items-center gap-2">
                {member.swim_test_passed ? "✅" : "❌"}
              </div>
            ) : (
              <span className="text-gray-400 text-xs italic">N/A</span>
            )}

            {/* Parent Note */}
            {member.age < 13 ? (
              <div className="flex items-center gap-2">
                {member.parent_note_on_file ? "✅" : "❌"}
              </div>
            ) : (
              <span className="text-gray-400 text-xs italic">N/A</span>
            )}

            <div className="flex justify-end">
              <button
                onClick={() =>
                  handleCheckin(member.id, member.is_checked_in || false)
                }
                className="bg-blue-600 text-white text-xs px-4 py-2 rounded hover:bg-blue-700"
              >
                {member.is_checked_in ? "Undo Check-In" : "Check In"}
              </button>
            </div>
          </li>
        ))}
      </ul>

      <div className="mt-10 border-t pt-6">
        {guestCount !== null && (
          <p className="mb-6 text-gray-700 font-medium">
            Guests this month: <span className="font-bold">{guestCount}</span>
          </p>
        )}

        <label htmlFor="guests" className="block font-semibold mb-1">
          Add Guests (comma-separated):
        </label>
        <input
          type="text"
          id="guests"
          value={guestNames}
          onChange={(e) => setGuestNames(e.target.value)}
          placeholder="e.g. John Doe, Jane Smith"
          className="w-full p-2 border rounded mb-2"
          disabled={activeMembershipId === null}
        />
        <button
          onClick={handleGuestSubmit}
          disabled={activeMembershipId === null}
          className={`px-4 py-2 rounded text-white ${
            activeMembershipId === null
              ? "bg-gray-400 cursor-not-allowed"
              : "bg-blue-600 hover:bg-blue-700"
          }`}
        >
          Add Guests
        </button>

        {guestFeedback && (
          <p className="mt-2 text-green-700 font-medium">{guestFeedback}</p>
        )}

        {todayGuests && todayGuests.length > 0 && (
          <div className="mt-4">
            <h3 className="font-semibold mb-2">Today&apos;s Guests</h3>
            <ul className="space-y-2">
              {todayGuests.map((guest) => (
                <li
                  key={guest.id}
                  className="flex justify-between items-center border rounded p-2 bg-gray-50"
                >
                  <span>{guest.guest_name}</span>
                  <button
                    onClick={() => handleDeleteGuest(guest.id)}
                    className="text-red-600 hover:underline text-sm"
                  >
                    Remove
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}
