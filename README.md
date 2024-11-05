# 🌟 Personal Mood Synchronizer App 🌟

## Watch Our Introduction Video
<p align="left">

  <a href="https://www.youtube.com/watch?v=7VcQHysm6qE&t=62s">
    <img src=".images/readme/youtube.png" alt="Restio" style="vertical-align:top; margin:6px 4px">
  </a>

</p>

## Project Overview
The **Restio** is designed to enhance mental health and promote self-improvement through positive daily habits. This innovative application empowers users by integrating mindfulness practices, personalized mood analysis, and community support, helping individuals cultivate a more positive mindset in their daily lives. 💪✨

## Features
- 🧘‍♀️ **Mindfulness and Meditation Exercises:** Engage users with guided practices to manage stress and improve emotional well-being.
- 🤖 **AI-Driven Mood Analysis:** Utilize artificial intelligence to provide personalized content based on users' emotional states, ensuring relevant support.
- 🤝 **Community Support:** Foster a collaborative environment where users can share experiences, uplift one another, and build a network of positivity.

## Why This App Matters
In a world where mental health challenges are increasingly prevalent, this app addresses the urgent need for accessible tools that help individuals navigate their emotions. By promoting positive habits and fostering community connections, we can create a supportive space for everyone to thrive.

## Getting Started
To get started with Restio:

1. ### Clone the repository:
   ```bash
   git clone https://github.com/pathfindermilan/restio.git
   ```
2. ### Run the fastapi AI app
   ```bash
   cd ai/
   cp .env.example .env
   vi .env
   ```
   * add your ARIA and ALLEGRO keys
   * save the env
   ```bash
   uv sync
   . .venv/bin/activate
   ```
   ```bash
   uvicorn api:app --host 127.0.0.1 --port 5000 --reload
   ```
    * more information about how to install uv python package manager 👉 [here](https://astral.sh/blog/uv)
    * more information about how to obtain the env keys, visit 👉 [Rhymes AI](https://rhymes.ai/)
3. ### Run the backend Go Gin gonic app
   #### Firstly run postgres instance
   ```bash
   docker run --name postgres-container -e POSTGRES_PASSWORD="password" -e POSTGRES_USER=postgres -e POSTGRES_DB=database_name -p 5432:5432 -d postgres
   ```
   ```bash
   cd backend/
   cp .env.example .env
   vi .env
   ```
   #### Add the following keys:
   * DB_NAME="database_name"
   * DB_PASSWORD="password"
   * DB_USER="postgres"
   * DB_PORT="5432"
   * DB_HOST="localhost"


   * JWT_SECRET=""
   * PORT=8002

   * SMTP_HOST=smtp.gmail.com
   * SMTP_PORT=587
   * SMTP_USE_TLS = True
   * SMTP_USERNAME="example@gmail.com" 
   * SMTP_PASSWORD="abcdefghijkl" # get the app password for example@gmail.com 👉 [here](https://myaccount.google.com/apppasswords)

   * FRONTEND_URL=http://localhost:3000

   * REDIS_HOST=localhost
   * REDIS_PORT=6379
   * REDIS_PASSWORD="password"
   * REDIS_DB=0
    
   * DESCRIBE_IMAGE_ENDPOINT=http://localhost:5000/describe-image
   * DESCRIBE_DOCUMENT_ENDPOINT=http://localhost:5000/describe-document
   * GENERATE_ANSWER=http://localhost:5000/generate-answer
   #### Install Go Lang
   - follow this tutorial 👉 [Install](https://go.dev/doc/install)
   #### Build the application
   ```bash
   go build -o server ./cmd/server/main.go
   ```
   #### Run the backend go lang application
   ```bash
   chmod +x server (optional)
   ./server
   ```
4. ### Run the frontend NextJS application
   ```bash
   cd frontend
   cp .env.example .env
   vi .env
   ```
   #### Change the key
   * NEXT_PUBLIC_SERVER=http://localhost:8002
   #### Install node and npm from 👉 [Install node and npm](https://nodejs.org/en/download/package-manager) or with the script below 👇
   ```bash
   sudo apt-get update
   
   if ! command_exists curl; then
      sudo apt-get install -y curl
   fi

   curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
   sudo apt-get install -y nodejs
   if [ ! -f "/usr/bin/node" ]; then
       sudo ln -s "$(which node)" /usr/bin/node
   fi

    if [ ! -f "/usr/bin/npm" ]; then
        sudo ln -s "$(which npm)" /usr/bin/npm
    fi

    # Verify installation
    echo "Verifying installation..."
    echo "Node.js version:"
    node --version
    echo "npm version:"
    npm --version

    # Check symbolic links
    echo "Checking symbolic links..."
    ls -l /usr/bin/node
    ls -l /usr/bin/npm

    # Set permissions
    echo "Setting permissions..."
    sudo chmod 755 /usr/bin/node
    sudo chmod 755 /usr/bin/npm

    # Configure npm global settings
    echo "Configuring npm global settings..."
    sudo mkdir -p /usr/local/lib/node_modules
    sudo chown -R $USER:$(id -gn $USER) /usr/local/lib/node_modules

    # Verify paths
    echo "Verifying paths..."
    which node
    which npm

    # Print installation completed message
    echo "Installation completed successfully!"
   ```
   #### Install node modules
   ```bash
   npm install
   ```
   #### Build the directory
   ```bash
   npm run build
   ```
   #### Start the app
   ```bash
   npm run dev
   ```
6. ### Your app is up and running on port 3000
   
## If you want to test our app deployed on GCloud 👇
   [RestiO](https://restio.xyz/)
## Overview
|    Part   |                Deployed                 |                Locally              |
|-----------|-----------------------------------------|-------------------------------------|
| Frontend  |     [Frontend](https://restio.xyz)      |         http://localhost:3000       |
|  Backend  |    [Backend](https://restio.website)    |         http://localhost:8002       |
|   GenAI   |   [Aria&Allegro](https://restio.site)   |         http://localhost:5000       |

## Technologies Used
This project utilizes the following technologies:
* Frontend: NextJS (v20.11.1)
* 🖥️ Backend: Go Gin Gonic, Fast Api
* 🗄️ Database: PostgresSQL
* 🤖 AI Integration: Area&Allegro
  
Each technology is chosen to ensure scalability, performance, and ease of use, providing a robust foundation for our app.
