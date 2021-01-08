+++
location = "San Francisco"
published_at = 2012-09-02T23:18:31-07:00
slug = "pg-max-connections"
title = "Max Connections for a Postgres Service"
+++

Once in a while, it's useful to know how many connections your Postgres service can support. For example, at Heroku we use this information to help alert us when any of our production-critical databases are approaching their connection limit.

Inspecting a Postgres configuration file will reveal a setting that specifies the maximum number of connections that its associated service will allow:

```
max_connections = 20
```

As with other settings, this can be checked by connecting to any running Postgres and executing the following query:

``` sql
select name, setting from pg_settings where name = 'max_connections';
```

**Protip:** you'll notice that for all our Postgres services at Heroku, from Dev to Ronin, and all the way to Mecha, the response will be `500`.
