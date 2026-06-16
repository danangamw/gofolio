-- Create "sys_configs" table
CREATE TABLE "sys_configs" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "key" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "value" text NOT NULL,
  "value_type" character varying(20) NOT NULL DEFAULT 'text',
  "group_name" character varying(100) NOT NULL DEFAULT 'general',
  "sort_order" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_sys_configs_key" to table: "sys_configs"
CREATE UNIQUE INDEX "idx_sys_configs_key" ON "sys_configs" ("key");
