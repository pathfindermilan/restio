'use client'
import { useState, useEffect } from 'react';
import { Smile, Brain, Heart, Users, Calendar } from 'lucide-react';
import { Header, Footer } from './components/HeaderFooter/HeadFooter';
import { motion } from 'framer-motion';
import { useAuth } from "@/utils/auth";
// Button Component with glass effect
const Button = ({ children, className = '', ...props }) => {
  return (
    <motion.button
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
      className={`inline-flex items-center justify-center rounded-md text-sm font-medium 
      transition-all duration-300 focus-visible:outline-none focus-visible:ring-2 
      focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 
      h-10 px-4 py-2 text-white backdrop-blur-sm bg-opacity-20 shadow-xl
      border border-white/10 ${className}`}
      {...props}
    >
      {children}
    </motion.button>
  );
};

// Card Component with glass effect
const Card = ({ children, className = '', ...props }) => {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      whileHover={{ scale: 1.02, translateY: -5 }}
      className={`rounded-xl border border-white/10 shadow-2xl 
      backdrop-blur-md bg-black/30 hover:bg-black/40 
      transition-all duration-300 ${className}`}
      {...props}
    >
      {children}
    </motion.div>
  );
};

export default function Home() {
  const auth = useAuth();

  console.log(auth);

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  return (
    
    <div className="min-h-screen bg-cover bg-center bg-fixed bg-black/80 " 
         style={{ backgroundImage: 'url("/4.webp")' }}>
      
      <header className="shadow-lg bg-black/80">
        <Header auth={auth}/>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
          className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6"
        >
          {/* Mood Tracking Card */}
          <Card className="p-4 sm:p-6">
            <div className="flex items-center gap-3 sm:gap-4 mb-3 sm:mb-4">
              <motion.div
                whileHover={{ rotate: 20 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Smile className="w-6 h-6 sm:w-8 sm:h-8 text-blue-400" />
              </motion.div>
              <h2 className="text-lg sm:text-xl font-semibold text-white/90">Mood Tracking</h2>
            </div>
            <p className="text-white/70 text-sm sm:text-base mb-4">
              Track your daily moods and emotions with our intuitive interface. 
              Get personalized insights based on your patterns.
            </p>
            <Button className="w-full bg-blue-500/50 hover:bg-blue-600/50">Track Mood</Button>
          </Card>

          {/* CBT Exercises Card */}
          <Card className="p-4 sm:p-6">
            <div className="flex items-center gap-3 sm:gap-4 mb-3 sm:mb-4">
              <motion.div
                whileHover={{ rotate: 180 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Brain className="w-6 h-6 sm:w-8 sm:h-8 text-purple-400" />
              </motion.div>
              <h2 className="text-lg sm:text-xl font-semibold text-white/90">CBT Exercises</h2>
            </div>
            <p className="text-white/70 text-sm sm:text-base mb-4">
              Access proven cognitive behavioral therapy techniques to help manage 
              stress, anxiety, and negative thoughts.
            </p>
            <Button className="w-full bg-purple-500/50 hover:bg-purple-600/50">Start Exercise</Button>
          </Card>

          {/* Mindfulness Card */}
          <Card className="p-4 sm:p-6">
            <div className="flex items-center gap-3 sm:gap-4 mb-3 sm:mb-4">
              <motion.div
                whileHover={{ scale: 1.2 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Heart className="w-6 h-6 sm:w-8 sm:h-8 text-red-400" />
              </motion.div>
              <h2 className="text-lg sm:text-xl font-semibold text-white/90">Mindfulness</h2>
            </div>
            <p className="text-white/70 text-sm sm:text-base mb-4">
              Practice guided meditation and relaxation techniques to reduce 
              stress and improve mental clarity.
            </p>
            <Button className="w-full bg-red-500/50 hover:bg-red-600/50">Begin Session</Button>
          </Card>

          {/* Community Support Card */}
          <Card className="p-4 sm:p-6">
            <div className="flex items-center gap-3 sm:gap-4 mb-3 sm:mb-4">
              <motion.div
                whileHover={{ y: -5 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Users className="w-6 h-6 sm:w-8 sm:h-8 text-green-400" />
              </motion.div>
              <h2 className="text-lg sm:text-xl font-semibold text-white/90">Community Support</h2>
            </div>
            <p className="text-white/70 text-sm sm:text-base mb-4">
              Connect with others on similar journeys and share experiences 
              in a safe, moderated environment.
            </p>
            <Button className="w-full bg-green-500/50 hover:bg-green-600/50">Join Community</Button>
          </Card>

          {/* Therapist Connection Card */}
          <Card className="p-4 sm:p-6">
            <div className="flex items-center gap-3 sm:gap-4 mb-3 sm:mb-4">
              <motion.div
                whileHover={{ rotate: 360 }}
                transition={{ type: "spring", stiffness: 300 }}
              >
                <Calendar className="w-6 h-6 sm:w-8 sm:h-8 text-amber-400" />
              </motion.div>
              <h2 className="text-lg sm:text-xl font-semibold text-white/90">Connect with Therapists</h2>
            </div>
            <p className="text-white/70 text-sm sm:text-base mb-4">
              Schedule sessions with licensed therapists when you need 
              professional support.
            </p>
            <Button className="w-full bg-amber-500/50 hover:bg-amber-600/50">Book Session</Button>
          </Card>
        </motion.div>
      </main>

      <footer className="bg-black/80 lg:fixed w-full bottom-0">
        <Footer />
      </footer>
    </div>
  );
}