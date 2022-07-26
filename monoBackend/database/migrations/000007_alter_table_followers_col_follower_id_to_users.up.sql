ALTER TABLE "Followers"
    ADD CONSTRAINT "Followers_Users_fk0" FOREIGN KEY ("follower_id") REFERENCES "Users" ("id")
        ON DELETE CASCADE ON UPDATE CASCADE;
