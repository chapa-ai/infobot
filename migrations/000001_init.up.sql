CREATE TABLE IF NOT EXISTS "currencies"
(
    "id" SERIAL UNIQUE NOT NULL,
    "symbol" text NOT NULL,
    "buy" text NOT NULL,
    "time" timestamp  NOT NULL,
    CONSTRAINT "Currencies_pkey" PRIMARY KEY ("id")
    );