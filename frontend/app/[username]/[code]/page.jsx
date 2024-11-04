"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import { KeyRound } from "lucide-react";
import axios from "axios";

export default function ActivateAccount({ params }) {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    const activateAccount = async () => {
      const username = params?.username;
      const verification_code = params?.code;

      if (!username || !verification_code) {
        setError("Missing activation parameters");
        setLoading(false);
        return;
      }

      try {
        await axios.post("/api/activate", {
          username,
          verification_code,
        });
        
        setSuccess(true);
      } catch (err) {
        setError(
          err.response?.data?.message || 
          "Failed to activate account. Please try again."
        );
      } finally {
        setLoading(false);
      }
    };

    activateAccount();
  }, [params]);

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
              Account Activation
            </h2>
          </motion.div>

          <div className="space-y-6">
            {loading && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="text-center"
              >
                <div className="animate-pulse text-lg text-indigo-400">
                  Activating your account...
                </div>
              </motion.div>
            )}

            {error && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="text-center text-red-400 bg-red-900/20 p-4 rounded-lg border border-red-900"
              >
                {error}
              </motion.div>
            )}

            {success && (
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                className="text-center text-green-400 bg-green-900/20 p-4 rounded-lg border border-green-900"
              >
                Your account has been successfully activated! You can now log in.
              </motion.div>
            )}
          </div>

          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.8 }}
            className="mt-8 text-center text-sm text-gray-400"
          >
            {success ? (
              <Link
                href="/login"
                className="font-medium text-indigo-400 hover:text-indigo-300 transition-colors duration-200"
              >
                Click here to log in
              </Link>
            ) : (
              <>
                Not registered yet?{" "}
                <Link
                  href="/register"
                  className="font-medium text-indigo-400 hover:text-indigo-300 transition-colors duration-200"
                >
                  Create an account
                </Link>
              </>
            )}
          </motion.p>
        </div>
      </motion.div>
    </div>
  );
}