import { NextResponse } from 'next/server';

export async function GET(request) {
  const { searchParams } = new URL(request.url);
  const mood = parseInt(searchParams.get('mood'));

  // In a real app, you'd have a more sophisticated way of selecting exercises
  const exercise = {
    description: "Let's practice reframing negative thoughts.",
    steps: [
      "Identify a negative thought you've had recently.",
      "Consider the evidence for and against this thought.",
      "Try to come up with a more balanced or realistic thought.",
      "Reflect on how this new perspective makes you feel.",
    ],
  };

  return NextResponse.json({ exercise });
}