DROP TABLE IF EXISTS exercise;

EXEC SQL INCLUDE '../schema/exercises/exercises_schema.sql'

-- generic exercises
INSERT INTO exercise (name, description)
VALUES 
	('Pushup', 'How to do a push-up 
              - Start in a high plank position with your hands slightly wider than shoulder-width apart.
              - Extend your legs and place your feet hip-width apart.
              - Squeeze your shoulder blades and brace your core.
              - Keeping your elbows close to your sides, bend them to lower your chest towards the floor.
              - Pause, then push your hands into the floor to return to the starting position.'),
	('Pullup', 'To perform a pull-up, 
                - grasp a pull-up bar with a slightly wider than shoulder-width grip, 
                - palms facing away from you, then engage your core and pull yourself up by pulling your elbows down towards your hips, 
                - aiming to bring your chest to the bar, before slowly lowering yourself back down to the starting position 
                - while keeping your body straight and maintaining control throughout the movement; 
                - focus on using your back muscles (lats) primarily, with assistance from your biceps.'),
    ('Squat', 'How to do a squat 
               - Stand with your feet shoulder-width apart and toes pointed slightly out
               - Keep your chest up and engage your core
               - Shift your weight onto your heels
               - Push your hips back into a sitting position
               - Bend your knees until your thighs are parallel to the floor
               - Push back up through your feet to return to the starting position')
RETURNING *;