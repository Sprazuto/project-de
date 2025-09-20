--
-- PostgreSQL database dump
--

SET statement_timeout = 0;

SET lock_timeout = 0;

SET client_encoding = 'UTF8';

SET standard_conforming_strings = on;

SET check_function_bodies = false;

SET client_min_messages = warning;

--
-- Name: golang_gin_db; Type: DATABASE; Schema: -; Owner: postgres
--
DROP DATABASE golang_gin_db;

CREATE DATABASE golang_gin_db
WITH
    TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';

ALTER DATABASE golang_gin_db OWNER TO postgres;

\connect golang_gin_db

SET statement_timeout = 0;

SET lock_timeout = 0;

SET client_encoding = 'UTF8';

SET standard_conforming_strings = on;

SET check_function_bodies = false;

SET client_min_messages = warning;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:
--
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;

--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:
--
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';

CREATE FUNCTION created_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	NEW.updated_at = EXTRACT(EPOCH FROM NOW());
	NEW.created_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;
$$;

ALTER FUNCTION public.created_at_column() OWNER TO postgres;

CREATE FUNCTION update_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = EXTRACT(EPOCH FROM NOW());
    RETURN NEW;
END;
$$;

ALTER FUNCTION public.update_at_column() OWNER TO postgres;

SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: article; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--
CREATE TABLE article (
    id integer NOT NULL,
    user_id integer,
    title character varying,
    content text,
    updated_at integer,
    created_at integer
);

ALTER TABLE article OWNER TO postgres;

--
-- Name: article_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE article_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE article_id_seq OWNER TO postgres;

ALTER SEQUENCE article_id_seq OWNED BY article.id;

--
-- Name: login_attempts; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--
CREATE TABLE login_attempts (
    id integer NOT NULL,
    user_id integer,
    success boolean DEFAULT false,
    attempt_time integer,
    updated_at integer,
    created_at integer
);

ALTER TABLE login_attempts OWNER TO postgres;

--
-- Name: login_attempts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE login_attempts_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE login_attempts_id_seq OWNER TO postgres;

ALTER SEQUENCE login_attempts_id_seq OWNED BY login_attempts.id;

--
-- Name: role_permissions; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--
CREATE TABLE role_permissions (
    id integer NOT NULL,
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    updated_at integer,
    created_at integer
);

ALTER TABLE role_permissions OWNER TO postgres;

--
-- Name: role_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE role_permissions_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE role_permissions_id_seq OWNER TO postgres;

ALTER SEQUENCE role_permissions_id_seq OWNED BY role_permissions.id;

--
-- Name: roles; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--
CREATE TABLE roles (
    id integer NOT NULL,
    name character varying NOT NULL,
    updated_at integer,
    created_at integer
);

ALTER TABLE roles OWNER TO postgres;

--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE roles_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE roles_id_seq OWNER TO postgres;

ALTER SEQUENCE roles_id_seq OWNED BY roles.id;

--
-- Name: permissions; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--
CREATE TABLE permissions (
    id integer NOT NULL,
    name character varying NOT NULL,
    updated_at integer,
    created_at integer
);

ALTER TABLE permissions OWNER TO postgres;

--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE permissions_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE permissions_id_seq OWNER TO postgres;

ALTER SEQUENCE permissions_id_seq OWNED BY permissions.id;

--
-- Name: user_roles; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--
CREATE TABLE user_roles (
    id integer NOT NULL,
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    updated_at integer,
    created_at integer
);

ALTER TABLE user_roles OWNER TO postgres;

--
-- Name: user_roles_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE user_roles_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE user_roles_id_seq OWNER TO postgres;

ALTER SEQUENCE user_roles_id_seq OWNED BY user_roles.id;

--
-- Name: user; Type: TABLE; Schema: public; Owner: postgres; Tablespace:
--
CREATE TABLE "user" (
    id integer NOT NULL,
    email character varying,
    username character varying NOT NULL,
    password character varying,
    name character varying,
    failed_attempts integer DEFAULT 0,
    locked_until integer DEFAULT 0,
    updated_at integer,
    created_at integer
);

ALTER TABLE "user" OWNER TO postgres;

--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--
CREATE SEQUENCE user_id_seq START
WITH
    1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER TABLE user_id_seq OWNER TO postgres;

ALTER SEQUENCE user_id_seq OWNED BY "user".id;

--
-- Name: id; Type: DEFAULT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY article
ALTER COLUMN id
SET DEFAULT nextval('article_id_seq'::regclass);

ALTER TABLE ONLY login_attempts
ALTER COLUMN id
SET DEFAULT nextval(
    'login_attempts_id_seq'::regclass
);

ALTER TABLE ONLY role_permissions
ALTER COLUMN id
SET DEFAULT nextval(
    'role_permissions_id_seq'::regclass
);

ALTER TABLE ONLY roles
ALTER COLUMN id
SET DEFAULT nextval('roles_id_seq'::regclass);

ALTER TABLE ONLY permissions
ALTER COLUMN id
SET DEFAULT nextval(
    'permissions_id_seq'::regclass
);

ALTER TABLE ONLY user_roles
ALTER COLUMN id
SET DEFAULT nextval('user_roles_id_seq'::regclass);

ALTER TABLE ONLY "user"
ALTER COLUMN id
SET DEFAULT nextval('user_id_seq'::regclass);

--
-- Data for Name: article; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY article (id, user_id, title, content, updated_at, created_at) FROM stdin;
\.

SELECT pg_catalog.setval ('article_id_seq', 1, false);

--
-- Data for Name: login_attempts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY login_attempts (id, user_id, success, attempt_time, updated_at, created_at) FROM stdin;
\.

SELECT pg_catalog.setval ( 'login_attempts_id_seq', 1, false );

--
-- Data for Name: role_permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY role_permissions (id, role_id, permission_id, updated_at, created_at) FROM stdin;
\.

SELECT pg_catalog.setval ( 'role_permissions_id_seq', 1, false );

--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--
COPY roles (id, name, updated_at, created_at) FROM stdin; \.

SELECT pg_catalog.setval ('roles_id_seq', 1, false);

--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: postgres
--
COPY permissions (id, name, updated_at, created_at) FROM stdin; \.

SELECT pg_catalog.setval ( 'permissions_id_seq', 1, false );

--
-- Data for Name: user_roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY user_roles (id, user_id, role_id, updated_at, created_at) FROM stdin;
\.

SELECT pg_catalog.setval ('user_roles_id_seq', 1, false);

--
-- Data for Name: user; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY "user" (id, email, username, password, name, failed_attempts, locked_until, updated_at, created_at) FROM stdin;
\.

SELECT pg_catalog.setval ('user_id_seq', 1, false);

--
-- Name: article article_id; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY article
ADD CONSTRAINT article_id PRIMARY KEY (id);

--
-- Name: login_attempts login_attempts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY login_attempts
ADD CONSTRAINT login_attempts_pkey PRIMARY KEY (id);

--
-- Name: role_permissions role_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY role_permissions
ADD CONSTRAINT role_permissions_pkey PRIMARY KEY (id);

--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY roles
ADD CONSTRAINT roles_pkey PRIMARY KEY (id);

--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY permissions
ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);

--
-- Name: user_roles user_roles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY user_roles
ADD CONSTRAINT user_roles_pkey PRIMARY KEY (id);

--
-- Name: user user_id; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY "user" ADD CONSTRAINT user_id PRIMARY KEY (id);

--
-- Name: user user_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres; Tablespace:
--
ALTER TABLE ONLY "user"
ADD CONSTRAINT user_username_key UNIQUE (username);

--
-- Name: article article_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY article
ADD CONSTRAINT article_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: login_attempts login_attempts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY login_attempts
ADD CONSTRAINT login_attempts_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: role_permissions role_permissions_permission_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY role_permissions
ADD CONSTRAINT role_permissions_permission_id_fkey FOREIGN KEY (permission_id) REFERENCES permissions (id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: role_permissions role_permissions_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY role_permissions
ADD CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles (id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: user_roles user_roles_role_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY user_roles
ADD CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles (id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: user_roles user_roles_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--
ALTER TABLE ONLY user_roles
ADD CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) ON UPDATE CASCADE ON DELETE CASCADE;

--
-- Name: article create_article_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER create_article_created_at BEFORE INSERT ON article FOR EACH ROW EXECUTE PROCEDURE created_at_column();

--
-- Name: login_attempts create_login_attempts_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER create_login_attempts_created_at BEFORE INSERT ON login_attempts FOR EACH ROW EXECUTE PROCEDURE created_at_column();

--
-- Name: role_permissions create_role_permissions_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER create_role_permissions_created_at BEFORE INSERT ON role_permissions FOR EACH ROW EXECUTE PROCEDURE created_at_column();

--
-- Name: roles create_roles_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER create_roles_created_at BEFORE INSERT ON roles FOR EACH ROW EXECUTE PROCEDURE created_at_column();

--
-- Name: permissions create_permissions_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER create_permissions_created_at BEFORE INSERT ON permissions FOR EACH ROW EXECUTE PROCEDURE created_at_column();

--
-- Name: user_roles create_user_roles_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER create_user_roles_created_at BEFORE INSERT ON user_roles FOR EACH ROW EXECUTE PROCEDURE created_at_column();

--
-- Name: user create_user_created_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER create_user_created_at BEFORE INSERT ON "user" FOR EACH ROW EXECUTE PROCEDURE created_at_column();

--
-- Name: article update_article_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER update_article_updated_at BEFORE UPDATE ON article FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: login_attempts update_login_attempts_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER update_login_attempts_updated_at BEFORE UPDATE ON login_attempts FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: role_permissions update_role_permissions_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER update_role_permissions_updated_at BEFORE UPDATE ON role_permissions FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: roles update_roles_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER update_roles_updated_at BEFORE UPDATE ON roles FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: permissions update_permissions_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER update_permissions_updated_at BEFORE UPDATE ON permissions FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: user_roles update_user_roles_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER update_user_roles_updated_at BEFORE UPDATE ON user_roles FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: user update_user_updated_at; Type: TRIGGER; Schema: public; Owner: postgres
--
CREATE TRIGGER update_user_updated_at BEFORE UPDATE ON "user" FOR EACH ROW EXECUTE PROCEDURE update_at_column();

--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--
REVOKE ALL ON SCHEMA public FROM PUBLIC;

REVOKE ALL ON SCHEMA public FROM postgres;

GRANT ALL ON SCHEMA public TO postgres;

GRANT ALL ON SCHEMA public TO PUBLIC;

--
-- PostgreSQL database dump complete
--
