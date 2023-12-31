ALTER TABLE movies ADD CONSTRAINT fk_genre_id FOREIGN KEY (genre_id) REFERENCES genres(id);
ALTER TABLE movie_images ADD CONSTRAINT fk_movie_id FOREIGN KEY (movie_id) REFERENCES movies(id);
ALTER TABLE users ADD CONSTRAINT fk_profile_id FOREIGN KEY (profile_id) REFERENCES profiles(id);