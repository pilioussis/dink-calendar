set -e
echo "Start"
docker compose exec dink go run main/**.go

echo "Copy pic to pi"
scp ./out/dither.bmp dink:/home/dink/prog/e-Paper/RaspberryPi_JetsonNano/python/pic/

echo "Run command"
ssh dink 'cd /home/dink/prog/e-Paper/RaspberryPi_JetsonNano/python/examples && python3 epd_7in3e_test.py'

# scp /Users/sarahmcpherson/prog/dink-calendar/waveshare-pi-demo/python/examples/epd_7in3e_test.py dink:/home/dink/prog/e-Paper/RaspberryPi_JetsonNano/python/examples/