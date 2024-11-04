// app/verify_user/VerifyUserForm.js
"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { motion } from "framer-motion";
import { KeyRound, ArrowRight } from "lucide-react";

export default function VerifyUserForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [username, setUsername] = useState("");
  const [verificationCode, setVerificationCode] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const usernameParam = searchParams?.get("username");
    if (usernameParam) {
      setUsername(usernameParam);
    } else {
      setError("Missing username parameter");
    }
  }, [searchParams]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      if (!username || !verificationCode) {
        throw new Error("Username and verification code are required");
      }

      // API call would go here
      await router.push(`/${username}/${verificationCode}`);
    } catch (err) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  if (error && !username) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center text-red-400 bg-red-900/50 p-4 rounded-md">
          {error}. Please check your activation email.
        </div>
      </div>
    );
  }

  return (
    <div
      className="min-h-screen bg-cover bg-center bg-fixed relative text-white flex items-center"
      style={{ backgroundImage: 'url("/1.jpg")' }}
    >
      <div className="absolute inset-0 bg-black bg-opacity-70"></div>
      
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="relative z-10 container mx-auto px-4 py-20"
      >
        <div className="max-w-md mx-auto bg-gray-900/60 backdrop-blur-lg rounded-xl p-8 shadow-2xl">
          <motion.div
            initial={{ scale: 0.8 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.2 }}
            className="flex flex-col items-center mb-8"
          >
            <div className="bg-indigo-500/20 p-3 rounded-full mb-4">
              <KeyRound className="w-8 h-8 text-indigo-400" />
            </div>
            <h2 className="text-2xl font-bold text-center tracking-tight bg-gradient-to-r from-indigo-400 to-purple-400 bg-clip-text text-transparent">
              Activate Your Account
            </h2>
            <p className="mt-2 text-gray-400">
              Username: <span className="text-indigo-400">{username}</span>
            </p>
          </motion.div>

          {error && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              className="mb-4 p-3 rounded-lg bg-red-900/50 text-red-400 text-sm"
            >
              {error}
            </motion.div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <motion.div
              initial={{ x: -20, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ delay: 0.4 }}
            >
              <label className="block text-sm font-medium mb-2 text-gray-300">
                Activation Code
              </label>
              <input
                type="text"
                value={verificationCode}
                onChange={(e) => setVerificationCode(e.target.value)}
                className="w-full px-4 py-3 rounded-lg bg-gray-800/50 border border-gray-700 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/30 transition-all duration-200 text-white placeholder-gray-400"
                placeholder="Enter your activation code"
                required
                autoFocus
                disabled={isLoading}
              />
            </motion.div>

            <motion.button
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.6 }}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              type="submit"
              disabled={isLoading}
              className="w-full flex items-center justify-center gap-2 py-3 px-4 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 rounded-lg font-semibold transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                "Processing..."
              ) : (
                <>
                  Activate Account
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </motion.button>
          </form>

          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.8 }}
            className="mt-8 text-center text-sm text-gray-400"
          >
            Not registered yet?{" "}
            <Link
              href="/register"
              className="font-medium text-indigo-400 hover:text-indigo-300 transition-colors duration-200"
            >
              Create an account
            </Link>
          </motion.p>
        </div>
      </motion.div>
    </div>
  );
}