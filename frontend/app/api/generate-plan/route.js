import { NextResponse } from 'next/server';

export async function POST(request) {
  const { moodHistory } = await request.json();

  // In a real app, you'd use AI to analyze the mood history and generate a plan
  const plan = {
    recommendations: [
      "Practice mindfulness meditation for 10 minutes each morning",
      "Engage in moderate exercise for 30 minutes, 3 times a week",
      "Keep a gratitude journal, writing down 3 things you're grateful for each day",
      "Limit screen time in the evening to improve sleep quality",
    ],
  };

  return NextResponse.json({ plan });
}