# xinga.me
> Api de Ofensas Gratuitas. Inspirado por [Guia do Xingamento das Galáxias](https://blog.bytequeeugosto.com.br/guia-do-xingamento-das-galaxias/)

## Como utilizar

### Web

http://xinga-me.appspot.com/

### Api

http://xinga-me.appspot.com/api

### Slack

Está disponível no slack através do Slash Command no endereço *http://xinga-me.appspot.com/slack*

### Docker

Tenha certeza de ter o [Docker](https://www.docker.com/get-started) instalado na sua máquina.

Para buildar a imagem Docker, rode o seguinte comando:

```shell
$ docker build . -t xinga-me
```

Para executar um container, defina uma variável de ambiente `PORT` e rode o seguinte comando:

```shell
$ export PORT=8000

$ docker run --env "PORT=$PORT" -p $PORT:$PORT xinga-me
```

### Docker compose

Para executar o docker-compose, defina uma variável de ambiente `PORT` e rode o seguinte comando:

```shell
$ export PORT=8000

$ docker-compose up
```

## Referencias

* Obtida do blog post [Byte Que Eu Gosto - Guia do Xingamento das Galáxias](https://blog.bytequeeugosto.com.br/guia-do-xingamento-das-galaxias/)
