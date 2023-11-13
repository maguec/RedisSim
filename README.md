# Redis Sim

This is a command line tool to simulate various scenarios running against a Redis compatable cluster.
The goal is to make it easy to show how a cluster responds under CPU pressure and Memory Pressure 

## Building

```
git clone https://github.com/maguec/RedisSim.git
cd RedisSim
# install go and make
sudo apt-get install -y golang-1.21 make
make
```

## Options 

To see what is available run

```
./RedisSim -h
```


## String Fill

Create a large number of keys in a Redis database of type [STRING](https://redis.io/docs/data-types/strings/)

For example the following create 1M keys in about 40 seconds

```
./RedisSim 	--port 30001 --clients 100 \
		--server localhost stringfill \
		--string-count 1000000 \
		--prefix stringprefix
```


## Exercise

Update a large number of keys in a Redis database of type [STRING](https://redis.io/docs/data-types/strings/)

For example the following read then update 100K keys two times in a loop

```
./RedisSim 	--port 30001 exercise  \
		--prefix stringprefix  --runs 2 \
		--key-count 100000 --ratio 1:1
```


## Load CSV

Load a CSV file into Redis type [HASH](https://redis.io/docs/data-types/hash/) and set the keyname 


If we have the following csv file

```
Name,Age,Profession,Team
Tim,35,Pitcher,Giants
Joe,71,Quarterback,49ers
Steph,34,Guard,Warriors
```

it can be loaded into Redis using the Name as the key with a prefix of athletes

```
./RedisSim 	loadcsv --port 30001  \
		--csv-file /tmp/test.csv \
		--key-field Name  --csv-prefix athletes
```


```
$ redis-cli -c -p 30001 hgetall athletes:Tim
1) "Name"
2) "Tim"
3) "Age"
4) "35"
5) "Profession"
6) "Pitcher"
7) "Team"
8) "Giants"
```





