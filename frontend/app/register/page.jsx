"use client";

import { useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import { Eye, EyeOff, User, Mail, Lock, ArrowRight, UserCircle } from "lucide-react";
import { useAuth } from "@/utils/auth";

export default function Register() {
  const auth = useAuth();
  const [formState, setFormState] = useState({
    email: "",
    username: "",
    password: "",
    re_password: "",
    first_name: "",
    last_name: "",
  });

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await auth.register(formState);
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
    <div 
      className="min-h-screen bg-cover bg-center bg-fixed flex items-center justify-center p-4"
      style={{ backgroundImage: 'url("/4.webp")' }}
    >
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
          <h2 className="text-3xl font-bold text-white mb-2">Create Account</h2>
          <p className="text-indigo-200 mb-8">Join us today</p>
        </motion.div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <motion.div variants={itemVariants}>
              <label className="block text-sm font-medium text-indigo-200 mb-2">
                First Name
              </label>
              <div className="relative">
                <UserCircle className="absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-300 h-5 w-5" />
                <input
                  type="text"
                  required
                  className="w-full pl-10 pr-4 py-3 bg-white/5 border border-indigo-300/20 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:border-transparent text-white placeholder-indigo-300"
                  placeholder="First name"
                  value={formState.first_name}
                  onChange={(e) => setFormState(state => ({ ...state, first_name: e.target.value }))}
                />
              </div>
            </motion.div>

            <motion.div variants={itemVariants}>
              <label className="block text-sm font-medium text-indigo-200 mb-2">
                Last Name
              </label>
              <div className="relative">
                <UserCircle className="absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-300 h-5 w-5" />
                <input
                  type="text"
                  required
                  className="w-full pl-10 pr-4 py-3 bg-white/5 border border-indigo-300/20 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:border-transparent text-white placeholder-indigo-300"
                  placeholder="Last name"
                  value={formState.last_name}
                  onChange={(e) => setFormState(state => ({ ...state, last_name: e.target.value }))}
                />
              </div>
            </motion.div>
          </div>

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
                placeholder="Choose a username"
                value={formState.username}
                onChange={(e) => setFormState(state => ({ ...state, username: e.target.value }))}
              />
            </div>
          </motion.div>

          <motion.div variants={itemVariants}>
            <label className="block text-sm font-medium text-indigo-200 mb-2">
              Email
            </label>
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-300 h-5 w-5" />
              <input
                type="email"
                required
                className="w-full pl-10 pr-4 py-3 bg-white/5 border border-indigo-300/20 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:border-transparent text-white placeholder-indigo-300"
                placeholder="Your email address"
                value={formState.email}
                onChange={(e) => setFormState(state => ({ ...state, email: e.target.value }))}
              />
            </div>
          </motion.div>

          <motion.div variants={itemVariants}>
            <label className="block text-sm font-medium text-indigo-200 mb-2">
              Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-300 h-5 w-5" />
              <input
                type={showPassword ? "text" : "password"}
                required
                className="w-full pl-10 pr-12 py-3 bg-white/5 border border-indigo-300/20 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:border-transparent text-white placeholder-indigo-300"
                placeholder="Create a password"
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

          <motion.div variants={itemVariants}>
            <label className="block text-sm font-medium text-indigo-200 mb-2">
              Confirm Password
            </label>
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-300 h-5 w-5" />
              <input
                type={showConfirmPassword ? "text" : "password"}
                required
                className="w-full pl-10 pr-12 py-3 bg-white/5 border border-indigo-300/20 rounded-lg focus:ring-2 focus:ring-indigo-400 focus:border-transparent text-white placeholder-indigo-300"
                placeholder="Confirm your password"
                value={formState.re_password}
                onChange={(e) => setFormState(state => ({ ...state, re_password: e.target.value }))}
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-indigo-300 hover:text-indigo-200"
              >
                {showConfirmPassword ? (
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
            className="mt-6"
          >
            <button
              type="submit"
              className="w-full bg-gradient-to-r from-indigo-500 to-purple-500 text-white py-3 px-4 rounded-lg font-medium flex items-center justify-center gap-2 hover:from-indigo-600 hover:to-purple-600 transition-all duration-300 shadow-lg hover:shadow-xl"
            >
              Create Account
              <ArrowRight className="h-5 w-5" />
            </button>
          </motion.div>
        </form>

        <motion.p
          variants={itemVariants}
          className="mt-8 text-center text-indigo-200"
        >
          Already have an account?{" "}
          <Link
            href="/login"
            className="font-semibold text-indigo-400 hover:text-indigo-300 transition-colors duration-200"
          >
            Sign in
          </Link>
        </motion.p>
      </motion.div>
    </div>
  );
}