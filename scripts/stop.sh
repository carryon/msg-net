ps -ef | grep msg-net | grep -v grep | awk '{print $2}' | xargs kill -9
