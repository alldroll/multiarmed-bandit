FROM golang

RUN mkdir -p /go/src/github.com/alldroll/multiarmed-bandit

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)

# TODO add dependency manager

RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u gonum.org/v1/gonum/stat/distuv

# RUN git clone https://github.com/vishnubob/wait-for-it.git
