#!/bin/sh
#usage precopy.sh <ContainerName> <RemoteLocation> <page-server port>
if [ "$#" -ne 3 ] ; then 
	echo "usage precopy.sh <ContainerName> <RemoteLocation> <page-server port>"
	exit 1
fi
#renaming & formatting variables
containerName=$1
remoteHost=$2
port=$3
containerID=$( docker inspect --format "{{.ID}}" $1 )
containerImageName=$( docker ps -a --filter name=$containerName --format \"{{.Image}}\" )
username=${remoteHost%@*}
ip=${remoteHost#*@}
containerImageName="${containerImageName%\"}"
containerImageName="${containerImageName#\"}"
check=$( ssh -i /root/.ssh/id_rsa $remoteHost docker images --format \"{{.Repository}}\"  $containerImageName)

#checking if image exists on remote host

if [ "$check" != "$containerImageName" ]
then
	echo "image not present on remote host"
	exit 1;
fi

#making directories for saving

mkdir memimage
directory=$( ssh $remoteHost pwd )
ssh $remoteHost mkdir memimage

#pre-dump 20 times
cd memimage
pwd
echo "pre-dump 1"	
ssh $remoteHost mkdir memimage/1
mkdir ./1
ssh $remoteHost nohup criu page-server --auto-dedup --images-dir $directory/memimage/1 --track-mem --port $port &
sleep 2
docker checkpoint --leave-running --predump --track-mem --image-dir=$( pwd )/1  --page-server --address=$ip --port=$port $containerName
scp $(pwd)/1/* $remoteHost:$directory/memimage/1/
sleep 2

for i in `seq 2 10`; 
do
	echo "pre-dump $i"
	mkdir $i
	ssh $remoteHost mkdir memimage/$i
	ssh $remoteHost nohup criu page-server --auto-dedup --images-dir $directory/memimage/$i --prev-images-dir ../$(( $i - 1 )) --port $port &
	sleep 2
	docker checkpoint --leave-running --predump --track-mem --image-dir=$( pwd )/$i --track-mem --prev-image-dir=$( pwd )/$(( $i - 1 )) --page-server --address=$ip --port=$port $containerName
	scp $(pwd)/$i/* $remoteHost:$directory/memimage/$i
	sleep 2
done

#final dump

mkdir final
ssh $remoteHost mkdir memimage/final
ssh $remoteHost nohup criu page-server --images-dir $directory/memimage/final --prev-images-dir ../3 --port $port &
sleep 2
docker checkpoint --image-dir=$( pwd )/final --prev-image-dir=$( pwd )/3 --page-server --address=$ip --port=$port $containerName

#sending rest of the files to the server

scp $(pwd)/final/* $remoteHost:$directory/memimage/final

#making tar of disk filse of container and sending it

foldername=$( cat /var/lib/docker/0.0/image/aufs/layerdb/mounts/$containerID/mount-id )
foldername2=$foldername"-init"
tar -cvzf imagedata.tar.gz -C /var/lib/docker/0.0/aufs/diff/ $foldername $foldername2
scp $(pwd)/imagedata.tar.gz $remoteHost:$directory/memimage

#restarting container

newContainerID=$( ssh $remoteHost docker create --name=$containerName $containerImageName )
newFolderName=$( ssh $remoteHost cat /var/lib/docker/0.0/image/aufs/layerdb/mounts/$newContainerID/mount-id )
newFloderName2=$newFolderName"-init"
ssh $remoteHost tar -xvzf $directory/memimage/imagedata.tar.gz
ssh $remoteHost cp -rv $directory/$foldername/*  /var/lib/docker/0.0/aufs/diff/$newFolderName/
ssh $remoteHost cp -rv $directory/$foldername2/*  /var/lib/docker/0.0/aufs/diff/$newFloderName2/
ssh $remoteHost docker restore --force --image-dir=$directory/memimage/final $containerName 

#housekeeping and deleting the files
cd ..
rm -rf $( pwd )/memimage
ssh $remoteHost rm -rf $directory/memimage
ssh $remoteHost rm -rf $directory/$foldername
ssh $remoteHost rm -rf $directory/$foldername2
