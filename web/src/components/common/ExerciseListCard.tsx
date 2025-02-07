import { Card, CardActionArea, CardContent, CardHeader, CardMedia, Typography } from "@mui/material";

export class Exercise {
    constructor(public id: number, public title: string, public description: string, public onClick: ((id: number) => void) | null) { }
}

export default function ExerciseListCard(
    exercise: Exercise
) {
    return (
        <Card sx={{ height: '100%' }}>
            <CardActionArea 
                onClick={() => exercise.onClick ? exercise.onClick(exercise.id) : undefined}
                sx={{ height: '100%' }}>
                <CardHeader
                    avatar={
                        <CardMedia
                            component="img"
                            sx={{ width: 64 }}
                            alt="No image"
                        />
                    }
                    title={<Typography align="left">{exercise.title}</Typography>}
                />
                <CardContent sx={{ height: '100%' }}>
                    <Typography variant="body2" sx={{
                        color: 'text.secondary',
                    }}>
                        {exercise.description}
                    </Typography>
                </CardContent>
            </CardActionArea>
        </Card>
    );
}