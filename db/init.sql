CREATE TABLE IF NOT EXISTS `students` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `email` varchar(255) UNIQUE NOT NULL,
  `suspended` boolean NOT NULL DEFAULT FALSE,
  INDEX (`email`)
);

CREATE TABLE IF NOT EXISTS `teachers` (
  `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `email` varchar(255) UNIQUE NOT NULL,
  INDEX (`email`)
);

CREATE TABLE IF NOT EXISTS `teachers_students` (
  `teacher_id` int NOT NULL,
  `student_id` int NOT NULL,
  PRIMARY KEY (`teacher_id`, `student_id`),
  INDEX (`teacher_id`, `student_id`)
);

INSERT INTO `students` (`email`) VALUES ('student1@gmail.com'), ('student2@gmail.com');
INSERT INTO `teachers` (`email`) VALUES ('teacherken@gmail.com'), ('teacherjoe@gmail.com');
