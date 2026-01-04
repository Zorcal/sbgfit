-- migrate:up
CREATE TYPE exercise_category AS ENUM ('cardio', 'strength', 'plyometric');
CREATE TYPE exercise_equipment AS ENUM ('bodyweight', 'kettlebell', 'rowing-machine', 'ski-erg', 'medicine-ball', 'dumbbells', 'barbell', 'sled', 'box', 'jump-rope', 'assault-bike');
CREATE TYPE primary_muscle AS ENUM ('chest', 'back', 'shoulders', 'biceps', 'triceps', 'forearms', 'core', 'abs', 'obliques', 'glutes', 'quads', 'hamstrings', 'calves', 'legs', 'full-body', 'grip');
CREATE TYPE exercise_tag AS ENUM ('crossfit', 'hyrox', 'beginner-friendly', 'advanced', 'conditioning', 'strength-endurance', 'power', 'core', 'functional', 'competition', 'plyometric');

CREATE TABLE sbgfit.exercises (
    id SERIAL PRIMARY KEY,
    external_id UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    category exercise_category,
    equipment_types exercise_equipment[],
    primary_muscles primary_muscle[],
    description TEXT,
    instructions TEXT[],
    tags exercise_tag[],
    created_by_user_id INTEGER REFERENCES sbgfit.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Index for filtering by exercise category.
-- Example: WHERE category = 'strength'
CREATE INDEX idx_exercises_category ON sbgfit.exercises(category);

-- Index for filtering by equipment type.
-- Example: WHERE 'bodyweight' = ANY(equipment_types)
CREATE INDEX idx_exercises_equipment ON sbgfit.exercises USING GIN(equipment_types);

-- Index for finding exercises created by specific user.
-- Example: WHERE created_by_user_id = 123
CREATE INDEX idx_exercises_created_by_user_id ON sbgfit.exercises(created_by_user_id);

-- Index for case-insensitive name pattern matching (ILIKE queries).
-- Example: WHERE LOWER(name) LIKE LOWER('%push%')
CREATE INDEX idx_exercises_name_ilike ON sbgfit.exercises(LOWER(name));

-- Index for sorting by natural name order.
-- Example: ORDER BY name COLLATE natsort
CREATE INDEX idx_exercises_name_sort ON sbgfit.exercises(name COLLATE natsort);

-- migrate:down
DROP TABLE sbgfit.exercises;
DROP TYPE exercise_tag;
DROP TYPE primary_muscle;
DROP TYPE exercise_equipment;
DROP TYPE exercise_category;
