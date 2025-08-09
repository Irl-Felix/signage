-- Seed Businesses
INSERT INTO Business (name, subscription_plan, status) VALUES
('TechNova Solutions', 'basic', 'active'),
('GlobalMart Retail', 'pro', 'active'),
('PixelWave Media', 'enterprise', 'pending'),
('FreshHarvest Foods', 'basic', 'suspended'),
('BlueSky Travel', 'pro', 'active'),
('UrbanEdge Apparel', 'enterprise', 'active'),
('EduCore Learning', 'basic', 'active'),
('HealthPlus Clinic', 'pro', 'pending');

-- Seed Locations
INSERT INTO Location (business_id, name, address, timezone) VALUES
((SELECT id FROM Business WHERE name='TechNova Solutions'), 'Headquarters', '123 Main St, Taipei', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='TechNova Solutions'), 'Branch North', '45 North Rd, Hsinchu', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='GlobalMart Retail'), 'Downtown Store', '456 Market Rd, Kaohsiung', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='GlobalMart Retail'), 'Mall Outlet', '88 Shopping Blvd, Tainan', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='PixelWave Media'), 'Creative Hub', '789 Studio Ln, Taichung', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='FreshHarvest Foods'), 'Main Market', '12 Farm Ave, Tainan', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='BlueSky Travel'), 'Airport Office', '1 Skyway Blvd, Taoyuan', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='UrbanEdge Apparel'), 'Flagship Store', '999 Fashion Ave, Taipei', 'Asia/Taipei'),
((SELECT id FROM Business WHERE name='EduCore Learning'), 'Main Campus', '5 Scholar St, Taipei', 'Asia/Taipei');
