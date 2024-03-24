CREATE TABLE IF NOT EXISTS projects
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS projects_id_idx ON projects(id);



CREATE TABLE IF NOT EXISTS goods
(
    id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    name VARCHAR(255) NOT NULL,
    description TEXT default '',
    priority INTEGER default 1,
    removed BOOLEAN default false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS goods_id_idx ON goods(id);
CREATE INDEX IF NOT EXISTS goods_project_id_idx ON goods(project_id);
CREATE INDEX IF NOT EXISTS goods_name_idx ON goods(name);

-- CREATE TABLE IF NOT EXISTS post (
--                       id int NOT NULL,
--                       title text,
--                       body text,
--                       PRIMARY KEY(id)
-- );