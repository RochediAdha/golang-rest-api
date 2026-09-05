CREATE TABLE IF NOT EXISTS role_privileges (
    id uuid NOT NULL,
    "roleId" uuid NOT NULL,
    "menuId" uuid NOT NULL,
    "privilegeId" uuid NOT NULL,
    "createdAt" timestamp(3) without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "createdBy" uuid,
    CONSTRAINT role_privileges_pkey PRIMARY KEY (id),
    CONSTRAINT "role_privileges_roleId_menuId_privilegeId_key" UNIQUE ("roleId", "menuId", "privilegeId"),
    CONSTRAINT "role_privileges_menuId_fkey" FOREIGN KEY ("menuId")
        REFERENCES public.menus (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT "role_privileges_privilegeId_fkey" FOREIGN KEY ("privilegeId")
        REFERENCES public.privileges (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    CONSTRAINT "role_privileges_roleId_fkey" FOREIGN KEY ("roleId")
        REFERENCES public.roles (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
);
