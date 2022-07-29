package db_middleware

import (
	"github.com/google/uuid"
	"said-and-dot-backend/internal/api/entities"
	api_errors "said-and-dot-backend/internal/api/routes/errors"
	"said-and-dot-backend/internal/database"
)

func GetUserByID(id string, db database.Database) (*entities.User, error) {
	var userExists bool
	if err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE id = $1)", id).Scan(&userExists); err != nil {
		return nil, err
	} else if !userExists {
		return nil, api_errors.ErrUserDoesNotExist
	}

	requiredUser := new(entities.User)

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	requiredUser.ID = parsedID

	if err := db.QueryRow(
		"SELECT username, first_name, last_name, email FROM Users WHERE id = $1", requiredUser.ID).Scan(
		&requiredUser.Username, &requiredUser.FirstName, &requiredUser.LastName, &requiredUser.Email); err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	return requiredUser, nil
}

func GetUserTweets(user *entities.User, db database.Database) (*[]entities.Tweet, error) {
	var foundTweetsRecords database.Rows

	userTweets := new([]entities.Tweet)

	foundTweetsRecords, err := db.Query(
		"SELECT id, tweet, created FROM Tweets WHERE user_id = $1", user.ID)
	if err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	defer foundTweetsRecords.Close()

	for foundTweetsRecords.Next() {
		var tweet entities.Tweet

		if err := foundTweetsRecords.Scan(&tweet.ID, &tweet.Tweet, &tweet.Created); err != nil {
			return nil, api_errors.ErrDatabaseError
		}
		tweet.UserID = user.ID

		tweetFavorites, err := GetTweetFavorites(tweet.ID, db)
		if err != nil {
			return nil, err
		}

		for _, favorite := range *tweetFavorites {
			*tweet.Likes = append(*tweet.Likes, favorite)
		}

		*userTweets = append(*userTweets, tweet)
	}

	//sort.Slice(userTweets, func(i, j int) (less bool) {
	//	return (*userTweets)[i].Created.Before((*userTweets)[j].Created)
	//})

	return userTweets, nil
}

func GetTweetFavorites(tweetID uuid.UUID, db database.Database) (*[]uuid.UUID, error) {
	tweetFavorites := new([]uuid.UUID)

	var foundTweetFavorites database.Rows

	foundTweetFavorites, err := db.Query(
		"SELECT user_id FROM Favorites WHERE tweet_id = $1", tweetID)
	if err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	defer foundTweetFavorites.Close()

	for foundTweetFavorites.Next() {
		var favorite uuid.UUID

		if err := foundTweetFavorites.Scan(&favorite); err != nil {
			return nil, api_errors.ErrDatabaseError
		}
		*tweetFavorites = append(*tweetFavorites, favorite)
	}

	return tweetFavorites, nil
}

func GetAllTweets(db database.Database) (*[]entities.TweetWithAuthor, error) {
	var tweetsRecords database.Rows

	allTweets := new([]entities.TweetWithAuthor)

	tweetsRecords, err := db.Query(
		"SELECT id, user_id, tweet, created FROM Tweets")
	if err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	defer tweetsRecords.Close()

	for tweetsRecords.Next() {
		curTweet := new(entities.Tweet)

		if err := tweetsRecords.Scan(&curTweet.ID, &curTweet.UserID, &curTweet.Tweet, &curTweet.Created); err != nil {
			return nil, api_errors.ErrDatabaseError
		}

		curTweetAuthor, err := GetUserByID(curTweet.UserID.String(), db)
		if err != nil {
			return nil, api_errors.ErrDatabaseError
		}

		*allTweets = append(*allTweets, entities.TweetWithAuthor{
			Tweet:  curTweet,
			Author: curTweetAuthor,
		})
	}
	//sort.Slice(allTweets, func(i, j int) (less bool) {
	//	return (*allTweets)[i].Created.Before((*allTweets)[j].Created)
	//})

	return allTweets, nil
}

func GetFollowersByUserID(userID uuid.UUID, db database.Database) (*[]uuid.UUID, error) {
	userFollowers := new([]uuid.UUID)

	var foundUserFollowers database.Rows

	foundUserFollowers, err := db.Query(
		"SELECT follower_id FROM followers WHERE user_id = $1", userID)
	if err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	defer foundUserFollowers.Close()

	for foundUserFollowers.Next() {
		var follower uuid.UUID

		if err := foundUserFollowers.Scan(&follower); err != nil {
			return nil, api_errors.ErrDatabaseError
		}
		*userFollowers = append(*userFollowers, follower)
	}

	return userFollowers, nil
}

func GetFollowsByUserID(userID uuid.UUID, db database.Database) (*[]uuid.UUID, error) {
	userFollows := new([]uuid.UUID)

	var foundUserFollows database.Rows

	foundUserFollows, err := db.Query(
		"SELECT user_id FROM followers WHERE follower_id = $1", userID)
	if err != nil {
		return nil, api_errors.ErrDatabaseError
	}
	defer foundUserFollows.Close()

	for foundUserFollows.Next() {
		var followedUser uuid.UUID

		if err := foundUserFollows.Scan(&followedUser); err != nil {
			return nil, api_errors.ErrDatabaseError
		}
		*userFollows = append(*userFollows, followedUser)
	}

	return userFollows, nil
}

func SaveNewTweet(tweet *entities.Tweet, db database.Database) error {
	if _, err := db.Exec(
		"INSERT INTO tweets (user_id, tweet, created)  VALUES ($1, $2, $3)",
		tweet.UserID, tweet.Tweet, tweet.Created); err != nil {
		return err
	}
	return nil
}
