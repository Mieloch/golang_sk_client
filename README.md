# golang_sk_client
 - instalacja GO
 
sudo add-apt-repository ppa:ubuntu-lxc/lxd-stable

sudo apt-get update

sudo apt-get install golang

 - stworzenie root path dla GO workspace

export GOPATH=$HOME/{workspace}

export PATH=$PATH:$GOPATH/bin


 - pobranie biblioteki UI wykorzystanej w projekcie

go get github.com/jroimartin/gocui


 - stworzenie folderu z projektem (tutaj git clone)

mkdir $GOPATH/src/github.com/{user}/{project_name}

 - budowanie projektu. Binarka tworzy się w $GOPATH/bin/{project_name}
 
 
cd  $GOPATH/src/github.com/{user}/{project_name}

go install
