# One more E-ink calendar
Hello, its me. I didn't like the other calendars out there so I made my own. If you like bleeding edge E-ink as much as I do, you might find this interesting.

### Let's lead with the piccies
Here's a video of a refresh cycle.


https://github.com/user-attachments/assets/c472b0ad-aced-4d1a-b8ad-b3a8ebeae87f


Haven't made my frame yet, here's a photoshopped pic.
![photoshopped render](/img/render.jpg)

### What is it?
- A spiffy new (as of Nov 2024) Spectra 6 E-ink panel, driven by a Pi + HAT.
- Integrates with Google Calendar.
- Battery powered with an RTC to save power.
- Tries not to look like a normal E-ink display.
- Written in Go using chromium as a renderer.

### What's the issue with google cal?
I like Google Calendar on my phone, but they've optimised their UX for business, not me and my lovely wife.
I want to be able to see:
- Eveything I'm doing
- Everything my wife is doing
- Everything we are doing together

Sharing personal calendars, creating a 3rd shared calendar and shared events almost get me what I want, except:
- Shared events show duplicates (when sharing personal calendars, shared events appear twice)
- No notifications (if using a 3rd shared calendar that we both own)

### My beef with google calendar UX
- I know the names of the weekdays, I'll save the $1000/cm² e-ink space for my precious events.
- The month view makes me grumpy. I want to see the next 4+ weeks always. I don't want to see the weeks in the past.

### A sad story (do not read this if you are easily rustled)
In an effort to make it look unlike a "techy device" I initially tried to go for a 1940's school classroom style, with washed out fountain pen ink on weathered paper.
![photoshopped render](/img/paper-full-color-cropped.png)

Using the "burn" blend mode in CSS with the right color spacing, you can acheive a washed-out ink look that retains the texture of the underlying paper in a (somewhat) realistic way. A fun little trick which you don't often come across outsie of photo/video editinig software.

This idea and solution was all before I got the panel. I knew in advance I'd have to dither it into the 6-color space of Spectra E-ink. But alas I goofed in the following ways:
- Assuming #00ff00 green would look like green. Alas it is a "dark seaweed" type color.
- Assuming the panel could render some gradation between white and the 6 colors. eg #ffffff (white) -> #8888ff (light blue) -> #0000ff (blue).
- Assuming the demo images weren't being sneaky. They're been carefully selected to show certain colors, and then color corrected in post.

Here's the 2-step dither I thought I was allowed.
![photoshopped render](/img/dither-2.png)

Here's the dither I was actually allowed (no steps):
![photoshopped render](/img/dither-1.png)

Here's the dither after the e-ink panel masacred it:
![photoshopped render](/img/masacre.jpeg)

I'm sure with some calibration you may be able to improve things, but I'll wait until color E-ink gets better and this ye-oldy design is possible.