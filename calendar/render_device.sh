set -e

echo "Create bin"
docker compose exec dink go build -o dink_bin main/**.go

echo "Copy bin to pi"
scp ./dink_bin dink:/home/dink/prog/dink-calendar/calendar

# echo "Run"
# ssh dink 'cd /home/dink/prog/dink-calendar/calendar && ./dink_bin serve'

# scp dink:/home/dink/prog/dink-calendar/calendar/out/cal.png ./out/cal.png
