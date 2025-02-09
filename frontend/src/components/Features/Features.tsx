import { FC } from 'react';
import './Features.css';

const Features: FC = () => {
    return (
        <section className="features">
            <div className="features-grid">
                <div className="feature-card">
                    <h3>Pension Planning</h3>
                    <p>Expert guidance and tools to help you plan for a comfortable retirement with confidence.</p>
                </div>
                <div className="feature-card">
                    <h3>ISA Management</h3>
                    <p>Maximize your tax-free savings with our easy-to-use ISA management platform.</p>
                </div>
                <div className="feature-card">
                    <h3>Smart Investing</h3>
                    <p>Access sophisticated investment strategies tailored to your goals and risk tolerance.</p>
                </div>
            </div>
        </section>
    );
};

export default Features;
