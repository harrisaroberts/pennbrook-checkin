-- Create the memberships table
CREATE TABLE memberships (
    id SERIAL PRIMARY KEY,
    quickbooks_id TEXT,
    membership_type TEXT NOT NULL,  -- family, sustaining, individual, etc.
    status TEXT NOT NULL,           -- paid, unpaid, inactive, etc.
    guest_limit_monthly INTEGER DEFAULT 10,
    guest_limit_total INTEGER DEFAULT 20,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the members table
CREATE TABLE members (
    id SERIAL PRIMARY KEY,
    membership_id INTEGER REFERENCES memberships(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    age INTEGER,
    member_type TEXT NOT NULL,         -- child, adult, caregiver, etc.
    swim_test_passed BOOLEAN DEFAULT FALSE,
    parent_note_on_file BOOLEAN DEFAULT FALSE,
    caregiver_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the guest visits table
CREATE TABLE guests (
    id SERIAL PRIMARY KEY,
    membership_id INTEGER REFERENCES memberships(id) ON DELETE CASCADE,
    guest_name TEXT NOT NULL,
    visit_date DATE NOT NULL DEFAULT CURRENT_DATE,
    entered_by TEXT,
    notes TEXT
);

CREATE TABLE IF NOT EXISTS checkins (
  id SERIAL PRIMARY KEY,
  member_id INT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
  checkin_date DATE NOT NULL DEFAULT CURRENT_DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (member_id, checkin_date)
);

