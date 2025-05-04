CREATE TABLE IF NOT EXISTS
    profiles (
        "guid" BYTEA,
        "name" TEXT NOT NULL,
        "surname" TEXT NOT NULL,
        "patronymic" TEXT,
        "age" INT,
        "gender" TEXT,
        "nationalize" TEXT
    );

CREATE INDEX name_and_surname ON profiles("name", "surname");