"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

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

export default function FamilyPage() {
    const { id } = useParams();
    const [family, setFamily] = useState<Member[]>([]);
    const [guestNames, setGuestNames] = useState("");
    const [guestCount, setGuestCount] = useState<number | null>(null);
    const [membershipId, setMembershipId] = useState<number | null>(null);
    const [guestFeedback, setGuestFeedback] = useState("");
    const [todayGuests, setTodayGuests] = useState<Guest[]>([]);

    type Guest = {
        id: number;
        guest_name: string;
        visit_date: string;
    };
    const fetchTodayGuests = async (membershipId: number) => {
        try {
            const res = await fetch(`http://localhost:8080/guests/today?membership_id=${membershipId}`);
            const data = await res.json();
            setTodayGuests(data);
        } catch (err) {
            console.error("Failed to fetch today's guests:", err);
        }
    };
    const handleDeleteGuest = async (guestId: number) => {
        try {
            const res = await fetch(`http://localhost:8080/guests/${guestId}`, {
                method: "DELETE",
            });

            if (res.status === 204) {
                setTodayGuests(prev => prev.filter(g => g.id !== guestId));
            } else {
                console.error("Failed to delete guest");
            }
        } catch (err) {
            console.error("Error deleting guest:", err);
        }
    };



    const fetchFamily = async () => {
        try {
            const res = await fetch(`http://localhost:8080/members/family?id=${id}`);
            const data = await res.json();
            setFamily(data);
            if (data.length > 0 && data[0].membership_id) setMembershipId(data[0].membership_id);
            console.log("📥 Family data received:", data);
        } catch (err) {
            console.error("Failed to load family data:", err);
        }
    };

    const fetchGuestCount = async (membershipId: number) => {
        try {
            const res = await fetch(`http://localhost:8080/guests/monthly-total?id=${membershipId}`);
            const data = await res.json();
            console.log("👥 Guest count fetched:", data);
            setGuestCount(data.guest_count);
        } catch (err) {
            console.error("Failed to fetch guest count:", err);
        }
    };

    useEffect(() => {
        fetchFamily();
    }, [id]);

    useEffect(() => {
        if (membershipId !== null && !isNaN(membershipId)) {
            fetchGuestCount(membershipId);
            fetchTodayGuests(membershipId);
        }
    }, [membershipId]);

    const handleCheckin = async (memberId: number, isCheckedIn: boolean) => {
        console.log("🟢 Check-in clicked:", { memberId, isCheckedIn });
        const url = `http://localhost:8080/checkins${isCheckedIn ? `/${memberId}` : ""}`;
        const res = await fetch(url, {
            method: isCheckedIn ? "DELETE" : "POST",
            headers: { "Content-Type": "application/json" },
            body: isCheckedIn ? undefined : JSON.stringify({ member_id: memberId }),
        });

        console.log("📦 Response status:", res.status);
        if (res.ok) {
            console.log("✅ Check-in updated successfully");
            fetchFamily();
        } else {
            console.warn("❌ Check-in failed");
        }
    };

    const handleGuestSubmit = async () => {
        if (!guestNames.trim()) {
            console.warn("⚠️ Guest name input is empty");
            return;
        }

        if (membershipId === null) {
            console.error("❌ No membership ID set");
            return;
        }

        const nameList = guestNames
        .split(",")
        .map((name) => name.trim())
        .filter(Boolean);

        const guestCountBeingAdded = nameList.length;
        console.log("🧾 Parsed guests:", nameList);

        if (guestCountBeingAdded === 0) {
            console.warn("⚠️ No valid guest names parsed");
            return;
        }

        try {
            const res = await fetch("http://localhost:8080/guests", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    membership_id: membershipId,
                    guest_names: guestNames,
                }),
            });

            if (res.ok) {
                console.log("✅ Guests added");
                setGuestNames("");
                setGuestFeedback(`${guestCountBeingAdded} guest${guestCountBeingAdded > 1 ? "s" : ""} signed in successfully.`);
                fetchGuestCount(membershipId);
                fetchTodayGuests(membershipId);
                setTimeout(() => setGuestFeedback(""), 4000);
            } else {
                const errorText = await res.text();
                console.error("❌ Server error:", res.status, errorText);
                setGuestFeedback("Failed to sign in guests.");
            }
        } catch (err) {
            console.error("🚨 Fetch failed:", err);
            setGuestFeedback("Something went wrong while adding guests.");
        }
    };

    return (
        <div>
            <h1 className="text-2xl font-bold mb-4">👨‍👩‍👧‍👦 Family Check-In</h1>

            {guestCount !== null && (
                <p className="mb-6 text-gray-700 font-medium">
                    Guests this month: <span className="font-bold">{guestCount}</span>
                </p>
            )}

            <ul className="space-y-4">
                {family.map((member) => (
                    <li
                        key={member.id}
                        className={`border rounded p-4 shadow ${member.is_checked_in ? "bg-green-100" : ""}`}
                    >
                        <div className="font-semibold text-lg">{member.name}</div>
                        <div className="text-sm text-gray-700">
                            {member.member_type}, Age {member.age}
                        </div>
                        <div className="text-sm text-gray-600">
                            Swim Test: {member.swim_test_passed ? "✅" : "❌"} | Parent Note:{" "}
                            {member.parent_note_on_file ? "✅" : "❌"}
                        </div>
                        <button
                            onClick={() => handleCheckin(member.id, member.is_checked_in || false)}
                            className="mt-2 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                        >
                            {member.is_checked_in ? "Undo Check-In" : "Check In"}
                        </button>
                    </li>
                ))}
            </ul>

            <div className="mt-10 border-t pt-6">
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
                    disabled={membershipId === null}
                />
                <button
                    onClick={handleGuestSubmit}
                    disabled={membershipId === null}
                    className={`px-4 py-2 rounded text-white ${
membershipId === null
? "bg-gray-400 cursor-not-allowed"
: "bg-blue-600 hover:bg-blue-700"
}`}
                >
                    Add Guests
                </button>

                {guestFeedback && (
                    <p className="mt-2 text-green-700 font-medium">{guestFeedback}</p>
                )}
                {todayGuests.length > 0 && (
                    <div className="mt-4">
                        <h3 className="font-semibold mb-2">Today's Guests</h3>
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

