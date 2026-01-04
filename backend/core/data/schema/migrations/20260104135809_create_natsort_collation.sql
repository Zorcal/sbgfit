-- migrate:up
-- ICU collation for case insensitive natural sorting
CREATE COLLATION IF NOT EXISTS public.natsort (
    provider = icu,
    locale = 'und-u-kn-ks-level2',
    deterministic = false
);


-- migrate:down
DROP COLLATION IF EXISTS public.natsort;
