+++
location = "Calgary"
published_at = 2009-03-16T00:00:00-06:00
slug = "how-to-play-1080p-hd-video-encoded-with-x264-in-an-mkv-container"
tiny_slug = "6"
title = "How to Play 1080p HD Video Encoded with x264 in an MKV Container"
+++

Full 1080p HD video in x264/H.264 is notoriously difficult to properly decode, due to both its enormous resolution and high compression. The fact that a huge number of these videos are stored in MKV containers doesn't help the situation, mostly because some common video players don't read MKVs as optimally as possible.

The past few years have seen a plateau in terms of CPU clockspeed, so we haven't been able to rely on ever-increasing processing power to solve the problem. Instead we have to be smart about how we play our HD videos, by moving decoding to our video cards or splitting the work across multiple cores.

This article contains a few tips to play your HD videos the right way (or at least the best ways I've found to date). I've labeled the article with 1080p, but these suggestions work just as well with 720p HD video or lower.

Networked Media Tank
--------------------

If you're willing to try a hardware solution, you can save yourself quite a bit of time installing codecs and messing around with configuration settings on your PC. The idea of a box attached to your TV that can stream media over your network has been around for a while, but it's something that in my opinion, has only really come to fruition recently.

For a while, a hacked Xbox running <acronym title="Xbox Media Center">XBMC</acronym> was far and away the best media player on the market. Unfortunately, the Xbox lacks the horsepower to play any real high-definition material, and newer consoles have never been as useful for media playback as the old Xbox was. Many alternatives are available today: modded Apple TV, Xbox 360, PS3, or your own HTPC, but these options are expensive, lack codec support, or have no chance of hitting 1080p playback. Nowadays, none of these options need even be considered after the relatively recent appearance of a new class of device on the market: the networked media tank (NMT).

Networked media tanks really are just as awesome as their name suggests. These devices are small computers that run [processors made by Sigma Designs](http://www.sigmadesigns.com/public/Products/selection_guide/selection_guide.html) specifically for jobs like decoding HD media (although apparently the official definition of an NMT is a box that runs [Syabas middleware](http://www.syabas.com/solution_nmt.html).) The first thing you'll notice when trying a 1080p x264 MKV on an NMT is just how smooth the playback is, either that, or sheer amazement at how a commercial device is playing an MKV out of the box (that is if you're used to pretty but borderline useless Apple TVs and such).

The most popular NMTs are the [Popcorn Hour](http://www.popcornhour.com/) and the [HDX](http://www.hdx1080.com/) (see a [full list of NMTs](http://www.networkedmediatank.com/wiki/index.php/Products),) both of which are sleek little boxes that can stream over the network, play from a USB drive, or play from a user-installed internal HDD. You can install other interesting add-ons like a BitTorrent client, NZB downloader or a tool to grab cover art for you. These features are all great, but what's the best part? NMTs tend to weigh in just a little over $200 US, in most cases below the comparatively inadequate competition from big names. Personally, I've only had direct experience with the HDX 1000, and found it absolutely outstanding.

Windows
-------

The one HD playback solution under Windows that I'm going to talk about is hardware decoding, where your heavy-lifting is offloaded to your DXVA(DirectX Video Acceleration)-enabled video card. We'll be using Media Player Classic Homecinema so you'll need at least an Nvidia Gefore 8 series, or ATI Radeon HD series card.

Follow these steps to get up and running:

1. Download and install your video card's latest drivers ([Nvidia drivers download page](http://www.Nvidia.com/Download/index.aspx?lang=en-us))
# Download and install [Haali's Media Splitter](http://haali.cs.msu.ru/mkv/), a very fast MKV reader
2. Download and install [Media Player Classic Homecinema](http://mpc-hc.sourceforge.net/) (MPC-HC), a player bundled with DXVA support
3. Under _Options_ &rarr; _Playback_ &rarr; _Output_, choose _VMR9 (renderless)_ if you're on Windows XP (as seen in Fig. 1 below) or _EVR_ if you're on Windows Vista (see [more information on MPC-HC DXVA support](http://mpc-hc.sourceforge.net/DXVASupport.html))
4. Under _Options_ &rarr; _Internal Filters_ &rarr; _Source Filters_, uncheck the options for _Matroska_. We do this to allow Haali's Media Splitter to read our MKV files, which is faster than MPC-HC's reader.

<div class="figure">
    <a href="https://d25zpof2afwnhk.cloudfront.net/how-to-play-1080p-hd-video-encoded-with-x264-in-an-mkv-container/mpc-hc-options-vmr9.png" title="Link to full-size image"><img src="https://d25zpof2afwnhk.cloudfront.net/how-to-play-1080p-hd-video-encoded-with-x264-in-an-mkv-container/mpc-hc-options-vmr9-small.png" alt="Settings to correctly enable VMR9 in Media Player Classic Homecinema" /></a>
    <p><strong>Fig. 1:</strong> <em>MPC-HC output settings for hardware-accelerated VMR9 playback in Windows</em></p>
</div>

Mac OS X
--------

This one is easy: use [Plex](http://plexapp.com/). Plex (also known as Plexapp) is a fork of XBMC for Intel-based Macintosh computers, and it is _awesome_. Apart from its beautiful and highly-intuitive interface, Plex will use multiple cores to decode HD video, allowing you to play 1080p videos even on your notebook.

Another option, but one I've admittedly never had much luck with, is to use VLC. VLC will be too slow to play HD video out of the box, but you can configure it to skip its x264 loop filter as shown in Fig. 2 below (remember to select the _All_ option from the radio buttons in the bottom left or you won't see these settings). Depending on your processor, this may speed up VLC enough to make it usable.

<div class="figure">
    <a href="https://d25zpof2afwnhk.cloudfront.net/how-to-play-1080p-hd-video-encoded-with-x264-in-an-mkv-container/vlc-options-skip-loop-filter.png" title="Link to full-size image"><img src="https://d25zpof2afwnhk.cloudfront.net/how-to-play-1080p-hd-video-encoded-with-x264-in-an-mkv-container/vlc-options-skip-loop-filter-small.png" alt="Settings to get better VLC performance on Mac OSX by skipping the loop filter" /></a>
    <p><strong>Fig. 2:</strong> <em>VLC settings for skip loop filter on Mac OSX</em></p>
</div>

Hardware decoding support has begun to appear in Mac OSX as well. Unfortunately, as things stand today, QuickTime is the only application able to access this functionality, and installing Perian to give QuickTime access to decent codec/container support will break hardware decoding (so you can't win). Many people are hoping that Apple's upcoming release of Mac OS X 10.6 (Snow Leopard), which is supposed to move a lot of computing to the <acronym title="Graphics Processing Unit">GPU</acronym>, will resolve this problem.

Linux
-----

As of November 2008, the proprietary Nvidia drivers for Linux (and other Unix-based systems) support video hardware acceleration in the form of an API called Video Decode and Presentation API for Unix (VDPAU). VDPAU is a Unix equivalent to DXVA on Windows and support for it has already been added or is being added to many popular Linux media players. See the [Wikipedia article on VDPAU](http://en.wikipedia.org/wiki/VDPAU) for more information on supporting players.

I haven't tried VDPAU out for myself so I'll hold off on any specific instructions for Linux until I get the chance.
