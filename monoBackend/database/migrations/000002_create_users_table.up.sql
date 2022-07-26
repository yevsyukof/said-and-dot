CREATE TABLE IF NOT EXISTS "Users"
(
    "id"            uuid                     PRIMARY KEY DEFAULT uuid_generate_v4(),
    "username"      varchar(60)              NOT NULL UNIQUE CHECK ( username != '' ),
    "password_hash" varchar(128)             NOT NULL CHECK ( password_hash != '' ),
    "first_name"    varchar(40)              NOT NULL CHECK ( first_name != '' ),
    "last_name"     varchar(40)              NOT NULL CHECK ( last_name != '' ),
    "email"         varchar(60)              NOT NULL UNIQUE CHECK ( email != '' ),
    "created"       timestamp with time zone NOT NULL,
    "updated"       timestamp with time zone NOT NULL CHECK ( updated > created )
)