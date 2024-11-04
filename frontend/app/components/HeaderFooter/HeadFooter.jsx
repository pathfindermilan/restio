import React, { useState } from 'react';
import { Menu, X, Moon, User, BarChart } from 'lucide-react';
import Link from 'next/link';
import { useAuth } from "@/utils/auth";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const auth = useAuth() || {}; // Provide a default empty object if auth is undefined

  const handleLogout = async (e) => {
    e.preventDefault();
    if (auth.logout) {
      await auth.logout();
    }
  };

  const navLinks = [
    { name: 'Mood Tracking', href: '/', icon: BarChart },
  ];

  return (
    <header className="bg-transparent shadow-sm">
      <nav className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-2">
              <Moon className="h-8 w-8 text-purple-600" />
              <span className="text-xl font-bold text-white">MoodSync</span>
            </Link>
          </div>

          {/* Desktop Navigation */}
          <div className="hidden md:flex md:items-center md:space-x-8">
            {navLinks.map((link) => {
              const Icon = link.icon;
              return (
                <Link
                  key={link.name}
                  href={link.href}
                  className="flex items-center space-x-1 text-white hover:text-purple-600 transition-colors"
                >
                  <Icon className="h-4 w-4" />
                  <span>{link.name}</span>
                </Link>
              );
            })}
            
            {auth.loading === true ? (
              <div className="text-white">Loading...</div>
            ) : (
              <>
                {auth?.user ? (
                  <div className="flex items-center space-x-4">
                     <button
                      onClick={handleLogout}
                      className="text-white hover:text-purple-600 transition-colors"
                    >
                      Logout
                    </button>
                    <Link
                      href="/profile"
                      className="flex items-center space-x-1 bg-purple-600 text-white px-4 py-2 rounded-full hover:bg-purple-700 transition-colors"
                    >
                      <User className="h-4 w-4" />
                      <span>Profile</span>
                    </Link>
                   
                  </div>
                ) : (
                  <div className="flex items-center space-x-4">
                    <Link
                      href="/login"
                      className="text-white hover:text-purple-600 transition-colors"
                    >
                      Login
                    </Link>
                    <Link
                      href="/register"
                      className="bg-purple-600 text-white px-4 py-2 rounded-full hover:bg-purple-700 transition-colors"
                    >
                      Sign up
                    </Link>
                  </div>
                )}
              </>
            )}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="text-white hover:text-purple-600"
            >
              {isMenuOpen ? (
                <X className="h-6 w-6" />
              ) : (
                <Menu className="h-6 w-6" />
              )}
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {isMenuOpen && (
          <div className="md:hidden py-4">
            <div className="flex flex-col space-y-4">
              {navLinks.map((link) => {
                const Icon = link.icon;
                return (
                  <Link
                    key={link.name}
                    href={link.href}
                    className="flex items-center space-x-2 text-white hover:text-purple-600 transition-colors"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    <Icon className="h-4 w-4" />
                    <span>{link.name}</span>
                  </Link>
                );
              })}
              
              {auth.loading === true ? (
                <div className="text-white">Loading...</div>
              ) : (
                <>
                  {auth?.user ? (
                    <>
                      <Link
                        href="/profile"
                        className="flex items-center space-x-2 text-white hover:text-purple-600 transition-colors"
                        onClick={() => setIsMenuOpen(false)}
                      >
                        <User className="h-4 w-4" />
                        <span>Profile</span>
                      </Link>
                      <button
                        onClick={(e) => {
                          handleLogout(e);
                          setIsMenuOpen(false);
                        }}
                        className="flex items-center space-x-2 text-white hover:text-purple-600 transition-colors"
                      >
                        Logout
                      </button>
                    </>
                  ) : (
                    <>
                      <Link
                        href="/login"
                        className="text-white hover:text-purple-600 transition-colors"
                        onClick={() => setIsMenuOpen(false)}
                      >
                        Login
                      </Link>
                      <Link
                        href="/register"
                        className="text-white hover:text-purple-600 transition-colors"
                        onClick={() => setIsMenuOpen(false)}
                      >
                        Sign up
                      </Link>
                    </>
                  )}
                </>
              )}
            </div>
          </div>
        )}
      </nav>
    </header>
  );
};

const Footer = () => {
  const footerLinks = {
    Features: [
      { name: 'Mood Tracking', href: '/tracking' },
      { name: 'CBT Exercises', href: '/cbt' },
      { name: 'Mindfulness', href: '/mindfulness' },
      { name: 'Sleep Monitoring', href: '/sleep' },
    ],
    Support: [
      { name: 'Crisis Help', href: '/crisis' },
      { name: 'Find Therapist', href: '/therapist' },
      { name: 'Community', href: '/community' },
      { name: 'FAQ', href: '/faq' },
    ],
    Company: [
      { name: 'About Us', href: '/about' },
      { name: 'Privacy Policy', href: '/privacy' },
      { name: 'Terms of Service', href: '/terms' },
      { name: 'Contact', href: '/contact' },
    ],
  };

  return (
    <footer className="bg-transparent">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-7">
        {/* Bottom section */}
        <div className="border-gray-200">
          <p className="text-center text-lg text-white">
            © {new Date().getFullYear()} MoodSync. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
};

export { Header, Footer };