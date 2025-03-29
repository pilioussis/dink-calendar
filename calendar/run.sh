set -e

cd /home/dink/prog/dink-calendar/calendar && ./dink_bin

cp ./out/dither.bmp /home/dink/prog/13.3inch_e-Paper_E/RaspberryPi/python/pic/

cd /home/dink/prog/13.3inch_e-Paper_E/RaspberryPi/python/examples && python epd_13in3E_test.py

# nohup ./dink_bin serve > /home/dink/log/dink.log 2>&1 &