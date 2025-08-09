INSERT INTO UserProfile (supabase_uid, first_name, last_name, email, phone, preferred_language, status) VALUES
(gen_random_uuid(), 'Alice', 'Chen', 'alice.chen@example.com', '+886912345678', 'en', 'active'),
(gen_random_uuid(), 'Bob', 'Wang', 'bob.wang@example.com', '+886987654321', 'zh', 'active'),
(gen_random_uuid(), 'Charlie', 'Lin', 'charlie.lin@example.com', '+886998877665', 'en', 'pending'),
(gen_random_uuid(), 'David', 'Liu', 'david.liu@example.com', '+886912398765', 'zh', 'suspended'),
(gen_random_uuid(), 'Emily', 'Huang', 'emily.huang@example.com', '+886976543210', 'en', 'active'),
(gen_random_uuid(), 'Frank', 'Ko', 'frank.ko@example.com', '+886900112233', 'en', 'active'),
(gen_random_uuid(), 'Grace', 'Chang', 'grace.chang@example.com', '+886987650123', 'zh', 'pending'),
(gen_random_uuid(), 'Henry', 'Wu', 'henry.wu@example.com', '+886911223344', 'en', 'active'),
(gen_random_uuid(), 'Ivy', 'Tsai', 'ivy.tsai@example.com', '+886923456789', 'zh', 'active'),
(gen_random_uuid(), 'Jack', 'Yeh', 'jack.yeh@example.com', '+886955667788', 'en', 'active');
