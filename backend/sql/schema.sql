DROP TABLE IF EXISTS public.completed_user_tasks;
DROP TABLE IF EXISTS public.tasks;
DROP TABLE IF EXISTS public.traders;
DROP TABLE IF EXISTS public.discord_users;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.discord_users (

    id BIGINT PRIMARY KEY

);

CREATE TABLE public.traders (

    id SERIAL PRIMARY KEY,
    name VARCHAR(11) UNIQUE NOT NULL

);

CREATE TABLE public.tasks (

    id SERIAL PRIMARY KEY,
    trader_id INT NOT NULL REFERENCES traders(id),
    name VARCHAR(100) NOT NULL,
    objectives TEXT NOT NULL,
    rewards TEXT NOT NULL,
    prerequisites TEXT,
    required_for_kappa BOOLEAN NOT NULL

);

CREATE TABLE public.completed_user_tasks (

    user_id BIGINT NOT NULL REFERENCES discord_users(id) ON DELETE CASCADE,
    task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, task_id)

);

ALTER TABLE public.discord_users ENABLE ROW LEVEL SECURITY;
CREATE POLICY "anon_all_permissions"
ON "public"."discord_users"
AS PERMISSIVE
TO anon
USING (true);

ALTER TABLE public.traders ENABLE ROW LEVEL SECURITY;
CREATE POLICY "anon_all_permissions"
ON "public"."traders"
AS PERMISSIVE
TO anon
USING (true);

ALTER TABLE public.tasks ENABLE ROW LEVEL SECURITY;
CREATE POLICY "anon_all_permissions"
ON "public"."tasks"
AS PERMISSIVE
TO anon
USING (true);

ALTER TABLE public.completed_user_tasks ENABLE ROW LEVEL SECURITY;
CREATE POLICY "anon_all_permissions"
ON "public"."completed_user_tasks"
AS PERMISSIVE
TO anon
USING (true);
