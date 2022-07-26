ALTER TABLE "Followers"
    ADD CONSTRAINT "Followers_Users_fk1" FOREIGN KEY ("user_id") REFERENCES "Users" ("id")
        ON DELETE CASCADE ON UPDATE CASCADE;
