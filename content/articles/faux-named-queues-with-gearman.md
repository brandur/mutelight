+++
location = "Calgary"
published_at = 2011-07-17T17:56:00-06:00
slug = "faux-named-queues-with-gearman"
tiny_slug = "named-gearman"
title = "Faux Named Queues with Gearman"
+++

A few months back, our architecture team deemed that a task queue service would be needed for a few upcoming projects, mostly so that we could replace our bad habit of running expensive background work via cron scripts. The software chosen for this role was [Gearman](http://gearman.org), a fast job queueing framework with good PHP support.

Gearman has been mostly pretty functional as far as our requirements were concerned, but came with one rather surprising quirk: any given Gearman server or cluster only supports a single job queue. I say that this is surprising because support for multiple queues is an integral feature of Gearman's alternatives such as RabbitMQ (channels), and Redis (different keyed sets).

Background
----------

To understand the reason that Gearman was built this way, a bit of background is in order. When a Gearman worker connects to a Gearman daemon, it registers the names of all the types of jobs it wants to be able to handle before entering its job loop and being sent work:

``` php
$worker = new GearmanWorker();
$worker->addServer();

// Tell the Gearman daemon that this worker will handle any jobs in its queue 
// called 'resizeImage' using the PHP function 'doImageResize'
$worker->addFunction('resizeImage', 'doImageResize');

while ($worker->work())
{
    // Finished processing a job
}

function doImageResize($job)
{
    // This is where the work happens
}
```

A result of the way this function registration works is that the Gearman daemon knows exactly which connected workers can handle which types of jobs. That is, if only worker #1 registered a function for `resizeImage`, then when a `resizeImage` job is encountered, it must be sent to worker #1.

Another aspect of Gearman's behaviour which isn't intuitively obvious is that if there's a job ready to be pulled off the queue, and has no connected workers that are able to handle it, that job stays in the queue ready to be run, but in the meantime other jobs in the queue continue to be processed. When a worker connects that can handle that previously unprocessable job, it gets passed off immediately.

In this way, jobs that belong to separate logical groups and which may have been stored in different queues in another system, can coexist in a single Gearman queue by simply having distinct Gearman workers.

Faux Named Queues
-----------------

Mostly due to internal IT logistics, we needed a number of people running the same application to be able to be able to share a single Gearman cluster. This proved to be a problem when person #1 would queue a job, and it would promptly be handled by person #2's worker.

Luckily, one of my colleagues realized that we could solve this problem by leveraging the same function registration behaviour detailed above. Clients could prepend a hostname before inserting any jobs:

``` php
$client = new GearmanClient();
$client->addServer();

doBackground($client, 'resizeImage', 'lolcat.jpg');

function doBackground($client, $job, $workload)
{
    if (!IS_PRODUCTION) {
        $job = $hostname . '-' . $job;
    }

    $client->doBackground($job, $workload);
}
```

Likewise, workers could prepend a hostname to any jobs they registered:

``` php
$worker = new GearmanWorker();
$worker->addServer();

addFunction($worker, 'resizeImage', 'doImageResize');

while ($worker->work())
{
    // Finished processing a job
}

function addFunction($worker, $job, $phpFunction)
{
    if (!IS_PRODUCTION) {
        $job = $hostname . '-' . $job;
    }

    $worker->addFunction($job, $phpFunction);
}
```

Using this technique, jobs are only processed by workers running under the same host as the client that queued them, effectively solving our sharing problem; and thanks to Gearman's behaviour when encountering unregistered jobs, users won't block each other, even if their worker isn't running.

A better solution might be to have everyone running local Gearman daemons, but this method serves as an effective alternative in cases where that's not possible.
