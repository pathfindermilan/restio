import { useState } from 'react';

export default function MindfulnessExercises() {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentExercise, setCurrentExercise] = useState(null);

  const exercises = [
    { id: 1, name: "5-Minute Breathing", duration: 300 },
    { id: 2, name: "Body Scan Meditation", duration: 600 },
    { id: 3, name: "Loving-Kindness Meditation", duration: 900 },
  ];

  const handleExerciseStart = (exercise) => {
    setCurrentExercise(exercise);
    setIsPlaying(true);
    // In a real app, you would start playing the audio here
  };

  return (
    <div className="my-8">
      <h2 className="text-2xl font-semibold mb-4">Mindfulness Exercises</h2>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {exercises.map((exercise) => (
          <div key={exercise.id} className="p-4 border rounded-lg">
            <h3 className="font-semibold">{exercise.name}</h3>
            <p>{exercise.duration / 60} minutes</p>
            <button
              onClick={() => handleExerciseStart(exercise)}
              className="mt-2 px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
            >
              Start
            </button>
          </div>
        ))}
      </div>
      {isPlaying && (
        <div className="mt-4 p-4 bg-gray-100 rounded-lg">
          <p>Now playing: {currentExercise.name}</p>
          {/* Add audio player controls here */}
        </div>
      )}
    </div>
  );
}