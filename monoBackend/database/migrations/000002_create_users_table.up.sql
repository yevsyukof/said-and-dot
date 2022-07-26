CREATE TABLE IF NOT EXISTS "Users"
(
    "id"         uuid                     PRIMARY KEY DEFAULT uuid_generate_v4(),
    "username"   varchar(80)              NOT NULL UNIQUE CHECK ( username != '' ),
    "first_name" varchar(80)              NOT NULL CHECK ( first_name != '' ),
    "last_name"  varchar(80)              NOT NULL CHECK ( last_name != '' ),
    "email"      varchar(128)             NOT NULL UNIQUE CHECK ( email != '' ),
    "created"    timestamp with time zone NOT NULL,
    "updated"    timestamp with time zone NOT NULL CHECK ( updated > created )
)