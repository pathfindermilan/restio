'use client'
import React, { useState, useRef, useEffect } from 'react';
import { Smile, ChevronDown, Send, Upload } from 'lucide-react';
import { Header } from './components/HeaderFooter/HeadFooter';
import { useAuth } from "@/utils/auth";
import { options, moodScale } from '@/lib/moddata';

export default function Home() {
  const auth = useAuth();
  const [imageFile, setImageFile] = useState(null);
  const [pdfFile, setPdfFile] = useState(null);
  const [age, setAge] = useState('');
  const [userText, setUserText] = useState('');
  const [showSuccess, setShowSuccess] = useState(false);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const [selectedOption, setSelectedOption] = useState(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [selectedMoodLevel, setSelectedMoodLevel] = useState(null);
  const cardRef = useRef(null);

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

  
  const handleSubmit = async (e) => {
    if (e) e.preventDefault();
    setIsProcessing(true);
   
    try {
        // Get auth token
        const token = localStorage.getItem("token");
      if (!token) {
        alert("Please login to submit");
        return;
      }
  
      const formData = new FormData();
  
      // Add user mood data
      formData.append("feeling_level", selectedMoodLevel);
      formData.append("query_text", userText);
      formData.append("age", age);
      formData.append("content_type", selectedOption);
  
      // Handle file uploads
      if (imageFile) {
        formData.append("image", imageFile);
      }
  
      if (pdfFile) {
        formData.append("document", pdfFile);
      }
      for (let [key, value] of formData.entries()) {
        console.log(key, value);
      }
      // Make API request with proper headers
      const response = await fetch('https://restio.website/api/sync', {
        method: 'POST',
        headers: {
          "Authorization": `Bearer ${token}`,
          "User-Agent": "insomnia/9.3.2",
        },
        body: formData
      });
  
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to submit mood data');
      }
  
      const result = await response.json();
      console.log('Submission successful:', result);
  
      // Show success message
      setShowSuccess(true);
      setTimeout(() => setShowSuccess(false), 3000);
  
      // Reset form
      resetForm();
  
    } catch (error) {
      console.error('Error submitting mood data:', error);
      // Show error message to user
      alert(error.message || 'Failed to submit mood data. Please try again.');
    } finally {
      setIsProcessing(false);
    }
  };
  
  // Helper function to reset form
  const resetForm = () => {
    setImageFile(null);
    setPdfFile(null);
    setAge('');
    setSelectedMoodLevel(null);
    setSelectedOption(null);
    setUserText('');
  };
  return (
    <div className="min-h-screen bg-cover bg-center bg-fixed bg-black/80 flex flex-col"
      style={{ backgroundImage: 'url("/4.webp")' }}>

      <header className="shadow-lg bg-black/80">
        <Header auth={auth} />
      </header>

      <main className="flex-1 flex items-center justify-center p-4">
        <div
          ref={cardRef}
          className="w-full max-w-3xl rounded-xl border border-white/10 shadow-2xl backdrop-blur-md bg-black/30 p-4 md:p-6 opacity-100"
        >
          {/* Header */}
          <div className="flex items-center gap-2 mb-4">
            <div className="hover:rotate-12 transition-transform">
              <Smile className="w-6 h-6 md:w-8 md:h-8 text-blue-400" />
            </div>
            <h2 className="text-xl md:text-2xl font-semibold text-white/90">Track Your Mood</h2>
          </div>

          {/* Content Type Dropdown */}
          <div className="space-y-2">
            <label className="block text-lg font-medium text-white/80">
              Select Content Type
            </label>
            <div className="relative">
              <button
                onClick={() => setIsDropdownOpen(!isDropdownOpen)}
                className="w-full p-4 backdrop-blur-sm border border-white/10 rounded-lg 
                    flex items-center justify-between text-white/80 hover:bg-white/5 
                    transition-colors"
              >
                {selectedOption ? (
                  <div className="flex items-center gap-2">
                    {options.find(opt => opt.id === selectedOption)?.icon}
                    <span className="text-lg">
                      {options.find(opt => opt.id === selectedOption)?.label}
                    </span>
                  </div>
                ) : (
                  <span className="text-white/50">Select content type...</span>
                )}
                <ChevronDown className={`w-5 h-5 transform transition-transform ${isDropdownOpen ? 'rotate-180' : ''}`} />
              </button>

              {isDropdownOpen && (
                <div className="absolute w-full mt-2 backdrop-blur-md bg-black/50 
                    border border-white/10 rounded-lg shadow-xl z-10">
                  {options.map((option) => (
                    <button
                      key={option.id}
                      onClick={() => {
                        setSelectedOption(option.id);
                        setIsDropdownOpen(false);
                      }}
                      className="w-full p-4 flex items-center gap-3 text-white/80 
                          hover:bg-white/10 transition-colors text-left"
                    >
                      {option.icon}
                      <span className="text-lg font-medium">{option.label}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* File Upload Section */}
            <div className="grid grid-cols-2 gap-2">
              <div className="border-2 border-dashed border-white/10 rounded-lg p-3 text-center">
                <input type="file" accept="image/*" onChange={handleImageUpload} className="hidden" id="imageUpload" />
                <label htmlFor="imageUpload" className="cursor-pointer">
                  <Upload className="mx-auto h-8 w-8 text-white/40" />
                  <p className="mt-1 text-sm md:text-base text-white/70">
                    {imageFile ? imageFile.name : 'Upload Image'}
                  </p>
                </label>
              </div>

              <div className="border-2 border-dashed border-white/10 rounded-lg p-3 text-center">
                <input type="file" accept=".pdf" onChange={handlePdfUpload} className="hidden" id="pdfUpload" />
                <label htmlFor="pdfUpload" className="cursor-pointer">
                  <Upload className="mx-auto h-8 w-8 text-white/40" />
                  <p className="mt-1 text-sm md:text-base text-white/70">
                    {pdfFile ? pdfFile.name : 'Upload PDF'}
                  </p>
                </label>
              </div>
            </div>

            {/* Age Input */}
            <div>
              <label className="block text-base md:text-lg font-medium text-white/80">Age</label>
              <input
                type="number"
                value={age}
                onChange={(e) => setAge(e.target.value)}
                className="w-full p-2 md:p-3 bg-black/30 border border-white/10 rounded-lg text-white text-base md:text-lg"
                placeholder="Enter your age"
              />
            </div>

            {/* Text Input */}
            <div>
              <label className="block text-base md:text-lg font-medium text-white/80">Your Thoughts</label>
              <textarea
                value={userText}
                onChange={(e) => setUserText(e.target.value)}
                className="w-full p-2 md:p-3 bg-black/30 border border-white/10 rounded-lg text-white text-base md:text-lg min-h-[100px]"
                placeholder="Share your thoughts and feelings..."
              />
            </div>

            {/* Mood Selection */}
            <div className="space-y-4">
            <label className="block text-lg font-medium text-white/80">
              How are you feeling today? (1-10)
            </label>
            <div className="grid grid-cols-5 gap-2">
              {moodScale.map(({ level, emoji, description }) => (
                <button
                  key={level}
                  onClick={() => setSelectedMoodLevel(level)}
                  className={`p-3 rounded-xl backdrop-blur-sm border border-white/10 
                    ${selectedMoodLevel === level ? 'bg-white/20' : 'bg-black/30'} 
                    hover:bg-white/10 transition-colors flex flex-col items-center`}
                >
                  <span className="text-2xl mb-1">{emoji}</span>
                  <span className="text-white/90 text-sm font-medium">{level}</span>
                  <span className="text-white/60 text-xs">{description}</span>
                </button>
              ))}
            </div>
            
            {selectedMoodLevel && (
              <div className="text-center text-white/80 mt-2">
                Selected mood level: {selectedMoodLevel}/10
              </div>
            )}
          </div>

            {/* Submit Button */}
            <button
            onClick={handleSubmit}
            className={`w-full backdrop-blur-sm border border-white/10 text-white 
              p-3 rounded-lg flex items-center justify-center gap-2 text-lg 
              font-medium transition-colors
              ${isProcessing ? 'bg-gray-500/30' :
              selectedMoodLevel ? 'bg-blue-500/30 hover:bg-blue-600/30' :
                'bg-gray-500/30'}`}
            disabled={isProcessing || !selectedMoodLevel}
          >
            <Send className="h-5 w-5" />
            {isProcessing ? 'Processing...' : 'Submit'}
          </button>
          </div>
        </div>

        {/* Success Message */}
        {showSuccess && (
          <div className="fixed bottom-4 right-4 backdrop-blur-md bg-green-500/30 
            border border-green-500/30 text-white p-3 rounded-lg shadow-lg">
            <p className="text-sm md:text-base font-medium">
              Mood data submitted successfully!
            </p>
          </div>
        )}
      </main>
    </div>
  );
}