docker compose -p awb pull
docker compose -p awb down
docker volume rm awb_shared_volume
docker compose -p awb up -d
# cd Desktop/AdvancedWeb-Backend/
# docker compose -p awb pull
# docker volume rm awb_shared_volume
# docker compose -p awb up -d --build