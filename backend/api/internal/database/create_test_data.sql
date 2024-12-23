-- Users
INSERT INTO Users (username, picture, bio, links, creation_date) VALUES
    ('dev_user1', 'https://example.com/dev_user1.jpg', 'Full-stack developer passionate about open-source projects.', '["https://github.com/dev_user1", "https://devuser1.com"]', '2023-12-13 00:00:00'),
    ('tech_writer2', 'https://example.com/tech_writer2.jpg', 'Technical writer and Python enthusiast.', '["https://blog.techwriter.com"]', '2022-12-13 00:00:00'),
    ('data_scientist3', 'https://example.com/data_scientist3.jpg', 'Data scientist with a passion for machine learning.', '["https://github.com/data_scientist3", "https://datascientist3.com"]', '2023-06-13 00:00:00'),
    ('backend_guru4', 'https://example.com/backend_guru4.jpg', 'Backend expert specializing in scalable systems.', '["https://github.com/backend_guru4"]', '2024-01-15 00:00:00'),
    ('ui_designer5', 'https://example.com/ui_designer5.jpg', 'UI/UX designer with a love for user-friendly apps.', '["https://portfolio.uidesigner5.com"]', '2023-05-10 00:00:00');

-- Projects
INSERT INTO Projects (name, description, status, likes, tags, links, owner, creation_date) VALUES
    ('OpenAPI Toolkit', 'A toolkit for generating and testing OpenAPI specs.', 1, 120, '["OpenAPI", "Go", "Tooling"]', '["https://github.com/dev_user1/openapi-toolkit"]', (SELECT id FROM Users WHERE username = 'dev_user1'), '2023-06-13 00:00:00'),
    ('DocuHelper', 'A library for streamlining technical documentation processes.', 2, 85, '["Documentation", "Python"]', '["https://github.com/tech_writer2/docuhelper"]', (SELECT id FROM Users WHERE username = 'tech_writer2'), '2021-12-13 00:00:00'),
    ('ML Research', 'Research repository for various machine learning algorithms.', 1, 45, '["Machine Learning", "Python", "Research"]', '["https://github.com/data_scientist3/ml-research"]', (SELECT id FROM Users WHERE username = 'data_scientist3'), '2024-09-13 00:00:00'),
    ('ScaleDB', 'A scalable database system for modern apps.', 1, 70, '["Database", "Scalability", "Backend"]', '["https://github.com/backend_guru4/scaledb"]', (SELECT id FROM Users WHERE username = 'backend_guru4'), '2024-03-15 00:00:00');

-- Posts
INSERT INTO Posts (content, project_id, creation_date, user_id, likes) VALUES
    ('Excited to release the first version of OpenAPI Toolkit!', (SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), '2024-09-13 00:00:00', (SELECT id FROM Users WHERE username = 'dev_user1'), 40),
    ('We''ve archived DocuHelper, but feel free to explore the code.', (SELECT id FROM Projects WHERE name = 'DocuHelper'), '2024-06-13 00:00:00', (SELECT id FROM Users WHERE username = 'tech_writer2'), 25),
    ('Updated ML Research repo with new algorithms for data analysis.', (SELECT id FROM Projects WHERE name = 'ML Research'), '2024-11-13 00:00:00', (SELECT id FROM Users WHERE username = 'data_scientist3'), 15);

-- Comments on Projects
INSERT INTO Comments (content, parent_comment_id, likes, creation_date, user_id, likes) VALUES
    ('This is a fantastic project! Can''t wait to contribute.', NULL, 5, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'dev_user1'), 0),
    ('I love the concept, but I think the documentation could be improved.', NULL, 3, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'tech_writer2'), 0),
    ('Great to see more open-source tools for API development!', NULL, 4, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'backend_guru4'), 0),
    ('I agree, but the API specs seem a bit too complex for beginners.', (SELECT id FROM Comments WHERE content = 'Great to see more open-source tools for API development!'), 2, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'data_scientist3'), 0),
    ('I hope this toolkit will integrate with other Go tools soon!', (SELECT id FROM Comments WHERE content = 'This is a fantastic project! Can''t wait to contribute.'), 1, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'ui_designer5'), 0);

-- ProjectComments relations (Mapping comments to projects)
INSERT INTO ProjectComments (project_id, comment_id, user_id) VALUES
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Comments WHERE content = 'This is a fantastic project! Can''t wait to contribute.'), (SELECT id FROM Users WHERE username = 'dev_user1')),
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Comments WHERE content = 'I love the concept, but I think the documentation could be improved.'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Comments WHERE content = 'Great to see more open-source tools for API development!'), (SELECT id FROM Users WHERE username = 'backend_guru4')),
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Comments WHERE content = 'I agree, but the API specs seem a bit too complex for beginners.'), (SELECT id FROM Users WHERE username = 'data_scientist3')),
    ((SELECT id FROM Projects WHERE name = 'OpenAPI Toolkit'), (SELECT id FROM Comments WHERE content = 'I hope this toolkit will integrate with other Go tools soon!'), (SELECT id FROM Users WHERE username = 'ui_designer5'));

-- Comments on Posts
INSERT INTO Comments (content, parent_comment_id, likes, creation_date, user_id) VALUES
    ('Awesome update! I''ll try it out.', NULL, 2, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'backend_guru4')),
    ('Thanks for sharing! Will this feature be extended soon?', NULL, 1, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'data_scientist3')),
    ('Great work, looking forward to more updates!', NULL, 4, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'ui_designer5')),
    ('Will this be compatible with earlier versions of OpenAPI?', (SELECT id FROM Comments WHERE content = 'Thanks for sharing! Will this feature be extended soon?'), 1, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ('I hope the next update addresses performance improvements.', (SELECT id FROM Comments WHERE content = 'Awesome update! I''ll try it out.'), 3, '2024-12-23 00:00:00', (SELECT id FROM Users WHERE username = 'data_scientist3'));

-- PostComments relations (Mapping comments to posts)
INSERT INTO PostComments (post_id, comment_id, user_id) VALUES
    ((SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Comments WHERE content = 'Awesome update! I''ll try it out.'), (SELECT id FROM Users WHERE username = 'backend_guru4')),
    ((SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Comments WHERE content = 'Thanks for sharing! Will this feature be extended soon?'), (SELECT id FROM Users WHERE username = 'data_scientist3')),
    ((SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Comments WHERE content = 'Great work, looking forward to more updates!'), (SELECT id FROM Users WHERE username = 'ui_designer5')),
    ((SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Comments WHERE content = 'Will this be compatible with earlier versions of OpenAPI?'), (SELECT id FROM Users WHERE username = 'tech_writer2')),
    ((SELECT id FROM Posts WHERE content = 'Excited to release the first version of OpenAPI Toolkit!'), (SELECT id FROM Comments WHERE content = 'I hope the next update addresses performance improvements.'), (SELECT id FROM Users WHERE username = 'data_scientist3'));
