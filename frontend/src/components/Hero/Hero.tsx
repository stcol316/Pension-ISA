import { FC } from 'react';
import './Hero.css';

const Hero: FC = () => {
    return (
        <div className="hero">
            <h1>Making workplace pensions and savings easy</h1>
            <p>Secure your future with smart investments</p>
            <button className="cta-button">Get Started</button>
        </div>
    );
};

export default Hero;
