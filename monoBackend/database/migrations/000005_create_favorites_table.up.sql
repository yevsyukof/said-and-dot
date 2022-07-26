CREATE TABLE IF NOT EXISTS "Favorites"
(
    "user_id"  uuid NOT NULL,
    "tweet_id" uuid NOT NULL,
    CONSTRAINT "Favorites_user_id_tweet_id_unique" UNIQUE ("user_id", "tweet_id")
)
