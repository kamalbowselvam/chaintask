FROM mcr.microsoft.com/devcontainers/go:0-1-bullseye

# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

# [Optional] Uncomment the next lines to use go get to install anything else you need
# USER vscode
# RUN go get -x <your-dependency-or-tool>
# USER root

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz && mv migrate /usr/bin/migrate
RUN apt update && apt install postgresql-client -y
RUN curl -L https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz | tar xvz && mv ngrok /usr/bin/ngrok
USER vscode
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/golang/mock/mockgen@v1.6.0