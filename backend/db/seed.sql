INSERT INTO users (email, password_hash, first_name, last_name, role, is_active)
VALUES (
	'admin@gmail.com',
	crypt('Admin@123', gen_salt('bf')),
	'admin',
	'admin',
	'admin',
	true
)
ON CONFLICT (email) DO NOTHING;
