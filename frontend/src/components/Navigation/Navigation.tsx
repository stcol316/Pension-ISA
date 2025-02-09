import { FC } from 'react';

const Navigation: FC = () => {
    return(
        <nav className="nav">
            <div className="logo">
                <h1>Pension & ISA</h1>
            </div>
            <ul className="nav-links">
                <li><a href="#features">Features</a></li>
                <li><a href="#about">About</a></li>
                <li><a href="#contact">Contact</a></li>
            </ul>
        </nav>
    );
};

export default Navigation