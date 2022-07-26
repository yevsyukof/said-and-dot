ALTER TABLE "Favorites"
    ADD CONSTRAINT "Favorites_Tweets_fk0" FOREIGN KEY ("tweet_id") REFERENCES "Tweets" ("id")
        ON DELETE CASCADE ON UPDATE CASCADE;
