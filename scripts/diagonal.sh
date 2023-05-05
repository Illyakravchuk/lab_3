#!/bin/bash

# Вказуємо параметри руху вікна і фігури
startWinPositon=100
finishWinPositon=700
windowSize=800
figureRadius=100
maxDistance=600
jmp=10
timeInterval=0.01
currPos=$startWinPositon
distance=0

# Відпрвляємо запит для створення вікна і фігури
curl -d "figure $startWinPositon $startWinPositon" http://localhost:17000 

sleep $timeInterval
 
while true; do
    # Обноляєм координати фігури
    curl -d "update" http://localhost:17000 

    # Рухаємо фігуру вліво
    while (( currPos < finishWinPositon )); do
        if ((distance >= maxDistance)); then
            curl -d "move $((maxDistance-distance)) $((maxDistance-distance))" http://localhost:17000 
            currPos=$finishWinPositon
            distance=$maxDistance
        else
            curl -d "move $jmp $jmp" http://localhost:17000 
            currPos=$((currPos + jmp))
            distance=$((distance + jmp))
        fi
        curl -d "update" http://localhost:17000 
        sleep $timeInterval
    done

    # Рухаємо фігуру вправо
    while (( currPos > startWinPositon )); do
        if ((distance - jmp <= 0)); then
            curl -d "move $((-distance)) $((-distance))" http://localhost:17000 
            currPos=$startWinPositon
            distance=0
        else
            curl -d "move $((-jmp)) $((-jmp))" http://localhost:17000 
            currPos=$((currPos - jmp))
            distance=$((distance - jmp))
        fi
        curl -d "update" http://localhost:17000 
        sleep $timeInterval
    done
done