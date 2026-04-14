DROP INDEX IF EXISTS forum.idx_comments_post_id;
DROP INDEX IF EXISTS forum.idx_comments_user_id;
DROP INDEX IF EXISTS forum.idx_posts_board_id;
DROP INDEX IF EXISTS forum.idx_posts_user_id;
DROP INDEX IF EXISTS forum.idx_profiles_user_id;

DROP TABLE IF EXISTS forum.comments;
DROP TABLE IF EXISTS forum.posts;
DROP TABLE IF EXISTS forum.boards;
DROP TABLE IF EXISTS forum.profiles;
DROP TABLE IF EXISTS forum.users;

DROP SCHEMA IF EXISTS forum CASCADE;