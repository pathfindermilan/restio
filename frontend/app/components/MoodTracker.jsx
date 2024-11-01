import { useState } from 'react';

export default function MoodTracker({ onMoodUpdate }) {
  const [mood, setMood] = useState(5);

  const handleMoodChange = (e) => {
    const newMood = parseInt(e.target.value, 10);
    setMood(newMood);
    onMoodUpdate(newMood);
  };

  const getMoodEmoji = (moodValue) => {
    if (moodValue <= 2) return '😢';
    if (moodValue <= 4) return '😕';
    if (moodValue <= 6) return '😐';
    if (moodValue <= 8) return '🙂';
    return '😄';
  };

  return (
    <div className="my-8 p-6 bg-gray-100 rounded-lg">
      <h2 className="text-2xl font-semibold mb-4">How are you feeling today?</h2>
      <div className="flex items-center mb-4">
        <input
          type="range"
          min="1"
          max="10"
          value={mood}
          onChange={handleMoodChange}
          className="w-full mr-4"
        />
        <span className="text-4xl">{getMoodEmoji(mood)}</span>
      </div>
      <p className="text-center">Your mood: {mood}/10</p>
    </div>
  );
}