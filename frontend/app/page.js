'use client'
import { useState, useEffect } from 'react';
import Head from 'next/head';
import MoodTracker from '@/app/components/MoodTracker';
import CBTExercises from '@/app/components/CBTExercises';
import MindfulnessExercises from '@/app/components/MindfulnessExercises';
import PersonalizedPlan from '@/app/components/PersonalizedPlan';
import TherapistConnect from '@/app/components/TherapistConnect';

export default function Home() {
  const [currentMood, setCurrentMood] = useState(null);
  const [moodHistory, setMoodHistory] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchMoodHistory = async () => {
      try {
        const response = await fetch('/api/mood');
        if (!response.ok) {
          throw new Error('Failed to fetch mood history');
        }
        const data = await response.json();
        setMoodHistory(data.moodHistory);
      } catch (err) {
        setError(err.message);
      } finally {
        setIsLoading(false);
      }
    };

    fetchMoodHistory();
  }, []);

  const handleMoodUpdate = async (mood) => {
    setCurrentMood(mood);
    const newMoodEntry = { mood, timestamp: new Date() };
    
    try {
      const response = await fetch('/api/mood', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newMoodEntry),
      });

      if (!response.ok) {
        throw new Error('Failed to save mood');
      }

      setMoodHistory([...moodHistory, newMoodEntry]);
    } catch (err) {
      setError(err.message);
    }
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div className="container mx-auto px-4">
      <Head>
        <title>Mood Sync</title>
        <meta name="description" content="Track and improve your mental health" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className="py-8">
        <h1 className="text-4xl font-bold mb-8 text-center">Mood Sync</h1>
        
        <MoodTracker onMoodUpdate={handleMoodUpdate} />
        
        {currentMood && (
          <CBTExercises currentMood={currentMood} />
        )}
        
        <MindfulnessExercises />
        
        {moodHistory.length > 0 && (
          <PersonalizedPlan moodHistory={moodHistory} />
        )}
        
        <TherapistConnect />
      </main>

      <footer className="mt-8 text-center text-gray-500">
        <p>© 2023 Mood Sync. All rights reserved.</p>
      </footer>
    </div>
  );
}