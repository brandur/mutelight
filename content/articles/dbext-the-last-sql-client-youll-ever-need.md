+++
location = "Calgary"
published_at = 2011-05-03T17:18:00-06:00
series_permalink = "masters-of-vim"
slug = "dbext-the-last-sql-client-youll-ever-need"
tiny_slug = "dbext"
title = "dbext: The Last SQL Client You'll Ever Need"
+++

As the first part of my [Masters of Vim series](/series/masters-of-vim), I'd like to introduce **dbext**, a Vim plugin that's potentially life-changing after fighting through the somewhat daunting initial configuration and usage.

This short article is no substitute for the extensive documentation that comes bundle with the plugin. Additional information and answers to any lingering questions can usually be found there.

Installation
------------

Install the plugin by downloading the [archive from Vim.org](http://www.vim.org/scripts/script.php?script_id=356), or clone the [Git mirror](https://github.com/vim-scripts/dbext.vim).

Additionally, the script itself relies on the command line clients of the various DBMSes that it supports. For dbext to work, refer to the table below and make sure that the binary appropriate for the system you'd like to work with is available. Often these clients are bundled with the DBMS itself, but otherwise should be available through your favorite package manager.

<table>
    <caption>Required command line client binaries for each DBMS supported by dbext</caption>
    <tr>
        <th style="width: 50%;">DBMS</th>
        <th style="width: 50%;">Client Bin</th>
    <tr>
    <tr>
        <td>ASA</td>
        <td><code>dbisql</code></td>
    </tr>
    <tr>
        <td>ASE (Sybase)</td>
        <td><code>isql</code></td>
    </tr>
    <tr>
        <td>DB2</td>
        <td><code>db2batch</code> or <code>db2cmd</code></td>
    </tr>
    <tr>
        <td>Ingres</td>
        <td><code>sql</code></td>
    </tr>
    <tr>
        <td>MySql</td>
        <td><code>isql</code></td>
    </tr>
    <tr>
        <td>Oracle</td>
        <td><code>sqlplus</code></td>
    </tr>
    <tr>
        <td>PostgreSQL</td>
        <td><code>psql</code></td>
    </tr>
    <tr>
        <td>Microsoft SQL Server</td>
        <td><code>osql</code></td>
    </tr>
    <tr>
        <td>SQLite</td>
        <td><code>sqlite</code></td>
    </tr>
</table>

Configuration
-------------

dbext needs the connection information of the database to operate on. It provides an interactive prompt to enter it, but I'd strongly recommend configuring your `.vimrc` instead.

Any number of profiles (connections) can be configured using this basic format:

``` vim
let g:dbext_default_profile_<profile_name> = '<connection string>'
```

For example:

```
" MySQL
let g:dbext_default_profile_mysql_local = 'type=MYSQL:user=root:passwd=whatever:dbname=mysql'

" SQLite
let g:dbext_default_profile_sqlite_for_rails = 'type=SQLITE:dbname=/path/to/my/sqlite.db'

" Microsoft SQL Server
let g:dbext_default_profile_microsoft_production = 'type=SQLSRV:user=sa:passwd=whatever:host=localhost'
```

These profile names will be used later to select a database on which to run query, so it's recommended that you make them somewhat logical.

Usage
-----

Open Vim pointing to a `*.sql` file. The SQL extension is not necessary to use dbext, but it's useful to get syntax highlighting. Write in a simple query relevant to your database:

``` sql
select * from user limit 100
```

Move your cursor to anywhere on the line you entered and type `<leader>sel` for `s` SQL, `e` execute, `l` line (from here on out I'm going to assume that `leader` is `\`). dbext will prompt you for a connection:

```
0. None
1. mysql_local
2. sqlite_for_rails
3. microsoft_production
[Optional] Enter profile #: 0
```

Enter the number corresponding to the connection for which your query will run with the results appearing in a split below (`C-w =` to even out the heights of each split).

```
+--------+----------+
| userID | username |
+--------+----------+
|      1 | bob      |
|      2 | joe      |
|      3 | jen      |
+--------+----------+
```

The `\se` (`s` SQL, `e` execute) command is useful for multiline queries. It searches backwards for the beginning of a query by looking for certain SQL command keywords (e.g. `SELECT`) and searches forwards for the end of a query by looking for your connection's command terminator (e.g. `;`). For example, the following SQL should execute with `\se` no matter where your cursor is on it:

``` sql
select *
from user 
limit 100;
```

`\st` selects everything from the table whose name is under your cursor (e.g. `user` in the previous example).

`\sT` selects from the table under your cursor, but prompts for the number of rows to select. This is a much safer alternative to `\st` when working with  a lot of data.

`\stw` selects from the table under your cursor, but prompts for a where clause (don't include the keyword `WHERE`).

`\sta` prompts for table name, then selects from that table.

### Schema

`\sdt` describes the table whose name is under your cursor:

```
+----------+-----------------------+
| Field    | Type                  |
+----------+-----------------------+
| userID   | mediumint(8) unsigned |
| username | varchar(30)           |
+----------+-----------------------+
```

`\sdp` is very similar, but instead describes the stored procedure under your cursor.

The three command directives `:DBListTable`, `:DBListProcedure`, and `:DBListView` list the database's tables, stored procedures, and views respectively. I find it useful to map the table list to its own key:

``` vim
map <leader>l :DBListTable<CR>
```

`\slc` copies each of the column names in the table under your cursor to the unnamed register in a format like `userID, username`. This is useful for constructing select queries on tables.

### The Results Buffer 

There are a few useful shortcuts specific to when your cursor is in the results buffer.

* `R` will re-run the command which populated the current results.
* `q` will quickly close the results.

Summary
-------

Recall that this article was written to be a very minimal introduction to dbext to highlight what are (in my opinion) some of its most useful features. Refer to dbext's excellent documentation to discover its full potential.
