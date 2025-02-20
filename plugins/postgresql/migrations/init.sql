--
-- PostgreSQL database dump
--

-- Dumped from database version 17.2 (Debian 17.2-1.pgdg120+1)
-- Dumped by pg_dump version 17.0

-- Started on 2025-02-14 17:36:17 GMT

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 220 (class 1259 OID 16483)
-- Name: migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.migrations (
    current_version integer
);


ALTER TABLE public.migrations OWNER TO postgres;

--
-- TOC entry 217 (class 1259 OID 16456)
-- Name: settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.settings (
    setting_key character varying(255) NOT NULL,
    default_value character varying(255)
);


ALTER TABLE public.settings OWNER TO postgres;

--
-- TOC entry 218 (class 1259 OID 16461)
-- Name: topic_settings; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.topic_settings (
    topic_id uuid,
    setting_key character varying(255),
    value character varying(255)
);


ALTER TABLE public.topic_settings OWNER TO postgres;

--
-- TOC entry 219 (class 1259 OID 16466)
-- Name: topics; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.topics (
    topic_id uuid NOT NULL,
    topic_name character varying(255),
    num_partitions integer,
    replication_factor integer
);


ALTER TABLE public.topics OWNER TO postgres;

--
-- TOC entry 3218 (class 2606 OID 16470)
-- Name: settings settings_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.settings
    ADD CONSTRAINT settings_pkey PRIMARY KEY (setting_key);


--
-- TOC entry 3220 (class 2606 OID 16472)
-- Name: topics topics_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT topics_pkey PRIMARY KEY (topic_id);


--
-- TOC entry 3222 (class 2606 OID 16489)
-- Name: topics uc_topic_name; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.topics
    ADD CONSTRAINT uc_topic_name UNIQUE (topic_name);


--
-- TOC entry 3223 (class 2606 OID 16473)
-- Name: topic_settings fk_settings_key; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.topic_settings
    ADD CONSTRAINT fk_settings_key FOREIGN KEY (setting_key) REFERENCES public.settings(setting_key) NOT VALID;


--
-- TOC entry 3224 (class 2606 OID 16478)
-- Name: topic_settings fk_topic_id; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.topic_settings
    ADD CONSTRAINT fk_topic_id FOREIGN KEY (topic_id) REFERENCES public.topics(topic_id) NOT VALID;


-- Completed on 2025-02-14 17:36:17 GMT

--
-- PostgreSQL database dump complete
--

