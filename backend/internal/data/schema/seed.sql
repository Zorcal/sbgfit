BEGIN;

-- Seed lookup tables first

INSERT INTO sbgfit.exercise_categories (code, name) VALUES
('cardio', 'Cardio'),
('strength', 'Strength'),
('plyometric', 'Plyometric')
ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO sbgfit.equipment_types (code, name) VALUES
('bodyweight', 'Bodyweight'),
('kettlebell', 'Kettlebell'),
('rowing-machine', 'Rowing Machine'),
('ski-erg', 'Ski Erg'),
('medicine-ball', 'Medicine Ball'),
('dumbbells', 'Dumbbells'),
('barbell', 'Barbell'),
('sled', 'Sled'),
('box', 'Box'),
('jump-rope', 'Jump Rope'),
('assault-bike', 'Assault Bike')
ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO sbgfit.primary_muscles (code, name) VALUES
('chest', 'Chest'),
('back', 'Back'),
('shoulders', 'Shoulders'),
('biceps', 'Biceps'),
('triceps', 'Triceps'),
('forearms', 'Forearms'),
('core', 'Core'),
('abs', 'Abs'),
('obliques', 'Obliques'),
('glutes', 'Glutes'),
('quads', 'Quads'),
('hamstrings', 'Hamstrings'),
('calves', 'Calves'),
('legs', 'Legs'),
('full-body', 'Full Body'),
('grip', 'Grip')
ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO sbgfit.exercise_tags (code, name) VALUES
('crossfit', 'CrossFit'),
('hyrox', 'Hyrox'),
('beginner-friendly', 'Beginner Friendly'),
('advanced', 'Advanced'),
('conditioning', 'Conditioning'),
('strength-endurance', 'Strength Endurance'),
('power', 'Power'),
('core', 'Core'),
('functional', 'Functional'),
('competition', 'Competition'),
('plyometric', 'Plyometric')
ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name;

-- Helper function to insert exercise with all relationships
CREATE OR REPLACE FUNCTION insert_exercise(
    p_external_id UUID,
    p_name TEXT,
    p_category_code TEXT,
    p_description TEXT,
    p_instructions TEXT[],
    p_equipment_codes TEXT[] DEFAULT ARRAY[]::TEXT[],
    p_muscle_codes TEXT[] DEFAULT ARRAY[]::TEXT[],
    p_tag_codes TEXT[] DEFAULT ARRAY[]::TEXT[]
) RETURNS INTEGER AS $$
DECLARE
    current_exercise_id INTEGER;
    equipment_code TEXT;
    muscle_code TEXT;
    tag_code TEXT;
BEGIN
    -- Insert or update main exercise
    INSERT INTO sbgfit.exercises (external_id, name, category_id, description, instructions)
    VALUES (
        p_external_id,
        p_name,
        (SELECT id FROM sbgfit.exercise_categories WHERE code = p_category_code),
        p_description,
        p_instructions
    )
    ON CONFLICT (external_id) DO UPDATE SET
        name = EXCLUDED.name,
        category_id = EXCLUDED.category_id,
        description = EXCLUDED.description,
        instructions = EXCLUDED.instructions,
        updated_at = CURRENT_TIMESTAMP
    RETURNING id INTO current_exercise_id;

    -- Replace all relationships (delete existing, insert new)
    DELETE FROM sbgfit.exercise_equipment WHERE exercise_id = current_exercise_id;
    DELETE FROM sbgfit.exercise_primary_muscles WHERE exercise_id = current_exercise_id;  
    DELETE FROM sbgfit.exercise_exercise_tags WHERE exercise_id = current_exercise_id;

    -- Insert equipment relationships
    FOREACH equipment_code IN ARRAY p_equipment_codes
    LOOP
        INSERT INTO sbgfit.exercise_equipment (exercise_id, equipment_type_id)
        VALUES (current_exercise_id, (SELECT id FROM sbgfit.equipment_types WHERE code = equipment_code));
    END LOOP;

    -- Insert muscle relationships
    FOREACH muscle_code IN ARRAY p_muscle_codes
    LOOP
        INSERT INTO sbgfit.exercise_primary_muscles (exercise_id, primary_muscle_id)
        VALUES (current_exercise_id, (SELECT id FROM sbgfit.primary_muscles WHERE code = muscle_code));
    END LOOP;

    -- Insert tag relationships
    FOREACH tag_code IN ARRAY p_tag_codes
    LOOP
        INSERT INTO sbgfit.exercise_exercise_tags (exercise_id, exercise_tag_id)
        VALUES (current_exercise_id, (SELECT id FROM sbgfit.exercise_tags WHERE code = tag_code));
    END LOOP;

    RETURN current_exercise_id;
END;
$$ LANGUAGE plpgsql;

-- Global exercises using helper function

-- Burpees
SELECT insert_exercise(
    '01234567-89ab-cdef-0123-456789abcdef',
    'Burpees',
    'cardio',
    'From standing, squat down, jump back to plank, do a push-up, jump feet back to squat, then jump up with arms overhead',
    ARRAY[
        'Start standing',
        'Squat down hands on ground',
        'Jump back to plank',
        'Do push-up',
        'Jump feet to squat',
        'Jump up arms overhead'
    ],
    ARRAY['bodyweight'],
    ARRAY['full-body'],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'functional', 'competition']
);

-- Kettlebell Swings
SELECT insert_exercise(
    '11111111-1111-1111-1111-111111111111',
    'Kettlebell Swings',
    'strength',
    'Hip-hinge movement swinging kettlebell from between legs to chest height',
    ARRAY[
        'Stand with kettlebell',
        'Hinge at hips grab bell',
        'Drive hips forward swing up',
        'Let bell swing back',
        'Repeat motion'
    ],
    ARRAY['kettlebell'],
    ARRAY['glutes', 'hamstrings', 'core'],
    ARRAY['crossfit', 'power', 'functional']
);

-- Rowing
SELECT insert_exercise(
    '22222222-2222-2222-2222-222222222222',
    'Rowing',
    'cardio',
    'Full-body cardio movement on rowing machine',
    ARRAY[
        'Sit on machine feet strapped',
        'Grab handle',
        'Push legs lean back',
        'Pull to chest',
        'Reverse movement'
    ],
    ARRAY['rowing-machine'],
    ARRAY['back', 'legs', 'core'],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'core']
);

-- Ski Erg
SELECT insert_exercise(
    '33333333-3333-3333-3333-333333333333',
    'Ski Erg',
    'cardio',
    'Upper body cardio movement mimicking cross-country skiing',
    ARRAY[
        'Stand feet hip-width',
        'Grab handles overhead',
        'Pull down skiing motion',
        'Return overhead',
        'Maintain rhythm'
    ],
    ARRAY['ski-erg'],
    ARRAY['shoulders', 'core', 'legs'],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'core']
);

-- Wall Balls
SELECT insert_exercise(
    '44444444-4444-4444-4444-444444444444',
    'Wall Balls',
    'strength',
    'Squat and throw medicine ball to target on wall',
    ARRAY[
        'Hold ball at chest',
        'Squat keeping chest up',
        'Drive up throw to target',
        'Catch ball squat again',
        'Repeat continuously'
    ],
    ARRAY['medicine-ball'],
    ARRAY['legs', 'shoulders', 'core'],
    ARRAY['crossfit', 'power', 'functional']
);

-- Farmers Walk
SELECT insert_exercise(
    '55555555-5555-5555-5555-555555555555',
    'Farmers Walk',
    'strength',
    'Walk while carrying heavy weights in each hand',
    ARRAY[
        'Pick up weights',
        'Stand tall shoulders back',
        'Walk maintaining posture',
        'Keep core tight',
        'Set down safely'
    ],
    ARRAY['dumbbells'],
    ARRAY['grip', 'core', 'legs'],
    ARRAY['hyrox', 'strength-endurance', 'functional']
);

-- Sled Push
SELECT insert_exercise(
    '66666666-6666-6666-6666-666666666666',
    'Sled Push',
    'strength',
    'Push weighted sled across floor',
    ARRAY[
        'Hands on handles',
        'Lean forward straight back',
        'Drive with legs forward',
        'Maintain pace',
        'Keep core engaged'
    ],
    ARRAY['sled'],
    ARRAY['legs', 'glutes', 'core'],
    ARRAY['hyrox', 'strength-endurance', 'functional']
);

-- Sled Pull
SELECT insert_exercise(
    '77777777-7777-7777-7777-777777777777',
    'Sled Pull',
    'strength',
    'Pull weighted sled toward you',
    ARRAY[
        'Grab rope or handles',
        'Lean back slightly',
        'Pull hand over hand',
        'Reset position',
        'Maintain rhythm'
    ],
    ARRAY['sled'],
    ARRAY['back', 'biceps', 'core'],
    ARRAY['hyrox', 'strength-endurance', 'functional']
);

-- Box Jumps
SELECT insert_exercise(
    '88888888-8888-8888-8888-888888888888',
    'Box Jumps',
    'plyometric',
    'Jump onto elevated box or platform',
    ARRAY[
        'Stand in front of box',
        'Swing arms bend knees',
        'Jump up land softly',
        'Stand upright on box',
        'Step down safely'
    ],
    ARRAY['box'],
    ARRAY['legs', 'glutes'],
    ARRAY['crossfit', 'power', 'plyometric']
);

-- Lunges
SELECT insert_exercise(
    '99999999-9999-9999-9999-999999999999',
    'Lunges',
    'strength',
    'Single-leg strength movement stepping forward into lunge position',
    ARRAY[
        'Stand feet hip-width',
        'Step forward to lunge',
        'Lower back knee down',
        'Push through front heel',
        'Alternate or complete side'
    ],
    ARRAY['bodyweight'],
    ARRAY['legs', 'glutes'],
    ARRAY['crossfit', 'hyrox', 'beginner-friendly', 'functional']
);

-- Pull-ups
SELECT insert_exercise(
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'Pull-ups',
    'strength',
    'Hanging from a bar and pulling body up until chin clears the bar',
    ARRAY[
        'Hang from pull-up bar',
        'Pull body up',
        'Chin over bar',
        'Lower with control',
        'Repeat'
    ],
    ARRAY['bodyweight'],
    ARRAY['back', 'biceps'],
    ARRAY['crossfit', 'functional', 'beginner-friendly']
);

-- Push-ups
SELECT insert_exercise(
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'Push-ups',
    'strength',
    'Classic bodyweight exercise targeting chest, shoulders, and triceps',
    ARRAY[
        'Start in plank position',
        'Lower chest to ground',
        'Push back to start',
        'Keep body straight',
        'Repeat'
    ],
    ARRAY['bodyweight'],
    ARRAY['chest', 'shoulders', 'triceps'],
    ARRAY['crossfit', 'beginner-friendly', 'functional']
);

-- Dumbbell Deadlifts
SELECT insert_exercise(
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'Dumbbell Deadlifts',
    'strength',
    'Hip hinge movement lifting dumbbells from ground to standing position',
    ARRAY[
        'Stand with feet hip-width',
        'Hinge at hips',
        'Grab weights',
        'Drive hips forward',
        'Stand tall'
    ],
    ARRAY['dumbbells'],
    ARRAY['back', 'glutes', 'hamstrings'],
    ARRAY['crossfit', 'functional', 'strength-endurance']
);

-- Air Squats
SELECT insert_exercise(
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'Air Squats',
    'strength',
    'Bodyweight squat focusing on proper hip and knee movement',
    ARRAY[
        'Stand with feet shoulder-width',
        'Lower hips back and down',
        'Keep chest up',
        'Drive through heels',
        'Return to standing'
    ],
    ARRAY['bodyweight'],
    ARRAY['legs', 'glutes'],
    ARRAY['crossfit', 'beginner-friendly', 'functional']
);

-- Dumbbell Thrusters
SELECT insert_exercise(
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    'Dumbbell Thrusters',
    'strength',
    'Combination squat to overhead press with dumbbells',
    ARRAY[
        'Hold weights at shoulders',
        'Squat down',
        'Drive up explosively',
        'Press weights overhead',
        'Lower to shoulders'
    ],
    ARRAY['dumbbells'],
    ARRAY['legs', 'shoulders', 'core'],
    ARRAY['crossfit', 'functional', 'conditioning']
);

-- Double Unders
SELECT insert_exercise(
    'ffffffff-ffff-ffff-ffff-ffffffffffff',
    'Double Unders',
    'cardio',
    'Jump rope where rope passes under feet twice per jump',
    ARRAY[
        'Hold rope handles',
        'Jump higher than normal',
        'Spin rope faster',
        'Land on balls of feet',
        'Keep rhythm consistent'
    ],
    ARRAY['jump-rope'],
    ARRAY['legs', 'core'],
    ARRAY['crossfit', 'conditioning', 'advanced']
);

-- Mountain Climbers
SELECT insert_exercise(
    '10101010-1010-1010-1010-101010101010',
    'Mountain Climbers',
    'cardio',
    'Dynamic plank position with alternating knee drives',
    ARRAY[
        'Start in plank position',
        'Drive right knee to chest',
        'Switch legs quickly',
        'Keep hips level',
        'Maintain fast pace'
    ],
    ARRAY['bodyweight'],
    ARRAY['core', 'legs'],
    ARRAY['crossfit', 'conditioning', 'core']
);

-- Turkish Get-ups
SELECT insert_exercise(
    '11111111-2222-3333-4444-555555555555',
    'Turkish Get-ups',
    'strength',
    'Complex movement from lying to standing while holding weight overhead',
    ARRAY[
        'Lie on back with weight up',
        'Roll to elbow',
        'Push to hand',
        'Bridge hips up',
        'Stand up slowly'
    ],
    ARRAY['kettlebell'],
    ARRAY['core', 'shoulders', 'full-body'],
    ARRAY['functional', 'advanced', 'core']
);

-- Dumbbell Bench Press
SELECT insert_exercise(
    '22222222-3333-4444-5555-666666666666',
    'Dumbbell Bench Press',
    'strength',
    'Upper body pressing movement with dumbbells for chest development',
    ARRAY[
        'Lie on bench',
        'Lower weights to chest',
        'Press up explosively',
        'Keep back flat',
        'Control the weight'
    ],
    ARRAY['dumbbells'],
    ARRAY['chest', 'shoulders', 'triceps'],
    ARRAY['functional', 'strength-endurance']
);

-- Dumbbell Bent-over Rows
SELECT insert_exercise(
    '33333333-4444-5555-6666-777777777777',
    'Dumbbell Bent-over Rows',
    'strength',
    'Pulling movement with dumbbells targeting back muscles and posterior chain',
    ARRAY[
        'Hinge at hips',
        'Hold weights with arms extended',
        'Pull weights to torso',
        'Squeeze shoulder blades',
        'Lower with control'
    ],
    ARRAY['dumbbells'],
    ARRAY['back', 'biceps'],
    ARRAY['functional', 'strength-endurance']
);

-- Dumbbell Overhead Press
SELECT insert_exercise(
    '44444444-5555-6666-7777-888888888888',
    'Dumbbell Overhead Press',
    'strength',
    'Pressing dumbbells overhead while standing',
    ARRAY[
        'Hold weights at shoulders',
        'Brace core',
        'Press straight overhead',
        'Lock out arms',
        'Lower with control'
    ],
    ARRAY['dumbbells'],
    ARRAY['shoulders', 'triceps', 'core'],
    ARRAY['crossfit', 'functional', 'strength-endurance']
);

-- Russian Twists
SELECT insert_exercise(
    '55555555-6666-7777-8888-999999999999',
    'Russian Twists',
    'strength',
    'Rotational core exercise targeting obliques',
    ARRAY[
        'Sit with knees bent',
        'Lean back slightly',
        'Rotate torso side to side',
        'Touch ball to ground',
        'Keep feet off ground'
    ],
    ARRAY['medicine-ball'],
    ARRAY['core', 'obliques'],
    ARRAY['core', 'functional']
);

-- Plank
SELECT insert_exercise(
    '66666666-7777-8888-9999-aaaaaaaaaaaa',
    'Plank',
    'strength',
    'Isometric hold strengthening core and stabilizer muscles',
    ARRAY[
        'Start in push-up position',
        'Lower to forearms',
        'Keep body straight',
        'Engage core',
        'Hold position'
    ],
    ARRAY['bodyweight'],
    ARRAY['core', 'abs'],
    ARRAY['beginner-friendly', 'core', 'functional']
);

-- Dips
SELECT insert_exercise(
    '77777777-8888-9999-aaaa-bbbbbbbbbbbb',
    'Dips',
    'strength',
    'Bodyweight exercise targeting triceps and chest',
    ARRAY[
        'Support body on parallel bars',
        'Lower body down',
        'Push back to start',
        'Keep body upright',
        'Control the movement'
    ],
    ARRAY['bodyweight'],
    ARRAY['triceps', 'chest', 'shoulders'],
    ARRAY['functional', 'strength-endurance']
);

-- Running
SELECT insert_exercise(
    '88888888-9999-aaaa-bbbb-cccccccccccc',
    'Running',
    'cardio',
    'Running at various intensities for cardiovascular conditioning',
    ARRAY[
        'Maintain proper running form',
        'Land on mid-foot',
        'Keep cadence high',
        'Breathe rhythmically',
        'Vary pace as needed'
    ],
    ARRAY['bodyweight'],
    ARRAY['legs', 'core'],
    ARRAY['hyrox', 'conditioning', 'beginner-friendly']
);

-- Sandbag Carry
SELECT insert_exercise(
    '99999999-aaaa-bbbb-cccc-dddddddddddd',
    'Sandbag Carry',
    'strength',
    'Carrying heavy sandbag for distance or time',
    ARRAY[
        'Pick up sandbag',
        'Hold close to body',
        'Walk with good posture',
        'Keep core engaged',
        'Set down safely'
    ],
    ARRAY['sled'],
    ARRAY['core', 'legs', 'grip'],
    ARRAY['hyrox', 'functional', 'strength-endurance']
);

-- Barbell Back Squat
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000001',
    'Barbell Back Squat',
    'strength',
    'Fundamental squatting movement with barbell on back',
    ARRAY[
        'Position barbell on upper back',
        'Stand with feet shoulder-width',
        'Descend by sitting back',
        'Drive through heels to stand',
        'Keep chest up throughout'
    ],
    ARRAY['barbell'],
    ARRAY['legs', 'glutes', 'core'],
    ARRAY['crossfit', 'functional', 'strength-endurance']
);

-- Barbell Deadlift
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000002',
    'Barbell Deadlift',
    'strength',
    'Hip hinge movement lifting barbell from ground to standing',
    ARRAY[
        'Stand with feet hip-width',
        'Grip barbell with hands outside legs',
        'Hinge at hips and knees',
        'Drive through heels and hips',
        'Stand tall with shoulders back'
    ],
    ARRAY['barbell'],
    ARRAY['back', 'glutes', 'hamstrings', 'grip'],
    ARRAY['crossfit', 'functional', 'strength-endurance']
);

-- Barbell Bench Press
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000003',
    'Barbell Bench Press',
    'strength',
    'Classic upper body pressing movement with barbell',
    ARRAY[
        'Lie on bench with barbell racked',
        'Grip barbell slightly wider than shoulders',
        'Lower bar to chest with control',
        'Press bar straight up',
        'Lock out arms at top'
    ],
    ARRAY['barbell'],
    ARRAY['chest', 'shoulders', 'triceps'],
    ARRAY['functional', 'strength-endurance']
);

-- Barbell Thrusters
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000004',
    'Barbell Thrusters',
    'strength',
    'Combination front squat to overhead press with barbell',
    ARRAY[
        'Hold barbell in front rack position',
        'Perform front squat',
        'Drive up explosively',
        'Press barbell overhead',
        'Lower to front rack position'
    ],
    ARRAY['barbell'],
    ARRAY['legs', 'shoulders', 'core', 'full-body'],
    ARRAY['crossfit', 'functional', 'conditioning']
);

-- Barbell Bent-over Rows
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000005',
    'Barbell Bent-over Rows',
    'strength',
    'Pulling movement with barbell targeting back muscles',
    ARRAY[
        'Hinge at hips holding barbell',
        'Keep back straight and core tight',
        'Pull barbell to lower chest',
        'Squeeze shoulder blades together',
        'Lower with control'
    ],
    ARRAY['barbell'],
    ARRAY['back', 'biceps', 'core'],
    ARRAY['functional', 'strength-endurance']
);

-- Barbell Overhead Press
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000006',
    'Barbell Overhead Press',
    'strength',
    'Standing overhead press with barbell',
    ARRAY[
        'Hold barbell at shoulder height',
        'Stand with feet hip-width',
        'Brace core and glutes',
        'Press barbell straight overhead',
        'Lower to starting position'
    ],
    ARRAY['barbell'],
    ARRAY['shoulders', 'triceps', 'core'],
    ARRAY['crossfit', 'functional', 'strength-endurance']
);

-- Clean and Jerk
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000007',
    'Clean and Jerk',
    'strength',
    'Olympic weightlifting movement from ground to overhead',
    ARRAY[
        'Deadlift barbell to hips',
        'Explosively extend hips and knees',
        'Pull barbell to front rack',
        'Dip and drive to press overhead',
        'Lock out arms and stabilize'
    ],
    ARRAY['barbell'],
    ARRAY['full-body', 'legs', 'shoulders', 'back'],
    ARRAY['crossfit', 'advanced', 'power', 'competition']
);

-- Barbell Front Squat
SELECT insert_exercise(
    'b0000000-0000-0000-0000-000000000008',
    'Barbell Front Squat',
    'strength',
    'Squat with barbell held in front rack position',
    ARRAY[
        'Position barbell in front rack',
        'Keep elbows up and chest proud',
        'Descend into squat position',
        'Drive through heels to stand',
        'Maintain upright torso'
    ],
    ARRAY['barbell'],
    ARRAY['legs', 'glutes', 'core'],
    ARRAY['crossfit', 'functional', 'advanced']
);

-- Assault Bike
SELECT insert_exercise(
    'a5000000-0000-0000-0000-000000000001',
    'Assault Bike',
    'cardio',
    'High-intensity cardio using air resistance bike with moving handles',
    ARRAY[
        'Sit on bike with feet on pedals',
        'Grip moving handles',
        'Push and pull with arms',
        'Pedal with legs simultaneously',
        'Maintain steady breathing'
    ],
    ARRAY['assault-bike'],
    ARRAY['legs', 'core', 'full-body'],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'advanced']
);

-- Clean up helper function
DROP FUNCTION insert_exercise;

COMMIT;
