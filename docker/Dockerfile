FROM golang

RUN mkdir -p /go/src/github.com/alldroll/multiarmed-bandit

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/go-sql-driver/mysql
# RUN git clone https://github.com/vishnubob/wait-for-it.git
