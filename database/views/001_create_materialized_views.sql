CREATE MATERIALIZED VIEW customer_fund_totals AS
SELECT 
    rc.id as customer_id,
    rc.first_name,
    rc.last_name,
    rc.email,
    f.id as fund_id,
    f.name as fund_name,
    SUM(i.amount) as total_investment
FROM retail_customers rc
JOIN investments i ON rc.id = i.customer_id
JOIN funds f ON f.id = i.fund_id
GROUP BY 
    rc.id, 
    rc.first_name, 
    rc.last_name, 
    rc.email,
    f.id,
    f.name;

