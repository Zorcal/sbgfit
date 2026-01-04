BEGIN;

-- Global exercises.
INSERT INTO sbgfit.exercises (
    external_id,
    name,
    category,
    equipment_types,
    primary_muscles,
    description,
    instructions,
    tags
) VALUES
-- Burpees
(
    '01234567-89ab-cdef-0123-456789abcdef',
    'Burpees',
    'cardio',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['full-body']::primary_muscle[],
    'From standing, squat down, jump back to plank, do a push-up, jump feet back to squat, then jump up with arms overhead',
    ARRAY[
        'Start standing',
        'Squat down hands on ground',
        'Jump back to plank',
        'Do push-up',
        'Jump feet to squat',
        'Jump up arms overhead'
    ],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'functional', 'competition']::exercise_tag[]
),

-- Kettlebell Swings
(
    '11111111-1111-1111-1111-111111111111',
    'Kettlebell Swings',
    'strength',
    ARRAY['kettlebell']::exercise_equipment[],
    ARRAY['glutes', 'hamstrings', 'core']::primary_muscle[],
    'Hip-hinge movement swinging kettlebell from between legs to chest height',
    ARRAY[
        'Stand with kettlebell',
        'Hinge at hips grab bell',
        'Drive hips forward swing up',
        'Let bell swing back',
        'Repeat motion'
    ],
    ARRAY['crossfit', 'power', 'functional']::exercise_tag[]
),

-- Rowing
(
    '22222222-2222-2222-2222-222222222222',
    'Rowing',
    'cardio',
    ARRAY['rowing-machine']::exercise_equipment[],
    ARRAY['back', 'legs', 'core']::primary_muscle[],
    'Full-body cardio movement on rowing machine',
    ARRAY[
        'Sit on machine feet strapped',
        'Grab handle',
        'Push legs lean back',
        'Pull to chest',
        'Reverse movement'
    ],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'core']::exercise_tag[]
),

-- Ski Erg
(
    '33333333-3333-3333-3333-333333333333',
    'Ski Erg',
    'cardio',
    ARRAY['ski-erg']::exercise_equipment[],
    ARRAY['shoulders', 'core', 'legs']::primary_muscle[],
    'Upper body cardio movement mimicking cross-country skiing',
    ARRAY[
        'Stand feet hip-width',
        'Grab handles overhead',
        'Pull down skiing motion',
        'Return overhead',
        'Maintain rhythm'
    ],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'core']::exercise_tag[]
),

-- Wall Balls
(
    '44444444-4444-4444-4444-444444444444',
    'Wall Balls',
    'strength',
    ARRAY['medicine-ball']::exercise_equipment[],
    ARRAY['legs', 'shoulders', 'core']::primary_muscle[],
    'Squat and throw medicine ball to target on wall',
    ARRAY[
        'Hold ball at chest',
        'Squat keeping chest up',
        'Drive up throw to target',
        'Catch ball squat again',
        'Repeat continuously'
    ],
    ARRAY['crossfit', 'power', 'functional']::exercise_tag[]
),

-- Farmers Walk
(
    '55555555-5555-5555-5555-555555555555',
    'Farmers Walk',
    'strength',
    ARRAY['dumbbells']::exercise_equipment[],
    ARRAY['grip', 'core', 'legs']::primary_muscle[],
    'Walk while carrying heavy weights in each hand',
    ARRAY[
        'Pick up weights',
        'Stand tall shoulders back',
        'Walk maintaining posture',
        'Keep core tight',
        'Set down safely'
    ],
    ARRAY['hyrox', 'strength-endurance', 'functional']::exercise_tag[]
),

-- Sled Push
(
    '66666666-6666-6666-6666-666666666666',
    'Sled Push',
    'strength',
    ARRAY['sled']::exercise_equipment[],
    ARRAY['legs', 'glutes', 'core']::primary_muscle[],
    'Push weighted sled across floor',
    ARRAY[
        'Hands on handles',
        'Lean forward straight back',
        'Drive with legs forward',
        'Maintain pace',
        'Keep core engaged'
    ],
    ARRAY['hyrox', 'strength-endurance', 'functional']::exercise_tag[]
),

-- Sled Pull
(
    '77777777-7777-7777-7777-777777777777',
    'Sled Pull',
    'strength',
    ARRAY['sled']::exercise_equipment[],
    ARRAY['back', 'biceps', 'core']::primary_muscle[],
    'Pull weighted sled toward you',
    ARRAY[
        'Grab rope or handles',
        'Lean back slightly',
        'Pull hand over hand',
        'Reset position',
        'Maintain rhythm'
    ],
    ARRAY['hyrox', 'strength-endurance', 'functional']::exercise_tag[]
),

-- Box Jumps
(
    '88888888-8888-8888-8888-888888888888',
    'Box Jumps',
    'plyometric',
    ARRAY['box']::exercise_equipment[],
    ARRAY['legs', 'glutes']::primary_muscle[],
    'Jump onto elevated box or platform',
    ARRAY[
        'Stand in front of box',
        'Swing arms bend knees',
        'Jump up land softly',
        'Stand upright on box',
        'Step down safely'
    ],
    ARRAY['crossfit', 'power', 'plyometric']::exercise_tag[]
),

-- Lunges
(
    '99999999-9999-9999-9999-999999999999',
    'Lunges',
    'strength',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['legs', 'glutes']::primary_muscle[],
    'Single-leg strength movement stepping forward into lunge position',
    ARRAY[
        'Stand feet hip-width',
        'Step forward to lunge',
        'Lower back knee down',
        'Push through front heel',
        'Alternate or complete side'
    ],
    ARRAY['crossfit', 'hyrox', 'beginner-friendly', 'functional']::exercise_tag[]
),

-- Pull-ups
(
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    'Pull-ups',
    'strength',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['back', 'biceps']::primary_muscle[],
    'Hanging from a bar and pulling body up until chin clears the bar',
    ARRAY[
        'Hang from pull-up bar',
        'Pull body up',
        'Chin over bar',
        'Lower with control',
        'Repeat'
    ],
    ARRAY['crossfit', 'functional', 'beginner-friendly']::exercise_tag[]
),

-- Push-ups
(
    'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
    'Push-ups',
    'strength',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['chest', 'shoulders', 'triceps']::primary_muscle[],
    'Classic bodyweight exercise targeting chest, shoulders, and triceps',
    ARRAY[
        'Start in plank position',
        'Lower chest to ground',
        'Push back to start',
        'Keep body straight',
        'Repeat'
    ],
    ARRAY['crossfit', 'beginner-friendly', 'functional']::exercise_tag[]
),

-- Dumbbell Deadlifts
(
    'cccccccc-cccc-cccc-cccc-cccccccccccc',
    'Dumbbell Deadlifts',
    'strength',
    ARRAY['dumbbells']::exercise_equipment[],
    ARRAY['back', 'glutes', 'hamstrings']::primary_muscle[],
    'Hip hinge movement lifting dumbbells from ground to standing position',
    ARRAY[
        'Stand with feet hip-width',
        'Hinge at hips',
        'Grab weights',
        'Drive hips forward',
        'Stand tall'
    ],
    ARRAY['crossfit', 'functional', 'strength-endurance']::exercise_tag[]
),

-- Air Squats
(
    'dddddddd-dddd-dddd-dddd-dddddddddddd',
    'Air Squats',
    'strength',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['legs', 'glutes']::primary_muscle[],
    'Bodyweight squat focusing on proper hip and knee movement',
    ARRAY[
        'Stand with feet shoulder-width',
        'Lower hips back and down',
        'Keep chest up',
        'Drive through heels',
        'Return to standing'
    ],
    ARRAY['crossfit', 'beginner-friendly', 'functional']::exercise_tag[]
),

-- Dumbbell Thrusters
(
    'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
    'Dumbbell Thrusters',
    'strength',
    ARRAY['dumbbells']::exercise_equipment[],
    ARRAY['legs', 'shoulders', 'core']::primary_muscle[],
    'Combination squat to overhead press with dumbbells',
    ARRAY[
        'Hold weights at shoulders',
        'Squat down',
        'Drive up explosively',
        'Press weights overhead',
        'Lower to shoulders'
    ],
    ARRAY['crossfit', 'functional', 'conditioning']::exercise_tag[]
),

-- Double Unders
(
    'ffffffff-ffff-ffff-ffff-ffffffffffff',
    'Double Unders',
    'cardio',
    ARRAY['jump-rope']::exercise_equipment[],
    ARRAY['legs', 'core']::primary_muscle[],
    'Jump rope where rope passes under feet twice per jump',
    ARRAY[
        'Hold rope handles',
        'Jump higher than normal',
        'Spin rope faster',
        'Land on balls of feet',
        'Keep rhythm consistent'
    ],
    ARRAY['crossfit', 'conditioning', 'advanced']::exercise_tag[]
),

-- Mountain Climbers
(
    '10101010-1010-1010-1010-101010101010',
    'Mountain Climbers',
    'cardio',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['core', 'legs']::primary_muscle[],
    'Dynamic plank position with alternating knee drives',
    ARRAY[
        'Start in plank position',
        'Drive right knee to chest',
        'Switch legs quickly',
        'Keep hips level',
        'Maintain fast pace'
    ],
    ARRAY['crossfit', 'conditioning', 'core']::exercise_tag[]
),

-- Turkish Get-ups
(
    '11111111-2222-3333-4444-555555555555',
    'Turkish Get-ups',
    'strength',
    ARRAY['kettlebell']::exercise_equipment[],
    ARRAY['core', 'shoulders', 'full-body']::primary_muscle[],
    'Complex movement from lying to standing while holding weight overhead',
    ARRAY[
        'Lie on back with weight up',
        'Roll to elbow',
        'Push to hand',
        'Bridge hips up',
        'Stand up slowly'
    ],
    ARRAY['functional', 'advanced', 'core']::exercise_tag[]
),

-- Dumbbell Bench Press
(
    '22222222-3333-4444-5555-666666666666',
    'Dumbbell Bench Press',
    'strength',
    ARRAY['dumbbells']::exercise_equipment[],
    ARRAY['chest', 'shoulders', 'triceps']::primary_muscle[],
    'Upper body pressing movement with dumbbells for chest development',
    ARRAY[
        'Lie on bench',
        'Lower weights to chest',
        'Press up explosively',
        'Keep back flat',
        'Control the weight'
    ],
    ARRAY['functional', 'strength-endurance']::exercise_tag[]
),

-- Dumbbell Bent-over Rows
(
    '33333333-4444-5555-6666-777777777777',
    'Dumbbell Bent-over Rows',
    'strength',
    ARRAY['dumbbells']::exercise_equipment[],
    ARRAY['back', 'biceps']::primary_muscle[],
    'Pulling movement with dumbbells targeting back muscles and posterior chain',
    ARRAY[
        'Hinge at hips',
        'Hold weights with arms extended',
        'Pull weights to torso',
        'Squeeze shoulder blades',
        'Lower with control'
    ],
    ARRAY['functional', 'strength-endurance']::exercise_tag[]
),

-- Dumbbell Overhead Press
(
    '44444444-5555-6666-7777-888888888888',
    'Dumbbell Overhead Press',
    'strength',
    ARRAY['dumbbells']::exercise_equipment[],
    ARRAY['shoulders', 'triceps', 'core']::primary_muscle[],
    'Pressing dumbbells overhead while standing',
    ARRAY[
        'Hold weights at shoulders',
        'Brace core',
        'Press straight overhead',
        'Lock out arms',
        'Lower with control'
    ],
    ARRAY['crossfit', 'functional', 'strength-endurance']::exercise_tag[]
),

-- Russian Twists
(
    '55555555-6666-7777-8888-999999999999',
    'Russian Twists',
    'strength',
    ARRAY['medicine-ball']::exercise_equipment[],
    ARRAY['core', 'obliques']::primary_muscle[],
    'Rotational core exercise targeting obliques',
    ARRAY[
        'Sit with knees bent',
        'Lean back slightly',
        'Rotate torso side to side',
        'Touch ball to ground',
        'Keep feet off ground'
    ],
    ARRAY['core', 'functional']::exercise_tag[]
),

-- Plank
(
    '66666666-7777-8888-9999-aaaaaaaaaaaa',
    'Plank',
    'strength',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['core', 'abs']::primary_muscle[],
    'Isometric hold strengthening core and stabilizer muscles',
    ARRAY[
        'Start in push-up position',
        'Lower to forearms',
        'Keep body straight',
        'Engage core',
        'Hold position'
    ],
    ARRAY['beginner-friendly', 'core', 'functional']::exercise_tag[]
),

-- Dips
(
    '77777777-8888-9999-aaaa-bbbbbbbbbbbb',
    'Dips',
    'strength',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['triceps', 'chest', 'shoulders']::primary_muscle[],
    'Bodyweight exercise targeting triceps and chest',
    ARRAY[
        'Support body on parallel bars',
        'Lower body down',
        'Push back to start',
        'Keep body upright',
        'Control the movement'
    ],
    ARRAY['functional', 'strength-endurance']::exercise_tag[]
),

-- High-intensity Interval Running
(
    '88888888-9999-aaaa-bbbb-cccccccccccc',
    'Running',
    'cardio',
    ARRAY['bodyweight']::exercise_equipment[],
    ARRAY['legs', 'core']::primary_muscle[],
    'Running at various intensities for cardiovascular conditioning',
    ARRAY[
        'Maintain proper running form',
        'Land on mid-foot',
        'Keep cadence high',
        'Breathe rhythmically',
        'Vary pace as needed'
    ],
    ARRAY['hyrox', 'conditioning', 'beginner-friendly']::exercise_tag[]
),

-- Sandbag Carry
(
    '99999999-aaaa-bbbb-cccc-dddddddddddd',
    'Sandbag Carry',
    'strength',
    ARRAY['sled']::exercise_equipment[],
    ARRAY['core', 'legs', 'grip']::primary_muscle[],
    'Carrying heavy sandbag for distance or time',
    ARRAY[
        'Pick up sandbag',
        'Hold close to body',
        'Walk with good posture',
        'Keep core engaged',
        'Set down safely'
    ],
    ARRAY['hyrox', 'functional', 'strength-endurance']::exercise_tag[]
),

-- Barbell Back Squat
(
    'b0000000-0000-0000-0000-000000000001',
    'Barbell Back Squat',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['legs', 'glutes', 'core']::primary_muscle[],
    'Fundamental squatting movement with barbell on back',
    ARRAY[
        'Position barbell on upper back',
        'Stand with feet shoulder-width',
        'Descend by sitting back',
        'Drive through heels to stand',
        'Keep chest up throughout'
    ],
    ARRAY['crossfit', 'functional', 'strength-endurance']::exercise_tag[]
),

-- Barbell Deadlift
(
    'b0000000-0000-0000-0000-000000000002',
    'Barbell Deadlift',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['back', 'glutes', 'hamstrings', 'grip']::primary_muscle[],
    'Hip hinge movement lifting barbell from ground to standing',
    ARRAY[
        'Stand with feet hip-width',
        'Grip barbell with hands outside legs',
        'Hinge at hips and knees',
        'Drive through heels and hips',
        'Stand tall with shoulders back'
    ],
    ARRAY['crossfit', 'functional', 'strength-endurance']::exercise_tag[]
),

-- Barbell Bench Press
(
    'b0000000-0000-0000-0000-000000000003',
    'Barbell Bench Press',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['chest', 'shoulders', 'triceps']::primary_muscle[],
    'Classic upper body pressing movement with barbell',
    ARRAY[
        'Lie on bench with barbell racked',
        'Grip barbell slightly wider than shoulders',
        'Lower bar to chest with control',
        'Press bar straight up',
        'Lock out arms at top'
    ],
    ARRAY['functional', 'strength-endurance']::exercise_tag[]
),

-- Barbell Thrusters
(
    'b0000000-0000-0000-0000-000000000004',
    'Barbell Thrusters',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['legs', 'shoulders', 'core', 'full-body']::primary_muscle[],
    'Combination front squat to overhead press with barbell',
    ARRAY[
        'Hold barbell in front rack position',
        'Perform front squat',
        'Drive up explosively',
        'Press barbell overhead',
        'Lower to front rack position'
    ],
    ARRAY['crossfit', 'functional', 'conditioning']::exercise_tag[]
),

-- Barbell Bent-over Rows
(
    'b0000000-0000-0000-0000-000000000005',
    'Barbell Bent-over Rows',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['back', 'biceps', 'core']::primary_muscle[],
    'Pulling movement with barbell targeting back muscles',
    ARRAY[
        'Hinge at hips holding barbell',
        'Keep back straight and core tight',
        'Pull barbell to lower chest',
        'Squeeze shoulder blades together',
        'Lower with control'
    ],
    ARRAY['functional', 'strength-endurance']::exercise_tag[]
),

-- Barbell Overhead Press
(
    'b0000000-0000-0000-0000-000000000006',
    'Barbell Overhead Press',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['shoulders', 'triceps', 'core']::primary_muscle[],
    'Standing overhead press with barbell',
    ARRAY[
        'Hold barbell at shoulder height',
        'Stand with feet hip-width',
        'Brace core and glutes',
        'Press barbell straight overhead',
        'Lower to starting position'
    ],
    ARRAY['crossfit', 'functional', 'strength-endurance']::exercise_tag[]
),

-- Clean and Jerk
(
    'b0000000-0000-0000-0000-000000000007',
    'Clean and Jerk',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['full-body', 'legs', 'shoulders', 'back']::primary_muscle[],
    'Olympic weightlifting movement from ground to overhead',
    ARRAY[
        'Deadlift barbell to hips',
        'Explosively extend hips and knees',
        'Pull barbell to front rack',
        'Dip and drive to press overhead',
        'Lock out arms and stabilize'
    ],
    ARRAY['crossfit', 'advanced', 'power', 'competition']::exercise_tag[]
),

-- Barbell Front Squat
(
    'b0000000-0000-0000-0000-000000000008',
    'Barbell Front Squat',
    'strength',
    ARRAY['barbell']::exercise_equipment[],
    ARRAY['legs', 'glutes', 'core']::primary_muscle[],
    'Squat with barbell held in front rack position',
    ARRAY[
        'Position barbell in front rack',
        'Keep elbows up and chest proud',
        'Descend into squat position',
        'Drive through heels to stand',
        'Maintain upright torso'
    ],
    ARRAY['crossfit', 'functional', 'advanced']::exercise_tag[]
),

-- Assault Bike
(
    'a5000000-0000-0000-0000-000000000001',
    'Assault Bike',
    'cardio',
    ARRAY['assault-bike']::exercise_equipment[],
    ARRAY['legs', 'core', 'full-body']::primary_muscle[],
    'High-intensity cardio using air resistance bike with moving handles',
    ARRAY[
        'Sit on bike with feet on pedals',
        'Grip moving handles',
        'Push and pull with arms',
        'Pedal with legs simultaneously',
        'Maintain steady breathing'
    ],
    ARRAY['crossfit', 'hyrox', 'conditioning', 'advanced']::exercise_tag[]
)

ON CONFLICT (external_id) DO UPDATE SET
    name = EXCLUDED.name,
    category = EXCLUDED.category,
    equipment_types = EXCLUDED.equipment_types,
    primary_muscles = EXCLUDED.primary_muscles,
    description = EXCLUDED.description,
    instructions = EXCLUDED.instructions,
    tags = EXCLUDED.tags,
    updated_at = CURRENT_TIMESTAMP;

COMMIT;
