CREATE TABLE IF NOT EXISTS Tweets
(
    id      uuid                     PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid                     NOT NULL,
    tweet   varchar(512)             NOT NULL CHECK ( tweet != '' ),
    likes   bigint                   NOT NULL DEFAULT '0' CHECK ( likes >= 0 ),
    replies bigint                   NOT NULL DEFAULT '0' CHECK ( replies >= 0 ),
    created timestamp with time zone NOT NULL
)
