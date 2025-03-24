# One more E-ink calendar
Hello, its me. I didn't like the other calendars out there so I made my own. If you like bleeding edge E-ink as much as I do, you might find this interesting.

## Let's lead with the piccies
Here's a video of a refresh cycle.

https://github.com/user-attachments/assets/c472b0ad-aced-4d1a-b8ad-b3a8ebeae87f

Haven't made my frame yet, here's a photoshopped pic.
![photoshopped render](/img/render.jpg)

## What is it?
- A spiffy new (as of Nov 2024) Spectra 6 E-ink panel, driven by a Pi + HAT.
- Integrates with Google Calendar.
- Battery powered with an RTC to save power.
- Look (sort of) different to a normal E-ink display.
- Written in Go using chromium as a renderer.

## What's the issue with google cal?
Google Calendar have optimised their UX for business, not me and my lovely wife.
I want to be able to see:
- Eveything I'm doing
- Everything my wife is doing (not weird I promise, makes it easier to buy tickets)
- Everything we are doing together

Solutions like share personal calendars, creating a 3rd shared calendar and shared events _almost_ gets me what I need, except:
- Shared events show duplicates when sharing personal calendars to one another with shared events.
- No notifications when using a 3rd shared calendar that we both own.

## My beef with google calendar UX
- Not a fan of the month view. I want to see the next 4+ weeks regardless of month cutoff. I don't need to see past weeks.
- I know the names of the weekdays, no need to waste $9000/cm² e-ink space on it.

## Other learnings
- Using an RTC that completely cuts the power to the pi is surprisingly efficient. I wake the pi up 3 times a day for a quick render, then send it back to sleep. Should get at least a month on a single battery.
- The display refresh is comically slow, and makes a hissing noise.


## A sad story (do not read this if you are easily rustled)
Initially I had the tremendous idea to make the calendar look unlike a "techy device". Like in a 1940's school classroom type style, with washed out fountain pen ink on weathered paper.

Using the "burn" blend mode in CSS with the right color spacing, you can acheive a washed-out ink look that retains the texture of the underlying paper in a (somewhat) realistic way. A fun little trick which you don't often come across outsie of photo/video editinig software. I was pretty happy with how that turned out.

![photoshopped render](/img/paper-full-color-cropped.png)

I implemented this all before I got the panel, knowing I'd have to dither that image above into the 6-color space of Spectra E-ink. Alas I goofed in the following ways:
- Assuming **00ff00** green would look like green. Its a "dark seaweed" type color. The vegetables in the demo picture have certainly been exaggerated.
- Assuming the panel could render a stepped gradation between white and the 5 other colors. eg a one step gradation being white **ffffff** → light blue **8888ff** → blue **0000ff**.

Here's the 2-step dither I thought I was allowed.
![photoshopped render](/img/dither-2.png)

Here's the dither I was actually allowed (no steps):
![photoshopped render](/img/dither-1.png)

Here's the dither after the e-ink panel masacred it:
![photoshopped render](/img/masacre.jpeg)

I'm sure with some calibration you may be able to improve things, but I'll wait until color E-ink gets better and this ye-oldy design is possible.