-- Create one fake membership
INSERT INTO memberships (membership_type, status, quickbooks_id)
VALUES ('family', 'paid', 'QB-SEED')
RETURNING id;

-- Assuming that returns ID 1 (if not, check with SELECT * FROM memberships)

-- Now insert 100 members for that membership
DO $$
BEGIN
  FOR i IN 1..100 LOOP
    INSERT INTO members (
      membership_id, name, age, member_type, swim_test_passed, parent_note_on_file
    )
    VALUES (
      1,
      CONCAT('Member ', i),
      (10 + (random() * 40)::int),
      CASE WHEN i % 3 = 0 THEN 'adult' ELSE 'child' END,
      (random() > 0.5),
      (random() > 0.3)
    );
  END LOOP;
END$$;

