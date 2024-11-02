'use client'
import React, { useState, useRef } from 'react';
import { Upload, Mic, Send, Video, Laugh, Music, Dumbbell, ChevronDown } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const ModSyncForm = () => {
  const [imageFile, setImageFile] = useState(null);
  const [pdfFile, setPdfFile] = useState(null);
  const [age, setAge] = useState('');
  const [isRecording, setIsRecording] = useState(false);
  const [selectedEmoji, setSelectedEmoji] = useState('');
  const [showSuccess, setShowSuccess] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [selectedOption, setSelectedOption] = useState(null);
  const audioRef = useRef(null);

  const options = [
    { id: 'speech', label: 'Speech', icon: <Music className="w-5 h-5" /> },
    { id: 'exercise', label: 'Exercise', icon: <Dumbbell className="w-5 h-5" /> },
    { id: 'joke', label: 'Joke', icon: <Laugh className="w-5 h-5" /> },
    { id: 'video', label: 'Video', icon: <Video className="w-5 h-5" /> }
  ];
  
  const emojis = [
    { emoji: '😊', feeling: 'Happy' },
    { emoji: '😢', feeling: 'Sad' },
    { emoji: '😡', feeling: 'Angry' },
    { emoji: '😴', feeling: 'Tired' },
    { emoji: '🤔', feeling: 'Confused' },
    { emoji: '😎', feeling: 'Confident' }
  ];

  const handleImageUpload = (e) => {
    const file = e.target.files[0];
    if (file && file.type.startsWith('image/')) {
      setImageFile(file);
    }
  };

  const handlePdfUpload = (e) => {
    const file = e.target.files[0];
    if (file && file.type === 'application/pdf') {
      setPdfFile(file);
    }
  };

  const toggleRecording = () => {
    setIsRecording(!isRecording);
  };

  const handleSubmit = async () => {
    setShowSuccess(true);
    setTimeout(() => setShowSuccess(false), 3000);
    
    const formData = {
      image: imageFile,
      pdf: pdfFile,
      age: age,
      feeling: selectedEmoji,
      hasVoiceData: isRecording,
      selectedOption: selectedOption
    };
    
    console.log('Submitting data:', formData);
  };

  return (
    <div className="min-h-screen  p-8 font-sans bg-cover bg-center bg-fixed" style={{ backgroundImage: 'url("/4.webp")'} }>
      <div className="max-w-2xl mx-auto bg-white rounded-xl shadow-lg overflow-hidden">
        {/* Card Header */}
        <div className="p-6 border-b border-gray-200 bg-gradient-to-r from-blue-500 to-purple-500">
          <h2 className="text-3xl font-bold text-center text-white font-sans">
            ModSync Data Collection
          </h2>
        </div>
        
        {/* Card Content */}
        <div className="p-6 space-y-8">
          {/* Content Type Dropdown */}
          <div className="space-y-2">
            <label className="block text-lg font-semibold text-gray-700">
              Select Content Type
            </label>
            <div className="relative">
              <motion.button
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                className="w-full p-4 bg-white border-2 rounded-lg flex items-center justify-between text-gray-700 hover:border-blue-500 transition-colors"
                whileHover={{ scale: 1.01 }}
                whileTap={{ scale: 0.99 }}
              >
                {selectedOption ? (
                  <div className="flex items-center gap-2">
                    {options.find(opt => opt.id === selectedOption)?.icon}
                    <span className="text-lg">
                      {options.find(opt => opt.id === selectedOption)?.label}
                    </span>
                  </div>
                ) : (
                  <span className="text-gray-500">Select content type...</span>
                )}
                <motion.div
                  animate={{ rotate: isDropdownOpen ? 180 : 0 }}
                  transition={{ duration: 0.2 }}
                >
                  <ChevronDown className="w-5 h-5" />
                </motion.div>
              </motion.button>

              <AnimatePresence>
                {isDropdownOpen && (
                  <motion.div
                    initial={{ opacity: 0, y: -10 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -10 }}
                    transition={{ duration: 0.2 }}
                    className="absolute w-full mt-2 bg-white border rounded-lg shadow-xl z-10"
                  >
                    {options.map((option) => (
                      <motion.button
                        key={option.id}
                        onClick={() => {
                          setSelectedOption(option.id);
                          setIsDropdownOpen(false);
                        }}
                        className="w-full p-4 flex items-center gap-3 hover:bg-blue-50 transition-colors text-left"
                        whileHover={{ x: 4 }}
                      >
                        {option.icon}
                        <span className="text-lg font-medium">{option.label}</span>
                      </motion.button>
                    ))}
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          </div>

          {/* File Upload Section */}
          <div className="space-y-4">
            <motion.div
              className="border-2 border-dashed rounded-lg p-6 text-center hover:border-blue-500 transition-colors"
              whileHover={{ scale: 1.01 }}
              whileTap={{ scale: 0.99 }}
            >
              <input
                type="file"
                accept="image/*"
                onChange={handleImageUpload}
                className="hidden"
                id="imageUpload"
              />
              <label htmlFor="imageUpload" className="cursor-pointer">
                <Upload className="mx-auto h-12 w-12 text-gray-400" />
                <p className="mt-2 text-lg text-gray-600 font-medium">
                  {imageFile ? imageFile.name : 'Upload Image'}
                </p>
              </label>
            </motion.div>

            <motion.div
              className="border-2 border-dashed rounded-lg p-6 text-center hover:border-blue-500 transition-colors"
              whileHover={{ scale: 1.01 }}
              whileTap={{ scale: 0.99 }}
            >
              <input
                type="file"
                accept=".pdf"
                onChange={handlePdfUpload}
                className="hidden"
                id="pdfUpload"
              />
              <label htmlFor="pdfUpload" className="cursor-pointer">
                <Upload className="mx-auto h-12 w-12 text-gray-400" />
                <p className="mt-2 text-lg text-gray-600 font-medium">
                  {pdfFile ? pdfFile.name : 'Upload PDF'}
                </p>
              </label>
            </motion.div>
          </div>

          {/* Age Input */}
          <div className="space-y-2">
            <label className="block text-lg font-semibold text-gray-700">Age</label>
            <motion.input
              type="number"
              value={age}
              onChange={(e) => setAge(e.target.value)}
              className="w-full p-4 border-2 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-lg"
              placeholder="Enter your age"
              whileFocus={{ scale: 1.01 }}
            />
          </div>

          {/* Voice Recording */}
          <div className="text-center">
            <motion.button
              onClick={toggleRecording}
              className={`p-6 rounded-full ${
                isRecording ? 'bg-red-500' : 'bg-blue-500'
              } text-white shadow-lg`}
              whileHover={{ scale: 1.1 }}
              whileTap={{ scale: 0.9 }}
              animate={{
                scale: isRecording ? [1, 1.1, 1] : 1,
                transition: {
                  repeat: isRecording ? Infinity : 0,
                  duration: 1
                }
              }}
            >
              <Mic className="h-8 w-8" />
            </motion.button>
            <p className="mt-3 text-lg font-medium text-gray-600">
              {isRecording ? 'Recording...' : 'Click to Record Voice'}
            </p>
          </div>

          {/* Emoji Selection */}
          <div className="space-y-3">
            <label className="block text-lg font-semibold text-gray-700">
              How are you feeling?
            </label>
            <div className="flex flex-wrap gap-4 justify-center">
              {emojis.map(({ emoji, feeling }) => (
                <motion.button
                  key={feeling}
                  onClick={() => setSelectedEmoji(emoji)}
                  className={`p-3 rounded-xl ${
                    selectedEmoji === emoji ? 'bg-blue-100' : 'bg-gray-50'
                  }`}
                  whileHover={{ scale: 1.1 }}
                  whileTap={{ scale: 0.95 }}
                >
                  <span role="img" aria-label={feeling} className="text-4xl">
                    {emoji}
                  </span>
                  <p className="text-sm font-medium mt-2">{feeling}</p>
                </motion.button>
              ))}
            </div>
          </div>
        </div>

        {/* Card Footer */}
        <div className="p-6 border-t border-gray-200">
          <motion.button
            onClick={handleSubmit}
            className="w-full bg-gradient-to-r from-blue-500 to-purple-500 text-white p-4 rounded-lg flex items-center justify-center gap-3 text-lg font-semibold shadow-lg"
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <Send className="h-6 w-6" />
            Submit Data
          </motion.button>
        </div>
      </div>

      {/* Success Alert */}
      <AnimatePresence>
        {showSuccess && (
          <motion.div
            initial={{ opacity: 0, y: 50 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 50 }}
            className="fixed bottom-4 right-4 bg-green-500 text-white p-4 rounded-lg shadow-lg"
          >
            <p className="text-lg font-medium">
              Data submitted successfully! AI recommendations coming soon.
            </p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default ModSyncForm;