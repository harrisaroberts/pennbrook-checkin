-- Clear existing data
DELETE FROM members;
DELETE FROM memberships;

-- Seed 50 families with 4–6 members each
DO $$
DECLARE
  i INT;
  member_count INT;
  first_names TEXT[] := ARRAY[
    'Alex', 'Jordan', 'Taylor', 'Morgan', 'Casey', 'Riley', 'Jamie', 'Skyler', 'Avery', 'Parker',
    'Dakota', 'Emerson', 'Reagan', 'Kendall', 'Quinn', 'Rowan', 'Sawyer', 'Finley', 'Blake', 'Payton',
    'Charlie', 'Drew', 'Logan', 'Cameron', 'Sage', 'Harper', 'Jules', 'Tatum', 'Spencer', 'Remy'
  ];
  last_names TEXT[] := ARRAY[
    'Smith', 'Johnson', 'Brown', 'Davis', 'Miller', 'Wilson', 'Moore', 'Taylor', 'Anderson', 'Thomas',
    'Jackson', 'White', 'Harris', 'Martin', 'Thompson', 'Garcia', 'Martinez', 'Robinson', 'Clark', 'Rodriguez',
    'Lewis', 'Lee', 'Walker', 'Hall', 'Allen', 'Young', 'King', 'Wright', 'Scott', 'Green', 'Nelson', 'Reed',
    'Cook', 'Morgan', 'Bell', 'Murphy', 'Bailey', 'Rivera', 'Cooper', 'Richardson', 'Cox', 'Howard', 'Ward',
    'Torres', 'Peterson', 'Gray', 'Ramirez', 'James', 'Watson', 'Brooks', 'Kelly'
  ];
  shuffled_names TEXT[];
  selected_names TEXT[];
  family_id INT;
BEGIN
  FOR i IN 1..50 LOOP
    -- Insert a new family
    INSERT INTO memberships (membership_type, status, quickbooks_id)
    VALUES ('family', 'paid', 'QB-' || i)
    RETURNING id INTO family_id;

    -- Shuffle and slice 4–6 unique names
    shuffled_names := ARRAY(SELECT unnest(first_names) ORDER BY random());
    member_count := 4 + (random() * 2)::int;
    selected_names := shuffled_names[1:member_count];

    -- Insert members with those names
    FOR j IN 1..array_length(selected_names, 1) LOOP
      INSERT INTO members (
        membership_id,
        name,
        age,
        member_type,
        swim_test_passed,
        parent_note_on_file
      ) VALUES (
        family_id,
        selected_names[j] || ' ' || last_names[i],
        5 + (random() * 50)::int,
        CASE WHEN j = 1 THEN 'adult' ELSE 'child' END,
        (random() > 0.5),
        (random() > 0.3)
      );
    END LOOP;
  END LOOP;
END$$;

