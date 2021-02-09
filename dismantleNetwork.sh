# stopping the containers. This will stop all the running containers. 
docker stop $(docker ps -q)

# removing unused container,images,network and volumes 
docker system prune --volumes

# remove all certificates created.
cd channel-artifacts && rm -rf *

cd ../crypto-config && sudo rm -rf *

##===================================================================================
#  docker volume prune
# docker volume rm $(docker volume ls -q)
## Delete orphan container also
# removing all the images.Note images are created when we run 
# docker rmi $(docker images dev-* -q)
# removing all stoped containers. 
#docker rm $(docker ps -aq)