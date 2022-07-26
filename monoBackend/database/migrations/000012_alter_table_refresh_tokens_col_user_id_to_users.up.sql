ALTER TABLE "Refresh_tokens"
    ADD CONSTRAINT "Refresh_tokens_Users_fk0" FOREIGN KEY ("user_id") REFERENCES "Users" ("id")
        ON DELETE CASCADE ON UPDATE CASCADE;
