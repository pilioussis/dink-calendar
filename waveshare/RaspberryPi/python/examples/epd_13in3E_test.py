#!/usr/bin/python
# -*- coding:utf-8 -*-

import sys
import os
picdir = os.path.join(os.path.dirname(os.path.dirname(os.path.realpath(__file__))), 'pic')
libdir = os.path.join(os.path.dirname(os.path.dirname(os.path.realpath(__file__))), 'lib')
if os.path.exists(libdir):
    sys.path.append(libdir)


import epd13in3E
import time

from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont
from PIL import ImageColor

from PIL import Image

import json

epd = epd13in3E.EPD()
try:
    epd.Init()
    print("clearing...")
    epd.Clear()

    Himage = Image.open(os.path.join(picdir, 'dither.bmp'))
    epd.display(epd.getbuffer(Himage))
    input("Press Enter to continue...")

    print("clearing...")
    # epd.Clear()

    print("goto sleep...")
    epd.sleep()
except:
    print("exception raised...")
    epd.sleep()


