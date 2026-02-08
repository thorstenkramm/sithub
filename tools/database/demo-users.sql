-- Demo users for SitHub local development and E2E testing.
-- All users have the password: SitHubDemo2026!!
-- User source: internal (local authentication)
-- bcrypt cost: 12

INSERT OR IGNORE INTO users (id, email, display_name, password_hash, user_source, entra_id, is_admin, last_login, created_at, updated_at) VALUES
('7c937bdb-a9ec-4a07-b6fe-346be95b8e95', 'anna@sithub.local', 'Anna Admin', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 1, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('2a93a785-c52e-4cf5-b109-eaab0cb7d2d4', 'sam@sithub.local', 'Sam Superadmin', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 1, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('830232de-a6a0-41af-b33a-c459322dab43', 'alex@sithub.local', 'Alex Employee', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('480f24d3-65ad-4e98-85b7-3f6eec52a28c', 'dana@sithub.local', 'Dana Developer', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('dad3e32e-5a5a-4f61-92df-4372ea3e53ee', 'emma@sithub.local', 'Emma Engineer', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('0f420659-99d0-4aac-bd51-008ccd4b1e17', 'felix@sithub.local', 'Felix Finance', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('4f001b3e-af37-430f-8015-f82f9caa7f5c', 'gina@sithub.local', 'Gina Graphics', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('95c97dc2-775b-4092-916e-075e3a8af438', 'hans@sithub.local', 'Hans HR', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('282d37c5-2e77-4413-a154-ca32f564fd14', 'iris@sithub.local', 'Iris Intern', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('b9d15b08-c583-41d5-80e6-7fa23660bfde', 'jan@sithub.local', 'Jan Junior', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('709ee9c1-0991-45cb-8f4d-6c93737ccca4', 'kate@sithub.local', 'Kate Keeper', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('387e2eab-e19f-40e6-930f-43f60993d590', 'leo@sithub.local', 'Leo Lead', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('06dda637-69e5-432a-9cb4-15e076ee2b09', 'mia@sithub.local', 'Mia Manager', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('b707b839-8de3-48b8-83fb-3f1ebeb3117b', 'nico@sithub.local', 'Nico Network', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z'),
('ae04a026-9107-4c3b-ae95-52d4bae7f273', 'olivia@sithub.local', 'Olivia Ops', '$2a$12$H2due882dbw6aQDvnbwUeuoW731Lb6Ze0qOCmGU3RkBP5FWmMLfxy', 'internal', '', 0, '', '2026-01-01T00:00:00Z', '2026-01-01T00:00:00Z');
