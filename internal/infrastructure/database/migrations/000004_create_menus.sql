DO $$ BEGIN
    CREATE TYPE "MenuType" AS ENUM ('ITEM', 'GROUP', 'COLLAPSE');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS menus (
    id uuid NOT NULL,
    "parentId" uuid,
    code text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    path text COLLATE pg_catalog."default",
    icon text COLLATE pg_catalog."default",
    description text COLLATE pg_catalog."default",
    "sortOrder" integer NOT NULL DEFAULT 0,
    type "MenuType" NOT NULL DEFAULT 'ITEM'::"MenuType",
    "isActive" boolean NOT NULL DEFAULT true,
    "createdAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "deletedAt" timestamp(3) without time zone NULL,
    "createdBy" uuid,
    "updatedBy" uuid,
    "deletedBy" uuid,
    CONSTRAINT menus_pkey PRIMARY KEY (id),
    CONSTRAINT "menus_parentId_fkey" FOREIGN KEY ("parentId")
        REFERENCES public.menus (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
