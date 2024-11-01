import { useEffect, useState } from 'react';

export default function PersonalizedPlan({ moodHistory }) {
  const [plan, setPlan] = useState(null);

  useEffect(() => {
    const generatePlan = async () => {
      const response = await fetch('/api/generate-plan', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ moodHistory }),
      });
      const data = await response.json();
      setPlan(data.plan);
    };

    if (moodHistory.length > 0) {
      generatePlan();
    }
  }, [moodHistory]);

  if (!plan) return null;

  return (
    <div className="my-8 p-6 bg-purple-100 rounded-lg">
      <h2 className="text-2xl font-semibold mb-4">Your Personalized Plan</h2>
      <ul className="list-disc list-inside">
        {plan.recommendations.map((recommendation, index) => (
          <li key={index} className="mb-2">{recommendation}</li>
        ))}
      </ul>
    </div>
  );
}