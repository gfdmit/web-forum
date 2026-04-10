INSERT INTO forum.users (login, password_hash)
VALUES
    ('admin', 'admin');

INSERT INTO forum.profiles (user_id, university_id, firstname, lastname, middlename, birthday, faculty, grade, "group", status)
VALUES
    (1, '1', 'admin', 'admin', 'admin', '1.1.2001', 'admin', 'admin', 'admin', 'admin');

INSERT INTO forum.boards (name, description)
VALUES 
    ('b', 'Бред'),
    ('news', 'Ньюсач')
ON CONFLICT (name) DO NOTHING;

INSERT INTO forum.posts (user_id, board_id, title, text)
VALUES
    (1, 1, 'First post!', 'Hello everyone!'),
    (1, 2, 'New feature', 'Check out the update!'),
    (1, 2, 'Third post!', 'Wazzap!');

INSERT INTO forum.comments (user_id, post_id, text)
VALUES
    (1, 1, 'Welcome!'),
    (1, 1, 'Great start!');