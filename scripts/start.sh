echo "Input router number:"
read num

echo "Input router listen port"
read port

discovery=0.0.0.0:$port
for((i=0; i<$num; i++));do
    if [ $i -eq 0 ]; then
        echo "./msg-net router --id=$i --address=0.0.0.0:$port &"
        ./msg-net router --id=$i --address=0.0.0.0:$port &
    else
        echo "./msg-net router --id=$i --address=0.0.0.0:$[$port+i] --discovery=$discovery &"
        ./msg-net router --id=$i --address=0.0.0.0:$[$port+i] --discovery=$discovery &
    fi
done    