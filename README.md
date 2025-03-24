# One more E-ink calendar
Hello, its me. I didn't like the other calendars out there so I made my own. If you like bleeding edge E-ink as much as I do, you might find this interesting.

### Let's lead with the piccies
These panels have some of the best color representation (as of Nov 2024), although their refresh time is sluggish. (Not a problem if you only need to re-render a couple of times a day)
[![Watch the video](/img/render.jpg)](https://raw.githubusercontent.com/pilioussis/dink-calendar/main/img/output.mp4)

Photoshopped this while I make my frame.
![photoshopped render](/img/render.jpg)

Another go
[![Watch the video](/img/render.jpg)](/img/output.mp4)


### What is it?
- A spiffy new (as of Nov 2024) Spectra 6 E-ink panel, driven by a Pi + HAT.
- Integrates with Google Calendar.
- Battery powered with an RTC to save power.
- Tries not to look like a normal E-ink display.
- Written in Go using chromium as a renderer.

### What's the issue with google cal?
I like Google Calendar on my phone, but they have optimised their UX for business, not me and my lovely wife.
I want to be able to see:
- Eveything I'm doing
- Everything my wife is doing
- Everything we are doing together

Shared calendars, shared events and shared accounts don't quite cut the mustard. There's some sort of CAP theorem going on where:
- Shared events show duplicates (when sharing calendars, shared events appear twice)
- I get notiications (if using a 3rd shared calendar that be both own)
- Isn't asking my lovely wifey to do a short course in

### What is your beef with calendar UX?
- I know the names of the weekdays, I'll save the $1000/cm² e-ink space for my precious events.
- The month view makes me grumpy. I want to see the next 4+ weeks always. I don't want to see the weeks in the past.

### A sad story (Do not read this if you are easily rustled)
In an effort to make it look unlike a "techy device" I initially tried to go for a 1940's school classroom style, with washed out fountain pen ink on weathered paper.
![photoshopped render](/img/paper-full-color-cropped.png)

I used the rarely used "burn" blend mode in CSS, which you don't see often outside of fancy photo/video editinig software. With the right color spacing, you can acheive a washed-out ink look that retains the texture of the underlying paper in a (somewhat) realistic way.

This was all before I got the panel. I knew I would have to dither it into the 6-color space of Spectra E-ink. But alas I goofed in the following ways:
- Assuming the panel could render a gradation between white and the 6 colors. eg #ffffff (white) slowly turning into #0000ff (blue). No! Only #ffffff or #0000ff exactly.
- #00ff00 green would look like green. Alas it is a "washed up winter seaweed"
- Believing the demo images. They have been carefully selected to not show certain colors, and then color corrected in post. 

The 2-step dither I thought I was allowed (eg two color extra colors between #ffffff and #0000ff):
![photoshopped render](/img/dither-2.png)

The dither I was actually allowed (no steps):
![photoshopped render](/img/dither-1.png)

The dither after the e-ink panel masacred it:
![photoshopped render](/img/masacre.jpeg)

I'm sure with some calibration you may be able to get that looking closer to the original, but I chose to wait until color E-ink gets better and this ye-oldy design is possible.







