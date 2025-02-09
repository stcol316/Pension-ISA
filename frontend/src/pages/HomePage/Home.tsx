import { FC } from 'react';
import Features from '../../components/Features/Features';
import Hero from '../../components/Hero/Hero';
import './Home.css';

const Home: FC = () => {
  return (
    <div className="home-page">
      <Hero />
      <Features />
    </div>
  );
};

export default Home;
