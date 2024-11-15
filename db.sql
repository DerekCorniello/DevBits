-- Prelimenary db
-- User Table
CREATE TABLE Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    profile_pic TEXT,
    bio TEXT,
    links TEXT[]
);

-- Project Table
CREATE TABLE Projects (
    project_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50),
    links TEXT[],
    tags TEXT[],
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Post Table
CREATE TABLE Posts (
    post_id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    project_id INT REFERENCES Projects(project_id) ON DELETE CASCADE,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Comment Table
CREATE TABLE Comments (
    comment_id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    post_id INT REFERENCES Posts(post_id) ON DELETE CASCADE,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE
);

-- Likes for Projects
CREATE TABLE ProjectLikes (
    project_id INT REFERENCES Projects(project_id) ON DELETE CASCADE,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, user_id)
);

-- Likes for Posts
CREATE TABLE PostLikes (
    post_id INT REFERENCES Posts(post_id) ON DELETE CASCADE,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, user_id)
);

-- Likes for Comments
CREATE TABLE CommentLikes (
    comment_id INT REFERENCES Comments(comment_id) ON DELETE CASCADE,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (comment_id, user_id)
);

-- Follows between Users (User Following)
CREATE TABLE UserFollows (
    follower_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    followed_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (follower_id, followed_id),
    CHECK (follower_id != followed_id) -- Prevents a user from following themselves
);

-- Follows for Projects (User Following a Project)
CREATE TABLE ProjectFollows (
    project_id INT REFERENCES Projects(project_id) ON DELETE CASCADE,
    user_id INT REFERENCES Users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, user_id)
);
