set -e

cd /home/dink/prog/dink-calendar/calendar && ./dink_bin
nohup ./dink_bin serve > /home/dink/log/dink.log 2>&1 &