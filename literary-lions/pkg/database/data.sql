CREATE TABLE "users" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "username" TEXT NOT NULL,
  "email" TEXT NOT NULL UNIQUE,
  "password" TEXT NOT NULL
);

CREATE TABLE "posts" (
  "id" INTEGER PRIMARY KEY,
  "title" TEXT NOT NULL,
  "topic" TEXT NOT NULL,
  "content" TEXT NOT NULL,
  "user_id" INTEGER,
  "username" TEXT,
  "views" INTEGER DEFAULT 0, 
  "date" DATETIME DEFAULT CURRENT_TIMESTAMP,
  "edited_at" DATETIME,
  "likes" INTEGER DEFAULT 0 CHECK (likes >= 0),
  "dislikes" INTEGER DEFAULT 0 CHECK (dislikes >= 0),
  FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

CREATE TABLE "comments" (
  "id" INTEGER PRIMARY KEY,
  "content" TEXT NOT NULL,
  "created_at" DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  "edited_at" DATETIME,
  "post_id" INTEGER,
  "comment_id" INTEGER,
  "user_id" INTEGER NOT NULL,
  "post_title" TEXT NOT NULL,
  "likes" INTEGER DEFAULT 0 CHECK (likes >= 0),
  "dislikes" INTEGER DEFAULT 0 CHECK (dislikes >= 0),
  FOREIGN KEY("post_id") REFERENCES "posts"("id") ON DELETE CASCADE,
  FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE
);

CREATE TABLE "comment_likes" (
  "id" INTEGER PRIMARY KEY,
  "like_type" TEXT CHECK ("like_type" IN ('like', 'dislike')) NOT NULL,
  "user_id" INTEGER NOT NULL,
  "comment_id" INTEGER NOT NULL,
  FOREIGN KEY("user_id") REFERENCES "users"("id"),
  FOREIGN KEY("comment_id") REFERENCES "comments"("id"),
  UNIQUE("user_id", "comment_id")
);

CREATE TABLE "post_likes" (
  "id" INTEGER PRIMARY KEY,
  "like_type" TEXT CHECK ("like_type" IN ('like', 'dislike')) NOT NULL,
  "user_id" INTEGER NOT NULL,
  "post_id" INTEGER NOT NULL,
  FOREIGN KEY("user_id") REFERENCES "users"("id"),
  FOREIGN KEY("post_id") REFERENCES "posts"("id"),
  UNIQUE("user_id", "post_id")
);
