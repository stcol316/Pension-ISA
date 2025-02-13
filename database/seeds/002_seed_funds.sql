INSERT INTO risk_levels (id, name, description) VALUES
    (1, 'Low', 'Conservative investments focusing on capital preservation. Typically includes high-grade bonds and stable assets with lower potential returns but minimal risk of loss.'),
    (2, 'Medium', 'Balanced approach combining stability and growth. Mix of bonds and equities aiming for moderate long-term returns while managing volatility.'),
    (3, 'High', 'Growth-focused investments accepting larger short-term fluctuations for potentially higher long-term returns. Primarily equities and higher-risk assets.');

INSERT INTO funds (name, description, risk_level_id) VALUES 
    ('Ethical Bond Fund', 'Fixed income investments meeting strict ethical criteria', 1),
    ('Balanced Growth Fund', 'Balanced portfolio of 60% stocks and 40% bonds', 2),
    ('Emerging Markets Fund', 'Focus on high-growth potential markets in developing economies', 3);