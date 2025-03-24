# One more e-ink calendar
Hello. I made another e-ink calendar. If you like display tech as much as I do, you might find this interesting.

Let's start with a video of its leisurely 19s refresh cycle.

https://github.com/user-attachments/assets/c472b0ad-aced-4d1a-b8ad-b3a8ebeae87f


## What is it?
- A spiffy new (as of Nov 2024) Spectra 6 E-ink panel, driven by a Pi + HAT.
- Integrates with Google Calendar.
- Battery powered with an RTC to save power.
- Written in Go using chromium as the renderer.

## The issue with sharing events
Google Calendar have optimized their UX for business, not me and my lovely wife.
We want to be able to see:
- Everything I'm doing
- Everything she's is doing
- Everything we are doing together

Solutions like share personal calendars, creating a 3rd shared calendar and shared events _almost_ get me what I need, except:
- Shared events show duplicates when personal calendars are shared to one another.
- No invite notifications when using a 3rd shared calendar that we both own.
- I think you can do it with extra google accounts, but nobody wants that.

## My beef with calendar UX
- Not a fan of the month view. I want to see the next 4+ weeks regardless of where we are in the month. I don't need to see past weeks.
- No need to waste precious e-ink space on weekday names. When [this Spectra 6 panel](https://buy-lcd.com/products/gdep253c02) comes down in price I'll consider it.

## Other learnings
- Using an RTC that completely cuts the power to the pi is surprisingly efficient. I wake the pi up 4 times a day, do a render (if necessary), then send it back to sleep. Should get at least a month on a 18Ah battery.
- The display refresh is comically slow, and makes a hissing noise.

## The final product
Here's a photoshopped pic I made to try different frames before I make one.
![photoshopped render](/img/render.jpg)

## A sad story of what could have been
While I was waiting for the panel to ship, I had the tremendous idea to make the calendar look unlike a "techy device". Like in a 1940's school classroom type style, with washed out fountain pen ink on weathered paper.

With the infrequently used "burn" blend mode in CSS, you can achieve a washed-out ink look that retains the texture of the underlying paper in a (kinda) realistic way.

![web render](/img/paper-full-color-cropped.png)

I implemented it all before the panel arrived, knowing I'd have to dither that image into the meagre 6-color space of Spectra e-ink. Alas I goofed in the following ways:
- Assuming **00ff00** green would look like green. Its a "dark seaweed" color, at best. The manufacturer is definitely doing tricks with [the demo pictures of vegetables](https://img201.yun300.cn/repository/image/c048f814-6525-4e4d-a2d0-cd9d67d80459.jpg_640xaf.jpg?tenantId=160096&viewType=1&k=1734319076000).
- Assuming the panel could render a stepped gradation between white and the 5 other colors. eg a one step gradation being white **ffffff** → light blue (the step) **8888ff** → blue **0000ff**.

Here's the 2-step dither I thought I could get away with.
![2 step render](/img/dither-2.png)

Here's the dither restricted to the 6-Color space:
![0 step render](/img/dither-1.png)

And here's the dither after the e-ink panel massacred it:
![massacred render](/img/massacred.jpeg)

Please let me know when e-ink tech gets better and this ye-olde design is possible.