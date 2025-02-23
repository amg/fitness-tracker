import { useState, useEffect, useCallback } from 'react';
import { Box } from "@mui/material"
import Grid from '@mui/material/Grid2';
import GenericError from '../common/GenericError';
import ExerciseListCard, { Exercise } from '../common/ExerciseListCard';
import { useGlobalAuthContext } from '../../helpers/authContext';

const nodeApiBaseUrl = window.env.NODE_API_BASE_URL

function RecordExercisePage() {
    const authContext = useGlobalAuthContext();

    const [data, setData] = useState<Array<Exercise> | null>(null)

    const fetchExercises = useCallback(() => {
        fetch(nodeApiBaseUrl + '/node/exercises', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
            mode: "cors",
            credentials: "include",
        })
            .then(response => {
                if (response.status === 401) {
                    authContext.setAuthState(null)
                    throw new Error(`Response status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setData(data);
            })
            .catch(error => {
                setData(null);
                console.info('Caught error:', error);
            });
    }, [authContext])

    useEffect(() => {
        fetchExercises()
    }, [fetchExercises])

    return <>
        <h2>Select the exercise</h2>
        <br />
        {
            data ? (
                <Box sx={{ flexGrow: 1 }}>
                    <h3>All available</h3>
                    <Grid container direction={"row"} spacing={{ xs: 2, md: 3 }} columns={{ xs: 4, sm: 8, md: 12 }}
                    sx={{ 
                        alignItems: "stretch",
                        justifyContent: "center",
                    }}>
                        {data.map((exercise, index) => (
                            <Grid key={index} size={{ xs: 2, sm: 4, md: 4 }}>
                                <ExerciseListCard {...exercise} onClick={(id: number) => {
                                    console.log("Clicked on exercise: %d", id)
                                }} />
                            </Grid>
                        ))}
                    </Grid>
                </Box>
            ) : (
                <GenericError message={"Couldn't load the exercises"} />
            )
        }
    </>
}

// function mockExerciseData() {
//     return [
//         {
//             id: 1,
//             title: "Push ups",
//             description: `
//             How to do a push-up 
//               - Start in a high plank position with your hands slightly wider than shoulder-width apart.
//               - Extend your legs and place your feet hip-width apart.
//               - Squeeze your shoulder blades and brace your core.
//               - Keeping your elbows close to your sides, bend them to lower your chest towards the floor.
//               - Pause, then push your hands into the floor to return to the starting position.
//             `,
//             onClick: null,
//         },
//         {
//             id: 2,
//             title: "Pull up",
//             description: `
//                 To perform a pull-up, grasp a pull-up bar with a slightly wider than shoulder-width grip, 
//                 palms facing away from you, then engage your core and pull yourself up by pulling your elbows down towards your hips, 
//                 aiming to bring your chest to the bar, before slowly lowering yourself back down to the starting position 
//                 while keeping your body straight and maintaining control throughout the movement; 
//                 focus on using your back muscles (lats) primarily, with assistance from your biceps.
//             `,
//             onClick: null,
//         },
//         {
//             id: 3,
//             title: "Squats",
//             description: `
//             How to do a squat 
//                - Stand with your feet shoulder-width apart and toes pointed slightly out
//                - Keep your chest up and engage your core
//                - Shift your weight onto your heels
//                - Push your hips back into a sitting position
//                - Bend your knees until your thighs are parallel to the floor
//                - Push back up through your feet to return to the starting position
//             `,
//             onClick: null,
//         },
//     ];
// }

export default RecordExercisePage