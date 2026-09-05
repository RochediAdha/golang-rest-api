CREATE TABLE IF NOT EXISTS privileges (
    id uuid NOT NULL,
    code text COLLATE pg_catalog."default" NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    description text COLLATE pg_catalog."default",
    "isActive" boolean NOT NULL DEFAULT true,
    "createdAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" timestamp(3) without time zone NOT NULL,
    "createdBy" uuid,
    "updatedBy" uuid,
    CONSTRAINT privileges_pkey PRIMARY KEY (id),
    CONSTRAINT privileges_code_key UNIQUE (code)
);
