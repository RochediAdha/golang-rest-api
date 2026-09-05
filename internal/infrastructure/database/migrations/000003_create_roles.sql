CREATE TABLE IF NOT EXISTS roles (
    id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    "isActive" boolean NOT NULL DEFAULT true,
    "createdAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "deletedAt" timestamp(3) without time zone NULL,
    "createdBy" uuid,
    "updatedBy" uuid,
    "deletedBy" uuid,
    CONSTRAINT roles_pkey PRIMARY KEY (id)
);
