## Tokoku
![banner](https://media.discordapp.net/attachments/761056621849477160/1217157858492289124/cat_anime-girl.png?ex=660c3c17&is=65f9c717&hm=75326166d5e93a3455fabb89b7bc9b429d7957ef836c02f1d1aee92de21a79c9&=&format=webp&quality=lossless&width=935&height=526)

## Background Project
This project is just an example of a project regarding Microservices. Actually, I've made it before, you can see it here [**Alpha**](https://github.com/rulanugrh/alpha) but this project is an improvisation from before. I also added several features as a result of my learning, and there is additional observability of each existing service possible :u

## Installation
First you have to do 3 steps below
```bash
# clone my project into your local
$ git clone https://github.com/rulanugrh/tokoku

# after clone you can change to folder tokoku
$ cd tokoku

# because this is microservices, you can copy all .env.example to .env
# product service
$ cp product/.env.example product/.env

# user service
$ cp user/.env.example user/.env

# order service
$ cp order/.env.example order/.env

# webhook payment
$ cp webhook/.env.example webhook/.env
```

Next, you can run the program with `build.sh` file
```bash
# change permission of file build.sh
$ chmod +x build.sh

# if you dont have docker you can running this
$ ./build.sh requirements

# if you have docker you can running this
$ ./build.sh build
```
