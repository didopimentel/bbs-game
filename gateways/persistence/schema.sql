
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

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TYPE public.battle_action_type AS ENUM (
    'weapon',
    'spell',
    'item',
    'escape'
);

CREATE TYPE public.battle_participant_type AS ENUM (
    'creature',
    'player',
    'environment'
);

CREATE TYPE public.damage_type AS ENUM (
    'normal',
    'magic'
);

CREATE TABLE public.accounts (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    password text NOT NULL,
    email text NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE public.battle_actions (
    battle_id uuid NOT NULL,
    causer_id uuid NOT NULL,
    target_id uuid NOT NULL,
    action_type public.battle_action_type,
    value integer
);

CREATE TABLE public.battle_creatures (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    level integer DEFAULT 1 NOT NULL,
    experience integer DEFAULT 0 NOT NULL,
    damage text NOT NULL,
    hp integer DEFAULT 1 NOT NULL,
    CONSTRAINT battle_creatures_name_check CHECK ((name <> ''::text))
);

CREATE TABLE public.battle_participants (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    battle_id uuid NOT NULL,
    participant_type public.battle_participant_type,
    participant_id uuid
);

CREATE TABLE public.battles (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    finished boolean DEFAULT false
);

CREATE TABLE public.creatures (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    level integer DEFAULT 1 NOT NULL,
    experience integer DEFAULT 0 NOT NULL,
    damage text NOT NULL,
    hp integer DEFAULT 1 NOT NULL,
    CONSTRAINT creatures_name_check CHECK ((name <> ''::text))
);

CREATE TABLE public.drops (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    item_id uuid NOT NULL,
    rate double precision NOT NULL,
    quantity integer NOT NULL
);

CREATE TABLE public.equipment (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    player_id uuid NOT NULL,
    main_hand uuid,
    chest uuid
);

CREATE TABLE public.inventory_items (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    player_id uuid NOT NULL,
    item_id uuid NOT NULL
);

CREATE TABLE public.items (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    type text NOT NULL,
    damage_amount integer,
    heal_amount integer,
    damage_type public.damage_type
);

CREATE TABLE public.npc (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    CONSTRAINT npc_name_check CHECK ((name <> ''::text))
);

CREATE TABLE public.players (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name text NOT NULL,
    level integer DEFAULT 1 NOT NULL,
    experience integer DEFAULT 0 NOT NULL,
    damage text NOT NULL,
    hp integer DEFAULT 1 NOT NULL,
    total_hp integer DEFAULT 1 NOT NULL,
    account_id uuid NOT NULL,
    created_at timestamp without time zone,
    updated_at timestamp without time zone,
    CONSTRAINT players_name_check CHECK ((name <> ''::text))
);

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);

CREATE TABLE public.shop_items (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    npc_id uuid NOT NULL,
    item_id uuid NOT NULL,
    price integer NOT NULL
);

ALTER TABLE ONLY public.accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.battle_creatures
    ADD CONSTRAINT battle_creatures_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.battle_participants
    ADD CONSTRAINT battle_participants_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.battles
    ADD CONSTRAINT battles_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.creatures
    ADD CONSTRAINT creatures_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.drops
    ADD CONSTRAINT drops_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.inventory_items
    ADD CONSTRAINT inventory_items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.npc
    ADD CONSTRAINT npc_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);

ALTER TABLE ONLY public.shop_items
    ADD CONSTRAINT shop_items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.battle_actions
    ADD CONSTRAINT battle_actions_battle_id_fkey FOREIGN KEY (battle_id) REFERENCES public.battles(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.battle_actions
    ADD CONSTRAINT battle_actions_causer_id_fkey FOREIGN KEY (causer_id) REFERENCES public.battle_participants(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.battle_actions
    ADD CONSTRAINT battle_actions_target_id_fkey FOREIGN KEY (target_id) REFERENCES public.battle_participants(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.battle_participants
    ADD CONSTRAINT battle_participants_battle_id_fkey FOREIGN KEY (battle_id) REFERENCES public.battles(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.drops
    ADD CONSTRAINT drops_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_chest_fkey FOREIGN KEY (chest) REFERENCES public.items(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_main_hand_fkey FOREIGN KEY (main_hand) REFERENCES public.items(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.equipment
    ADD CONSTRAINT equipment_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.inventory_items
    ADD CONSTRAINT inventory_items_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.inventory_items
    ADD CONSTRAINT inventory_items_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.shop_items
    ADD CONSTRAINT shop_items_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.items(id) ON DELETE CASCADE;

ALTER TABLE ONLY public.shop_items
    ADD CONSTRAINT shop_items_npc_id_fkey FOREIGN KEY (npc_id) REFERENCES public.npc(id) ON DELETE CASCADE;

