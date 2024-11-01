import { useState, useEffect } from 'react';

export default function CBTExercises({ currentMood }) {
  const [exercise, setExercise] = useState(null);

  useEffect(() => {
    const fetchExercise = async () => {
      const response = await fetch(`/api/cbt?mood=${currentMood}`);
      const data = await response.json();
      setExercise(data.exercise);
    };

    if (currentMood) {
      fetchExercise();
    }
  }, [currentMood]);

  if (!exercise) return null;

  return (
    <div className="my-8 p-6 bg-blue-100 rounded-lg">
      <h2 className="text-2xl font-semibold mb-4">CBT Exercise</h2>
      <p className="mb-4">{exercise.description}</p>
      <ol className="list-decimal list-inside">
        {exercise.steps.map((step, index) => (
          <li key={index} className="mb-2">{step}</li>
        ))}
      </ol>
    </div>
  );
}