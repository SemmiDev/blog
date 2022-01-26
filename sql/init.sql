-- drop all table in blog database

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS images CASCADE;
DROP TABLE IF EXISTS tags CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS posts_tags CASCADE;
DROP TABLE IF EXISTS posts_categories CASCADE;
DROP TABLE IF EXISTS posts_images CASCADE;
DROP TABLE IF EXISTS favorites CASCADE;
DROP TABLE IF EXISTS comments CASCADE;

SELECT * FROM users;

CREATE TABLE users
(
    id                 VARCHAR(255)       NOT NULL PRIMARY KEY,
    name               VARCHAR(50)        NOT NULL,
    nickname           VARCHAR(50)        NOT NULL,
    email              VARCHAR(50) UNIQUE NOT NULL,
    password           BYTEA              NOT NULL,
    bio                VARCHAR(50)        DEFAULT '',
    image              VARCHAR(255)       NOT NULL DEFAULT 'user-default-image.png',
    created_at         TIMESTAMP          NOT NULL DEFAULT NOW()
);

CREATE INDEX ON users (nickname);
--
-- CREATE TABLE posts
-- (
--     id           VARCHAR(255)        NOT NULL PRIMARY KEY,
--     author_id    VARCHAR(255)        NOT NULL,
--     title        VARCHAR(100)        NOT NULL,
--     slug         VARCHAR(100) UNIQUE NOT NULL,
--     excerpt      VARCHAR(100)        NOT NULL,
--     content      TEXT                NOT NULL,
--     is_published BOOLEAN             NOT NULL DEFAULT FALSE,
--     published_at TIMESTAMP           NOT NULL DEFAULT NOW(),
--     created_at   TIMESTAMP           NOT NULL DEFAULT NOW(),
--     updated_at   TIMESTAMP           NOT NULL DEFAULT NOW(),
--     FOREIGN KEY (author_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
-- );
--
-- CREATE INDEX ON posts (title);
--
-- CREATE TABLE images
-- (
--     id   VARCHAR(255)        NOT NULL PRIMARY KEY,
--     name VARCHAR(255) UNIQUE NOT NULL,
--     path VARCHAR(255)        NOT NULL
-- );
--
-- CREATE TABLE posts_images
-- (
--     post_id  VARCHAR(255) NOT NULL,
--     image_id VARCHAR(255) NOT NULL,
--     PRIMARY KEY (post_id, image_id),
--     FOREIGN KEY (post_id) REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE,
--     FOREIGN KEY (image_id) REFERENCES images (id) ON UPDATE CASCADE ON DELETE CASCADE
-- );
--
-- CREATE TABLE tags
-- (
--     id   VARCHAR(255)        NOT NULL PRIMARY KEY,
--     name VARCHAR(255) UNIQUE NOT NULL,
--     slug VARCHAR(255) UNIQUE NOT NULL
-- );
--
-- CREATE TABLE posts_tags
-- (
--     post_id VARCHAR(255) NOT NULL,
--     tag_id  VARCHAR(255) NOT NULL,
--     PRIMARY KEY (post_id, tag_id),
--     FOREIGN KEY (post_id) REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE,
--     FOREIGN KEY (tag_id) REFERENCES tags (id) ON UPDATE CASCADE ON DELETE CASCADE
-- );
--
-- CREATE TABLE categories
-- (
--     id   VARCHAR(255)        NOT NULL PRIMARY KEY,
--     name VARCHAR(255) UNIQUE NOT NULL,
--     slug VARCHAR(255) UNIQUE NOT NULL
-- );
--
-- CREATE TABLE posts_categories
-- (
--     post_id     VARCHAR(255) NOT NULL,
--     category_id VARCHAR(255) NOT NULL,
--     PRIMARY KEY (post_id, category_id),
--     FOREIGN KEY (post_id) REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE,
--     FOREIGN KEY (category_id) REFERENCES categories (id) ON UPDATE CASCADE ON DELETE CASCADE
-- );
--
-- CREATE TABLE comments
-- (
--     id         VARCHAR(255) NOT NULL PRIMARY KEY,
--     user_id    VARCHAR(255) NOT NULL,
--     post_id    VARCHAR(255) NOT NULL,
--     content    TEXT         NOT NULL,
--     created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
--     updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
--     FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
--     FOREIGN KEY (post_id) REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE
-- );
--
-- CREATE TABLE favorites
-- (
--     id      VARCHAR(255) NOT NULL PRIMARY KEY,
--     user_id VARCHAR(255) NOT NULL,
--     post_id VARCHAR(255) NOT NULL,
--     FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE,
--     FOREIGN KEY (post_id) REFERENCES posts (id) ON UPDATE CASCADE ON DELETE CASCADE
-- );
--
--
-- -- get all posts, tags, images, and categories where id = '1' with join table
-- SELECT posts.id, posts.title, posts.slug, posts.excerpt, posts.content, posts.is_published, posts.published_at, posts.created_at, posts.updated_at,
--        tags.id, tags.name, tags.slug,
--        images.id, images.name, images.path,
--        categories.id, categories.name, categories.slug
-- FROM posts
-- LEFT JOIN posts_tags ON posts.id = posts_tags.post_id
-- LEFT JOIN tags ON posts_tags.tag_id = tags.id
-- LEFT JOIN posts_categories ON posts.id = posts_categories.post_id
-- LEFT JOIN categories ON posts_categories.category_id = categories.id
-- LEFT JOIN posts_images ON posts.id = posts_images.post_id
-- LEFT JOIN images ON posts_images.image_id = images.id
-- WHERE posts.id = '69ae1e1d-5571-4ffe-94fa-21f71880d649';
--
-- SELECT table_name
-- FROM information_schema.tables
-- WHERE table_schema = 'public';
--
-- SELECT * FROM users;
-- SELECT * FROM tags;
-- SELECT * FROM categories;
-- SELECT * FROM posts;
-- SELECT * FROM images;
-- SELECT * FROM posts_tags;
-- SELECT * FROM posts_categories;
-- SELECT * FROM posts_images;
-- SELECT *
-- FROM posts;
-- SELECT *
-- FROM tags
-- WHERE author_id = 'f8090b56-4a18-4cdc-94cb-a5a09422140c';
--
-- SELECT id,
--        author_id,
--        title,
--        slug,
--        excerpt,
--        content,
--        is_published,
--        published_at,
--        created_at,
--        updated_at
-- FROM posts
-- WHERE author_id = 'f8090b56-4a18-4cdc-94cb-a5a09422140c';
--
-- SELECT * FROM  users;
--
-- DELETE FROM users;
--
-- SELECT count(*) FROM users WHERE email = 'sam';