"use client";

import { useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import { Eye, EyeOff, User, Lock, ArrowRight, HelpCircle } from "lucide-react";
import { useAuth } from "@/utils/auth";

export default function Login() {
  const auth = useAuth();
  const [formState, setFormState] = useState({
    username: "",
    password: "",
  });
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await auth.login(formState);
    } catch (error) {
      alert(error);
    }
  };

  const containerVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.6,
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { opacity: 0, x: -20 },
    visible: {
      opacity: 1,
      x: 0,
      transition: { duration: 0.5 }
    }
  };

  return (
    <div className="min-h-screen bg-cover bg-center bg-fixed flex items-center justify-center p-4"
    style={{ backgroundImage: 'url("/4.webp")'} }>
      <motion.div
        initial="hidden"
        animate="visible"
        variants={containerVariants}
        className="w-full max-w-md bg-white/10 backdrop-blur-lg rounded-2xl p-8 shadow-2xl"
      >
        <motion.div
          variants={itemVariants}
          className="text-center"
        >
          <h2 className="text-3xl font-bold text-white mb-2">Welcome Back</h2>
          <p className="text-indigo-200 mb-8">Sign in to your account</p>
        </motion.div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <motion.div variants={itemVariants}>
            <label className="block text-sm font-medium text-indigo-200 mb-2">
              Username
            </label>
            <div className="relative">
              <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-300 h-5 w-5" />
              <input
                type="text"
                required
                className="w-full pl-10 pr-4 py-3 bg-white/5 border border-indigo-300/20 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:border-transparent text-white placeholder-indigo-300"
                placeholder="Enter your username"
                value={formState.username}
                onChange={(e) => setFormState(state => ({ ...state, username: e.target.value }))}
              />
            </div>
          </motion.div>

          <motion.div variants={itemVariants}>
            <div className="flex justify-between mb-2">
              <label className="text-sm font-medium text-indigo-200">
                Password
              </label>
              <Link
                href="#"
                className="text-sm text-indigo-300 hover:text-indigo-200 flex items-center gap-1"
              >
                <HelpCircle className="h-4 w-4" />
                Forgot password?
              </Link>
            </div>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-300 h-5 w-5" />
              <input
                type={showPassword ? "text" : "password"}
                required
                className="w-full pl-10 pr-12 py-3 bg-white/5 border border-indigo-300/20 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:border-transparent text-white placeholder-indigo-300"
                placeholder="Enter your password"
                value={formState.password}
                onChange={(e) => setFormState(state => ({ ...state, password: e.target.value }))}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-indigo-300 hover:text-indigo-200"
              >
                {showPassword ? (
                  <EyeOff className="h-5 w-5" />
                ) : (
                  <Eye className="h-5 w-5" />
                )}
              </button>
            </div>
          </motion.div>

          <motion.div
            variants={itemVariants}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <button
              type="submit"
              className="w-full bg-gradient-to-r from-indigo-500 to-purple-500 text-white py-3 px-4 rounded-lg font-medium flex items-center justify-center gap-2 hover:from-indigo-600 hover:to-purple-600 transition-all duration-300 shadow-lg hover:shadow-xl"
            >
              Sign in
              <ArrowRight className="h-5 w-5" />
            </button>
          </motion.div>
        </form>

        <motion.p
          variants={itemVariants}
          className="mt-8 text-center text-indigo-200"
        >
          Not a member?{" "}
          <Link
            href="/register"
            className="font-semibold text-indigo-400 hover:text-indigo-300 transition-colors duration-200"
          >
            Register here
          </Link>
        </motion.p>
      </motion.div>
    </div>
  );
}