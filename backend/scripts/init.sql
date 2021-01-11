DROP TABLE IF EXISTS task_labels CASCADE;
DROP TABLE IF EXISTS labels CASCADE;
DROP TABLE IF EXISTS tasks CASCADE;
DROP TABLE IF EXISTS task_lists CASCADE;
DROP TABLE IF EXISTS board_users CASCADE;
DROP TABLE IF EXISTS boards CASCADE;
DROP TABLE IF EXISTS project_users CASCADE;
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS datetimes CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tokens CASCADE;
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    nickname varchar(32) UNIQUE NOT NULL,
    firstname varchar(32) NOT NULL DEFAULT '',
    lastname varchar(32) NOT NULL DEFAULT '',
    email varchar(32) NOT NULL,
    phone varchar(32) NOT NULL DEFAULT '',
    password varchar(32) NOT NULL,
    avatar varchar(100) NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS permissions (
    id serial PRIMARY KEY,
    read boolean NOT NULL DEFAULT false,
    write boolean NOT NULL DEFAULT false,
    admin boolean NOT NULL DEFAULT false
);
CREATE TABLE IF NOT EXISTS datetimes (
    id serial PRIMARY KEY,
    created bigint NOT NULL,
    updated bigint NOT NULL,
    accessed bigint NOT NULL
);
CREATE TABLE IF NOT EXISTS projects (
    id serial PRIMARY KEY,
    owner_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    default_permissions_id int REFERENCES permissions (id) ON DELETE CASCADE NOT NULL,
    datetimes_id int REFERENCES datetimes (id) ON DELETE CASCADE NOT NULL,
    title varchar(50) NOT NULL,
    description text
);
CREATE TABLE IF NOT EXISTS project_users (
    id serial PRIMARY KEY,
    user_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    project_id int REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    permissions_id int REFERENCES permissions (id) ON DELETE CASCADE NOT NULL
);
CREATE TABLE IF NOT EXISTS boards (
    id serial PRIMARY KEY,
    project_id int REFERENCES projects (id) ON DELETE CASCADE NOT NULL,
    owner_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    default_permissions_id int REFERENCES permissions (id) ON DELETE CASCADE NOT NULL,
    datetimes_id int references datetimes (id) ON DELETE CASCADE NOT NULL,
    title varchar(50) NOT NULL
);
CREATE TABLE IF NOT EXISTS board_users (
    id serial PRIMARY KEY,
    user_id int REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    board_id int REFERENCES boards (id) ON DELETE CASCADE NOT NULL,
    permissions_id int REFERENCES permissions (id) ON DELETE CASCADE NOT NULL
);
CREATE TABLE IF NOT EXISTS task_lists (
    id serial PRIMARY KEY,
    board_id int REFERENCES boards (id) ON DELETE CASCADE NOT NULL,
    title varchar(30) NOT NULL,
    position int NOT NULL
);
CREATE TABLE IF NOT EXISTS tasks (
    id serial PRIMARY KEY,
    list_id int REFERENCES task_lists (id) ON DELETE CASCADE NOT NULL,
    title varchar(30) NOT NULL,
    description text,
    datetimes_id int REFERENCES datetimes (id) ON DELETE CASCADE NOT NULL,
    position smallint NOT NULL
);
CREATE TABLE IF NOT EXISTS tokens (
    id serial PRIMARY KEY,
    jwt text NOT NULL
);
CREATE TABLE IF NOT EXISTS labels (
    id serial PRIMARY KEY,
    board_id int REFERENCES boards (id) ON DELETE CASCADE NOT NULL,
    name varchar(30) NOT NULL,
    color int NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS task_labels (
    id serial PRIMARY KEY,
    task_id int REFERENCES tasks (id) ON DELETE CASCADE NOT NULL,
    label_id int REFERENCES labels (id) ON DELETE CASCADE NOT NULL
);
-- USERS
-- 1
INSERT INTO users (nickname, email, avatar, password)
VALUES ('alex', 'alex@mail.ru', '', 'qwerty');
-- 2
INSERT INTO users (nickname, email, avatar, password)
VALUES ('test_user', 'test_user@mail.ru', '', 'qwerty');
-- 3
INSERT INTO users (nickname, email, avatar, password)
VALUES ('nick1', 'nick1@mail.ru', '', 'qwerty');
-- PROJECTS id=1
-- 1
INSERT INTO permissions (read, write, admin)
VALUES (true, true, true);
-- 2
INSERT INTO permissions (read, write, admin)
VALUES (true, true, false);
-- 1
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 1
INSERT INTO projects (
        owner_id,
        default_permissions_id,
        datetimes_id,
        title,
        description
    )
VALUES (
        1,
        2,
        1,
        'First project',
        'This is the first project'
    );
-- 1
INSERT INTO project_users (user_id, project_id, permissions_id)
VALUES (1, 1, 1);
-- BOARD id=1
-- 3
INSERT INTO permissions (read, write, admin)
VALUES (true, true, true);
-- 4
INSERT INTO permissions (read, write, admin)
VALUES (true, true, false);
-- 2
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 1
INSERT INTO boards (
        project_id,
        owner_id,
        default_permissions_id,
        datetimes_id,
        title
    )
VALUES (1, 1, 4, 2, 'First board');
-- 1
INSERT INTO board_users (user_id, board_id, permissions_id)
VALUES (1, 1, 3);
--
-- 5
-- INSERT INTO permissions (read, write, admin)
-- VALUES (true, true, true);
-- -- 3
-- INSERT INTO datetimes (created, updated, accessed)
-- VALUES (1605925262, 1605925262, 1605925262);
-- -- 2
-- INSERT INTO project_users (user_id, project_id, permissions_id)
-- VALUES (2, 1, 5);
-- LISTS
-- 1
INSERT INTO task_lists (board_id, title, position)
VALUES (1, 'Not Stated', 0);
-- 2
INSERT INTO task_lists (board_id, title, position)
VALUES (1, 'SSSSSSSSS', 1);
-- 3
INSERT INTO task_lists (board_id, title, position)
VALUES (1, 'herbfneifj', 2);
-- TASKS
-- 3
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 4
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 5
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 1
INSERT INTO tasks (list_id, title, description, datetimes_id, position)
VALUES (1, 'First task', 'This is the first task', 3, 0);
-- 2
INSERT INTO tasks (list_id, title, description, datetimes_id, position)
VALUES (1, 'Second task', 'This is the second task', 4, 1);
-- 3
INSERT INTO tasks (list_id, title, description, datetimes_id, position)
VALUES (1, 'Third task', 'This is the third task', 5, 2);
-- 6
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 7
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 8
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 4
INSERT INTO tasks (list_id, title, description, datetimes_id, position)
VALUES (2, 'FIRST TASK', 'This is the first task in second list', 6, 0);
-- 5
INSERT INTO tasks (list_id, title, description, datetimes_id, position)
VALUES (2, 'SECOND TASK', 'This is the second task in second list', 7, 1);
-- 6
INSERT INTO tasks (list_id, title, description, datetimes_id, position)
VALUES (2, 'THIRD TASK', 'This is the third task in second list', 8, 2);
-- BOARD id=2
-- 5
INSERT INTO permissions (read, write, admin)
VALUES (true, true, true);
-- 6
INSERT INTO permissions (read, write, admin)
VALUES (true, true, false);
-- 9
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 2
INSERT INTO boards (
        project_id,
        owner_id,
        default_permissions_id,
        datetimes_id,
        title
    )
VALUES (1, 2, 6, 2, 'Second board');
-- 2
INSERT INTO board_users (user_id, board_id, permissions_id)
VALUES (2, 2, 5);
-- PROJECT id=2
-- 7
INSERT INTO permissions (read, write, admin)
VALUES (true, true, true);
-- 8
INSERT INTO permissions (read, write, admin)
VALUES (true, true, false);
-- 10
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 2
INSERT INTO projects (
        owner_id,
        default_permissions_id,
        datetimes_id,
        title,
        description
    )
VALUES (
        2,
        8,
        10,
        'Second project',
        'This is the second project'
    );
-- 2
INSERT INTO project_users (user_id, project_id, permissions_id)
VALUES (2, 2, 7);
-- PROJECT_USERS id=3
-- 9
INSERT INTO permissions (read, write, admin)
VALUES (true, true, false);
-- 3
INSERT INTO project_users (user_id, project_id, permissions_id)
VALUES (2, 1, 9);
-- PROJECT id=3
-- 10
INSERT INTO permissions (read, write, admin)
VALUES (true, true, true);
-- 11
INSERT INTO permissions (read, write, admin)
VALUES (true, true, false);
-- 11
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 3
INSERT INTO projects (
        owner_id,
        default_permissions_id,
        datetimes_id,
        title,
        description
    )
VALUES (
        3,
        11,
        11,
        'Third project',
        'This is the thrd project'
    );
-- 4
INSERT INTO project_users (user_id, project_id, permissions_id)
VALUES (3, 3, 10);
-- PROJECT_USERS id=4
-- 12
INSERT INTO permissions (read, write, admin)
VALUES (true, false, false);
-- 3
INSERT INTO project_users (user_id, project_id, permissions_id)
VALUES (3, 1, 12);
-- BOARDS id=3
-- 13
INSERT INTO permissions (read, write, admin)
VALUES (true, true, true);
-- 14
INSERT INTO permissions (read, write, admin)
VALUES (true, true, false);
-- 12
INSERT INTO datetimes (created, updated, accessed)
VALUES (1605925262, 1605925262, 1605925262);
-- 3
INSERT INTO boards (
        project_id,
        owner_id,
        default_permissions_id,
        datetimes_id,
        title
    )
VALUES (3, 3, 14, 12, 'First board in the second project');
-- 3
INSERT INTO board_users (user_id, board_id, permissions_id)
VALUES (3, 3, 13);
-- BOARD_USERS id=3
-- 15
INSERT INTO permissions (read, write, admin)
VALUES (true, true, true);
-- 4
INSERT INTO board_users (user_id, board_id, permissions_id) -- автор проекта
VALUES (1, 3, 15);
-- BOARD_USERS id=4
-- 16
INSERT INTO permissions (read, write, admin)
VALUES (true, false, false);
-- 5
INSERT INTO board_users (user_id, board_id, permissions_id)
VALUES (2, 3, 16);
-- USERS id=4
INSERT INTO users (nickname, email, avatar, password)
VALUES ('ivan', 'ivan@mail.ru', '', 'qwerty');
-- PROJECT_USERS id=5
-- 17
INSERT INTO permissions (read, write, admin)
VALUES (true, false, false);
-- 3
INSERT INTO project_users (user_id, project_id, permissions_id)
VALUES (4, 1, 17);