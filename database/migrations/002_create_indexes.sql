-- Note: Indexing for faster reads
-- Due to the 'leftmost index rule' this will provide faster lookups on
-- getting all customer investments as well as all customer investments by fund
CREATE INDEX idx_investments_customer_fund ON investments(customer_id, fund_id);
--Note: Materialized view index
CREATE UNIQUE INDEX idx_customer_fund_totals ON customer_fund_totals(customer_id, fund_id);