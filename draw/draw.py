import sys
from PIL import Image
from lib import epd13in3E

epd = epd13in3E.EPD()

if len(sys.argv) > 1:
    filename = sys.argv[1]
else:
    print("No filename provided. Exiting.")
    sys.exit(1)


try:
    epd.Init()
    print("clearing panel")
    epd.Clear()

    Himage = Image.open(filename)
    epd.display(epd.getbuffer(Himage))

    print("sending panel to sleep")
    epd.sleep()
except:
    print("exception raised when drawing to panel")
    epd.sleep()


