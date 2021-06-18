build:
	docker build --force-rm -t tnsmith/yotta-test .
push:
	docker push tnsmith/yotta-test
start:
	docker-compose up -d yotta
	-docker exec --user ydbadm -it trn-micro_yotta_1 bash -c ". /ydbdir/ydbenv && mupip set -rec=500000 -reg DATA" 
	-docker exec --user ydbadm -it trn-micro_yotta_1 bash -c ". /ydbdir/ydbenv && mupip set -key=500 -reg DATA" 
	-docker exec --user ydbadm -it trn-micro_yotta_1 bash -c ". /ydbdir/ydbenv && mupip set -null_subscripts=TRUE -region DATA" 
	docker-compose up -d app
stop:
	docker-compose down --volumes
	docker container prune -f
	docker volume prune -f
	docker image prune -f