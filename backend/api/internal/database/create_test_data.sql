-- Insert Users
INSERT INTO Users (username, profile_pic, bio, links, creation_date) VALUES
    ('dev_user1', 'https://example.com/dev_user1.jpg', 'Full-stack developer passionate about open-source projects.', '["https://github.com/dev_user1", "https://devuser1.com"]', datetime('now', '-1 year')),
    ('tech_writer2', 'https://example.com/tech_writer2.jpg', 'Technical writer and Python enthusiast.', '["https://blog.techwriter.com"]', datetime('now', '-2 years')),
    ('data_scientist3', 'https://example.com/data_scientist3.jpg', 'Data scientist with a passion for machine learning.', '["https://github.com/data_scientist3", "https://datascientist3.com"]', datetime('now', '-1.5 years'));

-- Insert Projects
INSERT INTO Projects (name, description, status, likes, tags, links, owner, creation_date) VALUES
    ('OpenAPI Toolkit', 'A toolkit for generating and testing OpenAPI specs.', 1, 120, '["OpenAPI", "Go", "Tooling"]', '["https://github.com/dev_user1/openapi-toolkit"]', (SELECT id FROM Users WHERE username = 'dev_user1'), datetime('now', '-1.5 years')),
    ('DocuHelper', 'A library for streamlining technical documentation processes.', 2, 85, '["Documentation", "Python"]', '["https://github.com/tech_writer2/docuhelper"]', (SELECT id FROM Users WHERE username = 'tech_writer2'), datetime('now', '-3 years')),
    ('ML Research', 'Research repository for various machine learning algorithms.', 1, 45, '["Machine Learning", "Python", "Research"]', '["https://github.com/data_scientist3/ml-research"]', (SELECT id FROM Users WHERE username = 'data_scientist3'), datetime('now', '-3 months'));

-- Insert Posts
INSERT INTO Posts (content, project_id, creation_date, user_id, likes) VALUES
    ('Excited to release the first version of OpenAPI Toolkit!', (SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), datetime('now', '-3 months'), (SELECT id FROM Users WHERE username = 'dev_user1'), 40),
    ('We''ve archived DocuHelper, but feel free to explore the code.', (SELECT id FROM Projects WHERE name = 'DocuHelper'), datetime('now', '-6 months'), (SELECT id FROM Users WHERE username = 'tech_writer2'), 25),
    ('Updated ML Research repo with new algorithms for data analysis.', (SELECT id FROM Projects WHERE name = 'ML Research'), datetime('now', '-1 month'), (SELECT id FROM Users WHERE username = 'data_scientist3'), 15);

-- Insert Comments
INSERT INTO Comments (content, post_id, parent_comment_id, creation_date, user_id) VALUES
    ('This is amazing! Can''t wait to try it out.', (SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), NULL, datetime('now', '-2 months'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ('Thanks for the kind words!', (SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Comments WHERE content = 'This is amazing! Can''t wait to try it out.'), datetime('now', '-1 month'), (SELECT id FROM Users WHERE username = 'dev_user1')),
    ('Looks great! I''ll test it and report back.', (SELECT id FROM Posts WHERE content = 'Updated ML Research repo with new algorithms for data analysis.'), NULL, datetime('now', '-1 month'), (SELECT id FROM Users WHERE username = 'data_scientist3'));

-- Insert Project Likes
INSERT INTO ProjectLikes (project_id, user_id) VALUES
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Users WHERE username = 'dev_user1')),
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ((SELECT id FROM Projects WHERE name = 'DocuHelper'), (SELECT id FROM Users WHERE username = 'dev_user1')),
    ((SELECT id FROM Projects WHERE name = 'ML Research'), (SELECT id FROM Users WHERE username = 'data_scientist3'));

-- Insert Post Likes
INSERT INTO PostLikes (post_id, user_id) VALUES
    ((SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Users WHERE username = 'dev_user1')),
    ((SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ((SELECT id FROM Posts WHERE content = 'We''ve archived DocuHelper, but feel free to explore the code.'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ((SELECT id FROM Posts WHERE content = 'Updated ML Research repo with new algorithms for data analysis.'), (SELECT id FROM Users WHERE username = 'data_scientist3'));

-- Insert Comment Likes
INSERT INTO CommentLikes (comment_id, user_id) VALUES
    ((SELECT id FROM Comments WHERE content = 'This is amazing! Can''t wait to try it out.'), (SELECT id FROM Users WHERE username = 'dev_user1')),
    ((SELECT id FROM Comments WHERE content = 'Thanks for the kind words!'), (SELECT id FROM Users WHERE username = 'tech_writer2'));

-- Insert User Follows
INSERT INTO UserFollows (follower_id, follows_id) VALUES
    ((SELECT id FROM Users WHERE username = 'dev_user1'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ((SELECT id FROM Users WHERE username = 'tech_writer2'), (SELECT id FROM Users WHERE username = 'data_scientist3')),
    ((SELECT id FROM Users WHERE username = 'dev_user1'), (SELECT id FROM Users WHERE username = 'data_scientist3')),
    ((SELECT id FROM Users WHERE username = 'data_scientist3'), (SELECT id FROM Users WHERE username = 'dev_user1'));

-- Insert Project Follows
INSERT INTO ProjectFollows (project_id, user_id) VALUES
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ((SELECT id FROM Projects WHERE name = 'ML Research'), (SELECT id FROM Users WHERE username = 'dev_user1'));
