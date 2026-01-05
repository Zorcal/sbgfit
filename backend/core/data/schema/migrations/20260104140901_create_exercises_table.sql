-- migrate:up

-- Lookup tables

CREATE TABLE sbgfit.exercise_categories (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE sbgfit.equipment_types (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE sbgfit.primary_muscles (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE sbgfit.exercise_tags (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL
);

-- Main exercises table

CREATE TABLE sbgfit.exercises (
    id SERIAL PRIMARY KEY,
    external_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    category_id INTEGER REFERENCES sbgfit.exercise_categories(id),
    description TEXT,
    instructions TEXT[],
    created_by_user_id INTEGER REFERENCES sbgfit.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Junction tables

CREATE TABLE sbgfit.exercise_equipment (
    exercise_id INTEGER REFERENCES sbgfit.exercises(id) ON DELETE CASCADE,
    equipment_type_id INTEGER REFERENCES sbgfit.equipment_types(id),
    PRIMARY KEY (exercise_id, equipment_type_id)
);

CREATE TABLE sbgfit.exercise_primary_muscles (
    exercise_id INTEGER REFERENCES sbgfit.exercises(id) ON DELETE CASCADE,
    primary_muscle_id INTEGER REFERENCES sbgfit.primary_muscles(id),
    PRIMARY KEY (exercise_id, primary_muscle_id)
);

CREATE TABLE sbgfit.exercise_exercise_tags (
    exercise_id INTEGER REFERENCES sbgfit.exercises(id) ON DELETE CASCADE,
    exercise_tag_id INTEGER REFERENCES sbgfit.exercise_tags(id),
    PRIMARY KEY (exercise_id, exercise_tag_id)
);

-- Indexes for filtering
CREATE INDEX idx_exercises_category_id ON sbgfit.exercises(category_id);
CREATE INDEX idx_exercises_created_by_user_id ON sbgfit.exercises(created_by_user_id);
CREATE INDEX idx_exercises_name_ilike ON sbgfit.exercises(LOWER(name));
CREATE INDEX idx_exercises_name_sort ON sbgfit.exercises(name COLLATE natsort);

-- Junction table indexes
CREATE INDEX idx_exercise_equipment_exercise_id ON sbgfit.exercise_equipment(exercise_id);
CREATE INDEX idx_exercise_equipment_type_id ON sbgfit.exercise_equipment(equipment_type_id);
CREATE INDEX idx_exercise_muscles_exercise_id ON sbgfit.exercise_primary_muscles(exercise_id);
CREATE INDEX idx_exercise_muscles_muscle_id ON sbgfit.exercise_primary_muscles(primary_muscle_id);
CREATE INDEX idx_exercise_tags_exercise_id ON sbgfit.exercise_exercise_tags(exercise_id);
CREATE INDEX idx_exercise_tags_tag_id ON sbgfit.exercise_exercise_tags(exercise_tag_id);

CREATE VIEW sbgfit.exercise_details AS
SELECT
    e.id,
    e.external_id,
    e.name,
    ec.code as category,
    e.description,
    e.instructions,
    COALESCE(
        ARRAY_AGG(DISTINCT et.code) FILTER (WHERE et.code IS NOT NULL),
        ARRAY[]::TEXT[]
    ) AS equipment_types,
    COALESCE(
        ARRAY_AGG(DISTINCT pm.code) FILTER (WHERE pm.code IS NOT NULL),
        ARRAY[]::TEXT[]
    ) AS primary_muscles,
    COALESCE(
        ARRAY_AGG(DISTINCT tag.code) FILTER (WHERE tag.code IS NOT NULL),
        ARRAY[]::TEXT[]
    ) AS tags,
    e.created_by_user_id,
    e.created_at,
    e.updated_at
FROM sbgfit.exercises e
LEFT JOIN sbgfit.exercise_categories ec ON e.category_id = ec.id
LEFT JOIN sbgfit.exercise_equipment ee ON e.id = ee.exercise_id
LEFT JOIN sbgfit.equipment_types et ON ee.equipment_type_id = et.id
LEFT JOIN sbgfit.exercise_primary_muscles epm ON e.id = epm.exercise_id
LEFT JOIN sbgfit.primary_muscles pm ON epm.primary_muscle_id = pm.id
LEFT JOIN sbgfit.exercise_exercise_tags eet ON e.id = eet.exercise_id
LEFT JOIN sbgfit.exercise_tags tag ON eet.exercise_tag_id = tag.id
GROUP BY e.id, e.external_id, e.name, ec.code, e.description, e.instructions, e.created_by_user_id, e.created_at, e.updated_at;

-- migrate:down
DROP VIEW sbgfit.exercise_details;
DROP TABLE sbgfit.exercise_exercise_tags;
DROP TABLE sbgfit.exercise_primary_muscles;
DROP TABLE sbgfit.exercise_equipment;
DROP TABLE sbgfit.exercises;
DROP TABLE sbgfit.exercise_tags;
DROP TABLE sbgfit.primary_muscles;
DROP TABLE sbgfit.equipment_types;
DROP TABLE sbgfit.exercise_categories;
