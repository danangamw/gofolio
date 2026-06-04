-- Create "blogs" table
CREATE TABLE "blogs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "title" character varying(255) NOT NULL,
  "slug" character varying(255) NOT NULL,
  "content" text NOT NULL,
  "excerpt" text NULL,
  "status" character varying(20) NOT NULL DEFAULT 'draft',
  "published_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_blogs_slug" to table: "blogs"
CREATE UNIQUE INDEX "idx_blogs_slug" ON "blogs" ("slug");
-- Create "portfolios" table
CREATE TABLE "portfolios" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "title" character varying(255) NOT NULL,
  "description" text NULL,
  "image_url" character varying(500) NULL,
  "tech_stack" text[] NULL,
  "project_url" character varying(500) NULL,
  "repository_url" character varying(500) NULL,
  "sort_order" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "username" character varying(100) NOT NULL,
  "password_hash" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- Create "sessions" table
CREATE TABLE "sessions" (
  "id" character varying(128) NOT NULL,
  "user_id" uuid NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "last_active_at" timestamptz NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_sessions" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
