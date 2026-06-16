-- Modify "blogs" table
ALTER TABLE "blogs" ADD COLUMN "category" character varying(100) NOT NULL DEFAULT '';
-- Modify "portfolios" table
ALTER TABLE "portfolios" ADD COLUMN "icon" character varying(50) NOT NULL DEFAULT '';
