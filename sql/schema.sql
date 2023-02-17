CREATE TABLE `hours`
(
    hour         TIMESTAMP                                                 NOT NULL,
    availability ENUM ('available', 'not_available', 'training_scheduled') NOT NULL,
    PRIMARY KEY (hour)
);

CREATE TABLE `trainings`
(
    uuid             VARCHAR(32)  NOT NULL,
    user_uuid        VARCHAR(32)  NOT NULL,
    user             VARCHAR(255) NOT NULL,
    time             TIMESTAMP    NOT NULL,
    notes            TEXT         NOT NULL,
    proposed_time    TIMESTAMP,
    move_proposed_by VARCHAR(255),
    PRIMARY KEY (uuid)
);
