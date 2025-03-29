set -e

# echo "Update python demo file"
# scp /Users/dean/prog/di/nk-calendar/waveshare/RaspberryPi/python/examples/epd_13in3E_test.py dink:/home/dink/prog/13.3inch_e-Paper_E/RaspberryPi/python/examples/
# scp /Users/dean/prog/dink-calendar/waveshare-pi-demo/python/examples/epd_7in3e_test.py dink:/home/dink/prog/e-Paper/RaspberryPi_JetsonNano/python/examples/

echo "Generating calendar"
docker compose exec dink go run main/**.go

echo "Copy pic to pi"
scp ./out/dither.bmp dink:/home/dink/prog/dink-calendar/calendar/out/

echo "Run command"
ssh dink 'cd /home/dink/prog/dink-calendar/draw && python draw.py /home/dink/prog/dink-calendar/calendar/out/dither.bmp'

