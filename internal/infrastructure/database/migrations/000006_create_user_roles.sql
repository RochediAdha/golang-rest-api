CREATE TABLE IF NOT EXISTS user_roles (
    id uuid NOT NULL,
    "userId" uuid NOT NULL,
    "roleId" uuid NOT NULL,
    number text COLLATE pg_catalog."default",
    "createdAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "createdBy" uuid,
    CONSTRAINT user_roles_pkey PRIMARY KEY (id),
    CONSTRAINT "user_roles_userId_roleId_key" UNIQUE ("userId", "roleId"),
    CONSTRAINT "user_roles_roleId_fkey" FOREIGN KEY ("roleId")
        REFERENCES public.roles (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT "user_roles_userId_fkey" FOREIGN KEY ("userId")
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS "user_roles_roleId_number_key"
    ON user_roles ("roleId", number)
    WHERE number IS NOT NULL;
