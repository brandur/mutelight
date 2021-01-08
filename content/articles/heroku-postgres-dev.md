+++
location = "Berlin"
published_at = 2012-06-07T14:44:00+02:00
slug = "heroku-postgres-dev"
title = "Upgrade to Heroku Postgres Dev"
+++

Our [Heroku Postgres Dev](https://postgres.heroku.com/) plans recently went into public beta. They're very exciting work, providing the full power of a true Postgres database for development applications, and for free!

Some key features of the new dev plan is that databases under it are 9.1 (up from 8.3 which is what the shared databases ran under), support hstore, and can be managed remotely using `heroku pg:psql` or any other Postgres client.

However, since the dev plan adds a brand new database, the default is to end up with an empty store with none of your previous application data. If you're like me, and not too familiar with Heroku Postgres, it might not be immediately obvious how to seemlessly get your data migrated over. Lucky for you though, you're on Heroku! Using pgbackups, there's a very simple way to move your data between databases and produce a backup as a convenient byproduct.

Add the `pgbackups` addon and capture a backup of your current shared database:

    heroku addons:add pgbackups
    heroku pgbackups:capture

The Heroku command will tell you that a backup was produced with a name like `b001`. Now add your new Postgres dev database:

    heroku addons:add heroku-postgresql:dev

The name of your new database will come back as a token like `HEROKU_POSTGRESQL_CYAN`. It's attached to your app, but not yet acting as its primary database.

Now all that's left to do is restore the backup you made, and make it your primary:

    heroku pgbackups:restore HEROKU_POSTGRESQL_CYAN b001
    heroku pg:promote HEROKU_POSTGRESQL_CYAN

Open a psql session to the new Postgres dev instance and check that all your data is properly in place:

    heroku pg:psql

Optionally, you can destroy your old shared database:

    heroku addons:remove shared-database
