ALTER TABLE Favorites
    ADD CONSTRAINT "Favorites_Users_fk0" FOREIGN KEY (user_id) REFERENCES Users (id)
        ON DELETE CASCADE ON UPDATE CASCADE;
