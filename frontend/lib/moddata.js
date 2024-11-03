import {Music, Dumbbell, Laugh, Video } from 'lucide-react';
export const options = [
    { id: 'speech', label: 'Speech', icon: <Music className="w-4 h-4 md:w-5 md:h-5" /> },
    { id: 'exercise', label: 'Exercise', icon: <Dumbbell className="w-4 h-4 md:w-5 md:h-5" /> },
    { id: 'joke', label: 'Joke', icon: <Laugh className="w-4 h-4 md:w-5 md:h-5" /> },
    { id: 'video', label: 'Video', icon: <Video className="w-4 h-4 md:w-5 md:h-5" /> }
  ];
  

  export const moodScale = [
    { level: 1, emoji: "😢", description: "Very Sad" },
    { level: 2, emoji: "😥", description: "Sad" },
    { level: 3, emoji: "😕", description: "Slightly Down" },
    { level: 4, emoji: "😐", description: "Neutral-Low" },
    { level: 5, emoji: "😊", description: "Okay" },
    { level: 6, emoji: "🙂", description: "Slightly Good" },
    { level: 7, emoji: "😃", description: "Good" },
    { level: 8, emoji: "😄", description: "Very Good" },
    { level: 9, emoji: "🤗", description: "Great" },
    { level: 10, emoji: "🥳", description: "Excellent" }
  ];
  