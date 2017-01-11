# golang_sk_client
// tworzy root dla GO workspace

export GOROOT=$HOME/{go_workspace}

export PATH=$PATH:$GOROOT/bin

// pobiera biblioteke UI potrzebna do zbudowania

go get github.com/jroimartin/gocui

// tworzy folder z projektem

mkdir $GOPATH/src/github.com/{user}/{project_name}

cd  $GOPATH/src/github.com/{user}/{project_name}

// buduje projekt jako  $GOPATH/bin/{project_name}

go install
