--
-- Table structure for table `users`
--
CREATE TABLE `users` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL,
    `role` VARCHAR(50) NOT NULL DEFAULT 'Agent',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `password_reset_token` VARCHAR(255),
    `password_reset_expires` DATETIME
);


--
-- Table structure for table `tickets`
--
CREATE TABLE `tickets` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `title` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'Open',
    `priority` VARCHAR(50) NOT NULL DEFAULT 'Medium',
    `requester_id` INT NOT NULL,
    `assignee_id` INT,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `closed_at` DATETIME,
    FOREIGN KEY (`requester_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`assignee_id`) REFERENCES `users`(`id`) ON DELETE SET NULL ON UPDATE CASCADE
);

--
-- Table structure for table `comments`
--
CREATE TABLE `comments` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `ticket_id` INT NOT NULL,
    `author_id` INT NOT NULL, -- Assuming comments are primarily from agents (users)
    `body` TEXT NOT NULL,
    `is_internal` BOOLEAN DEFAULT FALSE,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`ticket_id`) REFERENCES `tickets`(`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (`author_id`) REFERENCES `users`(`id`) ON DELETE CASCADE ON UPDATE CASCADE
);
