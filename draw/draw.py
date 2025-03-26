import sys
import os
from PIL import Image
import epd13in3E

picdir = os.path.join(os.path.dirname(os.path.dirname(os.path.realpath(__file__))), 'pic')
libdir = os.path.join(os.path.dirname(os.path.dirname(os.path.realpath(__file__))), 'lib')

if os.path.exists(libdir):
    sys.path.append(libdir)

epd = epd13in3E.EPD()

try:
    epd.Init()
    print("clearing...")
    epd.Clear()

    Himage = Image.open(os.path.join(picdir, 'dither.bmp'))
    epd.display(epd.getbuffer(Himage))

    print("goto sleep...")
    epd.sleep()
except:
    print("exception raised...")
    epd.sleep()


