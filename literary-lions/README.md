# Literary lions

![moonlit lions](web/static/image/lions2.webp)

## Cloning repository

```bash
git clone https://gitea.koodsisu.fi/jerevaisanen/literary-lions.git [custom name]
```

## Installation with Docker, and running the server

To download and install docker, first visit

https://docs.docker.com/get-started/get-docker/

and follow the instructions for you operating system.

Once you have docker installed, you can build and run the container using
the makefile provided in the repository.

```bash
make all
```

Once the containerization is complete, the server will be up and running,
and can be accessed at

`http://localhost:8080/`

in your browser.

To stop the container, use the command

```bash
docker ps
```

to see all your docker containers, and their name and id. Then, you can
stop the container using the command

```bash
docker stop [id]
```

## Forum usage

In order to interact with the forum, users must be registered. Registered users can
create new posts, like/dislike posts and comments, and leave comments on posts and reply to comments.

To sign in or register, click on the lion icon seen on the right side of the header.
When signed in, users get to their profile using this same button.