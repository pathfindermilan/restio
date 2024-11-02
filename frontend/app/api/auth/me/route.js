const { NextResponse } = require("next/server");
const axios = require("axios").default;

export async function POST(request) {
  const body = await request.json();
  const { token } = body;

  if (!token) {
    return NextResponse.json({ error: "Missing Auth token" }, { status: 400 });
  }

  try {
    const response = await axios.get(
      `${process.env.NEXT_PUBLIC_SERVER}/api/profile`,
      {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
          "User-Agent": "insomnia/9.3.2",
        },
      }
    );

    console.log("ME response:");
    console.log(response.data);

    return NextResponse.json(response.data, { status: 200 });
  } catch (error) {
    return NextResponse.json({ error }, { status: 400 });
  }
}
