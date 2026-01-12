SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: sbgfit; Type: SCHEMA; Schema: -; Owner: -
--

CREATE SCHEMA sbgfit;


--
-- Name: natsort; Type: COLLATION; Schema: public; Owner: -
--

CREATE COLLATION public.natsort (provider = icu, deterministic = false, locale = 'und-u-kn-ks-level2');


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: equipment_types; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.equipment_types (
    id integer NOT NULL,
    code text NOT NULL,
    name text NOT NULL
);


--
-- Name: equipment_types_id_seq; Type: SEQUENCE; Schema: sbgfit; Owner: -
--

CREATE SEQUENCE sbgfit.equipment_types_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: equipment_types_id_seq; Type: SEQUENCE OWNED BY; Schema: sbgfit; Owner: -
--

ALTER SEQUENCE sbgfit.equipment_types_id_seq OWNED BY sbgfit.equipment_types.id;


--
-- Name: exercise_categories; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.exercise_categories (
    id integer NOT NULL,
    code text NOT NULL,
    name text NOT NULL
);


--
-- Name: exercise_categories_id_seq; Type: SEQUENCE; Schema: sbgfit; Owner: -
--

CREATE SEQUENCE sbgfit.exercise_categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exercise_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: sbgfit; Owner: -
--

ALTER SEQUENCE sbgfit.exercise_categories_id_seq OWNED BY sbgfit.exercise_categories.id;


--
-- Name: exercise_equipment; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.exercise_equipment (
    exercise_id integer NOT NULL,
    equipment_type_id integer NOT NULL
);


--
-- Name: exercise_exercise_tags; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.exercise_exercise_tags (
    exercise_id integer NOT NULL,
    exercise_tag_id integer NOT NULL
);


--
-- Name: exercise_primary_muscles; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.exercise_primary_muscles (
    exercise_id integer NOT NULL,
    primary_muscle_id integer NOT NULL
);


--
-- Name: exercise_tags; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.exercise_tags (
    id integer NOT NULL,
    code text NOT NULL,
    name text NOT NULL
);


--
-- Name: exercise_tags_id_seq; Type: SEQUENCE; Schema: sbgfit; Owner: -
--

CREATE SEQUENCE sbgfit.exercise_tags_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exercise_tags_id_seq; Type: SEQUENCE OWNED BY; Schema: sbgfit; Owner: -
--

ALTER SEQUENCE sbgfit.exercise_tags_id_seq OWNED BY sbgfit.exercise_tags.id;


--
-- Name: exercises; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.exercises (
    id integer NOT NULL,
    external_id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    category_id integer,
    description text,
    instructions text[],
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: exercises_id_seq; Type: SEQUENCE; Schema: sbgfit; Owner: -
--

CREATE SEQUENCE sbgfit.exercises_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: exercises_id_seq; Type: SEQUENCE OWNED BY; Schema: sbgfit; Owner: -
--

ALTER SEQUENCE sbgfit.exercises_id_seq OWNED BY sbgfit.exercises.id;


--
-- Name: primary_muscles; Type: TABLE; Schema: sbgfit; Owner: -
--

CREATE TABLE sbgfit.primary_muscles (
    id integer NOT NULL,
    code text NOT NULL,
    name text NOT NULL
);


--
-- Name: primary_muscles_id_seq; Type: SEQUENCE; Schema: sbgfit; Owner: -
--

CREATE SEQUENCE sbgfit.primary_muscles_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: primary_muscles_id_seq; Type: SEQUENCE OWNED BY; Schema: sbgfit; Owner: -
--

ALTER SEQUENCE sbgfit.primary_muscles_id_seq OWNED BY sbgfit.primary_muscles.id;


--
-- Name: equipment_types id; Type: DEFAULT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.equipment_types ALTER COLUMN id SET DEFAULT nextval('sbgfit.equipment_types_id_seq'::regclass);


--
-- Name: exercise_categories id; Type: DEFAULT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_categories ALTER COLUMN id SET DEFAULT nextval('sbgfit.exercise_categories_id_seq'::regclass);


--
-- Name: exercise_tags id; Type: DEFAULT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_tags ALTER COLUMN id SET DEFAULT nextval('sbgfit.exercise_tags_id_seq'::regclass);


--
-- Name: exercises id; Type: DEFAULT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercises ALTER COLUMN id SET DEFAULT nextval('sbgfit.exercises_id_seq'::regclass);


--
-- Name: primary_muscles id; Type: DEFAULT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.primary_muscles ALTER COLUMN id SET DEFAULT nextval('sbgfit.primary_muscles_id_seq'::regclass);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: equipment_types equipment_types_code_key; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.equipment_types
    ADD CONSTRAINT equipment_types_code_key UNIQUE (code);


--
-- Name: equipment_types equipment_types_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.equipment_types
    ADD CONSTRAINT equipment_types_pkey PRIMARY KEY (id);


--
-- Name: exercise_categories exercise_categories_code_key; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_categories
    ADD CONSTRAINT exercise_categories_code_key UNIQUE (code);


--
-- Name: exercise_categories exercise_categories_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_categories
    ADD CONSTRAINT exercise_categories_pkey PRIMARY KEY (id);


--
-- Name: exercise_equipment exercise_equipment_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_equipment
    ADD CONSTRAINT exercise_equipment_pkey PRIMARY KEY (exercise_id, equipment_type_id);


--
-- Name: exercise_exercise_tags exercise_exercise_tags_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_exercise_tags
    ADD CONSTRAINT exercise_exercise_tags_pkey PRIMARY KEY (exercise_id, exercise_tag_id);


--
-- Name: exercise_primary_muscles exercise_primary_muscles_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_primary_muscles
    ADD CONSTRAINT exercise_primary_muscles_pkey PRIMARY KEY (exercise_id, primary_muscle_id);


--
-- Name: exercise_tags exercise_tags_code_key; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_tags
    ADD CONSTRAINT exercise_tags_code_key UNIQUE (code);


--
-- Name: exercise_tags exercise_tags_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_tags
    ADD CONSTRAINT exercise_tags_pkey PRIMARY KEY (id);


--
-- Name: exercises exercises_external_id_key; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercises
    ADD CONSTRAINT exercises_external_id_key UNIQUE (external_id);


--
-- Name: exercises exercises_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercises
    ADD CONSTRAINT exercises_pkey PRIMARY KEY (id);


--
-- Name: primary_muscles primary_muscles_code_key; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.primary_muscles
    ADD CONSTRAINT primary_muscles_code_key UNIQUE (code);


--
-- Name: primary_muscles primary_muscles_pkey; Type: CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.primary_muscles
    ADD CONSTRAINT primary_muscles_pkey PRIMARY KEY (id);


--
-- Name: idx_exercise_equipment_exercise_id; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercise_equipment_exercise_id ON sbgfit.exercise_equipment USING btree (exercise_id);


--
-- Name: idx_exercise_equipment_type_id; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercise_equipment_type_id ON sbgfit.exercise_equipment USING btree (equipment_type_id);


--
-- Name: idx_exercise_muscles_exercise_id; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercise_muscles_exercise_id ON sbgfit.exercise_primary_muscles USING btree (exercise_id);


--
-- Name: idx_exercise_muscles_muscle_id; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercise_muscles_muscle_id ON sbgfit.exercise_primary_muscles USING btree (primary_muscle_id);


--
-- Name: idx_exercise_tags_exercise_id; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercise_tags_exercise_id ON sbgfit.exercise_exercise_tags USING btree (exercise_id);


--
-- Name: idx_exercise_tags_tag_id; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercise_tags_tag_id ON sbgfit.exercise_exercise_tags USING btree (exercise_tag_id);


--
-- Name: idx_exercises_category_id; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercises_category_id ON sbgfit.exercises USING btree (category_id);


--
-- Name: idx_exercises_name_ilike; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercises_name_ilike ON sbgfit.exercises USING btree (lower(name));


--
-- Name: idx_exercises_name_sort; Type: INDEX; Schema: sbgfit; Owner: -
--

CREATE INDEX idx_exercises_name_sort ON sbgfit.exercises USING btree (name COLLATE public.natsort);


--
-- Name: exercise_equipment exercise_equipment_equipment_type_id_fkey; Type: FK CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_equipment
    ADD CONSTRAINT exercise_equipment_equipment_type_id_fkey FOREIGN KEY (equipment_type_id) REFERENCES sbgfit.equipment_types(id);


--
-- Name: exercise_equipment exercise_equipment_exercise_id_fkey; Type: FK CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_equipment
    ADD CONSTRAINT exercise_equipment_exercise_id_fkey FOREIGN KEY (exercise_id) REFERENCES sbgfit.exercises(id) ON DELETE CASCADE;


--
-- Name: exercise_exercise_tags exercise_exercise_tags_exercise_id_fkey; Type: FK CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_exercise_tags
    ADD CONSTRAINT exercise_exercise_tags_exercise_id_fkey FOREIGN KEY (exercise_id) REFERENCES sbgfit.exercises(id) ON DELETE CASCADE;


--
-- Name: exercise_exercise_tags exercise_exercise_tags_exercise_tag_id_fkey; Type: FK CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_exercise_tags
    ADD CONSTRAINT exercise_exercise_tags_exercise_tag_id_fkey FOREIGN KEY (exercise_tag_id) REFERENCES sbgfit.exercise_tags(id);


--
-- Name: exercise_primary_muscles exercise_primary_muscles_exercise_id_fkey; Type: FK CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_primary_muscles
    ADD CONSTRAINT exercise_primary_muscles_exercise_id_fkey FOREIGN KEY (exercise_id) REFERENCES sbgfit.exercises(id) ON DELETE CASCADE;


--
-- Name: exercise_primary_muscles exercise_primary_muscles_primary_muscle_id_fkey; Type: FK CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercise_primary_muscles
    ADD CONSTRAINT exercise_primary_muscles_primary_muscle_id_fkey FOREIGN KEY (primary_muscle_id) REFERENCES sbgfit.primary_muscles(id);


--
-- Name: exercises exercises_category_id_fkey; Type: FK CONSTRAINT; Schema: sbgfit; Owner: -
--

ALTER TABLE ONLY sbgfit.exercises
    ADD CONSTRAINT exercises_category_id_fkey FOREIGN KEY (category_id) REFERENCES sbgfit.exercise_categories(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20260104120133'),
    ('20260104135809'),
    ('20260104140901');
