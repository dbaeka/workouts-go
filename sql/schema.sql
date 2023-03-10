DROP TABLE IF EXISTS `hours`;
CREATE TABLE IF NOT EXISTS `hours`
(
    hour         TIMESTAMP                                                 NOT NULL,
    availability ENUM ('available', 'not_available', 'training_scheduled') NOT NULL,
    PRIMARY KEY (hour)
);

DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users`
(
    uuid         VARCHAR(32)                 NOT NULL,
    balance      INT                         NOT NULL DEFAULT 0,
    display_name VARCHAR(100)                NOT NULL,
    role         ENUM ('attendee','trainer') NOT NULL,
    last_ip      VARCHAR(16),
    PRIMARY KEY (uuid)
);


DROP TABLE IF EXISTS `dates`;
CREATE TABLE IF NOT EXISTS `dates`
(
    date           DATE    NOT NULL,
    has_free_hours BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (date)
);

DROP TABLE IF EXISTS `trainings`;
CREATE TABLE IF NOT EXISTS `trainings`
(
    uuid             VARCHAR(32)  NOT NULL,
    user_uuid        VARCHAR(32)  NOT NULL,
    user             VARCHAR(255) NOT NULL,
    time             TIMESTAMP    NOT NULL,
    notes            TEXT         NOT NULL,
    proposed_time    TIMESTAMP,
    move_proposed_by VARCHAR(255),
    canceled         BOOLEAN      NOT NULL DEFAULT FALSE,
    PRIMARY KEY (uuid)
);
