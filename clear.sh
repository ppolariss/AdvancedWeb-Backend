docker image ls | grep '<none>' | awk '{print $3}' | xargs docker image rm
