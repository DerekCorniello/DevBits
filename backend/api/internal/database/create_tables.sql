-- Drop tables if they already exist
DROP TABLE IF EXISTS UserLoginInfo;

DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS UserFollows;

DROP TABLE IF EXISTS Projects;
DROP TABLE IF EXISTS ProjectLikes;
DROP TABLE IF EXISTS ProjectFollows;
DROP TABLE IF EXISTS ProjectComments;

DROP TABLE IF EXISTS Posts;
DROP TABLE IF EXISTS PostLikes;
DROP TABLE IF EXISTS PostComments;

DROP TABLE IF EXISTS Comments;
DROP TABLE IF EXISTS CommentLikes;

-- UserLoginInfo
CREATE TABLE UserLoginInfo (
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    PRIMARY KEY(username)
);

-- Users Table
CREATE TABLE Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    picture TEXT,
    bio TEXT,
    links JSON,
    creation_date TIMESTAMP NOT NULL
);

-- Projects Table
CREATE TABLE Projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status INTEGER,
    likes INTEGER DEFAULT 0,
    links JSON,
    tags JSON,
    owner INTEGER NOT NULL,
    creation_date TIMESTAMP NOT NULL,
    FOREIGN KEY (owner) REFERENCES Users(id) ON DELETE CASCADE
);

-- Project Comments Table (Normalizing comments relationship)
CREATE TABLE ProjectComments (
    project_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES Comments(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, comment_id)
);

-- Posts Table
CREATE TABLE Posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    creation_date TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL,
    likes INTEGER DEFAULT 0,
    FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Project Comments Table (Normalizing comments relationship)
CREATE TABLE PostComments (
    post_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES Comments(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, comment_id)
);

-- Comments Table
CREATE TABLE Comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    parent_comment_id INTEGER,
    creation_date TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES Comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Likes for Projects
CREATE TABLE ProjectLikes (
    project_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (project_id, user_id),
    FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Likes for Posts
CREATE TABLE PostLikes (
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES Posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Likes for Comments
CREATE TABLE CommentLikes (
    comment_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (comment_id, user_id),
    FOREIGN KEY (comment_id) REFERENCES Comments(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);

-- Follows between Users (User Following)
CREATE TABLE UserFollows (
    follower_id INTEGER NOT NULL,
    follows_id INTEGER NOT NULL,
    PRIMARY KEY (follower_id, follows_id),
    FOREIGN KEY (follower_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (follows_id) REFERENCES Users(id) ON DELETE CASCADE,
    CHECK (follower_id != follows_id)
);

-- Follows for Projects (User Following a Project)
CREATE TABLE ProjectFollows (
    project_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    PRIMARY KEY (project_id, user_id),
    FOREIGN KEY (project_id) REFERENCES Projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
);
