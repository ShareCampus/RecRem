DROP TABLE IF EXISTS `openai`;

CREATE TABLE
    `openai` (
        `id` int NOT NULL AUTO_INCREMENT,
        `token` varchar(64) NOT NULL,
        `model` varchar(64) NOT NULL,
        PRIMARY KEY(`id`)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
