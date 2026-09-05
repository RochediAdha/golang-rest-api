CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL,
    username text COLLATE pg_catalog."default" NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    "isActive" boolean NOT NULL DEFAULT true,
    "createdAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "createdBy" uuid,
    "updatedBy" uuid,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
