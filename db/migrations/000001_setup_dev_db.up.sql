CREATE USER proguser WITH PASSWORD 'progpass';

CREATE TYPE networktype AS ENUM ('evm', 'dot', 'sol', 'tm');

CREATE TABLE IF NOT EXISTS blocks(
    nwtype networktype NOT NULL,
    network smallint UNIQUE NOT NULL,
    height bigint NOT NULL,
    inserted_at timestamptz NOT NULL,
    synopsis jsonb NOT NULL,
    PRIMARY KEY(nwtype, network),
    UNIQUE(nwtype, network, height)
);

CREATE TABLE IF NOT EXISTS pool_actions_geth(
    nwtype networktype NOT NULL,
    network smallint UNIQUE NOT NULL,
    time timestamptz NOT NULL,
    inserted_at timestamptz NOT NULL,
    token0 VARCHAR(40) NOT NULL,
    token1 VARCHAR(40) NOT NULL,
    pair VARCHAR(40) NOT NULL,
    amount0 numeric NOT NULL,
    amount1 numeric NOT NULL,
    amountusd numeric,
    reserves0 numeric,
    reserves1 numeric,
    reservesusd numeric,
    type VARCHAR(20) NOT NULL,
    sender VARCHAR(40) NOT NULL,
    exchange VARCHAR(40),
    id VARCHAR(40),
    transaction VARCHAR(64) NOT NULL,
    slippage numeric,
    height bigint NOT NULL
);

GRANT SELECT, INSERT ON pool_actions_geth TO proguser;
GRANT SELECT, INSERT ON blocks TO proguser;