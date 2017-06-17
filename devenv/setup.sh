wget -q https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz && sudo tar -zxf go1.8.linux-amd64.tar.gz -C /usr/local && rm -f go1.8.linux-amd64.tar.gz

cat <<EOF >/home/vagrant/golang.sh
export GOROOT="/usr/local/go"
export PATH="/usr/local/go/bin:$PATH"
export GOPATH="/opt/gopath"
EOF

source /home/vagrant/golang.sh

# Ensure permissions are set for GOPATH
sudo chown -R vagrant:vagrant $GOPATH

cd $GOPATH/src/github.com/bocheninc/msg-net
go clean
go install

sudo mv $GOPATH/bin/msg-net /usr/local/bin/

