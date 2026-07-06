-- +goose Up

-- enum-типы: значения буква-в-букву доменные константы
-- (domain.AccountStatus / domain.AccountType) — иначе Reconstitute* не смаппит.
CREATE TYPE account_status AS ENUM ('Active', 'Frozen', 'Closed');
CREATE TYPE account_type   AS ENUM ('UserAccount', 'SystemAccount');

-- accounts: снапшот-агрегат (баланс + version). id = TEXT (1:1 с domain.AccountID).
CREATE TABLE accounts (
    id             TEXT           PRIMARY KEY,
    user_id        TEXT,                     -- NULL у system-счёта
    currency       TEXT           NOT NULL,  -- код валюты; minorUnits/symbol живут в registry
    balance_amount BIGINT         NOT NULL,  -- минимальные единицы, signed; никогда не float
    status         account_status NOT NULL,
    account_type   account_type   NOT NULL,
    version        BIGINT         NOT NULL,  -- оптимистичная блокировка (проверяется в pgx-репо)
    created_at     TIMESTAMPTZ    NOT NULL   -- доменное время, пишем явно (без DEFAULT now())
);

-- один system-счёт на валюту — страхует FindSystemAccount от двух контр-счётов в одной валюте.
CREATE UNIQUE INDEX ux_accounts_system_per_currency
    ON accounts (currency)
    WHERE account_type = 'SystemAccount';

-- transactions: иммутабельная проводка. key — клиентский ключ идемпотентности.
CREATE TABLE transactions (
    id          TEXT        PRIMARY KEY,
    key         TEXT        UNIQUE,               -- NULL = без ключа; NULL'ы различны → повторы не конфликтуют
    description TEXT        NOT NULL DEFAULT '',  -- пустая строка валидна; Go-строка не бывает nil
    created_at  TIMESTAMPTZ NOT NULL
);

-- postings: строки проводки (VO). Σ amount внутри одной transaction = 0 —
-- межстрочный инвариант домена, в DDL не enforce'им (живёт в агрегате Transaction).
CREATE TABLE postings (
    id             BIGSERIAL PRIMARY KEY,   -- порядок постингов при чтении: ORDER BY id
    transaction_id TEXT      NOT NULL REFERENCES transactions (id),
    account_id     TEXT      NOT NULL REFERENCES accounts (id),
    amount         BIGINT    NOT NULL CHECK (amount <> 0),  -- зеркалит ErrZeroPosting (per-row инвариант)
    currency       TEXT      NOT NULL
);

-- загрузка постингов транзакции при FindByID.
CREATE INDEX ix_postings_transaction_id ON postings (transaction_id);

-- +goose Down

-- зеркально Up: сначала таблицы (индексы уходят вместе с ними), потом типы.
-- postings первым — зависит по FK от transactions и accounts.
DROP TABLE postings;
DROP TABLE transactions;
DROP TABLE accounts;
DROP TYPE account_type;
DROP TYPE account_status;
