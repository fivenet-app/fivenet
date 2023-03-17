-- Table data: arpanet_accounts
-- Add 3 accounts with password `test-password-1`, `test-password-2` and `test-password-3`

INSERT INTO `arpanet_accounts`
(`enabled`, `username`, `password`, `license`)
VALUES (1, 'user-1', '$2y$10$z0m12MxuuGmCaz2f7NuFReUG6jD4U3O.6tmqlQCdCIjgbMtw0.qOm', '3c7681d6f7ad895eb7b1cc05cf895c7f1d1622c4');
INSERT INTO `arpanet_accounts`
(`enabled`, `username`, `password`, `license`)
VALUES (1, 'user-2', '$2y$10$z0m12MxuuGmCaz2f7NuFReUG6jD4U3O.6tmqlQCdCIjgbMtw0.qOm', 'fcee377a1fda007a8d2cc764a0a272e04d8c5d57');
INSERT INTO `arpanet_accounts`
(`enabled`, `username`, `password`, `license`)
VALUES (3, 'user-3', '$2y$10$z0m12MxuuGmCaz2f7NuFReUG6jD4U3O.6tmqlQCdCIjgbMtw0.qOm', 'db7e039146d5bf1b6781e7bc1bef31f0bb1298ea');
