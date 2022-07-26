CREATE TABLE IF NOT EXISTS "Followers"
(
    "follower_id" uuid NOT NULL CHECK (user_id != follower_id),
    "user_id"     uuid NOT NULL,
    CONSTRAINT "Followers_user_id_follower_id_unique" UNIQUE ("user_id", "follower_id")
)