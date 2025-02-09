import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Navigation from './components/Navigation/Navigation';
import Home from './pages/HomePage/Home';
import Features from './components/Features/Features';
import './App.css';

const App: React.FC = () => {
  return (
    <BrowserRouter>
      <div className="app">
        <Navigation />
        <main>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/features" element={<Features />} />
            <Route path="*" element={
              <div className="not-found">
                <h1>404 - Page Not Found</h1>
              </div>
            } />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
};

export default App;
