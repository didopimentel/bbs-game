CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE accounts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    password text NOT NULL,
    email text NOT NULL,
    created_at timestamp,
    updated_at timestamp
);

CREATE TYPE damage_type AS ENUM (
    'normal',
    'magic'
);

CREATE TABLE players (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    level int DEFAULT 1 NOT NULL,
    experience int DEFAULT 0 NOT NULL,
    damage text NOT NULL,
    hp int DEFAULT 1 NOT NULL,
    total_hp int DEFAULT 1 NOT NULL,
    account_id uuid NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    created_at timestamp,
    updated_at timestamp
    CHECK (name <> '')
);

CREATE TABLE items (
                       id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name text NOT NULL,
                       type text NOT NULL,
                       damage_amount int,
                       heal_amount int,
                       damage_type damage_type
);


CREATE TABLE inventory_items (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    player_id uuid NOT NULL REFERENCES players (id) ON DELETE CASCADE,
    item_id uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE equipment (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
   player_id uuid NOT NULL REFERENCES players (id) ON DELETE CASCADE,
   main_hand uuid REFERENCES items (id) ON DELETE CASCADE,
   chest uuid REFERENCES items (id) ON DELETE CASCADE
);

CREATE TABLE creatures (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    level int DEFAULT 1 NOT NULL,
    experience int DEFAULT 0 NOT NULL,
    damage text NOT NULL,
    hp int DEFAULT 1 NOT NULL,
    CHECK (name <> '')
);

CREATE TABLE battle_creatures (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
   name text NOT NULL,
   level int DEFAULT 1 NOT NULL,
   experience int DEFAULT 0 NOT NULL,
   damage text NOT NULL,
   hp int DEFAULT 1 NOT NULL,
   CHECK (name <> '')
);

CREATE TABLE drops (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
   item_id uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE,
   rate float NOT NULL,
   quantity int NOT NULL
);

CREATE TABLE npc (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
   name text NOT NULL,
   CHECK (name <> '')
);

CREATE TABLE shop_items (
   id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
   npc_id uuid NOT NULL REFERENCES npc (id) ON DELETE CASCADE,
   item_id uuid NOT NULL REFERENCES items (id) ON DELETE CASCADE,
   price int NOT NULL
);

CREATE TYPE battle_participant_type AS ENUM (
    'creature',
    'player',
    'environment'
);

CREATE TABLE battles (
     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
     finished boolean DEFAULT FALSE
);

CREATE TABLE battle_participants (
     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
     battle_id uuid NOT NULL REFERENCES battles (id) ON DELETE CASCADE,
     participant_type battle_participant_type,
     participant_id uuid
);


CREATE TYPE battle_action_type AS ENUM (
    'weapon',
    'spell',
    'item',
    'escape'
);

CREATE TABLE battle_actions (
    battle_id uuid NOT NULL REFERENCES battles (id) ON DELETE CASCADE,
    causer_id uuid NOT NULL REFERENCES battle_participants (id) ON DELETE CASCADE,
    target_id uuid NOT NULL REFERENCES battle_participants (id) ON DELETE CASCADE,
    action_type battle_action_type,
    value int
);

INSERT INTO creatures (name, level, experience, damage, hp) VALUES
    ('bug', 1, 20, '1d6', 50),
    ('cobra', 1, 20, '1d6', 70),
    ('wolf', 2, 35, '1d6+5',  90),
    ('wasp', 2, 35, '1d6+4', 110),
    ('scorpion', 3, 60, '2d6', 150),
    ('minotaur', 3, 60, '2d6+2', 140),
    ('bandit', 3, 60, '2d6+1', 150);