-- Table data: arpanet_accounts
-- Add 3 accounts with password `password`

INSERT INTO `arpanet_accounts`
(`id`, `enabled`, `username`, `password`, `license`)
VALUES (1, 1, 'user-1', '$2y$10$QHt2PpQ3kYheZZTASOLY5uzpzoi30O9oYijIZabSE78a8yqfp7mjW', '3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO `arpanet_accounts`
(`id`, `enabled`, `username`, `password`, `license`)
VALUES (2, 1, 'user-2', '$2y$10$QHt2PpQ3kYheZZTASOLY5uzpzoi30O9oYijIZabSE78a8yqfp7mjW', 'fcee377a1fda007a8d2cc764a0a272e04d8c5d57');
INSERT INTO `arpanet_accounts`
(`id`, `enabled`, `username`, `password`, `license`)
VALUES (3, 1, 'user-3', '$2y$10$QHt2PpQ3kYheZZTASOLY5uzpzoi30O9oYijIZabSE78a8yqfp7mjW', 'db7e039146d5bf1b6781e7bc1bef31f0bb1298ea');
