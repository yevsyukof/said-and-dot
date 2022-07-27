CREATE TABLE IF NOT EXISTS "Refresh_tokens"
(
    "user_id" uuid        NOT NULL,
    "token"   varchar(256) NOT NULL CHECK ( token != '' ),
    CONSTRAINT "Refresh_tokens_user_id_token_unique" UNIQUE ("user_id", "token")
)