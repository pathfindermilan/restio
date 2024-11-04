import { NextResponse } from 'next/server';

export async function GET() {
  // In a real app, you'd fetch this from a database
  const moodHistory = [
    { mood: 7, timestamp: new Date('2023-05-01T10:00:00Z') },
    { mood: 5, timestamp: new Date('2023-05-02T11:00:00Z') },
    // ... more mood entries
  ];

  return NextResponse.json({ moodHistory });
}

export async function POST(request) {
  const { mood, timestamp } = await request.json();

  // In a real app, you'd save this to a database
  console.log('Saving mood:', mood, 'at', timestamp);

  return NextResponse.json({ message: 'Mood saved successfully' });
}
